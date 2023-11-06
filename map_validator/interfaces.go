package map_validator

import "net/http"

type validatorType interface {
	int | string
}

type setRoleOperationType interface {
	SetRules(validations map[string]Rules) *dataState
}

type loadOperationType interface {
	LoadJsonHttp(r *http.Request) (*finalOperation, error)
	LoadFormHttp(r *http.Request) (*finalOperation, error)
	Load(data map[string]interface{}) *finalOperation
}

type finalOperationType interface {
	RunValidate() (*extraOperation, error)
}

type extraOperationType interface {
	Bind(i interface{}) error
}
