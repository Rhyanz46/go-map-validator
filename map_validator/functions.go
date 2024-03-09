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
	fieldVar := "${field}"
	expectedTypeVar := "${expected_type}"
	actualTypeVar := "${actual_type}"
	if strings.Contains(msg, fieldVar) {
		if meta.Field != nil {
			v := *meta.Field
			msg = strings.ReplaceAll(msg, fieldVar, v)
		}
	}
	if strings.Contains(msg, expectedTypeVar) {
		if meta.ExpectedType != nil {
			v := *meta.ExpectedType
			msg = strings.ReplaceAll(msg, expectedTypeVar, v.String())
		}
	}
	if strings.Contains(msg, actualTypeVar) {
		if meta.ActualType != nil {
			v := *meta.ActualType
			msg = strings.ReplaceAll(msg, actualTypeVar, v.String())
		}
	}
	return errors.New(msg)
}

func validateRecursive(wrapper *RulesWrapper, key string, data map[string]interface{}, rule Rules, loadedFrom loadFromType) (interface{}, error) {
	if wrapper != nil && wrapper.Setting.Strict {
		var allowedKeys []string
		keys := getAllKeys(data)
		for XKey, _ := range wrapper.Rules {
			allowedKeys = append(allowedKeys, XKey)
		}
		for _, XKey := range keys {
			if !isDataInList(XKey, allowedKeys) {
				return nil, errors.New(fmt.Sprintf("'%s' is not allowed key", XKey))
			}
		}
	}
	res, err := validate(key, data, rule, loadedFrom)
	if err != nil {
		return nil, err
	}

	if res != nil && len(rule.Unique) > 0 {
		for _, unique := range rule.Unique {
			newUniqueValues := make(map[string]map[string]interface{})
			if wrapper.uniqueValues == nil {
				wrapper.uniqueValues = &newUniqueValues
			} else {
				newUniqueValues = *wrapper.uniqueValues
			}
			for keyX, val := range newUniqueValues[unique] {
				if val == res {
					return nil, errors.New(fmt.Sprintf("value of '%s' and '%s' fields must be different", keyX, key))
				}
			}
			if newUniqueValues[unique] != nil {
				newUniqueValues[unique][key] = res
			} else {
				newUniqueValues[unique] = map[string]interface{}{key: res}
			}
			wrapper.uniqueValues = &newUniqueValues
		}
	}

	// if list
	if rule.Object != nil && res != nil {
		for keyX, ruleX := range rule.Object.Rules {
			_, err = validateRecursive(rule.Object, keyX, res.(map[string]interface{}), ruleX, fromJSONEncoder)
			if err != nil {
				//if rule.CustomMsg != nil #TODO: custom message for nested object
				return nil, err
			}
			//filledFields = append(filledFields, newFilledFields...) #TODO: get children fields fill or not filled
			//nullFields = append(nullFields, newNullFields...)
		}
	}

	// if object
	if rule.ListObject != nil && res != nil {
		listRes := res.([]interface{})
		for _, xRes := range listRes {
			for keyX, ruleX := range rule.ListObject.Rules {
				_, err = validateRecursive(rule.Object, keyX, xRes.(map[string]interface{}), ruleX, fromJSONEncoder)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return res, nil
}

func validate(field string, dataTemp map[string]interface{}, validator Rules, dataFrom loadFromType) (interface{}, error) {
	//var oldIntType reflect.Kind
	data := dataTemp[field]
	var sliceData []interface{}
	//var isListObject bool

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

	//if validator.ListObject != nil {
	//	res, err := toInterfaceSlice(data)
	//	if err != nil {
	//		return nil, errors.New("field '" + field + "' is not valid object")
	//	}
	//	return res, nil
	//}

	// validatorType type validation
	dataType := reflect.TypeOf(data).Kind()
	handleIntOnHttpJson := (dataFrom == fromHttpJson || dataFrom == fromJSONEncoder) && isIntegerFamily(validator.Type) && isIntegerFamily(dataType)
	customData := !(!validator.UUID &&
		!validator.IPV4 &&
		!validator.UUIDToString &&
		!validator.IPv4OptionalPrefix &&
		!validator.Email &&
		validator.Enum == nil &&
		validator.Object == nil &&
		validator.ListObject == nil &&
		!validator.AnonymousObject &&
		!validator.File &&
		!validator.IPV4Network &&
		validator.RegexString == "")

	//if dataType == reflect.Slice && !validator.Null && len(toInterfaceSlice(data)) == 0 {
	//	//if dataType == reflect.Slice && !validator.Null {
	//	return nil, errors.New("you need to input validatorType in '" + field + "' field")
	//}

	if dataType != validator.Type && !customData && !handleIntOnHttpJson {
		if (dataFrom == fromHttpJson || dataFrom == fromJSONEncoder) && isIntegerFamily(validator.Type) {
			validator.Type = reflect.Int
		}
		if validator.CustomMsg.OnTypeNotMatch != nil {
			return nil, buildMessage(*validator.CustomMsg.OnTypeNotMatch, MessageMeta{
				Field:        &field,
				ExpectedType: &validator.Type,
				ActualType:   &dataType,
			})
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
				return nil, buildMessage(*validator.CustomMsg.OnRegexString, MessageMeta{Field: &field})
			}
			return nil, errors.New("the field '" + field + "' should be string")
		}
		regex := regexp.MustCompile(validator.RegexString)
		if !regex.MatchString(data.(string)) {
			if validator.CustomMsg.OnRegexString != nil {
				return nil, buildMessage(*validator.CustomMsg.OnRegexString, MessageMeta{Field: &field})
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

	if validator.ListObject != nil {
		validator.Type = reflect.Slice
	}

	if validator.Type == reflect.Slice {
		sliceDataX, ok := toInterfaceSlice(data)
		if validator.ListObject != nil {
			if !ok {
				return nil, errors.New("field '" + field + "' is not valid list object")
			}
		}
		if !ok {
			return nil, errors.New("field '" + field + "' is not valid list")
		}
		sliceData = sliceDataX
		data = sliceDataX
	}

	if validator.AnonymousObject || validator.Object != nil {
		res, err := toMapStringInterface(data)
		if err != nil {
			return nil, errors.New("field '" + field + "' is not valid object")
		}
		return res, nil
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
		} else if reflect.Slice == dataType {
			if len(sliceData) < *validator.Min {
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
		} else if reflect.Slice == dataType {
			if len(sliceData) > *validator.Max {
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

func toInterfaceSlice(slice interface{}) ([]interface{}, bool) {

	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, false
	}

	if s.IsNil() {
		return nil, true
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, true
}

func toMapStringInterface(data interface{}) (map[string]interface{}, error) {
	m, ioData := data, new(bytes.Buffer)
	var res map[string]interface{}
	err := json.NewEncoder(ioData).Encode(&m)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(ioData).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
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
