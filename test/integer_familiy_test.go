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

// Test integer enum with multiple values from HTTP JSON
func TestMultipleIntegerEnumWithHttpRequest(t *testing.T) {
	jsonStr := `{
		"port": 443,
		"protocol": 1,
		"priority": 5
	}`

	rules := map_validator.BuildRoles().
		SetRule("port", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int{80, 443, 8080}},
		}).
		SetRule("protocol", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int{0, 1, 2}},
		}).
		SetRule("priority", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int{1, 3, 5, 7, 9}},
		})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %s", err)
	}

	extra, err := jsonHttp.RunValidate()
	if err != nil {
		t.Fatalf("Expected no fail, but it failed: %s", err.Error())
	}

	// Verify data integrity
	data := extra.GetData()
	if data["port"] != float64(443) {
		t.Errorf("Expected port to be 443, got %v", data["port"])
	}
	if data["protocol"] != float64(1) {
		t.Errorf("Expected protocol to be 1, got %v", data["protocol"])
	}
	if data["priority"] != float64(5) {
		t.Errorf("Expected priority to be 5, got %v", data["priority"])
	}
}

// Test invalid integer enum value from HTTP JSON should fail
func TestInvalidIntegerEnumWithHttpRequest(t *testing.T) {
	jsonStr := `{"port": 8888}`

	rules := map_validator.BuildRoles().SetRule("port", map_validator.Rules{
		Enum: &map_validator.EnumField[any]{
			Items: []int{80, 443, 8080},
		},
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %s", err)
	}

	_, err = jsonHttp.RunValidate()
	if err == nil {
		t.Fatal("Expected validation to fail for invalid enum value, but it passed")
	}

	expectedErrorSubstring := "not in enum list"
	if err != nil && !contains(err.Error(), expectedErrorSubstring) {
		t.Errorf("Expected error to contain '%s', but got: %s", expectedErrorSubstring, err.Error())
	}
}

// Test decimal number should fail for integer enum from HTTP JSON
func TestDecimalNumberIntegerEnumWithHttpRequest(t *testing.T) {
	jsonStr := `{"port": 80.5}`

	rules := map_validator.BuildRoles().SetRule("port", map_validator.Rules{
		Enum: &map_validator.EnumField[any]{
			Items: []int{80, 443},
		},
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %s", err)
	}

	_, err = jsonHttp.RunValidate()
	if err == nil {
		t.Fatal("Expected validation to fail for decimal number in integer enum, but it passed")
	}
}

// Test mixed Type and Enum validation for integer family from HTTP JSON
func TestMixedTypeAndIntegerEnumWithHttpRequest(t *testing.T) {
	jsonStr := `{
		"port": 443,
		"timeout": 30
	}`

	rules := map_validator.BuildRoles().
		SetRule("port", map_validator.Rules{
			Type: reflect.Int,
			Enum: &map_validator.EnumField[any]{Items: []int{80, 443, 8080}},
		}).
		SetRule("timeout", map_validator.Rules{
			Type: reflect.Int,
			Enum: &map_validator.EnumField[any]{Items: []int{10, 30, 60, 120}},
		})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %s", err)
	}

	extra, err := jsonHttp.RunValidate()
	if err != nil {
		t.Fatalf("Expected no fail, but it failed: %s", err.Error())
	}

	// Verify data integrity
	data := extra.GetData()
	if data["port"] != float64(443) {
		t.Errorf("Expected port to be 443, got %v", data["port"])
	}
	if data["timeout"] != float64(30) {
		t.Errorf("Expected timeout to be 30, got %v", data["timeout"])
	}
}

// Test large integer enum values from HTTP JSON  
func TestLargeIntegerEnumWithHttpRequest(t *testing.T) {
	jsonStr := `{
		"user_id": 1000000,
		"session_id": 9999999
	}`

	rules := map_validator.BuildRoles().
		SetRule("user_id", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int{1000000, 2000000, 3000000}},
		}).
		SetRule("session_id", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int{1111111, 9999999, 8888888}},
		})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %s", err)
	}

	extra, err := jsonHttp.RunValidate()
	if err != nil {
		t.Fatalf("Expected no fail, but it failed: %s", err.Error())
	}

	// Verify data integrity
	data := extra.GetData()
	if data["user_id"] != float64(1000000) {
		t.Errorf("Expected user_id to be 1000000, got %v", data["user_id"])
	}
	if data["session_id"] != float64(9999999) {
		t.Errorf("Expected session_id to be 9999999, got %v", data["session_id"])
	}
}

// Helper function for string contains check (if not already defined elsewhere)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || indexOfSubstring(s, substr) >= 0)
}

func indexOfSubstring(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
