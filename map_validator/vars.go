package map_validator

import "errors"

var ErrNoData = errors.New("no validatorType")
var ErrInvalidFormat = errors.New("validatorType format invalid")
var ErrUnsupportType = errors.New("type is not support")
