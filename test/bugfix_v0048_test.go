package test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

// v0.0.48 regression tests — each targets a defect confirmed by probe before
// the fix.

// --- 1. IfNull/Default values are validated like any other value ---

func TestDefaultGoesThroughValidation(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("status", map_validator.StrEnum("a", "b").Default("zzz")).
		Done()
	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, err = op.RunValidate(); err == nil {
		t.Error("invalid enum default 'zzz' should fail validation")
	}

	op, err = map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{"status": "a"})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, err = op.RunValidate(); err != nil {
		t.Errorf("explicit valid value should pass, got: %v", err)
	}

	intRules := map_validator.BuildRoles().
		SetRule("n", map_validator.Int().Default("hello")).
		Done()
	op, err = map_validator.NewValidateBuilder().SetRules(intRules).Load(map[string]interface{}{})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, err = op.RunValidate(); err == nil {
		t.Error("non-integer default 'hello' should fail validation")
	}

	okRules := map_validator.BuildRoles().
		SetRule("status", map_validator.StrEnum("a", "b").Default("a")).
		Done()
	op, err = map_validator.NewValidateBuilder().SetRules(okRules).Load(map[string]interface{}{})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	extra, err := op.RunValidate()
	if err != nil {
		t.Fatalf("valid default should pass, got: %v", err)
	}
	if extra.GetData()["status"] != "a" {
		t.Errorf("default should be applied, got %v", extra.GetData()["status"])
	}
}

// --- 2. Misconfigured Enum (non-slice / nil Items) is an error, not a silent pass ---

func TestMisconfiguredEnumRejected(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("x", map_validator.Rules{Type: reflect.String, Enum: &map_validator.EnumField[any]{Items: "abc"}}).
		Done()
	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{"x": "totally-invalid"})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, err = op.RunValidate(); err == nil {
		t.Error("enum with non-slice Items must not silently accept any value")
	}

	// nil Items previously panicked on reflect.TypeOf(nil).Kind()
	nilRules := map_validator.BuildRoles().
		SetRule("x", map_validator.Rules{Type: reflect.String, Enum: &map_validator.EnumField[any]{}}).
		Done()
	op, err = map_validator.NewValidateBuilder().SetRules(nilRules).Load(map[string]interface{}{"x": "anything"})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, err = op.RunValidate(); err == nil {
		t.Error("enum with nil Items should error, not panic")
	}
}

// --- 3. Conditional required works even when the dependency has no rule ---

func TestConditionalRequiredWithUndeclaredDependency(t *testing.T) {
	// RequiredWithout: flavor is required when custom_flavor is absent.
	// custom_flavor is intentionally NOT declared as a rule.
	rules := map_validator.BuildRoles().
		SetRule("flavor", map_validator.Rules{Type: reflect.String, RequiredWithout: []string{"custom_flavor"}}).
		Done()
	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	_, err = op.RunValidate()
	if err == nil {
		t.Error("flavor should be required when its undeclared dependency is absent")
	}
	if err != nil && !strings.Contains(err.Error(), "custom_flavor") {
		t.Errorf("error should mention the dependency, got: %v", err)
	}

	// dependency present in the payload (still undeclared) → flavor optional
	op, err = map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{"custom_flavor": "vanilla"})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, err = op.RunValidate(); err != nil {
		t.Errorf("flavor should be optional when dependency is present, got: %v", err)
	}

	// RequiredIf: c is required when a is filled. 'a' is not a rule.
	ifRules := map_validator.BuildRoles().
		SetRule("c", map_validator.Rules{Type: reflect.String, RequiredIf: []string{"a"}}).
		Done()
	op, err = map_validator.NewValidateBuilder().SetRules(ifRules).Load(map[string]interface{}{"a": "x"})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	_, err = op.RunValidate()
	if err == nil {
		t.Error("c should be required when its undeclared dependency is filled")
	}

	op, err = map_validator.NewValidateBuilder().SetRules(ifRules).Load(map[string]interface{}{"a": "x", "c": "y"})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, err = op.RunValidate(); err != nil {
		t.Errorf("c present → no error expected, got: %v", err)
	}

	// declared-dependency behavior is preserved
	declared := map_validator.BuildRoles().
		SetRule("a", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("c", map_validator.Rules{Type: reflect.String, RequiredIf: []string{"a"}}).
		Done()
	op, _ = map_validator.NewValidateBuilder().SetRules(declared).Load(map[string]interface{}{"a": "x"})
	if _, err = op.RunValidate(); err == nil {
		t.Error("declared dependency filled + c absent should still error")
	}
	op, _ = map_validator.NewValidateBuilder().SetRules(declared).Load(map[string]interface{}{})
	if _, err = op.RunValidate(); err != nil {
		t.Errorf("declared dependency absent → c optional, got: %v", err)
	}
}

// --- 4. JSON integers beyond float64's exact range survive to Bind ---

func TestJSONBigNumberPrecision(t *testing.T) {
	type dto struct {
		ID int64 `json:"id"`
	}
	rules := map_validator.BuildRoles().SetRule("id", map_validator.Int64()).Done()

	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"id": 9007199254740993}`))
	op, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	extra, err := op.RunValidate()
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}
	var got dto
	if err = extra.Bind(&got); err != nil {
		t.Fatalf("bind error: %v", err)
	}
	if got.ID != 9007199254740993 {
		t.Errorf("precision corrupted: want 9007199254740993, got %d", got.ID)
	}

	req = httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"id": 9223372036854775807}`))
	op, err = map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	extra, err = op.RunValidate()
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}
	got = dto{}
	if err = extra.Bind(&got); err != nil {
		t.Fatalf("bind error: %v", err)
	}
	if got.ID != 9223372036854775807 {
		t.Errorf("want math.MaxInt64, got %d", got.ID)
	}
}

// --- 5. IPV4Network works as a standalone flag rule (no Type needed) ---

func TestIPV4NetworkStandaloneRule(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("net", map_validator.Rules{IPV4Network: true}).
		Done()
	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{"net": "10.0.0.0"})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, err = op.RunValidate(); err != nil {
		t.Errorf("IPV4Network without explicit Type should accept a network address, got: %v", err)
	}

	op, err = map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{"net": "10.0.0.1"})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	_, err = op.RunValidate()
	if err == nil {
		t.Error("10.0.0.1 is not a network address (.0) and must be rejected")
	}
}

// --- 6. IPv4-mapped IPv6 strings are not accepted as IPv4 ---

func TestIPv4RejectsMappedIPv6(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("ip", map_validator.IPv4()).
		SetRule("net", map_validator.Rules{IPV4Network: true}).
		Done()
	payload := map[string]interface{}{"ip": "::ffff:1.2.3.4", "net": "::ffff:1.2.3.0"}
	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(payload)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, err = op.RunValidate(); err == nil {
		t.Error("IPv4-mapped IPv6 strings must be rejected")
	}

	op, _ = map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{"ip": "192.168.1.1", "net": "192.168.1.0"})
	if _, err = op.RunValidate(); err != nil {
		t.Errorf("valid dotted-quad IPv4/network must still pass, got: %v", err)
	}
}

// --- 7. Min/Max constrain the item count of ListObject rules ---

func TestListObjectContainerMinMax(t *testing.T) {
	item := map_validator.BuildRoles().SetRule("n", map_validator.Int()).Done()
	rules := map_validator.BuildRoles().
		SetRule("goods", map_validator.Rules{ListObject: item, Min: map_validator.SetTotal(3), Max: map_validator.SetTotal(4)}).
		Done()

	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{
		"goods": []interface{}{map[string]interface{}{"n": 1}},
	})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	_, err = op.RunValidate()
	if err == nil {
		t.Error("1 item should fail Min=3")
	}

	op, _ = map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{
		"goods": []interface{}{
			map[string]interface{}{"n": 1},
			map[string]interface{}{"n": 2},
			map[string]interface{}{"n": 3},
		},
	})
	if _, err = op.RunValidate(); err != nil {
		t.Errorf("3 items within [3,4] should pass, got: %v", err)
	}
}

// --- 8. List(NestedObject(w)) validates items like ListOfObject ---

func TestListWithNestedObjectRules(t *testing.T) {
	item := map_validator.BuildRoles().
		SetRule("name", map_validator.Str()).
		SetRule("quantity", map_validator.Int().WithMin(1)).
		Done()
	rules := map_validator.BuildRoles().
		SetRule("goods", map_validator.List(map_validator.NestedObject(item)).WithMin(2)).
		Done()

	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{
		"goods": []interface{}{map[string]interface{}{"name": "Apple", "quantity": 2}},
	})
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	_, err = op.RunValidate()
	if err == nil {
		t.Error("1 item should fail WithMin(2)")
	}

	op, _ = map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{
		"goods": []interface{}{
			map[string]interface{}{"name": "Apple", "quantity": 2},
			map[string]interface{}{"name": "Pear", "quantity": 0},
		},
	})
	_, err = op.RunValidate()
	if err == nil {
		t.Error("item with quantity=0 must be rejected by the item rule")
	}

	op, _ = map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{
		"goods": []interface{}{
			map[string]interface{}{"name": "Apple", "quantity": 2},
			map[string]interface{}{"name": "Pear", "quantity": 5},
		},
	})
	if _, err = op.RunValidate(); err != nil {
		t.Errorf("valid list should pass, got: %v", err)
	}
}

// --- 9. Min/Max apply to format rules (UUID, enum, regex), not only Email ---

func TestFormatRuleMinMaxEnforced(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("u", map_validator.UUID().WithMax(10)).
		SetRule("s", map_validator.StrEnum("abcdef").WithMax(3)).
		SetRule("r", map_validator.Str().Regex(`^\d+$`).WithMax(2)).
		SetRule("e", map_validator.Email().WithMax(100)).
		SetRule("ip", map_validator.IPv4().WithMax(15)).
		Done()
	payload := map[string]interface{}{
		"u":  "3b241101-e2bb-4255-8caf-4136c566a962",
		"s":  "abcdef",
		"r":  "123",
		"e":  "a@b.c",
		"ip": "192.168.1.1",
	}
	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(payload)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	_, err = op.RunValidate()
	if err == nil {
		t.Fatal("expected a Min/Max violation for u, s or r")
	}
	msg := err.Error()
	if !strings.Contains(msg, "u") && !strings.Contains(msg, "s") && !strings.Contains(msg, "r") {
		t.Errorf("error should come from the over-length format fields, got: %v", msg)
	}

	payload["u"] = "3b241101-e2bb-4255-8caf-4136c566a962"
	payload["s"] = "abcdef"
	payload["r"] = "12"
	okRules := map_validator.BuildRoles().
		SetRule("u", map_validator.UUID().WithMax(40)).
		SetRule("s", map_validator.StrEnum("abc", "abcdef").WithMax(6)).
		SetRule("r", map_validator.Str().Regex(`^\d+$`).WithMax(2)).
		SetRule("e", map_validator.Email().WithMax(100)).
		SetRule("ip", map_validator.IPv4().WithMax(15)).
		Done()
	op, _ = map_validator.NewValidateBuilder().SetRules(okRules).Load(payload)
	if _, err = op.RunValidate(); err != nil {
		t.Errorf("within-bound values should pass, got: %v", err)
	}
}

// --- 10. Bind keeps multipart files (no JSON round-trip on FileRequest) ---

func TestBindPreservesMultipartFile(t *testing.T) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile("doc", "a.txt")
	fw.Write([]byte("hello"))
	w.WriteField("name", "x")
	w.Close()
	req := httptest.NewRequest("POST", "/", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	rules := map_validator.BuildRoles().
		SetRule("doc", map_validator.Rules{File: true, Null: true}).
		SetRule("name", map_validator.Str()).
		Done()
	op, err := map_validator.NewValidateBuilder().SetRules(rules).LoadFormHttp(req)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	extra, err := op.RunValidate()
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}

	type dto struct {
		Doc  map_validator.FileRequest `json:"doc"`
		Name string                    `json:"name"`
	}
	var got dto
	if err = extra.Bind(&got); err != nil {
		t.Fatalf("bind should work with file fields, got: %v", err)
	}
	if got.Doc.File == nil {
		t.Fatal("Doc.File should be bound")
	}
	if got.Doc.FileInfo == nil || got.Doc.FileInfo.Filename != "a.txt" {
		t.Errorf("FileInfo.Filename should be a.txt, got: %v", got.Doc.FileInfo)
	}
	content, err := io.ReadAll(got.Doc.File)
	if err != nil {
		t.Fatalf("read uploaded file: %v", err)
	}
	if string(content) != "hello" {
		t.Errorf("file content should be 'hello', got %q", string(content))
	}
	got.Doc.File.Close()
	if got.Name != "x" {
		t.Errorf("Name should bind as usual, got %q", got.Name)
	}
}

// --- 11. Fractional JSON numbers are rejected by integer rules at validation ---

func TestFractionalIntRejectedFromJson(t *testing.T) {
	rules := map_validator.BuildRoles().SetRule("qty", map_validator.Int()).Done()

	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"qty": 2.5}`))
	op, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	_, err = op.RunValidate()
	if err == nil {
		t.Error("2.5 must fail an Int rule at validation, not at Bind")
	}

	req = httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"qty": 2.0}`))
	op, err = map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, err = op.RunValidate(); err != nil {
		t.Errorf("integral 2.0 should pass an Int rule, got: %v", err)
	}
}

// --- 12. nil payload / null body produce validation errors, not "no data" ---

func TestLoadNilAndNullBodyAreEmptyPayloads(t *testing.T) {
	rules := map_validator.BuildRoles().SetRule("n", map_validator.Int()).Done()

	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(nil)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	_, err = op.RunValidate()
	if err == nil {
		t.Fatal("nil payload should still require field n")
	}
	if strings.Contains(err.Error(), "no data to Validate") {
		t.Errorf("misleading 'no data' error for nil payload: %v", err)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(`null`))
	op, err = map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	_, err = op.RunValidate()
	if err == nil {
		t.Fatal("null body should still require field n")
	}
	if strings.Contains(err.Error(), "no data to Validate") {
		t.Errorf("misleading 'no data' error for null body: %v", err)
	}
}

// --- 13. The first reported error is deterministic (sorted rule keys) ---

func TestDeterministicFirstError(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("aaa", map_validator.Int()).
		SetRule("bbb", map_validator.Email()).
		Done()
	payload := map[string]interface{}{"aaa": "not-a-number", "bbb": "not-an-email"}

	var first string
	for i := 0; i < 20; i++ {
		op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(payload)
		if err != nil {
			t.Fatalf("load error: %v", err)
		}
		_, err = op.RunValidate()
		if err != nil {
			first = err.Error()
		}
		if first != "the field 'aaa' should be 'int'" {
			t.Fatalf("run %d: expected deterministic 'aaa' error, got: %v", i, first)
		}
	}
}
