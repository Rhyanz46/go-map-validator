package map_validator

import "errors"

var (
	ErrNoData            = errors.New("no validatorType")
	ErrInvalidJsonFormat = errors.New("is not valid json")
	ErrUnsupportType     = errors.New("type is not support")
)

type loadFromType int

const (
	fromHttpJson loadFromType = iota
	fromHttpMultipartForm
	fromMapString
)
