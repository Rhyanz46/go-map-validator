package test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

func TestRequiredWithout(t *testing.T) {
	role := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{Type: reflect.String}).
		SetRule("flavor", map_validator.Rules{Type: reflect.String, RequiredWithout: []string{"custom_flavor"}}).
		SetRule("custom_flavor", map_validator.Rules{Type: reflect.String, RequiredWithout: []string{"flavor"}})

	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(map[string]interface{}{
		"name": "SSD",
	})
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}

	expected := "is null you need to put value in "
	_, err = check.RunValidate()
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error with text at least %s, but we got %s", expected, err)
	}
}

func TestChildRequiredWithout(t *testing.T) {
	roleChild := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{Type: reflect.String}).
		SetRule("expired", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("flavor", map_validator.Rules{Type: reflect.String, RequiredWithout: []string{"custom_flavor", "size"}}).
		SetRule("custom_flavor", map_validator.Rules{Type: reflect.String, RequiredWithout: []string{"flavor", "size"}}).
		SetRule("size", map_validator.Rules{Type: reflect.Int, RequiredWithout: []string{"flavor", "custom_flavor"}})

	role := map_validator.
		BuildRoles().
		SetRule("data", map_validator.Rules{Object: roleChild})

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
	expected := "is null you need to put value in "
	_, err = check.RunValidate()
	if err == nil || !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error with text at least %s, but we got %s", expected, err)
	}
}

func TestClearChildRequiredWithout(t *testing.T) {
	roleChild := map_validator.BuildRoles().
		SetRule("name", map_validator.Rules{Type: reflect.String}).
		SetRule("expired", map_validator.Rules{Type: reflect.String, Null: true}).
		SetRule("flavor", map_validator.Rules{Type: reflect.String, RequiredWithout: []string{"custom_flavor", "size"}}).
		SetRule("custom_flavor", map_validator.Rules{Type: reflect.String, RequiredWithout: []string{"flavor", "size"}}).
		SetRule("size", map_validator.Rules{Type: reflect.Int, RequiredWithout: []string{"flavor", "custom_flavor"}})

	role := map_validator.BuildRoles().
		SetRule("data", map_validator.Rules{Object: roleChild})

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
