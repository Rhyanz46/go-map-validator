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
	UniqueOrigin      *string
	UniqueTarget      *string
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
	OnUnique      *string
}

func (cm *CustomMsg) uniqueNotNil() bool {
	return cm.OnUnique != nil
}

func (cm *CustomMsg) maxNotNil() bool {
	return cm.OnMax != nil
}

func (cm *CustomMsg) minNotNil() bool {
	return cm.OnMin != nil
}

func (cm *CustomMsg) regexNotNil() bool {
	return cm.OnRegexString != nil
}

func (cm *CustomMsg) typeNotMatchNotNil() bool {
	return cm.OnTypeNotMatch != nil
}

func (cm *CustomMsg) isNotNil() (notNil bool) {
	if cm == nil {
		return
	}
	if cm.OnTypeNotMatch != nil {
		notNil = true
	}
	if cm.OnMax != nil {
		notNil = true
	}
	if cm.OnMin != nil {
		notNil = true
	}
	if cm.OnRegexString != nil {
		notNil = true
	}
	if cm.OnUnique != nil {
		notNil = true
	}
	return notNil
}

type Setting struct {
	Strict bool
}

// rulesWrapper implements RulesWrapper
type rulesWrapper struct {
	Rules           map[string]Rules
	ListRules       ListRules
	isListRules     bool
	Setting         Setting
	uniqueValues    *map[string]map[string]interface{}
	filledField     *[]string
	nullFields      *[]string
	requiredWithout *map[string][]string
	requiredIf      *map[string][]string
	manipulator     []manipulator
}

type ListRules struct {
	Min    *int64
	Max    *int64
	Unique bool
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

	RequiredWithout []string
	RequiredIf      []string
	Object          RulesWrapper
	ListObject      RulesWrapper
	List            ListRulesWrapper

	CustomMsg CustomMsg // will support soon
}

type FileRequest struct {
	File     multipart.File
	FileInfo *multipart.FileHeader
}

type ruleState struct {
	rules              RulesWrapper
	extension          []ExtensionType
	strictAllowedValue bool
}

type dataState struct {
	rules              RulesWrapper
	extension          []ExtensionType
	strictAllowedValue bool
}

type finalOperation struct {
	rules      RulesWrapper
	loadedFrom loadFromType
	extension  []ExtensionType
	data       map[string]interface{}
}

type ExtraOperationData struct {
	rules        RulesWrapper
	loadedFrom   *loadFromType
	data         *map[string]interface{}
	filledFields []string
	nullFields   []string
}
