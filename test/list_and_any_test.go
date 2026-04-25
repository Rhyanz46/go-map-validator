package test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
	"github.com/google/uuid"
)

// =====================================================================
// List(elem) shortcut — primitive list helper
// =====================================================================

type listOfStrings struct {
	Tags []string `json:"tags"`
}

type listOfInts struct {
	IDs []int `json:"ids"`
}

// TestList_String_HappyPath: List(Str()) accepts []string and binds correctly.
func TestList_String_HappyPath(t *testing.T) {
	body := `{"tags": ["go", "validator", "test"]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str())).
		Done()

	got, err := map_validator.ValidateJSON[listOfStrings](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Tags) != 3 {
		t.Fatalf("want 3 tags, got %d", len(got.Tags))
	}
	if got.Tags[0] != "go" || got.Tags[1] != "validator" || got.Tags[2] != "test" {
		t.Errorf("unexpected tags: %v", got.Tags)
	}
}

// TestList_Int_HappyPath: List(Int()) accepts []int and binds correctly.
func TestList_Int_HappyPath(t *testing.T) {
	body := `{"ids": [1, 2, 3, 4]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("ids", map_validator.List(map_validator.Int())).
		Done()

	got, err := map_validator.ValidateJSON[listOfInts](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.IDs) != 4 {
		t.Fatalf("want 4 ids, got %d", len(got.IDs))
	}
}

// TestList_ElementMax_FailsWhenExceeded: WithMax on element rule constrains item length.
func TestList_ElementMax_FailsWhenExceeded(t *testing.T) {
	body := `{"tags": ["ok", "way-too-long-string"]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str().WithMax(5))).
		Done()

	_, err := map_validator.ValidateJSON[listOfStrings](req, rules)
	if err == nil {
		t.Fatal("expected error for element exceeding max length")
	}
	if !strings.Contains(err.Error(), "tags") {
		t.Errorf("expected error mentioning 'tags', got %q", err.Error())
	}
}

// TestList_ElementMin_FailsWhenBelow: WithMin on element rule.
func TestList_ElementMin_FailsWhenBelow(t *testing.T) {
	body := `{"tags": ["okay", "no"]}` // "no" below min 3
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str().WithMin(3))).
		Done()

	_, err := map_validator.ValidateJSON[listOfStrings](req, rules)
	if err == nil {
		t.Fatal("expected error for element below min length")
	}
}

// TestList_ContainerMin_FailsWhenEmpty: WithMin on List itself = list count constraint.
func TestList_ContainerMin_FailsWhenEmpty(t *testing.T) {
	body := `{"tags": []}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str()).WithMin(1)).
		Done()

	_, err := map_validator.ValidateJSON[listOfStrings](req, rules)
	if err == nil {
		t.Fatal("expected error for empty list when container min=1")
	}
}

// TestList_ContainerMax_FailsWhenExceeded: WithMax on List = list count limit.
func TestList_ContainerMax_FailsWhenExceeded(t *testing.T) {
	body := `{"tags": ["a", "b", "c", "d", "e"]}` // 5 items vs max 3
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str()).WithMax(3)).
		Done()

	_, err := map_validator.ValidateJSON[listOfStrings](req, rules)
	if err == nil {
		t.Fatal("expected error for list exceeding container max")
	}
}

// TestList_ElementAndContainerConstraints: both kinds of constraint applied independently.
func TestList_ElementAndContainerConstraints(t *testing.T) {
	body := `{"tags": ["go", "rs"]}` // 2 items, each 2 chars
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str().Between(1, 3)).Between(1, 5)).
		Done()

	got, err := map_validator.ValidateJSON[listOfStrings](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Tags) != 2 {
		t.Errorf("want 2 tags, got %d", len(got.Tags))
	}
}

// TestList_EmptyAllowedByDefault: list without container Min accepts empty.
func TestList_EmptyAllowedByDefault(t *testing.T) {
	body := `{"tags": []}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str())).
		Done()

	got, err := map_validator.ValidateJSON[listOfStrings](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Tags) != 0 {
		t.Errorf("want 0 tags, got %d", len(got.Tags))
	}
}

// TestList_Nullable_AllowsMissingField: Nullable allows the whole field to be absent.
func TestList_Nullable_AllowsMissingField(t *testing.T) {
	body := `{}` // tags missing entirely
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str()).Nullable()).
		Done()

	_, err := map_validator.ValidateJSON[listOfStrings](req, rules)
	if err != nil {
		t.Fatalf("Nullable list should allow missing field, got: %v", err)
	}
}

// TestList_RequiredWhenNotNullable: missing list field without Nullable → error.
func TestList_RequiredWhenNotNullable(t *testing.T) {
	body := `{}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str())).
		Done()

	_, err := map_validator.ValidateJSON[listOfStrings](req, rules)
	if err == nil {
		t.Fatal("expected error for missing required list field")
	}
}

// =====================================================================
// Any() escape hatch — passthrough field
// =====================================================================

type withMetadata struct {
	Title    string                 `json:"title"`
	Metadata map[string]interface{} `json:"metadata"`
}

type withAnything struct {
	Anything interface{} `json:"anything"`
}

// TestAny_PreservesMap: Any() lets a heterogen map survive validation and bind.
func TestAny_PreservesMap(t *testing.T) {
	body := `{"title": "hello", "metadata": {"foo": "bar", "n": 42, "nested": {"k": "v"}}}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("title", map_validator.Str()).
		SetRule("metadata", map_validator.Any()).
		Done()

	got, err := map_validator.ValidateJSON[withMetadata](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Metadata == nil {
		t.Fatal("metadata should not be nil after Any() passthrough")
	}
	if got.Metadata["foo"] != "bar" {
		t.Errorf("metadata.foo: want 'bar', got %v", got.Metadata["foo"])
	}
	nested, ok := got.Metadata["nested"].(map[string]interface{})
	if !ok {
		t.Fatalf("nested not preserved as map: %T", got.Metadata["nested"])
	}
	if nested["k"] != "v" {
		t.Errorf("nested.k: want 'v', got %v", nested["k"])
	}
}

// TestAny_PreservesArray: Any() preserves heterogen JSON array.
func TestAny_PreservesArray(t *testing.T) {
	body := `{"anything": [1, "two", true, null]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("anything", map_validator.Any()).
		Done()

	got, err := map_validator.ValidateJSON[withAnything](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr, ok := got.Anything.([]interface{})
	if !ok {
		t.Fatalf("anything not preserved as slice: %T", got.Anything)
	}
	if len(arr) != 4 {
		t.Errorf("want 4 elements, got %d", len(arr))
	}
}

// TestAny_PreservesPrimitive: Any() also passes through primitives.
func TestAny_PreservesPrimitive(t *testing.T) {
	body := `{"anything": "hello"}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("anything", map_validator.Any()).
		Done()

	got, err := map_validator.ValidateJSON[withAnything](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Anything != "hello" {
		t.Errorf("want 'hello', got %v", got.Anything)
	}
}

// TestAny_RequiresFieldByDefault: bare Any() requires the field to be present.
func TestAny_RequiresFieldByDefault(t *testing.T) {
	body := `{}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("metadata", map_validator.Any()).
		Done()

	_, err := map_validator.ValidateJSON[withMetadata](req, rules)
	if err == nil {
		t.Fatal("bare Any() should require field by default")
	}
}

// TestAny_NullableAllowsMissingField: Any().Nullable() makes field optional.
func TestAny_NullableAllowsMissingField(t *testing.T) {
	body := `{}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("metadata", map_validator.Any().Nullable()).
		Done()

	_, err := map_validator.ValidateJSON[withMetadata](req, rules)
	if err != nil {
		t.Fatalf("Any().Nullable() should allow missing field, got: %v", err)
	}
}

// =====================================================================
// Regression — silent-drop of undeclared fields IS the documented
// whitelist behavior. List() and Any() are the escape hatches.
// =====================================================================

type silentDropRegression struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

// TestRegression_UndeclaredFieldDropped: documents intentional whitelist binding.
// Without a rule for a struct field, Bind output zero-values that field.
func TestRegression_UndeclaredFieldDropped(t *testing.T) {
	body := `{"title": "hello", "tags": ["a", "b"]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("title", map_validator.Str()).
		Done() // tags intentionally NOT declared

	got, err := map_validator.ValidateJSON[silentDropRegression](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title != "hello" {
		t.Errorf("title: want 'hello', got %q", got.Title)
	}
	if got.Tags != nil {
		t.Errorf("undeclared field should be dropped (whitelist binding), got %v", got.Tags)
	}
}

// TestRegression_DeclaredListPreserved: rule via List() keeps the field.
func TestRegression_DeclaredListPreserved(t *testing.T) {
	body := `{"title": "hello", "tags": ["a", "b"]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("title", map_validator.Str()).
		SetRule("tags", map_validator.List(map_validator.Str())).
		Done()

	got, err := map_validator.ValidateJSON[silentDropRegression](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Tags) != 2 {
		t.Errorf("want 2 tags, got %d: %v", len(got.Tags), got.Tags)
	}
}

// TestRegression_DeclaredAnyPreserved: rule via Any() keeps the field as raw value.
func TestRegression_DeclaredAnyPreserved(t *testing.T) {
	type anyTagsStruct struct {
		Title string      `json:"title"`
		Tags  interface{} `json:"tags"`
	}
	body := `{"title": "hello", "tags": ["a", "b"]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("title", map_validator.Str()).
		SetRule("tags", map_validator.Any()).
		Done()

	got, err := map_validator.ValidateJSON[anyTagsStruct](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Tags == nil {
		t.Fatal("tags should be preserved by Any()")
	}
	arr, ok := got.Tags.([]interface{})
	if !ok {
		t.Fatalf("tags not slice: %T", got.Tags)
	}
	if len(arr) != 2 {
		t.Errorf("want 2 elements, got %d", len(arr))
	}
}

// =====================================================================
// CORE — production scenarios from feedback
// =====================================================================

type listOfUUIDStrings struct {
	IDs []string `json:"ids"`
}

// TestList_UUID — List(UUID()) accepts []string of valid UUIDs.
func TestList_UUID(t *testing.T) {
	body := `{"ids": ["123e4567-e89b-12d3-a456-426614174000", "550e8400-e29b-41d4-a716-446655440000"]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("ids", map_validator.List(map_validator.UUID())).
		Done()

	got, err := map_validator.ValidateJSON[listOfUUIDStrings](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.IDs) != 2 {
		t.Fatalf("want 2 ids, got %d", len(got.IDs))
	}
	if got.IDs[0] != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("ids[0]=%q", got.IDs[0])
	}
}

// TestList_UUID_InvalidFails — invalid UUID in list → error.
func TestList_UUID_InvalidFails(t *testing.T) {
	body := `{"ids": ["123e4567-e89b-12d3-a456-426614174000", "not-a-uuid"]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("ids", map_validator.List(map_validator.UUID())).
		Done()

	_, err := map_validator.ValidateJSON[listOfUUIDStrings](req, rules)
	if err == nil {
		t.Fatal("expected error for invalid UUID in list")
	}
}

// TestList_StrEnum — List(StrEnum(...)) constrains element values.
func TestList_StrEnum(t *testing.T) {
	type colors struct {
		Colors []string `json:"colors"`
	}
	rules := map_validator.BuildRoles().
		SetRule("colors", map_validator.List(map_validator.StrEnum("red", "blue", "green"))).
		Done()

	// happy path
	body := `{"colors": ["red", "blue"]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	got, err := map_validator.ValidateJSON[colors](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Colors) != 2 || got.Colors[0] != "red" {
		t.Errorf("colors=%v", got.Colors)
	}

	// invalid value
	bad := `{"colors": ["red", "purple"]}`
	badReq := httptest.NewRequest("POST", "/test", bytes.NewBufferString(bad))
	if _, err := map_validator.ValidateJSON[colors](badReq, rules); err == nil {
		t.Error("expected error for value not in enum")
	}
}

// TestList_BodyFieldNotArray — single string instead of array → error.
func TestList_BodyFieldNotArray(t *testing.T) {
	body := `{"tags": "single-string"}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str())).
		Done()

	_, err := map_validator.ValidateJSON[listOfStrings](req, rules)
	if err == nil {
		t.Fatal("expected error: body field is not array")
	}
}

// TestList_MixedTypesInArray — array with wrong element type → error.
func TestList_MixedTypesInArray(t *testing.T) {
	body := `{"ids": [1, 2, "three"]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("ids", map_validator.List(map_validator.Int())).
		Done()

	_, err := map_validator.ValidateJSON[listOfInts](req, rules)
	if err == nil {
		t.Fatal("expected error: element type mismatch")
	}
}

// TestList_CustomMessage — element-level custom message propagates.
func TestList_CustomMessage(t *testing.T) {
	body := `{"tags": ["abcdef"]}` // exceeds max 3
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str().WithMax(3).WithMsg(map_validator.CustomMsg{
			OnMax: map_validator.SetMessage("element too long: ${actual_length} chars"),
		}))).
		Done()

	_, err := map_validator.ValidateJSON[listOfStrings](req, rules)
	if err == nil {
		t.Fatal("expected error from custom message")
	}
	if !strings.Contains(err.Error(), "element too long") {
		t.Errorf("expected custom message in error, got %q", err.Error())
	}
}

// TestList_ConcurrentSharedRules — 20 goroutines paralel pakai rule yang sama.
func TestList_ConcurrentSharedRules(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str().WithMax(20))).
		Done()

	const workers = 20
	var wg sync.WaitGroup
	errs := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			body := `{"tags": ["alpha", "beta", "gamma"]}`
			req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")

			got, err := map_validator.ValidateJSON[listOfStrings](req, rules)
			if err != nil {
				errs <- err
				return
			}
			if len(got.Tags) != 3 {
				errs <- fmt.Errorf("want 3 tags, got %d", len(got.Tags))
			}
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("concurrent call failed: %v", err)
	}
}

// TestAny_ExplicitNullValue — explicit null vs missing field semantics.
func TestAny_ExplicitNullValue(t *testing.T) {
	type nullableMeta struct {
		Metadata interface{} `json:"metadata"`
	}

	// Any() (not nullable) + explicit null → error
	body := `{"metadata": null}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	rules := map_validator.BuildRoles().
		SetRule("metadata", map_validator.Any()).
		Done()
	if _, err := map_validator.ValidateJSON[nullableMeta](req, rules); err == nil {
		t.Error("Any() (required) + null body should error")
	}

	// Any().Nullable() + explicit null → success, bind to nil
	rulesNullable := map_validator.BuildRoles().
		SetRule("metadata", map_validator.Any().Nullable()).
		Done()
	req2 := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	got, err := map_validator.ValidateJSON[nullableMeta](req2, rulesNullable)
	if err != nil {
		t.Fatalf("Any().Nullable() + null body should succeed, got: %v", err)
	}
	if got.Metadata != nil {
		t.Errorf("metadata should be nil, got %v", got.Metadata)
	}
}

// TestList_InsideNestedObject — search.tags pattern.
func TestList_InsideNestedObject(t *testing.T) {
	type search struct {
		Filter struct {
			Query string   `json:"query"`
			Tags  []string `json:"tags"`
		} `json:"filter"`
	}

	filterRules := map_validator.BuildRoles().
		SetRule("query", map_validator.Str()).
		SetRule("tags", map_validator.List(map_validator.Str())).
		Done()
	rules := map_validator.BuildRoles().
		SetRule("filter", map_validator.NestedObject(filterRules)).
		Done()

	body := `{"filter": {"query": "hello", "tags": ["go", "rust"]}}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	got, err := map_validator.ValidateJSON[search](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Filter.Query != "hello" {
		t.Errorf("query=%q", got.Filter.Query)
	}
	if len(got.Filter.Tags) != 2 {
		t.Errorf("want 2 tags, got %d", len(got.Filter.Tags))
	}
}

// TestRegression_MassAssignmentSafety — undeclared field doesn't leak via bind.
func TestRegression_MassAssignmentSafety(t *testing.T) {
	type user struct {
		Name    string `json:"name"`
		IsAdmin bool   `json:"is_admin"`
	}
	rules := map_validator.BuildRoles().
		SetRule("name", map_validator.Str()).
		Done() // is_admin intentionally not declared

	body := `{"name": "user", "is_admin": true}` // attacker injects is_admin
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	got, err := map_validator.ValidateJSON[user](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "user" {
		t.Errorf("name=%q", got.Name)
	}
	if got.IsAdmin {
		t.Error("SECURITY: undeclared is_admin should NOT bind from body (whitelist binding)")
	}
}

// =====================================================================
// WORTH-ADDING — element regex, email list, null element, typed slice
// =====================================================================

// TestList_ElementRegex — element-level regex validation.
func TestList_ElementRegex(t *testing.T) {
	rules := map_validator.BuildRoles().
		SetRule("slugs", map_validator.List(map_validator.Str().Regex("^[a-z]+$"))).
		Done()

	// happy
	good := `{"slugs": ["abc", "hello"]}`
	type slugged struct {
		Slugs []string `json:"slugs"`
	}
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(good))
	if _, err := map_validator.ValidateJSON[slugged](req, rules); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// fail (uppercase not allowed)
	bad := `{"slugs": ["abc", "ABC"]}`
	badReq := httptest.NewRequest("POST", "/test", bytes.NewBufferString(bad))
	if _, err := map_validator.ValidateJSON[slugged](badReq, rules); err == nil {
		t.Error("expected error for element failing regex")
	}
}

// TestList_Email — List(Email()) for arrays of emails.
func TestList_Email(t *testing.T) {
	type emails struct {
		Recipients []string `json:"recipients"`
	}
	rules := map_validator.BuildRoles().
		SetRule("recipients", map_validator.List(map_validator.Email())).
		Done()

	good := `{"recipients": ["a@b.com", "c@d.org"]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(good))
	got, err := map_validator.ValidateJSON[emails](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Recipients) != 2 {
		t.Errorf("want 2 recipients, got %d", len(got.Recipients))
	}

	bad := `{"recipients": ["a@b.com", "not-email"]}`
	badReq := httptest.NewRequest("POST", "/test", bytes.NewBufferString(bad))
	if _, err := map_validator.ValidateJSON[emails](badReq, rules); err == nil {
		t.Error("expected error for invalid email in list")
	}
}

// TestList_NullElementInArray — null element semantics. Locks contract.
func TestList_NullElementInArray(t *testing.T) {
	body := `{"tags": ["a", null, "c"]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rules := map_validator.BuildRoles().
		SetRule("tags", map_validator.List(map_validator.Str())).
		Done()

	// Contract: null element in non-nullable list should error.
	_, err := map_validator.ValidateJSON[listOfStrings](req, rules)
	if err == nil {
		t.Fatal("expected error for null element in list (current contract)")
	}
}

// TestList_BindToTypedSlice — bind to []uuid.UUID typed (not []string).
func TestList_BindToTypedSlice(t *testing.T) {
	type withTypedIDs struct {
		IDs []uuid.UUID `json:"ids"`
	}
	rules := map_validator.BuildRoles().
		SetRule("ids", map_validator.List(map_validator.UUID())).
		Done()

	body := `{"ids": ["123e4567-e89b-12d3-a456-426614174000", "550e8400-e29b-41d4-a716-446655440000"]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	got, err := map_validator.ValidateJSON[withTypedIDs](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.IDs) != 2 {
		t.Fatalf("want 2 ids, got %d", len(got.IDs))
	}
	expected, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if got.IDs[0] != expected {
		t.Errorf("ids[0]=%v, want %v", got.IDs[0], expected)
	}
}

// =====================================================================
// NESTED COMPOSITIONS
// =====================================================================

// TestList_InsideListOfObject — order.goods[i].tags pattern.
func TestList_InsideListOfObject(t *testing.T) {
	type goodsItem struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}
	type order struct {
		Goods []goodsItem `json:"goods"`
	}

	itemRules := map_validator.BuildRoles().
		SetRule("name", map_validator.Str()).
		SetRule("tags", map_validator.List(map_validator.Str())).
		Done()
	rules := map_validator.BuildRoles().
		SetRule("goods", map_validator.ListOfObject(itemRules)).
		Done()

	body := `{"goods": [
		{"name": "Apple", "tags": ["fresh", "red"]},
		{"name": "Banana", "tags": ["yellow"]}
	]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	got, err := map_validator.ValidateJSON[order](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Goods) != 2 {
		t.Fatalf("want 2 goods, got %d", len(got.Goods))
	}
	if len(got.Goods[0].Tags) != 2 || got.Goods[0].Tags[0] != "fresh" {
		t.Errorf("goods[0].tags=%v", got.Goods[0].Tags)
	}
	if len(got.Goods[1].Tags) != 1 || got.Goods[1].Tags[0] != "yellow" {
		t.Errorf("goods[1].tags=%v", got.Goods[1].Tags)
	}
}

// TestAny_InsideNestedObject — config.settings (heterogen sub-object).
func TestAny_InsideNestedObject(t *testing.T) {
	type config struct {
		Config struct {
			Name     string                 `json:"name"`
			Settings map[string]interface{} `json:"settings"`
		} `json:"config"`
	}

	configRules := map_validator.BuildRoles().
		SetRule("name", map_validator.Str()).
		SetRule("settings", map_validator.Any()).
		Done()
	rules := map_validator.BuildRoles().
		SetRule("config", map_validator.NestedObject(configRules)).
		Done()

	body := `{"config": {"name": "webhook", "settings": {"timeout": 30, "retries": 3}}}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	got, err := map_validator.ValidateJSON[config](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Config.Name != "webhook" {
		t.Errorf("name=%q", got.Config.Name)
	}
	if got.Config.Settings == nil {
		t.Fatal("settings should not be nil")
	}
	if got.Config.Settings["timeout"] != float64(30) {
		t.Errorf("settings.timeout=%v", got.Config.Settings["timeout"])
	}
}

// TestAny_InsideListOfObject — events[].metadata (per-item heterogen).
func TestAny_InsideListOfObject(t *testing.T) {
	type event struct {
		ID       string                 `json:"id"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	type events struct {
		Events []event `json:"events"`
	}

	itemRules := map_validator.BuildRoles().
		SetRule("id", map_validator.UUID()).
		SetRule("metadata", map_validator.Any()).
		Done()
	rules := map_validator.BuildRoles().
		SetRule("events", map_validator.ListOfObject(itemRules)).
		Done()

	body := `{"events": [
		{"id": "123e4567-e89b-12d3-a456-426614174000", "metadata": {"source": "api"}},
		{"id": "550e8400-e29b-41d4-a716-446655440000", "metadata": {"source": "ui", "extra": [1,2,3]}}
	]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	got, err := map_validator.ValidateJSON[events](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Events) != 2 {
		t.Fatalf("want 2 events, got %d", len(got.Events))
	}
	if got.Events[0].Metadata["source"] != "api" {
		t.Errorf("events[0].metadata.source=%v", got.Events[0].Metadata["source"])
	}
	if got.Events[1].Metadata["extra"] == nil {
		t.Error("events[1].metadata.extra should be preserved")
	}
}

// TestRegression_NestedWhitelistDrop — nested level also strips undeclared.
func TestRegression_NestedWhitelistDrop(t *testing.T) {
	type item struct {
		Name        string `json:"name"`
		LeakedField string `json:"leaked_field"`
	}
	type wrapper struct {
		Items []item `json:"items"`
	}

	itemRules := map_validator.BuildRoles().
		SetRule("name", map_validator.Str()).
		Done() // leaked_field intentionally not declared
	rules := map_validator.BuildRoles().
		SetRule("items", map_validator.ListOfObject(itemRules)).
		Done()

	body := `{"items": [{"name": "x", "leaked_field": "danger"}, {"name": "y", "leaked_field": "more-danger"}]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	got, err := map_validator.ValidateJSON[wrapper](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Items) != 2 {
		t.Fatalf("want 2 items, got %d", len(got.Items))
	}
	for i, it := range got.Items {
		if it.LeakedField != "" {
			t.Errorf("items[%d].leaked_field should be stripped, got %q", i, it.LeakedField)
		}
	}
}

// =====================================================================
// DEEP NESTED — 3+ level compositions
// =====================================================================

// TestNested_ListOfObject_InsideNested — search.filter.matches[] pattern.
func TestNested_ListOfObject_InsideNested(t *testing.T) {
	type match struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	type filter struct {
		Query   string  `json:"query"`
		Matches []match `json:"matches"`
	}
	type search struct {
		Filter filter `json:"filter"`
	}

	matchRules := map_validator.BuildRoles().
		SetRule("id", map_validator.Int()).
		SetRule("name", map_validator.Str()).
		Done()
	filterRules := map_validator.BuildRoles().
		SetRule("query", map_validator.Str()).
		SetRule("matches", map_validator.ListOfObject(matchRules)).
		Done()
	rules := map_validator.BuildRoles().
		SetRule("filter", map_validator.NestedObject(filterRules)).
		Done()

	body := `{"filter": {"query": "hello", "matches": [{"id": 1, "name": "first"}, {"id": 2, "name": "second"}]}}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	got, err := map_validator.ValidateJSON[search](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Filter.Query != "hello" {
		t.Errorf("query=%q", got.Filter.Query)
	}
	if len(got.Filter.Matches) != 2 {
		t.Fatalf("want 2 matches, got %d", len(got.Filter.Matches))
	}
	if got.Filter.Matches[0].ID != 1 || got.Filter.Matches[0].Name != "first" {
		t.Errorf("matches[0]=%+v", got.Filter.Matches[0])
	}
}

// TestNested_NestedObject_InsideListItem — events[].context pattern.
func TestNested_NestedObject_InsideListItem(t *testing.T) {
	type ctx struct {
		UserID string `json:"user_id"`
		IP     string `json:"ip"`
	}
	type event struct {
		Type    string `json:"type"`
		Context ctx    `json:"context"`
	}
	type wrapper struct {
		Events []event `json:"events"`
	}

	ctxRules := map_validator.BuildRoles().
		SetRule("user_id", map_validator.UUID()).
		SetRule("ip", map_validator.IPv4()).
		Done()
	eventRules := map_validator.BuildRoles().
		SetRule("type", map_validator.Str()).
		SetRule("context", map_validator.NestedObject(ctxRules)).
		Done()
	rules := map_validator.BuildRoles().
		SetRule("events", map_validator.ListOfObject(eventRules)).
		Done()

	body := `{"events": [
		{"type": "login", "context": {"user_id": "123e4567-e89b-12d3-a456-426614174000", "ip": "192.168.1.1"}},
		{"type": "logout", "context": {"user_id": "550e8400-e29b-41d4-a716-446655440000", "ip": "10.0.0.5"}}
	]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	got, err := map_validator.ValidateJSON[wrapper](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Events) != 2 {
		t.Fatalf("want 2 events, got %d", len(got.Events))
	}
	if got.Events[0].Type != "login" {
		t.Errorf("events[0].type=%q", got.Events[0].Type)
	}
	if got.Events[0].Context.IP != "192.168.1.1" {
		t.Errorf("events[0].context.ip=%q", got.Events[0].Context.IP)
	}
}

// TestRegression_DeepWhitelistDrop — undeclared field stripped at level 3+.
func TestRegression_DeepWhitelistDrop(t *testing.T) {
	type meta struct {
		Source      string `json:"source"`
		LeakedField string `json:"leaked_field"`
	}
	type item struct {
		Name     string `json:"name"`
		Metadata meta   `json:"metadata"`
	}
	type wrapper struct {
		Items []item `json:"items"`
	}

	metaRules := map_validator.BuildRoles().
		SetRule("source", map_validator.Str()).
		Done() // leaked_field at level 3 intentionally undeclared
	itemRules := map_validator.BuildRoles().
		SetRule("name", map_validator.Str()).
		SetRule("metadata", map_validator.NestedObject(metaRules)).
		Done()
	rules := map_validator.BuildRoles().
		SetRule("items", map_validator.ListOfObject(itemRules)).
		Done()

	body := `{"items": [
		{"name": "x", "metadata": {"source": "api", "leaked_field": "danger-deep"}}
	]}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	got, err := map_validator.ValidateJSON[wrapper](req, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("want 1 item, got %d", len(got.Items))
	}
	if got.Items[0].Metadata.Source != "api" {
		t.Errorf("metadata.source=%q", got.Items[0].Metadata.Source)
	}
	if got.Items[0].Metadata.LeakedField != "" {
		t.Errorf("SECURITY: deep undeclared field should be stripped, got %q", got.Items[0].Metadata.LeakedField)
	}
}
