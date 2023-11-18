package test

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"testing"
)

func TestUUID(t *testing.T) {
	payload := "69b9f8e0-c5a5-41e8-9671-b169058ce4bd"
	badPayload1 := "550e8400-e29b-41d4-a716-44665544000z"
	badPayload := 1
	invalidUUIDs := []string{
		"550e8400-e29b-41d4-a716-44665544000z",
		"abcdefg-hijk-lmno-pqrs-tuvwxyz01234",
		"not-uuid-format-at-all",
		"123e4567-e89b-cdef-ghij-123456789012",
		"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaaa",
		"01234567-89ab-cdef-ghij-ijklmnopqrst",
		"invalid-uuid-string-1234",
		"xyz12345-6789-abcd-efgh-ijklmnopqrst",
	}
	validUUIDs := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"6fa459ea-ee8a-3ca4-894e-db77e160355e",
		"f47ac10b-58cc-4372-a567-0e02b2c3d479",
		"123e4567-e89b-cdef-0123-456789abcdef",
		"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		"01234567-89ab-cdef-0123-456789abcdef",
		"abcdefab-cdef-abcd-efab-cdefabcdefab",
		"00000000-0000-0000-0000-000000000000",
		"deadbeef-0123-4567-89ab-deadbeef0123",
		"01234567-89ab-cdef-0123-456789abcdef",
	}
	validRole := map_validator.Rules{UUID: true}
	check, err := map_validator.NewValidateBuilder().SetRule(validRole).Load(badPayload)
	if err != nil {
		t.Errorf("Expected have an error, but you got no error : %s", err)
	}
	_, err = check.RunValidate()
	expect := "the field 'data' it's not valid uuid"
	if err.Error() != expect {
		t.Errorf("Expected %s, but you got no error : %s", expect, err)
	}

	check, err = map_validator.NewValidateBuilder().SetRule(validRole).Load(badPayload1)
	if err != nil {
		t.Errorf("Expected have an error, but you got no error : %s", err)
	}
	_, err = check.RunValidate()
	expect = "the field 'data' it's not valid uuid"
	if err.Error() != expect {
		t.Errorf("Expected %s, but you got no error : %s", expect, err)
	}

	for _, uuidInvalid := range invalidUUIDs {
		check, err := map_validator.NewValidateBuilder().SetRule(validRole).Load(uuidInvalid)
		if err != nil {
			t.Errorf("Expected have an error, but you got no error : %s", err)
		}
		_, err = check.RunValidate()
		expect = "the field 'data' it's not valid uuid"
		if err.Error() != expect {
			t.Errorf("Expected %s, but you got no error : %s", expect, err)
			break
		}
	}

	for _, uuidValid := range validUUIDs {
		check, err := map_validator.NewValidateBuilder().SetRule(validRole).Load(uuidValid)
		if err != nil {
			t.Errorf("Expected have an error, but you got no error : %s", err)
		}
		_, err = check.RunValidate()
		expect = "the field 'data' it's not valid uuid"
		if err != nil {
			t.Errorf("Expected have an error, but you got no error : %s", err)
		}
	}

	check, err = map_validator.NewValidateBuilder().SetRule(validRole).Load(payload)
	if err != nil {
		t.Errorf("Expected have an error, but you got no error : %s", err)
	}
	_, err = check.RunValidate()
	if err != nil {
		t.Errorf("Expected have an error, but you got no error : %s", err)
	}
}
