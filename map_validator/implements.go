package map_validator

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

func NewValidateBuilder() *ruleState {
	return &ruleState{}
}

func (state *ruleState) SetRules(validations map[string]Rules) *dataState {
	state.rules = validations
	return &dataState{
		rules: &state.rules,
	}
}

func (state *dataState) Load(data map[string]interface{}) *finalOperation {
	return &finalOperation{
		rules:      state.rules,
		loadedFrom: fromMapString,
		data:       data,
	}
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
	return &finalOperation{
		rules:      state.rules,
		loadedFrom: fromHttpJson,
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
	return &finalOperation{
		rules:      state.rules,
		loadedFrom: fromHttpMultipartForm,
		data:       mapData,
	}, nil
}

func (state *finalOperation) RunValidate() (*extraOperation, error) {
	if state == nil || state.data == nil {
		return nil, errors.New("no data to Validate because last progress is error")
	}
	for key, validationData := range *state.rules {
		_, err := validate(key, state.data, validationData, state.loadedFrom)
		if err != nil {
			return nil, err
		}
	}
	return &extraOperation{
		rules:      state.rules,
		loadedFrom: &state.loadedFrom,
		data:       &state.data,
	}, nil
}

func (state *extraOperation) Bind(i interface{}) error {
	var data map[string]interface{}
	data = *state.data // this for memory allocation purpose
	if state == nil || state.data == nil {
		return errors.New("no data to Bind because last progress is error")
	}
	allKeysInMap := getAllKeys(data)
	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		panic("need struct pointer!")
	}

	el := val.Elem()
	t := val.Elem().Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("map_validator")
		if !isDataInList[string](tag, allKeysInMap) {
			continue
		}

		if tag == "" || !field.IsExported() || data[tag] == nil {
			continue
		}

		if field.Type.Kind() == reflect.Ptr && reflect.TypeOf(data[tag]).Kind() == field.Type.Elem().Kind() {
			err := convertValue(data[tag], field.Type.Elem().Kind(), el.Field(i), true)
			if err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.TypeOf(data[tag]).Kind() &&
			field.Type.Kind() != reflect.Struct {
			err := convertValue(data[tag], field.Type.Kind(), el.Field(i), false)
			if err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Interface {
			if reflect.TypeOf(data[tag]).Kind() == reflect.Map {
				err := convertValue(data[tag], field.Type.Kind(), el.Field(i), false)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
