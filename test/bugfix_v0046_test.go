package test

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

// Bug 1: LoadFormHttp rejected flag-based rules (Email, UUID, IPv4, Regex)
// with ErrUnsupportType because their Type is reflect.Invalid.
func TestFormFlagBasedRules(t *testing.T) {
	form := url.Values{}
	form.Set("email", "dev@example.com")
	form.Set("nickname", "dev_01")
	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rules := map_validator.BuildRoles().
		SetRule("email", map_validator.Email()).
		SetRule("nickname", map_validator.Str().Regex("^[a-z0-9_]+$")).
		Done()

	op, err := map_validator.NewValidateBuilder().SetRules(rules).LoadFormHttp(req)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	if _, err := op.RunValidate(); err != nil {
		t.Fatalf("validate: %s", err)
	}

	// invalid email must still be rejected through the same path
	form2 := url.Values{}
	form2.Set("email", "not-an-email")
	form2.Set("nickname", "dev_01")
	req2 := httptest.NewRequest("POST", "/", strings.NewReader(form2.Encode()))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	op2, err := map_validator.NewValidateBuilder().SetRules(rules).LoadFormHttp(req2)
	if err != nil {
		t.Fatalf("load2: %s", err)
	}
	if _, err := op2.RunValidate(); err == nil {
		t.Fatal("expected invalid email to fail validation")
	}

	// NestedObject and List rules must be rejected early with ErrUnsupportType
	nestedRules := map_validator.BuildRoles().
		SetRule("nested", map_validator.NestedObject(map_validator.BuildRoles().SetRule("a", map_validator.Str()).Done())).
		Done()
	if _, err := map_validator.NewValidateBuilder().SetRules(nestedRules).LoadFormHttp(req); err != map_validator.ErrUnsupportType {
		t.Fatalf("expected ErrUnsupportType for form+NestedObject, got %v", err)
	}

	listRules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str())).
		Done()
	if _, err := map_validator.NewValidateBuilder().SetRules(listRules).LoadFormHttp(req); err != map_validator.ErrUnsupportType {
		t.Fatalf("expected ErrUnsupportType for form+List, got %v", err)
	}
}

// Bug 2: WithMax on float fields truncated the value (3.9 -> 3) so values in
// (max, max+1) wrongly passed.
func TestFloatMaxNoTruncation(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("price", map_validator.Float64().WithMax(3)).
		Done()

	op, err := map_validator.NewValidateBuilder().SetRules(rules).
		Load(map[string]interface{}{"price": 3.9})
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	if _, err := op.RunValidate(); err == nil {
		t.Fatal("expected 3.9 to fail WithMax(3)")
	}

	// boundary values must still pass
	rules2 := map_validator.BuildRoles().
		SetRule("price", map_validator.Float64().Between(1, 3)).
		Done()
	op2, err := map_validator.NewValidateBuilder().SetRules(rules2).
		Load(map[string]interface{}{"price": 3.0})
	if err != nil {
		t.Fatalf("load2: %s", err)
	}
	if _, err := op2.RunValidate(); err != nil {
		t.Fatalf("expected 3.0 to pass Between(1, 3): %s", err)
	}
}

// Bug 3: manipulators declared on ListObject item rules never ran.
func TestListObjectItemManipulator(t *testing.T) {
	upper := func(data interface{}) (interface{}, error) {
		return strings.ToUpper(data.(string)), nil
	}
	itemRules := map_validator.BuildRoles().
		SetRule("name", map_validator.Str()).
		SetManipulator("name", upper).
		Done()
	rules := map_validator.BuildRoles().
		SetRule("items", map_validator.ListOfObject(itemRules)).
		Done()

	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{
		"items": []interface{}{map[string]interface{}{"name": "abc"}},
	})
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	extra, err := op.RunValidate()
	if err != nil {
		t.Fatalf("validate: %s", err)
	}
	items := extra.GetData()["items"].([]interface{})
	first := items[0].(map[string]interface{})
	if first["name"] != "ABC" {
		t.Fatalf("expected manipulator to produce ABC, got %v", first["name"])
	}
}

// Bug 4: Unique across sibling fields inside ListObject items was never checked.
func TestListObjectItemUnique(t *testing.T) {
	itemRules := map_validator.BuildRoles().
		SetRule("password", map_validator.Str().UniqueFrom("confirm")).
		SetRule("confirm", map_validator.Str()).
		Done()
	rules := map_validator.BuildRoles().
		SetRule("items", map_validator.ListOfObject(itemRules)).
		Done()

	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{
		"items": []interface{}{map[string]interface{}{"password": "same", "confirm": "same"}},
	})
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	if _, err := op.RunValidate(); err == nil {
		t.Fatal("expected unique violation inside ListObject item to fail")
	}
}

// Bug 5: ListRules.Unique was accepted but never enforced.
func TestListRulesUniqueEnforced(t *testing.T) {
	rule := map_validator.List(map_validator.Str())
	// emulate BuildListRoles().SetListRule(ListRules{Unique: true})
	rule.List.SetListRule(struct {
		Min    *int64
		Max    *int64
		Unique bool
	}{Unique: true})
	rules := map_validator.BuildRoles().SetRule("tags", rule).Done()

	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{
		"tags": []interface{}{"a", "b", "a"},
	})
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	if _, err := op.RunValidate(); err == nil {
		t.Fatal("expected duplicate elements to fail with ListRules.Unique")
	}

	// distinct elements must still pass
	op2, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{
		"tags": []interface{}{"a", "b", "c"},
	})
	if err != nil {
		t.Fatalf("load2: %s", err)
	}
	if _, err := op2.RunValidate(); err != nil {
		t.Fatalf("expected distinct elements to pass: %s", err)
	}
}

// Bug 6: a rule carrying both Object and ListObject panicked on type assertion.
func TestObjectPlusListObjectNoPanic(t *testing.T) {
	inner := map_validator.BuildRoles().SetRule("x", map_validator.Str()).Done()
	rules := map_validator.BuildRoles().SetRule("obj", map_validator.Rules{
		Object:     inner,
		ListObject: inner,
	}).Done()

	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{
		"obj": []interface{}{map[string]interface{}{"x": "a"}},
	})
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked: %v", r)
		}
	}()
	// an error is acceptable; a panic is not
	_, _ = op.RunValidate()
}
