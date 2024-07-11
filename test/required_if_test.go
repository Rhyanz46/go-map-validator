package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"reflect"
	"strings"
	"testing"
)

func TestRequiredIf(t *testing.T) {
	check, err := map_validator.NewValidateBuilder().SetRules(map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"name":          {Type: reflect.String},
			"flavor":        {Type: reflect.String, RequiredWithout: []string{"name"}},
			"custom_flavor": {Type: reflect.String, RequiredIf: []string{"name"}},
		},
	}).Load(map[string]interface{}{
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
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"data": {Object: &map_validator.RulesWrapper{
				Rules: map[string]map_validator.Rules{
					"unit_size": {Type: reflect.Int, RequiredIf: []string{"size"}},
					"size":      {Type: reflect.Int, RequiredIf: []string{"unit_size"}},
				},
			}},
		},
	}
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
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"data": {Object: &map_validator.RulesWrapper{
				Rules: map[string]map_validator.Rules{
					"unit_size": {Type: reflect.Int, RequiredIf: []string{"size"}},
					"size":      {Type: reflect.Int, RequiredIf: []string{"unit_size"}},
				},
			}},
		},
	}
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
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"data": {Object: &map_validator.RulesWrapper{
				Rules: map[string]map_validator.Rules{
					"name":          {Type: reflect.String},
					"flavor":        {Type: reflect.String, RequiredWithout: []string{"custom_flavor", "size"}},
					"custom_flavor": {Type: reflect.String, RequiredWithout: []string{"flavor", "size"}},
					"unit_size":     {Enum: &map_validator.EnumField[any]{Items: []string{"GB", "MB", "TB"}}, RequiredIf: []string{"size"}},
					"size":          {Type: reflect.Int, RequiredWithout: []string{"flavor", "custom_flavor"}, RequiredIf: []string{"unit_size"}},
				},
			}},
		},
	}
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
