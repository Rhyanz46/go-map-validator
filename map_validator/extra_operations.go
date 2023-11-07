package map_validator

import (
	"errors"
	"reflect"
)

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

func (state *extraOperation) GetFilledField() []string {
	if len(state.filledFields) > 0 {
		return state.filledFields
	}
	return []string{}
}

func (state *extraOperation) GetNullField() []string {
	if len(state.nullFields) > 0 {
		return state.nullFields
	}
	return []string{}
}
