package map_validator

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func NewValidateBuilder() *ruleState {
	return &ruleState{}
}

func (state *ruleState) SetRules(validations map[string]Rules) *dataState {
	state.rules = validations
	return &dataState{
		ruleState: state,
		data:      nil,
	}
}

func (state *dataState) Load(data map[string]interface{}) *finalOperation {
	state.data = data
	return &finalOperation{state}
}

func (state *dataState) LoadJsonHttp(r *http.Request) (*finalOperation, error) {
	var mapData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&mapData)
	if err != nil {
		if err.Error() == "EOF" {
			return nil, ErrNoData
		}
		return nil, ErrInvalidFormat
	}
	return &finalOperation{state}, nil
}

func (state *dataState) LoadFormHttp(r *http.Request) (*finalOperation, error) {
	var mapData map[string]interface{}
	allowType := []reflect.Kind{reflect.String, reflect.Int, reflect.Bool}
	for key, rule := range state.rules {
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
	return &finalOperation{state}, nil
}

func (state *finalOperation) RunValidate() error {
	for key, validationData := range state.rules {
		_, err := validate(key, state.data, validationData)
		if err != nil {
			return err
		}
	}
	return nil
}
