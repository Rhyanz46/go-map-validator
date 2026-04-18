package test

import (
	"bytes"
	"fmt"
	"net/http/httptest"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

// Example of using ValidateJSON to validate an HTTP JSON request and bind
// the result into a typed struct in a single call.
func ExampleValidateJSON() {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	rules := map_validator.BuildRoles().
		SetRule("email", map_validator.Email().WithMax(255)).
		SetRule("password", map_validator.Str().Between(6, 64)).
		Done()

	body := `{"email": "dev@example.com", "password": "secret123"}`
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	req1, err := map_validator.ValidateJSON[LoginRequest](req, rules)
	if err != nil {
		fmt.Println("validation failed:", err)
		return
	}
	fmt.Println(req1.Email)
	// Output: dev@example.com
}
