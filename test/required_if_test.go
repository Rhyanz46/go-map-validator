package test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

func TestRequiredIf(t *testing.T) {
	role := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{Type: reflect.String}).
		SetRule("flavor", map_validator.Rules{Type: reflect.String, RequiredWithout: []string{"name"}}).
		SetRule("custom_flavor", map_validator.Rules{Type: reflect.String, RequiredIf: []string{"name"}})
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(map[string]interface{}{
		"name": "SSD",
		//"custom_flavor": "large",
	})
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}

	expected := "' is filled you need to put value in ["
	_, err = check.RunValidate()
	if err == nil || !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error contains '%s', but we got %s", expected, err)
	}
}

func TestChildRequiredIfStandardError(t *testing.T) {
	roleChild := map_validator.BuildRoles().
		SetRule("unit_size", map_validator.Rules{Type: reflect.Int, RequiredIf: []string{"size"}}).
		SetRule("size", map_validator.Rules{Type: reflect.Int, RequiredIf: []string{"unit_size"}})

	role := map_validator.BuildRoles().
		SetRule("data", map_validator.Rules{Object: roleChild})
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"size": 1,
			//"unit_size": 2,
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	expected := "' is filled you need to put value in ["
	_, err = check.RunValidate()
	if err == nil || !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error contains '%s', but we got %s", expected, err)
	}
}

func TestChildRequiredIfStandardOk(t *testing.T) {
	roleChild := map_validator.BuildRoles().
		SetRule("unit_size", map_validator.Rules{Type: reflect.Int, RequiredIf: []string{"size"}}).
		SetRule("size", map_validator.Rules{Type: reflect.Int, RequiredIf: []string{"unit_size"}})

	role := map_validator.BuildRoles().SetRule("data", map_validator.Rules{Object: roleChild})
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"size":      1,
			"unit_size": 2,
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	_, err = check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
}

func TestChildRequiredIf(t *testing.T) {
	childRole := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{Type: reflect.String}).
		SetRule("flavor", map_validator.Rules{Type: reflect.String, RequiredWithout: []string{"custom_flavor", "size"}}).
		SetRule("custom_flavor", map_validator.Rules{Type: reflect.String, RequiredWithout: []string{"flavor", "size"}}).
		SetRule("unit_size", map_validator.Rules{Enum: &map_validator.EnumField[any]{Items: []string{"GB", "MB", "TB"}}, RequiredIf: []string{"size"}}).
		SetRule("size", map_validator.Rules{Type: reflect.Int, RequiredWithout: []string{"flavor", "custom_flavor"}, RequiredIf: []string{"unit_size"}})

	role := map_validator.BuildRoles().
		SetRule("data", map_validator.Rules{Object: childRole})
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"name":   "sabalong",
			"flavor": "normal",
			"size":   1,
			//"unit_size": 2,
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	expected := "' is filled you need to put value in ["
	_, err = check.RunValidate()

	if err == nil || !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error contains '%s', but we got %s", expected, err)
	}
}
