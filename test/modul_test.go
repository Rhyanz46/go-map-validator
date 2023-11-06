package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"reflect"
	"testing"
)

func TestMultipleValidation(t *testing.T) {
	payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS", "umur": 1, "menikah": true}
	err := map_validator.NewValidateBuilder().SetRules(map[string]map_validator.Rules{
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"hoby":          {Type: reflect.String, Null: false},
		"menikah":       {Type: reflect.Bool, Null: false},
	}).Load(payload).RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}

	err = map_validator.NewValidateBuilder().SetRules(map[string]map_validator.Rules{
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"hoby":          {Type: reflect.Int, Null: false},
		"menikah":       {Type: reflect.Bool, Null: false},
	}).Load(payload).RunValidate()
	if err == nil {
		t.Error("Expected have an error, but you got no error")
	} else {
		expected := "the field 'hoby' should be 'int'"
		if err.Error() != expected {
			t.Errorf("Expected :%s. But you got : %s", expected, err)
		}
	}
}
