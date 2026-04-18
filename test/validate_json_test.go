package test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

type validateJSONUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TestValidateJSONHappyPath: JSON valid + rules valid → struct terisi, tanpa error.
func TestValidateJSONHappyPath(t *testing.T) {
	body := `{"email": "dev@example.com", "password": "secret123"}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("email", map_validator.Email().WithMax(255)).
		SetRule("password", map_validator.Str().Between(6, 64)).
		Done()

	got, err := map_validator.ValidateJSON[validateJSONUser](req, rules)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Email != "dev@example.com" {
		t.Errorf("expected email 'dev@example.com', got %q", got.Email)
	}
	if got.Password != "secret123" {
		t.Errorf("expected password 'secret123', got %q", got.Password)
	}
}

// TestValidateJSONInvalidJSON: body JSON rusak → error (sebelum validasi).
func TestValidateJSONInvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{not valid`))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("email", map_validator.Email()).
		Done()

	_, err := map_validator.ValidateJSON[validateJSONUser](req, rules)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// TestValidateJSONMissingRequiredField: field wajib tidak ada → error dengan nama field.
func TestValidateJSONMissingRequiredField(t *testing.T) {
	body := `{"email": "dev@example.com"}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("email", map_validator.Email()).
		SetRule("password", map_validator.Str().WithMin(6)).
		Done()

	_, err := map_validator.ValidateJSON[validateJSONUser](req, rules)
	if err == nil {
		t.Fatal("expected error for missing password, got nil")
	}
	if !strings.Contains(err.Error(), "password") {
		t.Errorf("expected error mentioning 'password', got %q", err.Error())
	}
}

// TestValidateJSONTypeMismatch: tipe field tidak sesuai rules → error.
func TestValidateJSONTypeMismatch(t *testing.T) {
	type Age struct {
		Umur int `json:"umur"`
	}
	body := `{"umur": "bukan angka"}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("umur", map_validator.Int()).
		Done()

	_, err := map_validator.ValidateJSON[Age](req, rules)
	if err == nil {
		t.Fatal("expected error for type mismatch, got nil")
	}
}

// TestValidateJSONEmptyRules: rules kosong → ErrNoRules.
func TestValidateJSONEmptyRules(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")

	empty := map_validator.BuildRoles() // tidak ada SetRule

	_, err := map_validator.ValidateJSON[validateJSONUser](req, empty)
	if err == nil {
		t.Fatal("expected error for empty rules, got nil")
	}
	if err != map_validator.ErrNoRules {
		t.Errorf("expected ErrNoRules, got %v", err)
	}
}

// TestValidateJSONNilRequest: request nil → error (tidak panic).
func TestValidateJSONNilRequest(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("email", map_validator.Email()).
		Done()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected error return, got panic: %v", r)
		}
	}()

	_, err := map_validator.ValidateJSON[validateJSONUser](nil, rules)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
}

// TestValidateJSONNested: nested object ter-validasi dan ter-bind.
func TestValidateJSONNested(t *testing.T) {
	type Address struct {
		City string `json:"city"`
	}
	type User struct {
		Email   string  `json:"email"`
		Address Address `json:"address"`
	}

	addrRules := map_validator.BuildRoles().
		SetRule("city", map_validator.Str().Between(1, 100)).
		Done()
	rules := map_validator.BuildRoles().
		SetRule("email", map_validator.Email()).
		SetRule("address", map_validator.NestedObject(addrRules)).
		Done()

	body := `{"email": "a@b.com", "address": {"city": "Jakarta"}}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	got, err := map_validator.ValidateJSON[User](req, rules)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Email != "a@b.com" {
		t.Errorf("expected email 'a@b.com', got %q", got.Email)
	}
	if got.Address.City != "Jakarta" {
		t.Errorf("expected city 'Jakarta', got %q", got.Address.City)
	}
}

// TestValidateJSONSharedRulesAcrossCalls: regression test — rules di-share
// antar request harus aman setelah fix state mutable.
func TestValidateJSONSharedRulesAcrossCalls(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("email", map_validator.Email()).
		SetRule("password", map_validator.Str().WithMin(6)).
		Done()

	for i := 0; i < 3; i++ {
		body := `{"email": "dev@example.com", "password": "secret123"}`
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		got, err := map_validator.ValidateJSON[validateJSONUser](req, rules)
		if err != nil {
			t.Fatalf("call %d: expected no error, got %v", i, err)
		}
		if got.Email != "dev@example.com" {
			t.Errorf("call %d: expected email 'dev@example.com', got %q", i, got.Email)
		}
	}
}

// TestValidateJSONConcurrent: 20 goroutine pakai rules yang sama secara paralel.
// Proof utama bahwa fix state-mutable aman dari race. Wajib PASS di -race.
func TestValidateJSONConcurrent(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("email", map_validator.Email()).
		SetRule("password", map_validator.Str().WithMin(6)).
		Done()

	const workers = 20
	var wg sync.WaitGroup
	errs := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			body := `{"email": "dev@example.com", "password": "secret123"}`
			req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")

			got, err := map_validator.ValidateJSON[validateJSONUser](req, rules)
			if err != nil {
				errs <- err
				return
			}
			if got.Email != "dev@example.com" || got.Password != "secret123" {
				errs <- fmt.Errorf("unexpected bind result: %+v", got)
			}
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("concurrent call failed: %v", err)
	}
}

// TestValidateJSONEmptyBody: body kosong (EOF) → missing-field error, bukan invalid JSON.
func TestValidateJSONEmptyBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("email", map_validator.Email()).
		Done()

	_, err := map_validator.ValidateJSON[validateJSONUser](req, rules)
	if err == nil {
		t.Fatal("expected error for empty body, got nil")
	}
	if strings.Contains(err.Error(), "valid json") {
		t.Errorf("expected missing-field error, got invalid-json error: %q", err.Error())
	}
	if !strings.Contains(err.Error(), "email") {
		t.Errorf("expected error mentioning 'email', got %q", err.Error())
	}
}

// TestValidateJSONStrictMode: Setting{Strict:true} + body punya unknown key → error.
func TestValidateJSONStrictMode(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("email", map_validator.Email()).
		SetRule("password", map_validator.Str().WithMin(6)).
		SetSetting(map_validator.Setting{Strict: true}).
		Done()

	body := `{"email": "dev@example.com", "password": "secret123", "extra_field": "not allowed"}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	_, err := map_validator.ValidateJSON[validateJSONUser](req, rules)
	if err == nil {
		t.Fatal("expected error for unknown key in strict mode, got nil")
	}
	if !strings.Contains(err.Error(), "extra_field") {
		t.Errorf("expected error mentioning 'extra_field', got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "not allowed") {
		t.Errorf("expected error to say 'not allowed key', got %q", err.Error())
	}
}

// TestValidateJSONDefaultValue: field optional dengan .Nullable().Default("guest")
// saat body tidak kirim field → struct ter-populate nilai default.
func TestValidateJSONDefaultValue(t *testing.T) {
	type userWithRole struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	rules := map_validator.BuildRoles().
		SetRule("email", map_validator.Email()).
		SetRule("role", map_validator.Str().Nullable().Default("guest")).
		Done()

	body := `{"email": "dev@example.com"}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	got, err := map_validator.ValidateJSON[userWithRole](req, rules)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Role != "guest" {
		t.Errorf("expected role default 'guest', got %q", got.Role)
	}
}

// TestValidateJSONListObject: ListOfObject(itemRules), body array of object,
// tiap item ter-validasi dan ter-bind ke struct.
func TestValidateJSONListObject(t *testing.T) {
	type item struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
	}
	type order struct {
		Goods []item `json:"goods"`
	}

	itemRules := map_validator.BuildRoles().
		SetRule("name", map_validator.Str().Between(1, 50)).
		SetRule("quantity", map_validator.Int().WithMin(1)).
		Done()
	rules := map_validator.BuildRoles().
		SetRule("goods", map_validator.ListOfObject(itemRules)).
		Done()

	body := `{"goods": [{"name": "Apple", "quantity": 2}, {"name": "Banana", "quantity": 5}]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	got, err := map_validator.ValidateJSON[order](req, rules)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got.Goods) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got.Goods))
	}
	if got.Goods[0].Name != "Apple" || got.Goods[0].Quantity != 2 {
		t.Errorf("item[0] mismatch: %+v", got.Goods[0])
	}
	if got.Goods[1].Name != "Banana" || got.Goods[1].Quantity != 5 {
		t.Errorf("item[1] mismatch: %+v", got.Goods[1])
	}

	// negative case: item invalid (quantity < 1) → error
	badBody := `{"goods": [{"name": "Apple", "quantity": 0}]}`
	badReq := httptest.NewRequest("POST", "/test", bytes.NewBufferString(badBody))
	badReq.Header.Set("Content-Type", "application/json")
	if _, err := map_validator.ValidateJSON[order](badReq, rules); err == nil {
		t.Error("expected error for invalid item quantity, got nil")
	}
}
