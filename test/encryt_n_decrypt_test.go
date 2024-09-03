package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"reflect"
	"testing"
)

func TestEncryptCaseOne(t *testing.T) {
	var obj interface{}
	obj = map[string]interface{}{"nama": "arian", "kelas": float64(1)}
	payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS bro", "umur": 1, "menikah": true, "list_data": obj}
	type Data struct {
		JK       string      `map_validator:"jenis_kelamin" json:"jenis_kelamin"`
		Hoby     *string     `map_validator:"hoby" json:"hoby"`
		Menikah  bool        `map_validator:"menikah" json:"menikah"`
		ListData interface{} `map_validator:"list_data" json:"list_data"`
	}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"hoby":          {Type: reflect.String, Null: true},
			"menikah":       {Type: reflect.Bool, Null: false},
			"list_data":     {AnonymousObject: true},
		},
	}

	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	extraCheck, err := check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}

	testBind := &Data{}
	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %v", testBind.JK)
	}
	err = extraCheck.Bind(testBind)
	if err != nil {
		t.Errorf("Expected have an error, but you got no error : %s", err)
	}

	keys := getAllKeys(obj.(map[string]interface{}))
	keysRes := getAllKeys(testBind.ListData.(map[string]interface{}))
	if len(keys) != len(keysRes) {
		t.Errorf("Expected : %v But you got : %v", keys, keysRes)
	}
	for _, key := range keys {
		val := obj.(map[string]interface{})[key]
		valRes := testBind.ListData.(map[string]interface{})[key]
		if val != valRes {
			t.Errorf("Expected : %v But you got : %v", keys, keysRes)
		}
	}
}
