package test

import (
	"reflect"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

func TestInvalidRegexMessage(t *testing.T) {
	payload := map[string]interface{}{"hp": "62567888", "email": "devariansaputra.com"}
	validRole := map_validator.
		BuildRoles().
		SetRule("hp", map_validator.Rules{
			RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$`, CustomMsg: map_validator.CustomMsg{
				OnRegexString: map_validator.SetMessage("Your ${field} is not valid phone number"),
			}},
		)
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
	validRole = map_validator.BuildRoles().SetRule("email", map_validator.Rules{
		RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		CustomMsg:   map_validator.CustomMsg{OnRegexString: map_validator.SetMessage("Your email is not valid email format")},
	})
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
	validRole := map_validator.
		BuildRoles().
		SetRule("hp", map_validator.Rules{
			RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$`, CustomMsg: map_validator.CustomMsg{
				OnRegexString: map_validator.SetMessage("Your ${field} is not valid phone number"),
			}})

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
	validRole := map_validator.BuildRoles().SetRule("total", map_validator.Rules{
		Type: reflect.Int64,
		CustomMsg: map_validator.CustomMsg{
			OnTypeNotMatch: map_validator.SetMessage("Total must be a number, but your input is ${actual_type}"),
		},
	})
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
	validRole := map_validator.
		BuildRoles().
		SetRule("total", map_validator.Rules{
			Type: reflect.Int,
			CustomMsg: map_validator.CustomMsg{
				OnTypeNotMatch: map_validator.SetMessage("Total must be a number, but your input is ${actual_type}"),
			},
		})
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
	validRole := map_validator.BuildRoles().
		SetRule("total", map_validator.Rules{
			Type: reflect.Int,
			Max:  map_validator.SetTotal(3),
			Min:  map_validator.SetTotal(2),
			CustomMsg: map_validator.CustomMsg{
				OnMin: map_validator.SetMessage("The min size allowed is ${expected_min_length}., but your input is ${actual_length}"),
				OnMax: map_validator.SetMessage("The max size allowed is ${expected_max_length}., but your input is ${actual_length}"),
			},
		})
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

func TestCustomEnumMessage(t *testing.T) {
	payload := map[string]interface{}{"status": "invalid_status"}
	validRole := map_validator.BuildRoles().SetRule("status", map_validator.Rules{
		Type: reflect.String,
		Enum: &map_validator.EnumField[any]{Items: []string{"active", "inactive", "pending"}},
		CustomMsg: map_validator.CustomMsg{
			OnEnumValueNotMatch: map_validator.SetMessage("Field '${field}' expected ${expected_type}, got ${actual_type} - must be one of the allowed values"),
		},
	})
	
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err == nil {
		t.Error("Expected error, but got no error")
	}
	expected := "Field 'status' expected string, got string - must be one of the allowed values"
	if err.Error() != expected {
		t.Errorf("Expected '%s', but we got '%s'", expected, err.Error())
	}
}

func TestCustomEnumMessageWithHttpRequest(t *testing.T) {
	// Test with HTTP JSON (integer enum) - demonstrates type coercion from float64 to int
	payload := map[string]interface{}{"port": 999}
	validRole := map_validator.BuildRoles().SetRule("port", map_validator.Rules{
		Type: reflect.Int,
		Enum: &map_validator.EnumField[any]{Items: []int{80, 443, 8080}},
		CustomMsg: map_validator.CustomMsg{
			OnEnumValueNotMatch: map_validator.SetMessage("Field '${field}' expected ${expected_type} but got ${actual_type} - port not allowed"),
		},
	})
	
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err == nil {
		t.Error("Expected error, but got no error")
	}
	expected := "Field 'port' expected int but got int - port not allowed"
	if err.Error() != expected {
		t.Errorf("Expected '%s', but we got '%s'", expected, err.Error())
	}
}

func TestCustomEnumMessageTypeCoercion(t *testing.T) {
	// Test to demonstrate type coercion with HTTP JSON - float64 input to int enum
	// This shows when expected_type != actual_type
	payload := map[string]interface{}{"priority": 5.0} // float64 from JSON
	validRole := map_validator.BuildRoles().SetRule("priority", map_validator.Rules{
		Type: reflect.Int,
		Enum: &map_validator.EnumField[any]{Items: []int{1, 2, 3}}, // int enum
		CustomMsg: map_validator.CustomMsg{
			OnEnumValueNotMatch: map_validator.SetMessage("Priority '${field}': expected ${expected_type}, received ${actual_type}"),
		},
	})
	
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err == nil {
		t.Error("Expected error, but got no error")
	}
	expected := "Priority 'priority': expected int, received float64"
	if err.Error() != expected {
		t.Errorf("Expected '%s', but we got '%s'", expected, err.Error())
	}
}

func TestInvalidLengthMessageCaseOneT(t *testing.T) {
	rolesss := map_validator.BuildRoles().
		SetRule("project_id", map_validator.Rules{UUID: true}).
		SetRule("flavor", map_validator.Rules{UUID: true, RequiredWithout: []string{"custom_flavor"}}).
		SetRule("custom_flavor", map_validator.Rules{
			RequiredWithout: []string{"flavor"},
			Object: map_validator.BuildRoles().SetRule("size", map_validator.Rules{
				Type: reflect.Float64,
				Min:  map_validator.SetTotal(25),
				Max:  map_validator.SetTotal(9999999),
				CustomMsg: map_validator.CustomMsg{
					OnMin: map_validator.SetMessage("The minimum size allowed is 25 GB."),
				},
			}),
		}).
		SetRule("name", map_validator.Rules{Type: reflect.String, Max: map_validator.SetTotal(255)}).
		SetRule("description", map_validator.Rules{Type: reflect.String, Max: map_validator.SetTotal(500), Null: true}).
		SetRule("network_subnet", map_validator.Rules{IPV4Network: true, RequiredWithout: []string{"vpc"}}).
		SetRule("vpc", map_validator.Rules{Type: reflect.String, Max: map_validator.SetTotal(255), RequiredWithout: []string{"network_subnet"}}).
		SetRule("zone", map_validator.Rules{Type: reflect.String, Max: map_validator.SetTotal(255)}).
		SetRule("path", map_validator.Rules{Type: reflect.String, Max: map_validator.SetTotal(500)}).
		SetSetting(*&map_validator.Setting{Strict: true})
	validRole := map_validator.NewValidateBuilder().SetRules(rolesss)
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
