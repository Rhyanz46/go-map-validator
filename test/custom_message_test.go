package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"testing"
)

func TestInvalidRegexMessage(t *testing.T) {
	payload := map[string]interface{}{"hp": "62567888", "email": "devariansaputra.com"}
	validRole := map[string]map_validator.Rules{
		"hp": {RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$`, CustomMsg: map_validator.CustomMsg{
			OnRegexString: map_validator.SetMessage("Your ${field} is not valid phone number"),
		}},
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
	validRole = map[string]map_validator.Rules{
		"email": {
			RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
			CustomMsg:   map_validator.CustomMsg{OnRegexString: map_validator.SetMessage("Your email is not valid email format")},
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
