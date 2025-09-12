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

// Test int64 enum with HTTP JSON to ensure integer family coercion works
func TestInt64EnumWithHttpRequest(t *testing.T) {
	jsonStr := `{"status": 200, "code": 404}`

	rules := map_validator.BuildRoles().
		SetRule("status", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int64{200, 201, 400, 404, 500}},
		}).
		SetRule("code", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int64{200, 404, 500}},
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
	if data["status"] != float64(200) {
		t.Errorf("Expected status to be 200, got %v", data["status"])
	}
	if data["code"] != float64(404) {
		t.Errorf("Expected code to be 404, got %v", data["code"])
	}
}

// Test int32 enum with HTTP JSON
func TestInt32EnumWithHttpRequest(t *testing.T) {
	jsonStr := `{"priority": 1, "level": 5}`

	rules := map_validator.BuildRoles().
		SetRule("priority", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int32{1, 2, 3}},
		}).
		SetRule("level", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int32{1, 3, 5, 7, 9}},
		})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %s", err)
	}

	_, err = jsonHttp.RunValidate()
	if err != nil {
		t.Fatalf("Expected no fail, but it failed: %s", err.Error())
	}
}

// Test int16 enum with HTTP JSON
func TestInt16EnumWithHttpRequest(t *testing.T) {
	jsonStr := `{"port": 80, "timeout": 30}`

	rules := map_validator.BuildRoles().
		SetRule("port", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int16{80, 443, 8080}},
		}).
		SetRule("timeout", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int16{10, 30, 60}},
		})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %s", err)
	}

	_, err = jsonHttp.RunValidate()
	if err != nil {
		t.Fatalf("Expected no fail, but it failed: %s", err.Error())
	}
}

// Test int8 enum with HTTP JSON
func TestInt8EnumWithHttpRequest(t *testing.T) {
	jsonStr := `{"flag": 1, "mode": 5}`

	rules := map_validator.BuildRoles().
		SetRule("flag", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int8{0, 1, 2}},
		}).
		SetRule("mode", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int8{1, 3, 5, 7}},
		})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %s", err)
	}

	_, err = jsonHttp.RunValidate()
	if err != nil {
		t.Fatalf("Expected no fail, but it failed: %s", err.Error())
	}
}

// Test uint family enum with HTTP JSON
func TestUintEnumWithHttpRequest(t *testing.T) {
	jsonStr := `{"count": 100, "limit": 50}`

	rules := map_validator.BuildRoles().
		SetRule("count", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []uint{10, 50, 100, 200}},
		}).
		SetRule("limit", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []uint64{25, 50, 100}},
		})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %s", err)
	}

	_, err = jsonHttp.RunValidate()
	if err != nil {
		t.Fatalf("Expected no fail, but it failed: %s", err.Error())
	}
}

// Test negative number with uint enum should fail
func TestNegativeNumberUintEnumWithHttpRequest(t *testing.T) {
	jsonStr := `{"count": -1}`

	rules := map_validator.BuildRoles().
		SetRule("count", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []uint{1, 2, 3}},
		})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %s", err)
	}

	_, err = jsonHttp.RunValidate()
	if err == nil {
		t.Fatal("Expected validation to fail for negative number with uint enum, but it passed")
	}
}

// Test float32 enum with HTTP JSON  
func TestFloat32EnumWithHttpRequest(t *testing.T) {
	jsonStr := `{"rate": 1.5, "factor": 2.0}`

	rules := map_validator.BuildRoles().
		SetRule("rate", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []float32{1.0, 1.5, 2.0, 2.5}},
		}).
		SetRule("factor", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []float32{1.0, 2.0, 3.0}},
		})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %s", err)
	}

	_, err = jsonHttp.RunValidate()
	if err != nil {
		t.Fatalf("Expected no fail, but it failed: %s", err.Error())
	}
}

// Test float64 enum with HTTP JSON (should work naturally)
func TestFloat64EnumWithHttpRequest(t *testing.T) {
	jsonStr := `{"percentage": 85.5, "ratio": 0.75}`

	rules := map_validator.BuildRoles().
		SetRule("percentage", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []float64{25.0, 50.0, 75.0, 85.5, 100.0}},
		}).
		SetRule("ratio", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []float64{0.25, 0.5, 0.75, 1.0}},
		})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
	if err != nil {
		t.Fatalf("load error: %s", err)
	}

	_, err = jsonHttp.RunValidate()
	if err != nil {
		t.Fatalf("Expected no fail, but it failed: %s", err.Error())
	}
}

// Test mixed integer family types in single request
func TestMixedIntegerFamilyEnumWithHttpRequest(t *testing.T) {
	jsonStr := `{
		"id": 12345,
		"status": 200,
		"priority": 1,
		"level": 5,
		"flag": 1,
		"count": 100,
		"rate": 1.5
	}`

	rules := map_validator.BuildRoles().
		SetRule("id", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int64{12345, 67890}},
		}).
		SetRule("status", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int32{200, 400, 500}},
		}).
		SetRule("priority", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int16{1, 2, 3}},
		}).
		SetRule("level", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []int8{1, 3, 5, 7}},
		}).
		SetRule("flag", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []uint8{0, 1, 2}},
		}).
		SetRule("count", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []uint{50, 100, 200}},
		}).
		SetRule("rate", map_validator.Rules{
			Enum: &map_validator.EnumField[any]{Items: []float32{1.0, 1.5, 2.0}},
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

	// Verify all data integrity
	data := extra.GetData()
	if data["id"] != float64(12345) {
		t.Errorf("Expected id to be 12345, got %v", data["id"])
	}
	if data["status"] != float64(200) {
		t.Errorf("Expected status to be 200, got %v", data["status"])
	}
	if data["rate"] != float64(1.5) {
		t.Errorf("Expected rate to be 1.5, got %v", data["rate"])
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
