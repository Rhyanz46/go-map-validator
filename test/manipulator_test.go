package test

import (
	"reflect"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

func TestManipulate(t *testing.T) {
	data := map[string]interface{}{
		"name":        "arian\n king   saputra",
		"description": "test arian   keren bgt kan \t\n\n mantap bukan",
		"note":        "coba aja mungkin      bisa \t mantap",
		"kelas": map[string]interface{}{
			"name": "kelas ipa \t:12    ABC",
			"age":  12,
		},
	}

	trimAfterValidation := func(i interface{}) (result interface{}, e error) {
		x := i.(string)
		result = trimAndClean(x)
		return
	}

	roles := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{Type: reflect.String}).
		SetRule("description", map_validator.Rules{Type: reflect.String}).
		SetRule("note", map_validator.Rules{Type: reflect.String}).
		SetRule("kelas", map_validator.Rules{Null: true, Object: map_validator.
			BuildRoles().
			SetRule("name", map_validator.Rules{Type: reflect.String}).
			SetRule("age", map_validator.Rules{Type: reflect.Int}).
			SetManipulator("name", trimAfterValidation)}).
		SetManipulator("description", trimAfterValidation)

	xx, err := map_validator.NewValidateBuilder().SetRules(roles).Load(data)
	extraCheck, err := xx.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	keyRes := extraCheck.GetData()["kelas"].(map[string]interface{})["name"]
	if keyRes != "kelas ipa :12 ABC" {
		t.Errorf("Expected name to be arian king saputra, but got %s", keyRes)
	}
	if extraCheck.GetData()["description"] != "test arian keren bgt kan mantap bukan" {
		t.Errorf("Expected description to be test arian   keren bgt kan mantap bukan, but got %s", data["description"])
	}
}

func TestManipulateWithFullBuilderRoles(t *testing.T) {
	data := map[string]interface{}{
		"description": "test arian   keren bgt kan \n\n mantap bukan",
		"note":        "coba aja mungkin      bisa \t mantap",
	}

	trimAfterValidation := func(i interface{}) (result interface{}, e error) {
		x := i.(string)
		result = trimAndClean(x)
		return
	}

	roles := map_validator.BuildRoles().
		SetRule("description", map_validator.Rules{Type: reflect.String}).
		SetRule("note", map_validator.Rules{Type: reflect.String}).
		SetManipulator("description", trimAfterValidation)

	xx, err := map_validator.NewValidateBuilder().SetRules(roles).Load(data)
	extraCheck, err := xx.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	if extraCheck.GetData()["description"] != "test arian keren bgt kan mantap bukan" {
		t.Errorf("Expected description to be test arian   keren bgt kan mantap bukan, but got %s", data["description"])
	}
}

func TestManipulateOnNullableField(t *testing.T) {
	data := map[string]interface{}{
		"note": "coba aja mungkin      bisa \t mantap",
	}

	trimAfterValidation := func(i interface{}) (result interface{}, e error) {
		x := i.(string)
		result = trimAndClean(x)
		return
	}

	roles := map_validator.BuildRoles().
		SetRule("description", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("note", map_validator.Rules{Type: reflect.String}).
		SetManipulator("description", trimAfterValidation)

	xx, err := map_validator.NewValidateBuilder().SetRules(roles).Load(data)
	extraCheck, err := xx.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	if extraCheck.GetData()["description"] != nil {
		t.Errorf("Expected description to be nil, but got %s", data["description"])
	}
}

func TestManipulateOnMultiFields(t *testing.T) {
	data := map[string]interface{}{
		"name":        "arian\n king   saputra",
		"description": "test arian   keren bgt kan \t\n\n mantap bukan",
		"note":        "coba aja mungkin      bisa \t mantap",
	}

	trimAfterValidation := func(i interface{}) (result interface{}, e error) {
		x := i.(string)
		result = trimAndClean(x)
		return
	}

	roles := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{Type: reflect.String}).
		SetRule("description", map_validator.Rules{Type: reflect.String}).
		SetRule("note", map_validator.Rules{Type: reflect.String}).
		SetFieldsManipulator([]string{"description", "name"}, trimAfterValidation)

	xx, err := map_validator.NewValidateBuilder().SetRules(roles).Load(data)
	extraCheck, err := xx.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	if extraCheck.GetData()["name"] != "arian king saputra" {
		t.Errorf("Expected name to be arian king saputra, but got %s", extraCheck.GetData()["name"])
	}
	if extraCheck.GetData()["description"] != "test arian keren bgt kan mantap bukan" {
		t.Errorf("Expected description to be test arian   keren bgt kan mantap bukan, but got %s", extraCheck.GetData()["description"])
	}
}
