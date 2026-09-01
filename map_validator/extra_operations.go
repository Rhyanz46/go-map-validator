package map_validator

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

func (state *ExtraOperationData) Bind(i interface{}) error {
	var data map[string]interface{}
	if state == nil || state.data == nil {
		return errors.New("no data to Bind because last progress is error")
	}
	data = *state.data // this for memory allocation purpose

	// multipart.File is a stream and cannot survive a JSON round-trip —
	// carrying it around the marshal/unmarshal step keeps Bind working for
	// multipart forms instead of failing with an unmarshal error.
	files := map[string]interface{}{}
	clean := make(map[string]interface{}, len(data))
	for k, v := range data {
		if _, ok := v.(FileRequest); ok {
			files[k] = v
			continue
		}
		clean[k] = v
	}

	jsonStringData, err := json.Marshal(clean)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonStringData, i)
	if err != nil {
		return err
	}

	if len(files) > 0 {
		bindFileRequests(i, files)
	}

	return nil
}

// bindFileRequests injects FileRequest values into the struct fields whose
// json tag (or field name) matches the rule key. Only top-level fields are
// matched — multipart form fields are flat by nature.
func bindFileRequests(i interface{}, files map[string]interface{}) {
	el := reflect.ValueOf(i)
	if el.Kind() != reflect.Ptr || el.IsNil() {
		return
	}
	el = el.Elem()
	if el.Kind() != reflect.Struct {
		return
	}
	t := el.Type()
	for idx := 0; idx < t.NumField(); idx++ {
		field := t.Field(idx)
		if !field.IsExported() {
			continue
		}
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "-" {
			continue
		}
		if name == "" {
			name = field.Name
		}
		val, ok := files[name]
		if !ok {
			continue
		}
		fv := reflect.ValueOf(val)
		target := el.Field(idx)
		switch {
		case fv.Type().AssignableTo(target.Type()):
			target.Set(fv)
		case target.Kind() == reflect.Ptr && fv.Type().AssignableTo(target.Type().Elem()):
			p := reflect.New(target.Type().Elem())
			p.Elem().Set(fv)
			target.Set(p)
		}
	}
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
	if state == nil || state.data == nil {
		return map[string]interface{}{}
	}
	return *state.data
}
