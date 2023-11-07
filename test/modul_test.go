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
	validRole := map[string]map_validator.Rules{
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"hoby":          {Type: reflect.String, Null: false},
		"menikah":       {Type: reflect.Bool, Null: false},
	}
	validRoleOptionalMenikah := map[string]map_validator.Rules{
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"hoby":          {Type: reflect.String, Null: false},
		"menikah":       {Type: reflect.Bool, Null: true},
	}
	payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS bro", "umur": 1, "menikah": true}
	notFullPayload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS bro", "umur": 1}
	check := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
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
		t.Errorf("Expected have an error, but you got no error : %s", err)
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
		t.Errorf("Expected have an error, but you got no error : %s", err)
	} else {
		expected := "the field 'hoby' should be 'int'"
		if err.Error() != expected {
			t.Errorf("Expected :%s. But you got : %s", expected, err)
		}
	}

	check = map_validator.NewValidateBuilder().SetRules(validRole).Load(notFullPayload)
	extraCheck, err = check.RunValidate()
	expected := "we need 'menikah' field"
	if err.Error() != expected {
		t.Errorf("Expected : %s But you got : %s", expected, err)
	}

	testBind = &Data{}
	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
	}
	err = extraCheck.Bind(testBind)
	expected = "no data to Bind because last progress is error"
	if err != nil {
		if err.Error() != expected {
			t.Errorf("Expected : %s But you got : %s", expected, err)
		}
	}

	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
	}

	check = map_validator.NewValidateBuilder().SetRules(validRoleOptionalMenikah).Load(notFullPayload)
	extraCheck, err = check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}

	testBind = &Data{}
	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
	}
	err = extraCheck.Bind(testBind)
	expected = "no data to Bind because last progress is error"
	if err != nil {
		if err.Error() != expected {
			t.Errorf("Expected : %s But you got : %s", expected, err)
		}
	}

	if testBind.JK != notFullPayload["jenis_kelamin"] {
		t.Errorf("Expected : %s But you got : %s", notFullPayload["jenis_kelamin"], testBind.JK)
	}
}

func TestInterfaceFieldBinding(t *testing.T) {
	payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS bro", "umur": 1, "menikah": true}
	type Data struct {
		JK      string  `map_validator:"jenis_kelamin"`
		Hoby    *string `map_validator:"hoby"`
		Menikah bool    `map_validator:"menikah"`
	}
	validRole := map[string]map_validator.Rules{
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"hoby":          {Type: reflect.String, Null: true},
		"menikah":       {Type: reflect.Bool, Null: false},
	}

	check := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	extraCheck, err := check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}

	testBind := &Data{}
	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %v", testBind.JK)
	}
	err = extraCheck.Bind(testBind)
	if err != nil {
		t.Errorf("Expected have an error, but you got no error : %s", err)
	}

	if testBind.JK != payload["jenis_kelamin"] {
		t.Errorf("Expected : %s But you got : %s", payload["jenis_kelamin"], testBind.JK)
	}

	if testBind.JK != payload["jenis_kelamin"] {
		t.Errorf("Expected : %s But you got : %s", payload["jenis_kelamin"], testBind.JK)
	}
}
