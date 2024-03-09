package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"reflect"
	"strings"
	"testing"
)

func TestRequiredWithout(t *testing.T) {
	check, err := map_validator.NewValidateBuilder().SetRules(map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"name":          {Type: reflect.String},
			"flavor":        {Type: reflect.String, RequiredWithout: []string{"custom_flavor"}},
			"custom_flavor": {Type: reflect.String, RequiredWithout: []string{"flavor"}},
		},
	}).Load(map[string]interface{}{
		"name": "SSD",
	})
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}

	expected := "is null you need to put value in this"
	_, err = check.RunValidate()
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error with text at least %s, but we got %s", expected, err)
	}
}

func TestChildRequiredWithout(t *testing.T) {
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"data": {Object: &map_validator.RulesWrapper{
				Rules: map[string]map_validator.Rules{
					"name":          {Type: reflect.String},
					"expired":       {Type: reflect.String, Null: true},
					"flavor":        {Type: reflect.String, RequiredWithout: []string{"custom_flavor", "size"}},
					"custom_flavor": {Type: reflect.String, RequiredWithout: []string{"flavor", "size"}},
					"size":          {Type: reflect.Int, RequiredWithout: []string{"flavor", "custom_flavor"}},
				},
			}},
		},
	}
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"name": "sabalong",
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	expected := "is null you need to put value in this"
	_, err = check.RunValidate()
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error with text at least %s, but we got %s", expected, err)
	}
}

func TestClearChildRequiredWithout(t *testing.T) {
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"data": {Object: &map_validator.RulesWrapper{
				Rules: map[string]map_validator.Rules{
					"name":          {Type: reflect.String},
					"expired":       {Type: reflect.String, Null: true},
					"flavor":        {Type: reflect.String, RequiredWithout: []string{"custom_flavor", "size"}},
					"custom_flavor": {Type: reflect.String, RequiredWithout: []string{"flavor", "size"}},
					"size":          {Type: reflect.Int, RequiredWithout: []string{"flavor", "custom_flavor"}},
				},
			}},
		},
	}
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"name": "sabalong",
			"size": 63,
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
