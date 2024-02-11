package map_validator

import "net/http"

type validatorType interface {
	int | string
}

type setRoleOperationType interface {
	SetRules(validations map[string]Rules) *dataState
	StrictKeys() *ruleState
	AddExtension(extension ExtensionType) *ruleState
}

type loadOperationType interface {
	LoadJsonHttp(r *http.Request) (*finalOperation, error)
	LoadFormHttp(r *http.Request) (*finalOperation, error)
	Load(data map[string]interface{}) (*finalOperation, error)
}

type finalOperationType interface {
	RunValidate() (*ExtraOperationData, error)
}

type ExtraOperationType interface {
	Bind(i interface{}) error
	GetFilledField() []string
	GetNullField() []string
	GetData() map[string]interface{}
}

type ExtensionType interface {
	SetRoles(rules *map[string]Rules)
	BeforeLoad(data interface{}) error
	AfterLoad(data *map[string]interface{}) error
	BeforeValidation(data *map[string]interface{}) error
	AfterValidation(data *map[string]interface{}) error
	SetExtraData(data *ExtraOperationData) ExtensionType
}
