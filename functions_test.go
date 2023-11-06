package mapValidator

import (
	"reflect"
	"testing"
)

func TestValidateUUID(t *testing.T) {
	payload := map[string]interface{}{
		"field1": "123e4567-e89b-12d3-a456-426614174001",
	}

	validator := RequestDataValidator{
		UUID: true,
		Null: false,
	}

	_, err := Validate("field1", payload, validator)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateInvalidUUID(t *testing.T) {
	payload := map[string]interface{}{
		"field1": "invalid-uuid",
	}

	validator := RequestDataValidator{
		UUID: true,
		Null: false,
	}

	_, err := Validate("field1", payload, validator)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	expectedError := "the field 'field1' it's not valid uuid"
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expected error: %s, got: %v", expectedError, err)
	}
}

func TestValidateNotNull(t *testing.T) {
	payload := map[string]interface{}{
		"field1": "value",
	}

	validator := RequestDataValidator{
		Type: reflect.String,
		Null: false,
	}

	_, err := Validate("field1", payload, validator)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateNullNotAllowed(t *testing.T) {
	payload := map[string]interface{}{
		"field1": nil,
	}

	validator := RequestDataValidator{
		Type: reflect.String,
		Null: false,
	}

	_, err := Validate("field1", payload, validator)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	expectedError := "we need 'field1' field"
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expected error: %s, got: %v", expectedError, err)
	}
}

func TestValidateStringMaxLength(t *testing.T) {
	payload := map[string]interface{}{
		"field1": "1234567890",
	}

	validator := RequestDataValidator{
		Type: reflect.String,
		Max:  ToPointer[int](5),
	}

	_, err := Validate("field1", payload, validator)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	expectedError := "the field 'field1' should be or lower than 5 character"
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expected error: %s, got: %v", expectedError, err)
	}
}

func TestValidateStringMinLength(t *testing.T) {
	payload := map[string]interface{}{
		"field1": "123",
	}

	validator := RequestDataValidator{
		Type: reflect.String,
		Min:  ToPointer[int](5),
	}

	_, err := Validate("field1", payload, validator)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	expectedError := "the field 'field1' should be or greater than 5 character"
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expected error: %s, got: %v", expectedError, err)
	}
}

func TestValidateEmail(t *testing.T) {
	payload := map[string]interface{}{
		"email": "test@example.com",
	}

	validator := RequestDataValidator{
		Email: true,
	}

	_, err := Validate("email", payload, validator)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateInvalidEmail(t *testing.T) {
	payload := map[string]interface{}{
		"email": "invalid-email",
	}

	validator := RequestDataValidator{
		Email: true,
	}

	_, err := Validate("email", payload, validator)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	expectedError := "field email is not valid email"
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expected error: %s, got: %v", expectedError, err)
	}
}

func TestValidateIPV4(t *testing.T) {
	payload := map[string]interface{}{
		"ip_address": "192.168.1.1",
	}

	validator := RequestDataValidator{
		IPV4: true,
	}

	_, err := Validate("ip_address", payload, validator)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateInvalidIPV4(t *testing.T) {
	payload := map[string]interface{}{
		"ip_address": "invalid-ip",
	}

	validator := RequestDataValidator{
		IPV4: true,
	}

	_, err := Validate("ip_address", payload, validator)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	expectedError := "the field 'ip_address' it's not valid IP"
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expected error: %s, got: %v", expectedError, err)
	}
}

func TestEnumFieldCheck(t *testing.T) {
	payload := map[string]interface{}{"data": "arian", "jenis_kelamin": "laki-laki", "hoby": "Main PS"}
	_, err := Validate(
		"data", payload, RequestDataValidator{
			Null: false,
			Enum: &EnumField[any]{Items: []string{"arian", "aaa"}},
		},
	)
	if err != nil {
		t.Errorf("Test case 1 Error : %v", err)
	}

	_, err = Validate(
		"jenis_kelamin", payload, RequestDataValidator{
			Null: false,
			Enum: &EnumField[any]{Items: []string{"perempuan", "laki-laki"}},
		},
	)
	if err != nil {
		t.Errorf("Test case 2 Error : %v", err)
	}

	_, err = Validate(
		"jenis_kelamin", payload, RequestDataValidator{
			Null: false,
			Enum: &EnumField[any]{Items: []string{"bola", "badminton", "renang"}},
		},
	)
	if err == nil {
		t.Errorf("Test case 3 Error : this sould be error")
	}

}

func TestMultipleValidation(t *testing.T) {
	payload := map[string]interface{}{"jenis_kelamin": "laki-laki", "hoby": "Main PS", "umur": 1, "menikah": true}
	err := MultiValidate(payload, map[string]RequestDataValidator{
		"jenis_kelamin": {Enum: &EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"hoby":          {Type: reflect.String, Null: false},
		"menikah":       {Type: reflect.Bool, Null: false},
	})
	if err != nil {
		t.Errorf("Expected not have error, but got error : %s", err)
	}

	err = MultiValidate(payload, map[string]RequestDataValidator{
		"jenis_kelamin": {Enum: &EnumField[any]{Items: []string{"laki-laki", "perempuan"}}},
		"hoby":          {Type: reflect.Int, Null: false},
		"menikah":       {Type: reflect.Bool, Null: false},
	})
	if err == nil {
		t.Error("Expected have an error, but you got no error")
	}
}
