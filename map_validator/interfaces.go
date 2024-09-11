package map_validator

import "net/http"

type validatorType interface {
	int | string
}

type setRoleOperationType interface {
	SetRules(validations RulesWrapper) *dataState
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
	SetRoles(rules RulesWrapper)
	BeforeLoad(data interface{}) error
	AfterLoad(data *map[string]interface{}) error
	BeforeValidation(data *map[string]interface{}) error
	AfterValidation(data *map[string]interface{}) error
	SetExtraData(data *ExtraOperationData) ExtensionType
}

type ChainResultType interface {
	GetAllKeys() []string
	PrintHierarchyWithSeparator(separator string, currentPath string)
	ToMap() map[string]interface{}
	RunManipulator() error
	RunUniqueChecker()
	GetErrors() []error
}

type ChainerType interface {
	GetParentKey() string
	Next(index int) ChainerType
	Back() ChainerType
	Forward(index int) ChainerType
	SetKey(name string) ChainerType
	GetKey() string
	SetKeyValue(key string, value interface{}) ChainerType
	GetParentKeys() []string
	AddChild() ChainerType
	LoadFromMap(data map[string]interface{})
	SetValue(value interface{}) ChainerType
	GetValue() interface{}
	SetManipulator(manipulator *func(interface{}) (interface{}, error)) ChainerType
	SetUniques(uniques []string) ChainerType
	SetCustomMsg(customMsg *CustomMsg) ChainerType
	GetUniques() []string
	AddError(err error) ChainerType
	GetResult() ChainResultType

	GetChildren() []ChainerType
	GetParent() ChainerType
	GetBrothers() []ChainerType
}
