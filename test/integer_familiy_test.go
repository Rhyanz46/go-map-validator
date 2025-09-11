package test

import (
	"bytes"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

func TestIntegerEnumWithHttpRequest(t *testing.T) {
	jsonStr := `{"port": 80}`

	rules := map_validator.BuildRoles().SetRule("port", map_validator.Rules{
		Enum: &map_validator.EnumField[any]{
			Items: []int{80, 443},
		},
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error : %s", err)
	}

	_, err = jsonHttp.RunValidate()
	if err != nil {
		t.Fatalf("Expected no fail, but it fail, %s", err.Error())
	}

}

func TestStrinEnumWithHttpRequest(t *testing.T) {
	jsonStr := `{
		"port": 80
	}`

	rules := map_validator.BuildRoles().
		SetRule("port", map_validator.Rules{Type: reflect.Int})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error : %s", err)
	}

	_, err = jsonHttp.RunValidate()
	if err != nil {
		t.Fatalf("Expected no fail, but it fail, %s", err.Error())
	}

}

func TestIntegerEnumFamily(t *testing.T) {
	dataMap := map[string]interface{}{
		"port": 80,
	}

	rules := map_validator.BuildRoles().SetRule("port", map_validator.Rules{
		Enum: &map_validator.EnumField[any]{
			Items: []int{80, 443},
		},
	})

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).Load(dataMap)
	if err != nil {
		t.Fatalf("load error : %s", err)
	}

	_, err = jsonHttp.RunValidate()
	if err != nil {
		t.Fatalf("Expected no fail, but it fail, %s", err.Error())
	}

}
