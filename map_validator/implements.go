package map_validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

func NewValidateBuilder() *ruleState {
	return &ruleState{}
}

func (state *ruleState) SetRules(validations map[string]Rules) *optionalRolesState {
	if len(validations) == 0 {
		panic("you need to set roles")
	}
	return &optionalRolesState{
		data: &dataState{
			rules: &validations,
		},
	}
}

func (state *ruleState) SetRule(validation Rules) *oneDataState {
	return &oneDataState{
		rule: &validation,
	}
}

func (state *optionalRolesState) StrictKeys() *optionalRolesState {
	state.strictAllowedValue = true
	return state
}

func (state *optionalRolesState) Next() *dataState {
	return &dataState{
		rules: state.data.rules,
		optional: &optionalRolesState{
			strictAllowedValue: state.strictAllowedValue,
		},
	}
}

func (state *oneDataState) Load(data interface{}) (*finalOperation, error) {
	return &finalOperation{
		rules:      nil,
		rule:       state.rule,
		loadedFrom: fromOneValue,
		oneValue:   true,
		data:       map[string]interface{}{"data": data},
	}, nil
}

func (state *oneDataState) LoadFromHttp(r *http.Request, resWriter http.ResponseWriter) (*finalOperation, error) {
	return nil, nil
}

func (state *dataState) checkStrictKeys(data map[string]interface{}) error {
	var allowedKeys []string
	keys := getAllKeys(data)
	for key, _ := range *state.rules {
		allowedKeys = append(allowedKeys, key)
	}
	for _, key := range keys {
		if !isDataInList(key, allowedKeys) {
			return errors.New(fmt.Sprintf("'%s' is not allowed key", key))
		}
	}
	return nil
}
func (state *dataState) Load(data map[string]interface{}) (*finalOperation, error) {
	if state.optional != nil && state.optional.strictAllowedValue {
		if err := state.checkStrictKeys(data); err != nil {
			return nil, err
		}
	}
	return &finalOperation{
		rules:      state.rules,
		loadedFrom: fromMapString,
		oneValue:   false,
		rule:       nil,
		data:       data,
	}, nil
}

func (state *dataState) LoadJsonHttp(r *http.Request) (*finalOperation, error) {
	if state == nil {
		return nil, errors.New("no data to Load because last progress is error")
	}
	if r == nil {
		return nil, errors.New("no data to Load")
	}
	var mapData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&mapData)
	if err != nil {
		if err.Error() == "EOF" {
			return nil, ErrNoData
		}
		return nil, ErrInvalidFormat
	}
	if state.optional != nil && state.optional.strictAllowedValue {
		if err := state.checkStrictKeys(mapData); err != nil {
			return nil, err
		}
	}
	return &finalOperation{
		rules:      state.rules,
		loadedFrom: fromHttpJson,
		oneValue:   false,
		rule:       nil,
		data:       mapData,
	}, nil
}

func (state *dataState) LoadMultiPartFormHttp(r *http.Request, resWriter http.ResponseWriter) (*finalOperation, error) {
	if state == nil {
		return nil, errors.New("no data to Load because last progress is error")
	}
	if r == nil {
		return nil, errors.New("no data to Load")
	}
	var mapData map[string]interface{}
	allowType := []reflect.Kind{reflect.String, reflect.Int, reflect.Bool}
	for key, rule := range *state.rules {
		var isAllowType bool
		if rule.File {
			file, fileInfo, err := r.FormFile(key)
			if err != nil {
				mapData[key] = nil
			}
			if file == nil {
				mapData[key] = nil
			} else {
				mapData[key] = FileRequest{File: file, FileInfo: fileInfo}
			}
		} else {
			for _, allowItem := range allowType {
				if rule.Type == allowItem {
					isAllowType = true
					break
				}
			}
			if !isAllowType {
				return nil, ErrUnsupportType
			}
			value := r.FormValue(key)
			if value == "" {
				mapData[key] = nil
			} else {
				mapData[key] = value
			}
		}
	}
	if state.optional != nil && state.optional.strictAllowedValue {
		if err := state.checkStrictKeys(mapData); err != nil {
			return nil, err
		}
	}
	return &finalOperation{
		rules:      state.rules,
		loadedFrom: fromHttpMultipartForm,
		oneValue:   false,
		rule:       nil,
		data:       mapData,
	}, nil
}

func (state *finalOperation) RunValidate() (*extraOperation, error) {
	if state == nil || state.data == nil {
		return nil, errors.New("no data to Validate because last progress is error")
	}
	if state.oneValue {
		_, err := validate("data", state.data, *state.rule, state.loadedFrom)
		if err != nil {
			return nil, err
		}
		return &extraOperation{
			rule:       state.rule,
			rules:      nil,
			oneValue:   true,
			loadedFrom: &state.loadedFrom,
			data:       &state.data,
		}, nil
	}
	var filledFields []string
	var nullFields []string
	for key, validationData := range *state.rules {
		data, err := validate(key, state.data, validationData, state.loadedFrom)
		if err != nil {
			return nil, err
		}
		if data != nil {
			filledFields = append(filledFields, key)
		} else {
			nullFields = append(nullFields, key)
		}
	}
	return &extraOperation{
		rules:        state.rules,
		loadedFrom:   &state.loadedFrom,
		data:         &state.data,
		oneValue:     false,
		filledFields: filledFields,
		nullFields:   nullFields,
	}, nil
}
