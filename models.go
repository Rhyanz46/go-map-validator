package mapValidator

import "reflect"

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

type RequestDataValidator struct {
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
	IPv4OptionalPrefix bool

	CustomMsg CustomMsg // will support soon
}
