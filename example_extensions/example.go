package example_extensions

import "github.com/Rhyanz46/go-map-validator/map_validator"

type ExampleExtension struct {
	rules                map_validator.RulesWrapper
	data                 interface{}
	resetAfterValidation bool
	extraData            *map_validator.ExtraOperationData
}

func (e *ExampleExtension) SetRoles(rules map_validator.RulesWrapper) {
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
