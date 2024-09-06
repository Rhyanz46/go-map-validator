package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"reflect"
	"testing"
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
			SetManipulator("name", &trimAfterValidation)}).
		SetManipulator("description", &trimAfterValidation).
		Done()

	xx, err := map_validator.NewValidateBuilder().SetRules(roles).Load(data)
	_, err = xx.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
}

//func TestManipulateWithFullBuilderRoles(t *testing.T) {
//	data := map[string]interface{}{
//		"description": "test arian   keren bgt kan \n\n mantap bukan",
//		"note":        "coba aja mungkin      bisa \t mantap",
//	}
//
//	roles := map_validator.BuildRoles().
//		SetRule("description", map_validator.Rules{Type: reflect.String}).
//		SetRule("note", map_validator.Rules{Type: reflect.String}).
//		SetManipulator("description", nil).
//		Done()
//
//	xx, err := map_validator.NewValidateBuilder().SetRules(roles).Load(data)
//	_, err = xx.RunValidate()
//	if err != nil {
//		t.Errorf("Expected not have error, but got error : %s", err)
//		return
//	}
//}
