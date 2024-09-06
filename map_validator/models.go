package map_validator

import (
	"mime/multipart"
	"reflect"
)

type manipulator struct {
	Field string
	Func  *func(interface{}) (interface{}, error)
}

type MessageMeta struct {
	Field             *string
	ExpectedType      *reflect.Kind
	ActualLength      *int64
	ExpectedMinLength *int64
	ExpectedMaxLength *int64
	ActualType        *reflect.Kind
}

type EnumField[T any] struct {
	Items               T
	StringCaseSensitive bool // will support soon
}

type CustomMsg struct {
	OnTypeNotMatch *string
	//OnEnumValueNotMatch *string
	//OnNull              *string
	OnMax         *string
	OnMin         *string
	OnRegexString *string
}

type Setting struct {
	Strict bool
}

type RulesWrapper struct {
	Rules           map[string]Rules
	Setting         Setting
	uniqueValues    *map[string]map[string]interface{}
	filledField     *[]string
	nullFields      *[]string
	requiredWithout *map[string][]string
	requiredIf      *map[string][]string
	manipulator     []manipulator
}

type Rules struct {
	Null               bool
	NilIfNull          bool
	AnonymousObject    bool
	Email              bool
	Enum               *EnumField[any]
	Type               reflect.Kind
	Max                *int64
	Min                *int64
	IfNull             interface{}
	UUID               bool
	UUIDToString       bool
	IPV4               bool
	IPV4Network        bool
	IPv4OptionalPrefix bool
	File               bool
	RegexString        string
	Unique             []string
	RequiredWithout    []string
	RequiredIf         []string
	Object             *RulesWrapper
	ListObject         *RulesWrapper

	CustomMsg CustomMsg // will support soon
}

type FileRequest struct {
	File     multipart.File
	FileInfo *multipart.FileHeader
}

type ruleState struct {
	rules              *RulesWrapper
	extension          []ExtensionType
	strictAllowedValue bool
}

type dataState struct {
	rules              *RulesWrapper
	extension          []ExtensionType
	strictAllowedValue bool
}

type finalOperation struct {
	rules      *RulesWrapper
	loadedFrom loadFromType
	extension  []ExtensionType
	data       map[string]interface{}
}

type ExtraOperationData struct {
	rules        *RulesWrapper
	loadedFrom   *loadFromType
	data         *map[string]interface{}
	filledFields []string
	nullFields   []string
}
