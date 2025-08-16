package test

import (
	"bytes"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
	"github.com/google/uuid"
)

func TestInvalidNestedValidation(t *testing.T) {
	type Orang struct {
		Nama    string `map_validator:"nama" json:"nama"`
		Umur    int    `map_validator:"umur" json:"umur"`
		Hoby    string `map_validator:"hoby" json:"hoby"`
		JK      string `map_validator:"jenis_kelamin" json:"jenis_kelamin"`
		Menikah bool   `map_validator:"menikah" json:"menikah"`
		Anak    *Orang `map_validator:"anak" json:"anak"`
	}
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama":          {Type: reflect.String},
			"umur":          {Type: reflect.Int},
			"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
		},
	}
	role.Rules["anak"] = map_validator.Rules{Object: &role, Null: true}
	child1 := map[string]interface{}{
		"nama":          "Messi",
		"jenis_kelamin": "laki-laki",
		"hoby":          "football",
		"umur":          5,
		"menikah":       "aaa",
	}
	payload := map[string]interface{}{
		"nama":          "Arian Saputra",
		"jenis_kelamin": "laki-laki",
		"hoby":          "Main PS bro",
		"umur":          33,
		"menikah":       true,
		"anak":          child1,
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	extraCheck, err := check.RunValidate()
	expected := "the field 'menikah' should be 'bool'"
	if err.Error() != expected {
		t.Errorf("Expected :%s. But you got : %s", expected, err)
	}

	testBind := &Orang{}
	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
		return
	}
	err = extraCheck.Bind(testBind)
	expected = "no data to Bind because last progress is error"
	if err.Error() != expected {
		t.Errorf("Expected : %s But you got : %s", expected, err)
	}
}

func TestValidNestedValidation(t *testing.T) {
	type Orang struct {
		Nama    string `map_validator:"nama" json:"nama"`
		Umur    int    `map_validator:"umur" json:"umur"`
		Hoby    string `map_validator:"hoby" json:"hoby"`
		JK      string `map_validator:"jenis_kelamin" json:"jenis_kelamin"`
		Menikah bool   `map_validator:"menikah" json:"menikah"`
		Anak    *Orang `map_validator:"anak" json:"anak"`
	}
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama":          {Type: reflect.String},
			"umur":          {Type: reflect.Int},
			"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
		},
	}
	role.Rules["anak"] = map_validator.Rules{Object: &role, Null: true}
	child1 := map[string]interface{}{
		"nama":          "Messi",
		"jenis_kelamin": "laki-laki",
		"hoby":          "football",
		"umur":          5,
		"menikah":       false,
	}
	payload := map[string]interface{}{
		"nama":          "Arian Saputra",
		"jenis_kelamin": "laki-laki",
		"hoby":          "Main PS bro",
		"umur":          33,
		"menikah":       true,
		"anak":          child1,
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	extraCheck, err := check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}

	testBind := &Orang{}
	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
		return
	}
	err = extraCheck.Bind(testBind)
	if err != nil {
		t.Errorf("Expected have an error, but you got no error : %s", err)
		return
	}

	if testBind.JK != payload["jenis_kelamin"] {
		t.Errorf("Expected : %s But you got : %s", payload["jenis_kelamin"], testBind.JK)
		return
	}

	if testBind.Anak.Nama != child1["nama"] {
		t.Errorf("Expected : %s But you got : %s", child1["nama"], testBind.Anak.Nama)
		return
	}
}

func TestInvalidMultiNestedValidation(t *testing.T) {
	type Orang struct {
		Nama    string `map_validator:"nama" json:"nama"`
		Umur    int    `map_validator:"umur" json:"umur"`
		Hoby    string `map_validator:"hoby" json:"hoby"`
		JK      string `map_validator:"jenis_kelamin" json:"jenis_kelamin"`
		Menikah bool   `map_validator:"menikah" json:"menikah"`
		Anak    *Orang `map_validator:"anak" json:"anak"`
	}
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama":          {Type: reflect.String},
			"umur":          {Type: reflect.Int},
			"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
		},
	}
	role.Rules["anak"] = map_validator.Rules{Object: &role, Null: true}
	child1OfChild1 := map[string]interface{}{
		"nama":          "Ronaldo",
		"jenis_kelamin": "laki-laki",
		"hoby":          1212,
		"umur":          10,
		"menikah":       false,
	}
	child1 := map[string]interface{}{
		"nama":          "Messi",
		"jenis_kelamin": "laki-laki",
		"hoby":          "football",
		"umur":          20,
		"menikah":       true,
		"anak":          child1OfChild1,
	}
	payload := map[string]interface{}{
		"nama":          "Arian Saputra",
		"jenis_kelamin": "laki-laki",
		"hoby":          "Main PS bro",
		"umur":          33,
		"menikah":       true,
		"anak":          child1,
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	extraCheck, err := check.RunValidate()
	expected := "the field 'hoby' should be 'string'"
	if err.Error() != expected {
		t.Errorf("Expected :%s. But you got : %s", expected, err)
	}

	testBind := &Orang{}
	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
		return
	}
	err = extraCheck.Bind(testBind)
	expected = "no data to Bind because last progress is error"
	if err.Error() != expected {
		t.Errorf("Expected : %s But you got : %s", expected, err)
	}
}

func TestInvalidMultiMultiNestedValidation(t *testing.T) {
	type Orang struct {
		Nama    string `map_validator:"nama" json:"nama"`
		Umur    int    `map_validator:"umur" json:"umur"`
		Hoby    string `map_validator:"hoby" json:"hoby"`
		JK      string `map_validator:"jenis_kelamin" json:"jenis_kelamin"`
		Menikah bool   `map_validator:"menikah" json:"menikah"`
		Anak    *Orang `map_validator:"anak" json:"anak"`
	}
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama":          {Type: reflect.String},
			"umur":          {Type: reflect.Int},
			"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
		},
	}
	role.Rules["anak"] = map_validator.Rules{Object: &role, Null: true}
	child1OfChild1OfChild1 := map[string]interface{}{
		"nama":          false,
		"jenis_kelamin": "laki-laki",
		"hoby":          "run",
		"umur":          10,
		"menikah":       false,
	}
	child1OfChild1 := map[string]interface{}{
		"nama":          "Ronaldo",
		"jenis_kelamin": "laki-laki",
		"hoby":          "football",
		"umur":          10,
		"menikah":       true,
		"anak":          child1OfChild1OfChild1,
	}
	child1 := map[string]interface{}{
		"nama":          "Messi",
		"jenis_kelamin": "laki-laki",
		"hoby":          "football",
		"umur":          20,
		"menikah":       true,
		"anak":          child1OfChild1,
	}
	payload := map[string]interface{}{
		"nama":          "Arian Saputra",
		"jenis_kelamin": "laki-laki",
		"hoby":          "Main PS bro",
		"umur":          33,
		"menikah":       true,
		"anak":          child1,
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	extraCheck, err := check.RunValidate()
	expected := "the field 'nama' should be 'string'"
	if err.Error() != expected {
		t.Errorf("Expected :%s. But you got : %s", expected, err)
	}

	testBind := &Orang{}
	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
		return
	}
	err = extraCheck.Bind(testBind)
	expected = "no data to Bind because last progress is error"
	if err.Error() != expected {
		t.Errorf("Expected : %s But you got : %s", expected, err)
	}
}

func TestNotOptionalInvalidMultiNestedValidation(t *testing.T) {
	type Orang struct {
		Nama    string `map_validator:"nama" json:"nama"`
		Umur    int    `map_validator:"umur" json:"umur"`
		Hoby    string `map_validator:"hoby" json:"hoby"`
		JK      string `map_validator:"jenis_kelamin" json:"jenis_kelamin"`
		Menikah bool   `map_validator:"menikah" json:"menikah"`
		Anak    *Orang `map_validator:"anak" json:"anak"`
	}
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama":          {Type: reflect.String},
			"umur":          {Type: reflect.Int},
			"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
		},
	}
	roleChild := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama":          {Type: reflect.String},
			"umur":          {Type: reflect.Int},
			"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
		},
	}
	roleChild.Rules["anak"] = map_validator.Rules{Object: &roleChild, Null: true}
	role.Rules["anak"] = map_validator.Rules{Object: &roleChild}
	child1OfChild1 := map[string]interface{}{
		"nama":          "Ronaldo",
		"jenis_kelamin": "laki-laki",
		"hoby":          1212,
		"umur":          10,
		"menikah":       false,
	}
	child1 := map[string]interface{}{
		"nama":          "Messi",
		"jenis_kelamin": "laki-laki",
		"hoby":          "football",
		"umur":          20,
		"menikah":       true,
		"anak":          child1OfChild1,
	}
	payload := map[string]interface{}{
		"nama":          "Arian Saputra",
		"jenis_kelamin": "laki-laki",
		"hoby":          "Main PS bro",
		"umur":          33,
		"menikah":       true,
		"anak":          child1,
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	extraCheck, err := check.RunValidate()
	expected := "the field 'hoby' should be 'string'"
	if err.Error() != expected {
		t.Errorf("Expected :%s. But you got : %s", expected, err)
	}

	testBind := &Orang{}
	if testBind.JK != "" {
		t.Errorf("Expected : '' But you got : %s", testBind.JK)
		return
	}
	err = extraCheck.Bind(testBind)
	expected = "no data to Bind because last progress is error"
	if err.Error() != expected {
		t.Errorf("Expected : %s But you got : %s", expected, err)
	}
}

func TestInvalidListMultiNestedValidation(t *testing.T) {
	type Orang struct {
		Nama    string `map_validator:"nama" json:"nama"`
		Umur    int    `map_validator:"umur" json:"umur"`
		Hoby    string `map_validator:"hoby" json:"hoby"`
		JK      string `map_validator:"jenis_kelamin" json:"jenis_kelamin"`
		Menikah bool   `map_validator:"menikah" json:"menikah"`
		Anak    *Orang `map_validator:"anak" json:"anak"`
	}
	lastRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama":          {Type: reflect.String},
			"umur":          {Type: reflect.Int},
			"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
		},
	}
	secondRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama":          {Type: reflect.String},
			"umur":          {Type: reflect.Int},
			"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
			"anak-anak":     {ListObject: &lastRole},
		},
	}
	rootRole := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama":          {Type: reflect.String},
			"umur":          {Type: reflect.Int},
			"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
			"anak":          {Object: &secondRole, Null: true},
		},
	}
	validLastChild := []map[string]interface{}{
		{
			"nama":          "Messi",
			"jenis_kelamin": "laki-laki",
			"hoby":          "football",
			"umur":          20,
			"menikah":       true,
		},
	}
	invalidLastChild := []map[string]interface{}{
		{
			"nama":          1122,
			"jenis_kelamin": "laki-laki",
			"hoby":          "football",
			"umur":          20,
			"menikah":       true,
		},
	}
	lastChild := map[string]interface{}{
		"nama":          "Messi",
		"jenis_kelamin": "laki-laki",
		"hoby":          "football",
		"umur":          20,
		"menikah":       true,
	}
	secondChild := map[string]interface{}{
		"nama":          "Messi",
		"jenis_kelamin": "laki-laki",
		"hoby":          "football",
		"umur":          20,
		"menikah":       true,
		"anak-anak":     lastChild,
	}
	payload := map[string]interface{}{
		"nama":          "Arian Saputra",
		"jenis_kelamin": "laki-laki",
		"hoby":          "Main PS bro",
		"umur":          33,
		"menikah":       true,
		"anak":          secondChild,
	}
	check, err := map_validator.NewValidateBuilder().SetRules(rootRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	_, err = check.RunValidate()
	expected := "field 'anak-anak' is not valid list object"
	if err.Error() != expected {
		t.Errorf("Expected :%s. But you got : %s", expected, err)
	}

	secondChild["anak-anak"] = validLastChild
	payload["anak"] = secondChild
	check, err = map_validator.NewValidateBuilder().SetRules(rootRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	_, err = check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}

	secondChild["anak-anak"] = invalidLastChild
	payload["anak"] = secondChild
	check, err = map_validator.NewValidateBuilder().SetRules(rootRole).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	_, err = check.RunValidate()
	expected = "the field 'nama' should be 'string'"
	if err.Error() != expected {
		t.Errorf("Expected :%s. But you got : %s", expected, err)
	}
}

func TestNestedStrict(t *testing.T) {
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"nama":          {Type: reflect.String},
			"umur":          {Type: reflect.Int},
			"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
			"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
			"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
			"anak": {Object: &map_validator.RulesWrapper{
				Rules: map[string]map_validator.Rules{
					"nama": {Type: reflect.String},
					"umur": {Type: reflect.Int},
				},
				Setting: map_validator.Setting{Strict: true},
			}},
		},
	}
	payload := map[string]interface{}{
		"nama":          "Arian Saputra",
		"jenis_kelamin": "laki-laki",
		"hoby":          "Main PS bro",
		"umur":          33,
		"menikah":       true,
		"anak": map[string]interface{}{
			"nama":          "false",
			"jenis_kelamin": "laki-laki",
			"umur":          10,
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	_, err = check.RunValidate()
	expected := "'jenis_kelamin' is not allowed key"
	if err.Error() != expected {
		t.Errorf("Expected : %s. But you got : %s", expected, err)
	}
}

func TestNestedListDataHTTP(t *testing.T) {
	type GoodsRequest struct {
		Name        string  `json:"name" binding:"required"`
		Weight      float64 `json:"weight" binding:"required"`
		Quantity    int     `json:"quantity" binding:"required"`
		Description string  `json:"description"`
	}
	type CreateOrderRequest struct {
		SenderID                uuid.UUID      `json:"sender_id" binding:"required"`
		SenderAddress           string         `json:"sender_address" binding:"required"`
		SenderAddressCity       string         `json:"sender_address_city" binding:"required"`
		SenderAddressProvince   string         `json:"sender_address_province" binding:"required"`
		SenderLatitude          float64        `json:"sender_latitude"`
		SenderLongitude         float64        `json:"sender_longitude"`
		ReceiverName            string         `json:"receiver_name" binding:"required"`
		ReceiverPhone           string         `json:"receiver_phone" binding:"required"`
		ReceiverAddress         string         `json:"receiver_address" binding:"required"`
		ReceiverAddressCity     string         `json:"receiver_address_city" binding:"required"`
		ReceiverAddressProvince string         `json:"receiver_address_province" binding:"required"`
		ReceiverLatitude        float64        `json:"receiver_latitude"`
		ReceiverLongitude       float64        `json:"receiver_longitude"`
		Note                    string         `json:"note"`
		Goods                   []GoodsRequest `json:"goods" binding:"required"`
	}
	var request CreateOrderRequest
	jsonStr := `{
		"goods": [
			{
			"description": "string",
			"name": "string",
			"quantity": 1,
			"weight": 0
			}
		],
		"note": "string",
		"receiver_address": "string",
		"receiver_address_city": "string",
		"receiver_address_province": "string",
		"receiver_latitude": 0,
		"receiver_longitude": 0,
		"receiver_name": "string",
		"receiver_phone": "string",
		"sender_address": "string",
		"sender_address_city": "string",
		"sender_address_province": "string",
		"sender_id": "8aa3e797-2453-442f-b1d0-50f7d815bcaf",
		"sender_latitude": 0,
		"sender_longitude": 0
		}`
	// fake request
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(map_validator.BuildRoles().
		SetRule("sender_id", map_validator.Rules{Type: reflect.String, UUID: true}).
		SetRule("sender_address", map_validator.Rules{Type: reflect.String}).
		SetRule("sender_address_city", map_validator.Rules{Type: reflect.String}).
		SetRule("sender_address_province", map_validator.Rules{Type: reflect.String}).
		SetRule("sender_latitude", map_validator.Rules{Type: reflect.Float64}).
		SetRule("sender_longitude", map_validator.Rules{Type: reflect.Float64}).
		SetRule("receiver_name", map_validator.Rules{Type: reflect.String}).
		SetRule("receiver_phone", map_validator.Rules{Type: reflect.String}).
		SetRule("receiver_address", map_validator.Rules{Type: reflect.String}).
		SetRule("receiver_address_city", map_validator.Rules{Type: reflect.String}).
		SetRule("receiver_address_province", map_validator.Rules{Type: reflect.String}).
		SetRule("receiver_latitude", map_validator.Rules{Type: reflect.Float64}).
		SetRule("receiver_longitude", map_validator.Rules{Type: reflect.Float64}).
		SetRule("note", map_validator.Rules{Type: reflect.String}).
		SetRule("goods", map_validator.Rules{ListObject: map_validator.BuildRoles().
			SetRule("name", map_validator.Rules{Type: reflect.String}).
			SetRule("weight", map_validator.Rules{Type: reflect.Float64}).
			SetRule("quantity", map_validator.Rules{Type: reflect.Int, Min: map_validator.SetTotal(1)}).
			SetRule("description", map_validator.Rules{Type: reflect.String}),
		}).
		Done()).
		LoadJsonHttp(req)
	if err != nil {
		t.Errorf("load error : %s", err)
	}

	jsonData, err := jsonHttp.RunValidate()
	if err != nil {
		t.Errorf("validate err : %s", err)
	}

	if err = jsonData.Bind(&request); err != nil {
		t.Errorf("err : %s", err)
	}
}
