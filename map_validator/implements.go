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

func (state *ruleState) SetRules(validations map[string]Rules) *dataState {
	if len(validations) == 0 {
		panic("you need to set roles")
	}

	var tempExt []ExtensionType
	state.rules = validations

	for _, ex := range state.extension {
		ex.SetRoles(&state.rules)
		tempExt = append(tempExt, ex)
	}
	if len(tempExt) > 0 {
		state.extension = tempExt
	}
	return &dataState{
		rules:              &state.rules,
		extension:          state.extension,
		strictAllowedValue: state.strictAllowedValue,
	}
}

func (state *ruleState) StrictKeys() *ruleState {
	state.strictAllowedValue = true
	return state
}

func (state *ruleState) AddExtension(extension ExtensionType) *ruleState {
	state.extension = append(state.extension, extension)
	return state
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
	if state.strictAllowedValue {
		if err := state.checkStrictKeys(data); err != nil {
			return nil, err
		}
	}
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
	err := json.NewDecoder(r.Body).Decode(&mapData)
	if err != nil {
		if err.Error() != "EOF" {
			return nil, ErrInvalidJsonFormat
		}
		mapData = make(map[string]interface{})
	}
	if state.strictAllowedValue {
		if err := state.checkStrictKeys(mapData); err != nil {
			return nil, err
		}
	}
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
	if state.strictAllowedValue {
		if err := state.checkStrictKeys(mapData); err != nil {
			return nil, err
		}
	}
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
	for key, rule := range *state.rules {
		data, err := validateRecursive(key, state.data, rule, state.loadedFrom)
		if err != nil {
			return nil, err
		}
		//data, err := validate(key, state.data, rule, state.loadedFrom)
		//if err != nil {
		//	return nil, err
		//}
		//if rule.Object != nil {
		//	for keyX, ruleX := range *rule.Object {
		//		_, err := validate(keyX, data.(map[string]interface{}), ruleX, fromJSONEncoder)
		//		if err != nil {
		//			//if rule.CustomMsg != nil #TODO: custom message for nested object
		//			err = errors.New(fmt.Sprintf("error on object '%s' : '%s'", key, err))
		//			return nil, err
		//		}
		//	}
		//}
		if data != nil {
			filledFields = append(filledFields, key)
		} else {
			nullFields = append(nullFields, key)
		}
	}
	extraData := &ExtraOperationData{
		rules:        state.rules,
		loadedFrom:   &state.loadedFrom,
		data:         &state.data,
		filledFields: filledFields,
		nullFields:   nullFields,
	}
	for _, ex := range state.extension {
		err := ex.SetExtraData(extraData).AfterValidation(&state.data)
		if err != nil {
			return nil, err
		}
	}
	return extraData, nil
}
