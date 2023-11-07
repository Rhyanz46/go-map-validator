package map_validator

import "errors"

var (
	ErrNoData        = errors.New("no validatorType")
	ErrInvalidFormat = errors.New("validatorType format invalid")
	ErrUnsupportType = errors.New("type is not support")
)

type loadFromType int

const (
	fromHttpJson loadFromType = iota
	fromHttpMultipartForm
	fromMapString
)
