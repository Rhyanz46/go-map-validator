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
	// Only use indexed keys for ListObject, not for primitive List
	if rule.isList() && rule.ListObject != nil {
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

	res, err = validate(key, data, rule, loadedFrom)
	if err != nil {
		return nil, err
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
			if m, ok := xRes.(map[string]interface{}); ok {
				// Validate as object with the provided child rules
				tmpChain := newChainer().SetKey(chainKey)
				for keyX, ruleX := range rule.ListObject.getRules() {
					_, err = validateRecursive(tmpChain, rule.ListObject, keyX, m, ruleX, fromJSONEncoder)
					if err != nil {
						return nil, err
					}
				}
				// collect validated/manipulated item data back into the slice
				itemMapFull := tmpChain.GetResult().ToMap()
				filtered := make(map[string]interface{})
				for keyAllowed := range rule.ListObject.getRules() {
					if val, ok := itemMapFull[keyAllowed]; ok {
						filtered[keyAllowed] = val
					}
				}
				manipulated = append(manipulated, filtered)
			} else {
				// Fallback: treat as primitive element; validate against parent rule flags (e.g., UUID, Email)
				tmpRule := rule
				tmpRule.Object = nil
				tmpRule.ListObject = nil
				tmpRule.List = nil
				tmpPayload := map[string]interface{}{key: xRes}
				if _, err := validate(key, tmpPayload, tmpRule, fromJSONEncoder); err != nil {
					return nil, err
				}
				manipulated = append(manipulated, xRes)
			}
		}
		cChain.SetValue(manipulated)
	}

	return res, nil
}

// validate is the core field validator used across the package and tests
func validate(field string, dataTemp map[string]interface{}, validator Rules, dataFrom loadFromType) (interface{}, error) {
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

	// Keep original element-kind before any normalization
	originalElementKind := validator.Type

	// Pre-normalize list types to slice kind
	if validator.ListObject != nil || validator.List != nil {
		validator.Type = reflect.Slice
	}

	// Support legacy ListObject when List is not provided: enforce slice and return elements
	if validator.ListObject != nil && validator.List == nil {
		s, ok := toInterfaceSlice(data)
		if !ok {
			return nil, errors.New("field '" + field + "' is not valid list object")
		}
		return s, nil
	}

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

	// Early list handling to avoid container-level regex/enum/type checks
	if validator.List != nil {
		sliceDataX, ok := toInterfaceSlice(data)
		if !ok {
			return nil, errors.New("field '" + field + "' is not valid list")
		}

		// List of objects via Object rules or legacy ListObject
		if validator.Object != nil || validator.ListObject != nil {
			// ensure elements are objects for Object rules
			if validator.Object != nil {
				for _, it := range sliceDataX {
					if _, ok := it.(map[string]interface{}); !ok {
						return nil, errors.New("field '" + field + "' is not valid list object")
					}
				}
			}
			return sliceDataX, nil
		}

		// Primitive list: validate each element
		var elementMinPtr, elementMaxPtr *int64
		if lr, ok := validator.List.(*rulesWrapper); ok {
			// Treat ListRules.Min/Max as element content length constraints (string only)
			elementMinPtr = lr.ListRules.Min
			elementMaxPtr = lr.ListRules.Max
		}
		for _, it := range sliceDataX {
			tmpRule := validator
			tmpRule.List = nil
			tmpRule.ListObject = nil
			tmpRule.Object = nil
			// restore element type for per-item validation
			tmpRule.Type = originalElementKind
			// By default, do not carry container Min/Max into element checks
			tmpRule.Min = nil
			tmpRule.Max = nil
			// Apply element content constraints (pre-check) for string and numeric elements
			if it != nil {
				gotKind := reflect.TypeOf(it).Kind()
				// Resolve effective element kind: explicit Type if provided, else infer from value
				effectiveKind := originalElementKind
				if effectiveKind == reflect.Invalid {
					effectiveKind = gotKind
				}
				// String length constraints
				if effectiveKind == reflect.String && gotKind == reflect.String {
					if elementMinPtr != nil {
						if int64(utf8.RuneCountInString(it.(string))) < *elementMinPtr {
							return nil, fmt.Errorf("value in '%s' field should be or greater than %v", field, *elementMinPtr)
						}
					}
					if elementMaxPtr != nil {
						if int64(utf8.RuneCountInString(it.(string))) > *elementMaxPtr {
							return nil, fmt.Errorf("value in '%s' field should be or lower than %v", field, *elementMaxPtr)
						}
					}
				}
				// Numeric value constraints
				if isIntegerFamily(effectiveKind) && isIntegerFamily(gotKind) {
					// normalize to float64 for comparison
					var num float64
					switch v := it.(type) {
					case int:
						num = float64(v)
					case int8:
						num = float64(v)
					case int16:
						num = float64(v)
					case int32:
						num = float64(v)
					case int64:
						num = float64(v)
					case uint:
						num = float64(v)
					case uint8:
						num = float64(v)
					case uint16:
						num = float64(v)
					case uint32:
						num = float64(v)
					case uint64:
						num = float64(v)
					case float32:
						num = float64(v)
					case float64:
						num = v
					default:
						// fallback: let validate handle
						num = 0
					}
					if elementMinPtr != nil && num < float64(*elementMinPtr) {
						return nil, fmt.Errorf("value in '%s' field should be or greater than %v", field, *elementMinPtr)
					}
					if elementMaxPtr != nil && num > float64(*elementMaxPtr) {
						return nil, fmt.Errorf("value in '%s' field should be or lower than %v", field, *elementMaxPtr)
					}
				}
			}
			// Pre-check element type mismatch to craft a clearer wording
			// Only when element Type is explicitly set (avoid interfering with Enum/UUID/Regex-only rules)
			if it != nil && tmpRule.Type != reflect.Invalid {
				gotKind := reflect.TypeOf(it).Kind()
				expectedKind := tmpRule.Type
				allowIntCoerce := (dataFrom == fromHttpJson || dataFrom == fromJSONEncoder) && isIntegerFamily(expectedKind) && isIntegerFamily(gotKind)
				if gotKind != expectedKind && !allowIntCoerce {
					// Map kind to human-friendly noun (e.g., int/uint/float -> integer)
					noun := expectedKind.String()
					if isIntegerFamily(expectedKind) {
						noun = "integer"
					}
					return nil, fmt.Errorf("value in '%s' field should be %s", field, noun)
				}
			}

			tmpPayload := map[string]interface{}{field: it}
			if _, err := validate(field, tmpPayload, tmpRule, dataFrom); err != nil {
				return nil, err
			}
		}
		// list-size Min/Max come from outer rule (container size)
		var minPtr, maxPtr *int64
		if validator.Min != nil {
			minPtr = validator.Min
		}
		if validator.Max != nil {
			maxPtr = validator.Max
		}
		listLen := int64(len(sliceDataX))
		if minPtr != nil && listLen < *minPtr {
			return nil, fmt.Errorf("the field '%s' should be or greater than %v", field, *minPtr)
		}
		if maxPtr != nil && listLen > *maxPtr {
			return nil, fmt.Errorf("the field '%s' should be or lower than %v", field, *maxPtr)
		}
		return sliceDataX, nil
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

	// Helper function to build enum error message with custom or default text
	buildEnumErrorMessage := func(enumValues interface{}, enumType reflect.Type, actualType reflect.Kind) error {
		if validator.CustomMsg.OnEnumValueNotMatch != nil {
			expectedType := enumType.Elem().Kind()
			return buildMessage(*validator.CustomMsg.OnEnumValueNotMatch, MessageMeta{
				Field:        &field,
				ExpectedType: &expectedType,
				ActualType:   &actualType,
			})
		}
		return fmt.Errorf("the field '%s' value is not in enum list%v", field, enumValues)
	}

	if validator.Enum != nil {
		enumType := reflect.TypeOf(validator.Enum.Items)
		if enumType.Kind() == reflect.Slice {
			enumValue := reflect.ValueOf(validator.Enum.Items)
			// Handle integer family coercion for HTTP JSON like regular type validation
			if dataType != enumType.Elem().Kind() {
				// Allow type mismatch for integer family from HTTP JSON
				if (dataFrom == fromHttpJson || dataFrom == fromJSONEncoder) &&
					isIntegerFamily(enumType.Elem().Kind()) && isIntegerFamily(dataType) {
					// Type coercion will be handled in the switch cases below
				} else {
					// Use custom message for type mismatch if available, otherwise use custom enum message or default
					if validator.CustomMsg.OnTypeNotMatch != nil {
						expectedType := enumType.Elem().Kind()
						return nil, buildMessage(*validator.CustomMsg.OnTypeNotMatch, MessageMeta{
							Field:        &field,
							ExpectedType: &expectedType,
							ActualType:   &dataType,
						})
					}
					return nil, buildEnumErrorMessage(nil, enumType, dataType)
				}
			}

			// Handle cross-type enum validation for integer family from HTTP JSON
			if (dataFrom == fromHttpJson || dataFrom == fromJSONEncoder) &&
				isIntegerFamily(enumType.Elem().Kind()) && isIntegerFamily(dataType) &&
				dataType != enumType.Elem().Kind() {
				// Convert float64 JSON data to compare with integer family enum items
				if dataType == reflect.Float64 {
					floatData := data.(float64)
					enumKind := enumType.Elem().Kind()

					// For float32 and float64, no integer conversion check needed
					if enumKind == reflect.Float32 || enumKind == reflect.Float64 {
						// Handle float enum types
						switch enumKind {
						case reflect.Float32:
							var values []float32
							for i := 0; i < enumValue.Len(); i++ {
								values = append(values, float32(enumValue.Index(i).Float()))
							}
							if !valueInList[float32](values, float32(floatData), func(a, b float32) bool { return a == b }) {
								return nil, buildEnumErrorMessage(values, enumType, dataType)
							}
						case reflect.Float64:
							var values []float64
							for i := 0; i < enumValue.Len(); i++ {
								values = append(values, enumValue.Index(i).Float())
							}
							if !valueInList[float64](values, floatData, isEqualFloat64) {
								return nil, buildEnumErrorMessage(values, enumType, dataType)
							}
						}
						return data, nil
					}

					// Check if float64 can be safely converted to integer (no decimal part)
					if floatData == float64(int64(floatData)) {
						// Handle different integer enum types
						switch enumType.Elem().Kind() {
						case reflect.Int:
							var values []int
							for i := 0; i < enumValue.Len(); i++ {
								values = append(values, int(enumValue.Index(i).Int()))
							}
							if !valueInList[int](values, int(floatData), isEqualInt) {
								return nil, buildEnumErrorMessage(values, enumType, dataType)
							}
						case reflect.Int64:
							var values []int64
							for i := 0; i < enumValue.Len(); i++ {
								values = append(values, enumValue.Index(i).Int())
							}
							if !valueInList[int64](values, int64(floatData), isEqualInt64) {
								return nil, buildEnumErrorMessage(values, enumType, dataType)
							}
						case reflect.Int32:
							var values []int32
							for i := 0; i < enumValue.Len(); i++ {
								values = append(values, int32(enumValue.Index(i).Int()))
							}
							if !valueInList[int32](values, int32(floatData), func(a, b int32) bool { return a == b }) {
								return nil, buildEnumErrorMessage(values, enumType, dataType)
							}
						case reflect.Int16:
							var values []int16
							for i := 0; i < enumValue.Len(); i++ {
								values = append(values, int16(enumValue.Index(i).Int()))
							}
							if !valueInList[int16](values, int16(floatData), func(a, b int16) bool { return a == b }) {
								return nil, buildEnumErrorMessage(values, enumType, dataType)
							}
						case reflect.Int8:
							var values []int8
							for i := 0; i < enumValue.Len(); i++ {
								values = append(values, int8(enumValue.Index(i).Int()))
							}
							if !valueInList[int8](values, int8(floatData), func(a, b int8) bool { return a == b }) {
								return nil, buildEnumErrorMessage(values, enumType, dataType)
							}
						case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
							// Handle unsigned integers
							var values []uint64
							for i := 0; i < enumValue.Len(); i++ {
								values = append(values, enumValue.Index(i).Uint())
							}
							if floatData < 0 || !valueInList[uint64](values, uint64(floatData), func(a, b uint64) bool { return a == b }) {
								return nil, buildEnumErrorMessage(values, enumType, dataType)
							}
						}
						return data, nil
					} else {
						// Float has decimal part, cannot convert to integer enum
						return nil, errors.New("the field '" + field + "' should be '" + enumType.Elem().Kind().String() + "'")
					}
				}
			}

			switch dataType {
			case reflect.Int:
				var values []int
				for i := 0; i < enumValue.Len(); i++ {
					values = append(values, int(enumValue.Index(i).Int()))
				}
				if !valueInList[int](values, data.(int), isEqualInt) {
					return nil, buildEnumErrorMessage(values, enumType, dataType)
				}
			case reflect.Int64:
				var values []int64
				for i := 0; i < enumValue.Len(); i++ {
					values = append(values, enumValue.Index(i).Int())
				}
				if !valueInList[int64](values, data.(int64), isEqualInt64) {
					return nil, buildEnumErrorMessage(values, enumType, dataType)
				}
			case reflect.Float64:
				var values []float64
				for i := 0; i < enumValue.Len(); i++ {
					values = append(values, enumValue.Index(i).Float())
				}
				if !valueInList[float64](values, data.(float64), isEqualFloat64) {
					return nil, buildEnumErrorMessage(values, enumType, dataType)
				}
			case reflect.String:
				var values []string
				for i := 0; i < enumValue.Len(); i++ {
					values = append(values, enumValue.Index(i).String())
				}
				if !valueInList[string](values, data.(string), isEqualString) {
					return nil, buildEnumErrorMessage(values, enumType, dataType)
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

	// legacy ListObject fallback occurs via early list handling
	if validator.AnonymousObject || validator.Object != nil {
		res, err := toMapStringInterface(data)
		if err != nil {
			return nil, errors.New("field '" + field + "' is not valid object")
		}
		return res, nil
	}

	// Ensure sliceData is available for legacy slice length checks
	if sliceData == nil && data != nil && reflect.TypeOf(data).Kind() == reflect.Slice {
		if s, ok := toInterfaceSlice(data); ok {
			sliceData = s
		}
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
