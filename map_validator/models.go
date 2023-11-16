package map_validator

import (
	"mime/multipart"
	"reflect"
)

type EnumField[T any] struct {
	Items               T
	StringCaseSensitive bool // will support soon
}

type CustomMsg struct {
	OnTypeNotMatch      *string
	OnEnumValueNotMatch *string
	OnNull              *string
	OnMax               *string
	OnMin               *string
}

type Rules struct {
	Null               bool
	NilIfNull          bool
	IsMapInterface     bool
	Email              bool
	Enum               *EnumField[any]
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

	CustomMsg CustomMsg // will support soon
}

type FileRequest struct {
	File     multipart.File
	FileInfo *multipart.FileHeader
}

type ruleState struct {
}

type dataState struct {
	rules    *map[string]Rules
	optional *optionalRolesState
}

type oneDataState struct {
	rule *Rules
}

type optionalRolesState struct {
	data               *dataState
	strictAllowedValue bool
}

type finalOperation struct {
	rules      *map[string]Rules
	rule       *Rules
	loadedFrom loadFromType
	oneValue   bool
	data       map[string]interface{}
}

type extraOperation struct {
	rules        *map[string]Rules
	rule         *Rules
	oneValue     bool
	loadedFrom   *loadFromType
	data         *map[string]interface{}
	filledFields []string
	nullFields   []string
}
