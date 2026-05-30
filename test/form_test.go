package test

import (
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

func TestForm_IntCoerced(t *testing.T) {
	form := url.Values{}
	form.Set("age", "30")
	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	type out struct {
		Age int `json:"age"`
	}
	rules := map_validator.BuildRoles().
		SetRule("age", map_validator.Rules{Type: reflect.Int})

	op, err := map_validator.NewValidateBuilder().SetRules(rules).LoadFormHttp(req)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	extra, err := op.RunValidate()
	if err != nil {
		t.Fatalf("expected int form value to validate, got: %s", err)
	}
	var o out
	if err := extra.Bind(&o); err != nil {
		t.Fatalf("bind: %s", err)
	}
	if o.Age != 30 {
		t.Errorf("expected age 30, got %d", o.Age)
	}
}

func TestForm_IntWithMaxAndNegative(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("n", map_validator.Int().Between(0, 100))

	form := url.Values{}
	form.Set("n", "-5")
	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	op, err := map_validator.NewValidateBuilder().SetRules(rules).LoadFormHttp(req)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	if _, err = op.RunValidate(); err == nil {
		t.Error("expected -5 to fail Between(0,100) from form, but passed")
	}
}

func TestForm_BoolCoerced(t *testing.T) {
	form := url.Values{}
	form.Set("active", "true")
	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	type out struct {
		Active bool `json:"active"`
	}
	rules := map_validator.BuildRoles().
		SetRule("active", map_validator.Rules{Type: reflect.Bool})

	op, err := map_validator.NewValidateBuilder().SetRules(rules).LoadFormHttp(req)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	extra, err := op.RunValidate()
	if err != nil {
		t.Fatalf("expected bool form value to validate, got: %s", err)
	}
	var o out
	if err := extra.Bind(&o); err != nil {
		t.Fatalf("bind: %s", err)
	}
	if !o.Active {
		t.Errorf("expected active true, got %v", o.Active)
	}
}

func TestForm_FloatCoerced(t *testing.T) {
	form := url.Values{}
	form.Set("rate", "1.5")
	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rules := map_validator.BuildRoles().
		SetRule("rate", map_validator.Rules{Type: reflect.Float64})

	op, err := map_validator.NewValidateBuilder().SetRules(rules).LoadFormHttp(req)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	if _, err = op.RunValidate(); err != nil {
		t.Fatalf("expected float form value to validate, got: %s", err)
	}
}

func TestForm_InvalidIntFails(t *testing.T) {
	form := url.Values{}
	form.Set("age", "abc")
	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rules := map_validator.BuildRoles().
		SetRule("age", map_validator.Rules{Type: reflect.Int})

	op, err := map_validator.NewValidateBuilder().SetRules(rules).LoadFormHttp(req)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	if _, err = op.RunValidate(); err == nil {
		t.Error("expected 'abc' to fail int validation, but passed")
	}
}

func TestForm_StringStillWorks(t *testing.T) {
	form := url.Values{}
	form.Set("name", "bob")
	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	type out struct {
		Name string `json:"name"`
	}
	rules := map_validator.BuildRoles().
		SetRule("name", map_validator.Str().WithMax(50))

	op, err := map_validator.NewValidateBuilder().SetRules(rules).LoadFormHttp(req)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	extra, err := op.RunValidate()
	if err != nil {
		t.Fatalf("expected string form value to validate, got: %s", err)
	}
	var o out
	if err := extra.Bind(&o); err != nil {
		t.Fatalf("bind: %s", err)
	}
	if o.Name != "bob" {
		t.Errorf("expected name bob, got %q", o.Name)
	}
}
