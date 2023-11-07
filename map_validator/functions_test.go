package map_validator

import (
	"reflect"
	"testing"
)

func TestValidateUUID(t *testing.T) {
	payload := map[string]interface{}{
		"field1": "123e4567-e89b-12d3-a456-426614174001",
	}

	validator := Rules{
		UUID: true,
		Null: false,
	}

	_, err := validate("field1", payload, validator, fromHttpJson)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateInvalidUUID(t *testing.T) {
	payload := map[string]interface{}{
		"field1": "invalid-uuid",
	}

	validator := Rules{
		UUID: true,
		Null: false,
	}

	_, err := validate("field1", payload, validator, fromHttpJson)

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

	validator := Rules{
		Type: reflect.String,
		Null: false,
	}

	_, err := validate("field1", payload, validator, fromHttpJson)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateNullNotAllowed(t *testing.T) {
	payload := map[string]interface{}{
		"field1": nil,
	}

	validator := Rules{
		Type: reflect.String,
		Null: false,
	}

	_, err := validate("field1", payload, validator, fromHttpJson)

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

	validator := Rules{
		Type: reflect.String,
		Max:  SetTotal(5),
	}

	_, err := validate("field1", payload, validator, fromHttpJson)

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

	validator := Rules{
		Type: reflect.String,
		Min:  SetTotal(5),
	}

	_, err := validate("field1", payload, validator, fromHttpJson)

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

	validator := Rules{
		Email: true,
	}

	_, err := validate("email", payload, validator, fromHttpJson)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateInvalidEmail(t *testing.T) {
	payload := map[string]interface{}{
		"email": "invalid-email",
	}

	validator := Rules{
		Email: true,
	}

	_, err := validate("email", payload, validator, fromHttpJson)

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

	validator := Rules{
		IPV4: true,
	}

	_, err := validate("ip_address", payload, validator, fromHttpJson)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateInvalidIPV4(t *testing.T) {
	payload := map[string]interface{}{
		"ip_address": "invalid-ip",
	}

	validator := Rules{
		IPV4: true,
	}

	_, err := validate("ip_address", payload, validator, fromHttpJson)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	expectedError := "the field 'ip_address' it's not valid IP"
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expected error: %s, got: %v", expectedError, err)
	}
}

func TestEnumFieldCheck(t *testing.T) {
	payload := map[string]interface{}{"validatorType": "arian", "jenis_kelamin": "laki-laki", "hoby": "Main PS"}
	_, err := validate(
		"validatorType", payload, Rules{
			Null: false,
			Enum: &EnumField[any]{Items: []string{"arian", "aaa"}},
		}, fromHttpJson,
	)
	if err != nil {
		t.Errorf("Test case 1 Error : %v", err)
	}

	_, err = validate(
		"jenis_kelamin", payload, Rules{
			Null: false,
			Enum: &EnumField[any]{Items: []string{"perempuan", "laki-laki"}},
		}, fromHttpJson,
	)
	if err != nil {
		t.Errorf("Test case 2 Error : %v", err)
	}

	_, err = validate(
		"jenis_kelamin", payload, Rules{
			Null: false,
			Enum: &EnumField[any]{Items: []string{"bola", "badminton", "renang"}},
		}, fromHttpJson,
	)
	if err == nil {
		t.Errorf("Test case 3 Error : this sould be error")
	}

}

func TestIntFamily(t *testing.T) {
	payload := map[string]interface{}{"umur": 1, "harga": 1.3}
	_, err := validate(
		"umur", payload, Rules{
			Type: reflect.Int,
		}, fromMapString,
	)
	if err != nil {
		t.Errorf("Test case 1 Error : %v", err)
	}

	payload = map[string]interface{}{"umur": 1, "harga": 1.3}
	_, err = validate(
		"umur", payload, Rules{
			Type: reflect.Int,
		}, fromHttpJson,
	)
	if err != nil {
		t.Errorf("Test case 1 Error : %v", err)
	}

	_, err = validate(
		"harga", payload, Rules{
			Type: reflect.Float64,
		}, fromHttpJson,
	)
	if err != nil {
		t.Errorf("Test case 2 Error : %v", err)
	}

	payload = map[string]interface{}{"power": 133.3, "harga": 1.3}
	_, err = validate(
		"power", payload, Rules{
			Type: reflect.Int,
		}, fromHttpJson,
	)
	if err != nil {
		t.Errorf("Test case 1 Error : %v", err)
	}

	payload = map[string]interface{}{"power": 133.3, "harga": 1.3}
	_, err = validate(
		"power", payload, Rules{
			Type: reflect.Int16,
		}, fromHttpJson,
	)
	if err != nil {
		t.Errorf("Test case 1 Error : %v", err)
	}

	payload = map[string]interface{}{"power": "133.3", "harga": 1.3}
	_, err = validate(
		"power", payload, Rules{
			Type: reflect.Int16,
		}, fromHttpJson,
	)
	expected := "the field 'power' should be 'int'"
	if err.Error() != expected {
		t.Errorf("Expected : %s But you got : %s", expected, err)
	}

	_, err = validate(
		"harga", payload, Rules{
			Type: reflect.Float64,
		}, fromHttpJson,
	)
	if err != nil {
		t.Errorf("Test case 2 Error : %v", err)
	}

}
