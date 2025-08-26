package test

import (
	"reflect"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

// // TestValidList a positive test case for the new "List" rule.
func TestValidList(t *testing.T) {
	payload := map[string]interface{}{
		"tags":       []string{"RED", "BLUE"},
		"user_ids":   []string{"6ba7b810-9dad-11d1-80b4-00c04fd430c8", "6ba7b811-9dad-11d1-80b4-00c04fd430c8"},
		"numbers":    []int{10, 20, 30},
		"empty_list": []string{},
	}

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.Rules{List: map_validator.BuildListRoles(), Enum: &map_validator.EnumField[any]{Items: []string{"GREEN", "BLUE", "RED"}}, Max: map_validator.SetTotal(5)}).
		SetRule("user_ids", map_validator.Rules{List: map_validator.BuildListRoles(), UUID: true}).
		SetRule("numbers", map_validator.Rules{List: map_validator.BuildListRoles(), Type: reflect.Int, Min: map_validator.SetTotal(3)}).
		SetRule("empty_list", map_validator.Rules{List: map_validator.BuildListRoles(), Type: reflect.String})

	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(payload)
	if err != nil {
		t.Fatalf("Load error: %s", err)
	}

	if _, err = op.RunValidate(); err != nil {
		t.Errorf("Expected no error, but got error: %s", err)
	}
}

// // TestInvalidList a negative test case for the new "List" rule.
// func TestInvalidList(t *testing.T) {
// 	// Test case 1: Field is not a list
// 	payload1 := map[string]interface{}{"user_ids": "not-a-list"}
// 	rules1 := map_validator.BuildRoles().
// 		SetRule("user_ids", map_validator.Rules{List: map_validator.BuildListRoles(), UUID: true})

// 	op1, err := map_validator.NewValidateBuilder().SetRules(rules1).Load(payload1)
// 	if err != nil {
// 		t.Fatalf("Load error on payload1: %s", err)
// 	}
// 	_, err1 := op1.RunValidate()
// 	expectedErr1 := "the field 'user_ids' should be 'list'"
// 	if err1 == nil {
// 		t.Errorf("Expected error '%s', but got nil", expectedErr1)
// 	} else if err1.Error() != expectedErr1 {
// 		t.Errorf("Expected error '%s', but got '%s'", expectedErr1, err1.Error())
// 	}

// 	// Test case 2: One of the elements is invalid (UUID)
// 	payload2 := map[string]interface{}{"user_ids": []string{"6ba7b810-9dad-11d1-80b4-00c04fd430c8", "invalid-uuid"}}
// 	op2, err := map_validator.NewValidateBuilder().SetRules(rules1).Load(payload2)
// 	if err != nil {
// 		t.Fatalf("Load error on payload2: %s", err)
// 	}
// 	_, err2 := op2.RunValidate()
// 	// Ideal error message would be: "error in field 'user_ids' at index 1: the field 'user_ids[1]' is not valid uuid"
// 	if err2 == nil {
// 		t.Errorf("Expected error containing 'not valid uuid', but got nil")
// 	} else {
// 		// A less specific check for now, as the detailed error message is not yet implemented
// 		t.Logf("Note: A more specific error message is desired. Current error: %s", err2.Error())
// 		if !strings.Contains(err2.Error(), "not valid uuid") {
// 			t.Errorf("Expected error to contain 'not valid uuid', but it did not. Got: %s", err2.Error())
// 		}
// 	}

// 	// Test case 3: One of the elements is invalid (Max length)
// 	payload3 := map[string]interface{}{"tags": []string{"ok", "this-tag-is-too-long"}}
// 	rules3 := map_validator.BuildRoles().
// 		SetRule("tags", map_validator.Rules{List: map_validator.BuildListRoles(), Type: reflect.String, Max: map_validator.SetTotal(5)})
// 	op3, err := map_validator.NewValidateBuilder().SetRules(rules3).Load(payload3)
// 	if err != nil {
// 		t.Fatalf("Load error on payload3: %s", err)
// 	}
// 	_, err3 := op3.RunValidate()
// 	// Ideal error message would be: "error in field 'tags' at index 1: the field 'tags[1]' should be or lower than 5"
// 	if err3 == nil {
// 		t.Errorf("Expected error containing 'lower than 5', but got nil")
// 	} else {
// 		t.Logf("Note: A more specific error message is desired. Current error: %s", err3.Error())
// 		if !strings.Contains(err3.Error(), "lower than 5") {
// 			t.Errorf("Expected error to contain 'lower than 5', but it did not. Got: %s", err3.Error())
// 		}
// 	}
// }

// func TestListWithObjects(t *testing.T) {
// 	objectRules := map_validator.BuildRoles().
// 		SetRule("name", map_validator.Rules{Type: reflect.String, Min: map_validator.SetTotal(3)}).
// 		SetRule("age", map_validator.Rules{Type: reflect.Int, Min: map_validator.SetTotal(18)})

// 	rules := map_validator.BuildRoles().
// 		SetRule("users", map_validator.Rules{List: map_validator.BuildListRoles(), Object: objectRules})

// 	// Success case
// 	validPayload := map[string]interface{}{
// 		"users": []interface{}{
// 			map[string]interface{}{"name": "Alice", "age": 30},
// 			map[string]interface{}{"name": "Bob", "age": 25},
// 		},
// 	}
// 	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(validPayload)
// 	if err != nil {
// 		t.Fatalf("Load error on valid payload: %s", err)
// 	}
// 	if _, err = op.RunValidate(); err != nil {
// 		t.Errorf("Expected no error on valid payload, but got: %s", err)
// 	}

// 	// Failure case
// 	invalidPayload := map[string]interface{}{
// 		"users": []interface{}{
// 			map[string]interface{}{"name": "Alice", "age": 30},
// 			map[string]interface{}{"name": "Eve", "age": 17}, // Invalid age
// 		},
// 	}
// 	op2, err := map_validator.NewValidateBuilder().SetRules(rules).Load(invalidPayload)
// 	if err != nil {
// 		t.Fatalf("Load error on invalid payload: %s", err)
// 	}
// 	_, err2 := op2.RunValidate()
// 	if err2 == nil {
// 		t.Errorf("Expected error on invalid payload, but got nil")
// 	} else {
// 		if !strings.Contains(err2.Error(), "greater than 18") {
// 			t.Errorf("Expected error to be about age validation, but got: %s", err2.Error())
// 		}
// 	}
// }

// func TestListWithMinMax(t *testing.T) {
// 	rules := map_validator.BuildRoles().
// 		SetRule("tags", map_validator.Rules{List: map_validator.BuildListRoles(), Type: reflect.String, Min: map_validator.SetTotal(2), Max: map_validator.SetTotal(3)})

// 	// Success case
// 	validPayload := map[string]interface{}{"tags": []string{"go", "test"}}
// 	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(validPayload)
// 	if err != nil {
// 		t.Fatalf("Load error on valid payload: %s", err)
// 	}
// 	if _, err = op.RunValidate(); err != nil {
// 		t.Errorf("Expected no error on valid payload, but got: %s", err)
// 	}

// 	// Failure case (too few elements)
// 	invalidPayloadMin := map[string]interface{}{"tags": []string{"go"}}
// 	op2, err := map_validator.NewValidateBuilder().SetRules(rules).Load(invalidPayloadMin)
// 	if err != nil {
// 		t.Fatalf("Load error on invalid min payload: %s", err)
// 	}
// 	_, err2 := op2.RunValidate()
// 	if err2 == nil {
// 		t.Errorf("Expected error on min payload, but got nil")
// 	} else if !strings.Contains(err2.Error(), "greater than 2") {
// 		t.Errorf("Expected error to be about min elements, but got: %s", err2.Error())
// 	}

// 	// Failure case (too many elements)
// 	invalidPayloadMax := map[string]interface{}{"tags": []string{"go", "test", "tdd", "fail"}}
// 	op3, err := map_validator.NewValidateBuilder().SetRules(rules).Load(invalidPayloadMax)
// 	if err != nil {
// 		t.Fatalf("Load error on invalid max payload: %s", err)
// 	}
// 	_, err3 := op3.RunValidate()
// 	if err3 == nil {
// 		t.Errorf("Expected error on max payload, but got nil")
// 	} else if !strings.Contains(err3.Error(), "lower than 3") {
// 		t.Errorf("Expected error to be about max elements, but got: %s", err3.Error())
// 	}
// }

// func TestListWithHttpRequest(t *testing.T) {
// 	jsonStr := `{
// 		"name": "My Awesome List",
// 		"emails": ["test1@example.com", "test2@example.com", "not-an-email", "test4@example.com"]
// 	}`

// 	rules := map_validator.BuildRoles().
// 		SetRule("name", map_validator.Rules{Type: reflect.String}).
// 		SetRule("emails", map_validator.Rules{List: map_validator.BuildListRoles(), Email: true})

// 	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
// 	req.Header.Set("Content-Type", "application/json")

// 	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
// 	if err != nil {
// 		t.Fatalf("load error : %s", err)
// 	}

// 	_, err = jsonHttp.RunValidate()
// 	if err == nil {
// 		t.Fatal("Expected validation to fail, but it passed")
// 	}

// 	// Check if the error message is correct
// 	if !strings.Contains(err.Error(), "not valid email") {
// 		t.Errorf("Expected error to be about email validation, but got: %s", err.Error())
// 	}
// }

// func TestListWithNull(t *testing.T) {
// 	// Scenario 1: Null is allowed and field is provided as nil
// 	rules := map_validator.BuildRoles().
// 		SetRule("tags", map_validator.Rules{List: map_validator.BuildListRoles(), Type: reflect.String, Null: true})

// 	payload1 := map[string]interface{}{"tags": nil}
// 	op1, err := map_validator.NewValidateBuilder().SetRules(rules).Load(payload1)
// 	if err != nil {
// 		t.Fatalf("Load error on payload1: %s", err)
// 	}
// 	if _, err = op1.RunValidate(); err != nil {
// 		t.Errorf("Scenario 1 failed: Expected no error, but got %s", err)
// 	}

// 	// Scenario 2: Null is allowed and field is missing
// 	payload2 := map[string]interface{}{}
// 	op2, err := map_validator.NewValidateBuilder().SetRules(rules).Load(payload2)
// 	if err != nil {
// 		t.Fatalf("Load error on payload2: %s", err)
// 	}
// 	if _, err = op2.RunValidate(); err != nil {
// 		t.Errorf("Scenario 2 failed: Expected no error, but got %s", err)
// 	}

// 	// Scenario 3: Null is not allowed and field is missing
// 	rules2 := map_validator.BuildRoles().
// 		SetRule("tags", map_validator.Rules{List: map_validator.BuildListRoles(), Type: reflect.String, Null: false})

// 	payload3 := map[string]interface{}{}
// 	op3, err := map_validator.NewValidateBuilder().SetRules(rules2).Load(payload3)
// 	if err != nil {
// 		t.Fatalf("Load error on payload3: %s", err)
// 	}
// 	_, err3 := op3.RunValidate()
// 	expectedErr := "we need 'tags' field"
// 	if err3 == nil {
// 		t.Errorf("Scenario 3 failed: Expected error '%s', but got nil", expectedErr)
// 	} else if err3.Error() != expectedErr {
// 		t.Errorf("Scenario 3 failed: Expected error '%s', but got '%s'", expectedErr, err3.Error())
// 	}
// }
