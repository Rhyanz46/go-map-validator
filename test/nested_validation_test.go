package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"reflect"
	"testing"
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
	role := map[string]map_validator.Rules{
		"nama":          {Type: reflect.String},
		"umur":          {Type: reflect.Int},
		"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
	}
	role["anak"] = map_validator.Rules{Object: &role, Null: true}
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
	role := map[string]map_validator.Rules{
		"nama":          {Type: reflect.String},
		"umur":          {Type: reflect.Int},
		"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
	}
	role["anak"] = map_validator.Rules{Object: &role, Null: true}
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
	role := map[string]map_validator.Rules{
		"nama":          {Type: reflect.String},
		"umur":          {Type: reflect.Int},
		"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
	}
	role["anak"] = map_validator.Rules{Object: &role, Null: true}
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
	role := map[string]map_validator.Rules{
		"nama":          {Type: reflect.String},
		"umur":          {Type: reflect.Int},
		"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
	}
	role["anak"] = map_validator.Rules{Object: &role, Null: true}
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
	role := map[string]map_validator.Rules{
		"nama":          {Type: reflect.String},
		"umur":          {Type: reflect.Int},
		"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
	}
	roleChild := map[string]map_validator.Rules{
		"nama":          {Type: reflect.String},
		"umur":          {Type: reflect.Int},
		"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
	}
	roleChild["anak"] = map_validator.Rules{Object: &roleChild, Null: true}
	role["anak"] = map_validator.Rules{Object: &roleChild}
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
	lastRole := map[string]map_validator.Rules{
		"nama":          {Type: reflect.String},
		"umur":          {Type: reflect.Int},
		"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
	}
	secondRole := map[string]map_validator.Rules{
		"nama":          {Type: reflect.String},
		"umur":          {Type: reflect.Int},
		"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
		"anak-anak":     {ListObject: &lastRole},
	}
	rootRole := map[string]map_validator.Rules{
		"nama":          {Type: reflect.String},
		"umur":          {Type: reflect.Int},
		"hoby":          {Type: reflect.String, Null: true, IfNull: "-"},
		"jenis_kelamin": {Enum: &map_validator.EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"menikah":       {Type: reflect.Bool, Null: true, IfNull: false},
		"anak":          {Object: &secondRole, Null: true},
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
