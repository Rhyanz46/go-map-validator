package mapValidator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

type data interface {
	int | string
}

func IsIPv4PrefixValid(prefix string) (res bool) {
	for _, allowPrefix := range []string{"8", "16", "24", "32"} {
		if prefix == allowPrefix {
			res = true
			break
		}
	}
	return
}

func IsEmail(email string) bool {
	ok := strings.Contains(email, "@")
	if !ok {
		return false
	}
	ok = strings.Contains(strings.Split(email, "@")[1], ".")
	return ok
}

func IsIPv4Valid(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

func IsIPv4NetworkValid(ip string) bool {
	if IsIPv4Valid(ip) {
		ipString := strings.Split(ip, ".")
		if ipString[3] == "0" {
			return true
		}
	}
	return false
}

func Validate(field string, dataTemp map[string]interface{}, validator RequestDataValidator) (interface{}, error) {
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
			return nil, errors.New("data in '" + field + "' field is not valid object")
		}
		err = json.NewDecoder(ioData).Decode(&res)
		if err != nil {
			return nil, errors.New("data in '" + field + "' field is not valid object")
		}
		return res, nil
	}

	// data type validation
	dataType := reflect.TypeOf(data).Kind()
	if validator.Type == reflect.Int {
		validator.Type = reflect.Float64
	}

	if dataType == reflect.Slice && !validator.Null && len(ToInterfaceSlice(data)) == 0 {
		return nil, errors.New("you need to input data in '" + field + "' field")
	}

	if !validator.UUID && !validator.IPV4 && dataType != validator.Type && !validator.UUIDToString && !validator.IPv4OptionalPrefix && !validator.Email {
		return nil, errors.New("the field '" + field + "' should be '" + validator.Type.String() + "'")
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
		if reflect.TypeOf(data).Kind() != reflect.String || !IsEmail(data.(string)) {
			return nil, errors.New("field " + field + " is not valid email")
		}
	}

	if validator.IPV4 {
		errMsg := errors.New("the field '" + field + "' it's not valid IP")
		stringIp, ok := data.(string)
		if !ok {
			return nil, errMsg
		}
		if !IsIPv4Valid(stringIp) {
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
			if !IsIPv4Valid(splitIp[0]) {
				return nil, errMsg
			}
			return stringIp, nil
		} else if len(splitIp) == 2 {
			if !IsIPv4Valid(splitIp[0]) {
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

func ToPointer[T data](data T) (res *T) {
	res = &data
	return
}

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
