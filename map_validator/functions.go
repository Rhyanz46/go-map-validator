package map_validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func isEqualString(current, allowedField string) bool {
	return current == allowedField
}

func isEqualFloat64(current, allowedField float64) bool {
	return current == allowedField
}

func isEqualInt(current, allowedField int) bool {
	return current == allowedField
}

func isEqualInt64(current, allowedField int64) bool {
	return current == allowedField
}

func isEmail(email string) bool {
	ok := strings.Contains(email, "@")
	if !ok {
		return false
	}
	ok = strings.Contains(strings.Split(email, "@")[1], ".")
	return ok
}

func valueInList[T any](listData []T, data T, compare func(T, T) bool) bool {
	for _, currentValue := range listData {
		if compare(currentValue, data) {
			return true
		}
	}
	return false
}

func isIPv4Valid(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

func isIPv4NetworkValid(ip string) bool {
	if isIPv4Valid(ip) {
		ipString := strings.Split(ip, ".")
		if ipString[3] == "0" {
			return true
		}
	}
	return false
}

func buildMessage(msg string, meta MessageMeta) error {
	field := "${field}"
	if strings.Contains(msg, field) {
		msg = strings.ReplaceAll(msg, field, meta.Field)
	}
	return errors.New(msg)
}

func validate(field string, dataTemp map[string]interface{}, validator Rules, dataFrom loadFromType) (interface{}, error) {
	//var oldIntType reflect.Kind
	data := dataTemp[field]

	// null validation
	if !validator.Null && data == nil {
		return nil, errors.New("we need '" + field + "' field")
	} else if validator.Null && data == nil {
		if !validator.NilIfNull && validator.IfNull != nil {
			return validator.IfNull, nil
		}
		if validator.Type == reflect.Bool {
			return false, nil
		}
		return nil, nil
	}

	if validator.IsMapInterface {
		m, ioData := data, new(bytes.Buffer)
		var res map[string]interface{}
		err := json.NewEncoder(ioData).Encode(&m)
		if err != nil {
			return nil, errors.New("validatorType in '" + field + "' field is not valid object")
		}
		err = json.NewDecoder(ioData).Decode(&res)
		if err != nil {
			return nil, errors.New("validatorType in '" + field + "' field is not valid object")
		}
		return res, nil
	}

	// validatorType type validation
	dataType := reflect.TypeOf(data).Kind()
	handleIntOnHttpJson := dataFrom == fromHttpJson && isIntegerFamily(validator.Type) && isIntegerFamily(dataType)
	customData := !(!validator.UUID &&
		!validator.IPV4 &&
		!validator.UUIDToString &&
		!validator.IPv4OptionalPrefix &&
		!validator.Email &&
		validator.Enum == nil &&
		!validator.File &&
		!validator.IPV4Network &&
		validator.RegexString == "")

	if dataType == reflect.Slice && !validator.Null && len(ToInterfaceSlice(data)) == 0 {
		return nil, errors.New("you need to input validatorType in '" + field + "' field")
	}

	if dataType != validator.Type && !customData && !handleIntOnHttpJson {
		if dataFrom == fromHttpJson && isIntegerFamily(validator.Type) {
			validator.Type = reflect.Int
		}
		return nil, errors.New("the field '" + field + "' should be '" + validator.Type.String() + "'")
	}

	if validator.File {
		//this will return FileRequest
		return data, nil
	}

	if validator.RegexString != "" {
		if dataType != reflect.String {
			if validator.CustomMsg.OnRegexString != nil {
				return nil, buildMessage(*validator.CustomMsg.OnRegexString, MessageMeta{Field: field})
			}
			return nil, errors.New("the field '" + field + "' should be string")
		}
		regex := regexp.MustCompile(validator.RegexString)
		if !regex.MatchString(data.(string)) {
			if validator.CustomMsg.OnRegexString != nil {
				return nil, buildMessage(*validator.CustomMsg.OnRegexString, MessageMeta{Field: field})
			}
			return nil, errors.New("the field '" + field + "' is not valid regex")
		}
		return data, nil
	}

	if validator.Enum != nil {
		enumType := reflect.TypeOf(validator.Enum.Items)
		if enumType.Kind() == reflect.Slice {
			enumValue := reflect.ValueOf(validator.Enum.Items)
			if dataType != enumType.Elem().Kind() {
				return nil, errors.New("the field '" + field + "' should be '" + enumType.Elem().Kind().String() + "'")
			}
			switch dataType {
			case reflect.Int:
				var values []int
				for i := 0; i < enumValue.Len(); i++ {
					values = append(values, int(enumValue.Index(i).Int()))
				}
				if !valueInList[int](values, data.(int), isEqualInt) {
					return nil, errors.New(fmt.Sprintf("the field '%s' value is not in enum list%v", field, values))
				}
			case reflect.Int64:
				var values []int64
				for i := 0; i < enumValue.Len(); i++ {
					values = append(values, enumValue.Index(i).Int())
				}
				if !valueInList[int64](values, data.(int64), isEqualInt64) {
					return nil, errors.New(fmt.Sprintf("the field '%s' value is not in enum list%v", field, values))
				}
			case reflect.Float64:
				var values []float64
				for i := 0; i < enumValue.Len(); i++ {
					values = append(values, enumValue.Index(i).Float())
				}
				if !valueInList[float64](values, data.(float64), isEqualFloat64) {
					return nil, errors.New(fmt.Sprintf("the field '%s' value is not in enum list%v", field, values))
				}
			case reflect.String:
				var values []string
				for i := 0; i < enumValue.Len(); i++ {
					values = append(values, enumValue.Index(i).String())
				}
				if !valueInList[string](values, data.(string), isEqualString) {
					return nil, errors.New(fmt.Sprintf("the field '%s' value is not in enum list%v", field, values))
				}
			default:
				panic("not support type validatorType for enum value")
			}
		}
		return data, nil
	}

	if validator.UUIDToString {
		validator.UUID = true
	}

	if validator.UUID {
		errMsg := errors.New("the field '" + field + "' it's not valid uuid")
		stringUuid, ok := data.(string)
		if !ok {
			return nil, errMsg
		}
		dataUuid, err := uuid.Parse(stringUuid)
		if err != nil {
			return nil, errMsg
		}
		if validator.UUIDToString {
			return stringUuid, nil
		}
		return dataUuid, nil
	}

	if validator.Email {
		if reflect.TypeOf(data).Kind() != reflect.String || !isEmail(data.(string)) {
			return nil, errors.New("field " + field + " is not valid email")
		}
	}

	if validator.IPV4 {
		errMsg := errors.New("the field '" + field + "' it's not valid IP")
		stringIp, ok := data.(string)
		if !ok {
			return nil, errMsg
		}
		if !isIPv4Valid(stringIp) {
			return nil, errMsg
		}
		return stringIp, nil
	}

	if validator.IPV4Network {
		errMsg := errors.New("the field '" + field + "' it's not valid IP Network")
		stringIp, ok := data.(string)
		if !ok {
			return nil, errMsg
		}
		if !isIPv4NetworkValid(stringIp) {
			return nil, errMsg
		}
		return stringIp, nil
	}

	if validator.IPv4OptionalPrefix {
		errMsg := errors.New("the field '" + field + "' it's not valid IP")
		stringIp, ok := data.(string)
		if !ok {
			return nil, errMsg
		}
		splitIp := strings.Split(stringIp, "/")
		if len(splitIp) > 2 {
			return nil, errMsg
		} else if len(splitIp) == 1 {
			if !isIPv4Valid(splitIp[0]) {
				return nil, errMsg
			}
			return stringIp, nil
		} else if len(splitIp) == 2 {
			if !isIPv4Valid(splitIp[0]) {
				return nil, errMsg
			}
			prefix, err := strconv.Atoi(splitIp[1])
			if err != nil {
				return nil, errMsg
			}
			if prefix < 0 || prefix > 32 {
				return nil, errMsg
			}
			return stringIp, nil
		} else if len(stringIp) == 0 {
			return stringIp, nil
		}
		return nil, errMsg
	}

	if !validator.Null && reflect.TypeOf(data).Kind() == reflect.Bool && data == nil {
		return false, nil
	}

	if validator.Min != nil && data != nil {
		if reflect.String == dataType {
			if total := utf8.RuneCountInString(data.(string)); total < *validator.Min {
				return nil, errors.New(
					fmt.Sprintf("the field '%s' should be or greater than %v character", field, *validator.Min),
				)
			}
		} else if reflect.Float64 == dataType {
			if int(data.(float64)) < *validator.Min {
				return nil, errors.New(
					fmt.Sprintf("the field '%s' should be or greater than %v", field, *validator.Min),
				)
			}
		}
	}

	if validator.Max != nil && data != nil {
		if reflect.String == dataType {
			if total := utf8.RuneCountInString(data.(string)); total > *validator.Max {
				return nil, errors.New(
					fmt.Sprintf("the field '%s' should be or lower than %v character", field, *validator.Max),
				)
			}
		} else if reflect.Float64 == dataType {
			if int(data.(float64)) > *validator.Max {
				return nil, errors.New(
					fmt.Sprintf("the field '%s' should be or lower than %v", field, *validator.Max),
				)
			}
		}
	}

	return data, nil
}

func SetTotal(total int) *int {
	return &total
}

func SetMessage(msg string) *string { return &msg }

func ToInterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func convertValue(newValue interface{}, kind reflect.Kind, data reflect.Value, pointer bool) error {
	errNotSupport := errors.New("not support data")
	switch kind {
	case reflect.Int:
		converted, ok := newValue.(float64)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetInt(int64(converted))
		}
	case reflect.Int8:
		converted, ok := newValue.(int64)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetInt(converted)
		}
	case reflect.Int16:
		converted, ok := newValue.(int64)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetInt(converted)
		}
	case reflect.Int32:
		converted, ok := newValue.(int64)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetInt(converted)
		}
	case reflect.Int64:
		converted, ok := newValue.(int64)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetInt(converted)
		}
	case reflect.Uint:
		converted, ok := newValue.(float64)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetUint(uint64(converted))
		}
	case reflect.Uint8:
		converted, ok := newValue.(float64)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetUint(uint64(converted))
		}
	case reflect.Uint16:
		converted, ok := newValue.(float64)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetUint(uint64(converted))
		}
	case reflect.Uint32:
		converted, ok := newValue.(float64)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetUint(uint64(converted))
		}
	case reflect.Uint64:
		converted, ok := newValue.(float64)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetUint(uint64(converted))
		}
	case reflect.Float32:
		converted, ok := newValue.(float64)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetFloat(float64(converted))
		}
	case reflect.Float64:
		converted, ok := newValue.(float64)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetFloat(converted)
		}
	case reflect.String:
		converted, ok := newValue.(string)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetString(converted)
		}
	case reflect.Bool:
		converted, ok := newValue.(bool)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetBool(converted)
		}
	case reflect.Complex64:
		converted, ok := newValue.(complex128)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetComplex(complex128(complex64(converted)))
		}
	case reflect.Complex128:
		converted, ok := newValue.(complex128)
		if !ok {
			return errNotSupport
		}
		if pointer {
			data.Set(reflect.ValueOf(&converted))
		} else {
			data.SetComplex(converted)
		}
	case reflect.Interface:
		data.Set(reflect.ValueOf(newValue))
	default:
		return errNotSupport
	}
	return nil
}

func getAllKeys(data map[string]interface{}) (allKeysInMap []string) {
	for key, _ := range data {
		allKeysInMap = append(allKeysInMap, key)
	}
	return
}

func isDataInList[T validatorType](key T, data []T) (result bool) {
	for _, val := range data {
		if val == key {
			return true
		}
	}
	return
}

func isIntegerFamily(dataType reflect.Kind) bool {
	switch dataType {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Float32, reflect.Float64:
		return true
	}
	return false
}
