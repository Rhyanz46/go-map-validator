package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"reflect"
	"testing"
)

func TestMultipleValidation(t *testing.T) {
	type Data struct {
		JK      string `map_validator:"jenis_kelamin"`
		Hoby    string `map_validator:"hoby"`
		Menikah bool   `map_validator:"menikah"`
	}
	payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS bro", "umur": 1, "menikah": true}
	check := map_validator.NewValidateBuilder().SetRules(map[string]map_validator.Rules{
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"hoby":          {Type: reflect.String, Null: false},
		"menikah":       {Type: reflect.Bool, Null: false},
	}).Load(payload)
	extraCheck, err := check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}

	testBind := &Data{}
	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
	}
	err = extraCheck.Bind(testBind)
	if err != nil {
		return
	}

	if testBind.JK != payload["jenis_kelamin"] {
		t.Errorf("Expected : %s But you got : %s", payload["jenis_kelamin"], testBind.JK)
	}

	check = map_validator.NewValidateBuilder().SetRules(map[string]map_validator.Rules{
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"hoby":          {Type: reflect.Int, Null: false},
		"menikah":       {Type: reflect.Bool, Null: false},
	}).Load(payload)
	_, err = check.RunValidate()
	if err == nil {
		t.Error("Expected have an error, but you got no error")
	} else {
		expected := "the field 'hoby' should be 'int'"
		if err.Error() != expected {
			t.Errorf("Expected :%s. But you got : %s", expected, err)
		}
	}
}
