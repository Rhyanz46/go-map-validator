package map_validator

import (
	"mime/multipart"
	"reflect"
)

type MessageMeta struct {
	Field        *string
	ExpectedType *reflect.Kind
	ActualType   *reflect.Kind
}

type EnumField[T any] struct {
	Items               T
	StringCaseSensitive bool // will support soon
}

type CustomMsg struct {
	OnTypeNotMatch *string
	//OnEnumValueNotMatch *string
	//OnNull              *string
	//OnMax               *string
	//OnMin               *string
	OnRegexString *string
}

type Setting struct {
	Strict bool
}

type RulesWrapper struct {
	Rules   map[string]Rules
	Setting Setting
}

type Rules struct {
	Null               bool
	NilIfNull          bool
	IsMapInterface     bool
	Email              bool
	Enum               *EnumField[any] // new 🔥🔥🔥
	Type               reflect.Kind
	Max                *int
	Min                *int
	IfNull             interface{}
	UUID               bool
	UUIDToString       bool
	IPV4               bool
	IPV4Network        bool
	IPv4OptionalPrefix bool
	File               bool
	RegexString        string
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
