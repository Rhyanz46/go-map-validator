package test

import (
	"github.com/Rhyanz46/go-map-validator/example_extensions"
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"testing"
)

func TestExtension(t *testing.T) {
	payload := map[string]interface{}{"hp": "+62567888", "email": "dev@ariansaputra.com"}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"hp":    {RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$`},
			"email": {RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`},
		},
	}
	swaggerExt := example_extensions.ManipulatorExt().ResetAfterValidation()

	check, err := map_validator.
		NewValidateBuilder().
		AddExtension(swaggerExt).
		SetRules(validRole).
		Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	res, err := check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	if len(res.GetData()) > 0 {
		t.Errorf("Expected no data, but we got : %s", res.GetData())
	}
}

func TestInvalidExtension(t *testing.T) {
	payload := map[string]interface{}{"hp": "+62567888", "email": "dev@ariansaputra.com"}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"hp":    {RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$`},
			"email": {RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`},
		},
	}
	swaggerExt := example_extensions.ManipulatorExt()

	check, err := map_validator.
		NewValidateBuilder().
		AddExtension(swaggerExt).
		SetRules(validRole).
		Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	res, err := check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	if payload["hp"] != res.GetData()["hp"] {
		t.Errorf("Expected %s data, but we got : %s", payload["hp"], res.GetData()["hp"])
	}
}
