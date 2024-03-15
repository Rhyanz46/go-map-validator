package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"reflect"
	"testing"
)

func TestInvalidRegexMessage(t *testing.T) {
	payload := map[string]interface{}{"hp": "62567888", "email": "devariansaputra.com"}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"hp": {RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$`, CustomMsg: map_validator.CustomMsg{
				OnRegexString: map_validator.SetMessage("Your ${field} is not valid phone number"),
			}},
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err == nil {
		t.Error("Expected error, but got no error :")
	}
	expected := "Your hp is not valid phone number"
	if err.Error() != expected {
		t.Errorf("Expected '%s', but we got '%s' :", expected, err.Error())
	}
	validRole = map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"email": {
				RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
				CustomMsg:   map_validator.CustomMsg{OnRegexString: map_validator.SetMessage("Your email is not valid email format")},
			},
		},
	}
	expected = "Your email is not valid email format"
	check, err = map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err == nil {
		t.Error("Expected error, but got no error :")
	}
	if err.Error() != expected {
		t.Errorf("Expected '%s', but we got '%s' :", expected, err.Error())
	}
}

func TestValidRegexMessage(t *testing.T) {
	payload := map[string]interface{}{"hp": "+62567888", "email": "dev@ariansaputra.com"}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"hp": {RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$`, CustomMsg: map_validator.CustomMsg{
				OnRegexString: map_validator.SetMessage("Your ${field} is not valid phone number"),
			}},
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err != nil {
		t.Errorf("Expected no error, but we got error : %s ", err.Error())
	}
}

func TestInvalidTypeNotMatchMessage(t *testing.T) {
	payload := map[string]interface{}{"total": "2", "unit": "KG"}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"total": {
				Type: reflect.Int64,
				CustomMsg: map_validator.CustomMsg{
					OnTypeNotMatch: map_validator.SetMessage("Total must be a number, but your input is ${actual_type}"),
				},
			},
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err == nil {
		t.Error("Expected error, but got no error :")
	}
	expected := "Total must be a number, but your input is string"
	if err.Error() != expected {
		t.Errorf("Expected '%s', but we got '%s' :", expected, err.Error())
	}
}

func TestValidTypeNotMatchMessage(t *testing.T) {
	payload := map[string]interface{}{"total": 12, "unit": "KG"}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"total": {
				Type: reflect.Int,
				CustomMsg: map_validator.CustomMsg{
					OnTypeNotMatch: map_validator.SetMessage("Total must be a number, but your input is ${actual_type}"),
				},
			},
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
}

func TestInvalidLengthMessage(t *testing.T) {
	payload := map[string]interface{}{"total": 1, "unit": "KG"}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"total": {
				Type: reflect.Int,
				Max:  map_validator.SetTotal(3),
				Min:  map_validator.SetTotal(2),
				CustomMsg: map_validator.CustomMsg{
					OnMin: map_validator.SetMessage("The min size allowed is ${expected_min_length}., but your input is ${actual_length}"),
					OnMax: map_validator.SetMessage("The max size allowed is ${expected_max_length}., but your input is ${actual_length}"),
				},
			},
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	expected := "The min size allowed is 2., but your input is 1"
	_, err = check.RunValidate()
	if err.Error() != expected {
		t.Errorf("Expected %s, but got error : %s", expected, err)
	}

	payload = map[string]interface{}{"total": 12, "unit": "KG"}
	check, err = map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	expected = "The max size allowed is 3., but your input is 12"
	_, err = check.RunValidate()
	if err.Error() != expected {
		t.Errorf("Expected %s, but got error : %s", expected, err)
	}
}

func TestInvalidLengthMessageCaseOne(t *testing.T) {
	validRole := map_validator.NewValidateBuilder().SetRules(map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"project_id": {UUID: true},
			"flavor":     {UUID: true, RequiredWithout: []string{"custom_flavor"}},
			"custom_flavor": {RequiredWithout: []string{"flavor"}, Object: &map_validator.RulesWrapper{Rules: map[string]map_validator.Rules{
				"size": {
					Type: reflect.Float64,
					Min:  map_validator.SetTotal(25),
					Max:  map_validator.SetTotal(9999999),
					CustomMsg: map_validator.CustomMsg{
						OnMin: map_validator.SetMessage("The minimum size allowed is 25 GB."),
					},
				},
			}}},
			"name":           {Type: reflect.String, Max: map_validator.SetTotal(255)},
			"description":    {Type: reflect.String, Max: map_validator.SetTotal(500), Null: true},
			"network_subnet": {IPV4Network: true, RequiredWithout: []string{"vpc"}},
			"vpc":            {Type: reflect.String, Max: map_validator.SetTotal(255), RequiredWithout: []string{"network_subnet"}},
			"zone":           {Type: reflect.String, Max: map_validator.SetTotal(255)},
			"path":           {Type: reflect.String, Max: map_validator.SetTotal(500)},
		},
		Setting: map_validator.Setting{Strict: true},
	})
	check, err := validRole.Load(map[string]interface{}{
		"custom_flavor": map[string]interface{}{"size": 99999991},
		"zone":          "arjuna",
		"path":          "/aaa/aad",
		"name":          "121",
		"project_id":    "647d4b1b-b36b-4c6e-85b5-f4cabc7a8a78",
		"vpc":           "12.1.2.3",
	})
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	expected := "the field 'size' should be or lower than 9999999"
	_, err = check.RunValidate()
	if err.Error() != expected {
		t.Errorf("Expected %s, but got error : %s", expected, err)
	}
}
