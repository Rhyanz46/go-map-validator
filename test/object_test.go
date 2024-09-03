package test

import (
	"fmt"
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

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

func TestInvalidListMultiNestedValidationBind(t *testing.T) {
	type People struct {
		Name string `map_validator:"name" json:"name"`
		Age  int    `map_validator:"age" json:"age"`
	}
	type Class struct {
		RoomId  uuid.UUID `map_validator:"room_id" json:"room_id"`
		Name    string    `map_validator:"name"`
		Peoples []People  `json:"peoples"`
	}
	data := Class{}
	role := map_validator.RulesWrapper{
		Rules: map[string]map_validator.Rules{
			"room_id": {UUID: true},
			"name":    {Type: reflect.String},
			"peoples": {ListObject: &map_validator.RulesWrapper{
				Rules: map[string]map_validator.Rules{
					"name": {Type: reflect.String},
					"id":   {Type: reflect.Int},
				},
			}, Null: true},
		},
	}
	payload := map[string]interface{}{
		"name":    "English 1A",
		"room_id": "c18f068f-272c-491f-b379-5fd9bec8ada3",
		"peoples": []map[string]interface{}{
			{
				"name": "Arian",
				"id":   1,
			},
		},
	}
	check, err := map_validator.NewValidateBuilder().SetRules(role).Load(payload)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}
	jsonData, err := check.RunValidate()
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}

	err = jsonData.Bind(&data)
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
		return
	}

	fmt.Println(data)
}
