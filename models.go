package mapValidator

import "reflect"

type RequestDataValidator struct {
	Null           bool
	NilIfNull      bool
	IsMapInterface bool
	//Enum           []interface{} willl support soon
	Type               reflect.Kind
	Max                *int
	Min                *int
	IfNull             interface{}
	UUID               bool
	UUIDToString       bool
	IPV4               bool
	IPv4OptionalPrefix bool
}
