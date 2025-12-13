package map_validator

import "errors"

var (
	ErrNoData            = errors.New("no validatorType")
	ErrInvalidJsonFormat = errors.New("is not valid json")
	ErrUnsupportType     = errors.New("type is not support")
)

type LoadFromType int

const (
	FromHttpJson LoadFromType = iota
	FromHttpMultipartForm
	FromMapString
	FromJSONEncoder
)

// Keep backward compatibility with internal names
type loadFromType = LoadFromType

const (
	fromHttpJson          = FromHttpJson
	fromHttpMultipartForm = FromHttpMultipartForm
	fromMapString         = FromMapString
	fromJSONEncoder       = FromJSONEncoder
)

const (
	chainKey = "root"
)
