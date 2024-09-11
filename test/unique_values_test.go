package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"reflect"
	"testing"
)

func TestUniqueValue(t *testing.T) {
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"password":     {Type: reflect.String, Null: true},
			"new_password": {Type: reflect.String, Unique: []string{"password"}, Null: true},
		},
	}
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
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"password":     {Type: reflect.String, Null: true},
			"new_password": {Type: reflect.String, Unique: []string{"password"}, Null: true},
			"data": {Object: &map_validator.RulesWrapper{
				Rules: map[string]map_validator.Rules{
					"dt_password":     {Type: reflect.String, Null: true},
					"dt_new_password": {Type: reflect.String, Unique: []string{"password"}, Null: true},
				},
			}},
		},
	}
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
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"password":     {Type: reflect.String, Null: true},
			"new_password": {Type: reflect.String, Unique: []string{"password"}, Null: true},
		},
	}
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
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"data": {Object: &map_validator.RulesWrapper{
				Rules: map[string]map_validator.Rules{
					"name":         {Type: reflect.String, Null: true},
					"password":     {Type: reflect.String, Null: true},
					"new_password": {Type: reflect.String, Unique: []string{"password"}, Null: true},
				},
			}},
		},
	}
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
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"name":         {Type: reflect.String, Unique: []string{"basic", "password", "new_password"}, Null: true},
			"hoby":         {Type: reflect.String, Unique: []string{"basic"}, Null: true},
			"password":     {Type: reflect.String, Null: true},
			"new_password": {Type: reflect.String, Unique: []string{"password"}, Null: true},
		},
	}
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

func TestChildUniqueValueWithCustomMsg(t *testing.T) {
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"data": {Object: &map_validator.RulesWrapper{
				Rules: map[string]map_validator.Rules{
					"name":     {Type: reflect.String, Null: true},
					"password": {Type: reflect.String, Null: true},
					"new_password": {
						Type: reflect.String, Unique: []string{"password"}, Null: true,
						CustomMsg: map_validator.CustomMsg{
							OnUnique: map_validator.SetMessage("Nilai dari '${unique_origin}' tidak boleh sama dengan nilai '${unique_target}'"),
						},
					},
				},
			}},
		},
	}
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
