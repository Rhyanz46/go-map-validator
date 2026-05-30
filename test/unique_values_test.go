package test

import (
	"reflect"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

func TestUniqueValue(t *testing.T) {
	role := map_validator.BuildRoles().
		SetRule("password", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("new_password", map_validator.Rules{Type: reflect.String, Unique: []string{"password"}, Null: true})

	payload := map[string]interface{}{
		"password":     "sabalong",
		"new_password": "sabalong",
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	expected := "value of 'password' and 'new_password' fields must be different"
	expectedOr := "value of 'new_password' and 'password' fields must be different"
	_, err = check.RunValidate()
	if err == nil {
		t.Error("Expected error, but got no error :")
	}
	if err != nil && !(err.Error() == expected || err.Error() == expectedOr) {
		t.Errorf("Expected :%s. But you got : %s", expected, err)
	}
}

func TestUniqueValueInNested(t *testing.T) {
	roleChild := map_validator.BuildRoles().
		SetRule("dt_password", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("dt_new_password", map_validator.Rules{Type: reflect.String, Unique: []string{"password"}, Null: true})
	role := map_validator.BuildRoles().
		SetRule("password", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("new_password", map_validator.Rules{Type: reflect.String, Unique: []string{"password"}, Null: true}).
		SetRule("data", map_validator.Rules{Object: roleChild})

	payload := map[string]interface{}{
		"password":     "sabalong",
		"new_password": "sabalong",
		"data": map[string]interface{}{
			"dt_password":     "golang@123",
			"dt_new_password": "golang@123",
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	expected := "value of 'password' and 'new_password' fields must be different"
	expectedOr := "value of 'new_password' and 'password' fields must be different"
	_, err = check.RunValidate()
	if err == nil {
		t.Error("Expected error, but got no error :")
	}
	if err != nil && !(err.Error() == expected || err.Error() == expectedOr) {
		t.Errorf("Expected :%s. But you got : %s", expected, err)
	}
}

func TestNonUniqueValue(t *testing.T) {
	role := map_validator.BuildRoles().
		SetRule("password", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("new_password", map_validator.Rules{Type: reflect.String, Unique: []string{"password"}, Null: true})

	payload := map[string]interface{}{
		"password":     "sabalong",
		"new_password": "sabalong",
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	expected := "value of 'new_password' and 'password' fields must be different"
	_, err = check.RunValidate()
	if !(err.Error() == expected) {
		t.Errorf("Expected :%s. But you got : %s", expected, err)
	}
}

func TestChildUniqueValue(t *testing.T) {
	roleChild := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("password", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("new_password", map_validator.Rules{Type: reflect.String, Unique: []string{"password"}, Null: true})

	role := map_validator.BuildRoles().
		SetRule("data", map_validator.Rules{Object: roleChild})
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"password":     "sabalong",
			"new_password": "sabalong",
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	expected := "value of 'password' and 'new_password' fields must be different"
	expectedOr := "value of 'new_password' and 'password' fields must be different"
	_, err = check.RunValidate()
	if !(err.Error() == expected || err.Error() == expectedOr) {
		t.Errorf("Expected :%s. But you got : %s", expected, err)
	}
}

func TestUniqueManyValue(t *testing.T) {
	role := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{Type: reflect.String, Unique: []string{"basic", "password", "new_password"}, Null: true}).
		SetRule("hoby", map_validator.Rules{Type: reflect.String, Unique: []string{"basic"}, Null: true}).
		SetRule("password", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("new_password", map_validator.Rules{Type: reflect.String, Unique: []string{"password"}, Null: true})
	payload := map[string]interface{}{
		"name":         "sabalong_samalewa",
		"hoby":         "hoby",
		"password":     "sabalong",
		"new_password": "sabalong_samalewa",
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	expected := "value of 'name' and 'new_password' fields must be different"
	expectedOr := "value of 'new_password' and 'name' fields must be different"
	_, err = check.RunValidate()
	if !(err.Error() == expected || err.Error() == expectedOr) {
		t.Errorf("Expected :%s. But you got : %s", expected, err)
	}
}

// TestUniqueNonComparableValues: unique check must not panic when field values
// are non-comparable (slice/map); equal values should still be flagged.
func TestUniqueNonComparableValues(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("a", map_validator.Any()).
		SetRule("b", map_validator.Any().UniqueFrom("a"))

	check, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{
		"a": []interface{}{1, 2, 3},
		"b": []interface{}{1, 2, 3},
	})
	if err != nil {
		t.Fatalf("load error: %s", err)
	}
	_, err = check.RunValidate()
	if err == nil {
		t.Fatal("expected unique violation for equal slices, got nil")
	}

	// Different slices must pass.
	check2, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{
		"a": []interface{}{1, 2, 3},
		"b": []interface{}{4, 5, 6},
	})
	if err != nil {
		t.Fatalf("load error: %s", err)
	}
	if _, err = check2.RunValidate(); err != nil {
		t.Fatalf("expected different slices to pass, got: %s", err)
	}
}

func TestChildUniqueValueWithCustomMsg(t *testing.T) {
	roleChild := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("password", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("new_password", map_validator.Rules{
			Type: reflect.String, Unique: []string{"password"}, Null: true,
			CustomMsg: map_validator.CustomMsg{
				OnUnique: map_validator.SetMessage("Nilai dari '${unique_origin}' tidak boleh sama dengan nilai '${unique_target}'"),
			},
		})
	role := map_validator.BuildRoles().
		SetRule("data", map_validator.Rules{Object: roleChild})
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"password":     "sabalong",
			"new_password": "sabalong",
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	expected := "Nilai dari 'new_password' tidak boleh sama dengan nilai 'password'"
	_, err = check.RunValidate()
	if err == nil {
		t.Error("Expected have an error, but got no error ")
		return
	}
	if err.Error() != expected {
		t.Errorf("Expected :%s. But you got : %s", expected, err)
	}
}
