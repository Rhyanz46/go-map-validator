package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"testing"
)

type ExampleExtension struct {
	rules                *map[string]map_validator.Rules
	data                 interface{}
	resetAfterValidation bool
	extraData            *map_validator.ExtraOperationData
}

func (e *ExampleExtension) SetRoles(rules *map[string]map_validator.Rules) {
	e.rules = rules
}

func (e *ExampleExtension) BeforeLoad(data interface{}) error {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (e *ExampleExtension) AfterLoad(data *map[string]interface{}) error {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (e *ExampleExtension) BeforeValidation(data *map[string]interface{}) error {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (e *ExampleExtension) AfterValidation(data *map[string]interface{}) error {
	//TODO implement me
	//panic("implement me")
	if e.resetAfterValidation {
		empty := map[string]interface{}{}
		*data = empty
	}
	return nil
}

func (e *ExampleExtension) SetExtraData(data *map_validator.ExtraOperationData) map_validator.ExtensionType {
	e.extraData = data
	return e
}

func (e *ExampleExtension) ResetAfterValidation() *ExampleExtension {
	e.resetAfterValidation = true
	return e
}

func ManipulatorExt() *ExampleExtension {
	return &ExampleExtension{}
}

func TestExtension(t *testing.T) {
	payload := map[string]interface{}{"hp": "+62567888", "email": "dev@ariansaputra.com"}
	validRole := map[string]map_validator.Rules{
		"hp":    {RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$`},
		"email": {RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`},
	}
	swaggerExt := ManipulatorExt().ResetAfterValidation()

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
	validRole := map[string]map_validator.Rules{
		"hp":    {RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$`},
		"email": {RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`},
	}
	swaggerExt := ManipulatorExt()

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
