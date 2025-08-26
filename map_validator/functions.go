package map_validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
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
	actualLengthVar := "${actual_length}"
	expectedMinLengthVar := "${expected_min_length}"
	expectedMaxLengthVar := "${expected_max_length}"
	uniqueOriginVar := "${unique_origin}"
	uniqueTargetVar := "${unique_target}"
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
	if strings.Contains(msg, actualLengthVar) {
		if meta.ActualLength != nil {
			v := *meta.ActualLength
			msg = strings.ReplaceAll(msg, actualLengthVar, fmt.Sprintf("%v", v))
		}
	}
	if strings.Contains(msg, expectedMinLengthVar) {
		if meta.ExpectedMinLength != nil {
			v := *meta.ExpectedMinLength
			msg = strings.ReplaceAll(msg, expectedMinLengthVar, fmt.Sprintf("%v", v))
		}
	}
	if strings.Contains(msg, expectedMaxLengthVar) {
		if meta.ExpectedMaxLength != nil {
			v := *meta.ExpectedMaxLength
			msg = strings.ReplaceAll(msg, expectedMaxLengthVar, fmt.Sprintf("%v", v))
		}
	}
	if strings.Contains(msg, uniqueOriginVar) {
		if meta.Field != nil {
			v := *meta.UniqueOrigin
			msg = strings.ReplaceAll(msg, uniqueOriginVar, v)
		}
	}
	if strings.Contains(msg, uniqueTargetVar) {
		if meta.Field != nil {
			v := *meta.UniqueTarget
			msg = strings.ReplaceAll(msg, uniqueTargetVar, v)
		}
	}
	return errors.New(msg)
}

func validateRecursive(pChain ChainerType, wrapper RulesWrapper, key string, data map[string]interface{}, rule Rules, loadedFrom loadFromType) (interface{}, error) {
	//child and parent chain
	var res interface{}
	var err error
	chainKey := key
	if rule.isList() {
		chainKey = fmt.Sprintf("%s[%d]", pChain.GetKey(), len(pChain.GetChildren())-1)
	}
	cChain := pChain.AddChild().SetKey(chainKey)
	var endOfLoop bool
	if wrapper != nil && wrapper.getSetting().Strict {
		var allowedKeys []string
		keys := getAllKeys(data)
		for XKey, _ := range wrapper.getRules() {
			allowedKeys = append(allowedKeys, XKey)
		}
		for _, XKey := range keys {
			if !isDataInList(XKey, allowedKeys) {
				return nil, fmt.Errorf("'%s' is not allowed key", XKey)
			}
		}
	}

	if rule.isList() {
		// res, err = validateInterface(data, rule, loadedFrom)
		// if err != nil {
		// 	return nil, err
		// }
		return nil, nil
	} else {
		res, err = validateMapInterface(key, data, rule, loadedFrom)
		if err != nil {
			return nil, err
		}
	}

	if res != nil {
		cChain.SetValue(res)
	}

	if wrapper != nil {
		// add unique values
		if res != nil && len(rule.Unique) > 0 {
			cChain.SetUniques(rule.Unique)
		}

		// add custom message values
		if res != nil && rule.CustomMsg.isNotNil() {
			cChain.SetCustomMsg(&rule.CustomMsg)
		}

		// put filled and null fields
		if wrapper.getFilledField() == nil {
			wrapper.setFilledField(&[]string{})
		}
		if wrapper.getNullFields() == nil {
			wrapper.setNullFields(&[]string{})
		}

		if res != nil {
			wrapper.appendFilledField(key)
		} else {
			wrapper.appendNullFields(key)
		}

		if len(wrapper.getRules()) == len(*wrapper.getNullFields())+len(*wrapper.getFilledField()) {
			endOfLoop = true
		}

		for _, mptr := range wrapper.getManipulator() {
			if key == mptr.Field {
				cChain.SetManipulator(mptr.Func)
			}
		}
	}

	// put required without values
	if wrapper != nil && len(rule.RequiredWithout) > 0 {
		for _, unique := range rule.RequiredWithout {
			if wrapper.getRequiredWithout() == nil {
				wrapper.setRequiredWithout(&map[string][]string{})
			}

			if _, exists := (*wrapper.getRequiredWithout())[unique]; !exists {
				(*wrapper.getRequiredWithout())[unique] = []string{}
			}
			(*wrapper.getRequiredWithout())[unique] = append((*wrapper.getRequiredWithout())[unique], key)
		}
	}

	if endOfLoop && wrapper != nil && wrapper.getRequiredWithout() != nil {
		for _, field := range *wrapper.getNullFields() {
			var required bool
			dependenciesField := (*wrapper.getRequiredWithout())[field]
			if len(dependenciesField) == 0 {
				continue
			}
			for _, XField := range dependenciesField {
				if isDataInList(XField, *wrapper.getFilledField()) {
					required = true
				}
			}
			if !required {
				return nil, fmt.Errorf("if field '%s' is null you need to put value in %v field", field, dependenciesField)
			}
		}
	}

	// put required if values
	if wrapper != nil && len(rule.RequiredIf) > 0 {
		for _, unique := range rule.RequiredIf {
			if wrapper.getRequiredIf() == nil {
				wrapper.setRequiredIf(&map[string][]string{})
			}

			if _, exists := (*wrapper.getRequiredIf())[unique]; !exists {
				(*wrapper.getRequiredIf())[unique] = []string{}
			}
			(*wrapper.getRequiredIf())[unique] = append((*wrapper.getRequiredIf())[unique], key)
		}
	}

	if endOfLoop && wrapper != nil && wrapper.getRequiredIf() != nil {
		for _, field := range *wrapper.getFilledField() {
			var required bool
			dependenciesField := (*wrapper.getRequiredIf())[field]
			if len(dependenciesField) == 0 {
				continue
			}
			for _, XField := range dependenciesField {
				if !isDataInList(XField, *wrapper.getNullFields()) {
					required = true
				}
			}
			if !required {
				return nil, fmt.Errorf("if field '%s' is filled you need to put value in %v field also", field, dependenciesField)
			}
		}
	}

	// if list
	if rule.Object != nil && res != nil {
		for keyX, ruleX := range rule.Object.getRules() {
			_, err = validateRecursive(cChain, rule.Object, keyX, res.(map[string]interface{}), ruleX, fromJSONEncoder)
			if err != nil {
				return nil, err
			}
		}
	}

	if rule.ListObject != nil && res != nil {
		listRes := res.([]interface{})
		var manipulated []interface{}
		for _, xRes := range listRes {
			tmpChain := newChainer().SetKey(chainKey)
			for keyX, ruleX := range rule.ListObject.getRules() {
				_, err = validateRecursive(tmpChain, rule.ListObject, keyX, xRes.(map[string]interface{}), ruleX, fromJSONEncoder)
				if err != nil {
					return nil, err
				}
			}
			// collect validated/manipulated item data back into the slice
			// ensure only fields defined in ListObject rules are included
			itemMapFull := tmpChain.GetResult().ToMap()
			filtered := make(map[string]interface{})
			for keyAllowed := range rule.ListObject.getRules() {
				if val, ok := itemMapFull[keyAllowed]; ok {
					filtered[keyAllowed] = val
				}
			}
			manipulated = append(manipulated, filtered)
		}
		cChain.SetValue(manipulated)
	}

	return res, nil
}

func validateInterface(dataTemp interface{}, validator Rules, dataFrom loadFromType) (interface{}, error) {
	panic("not implemented")
}

func validateMapInterface(field string, dataTemp map[string]interface{}, validator Rules, dataFrom loadFromType) (interface{}, error) {
	//var oldIntType reflect.Kind
	data := dataTemp[field]
	var sliceData []interface{}
	//var isListObject bool

	if len(validator.RequiredWithout) > 0 || len(validator.RequiredIf) > 0 {
		validator.Null = true
	}

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
		regex, err := regexp.Compile(validator.RegexString)
		if err != nil {
			return nil, err
		}
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
					return nil, fmt.Errorf("the field '%s' value is not in enum list%v", field, values)
				}
			case reflect.Int64:
				var values []int64
				for i := 0; i < enumValue.Len(); i++ {
					values = append(values, enumValue.Index(i).Int())
				}
				if !valueInList[int64](values, data.(int64), isEqualInt64) {
					return nil, fmt.Errorf("the field '%s' value is not in enum list%v", field, values)
				}
			case reflect.Float64:
				var values []float64
				for i := 0; i < enumValue.Len(); i++ {
					values = append(values, enumValue.Index(i).Float())
				}
				if !valueInList[float64](values, data.(float64), isEqualFloat64) {
					return nil, fmt.Errorf("the field '%s' value is not in enum list%v", field, values)
				}
			case reflect.String:
				var values []string
				for i := 0; i < enumValue.Len(); i++ {
					values = append(values, enumValue.Index(i).String())
				}
				if !valueInList[string](values, data.(string), isEqualString) {
					return nil, fmt.Errorf("the field '%s' value is not in enum list%v", field, values)
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
		var isErr bool
		var actualLength int64
		err := fmt.Errorf("the field '%s' should be or greater than %v", field, *validator.Min)
		if reflect.String == dataType {
			if total := utf8.RuneCountInString(data.(string)); int64(total) < *validator.Min {
				isErr = true
				actualLength = int64(total)
			}
		} else if isIntegerFamily(dataType) {
			strData := removeAfter(fmt.Sprintf("%v", data), "e+")
			num := extractInteger(strData)
			if num < *validator.Min {
				isErr = true
				actualLength = num
			}
		} else if reflect.Slice == dataType {
			total := int64(len(sliceData))
			if total < *validator.Min {
				isErr = true
				actualLength = total
			}
		}

		if isErr {
			if validator.CustomMsg.OnMin != nil {
				return nil, buildMessage(*validator.CustomMsg.OnMin, MessageMeta{
					Field:             &field,
					ExpectedMinLength: SetTotal(*validator.Min),
					ActualLength:      SetTotal(actualLength),
				})
			}
			return nil, err
		}
	}

	if validator.Max != nil && data != nil {
		var isErr bool
		var actualLength int64
		err := fmt.Errorf("the field '%s' should be or lower than %v", field, *validator.Max)
		if reflect.String == dataType {
			if total := utf8.RuneCountInString(data.(string)); int64(total) > *validator.Max {
				isErr = true
				actualLength = int64(total)
			}
		} else if isIntegerFamily(dataType) {
			strData := removeAfter(fmt.Sprintf("%v", data), "e+")
			num := extractInteger(strData)
			if num > *validator.Max {
				isErr = true
				actualLength = num
			}
		} else if reflect.Slice == dataType {
			total := int64(len(sliceData))
			if total > *validator.Max {
				isErr = true
				actualLength = total
			}
		}

		if isErr {
			if validator.CustomMsg.OnMax != nil {
				return nil, buildMessage(*validator.CustomMsg.OnMax, MessageMeta{
					Field:             &field,
					ExpectedMaxLength: SetTotal(*validator.Max),
					ActualLength:      SetTotal(actualLength),
				})
			}
			return nil, err
		}
	}

	return data, nil
}

func SetTotal(total int64) *int64 {
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

func removeAfter(data, after string) string {
	split := strings.Split(data, after)
	if len(split) > 0 {
		return split[0]
	}
	return data
}

func extractInteger(data string) int64 {
	var intStr string
	re := regexp.MustCompile("[0-9]+")
	resStr := re.FindAllString(data, -1)
	for _, val := range resStr {
		intStr += val
	}
	res, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return 0
	}
	return res
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
