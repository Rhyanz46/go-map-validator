package map_validator

import "net/http"

type validatorType interface {
	int | string
}

type setRoleOperationType interface {
	SetRules(validations map[string]Rules) *optionalRolesState
	SetRule(validations Rules) *dataState
}

type loadOperationType interface {
	LoadJsonHttp(r *http.Request) (*finalOperation, error)
	LoadMultiPartFormHttp(r *http.Request, resWriter http.ResponseWriter) (*finalOperation, error)
	Load(data map[string]interface{}) (*finalOperation, error)
}

type loadOneValueOperationType interface {
	Load(data interface{}) (*finalOperation, error)
	LoadFromHttp(r *http.Request, resWriter http.ResponseWriter) (*finalOperation, error)
}

type optionalRulesOperationType interface {
	StrictKeys() *optionalRolesState
	Next() *dataState
}

type finalOperationType interface {
	RunValidate() (*extraOperation, error)
}

type extraOperationType interface {
	Bind(i interface{}) error
	GetFilledField() []string
	GetNullField() []string
	GetData() map[string]interface{}
}
