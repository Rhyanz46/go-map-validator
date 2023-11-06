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

	CustomMsg CustomMsg // will support soon
}

type FileRequest struct {
	File     multipart.File
	FileInfo *multipart.FileHeader
}

type ruleState struct {
	rules map[string]Rules
}

type dataState struct {
	*ruleState
	data map[string]interface{}
}

type finalOperation struct {
	*dataState
}

type extraOperation struct {
	*finalOperation
}
