package map_validator

import (
	"encoding/json"
	"errors"
)

func (state *ExtraOperationData) Bind(i interface{}) error {
	var data map[string]interface{}
	if state == nil || state.data == nil {
		return errors.New("no data to Bind because last progress is error")
	}
	data = *state.data // this for memory allocation purpose
	//allKeysInMap := getAllKeys(data)
	//val := reflect.ValueOf(i)
	//if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
	//	panic("need struct pointer!")
	//}
	//
	//el := val.Elem()
	//t := val.Elem().Type()

	//for i := 0; i < t.NumField(); i++ {
	//	field := t.Field(i)
	//	tag := field.Tag.Get("map_validator")
	//	if !isDataInList[string](tag, allKeysInMap) {
	//		continue
	//	}
	//
	//	if tag == "" || !field.IsExported() || data[tag] == nil {
	//		continue
	//	}
	//
	//	if field.Type.Kind() == reflect.Ptr && reflect.TypeOf(data[tag]).Kind() == field.Type.Elem().Kind() {
	//		err := convertValue(data[tag], field.Type.Elem().Kind(), el.Field(i), true)
	//		if err != nil {
	//			return err
	//		}
	//	} else if field.Type.Kind() == reflect.TypeOf(data[tag]).Kind() &&
	//		field.Type.Kind() != reflect.Struct {
	//		err := convertValue(data[tag], field.Type.Kind(), el.Field(i), false)
	//		if err != nil {
	//			return err
	//		}
	//	} else if field.Type.Kind() == reflect.Interface {
	//		if reflect.TypeOf(data[tag]).Kind() == reflect.Map {
	//			err := convertValue(data[tag], field.Type.Kind(), el.Field(i), false)
	//			if err != nil {
	//				return err
	//			}
	//		}
	//	}
	//}

	jsonStringData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonStringData, i)
	if err != nil {
		return err
	}

	return nil
}

func (state *ExtraOperationData) GetFilledField() []string {
	if len(state.filledFields) > 0 {
		return state.filledFields
	}
	return []string{}
}

func (state *ExtraOperationData) GetNullField() []string {
	if len(state.nullFields) > 0 {
		return state.nullFields
	}
	return []string{}
}

func (state *ExtraOperationData) GetData() map[string]interface{} {
	return *state.data
}
