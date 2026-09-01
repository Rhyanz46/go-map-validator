package map_validator

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
)

func NewValidateBuilder() *ruleState {
	return &ruleState{}
}

func (state *ruleState) SetRules(validations RulesWrapper) *dataState {
	var tempExt []ExtensionType
	state.rules = validations

	for _, ex := range state.extension {
		if state.rules != nil {
			ex.SetRoles(state.rules)
			tempExt = append(tempExt, ex)
		}
	}
	if len(tempExt) > 0 {
		state.extension = tempExt
	}
	return &dataState{
		rules:              state.rules,
		extension:          state.extension,
		strictAllowedValue: state.strictAllowedValue,
	}
}

func (state *ruleState) AddExtension(extension ExtensionType) *ruleState {
	state.extension = append(state.extension, extension)
	return state
}

//	func (state *dataState) checkStrictKeys(data map[string]interface{}) error {
//		var allowedKeys []string
//		keys := getAllKeys(data)
//		for key, _ := range state.rules.Rules {
//			allowedKeys = append(allowedKeys, key)
//		}
//		for _, key := range keys {
//			if !isDataInList(key, allowedKeys) {
//				return errors.New(fmt.Sprintf("'%s' is not allowed key", key))
//			}
//		}
//		return nil
//	}
func (state *dataState) Load(data map[string]interface{}) (*finalOperation, error) {
	if state == nil || state.rules == nil || len(state.rules.getRules()) == 0 {
		return nil, ErrNoRules
	}
	if data == nil {
		// A nil payload behaves like an empty one so required-field errors
		// surface from validation instead of a misleading "no data" error.
		data = make(map[string]interface{})
	}
	//if state.strictAllowedValue {
	//	if err := state.checkStrictKeys(data); err != nil {
	//		return nil, err
	//	}
	//}
	for _, ex := range state.extension {
		err := ex.BeforeLoad(&data)
		if err != nil {
			return nil, err
		}
	}
	for _, ex := range state.extension {
		err := ex.AfterLoad(&data)
		if err != nil {
			return nil, err
		}
	}
	return &finalOperation{
		rules:      state.rules,
		loadedFrom: fromMapString,
		extension:  state.extension,
		data:       data,
	}, nil
}

func (state *dataState) LoadJsonHttp(r *http.Request) (*finalOperation, error) {
	if state == nil {
		return nil, errors.New("no data to Load because last progress is error")
	}
	if state.rules == nil || len(state.rules.getRules()) == 0 {
		return nil, ErrNoRules
	}
	if r == nil {
		return nil, errors.New("no data to Load")
	}
	for _, ex := range state.extension {
		err := ex.BeforeLoad(r)
		if err != nil {
			return nil, err
		}
	}
	var mapData map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	// UseNumber keeps integer literals beyond float64's exact range intact;
	// normalizeJSONNumbers then restores plain Go types.
	decoder.UseNumber()
	err := decoder.Decode(&mapData)
	if err != nil {
		if errors.Is(err, io.EOF) {
			mapData = make(map[string]interface{})
		} else {
			return nil, ErrInvalidJsonFormat
		}
	}
	if mapData == nil {
		// a body of literal `null` decodes to a nil map — treat it as empty
		mapData = make(map[string]interface{})
	}
	normalizeJSONNumbers(mapData)
	//if state.strictAllowedValue {
	//	if err := state.checkStrictKeys(mapData); err != nil {
	//		return nil, err
	//	}
	//}
	for _, ex := range state.extension {
		err := ex.AfterLoad(&mapData)
		if err != nil {
			return nil, err
		}
	}
	return &finalOperation{
		rules:      state.rules,
		loadedFrom: fromHttpJson,
		extension:  state.extension,
		data:       mapData,
	}, nil
}

func (state *dataState) LoadFormHttp(r *http.Request) (*finalOperation, error) {
	if state == nil {
		return nil, errors.New("no data to Load because last progress is error")
	}
	if state.rules == nil || len(state.rules.getRules()) == 0 {
		return nil, ErrNoRules
	}
	if r == nil {
		return nil, errors.New("no data to Load")
	}
	for _, ex := range state.extension {
		err := ex.BeforeLoad(r)
		if err != nil {
			return nil, err
		}
	}
	mapData := map[string]interface{}{}
	for key, rule := range state.rules.getRules() {
		if rule.File {
			file, fileInfo, err := r.FormFile(key)
			if err != nil || file == nil {
				mapData[key] = nil
			} else {
				mapData[key] = FileRequest{File: file, FileInfo: fileInfo}
			}
			continue
		}
		// Compound structures (Object, ListObject, List) have no flat HTTP
		// form representation — reject early with ErrUnsupportType.
		if rule.Object != nil || rule.ListObject != nil || rule.List != nil || rule.AnonymousObject {
			return nil, ErrUnsupportType
		}
		// Flag-based rules (Email, UUID, IPv4, Regex, Enum-only) carry a
		// zero-value Type (reflect.Invalid) because they validate string
		// content rather than a reflect kind — accept them as plain strings.
		isFlagBasedRule := rule.Type == reflect.Invalid
		if !isFlagBasedRule && rule.Type != reflect.String && rule.Type != reflect.Bool && !isIntegerFamily(rule.Type) {
			return nil, ErrUnsupportType
		}
		value := r.FormValue(key)
		if value == "" {
			mapData[key] = nil
		} else {
			mapData[key] = coerceFormValue(value, rule.Type)
		}
	}
	//if state.strictAllowedValue {
	//	if err := state.checkStrictKeys(mapData); err != nil {
	//		return nil, err
	//	}
	//}
	for _, ex := range state.extension {
		err := ex.AfterLoad(&mapData)
		if err != nil {
			return nil, err
		}
	}
	return &finalOperation{
		rules:      state.rules,
		loadedFrom: fromHttpMultipartForm,
		extension:  state.extension,
		data:       mapData,
	}, nil
}

func (state *finalOperation) RunValidate() (*ExtraOperationData, error) {
	initChain := newChainer().SetKey(chainKey)
	if state == nil || state.data == nil {
		return nil, errors.New("no data to Validate because last progress is error")
	}
	var filledFields []string
	var nullFields []string
	for _, ex := range state.extension {
		err := ex.BeforeValidation(&state.data)
		if err != nil {
			return nil, err
		}
	}
	topState := newWrapperRunState()
	if state.rules.getSetting().Strict {
		if err := checkStrictKeys(state.rules, state.data); err != nil {
			return nil, err
		}
	}
	rulesMap := state.rules.getRules()
	for _, key := range sortedKeys(rulesMap) {
		data, err := validateRecursive(initChain, state.rules, topState, key, state.data, rulesMap[key], state.loadedFrom)
		if err != nil {
			return nil, err
		}
		if data != nil {
			filledFields = append(filledFields, key)
		} else {
			nullFields = append(nullFields, key)
		}
	}

	chainRes := initChain.GetResult()
	err := chainRes.RunManipulator()
	if err != nil {
		return nil, err
	}

	chainRes.RunUniqueChecker()
	for _, err = range chainRes.GetErrors() {
		if err != nil {
			return nil, err
		}
	}

	manipulatedData := chainRes.ToMap()
	extraData := &ExtraOperationData{
		rules:        state.rules,
		loadedFrom:   &state.loadedFrom,
		data:         &manipulatedData,
		filledFields: filledFields,
		nullFields:   nullFields,
	}
	for _, ex := range state.extension {
		err := ex.SetExtraData(extraData).AfterValidation(&manipulatedData)
		if err != nil {
			return nil, err
		}
	}
	return extraData, nil
}

func (state *Rules) isList() bool {
	return state.List != nil
}
