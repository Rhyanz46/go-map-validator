package map_validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

func isEqualString(current, allowedField string) bool {
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
	if parsedIP == nil || parsedIP.To4() == nil {
		return false
	}
	// ::ffff:a.b.c.d parses as an IPv4-mapped IPv6 address whose To4() is
	// non-nil; only the dotted-quad form counts as IPv4 here.
	return parsedIP.String() == ip
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
	actualValueVar := "${actual_value}"
	enumValuesVar := "${enum_values}"
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
		if meta.UniqueOrigin != nil {
			v := *meta.UniqueOrigin
			msg = strings.ReplaceAll(msg, uniqueOriginVar, v)
		}
	}
	if strings.Contains(msg, uniqueTargetVar) {
		if meta.UniqueTarget != nil {
			v := *meta.UniqueTarget
			msg = strings.ReplaceAll(msg, uniqueTargetVar, v)
		}
	}
	if strings.Contains(msg, actualValueVar) {
		if meta.ActualValue != nil {
			msg = strings.ReplaceAll(msg, actualValueVar, *meta.ActualValue)
		}
	}
	if strings.Contains(msg, enumValuesVar) {
		if meta.EnumValues != nil {
			msg = strings.ReplaceAll(msg, enumValuesVar, *meta.EnumValues)
		}
	}
	return errors.New(msg)
}

// wrapperRunState holds the mutable per-scope state used during a single
// RunValidate call. It lives outside RulesWrapper so that rules can be shared
// across handlers, reused across runs, and even self-referenced (recursive
// Object rules) without state bleeding between scopes.
type wrapperRunState struct {
	filledField     []string
	nullFields      []string
	requiredWithout map[string][]string
	requiredIf      map[string][]string
}

func newWrapperRunState() *wrapperRunState {
	return &wrapperRunState{}
}

func validateRecursive(pChain ChainerType, wrapper RulesWrapper, state *wrapperRunState, key string, data map[string]interface{}, rule Rules, loadedFrom loadFromType) (interface{}, error) {
	//child and parent chain
	var res interface{}
	var err error
	nodeKey := key
	// Only use indexed keys for ListObject, not for primitive List
	if rule.isList() && rule.ListObject != nil {
		nodeKey = fmt.Sprintf("%s[%d]", pChain.GetKey(), len(pChain.GetChildren())-1)
	}
	cChain := pChain.AddChild().SetKey(nodeKey)
	var endOfLoop bool

	res, err = validate(key, data, rule, loadedFrom)
	if err != nil {
		return nil, err
	}

	if res != nil {
		cChain.SetValue(res)
	}

	if wrapper != nil && state != nil {
		// add unique values
		if res != nil && len(rule.Unique) > 0 {
			cChain.SetUniques(rule.Unique)
		}

		// add custom message values
		if res != nil && rule.CustomMsg.isNotNil() {
			cChain.SetCustomMsg(&rule.CustomMsg)
		}

		if res != nil {
			state.filledField = append(state.filledField, key)
		} else {
			state.nullFields = append(state.nullFields, key)
		}

		if len(wrapper.getRules()) == len(state.nullFields)+len(state.filledField) {
			endOfLoop = true
		}

		for _, mptr := range wrapper.getManipulator() {
			if key == mptr.Field {
				cChain.SetManipulator(mptr.Func)
			}
		}
	}

	// put required without values
	if state != nil && len(rule.RequiredWithout) > 0 {
		if state.requiredWithout == nil {
			state.requiredWithout = map[string][]string{}
		}
		for _, unique := range rule.RequiredWithout {
			state.requiredWithout[unique] = append(state.requiredWithout[unique], key)
		}
	}

	if endOfLoop && state != nil && state.requiredWithout != nil {
		deps := make([]string, 0, len(state.requiredWithout))
		for dep := range state.requiredWithout {
			deps = append(deps, dep)
		}
		sort.Strings(deps)
		for _, dep := range deps {
			// A dependency counts as filled when it is present in the scope
			// data with a non-nil value — even when it was never declared as
			// a rule. Restricting this to declared fields used to make the
			// whole check a silent no-op for undeclared dependencies.
			if depValue, ok := data[dep]; ok && depValue != nil {
				continue
			}
			dependenciesField := state.requiredWithout[dep]
			var required bool
			for _, XField := range dependenciesField {
				if isDataInList(XField, state.filledField) {
					required = true
				}
			}
			if !required {
				return nil, fmt.Errorf("if field '%s' is null you need to put value in %v field", dep, dependenciesField)
			}
		}
	}

	// put required if values
	if state != nil && len(rule.RequiredIf) > 0 {
		if state.requiredIf == nil {
			state.requiredIf = map[string][]string{}
		}
		for _, unique := range rule.RequiredIf {
			state.requiredIf[unique] = append(state.requiredIf[unique], key)
		}
	}

	if endOfLoop && state != nil && state.requiredIf != nil {
		deps := make([]string, 0, len(state.requiredIf))
		for dep := range state.requiredIf {
			deps = append(deps, dep)
		}
		sort.Strings(deps)
		for _, dep := range deps {
			if depValue, ok := data[dep]; !ok || depValue == nil {
				continue
			}
			dependenciesField := state.requiredIf[dep]
			var required bool
			for _, XField := range dependenciesField {
				if !isDataInList(XField, state.nullFields) {
					required = true
				}
			}
			if !required {
				return nil, fmt.Errorf("if field '%s' is filled you need to put value in %v field also", dep, dependenciesField)
			}
		}
	}

	// nested object — a rule that also carries List belongs to the
	// list-item branch below, where Object supplies the item rules
	if rule.Object != nil && rule.List == nil && res != nil {
		objRes, ok := res.(map[string]interface{})
		if !ok {
			// a rule carrying both Object and ListObject can land here with a
			// slice — report a validation error instead of panicking
			return nil, buildErrorMessage(key, "is not valid object")
		}
		if rule.Object.getSetting().Strict {
			if err := checkStrictKeys(rule.Object, objRes); err != nil {
				return nil, err
			}
		}
		innerState := newWrapperRunState()
		objRules := rule.Object.getRules()
		for _, keyX := range sortedKeys(objRules) {
			_, err = validateRecursive(cChain, rule.Object, innerState, keyX, objRes, objRules[keyX], fromJSONEncoder)
			if err != nil {
				return nil, err
			}
		}
	}

	// list of objects: ListObject, or List combined with Object rules
	innerWrapper := rule.ListObject
	if innerWrapper == nil && rule.List != nil && rule.Object != nil {
		innerWrapper = rule.Object
	}
	if innerWrapper != nil && res != nil {
		listRes, ok := res.([]interface{})
		if !ok {
			// symmetric guard for a rule carrying both Object and ListObject
			return nil, buildErrorMessage(key, "is not valid list object")
		}
		var manipulated []interface{}
		for _, xRes := range listRes {
			if m, ok := xRes.(map[string]interface{}); ok {
				// Validate as object with the provided child rules
				if innerWrapper.getSetting().Strict {
					if err := checkStrictKeys(innerWrapper, m); err != nil {
						return nil, err
					}
				}
				tmpChain := newChainer().SetKey(chainKey)
				itemState := newWrapperRunState()
				itemRules := innerWrapper.getRules()
				for _, keyX := range sortedKeys(itemRules) {
					_, err = validateRecursive(tmpChain, innerWrapper, itemState, keyX, m, itemRules[keyX], fromJSONEncoder)
					if err != nil {
						return nil, err
					}
				}
				// The item chain is detached from the root chain, so run its
				// manipulators and unique checks locally — otherwise they are
				// silently skipped by the root-level traversal in RunValidate.
				if err = tmpChain.GetResult().RunManipulator(); err != nil {
					return nil, err
				}
				tmpChain.GetResult().RunUniqueChecker()
				for _, itemErr := range tmpChain.GetResult().GetErrors() {
					if itemErr != nil {
						return nil, itemErr
					}
				}
				// collect validated/manipulated item data back into the slice
				itemMapFull := tmpChain.GetResult().ToMap()
				filtered := make(map[string]interface{})
				for keyAllowed := range itemRules {
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
		// Min/Max on a list-of-objects rule constrain the item count
		if rule.Min != nil && int64(len(manipulated)) < *rule.Min {
			return nil, buildErrorMessagef(key, "should be or greater than %v", *rule.Min)
		}
		if rule.Max != nil && int64(len(manipulated)) > *rule.Max {
			return nil, buildErrorMessagef(key, "should be or lower than %v", *rule.Max)
		}
		cChain.SetValue(manipulated)
	}

	return res, nil
}

// ValidateValue validates a single value against the given rules without field context
// This is the core validation logic that can be reused for different scenarios
// field parameter is optional - if not provided, "value" will be used in error messages
func ValidateValue(data interface{}, validator Rules, dataFrom LoadFromType, field ...string) (interface{}, error) {
	// Use provided field name or default to "value"
	fieldName := "value"
	if len(field) > 0 && field[0] != "" {
		fieldName = field[0]
	}

	return validateValueInternal(data, validator, loadFromType(dataFrom), fieldName)
}

// validate is the core field validator used across the package and tests
func validate(field string, dataTemp map[string]interface{}, validator Rules, dataFrom loadFromType) (interface{}, error) {
	// Handle RequiredWithout and RequiredIf fields by setting Null to true
	if len(validator.RequiredWithout) > 0 || len(validator.RequiredIf) > 0 {
		validator.Null = true
	}

	// Extract the data value for the field
	data := dataTemp[field]

	// Use the internal validation function
	return validateValueInternal(data, validator, dataFrom, field)
}

// buildErrorMessage creates a natural error message based on field name
func buildErrorMessage(field string, message string) error {
	if field == "value" {
		// For default field, use more natural message without "field" prefix
		return errors.New(message)
	}
	// For specific fields, use the traditional format
	return errors.New("the field '" + field + "' " + message)
}

// buildErrorMessagef creates a natural error message with formatting
func buildErrorMessagef(field string, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	if field == "value" {
		return errors.New(message)
	}
	return fmt.Errorf("the field '%s' "+format, append([]interface{}{field}, args...)...)
}

// validateValueInternal contains the actual validation logic
func validateValueInternal(data interface{}, validator Rules, dataFrom loadFromType, field string) (interface{}, error) {
	var sliceData []interface{}

	// null validation
	if !validator.Null && data == nil {
		if field == "value" {
			return nil, errors.New("value is required")
		}
		return nil, errors.New("we need '" + field + "' field")
	} else if validator.Null && data == nil {
		if !validator.NilIfNull && validator.IfNull != nil {
			// A default value goes through the same validation as any other
			// value — an invalid default is a rule error, not a silent bypass.
			return validateValueInternal(validator.IfNull, validator, dataFrom, field)
		}
		return nil, nil
	}

	// Any() short-circuit: skip all type/format/range checks. The field is
	// already known to be present (null check above passed), so we return
	// the value verbatim for downstream Bind.
	if validator.Any {
		return data, nil
	}

	//if validator.ListObject != nil {
	//	res, err := toInterfaceSlice(data)
	//	if err != nil {
	//		return nil, buildErrorMessage(field, "is not valid object")
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
			return nil, buildErrorMessage(field, "is not valid list object")
		}
		return s, nil
	}

	// validatorType type validation
	dataType := reflect.TypeOf(data).Kind()
	handleIntOnHttpJson := integerCoercion(dataFrom, validator.Type, dataType)
	customData := !(!validator.UUID &&
		!validator.IPV4 &&
		!validator.IPV4Network &&
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
		if coercesNumbers(dataFrom) && isIntegerFamily(validator.Type) {
			validator.Type = reflect.Int
		}
		if validator.CustomMsg.OnTypeNotMatch != nil {
			return nil, buildMessage(*validator.CustomMsg.OnTypeNotMatch, MessageMeta{
				Field:        &field,
				ExpectedType: &validator.Type,
				ActualType:   &dataType,
			})
		}
		return nil, buildErrorMessage(field, "should be '"+validator.Type.String()+"'")
	}

	// Integer rules must reject fractional values from coerced sources (JSON
	// and forms decode numbers loosely). Integral floats still pass, so
	// {"qty": 2} binds fine while {"qty": 2.5} fails here instead of at Bind.
	if handleIntOnHttpJson && isIntegerKind(validator.Type) && (dataType == reflect.Float32 || dataType == reflect.Float64) {
		if num := numericAsFloat64(data); num != math.Trunc(num) {
			if validator.CustomMsg.OnTypeNotMatch != nil {
				return nil, buildMessage(*validator.CustomMsg.OnTypeNotMatch, MessageMeta{
					Field:        &field,
					ExpectedType: &validator.Type,
					ActualType:   &dataType,
				})
			}
			return nil, buildErrorMessage(field, "should be '"+validator.Type.String()+"'")
		}
	}

	// Early list handling to avoid container-level regex/enum/type checks
	if validator.List != nil {
		sliceDataX, ok := toInterfaceSlice(data)
		if !ok {
			return nil, buildErrorMessage(field, "is not valid list")
		}

		// List of objects via Object rules or legacy ListObject
		if validator.Object != nil || validator.ListObject != nil {
			// ensure elements are objects for Object rules
			if validator.Object != nil {
				for _, it := range sliceDataX {
					if _, ok := it.(map[string]interface{}); !ok {
						return nil, buildErrorMessage(field, "is not valid list object")
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
					strItem, _ := asString(it)
					if elementMinPtr != nil {
						actualLen := int64(utf8.RuneCountInString(strItem))
						if actualLen < *elementMinPtr {
							if validator.CustomMsg.OnMin != nil {
								return nil, buildMessage(*validator.CustomMsg.OnMin, MessageMeta{
									Field:             &field,
									ExpectedMinLength: elementMinPtr,
									ActualLength:      anyPtr(actualLen),
								})
							}
							return nil, buildErrorMessagef(field, "should be or greater than %v", *elementMinPtr)
						}
					}
					if elementMaxPtr != nil {
						actualLen := int64(utf8.RuneCountInString(strItem))
						if actualLen > *elementMaxPtr {
							if validator.CustomMsg.OnMax != nil {
								return nil, buildMessage(*validator.CustomMsg.OnMax, MessageMeta{
									Field:             &field,
									ExpectedMaxLength: elementMaxPtr,
									ActualLength:      anyPtr(actualLen),
								})
							}
							return nil, buildErrorMessagef(field, "should be or lower than %v", *elementMaxPtr)
						}
					}
				}
				// Numeric value constraints
				if isIntegerFamily(effectiveKind) && isIntegerFamily(gotKind) {
					// normalize to float64 for comparison (named-type safe)
					num := numericAsFloat64(it)
					if elementMinPtr != nil && num < float64(*elementMinPtr) {
						if validator.CustomMsg.OnMin != nil {
							actualLen := int64(num)
							return nil, buildMessage(*validator.CustomMsg.OnMin, MessageMeta{
								Field:             &field,
								ExpectedMinLength: elementMinPtr,
								ActualLength:      anyPtr(actualLen),
							})
						}
						return nil, buildErrorMessagef(field, "should be or greater than %v", *elementMinPtr)
					}
					if elementMaxPtr != nil && num > float64(*elementMaxPtr) {
						if validator.CustomMsg.OnMax != nil {
							actualLen := int64(num)
							return nil, buildMessage(*validator.CustomMsg.OnMax, MessageMeta{
								Field:             &field,
								ExpectedMaxLength: elementMaxPtr,
								ActualLength:      anyPtr(actualLen),
							})
						}
						return nil, buildErrorMessagef(field, "should be or lower than %v", *elementMaxPtr)
					}
				}
			}
			// Pre-check element type mismatch to craft a clearer wording
			// Only when element Type is explicitly set (avoid interfering with Enum/UUID/Regex-only rules)
			if it != nil && tmpRule.Type != reflect.Invalid {
				gotKind := reflect.TypeOf(it).Kind()
				expectedKind := tmpRule.Type
				allowIntCoerce := integerCoercion(dataFrom, expectedKind, gotKind)
				if gotKind != expectedKind && !allowIntCoerce {
					// Map kind to human-friendly noun (e.g., int/uint/float -> integer)
					noun := expectedKind.String()
					if isIntegerFamily(expectedKind) {
						noun = "integer"
					}
					return nil, buildErrorMessagef(field, "should be %s", noun)
				}
			}

			// Recursive validation for each element
			if _, err := validateValueInternal(it, tmpRule, dataFrom, field); err != nil {
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
			return nil, buildErrorMessagef(field, "should be or greater than %v", *minPtr)
		}
		if maxPtr != nil && listLen > *maxPtr {
			return nil, buildErrorMessagef(field, "should be or lower than %v", *maxPtr)
		}
		// enforce element uniqueness when declared via ListRules.Unique
		if lr, ok := validator.List.(*rulesWrapper); ok && lr.ListRules.Unique {
			for i := 0; i < len(sliceDataX); i++ {
				for j := i + 1; j < len(sliceDataX); j++ {
					if reflect.DeepEqual(sliceDataX[i], sliceDataX[j]) {
						if validator.CustomMsg.OnUnique != nil {
							return nil, buildMessage(*validator.CustomMsg.OnUnique, MessageMeta{Field: &field})
						}
						return nil, buildErrorMessage(field, "values must be unique")
					}
				}
			}
		}
		return sliceDataX, nil
	}

	if validator.File {
		//this will return FileRequest
		return data, nil
	}

	// Format checks below may convert `result` (e.g. UUID → uuid.UUID) but
	// never return early, so Min/Max keep applying to the original value.
	result := data

	if validator.RegexString != "" {
		if dataType != reflect.String {
			if validator.CustomMsg.OnRegexString != nil {
				return nil, buildMessage(*validator.CustomMsg.OnRegexString, MessageMeta{Field: &field})
			}
			return nil, buildErrorMessage(field, "should be string")
		}
		regex, err := regexp.Compile(validator.RegexString)
		if err != nil {
			return nil, err
		}
		strData, _ := asString(data)
		if !regex.MatchString(strData) {
			if validator.CustomMsg.OnRegexString != nil {
				return nil, buildMessage(*validator.CustomMsg.OnRegexString, MessageMeta{Field: &field})
			}
			return nil, buildErrorMessage(field, "is not valid regex")
		}
	}

	if validator.Enum != nil {
		if err := validateEnumValue(data, dataType, validator, field); err != nil {
			return nil, err
		}
	}

	if validator.UUIDToString {
		validator.UUID = true
	}

	if validator.UUID {
		errMsg := buildErrorMessage(field, "is not valid uuid")
		stringUuid, ok := asString(data)
		if !ok {
			return nil, errMsg
		}
		dataUuid, err := uuid.Parse(stringUuid)
		if err != nil {
			return nil, errMsg
		}
		if validator.UUIDToString {
			result = stringUuid
		} else {
			result = dataUuid
		}
	}

	if validator.Email {
		strData, ok := asString(data)
		if !ok || !isEmail(strData) {
			return nil, buildErrorMessage(field, "is not valid email")
		}
	}

	if validator.IPV4 {
		errMsg := buildErrorMessage(field, "is not valid IP")
		stringIp, ok := asString(data)
		if !ok {
			return nil, errMsg
		}
		if !isIPv4Valid(stringIp) {
			return nil, errMsg
		}
		result = stringIp
	}

	if validator.IPV4Network {
		errMsg := buildErrorMessage(field, "is not valid IP Network")
		stringIp, ok := asString(data)
		if !ok {
			return nil, errMsg
		}
		if !isIPv4NetworkValid(stringIp) {
			return nil, errMsg
		}
		result = stringIp
	}

	if validator.IPv4OptionalPrefix {
		errMsg := buildErrorMessage(field, "is not valid IP")
		stringIp, ok := asString(data)
		if !ok {
			return nil, errMsg
		}
		splitIp := strings.Split(stringIp, "/")
		if len(splitIp) > 2 {
			return nil, errMsg
		}
		if !isIPv4Valid(splitIp[0]) {
			return nil, errMsg
		}
		if len(splitIp) == 2 {
			prefix, err := strconv.Atoi(splitIp[1])
			if err != nil {
				return nil, errMsg
			}
			if prefix < 0 || prefix > 32 {
				return nil, errMsg
			}
		}
		result = stringIp
	}

	// legacy ListObject fallback occurs via early list handling
	if validator.AnonymousObject || validator.Object != nil {
		res, err := toMapStringInterface(data)
		if err != nil {
			return nil, buildErrorMessage(field, "is not valid object")
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
		var actualLength interface{}
		err := buildErrorMessagef(field, "should be or greater than %v", *validator.Min)
		if reflect.String == dataType {
			strData, _ := asString(data)
			if total := utf8.RuneCountInString(strData); int64(total) < *validator.Min {
				isErr = true
				actualLength = int64(total)
			}
		} else if isIntegerFamily(dataType) {
			switch {
			case dataType == reflect.Float32 || dataType == reflect.Float64:
				// Compare floats without int64 truncation so 3.9 does not pass Max=3
				num := reflect.ValueOf(data).Float()
				if num < float64(*validator.Min) {
					isErr = true
					actualLength = int64(num)
				}
			case isUintKind(dataType):
				// Compare in uint64 space: int64(v.Uint()) wraps huge values negative
				u := reflect.ValueOf(data).Uint()
				if *validator.Min >= 0 && u < uint64(*validator.Min) {
					isErr = true
					actualLength = u
				}
			default:
				num := extractInteger(data)
				if num < *validator.Min {
					isErr = true
					actualLength = num
				}
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
					ActualLength:      anyPtr(actualLength),
				})
			}
			return nil, err
		}
	}

	if validator.Max != nil && data != nil {
		var isErr bool
		var actualLength interface{}
		err := buildErrorMessagef(field, "should be or lower than %v", *validator.Max)
		if reflect.String == dataType {
			strData, _ := asString(data)
			if total := utf8.RuneCountInString(strData); int64(total) > *validator.Max {
				isErr = true
				actualLength = int64(total)
			}
		} else if isIntegerFamily(dataType) {
			switch {
			case dataType == reflect.Float32 || dataType == reflect.Float64:
				// Compare floats without int64 truncation so 3.9 does not pass Max=3
				num := reflect.ValueOf(data).Float()
				if num > float64(*validator.Max) {
					isErr = true
					actualLength = int64(num)
				}
			case isUintKind(dataType):
				// Compare in uint64 space: int64(v.Uint()) wraps huge values negative
				u := reflect.ValueOf(data).Uint()
				if *validator.Max < 0 || u > uint64(*validator.Max) {
					isErr = true
					actualLength = u
				}
			default:
				num := extractInteger(data)
				if num > *validator.Max {
					isErr = true
					actualLength = num
				}
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
					ActualLength:      anyPtr(actualLength),
				})
			}
			return nil, err
		}
	}

	return result, nil
}

// validateEnumValue checks `data` against the rule's Enum configuration.
// A misconfigured enum (nil or non-slice Items) is a rule error — it used to
// silently disable the enum check (or panic on nil Items).
func validateEnumValue(data interface{}, dataType reflect.Kind, validator Rules, field string) error {
	// build enum error message with custom or default text
	buildEnumErrorMessage := func(enumValues interface{}, enumType reflect.Type, actualType reflect.Kind) error {
		if validator.CustomMsg.OnEnumValueNotMatch != nil {
			expectedType := enumType.Elem().Kind()
			actualValue := fmt.Sprintf("%v", data)
			enumValuesStr := fmt.Sprintf("%v", enumValues)
			return buildMessage(*validator.CustomMsg.OnEnumValueNotMatch, MessageMeta{
				Field:        &field,
				ExpectedType: &expectedType,
				ActualType:   &actualType,
				ActualValue:  &actualValue,
				EnumValues:   &enumValuesStr,
			})
		}
		return buildErrorMessagef(field, "value is not in enum list%v", enumValues)
	}

	enumType := reflect.TypeOf(validator.Enum.Items)
	if enumType == nil || enumType.Kind() != reflect.Slice {
		return buildErrorMessage(field, "has an invalid enum configuration: items must be a slice")
	}
	enumValue := reflect.ValueOf(validator.Enum.Items)
	elemKind := enumType.Elem().Kind()

	// Integer-family enums accept any integer-family value regardless of
	// concrete or named type, from any source (map, JSON, form). Values
	// are compared numerically so int64(2) matches IntEnum(1,2,3).
	if isIntegerFamily(elemKind) && isIntegerFamily(dataType) {
		dataNum := numericAsFloat64(data)
		isFloatEnum := elemKind == reflect.Float32 || elemKind == reflect.Float64
		// an integer enum cannot match a fractional value
		if !isFloatEnum && dataNum != float64(int64(dataNum)) {
			return buildErrorMessage(field, "should be '"+elemKind.String()+"'")
		}
		var matched bool
		for i := 0; i < enumValue.Len(); i++ {
			if numericAsFloat64(enumValue.Index(i).Interface()) == dataNum {
				matched = true
				break
			}
		}
		if !matched {
			return buildEnumErrorMessage(validator.Enum.Items, enumType, dataType)
		}
		return nil
	}

	if dataType != elemKind {
		// Use custom message for type mismatch if available, otherwise use custom enum message or default
		if validator.CustomMsg.OnTypeNotMatch != nil {
			expectedType := elemKind
			return buildMessage(*validator.CustomMsg.OnTypeNotMatch, MessageMeta{
				Field:        &field,
				ExpectedType: &expectedType,
				ActualType:   &dataType,
			})
		}
		return buildEnumErrorMessage(validator.Enum.Items, enumType, dataType)
	}

	switch dataType {
	case reflect.String:
		var values []string
		for i := 0; i < enumValue.Len(); i++ {
			values = append(values, enumValue.Index(i).String())
		}
		strData, _ := asString(data)
		if !valueInList[string](values, strData, isEqualString) {
			return buildEnumErrorMessage(values, enumType, dataType)
		}
	default:
		return buildErrorMessagef(field, "enum is not supported for type '%s'", dataType.String())
	}
	return nil
}

func SetTotal(total int64) *int64 {
	return &total
}

func SetMessage(msg string) *string { return &msg }

// anyPtr boxes a value into *interface{} for MessageMeta length fields that may
// hold int64 or uint64, so huge uint64 values format without int64 overflow.
func anyPtr(v interface{}) *interface{} { return &v }

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
	decoder := json.NewDecoder(ioData)
	// UseNumber keeps integer literals beyond float64's exact range intact;
	// normalizeJSONNumbers then restores plain Go types.
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}
	normalizeJSONNumbers(res)
	return res, nil
}

func extractInteger(data interface{}) int64 {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return int64(v.Float())
	}
	return 0
}

// asString extracts a string from data that may hold a named string type
// (e.g. type ID string). A hard data.(string) assertion panics on named
// types even though their reflect.Kind matches.
func asString(data interface{}) (string, bool) {
	if s, ok := data.(string); ok {
		return s, true
	}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.String {
		return v.String(), true
	}
	return "", false
}

// numericAsFloat64 converts any integer-family value (including named types)
// to float64 for comparison purposes.
func numericAsFloat64(data interface{}) float64 {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	}
	return 0
}

// isUintKind reports whether k is an unsigned integer kind. Unsigned values
// need their own comparison path: int64(v.Uint()) wraps values above
// math.MaxInt64 to negative, silently defeating Min/Max checks.
func isUintKind(k reflect.Kind) bool {
	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}

func getAllKeys(data map[string]interface{}) (allKeysInMap []string) {
	for key, _ := range data {
		allKeysInMap = append(allKeysInMap, key)
	}
	return
}

// maxExactFloatInt64 is the largest magnitude int64 that float64 represents
// exactly (2^53). Integer values beyond it lose precision through float64.
const maxExactFloatInt64 = int64(1) << 53

// sortedKeys returns the rule keys in deterministic order so validation — and
// therefore "the first error" — no longer depends on Go's randomized map
// iteration.
func sortedKeys(m map[string]Rules) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// checkStrictKeys rejects payload keys that have no matching rule at this
// object level.
func checkStrictKeys(wrapper RulesWrapper, data map[string]interface{}) error {
	allowed := make([]string, 0, len(wrapper.getRules()))
	for k := range wrapper.getRules() {
		allowed = append(allowed, k)
	}
	keys := getAllKeys(data)
	sort.Strings(keys)
	for _, k := range keys {
		if !isDataInList(k, allowed) {
			return fmt.Errorf("'%s' is not allowed key", k)
		}
	}
	return nil
}

// normalizeJSONNumbers converts json.Number values (produced by decoders with
// UseNumber) back into plain Go numbers. Integers within float64's exact range
// become float64 — matching the legacy encoding/json behavior — while larger
// integers stay int64 so IDs like 9007199254740993 are never corrupted.
func normalizeJSONNumbers(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		for k, val := range t {
			t[k] = normalizeJSONNumbers(val)
		}
		return t
	case []interface{}:
		for i, val := range t {
			t[i] = normalizeJSONNumbers(val)
		}
		return t
	case json.Number:
		if i, err := strconv.ParseInt(string(t), 10, 64); err == nil {
			if i > maxExactFloatInt64 || i < -maxExactFloatInt64 {
				return i
			}
			return float64(i)
		}
		if f, err := strconv.ParseFloat(string(t), 64); err == nil {
			return f
		}
		return string(t)
	}
	return v
}

func isDataInList[T validatorType](key T, data []T) (result bool) {
	for _, val := range data {
		if val == key {
			return true
		}
	}
	return
}

// coercesNumbers reports whether values from this source arrive pre-decoded as
// JSON-like primitives (numbers as float64, bools as bool) so integer-family
// type coercion applies. Multipart/urlencoded form values are normalized the
// same way in LoadFormHttp, so they qualify too.
func coercesNumbers(dataFrom loadFromType) bool {
	return dataFrom == fromHttpJson || dataFrom == fromJSONEncoder || dataFrom == fromHttpMultipartForm
}

// integerCoercion reports whether an integer rule should accept the given
// actual kind. JSON/form sources decode all numbers as float64, so integer
// rules there tolerate any integer-family kind (including float). Map sources
// keep exact types, so only integer kinds (int/int64/uint/...) cross-coerce —
// a float for an int rule is still rejected.
func integerCoercion(dataFrom loadFromType, expected, actual reflect.Kind) bool {
	if coercesNumbers(dataFrom) && isIntegerFamily(expected) && isIntegerFamily(actual) {
		return true
	}
	return dataFrom == fromMapString && isIntegerKind(expected) && isIntegerKind(actual)
}

// coerceFormValue normalizes a raw form string into the JSON-like primitive the
// validator expects for the declared kind: numbers become float64 (matching the
// JSON decoder), booleans become bool. Integers beyond float64's exact range
// stay int64 so they do not lose precision. On parse failure the raw string is
// returned so the normal type check reports a clear "should be 'int'" error.
func coerceFormValue(value string, kind reflect.Kind) interface{} {
	switch {
	case kind == reflect.Bool:
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	case isIntegerFamily(kind):
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			if i > maxExactFloatInt64 || i < -maxExactFloatInt64 {
				return i
			}
			return float64(i)
		}
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return value
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

// isIntegerKind is the integer-family predicate WITHOUT floats: true only for
// the signed/unsigned integer kinds. Used for map-source cross-coercion where
// floats must keep their exact (float64) type check.
func isIntegerKind(dataType reflect.Kind) bool {
	switch dataType {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}
