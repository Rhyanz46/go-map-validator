package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"reflect"
	"testing"
)

func TestMultipleValidation(t *testing.T) {
	type Data struct {
		JK      string `map_validator:"jenis_kelamin" json:"jenis_kelamin"`
		Hoby    string `map_validator:"hoby" json:"hoby"`
		Menikah bool   `map_validator:"menikah" json:"menikah"`
	}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"hoby":          {Type: reflect.String, Null: false},
			"menikah":       {Type: reflect.Bool, Null: false},
		},
	}
	validRoleOptionalMenikah := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"hoby":          {Type: reflect.String, Null: false},
			"menikah":       {Type: reflect.Bool, Null: true},
		},
	}
	payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS bro", "umur": 1, "menikah": true}
	notFullPayload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS bro", "umur": 1}
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
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
	}
	err = extraCheck.Bind(testBind)
	if err != nil {
		t.Errorf("Expected have an error, but you got no error : %s", err)
	}

	if testBind.JK != payload["jenis_kelamin"] {
		t.Errorf("Expected : %s But you got : %s", payload["jenis_kelamin"], testBind.JK)
	}

	check, err = map_validator.NewValidateBuilder().SetRules(map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"hoby":          {Type: reflect.Int, Null: false},
			"menikah":       {Type: reflect.Bool, Null: false},
		},
	}).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err == nil {
		t.Errorf("Expected have an error, but you got no error : %s", err)
	} else {
		expected := "the field 'hoby' should be 'int'"
		if err.Error() != expected {
			t.Errorf("Expected :%s. But you got : %s", expected, err)
		}
	}

	check, err = map_validator.NewValidateBuilder().SetRules(validRole).Load(notFullPayload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	extraCheck, err = check.RunValidate()
	expected := "we need 'menikah' field"
	if err.Error() != expected {
		t.Errorf("Expected : %s But you got : %s", expected, err)
	}

	testBind = &Data{}
	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
	}
	err = extraCheck.Bind(testBind)
	expected = "no data to Bind because last progress is error"
	if err != nil {
		if err.Error() != expected {
			t.Errorf("Expected : %s But you got : %s", expected, err)
		}
	}

	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
	}

	check, err = map_validator.NewValidateBuilder().SetRules(validRoleOptionalMenikah).Load(notFullPayload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	extraCheck, err = check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}

	testBind = &Data{}
	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
	}
	err = extraCheck.Bind(testBind)
	expected = "no data to Bind because last progress is error"
	if err != nil {
		if err.Error() != expected {
			t.Errorf("Expected : %s But you got : %s", expected, err)
		}
	}

	if testBind.JK != notFullPayload["jenis_kelamin"] {
		t.Errorf("Expected : %s But you got : %s", notFullPayload["jenis_kelamin"], testBind.JK)
	}
}

func TestPointerFieldBinding(t *testing.T) {
	payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS bro", "umur": 1, "menikah": true}
	type Data struct {
		JK      string  `map_validator:"jenis_kelamin" json:"jenis_kelamin"`
		Hoby    *string `map_validator:"hoby" json:"hoby"`
		Menikah bool    `map_validator:"menikah" json:"menikah"`
	}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"hoby":          {Type: reflect.String, Null: true},
			"menikah":       {Type: reflect.Bool, Null: false},
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

	if testBind.JK != payload["jenis_kelamin"] {
		t.Errorf("Expected : %s But you got : %s", payload["jenis_kelamin"], testBind.JK)
	}

	if *testBind.Hoby != payload["hoby"] {
		t.Errorf("Expected : %s But you got : %s", payload["hoby"], *testBind.Hoby)
	}
}

func TestInterfaceFieldBinding(t *testing.T) {
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

func TestFilledAndNullField(t *testing.T) {
	payload := map[string]interface{}{"nama": "arian", "umur": 1}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama": {Type: reflect.String},
			"hoby": {Type: reflect.String, Null: true},
			"umur": {Type: reflect.Int, Null: false},
		}}
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	extraCheck, err := check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	totalFilled := extraCheck.GetFilledField()
	totalNull := extraCheck.GetNullField()
	if len(totalFilled) != 2 {
		t.Errorf("Expected total data is 2, but we got %v data = (%v)", len(totalFilled), totalFilled)
	}
	if len(totalNull) != 1 {
		t.Errorf("Expected total data is 1, but we got %v data = (%v)", len(totalNull), totalNull)
	}
}

func TestGetMapData(t *testing.T) {
	payload := map[string]interface{}{"nama": "arian", "umur": 1}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama": {Type: reflect.String},
			"hoby": {Type: reflect.String, Null: true},
			"umur": {Type: reflect.Int, Null: false},
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
	expectedKeys := getAllKeys(payload)
	for _, key := range getAllKeys(extraCheck.GetData()) {
		if !isDataInList(key, expectedKeys) {
			t.Errorf("Key is not Expected : %v", key)
		}
	}
}

func TestStrict(t *testing.T) {
	payload := map[string]interface{}{"nama": "arian", "umur": 1, "favorite": "coklat"}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama": {Type: reflect.String},
			"hoby": {Type: reflect.String, Null: true},
			"umur": {Type: reflect.Int, Null: false},
		},
		Setting: map_validator.Setting{Strict: true},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	expected := "'favorite' is not allowed key"
	_, err = check.RunValidate()
	if err.Error() != expected {
		t.Errorf("Expected %s, but we got : %s", expected, err)
	}
}

func TestStrictTwo(t *testing.T) {
	payload := map[string]interface{}{"nama": "arian", "umur": 1}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama": {Type: reflect.String},
			"hoby": {Type: reflect.String, Null: true},
		},
		Setting: map_validator.Setting{Strict: true},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	expected := "'umur' is not allowed key"
	_, err = check.RunValidate()
	if err.Error() != expected {
		t.Errorf("Expected %s, but we got : %s", expected, err)
	}
}

func TestValidRegex(t *testing.T) {
	payload := map[string]interface{}{"hp": "+62567888", "email": "dev@ariansaputra.com"}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"hp":    {RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$`},
			"email": {RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`},
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
}

func TestInvalidRegex(t *testing.T) {
	payload := map[string]interface{}{"hp": "62567888", "email": "devariansaputra.com"}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"hp": {RegexString: `^\+(?:\d{2}[- ]?\d{6}|\d{11})$`},
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err == nil {
		t.Error("Expected error, but got no error :")
	}
	validRole = map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"email": {RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`},
		},
	}
	check, err = map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err == nil {
		t.Error("Expected error, but got no error :")
	}
}

func TestValidSlice(t *testing.T) {
	payload := map[string]interface{}{"hobby": []string{"reading", "football"}}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"hobby": {Type: reflect.Slice},
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	_, err = check.RunValidate()
	if err != nil {
		t.Error("Expected no error, but we got :", err)
	}

	//invalidPayload := map[string]interface{}{"hobby": 122}
	invalidPayload := map[string]interface{}{"aaa": 122}
	check, err = map_validator.NewValidateBuilder().SetRules(validRole).Load(invalidPayload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	expected := "we need 'hobby' field"
	_, err = check.RunValidate()
	if err.Error() != expected {
		t.Errorf("Expected %s, but we got : %s", expected, err)
	}

	invalidPayload = map[string]interface{}{"hobby": 122}
	check, err = map_validator.NewValidateBuilder().SetRules(validRole).Load(invalidPayload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	expected = "the field 'hobby' should be 'slice'"
	_, err = check.RunValidate()
	if err.Error() != expected {
		t.Errorf("Expected %s, but we got : %s", expected, err)
	}
	//fmt.Println(err)
}

func TestInvalidSlice(t *testing.T) {
	payload := map[string]interface{}{"hobby": []string{"reading", "football", "eat"}}
	validRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"hobby": {Type: reflect.Slice, Max: map_validator.SetTotal(2)},
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}
	expected := "the field 'hobby' should be or lower than 2"
	_, err = check.RunValidate()
	if err.Error() != expected {
		t.Errorf("Expected %s, but we got : %s", expected, err)
	}
}
