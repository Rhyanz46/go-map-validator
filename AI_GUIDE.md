# Map Validator Guide

> **AI AGENT REFERENCE**: This document serves as a comprehensive guide for implementing map_validator in Go projects. Follow these rules and patterns strictly for consistent, secure, and maintainable validation implementations.

## 🎯 Purpose

This guide provides:
- Complete implementation patterns for map_validator
- Architectural boundaries and rules
- Best practices for validation in Go projects
- Anti-patterns to avoid
- Performance and security considerations

## 🤖 For AI Agents

**When implementing map_validator:**
1. Read this entire guide first
2. Identify the controller layer in the project
3. Follow all rules strictly
4. Ask questions if uncertain
5. Never deviate from documented patterns

### ⚡ Efficiency Mode for AI Agents

**For quick/rapid development:**
- **Skip CustomMsg** - Use default error messages to save time
- **Basic validations only** - Type, Min, Max, required fields
- **Simple error handling** - Just return err.Error()
- **Skip manipulators** - Only add if specifically needed

**For production-quality code:**
- **Add CustomMsg** - For user-facing APIs and complex validations
- **Complete validations** - Include all relevant rules
- **Detailed error messages** - User-friendly explanations
- **Use manipulators** - TrimValidation for all string fields

**Example - Efficiency Mode:**
```go
// Quick implementation - no CustomMsg
.SetRule("name", map_validator.Rules{
    Type: reflect.String,
    Max:  map_validator.SetTotal(255),
})
```

**Example - Quality Mode:**
```go
// Production implementation - with CustomMsg
.SetRule("name", map_validator.Rules{
    Type:        reflect.String,
    Max:         map_validator.SetTotal(255),
    RegexString: constant.RegexExcludeSpecialChar,
    CustomMsg: map_validator.CustomMsg{
        OnMax:         map_validator.SetMessage("Name is too long (max 255 chars)"),
        OnRegexString: common_utils.ToPointer("Name cannot contain special characters"),
    },
})
```

## 📦 Installation

### Required Dependencies

Add this import to your Go file:
```go
import "github.com/Rhyanz46/go-map-validator/map_validator"
```

### Install the Package
```bash
go get github.com/Rhyanz46/go-map-validator/map_validator
```

### Project Structure Setup

**For Custom Utilities (Optional):**
```
your-project/
├── go.mod
├── pkg/
│   └── map_validator_utils/
│       └── utils.go          # Custom validation utilities (optional)
```

**Where to Use map_validator:**
- **REST Controllers** - HTTP request handlers that accept JSON
- **API Handlers** - Any function that processes HTTP requests
- **Web Controllers** - Controllers serving web API endpoints

**Where NOT to use map_validator:**
- Service layer (business logic)
- Repository layer (data access)
- Model/Entity definitions
- Utility functions

### Setup Custom Utilities (Optional)

**ONLY create if you need to use SetManipulator or SetFieldsManipulator. If not using manipulators, skip this step.**

Create `pkg/map_validator_utils/utils.go` (only if needed):
```go
package map_validator_utils

import (
    "errors"
    "your-project/pkg/string_utils" // Replace with your actual string utilities
)

func TrimValidation(data interface{}) (result interface{}, err error) {
    strData, ok := data.(string)
    if !ok {
        return nil, errors.New("data is not string")
    }
    result = string_utils.CleanDashes(string_utils.TrimAndClean(strData))
    return
}
```

**When you need this:**
- When using `.SetManipulator("field", map_validator_utils.TrimValidation)`
- When using `.SetFieldsManipulator([]string{"field1", "field2"}, map_validator_utils.TrimValidation)`

**When you DON'T need this:**
- Basic validation without field manipulation
- When not sanitizing input data
- Quick development/prototype mode

### Common Import Pattern
```go
import (
    "github.com/gin-gonic/gin"
    "github.com/Rhyanz46/go-map-validator/map_validator"
    "reflect" // Required for type definitions
    "your-project/pkg/gin_utils" // For consistent responses

    // Import map_validator_utils ONLY if using SetManipulator
    // "your-project/pkg/map_validator_utils"
)
```

### Verification

Test your installation with a simple validation:
```go
func TestInstallation(c *gin.Context) {
    roles := map_validator.BuildRoles()
        .SetRule("test", map_validator.Rules{
            Type: reflect.String,
            Max:  map_validator.SetTotal(10),
        })

    jsonDataRoles := map_validator.NewValidateBuilder().SetRules(roles)
    jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "map_validator is working!"})
}
```

If the test runs without import errors, your installation is successful!

## ⚠️ IMPORTANT RULES FOR AI/LLM

### 1. **USAGE SCOPE RESTRICTIONS**

**✅ ONLY use `map_validator` for:**
- **HTTP Request Validation** - Validating incoming API requests
- **REST Controllers** - Functions that handle HTTP endpoints
- **API Handlers** - Any function processing HTTP requests with JSON bodies

**❌ NEVER use `map_validator` for:**
- **Service Layer** - Business logic and data processing
- **Repository Layer** - Database operations
- **Model/Entity Validation** - Struct validation at model level
- **Internal Data Processing** - Validation between services
- **Database Constraints** - Data integrity rules

**Key Principle**: `map_validator` is for **INPUT VALIDATION** at the **BOUNDARY** of your system (HTTP layer), not for **BUSINESS VALIDATION** inside your system.

**Examples of Correct Usage:**
```go
// ✅ CORRECT - HTTP handler validating incoming request
func (h *UserHandler) CreateUser(c *gin.Context) {
    roles := map_validator.BuildRoles()...
    // Validate user input here
}

// ✅ CORRECT - API endpoint validating JSON payload
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    roles := map_validator.BuildRoles()...
    // Validate product data here
}
```

**Examples of Wrong Usage:**
```go
// ❌ WRONG - Don't validate in service layer
func (s *UserService) ProcessUser(data User) error {
    // Validation should happen in controller layer
}

// ❌ WRONG - Don't validate in repository
func (r *UserRepository) Save(user User) error {
    // Assume data is already validated
}
```

**If unsure where to place validation, ask yourself:**
1. Is this function handling an HTTP request? → **Use map_validator**
2. Is this processing business logic? → **Don't use map_validator**
3. Is this accessing database? → **Don't use map_validator**

### 2. **VALIDATION PURPOSE**
- **✅ VALID USES**:
  - Validating incoming HTTP requests
  - Sanitizing user input
  - Enforcing request format rules
  - Type checking and conversion
- **❌ INVALID USES**:
  - Business logic validation
  - Data persistence validation
  - Inter-service communication validation
  - Internal data structure validation

### 3. **ARCHITECTURAL BOUNDARIES**
```go
// ✅ CORRECT - In Controller (any controller layer folder)
func (h *restHandler) CreateResource(c *gin.Context) {
    roles := map_validator.BuildRoles()...
    // Validate HTTP request here
}

// ✅ CORRECT - In API Handler
func (h *apiHandler) CreateResource(c *gin.Context) {
    roles := map_validator.BuildRoles()...
    // Validate HTTP request here
}

// ✅ CORRECT - In HTTP Handler
func (h *httpHandler) CreateResource(c *gin.Context) {
    roles := map_validator.BuildRoles()...
    // Validate HTTP request here
}

// ❌ WRONG - In Service/Usecase
func (s *service) ProcessBusinessLogic(data Data) error {
    // Don't validate here!
    // Assume data is already validated
}

// ❌ WRONG - In Repository
func (r *repository) Save(data Data) error {
    // Don't validate here!
    // Assume data is already validated
}
```

### 4. **SHARED VALIDATION ANTI-PATTERN**
```go
// ❌ WRONG - Don't create shared validation functions
func GetProductValidationRules() *map_validator.Roles {
    return map_validator.BuildRoles()...
}

// ✅ CORRECT - Define validation inline in controller
func (h *restHandler) CreateProduct(c *gin.Context) {
    roles := map_validator.BuildRoles()... // Define here
}
```

### 5. **DATA FLOW PRINCIPLE**
```
HTTP Request → Controller (with map_validator) → Clean Data → Service Layer → Repository
                                    ↑
                              Only validation happens here
```

### 6. **VALIDATION VS BUSINESS LOGIC**
- **map_validator** = Input sanitization and format checking
- **Service Layer** = Business rule validation
- Example:
  ```go
  // Controller: Check if email format is valid
  .SetRule("email", map_validator.Rules{
      RegexString: constant.RegexEmail,
  })

  // Service: Check if email is unique in database
  if s.repository.EmailExists(email) {
      return errors.New("email already exists")
  }
  ```

## Table of Contents
1. [Installation](#installation)
2. [Basic Validation Pattern](#basic-validation-pattern)
3. [LoadJsonHttp Method](#loadjsonhttp-method)
4. [Validation Result Processing](#validation-result-processing)
5. [The Bind Method](#the-bind-method)
6. [Direct Data Access (Get Method)](#direct-data-access-get-method)
7. [Advanced Validation Features](#advanced-validation-features)
8. [Field Types and Rules](#field-types-and-rules)
9. [Manipulators](#manipulators)
10. [Custom Messages](#custom-messages)
11. [Complex Validations](#complex-validations)
12. [Settings](#settings)
13. [Error Handling](#error-handling)

## Basic Validation Pattern

Pattern standar yang digunakan di seluruh project:

```go
// 1. Build validation rules
roles := map_validator.BuildRoles()
    .SetRule("field_name", map_validator.Rules{...})
    .SetRule("another_field", map_validator.Rules{...})

// 2. Create validator
jsonDataRoles := map_validator.NewValidateBuilder().SetRules(roles)

// 3. Load JSON from HTTP request
jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)
if err != nil {
    c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
    return
}

// 4. Run validation
jsonData, err := jsonDataValidate.RunValidate()
if err != nil {
    c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
    return
}

// 5. Bind to struct
var requestStruct RequestStruct
jsonData.Bind(&requestStruct)
```

## LoadJsonHttp Method

The `LoadJsonHttp` method is used to load and parse JSON data from HTTP requests. This is a critical step that validates the JSON structure before applying validation rules.

### Basic Usage
```go
// Load JSON from HTTP request
jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)
if err != nil {
    c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
    return
}
```

### What LoadJsonHttp Does:
1. **Parses JSON body** from HTTP request
2. **Validates JSON syntax** - returns error for invalid JSON
3. **Applies strict mode rules** (if enabled)
4. **Prepares data for validation rules**
5. **Returns validation builder** for next step

### Common LoadJsonHttp Errors
```go
jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)
if err != nil {
    // Common errors:
    // - Invalid JSON syntax
    // - Missing required fields (in strict mode)
    // - Unknown fields (in strict mode)
    // - Invalid data types

    c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{
        Message: err.Error(), // Will show specific error
    })
    return
}
```

### LoadJsonHttp with Different HTTP Methods
```go
// For POST request (JSON body)
jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)

// For PUT/PATCH request (JSON body)
jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)

// For GET request - Don't use LoadJsonHttp
// Use query parameters validation instead (not covered by map_validator)
```

### Important Notes
- Always use `c.Request` as parameter
- Must be called after creating validation builder
- Must be called before `RunValidate()`
- Returns error for malformed JSON immediately
- Works with `gin.Context` from Gin framework

## Validation Result Processing

After successful validation, you have two options to access the validated data:
1. **Bind to struct** (recommended for complex requests)
2. **Direct access** (for simple cases)

## The Bind Method

Use the `Bind` method to map validated data to your struct (recommended approach).

### Basic Binding
```go
// Define your request struct
type CreateUserRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Age      int    `json:"age"`
}

// After validation, bind the data with error handling
var request CreateUserRequest
if err := jsonData.Bind(&request); err != nil {
    c.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{
        Message: "Failed to process request data",
    })
    return
}

// Now you can use the validated data
fmt.Printf("Name: %s, Email: %s", request.Name, request.Email)
```

### Common Bind Errors

The `Bind` method can fail in these scenarios:
```go
if err := jsonData.Bind(&request); err != nil {
    // Common bind errors:
    // - Type assertion errors (e.g., string field cannot be converted to int)
    // - Missing required fields in struct
    // - JSON structure mismatch
    // - Invalid data conversion

    c.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{
        Message: "Failed to process request data",
    })
    return
}
```

### Complete Bind Pattern
```go
// After successful validation
jsonData, err := jsonDataValidate.RunValidate()
if err != nil {
    c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
    return
}

// Bind to struct with proper error handling
var request CreateUserRequest
if err := jsonData.Bind(&request); err != nil {
    log.Errorf("Bind error: %v", err) // Log for debugging
    c.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{
        Message: "Failed to process request data",
    })
    return
}

// Now safe to use the request struct
h.userService.CreateUser(request)

### Type Conversion Examples
```go
type CreateRegistryRequest struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Public      bool    `json:"public"`
    LimitSize   float64 `json:"limit_size"`
    ProjectID   string  `json:"project_id"`
}

// After validation and binding
var request CreateRegistryRequest
jsonData.Bind(&request)

// Access the data with correct types
if request.Public {
    // Registry is public
}

if request.LimitSize > 0 {
    // Storage limit is set
}
```

### Binding Nested Objects
```go
type WebhookRequest struct {
    Name     string                 `json:"name"`
    URL      string                 `json:"url"`
    Config   map[string]interface{} `json:"config"`
    Events   []string               `json:"events"`
}

var request WebhookRequest
jsonData.Bind(&request)

// Access nested data
if timeout, ok := request.Config["timeout"].(float64); ok {
    // Use timeout value
}
```

## Direct Data Access (Get Method)

For simple cases or when you need to access specific fields without creating a struct, use the `Get()` method.

### Basic Get Usage
```go
// After validation
jsonData, err := jsonDataValidate.RunValidate()
if err != nil {
    // handle error
}

// Get individual fields
name := jsonData.Get("name").(string)
age := int(jsonData.Get("age").(float64)) // Numbers are float64 by default
isPublic := jsonData.Get("public").(bool)

fmt.Printf("Name: %s, Age: %d, Public: %t", name, age, isPublic)
```

### Type Assertions
```go
// Always assert to the correct type
stringValue := jsonData.Get("field_name").(string)
intValue := int(jsonData.Get("number_field").(float64))
boolValue := jsonData.Get("bool_field").(bool)

// For optional fields, check if nil
if jsonData.Get("optional_field") != nil {
    optionalValue := jsonData.Get("optional_field").(string)
}
```

## Advanced Validation Features

### 1. **RequiredIf - Conditional Requirements**
Field is required only when another field is present:
```go
.SetRule("password", map_validator.Rules{
    Type:      reflect.String,
    Min:       map_validator.SetTotal(8),
    Null:      true,
    RequiredIf: []string{"confirm_password"}, // Required if confirm_password exists
})
```

### 2. **IfNull - Default Values**
Provide default values for optional fields:
```go
.SetRule("color", map_validator.Rules{
    Type:        reflect.String,
    RegexString: `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`,
    Max:         map_validator.SetTotal(7),
    Null:        true,
    IfNull:      "#000000", // Default hex color
})

.SetRule("timeout", map_validator.Rules{
    Type:   reflect.Int,
    Min:    map_validator.SetTotal(5),
    Max:    map_validator.SetTotal(300),
    Null:   true,
    IfNull: 30, // Default timeout: 30 seconds
})
```

### 3. **Unique - Field Comparison**
Ensure field value is different from another field:
```go
.SetRule("new_password", map_validator.Rules{
    Type:   reflect.String,
    Min:    map_validator.SetTotal(8),
    Unique: []string{"old_password"}, // Must be different from old_password
})
```

**Real-world Example from UpdateMyPassword:**
```go
func (h *restUser) UpdateMyPassword(c *gin.Context) {
    // ...
    data := user_controller.UpdateMemberPassword{}
    jsonDataRoles := map_validator.NewValidateBuilder().SetRules(
        map_validator.BuildRoles().
            SetRule("old_password", map_validator.Rules{Type: reflect.String}).
            SetRule("new_password", map_validator.Rules{
                Type:   reflect.String,
                Unique: []string{"old_password"}, // new_password must NOT equal old_password
                Min:    map_validator.SetTotal(8),
            }),
    )
    // ...
}

// Request JSON that will PASS validation:
{
    "old_password": "oldPass123",
    "new_password": "newPass456"  // Different from old_password ✅
}

// Request JSON that will FAIL validation:
{
    "old_password": "oldPass123",
    "new_password": "oldPass123"  // Same as old_password ❌
}
// Error: new_password must be different from old_password
```

**Important Note**: `Unique` means "must be different from", NOT "must be the same as". For confirmation fields where values should match, use business logic validation in the service layer instead.

### 4. **IPV4 Validation**
Validate IPv4 address format:
```go
.SetRule("ip_address", map_validator.Rules{
    Type:  reflect.String,
    IPV4:  true,
    Max:   map_validator.SetTotal(15), // XXX.XXX.XXX.XXX format
})

// Field can be IP or other types
.SetRule("target", map_validator.Rules{
    Type: reflect.String,
    IPV4: true,        // Can be IPv4 address
    Null: true,        // Or can be null
})
```

### 5. **List Validation**
Validate arrays of values:
```go
// Array of strings with enum validation
.SetRule("tags", map_validator.Rules{
    Type: reflect.String,
    List: map_validator.BuildListRoles()
        .SetRule("[]", map_validator.Rules{
            Enum: &map_validator.EnumField[any]{
                Items: []string{"production", "staging", "development"},
            },
        }),
})

// Array of numbers
.SetRule("scores", map_validator.Rules{
    Type: reflect.String,
    List: map_validator.BuildListRoles()
        .SetRule("[]", map_validator.Rules{
            Type: reflect.Float64,
            Min:  map_validator.SetTotal(0),
            Max:  map_validator.SetTotal(100),
        }),
})
```

### 6. **ListObject - Array of Objects**
Validate array of objects:
```go
.SetRule("items", map_validator.Rules{
    Type: reflect.String,
    ListObject: map_validator.BuildRoles()
        .SetRule("name", map_validator.Rules{
            Type: reflect.String,
            Max:  map_validator.SetTotal(100),
        })
        .SetRule("quantity", map_validator.Rules{
            Type: reflect.Int,
            Min:  map_validator.SetTotal(1),
        })
        .SetRule("price", map_validator.Rules{
            Type: reflect.Float64,
            Min:  map_validator.SetTotal(0),
        }),
})
```

### 7. **Multiple Custom Error Messages**
Different messages for different validation failures:
```go
.SetRule("username", map_validator.Rules{
    Type:        reflect.String,
    Min:         map_validator.SetTotal(3),
    Max:         map_validator.SetTotal(20),
    RegexString: `^[a-zA-Z0-9_]+$`,
    CustomMsg: map_validator.CustomMsg{
        OnMin:         map_validator.SetMessage("Username must be at least 3 characters"),
        OnMax:         map_validator.SetMessage("Username cannot exceed 20 characters"),
        OnRegexString: map_validator.SetMessage("Username can only contain letters, numbers, and underscores"),
    },
})
```

## Field Types and Rules

### 1. String Fields

#### Basic String with Max Length
```go
.SetRule("name", map_validator.Rules{
    Type: reflect.String,
    Max:  map_validator.SetTotal(255),
})
```

#### String with Regex and Custom Message
```go
.SetRule("name", map_validator.Rules{
    Type:        reflect.String,
    Max:         map_validator.SetTotal(255),
    RegexString: constant.RegexExcludeSpecialCharSpace,
    CustomMsg: map_validator.CustomMsg{
        OnRegexString: common_utils.ToPointer("the name field should not contains special character and space"),
    },
})
```

#### String with Min and Max
```go
.SetRule("password", map_validator.Rules{
    Type: reflect.String,
    Min:  map_validator.SetTotal(8),
    Max:  map_validator.SetTotal(64),
})
```

### 2. Integer Fields

#### Basic Integer
```go
.SetRule("project_id", map_validator.Rules{
    Type: reflect.Int,
})
```

#### Integer with Range
```go
.SetRule("port", map_validator.Rules{
    Type: reflect.Int,
    Min:  map_validator.SetTotal(1),
    Max:  map_validator.SetTotal(65535),
})
```

#### Integer64 for Timestamps
```go
.SetRule("start_timestamp", map_validator.Rules{
    Type: reflect.Int64,
    Max:  map_validator.SetTotal(17225601611111),
})
```

### 3. Float Fields
```go
.SetRule("limit_size", map_validator.Rules{
    Type: reflect.Float64,
    Min:  map_validator.SetTotal(0.1),
    Max:  map_validator.SetTotal(1024.0),
})
```

### 4. Boolean Fields
```go
.SetRule("public", map_validator.Rules{
    Type: reflect.Bool,
})
```

### 5. UUID Fields
```go
.SetRule("registry_id", map_validator.Rules{
    UUID: true,
})
```

### 6. Enum Fields

#### String Enum
```go
.SetRule("billing_type", map_validator.Rules{
    Type: reflect.String,
    Enum: &map_validator.EnumField[any]{
        Items: []string{"PPU", "fixed", "Monthly"},
    },
})
```

#### Integer Enum
```go
.SetRule("status", map_validator.Rules{
    Type: reflect.Int,
    Enum: &map_validator.EnumField[any]{
        Items: []int{0, 1, 2}, // 0: inactive, 1: active, 2: pending
    },
})
```

### 7. Optional Fields with Default Values
```go
.SetRule("color", map_validator.Rules{
    Type:        reflect.String,
    Max:         map_validator.SetTotal(7),
    Null:        true,
    IfNull:      "#000000",
    RegexString: `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`,
})
```

### 8. Conditional Required Fields
```go
.SetRule("limit_size", map_validator.Rules{
    Type:      reflect.Float64,
    Required:  true,
    RequiredIf: []string{"enable_storage_limit"},
})
```

### 9. Field Validation Against Other Fields
```go
.SetRule("new_password", map_validator.Rules{
    Type:   reflect.String,
    Min:    map_validator.SetTotal(8),
    Unique: []string{"old_password"}, // Must be different from old_password
})
```

### 10. IP Address Validation
```go
.SetRule("ip_address", map_validator.Rules{
    Type:  reflect.String,
    IPV4:  true,
    Max:   map_validator.SetTotal(15),
})
```

### 11. Complex Field Combinations
```go
// Field with multiple validation types
.SetRule("repo_type", map_validator.Rules{
    IPV4: true,        // Can be IP address
    List: map_validator.BuildListRoles(),  // Or list of items
    Min:  map_validator.SetTotal(1),       // Minimum 1 item
    Null: true,        // Optional field
})
```

### 12. Timestamp with Specific Limits
```go
.SetRule("occur_at", map_validator.Rules{
    Type: reflect.Int64,
    Max:  map_validator.SetTotal(17225601611111), // Specific timestamp max
})
```

## Manipulators

### Single Field Manipulator
```go
.SetManipulator("name", map_validator_utils.TrimValidation)
.SetManipulator("description", map_validator_utils.TrimValidation)
```

### Multiple Fields Manipulator
```go
.SetFieldsManipulator([]string{
    "name",
    "description",
    "purpose",
    "ip_address",
}, map_validator_utils.TrimValidation)
```

### Custom Manipulator Function
```go
// In pkg/map_validator_utils/utils.go
func TrimValidation(data interface{}) (result interface{}, err error) {
    strData, ok := data.(string)
    if !ok {
        return nil, errors.New("data is not string")
    }
    result = string_utils.CleanDashes(string_utils.TrimAndClean(strData))
    return
}
```

## Custom Messages

> **Note**: Custom messages are **optional**. If not provided, map_validator will use default error messages. Use custom messages for better user experience and clearer error descriptions.

### Basic Validation (Without Custom Messages)

```go
// ✅ VALID - Custom messages are optional
.SetRule("name", map_validator.Rules{
    Type: reflect.String,
    Max:  map_validator.SetTotal(255),
})
// Will use default error message: "name must not be greater than 255"

.SetRule("email", map_validator.Rules{
    Type:        reflect.String,
    RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
})
// Will use default error message: "email format is invalid"
```

### Using common_utils.ToPointer (Pattern from codebase)

**From registry.go - Real Example:**
```go
.SetRule("name", map_validator.Rules{
    Type:        reflect.String,
    Max:         map_validator.SetTotal(255),
    RegexString: constant.RegexExcludeSpecialCharSpace,
    CustomMsg: map_validator.CustomMsg{
        OnRegexString: common_utils.ToPointer("the name field should not contains special character and space"),
    },
})
```

**More examples from actual code:**
```go
// From tag.go
.SetRule("repo_pattern", map_validator.Rules{
    Type:        reflect.String,
    RegexString: constant.RegexExcludeSpecialCharSpace,
    CustomMsg: map_validator.CustomMsg{
        OnRegexString: common_utils.ToPointer("Repository pattern should not contains special character and space"),
    },
})

// From registryHandler.go
.SetRule("color", map_validator.Rules{
    RegexString: `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`,
    CustomMsg: map_validator.CustomMsg{
        OnRegexString: map_validator.SetMessage("Color field should be in hex format"),
    },
})
```

### Using map_validator.SetMessage
```go
CustomMsg: map_validator.CustomMsg{
    OnRegexString: map_validator.SetMessage("Color field should be in hex format"),
}
```

### Multiple Custom Messages
```go
CustomMsg: map_validator.CustomMsg{
    OnMin:         map_validator.SetMessage("Username must be at least 3 characters"),
    OnMax:         map_validator.SetMessage("Username cannot exceed 20 characters"),
    OnRegexString: map_validator.SetMessage("Username can only contain letters, numbers, and underscores"),
}
```

### Available Custom Message Types
```go
CustomMsg: map_validator.CustomMsg{
    OnType:        map_validator.SetMessage("Field type is invalid"),           // When type mismatch
    OnMin:         map_validator.SetMessage("Value is too small"),              // When value < min
    OnMax:         map_validator.SetMessage("Value exceeds maximum"),           // When value > max
    OnRegexString: map_validator.SetMessage("Format is invalid"),               // When regex fails
    OnEnum:        map_validator.SetMessage("Value not in allowed list"),       // When enum validation fails
    OnNull:        map_validator.SetMessage("Field cannot be null"),           // When null not allowed
    OnRequired:    map_validator.SetMessage("Field is required"),              // When field missing
    OnRequiredIf:  map_validator.SetMessage("Field required when X is present"), // When RequiredIf fails
    OnUnique:      map_validator.SetMessage("Field must be different from X"),   // When Unique fails
}
```

### Best Practices for Custom Messages

#### 1. **Use common_utils.ToPointer for consistency** (Recommended in this codebase)
```go
// ✅ PREFERRED - Consistent with project patterns
CustomMsg: map_validator.CustomMsg{
    OnRegexString: common_utils.ToPointer("the name field should not contains special character and space"),
}

// ✅ ALSO OK - Using SetMessage
CustomMsg: map_validator.CustomMsg{
    OnRegexString: map_validator.SetMessage("Color field should be in hex format"),
}
```

#### 2. **Write User-Friendly Messages**
```go
// ❌ BAD - Technical jargon
OnRegexString: map_validator.SetMessage("Regex validation failed")

// ✅ GOOD - User understands
OnRegexString: map_validator.SetMessage("Name can only contain letters and numbers")
```

#### 3. **Be Specific About Requirements**
```go
// ❌ VAGUE
OnMin: map_validator.SetMessage("Invalid value")

// ✅ SPECIFIC
OnMin: map_validator.SetMessage("Password must be at least 8 characters long")
OnMax: map_validator.SetMessage("Username cannot exceed 20 characters")
```

#### 4. **Consistent Message Format**
```go
// Use consistent format across your project
"OnRegexString: common_utils.ToPointer("the [field] field should not contains special character")
"OnRegexString: common_utils.ToPointer("Repository pattern should not contains special character and space")
```

### When to Use Custom Messages

#### ✅ **Use Custom Messages When:**
- **User-facing fields** - Help users understand what went wrong
- **Complex validation** - Regex patterns need explanation
- **Business requirements** - Specific rules that need clarity
- **Consistency** - Maintain consistent error format across API

```go
// Example: User registration form
.SetRule("password", map_validator.Rules{
    Type: reflect.String,
    Min:  map_validator.SetTotal(8),
    Max:  map_validator.SetTotal(64),
    RegexString: `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]`,
    CustomMsg: map_validator.CustomMsg{
        OnMin:         map_validator.SetMessage("Password must be at least 8 characters"),
        OnRegexString: map_validator.SetMessage("Password must contain uppercase, lowercase, number, and special character"),
    },
})
```

#### ✅ **Skip Custom Messages When:**
- **Internal APIs** - Technical audience understands default errors
- **Simple validation** - Obvious requirements (email format, etc.)
- **Prototype/MVP** - Quick development, can add later
- **Standard CRUD** - Basic operations don't need special messages

```go
// Example: Internal admin operation
.SetRule("internal_id", map_validator.Rules{
    Type: reflect.String,
    Max:  map_validator.SetTotal(50),
    // No CustomMsg - default message is fine for internal use
})
```

## Complex Validations

### 1. Nested Object Validation
```go
.SetRule("metadata", map_validator.Rules{
    Type: reflect.String,
    Object: map_validator.BuildRoles()
        .SetRule("labels", map_validator.Rules{
            Type: reflect.String,
            Max:  map_validator.SetTotal(500),
            Null: true,
        })
        .SetRule("annotations", map_validator.Rules{
            Type: reflect.String,
            Max:  map_validator.SetTotal(1000),
            Null: true,
        }),
})
```

### 2. Webhook Event Validation (Complex Nested)
```go
.SetRule("event_data", map_validator.Rules{
    Type: reflect.String,
    Object: map_validator.BuildRoles()
        .SetRule("push_data", map_validator.Rules{
            Object: map_validator.BuildRoles()
                .SetRule("ref", map_validator.Rules{
                    Type: reflect.String,
                    Max:  map_validator.SetTotal(255),
                })
                .SetRule("ref_type", map_validator.Rules{
                    Type: reflect.String,
                    Max:  map_validator.SetTotal(50),
                })
                .SetRule("ref_full_name", map_validator.Rules{
                    Type: reflect.String,
                    Max:  map_validator.SetTotal(512),
                }),
        })
        .SetRule("repository", map_validator.Rules{
            Object: map_validator.BuildRoles()
                .SetRule("name", map_validator.Rules{
                    Type: reflect.String,
                    Max:  map_validator.SetTotal(255),
                })
                .SetRule("full_name", map_validator.Rules{
                    Type: reflect.String,
                    Max:  map_validator.SetTotal(512),
                })
                .SetRule("date_created", map_validator.Rules{
                    Type: reflect.Int64,
                    Null: true,
                }),
        }),
})
```

### 3. List Validation
```go
.SetRule("tags", map_validator.Rules{
    Type: reflect.String,
    ListObject: map_validator.BuildRoles()
        .SetRule("name", map_validator.Rules{
            Type: reflect.String,
            Max:  map_validator.SetTotal(128),
        })
        .SetRule("digest", map_validator.Rules{
            Type: reflect.String,
            Max:  map_validator.SetTotal(128),
        }),
})
```

### 4. Array of Strings with Enum
```go
.SetRule("allowed_operations", map_validator.Rules{
    Type: reflect.String,
    List: map_validator.BuildListRoles()
        .SetRule("[]", map_validator.Rules{
            Enum: &map_validator.EnumField[any]{
                Items: []string{"read", "write", "delete", "admin"},
            },
        }),
})
```

## Settings

### Strict Mode
```go
.SetSetting(*map_validator.BuildSetting().MakeStrict())
// atau
.SetSetting(map_validator.Setting{Strict: true})
```

### Complete Example with Settings
```go
roles := map_validator.BuildRoles()
    .SetRule("name", map_validator.Rules{Type: reflect.String})
    .SetRule("public", map_validator.Rules{Type: reflect.Bool})
    .SetSetting(*map_validator.BuildSetting().MakeStrict())

jsonDataRoles := map_validator.NewValidateBuilder().SetRules(roles)
```

## Error Handling

### Basic Error Handling
```go
jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)
if err != nil {
    c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{
        Message: err.Error(),
    })
    return
}

jsonData, err := jsonDataValidate.RunValidate()
if err != nil {
    c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{
        Message: err.Error(),
    })
    return
}
```

### Advanced Error Handling with Context
```go
jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)
if err != nil {
    // Log error for debugging
    log.Errorf("Validation error: %v", err)

    // Return user-friendly message
    c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{
        Message: "Invalid request format",
    })
    return
}

jsonData, err := jsonDataValidate.RunValidate()
if err != nil {
    // Log validation failure
    log.Warnf("Field validation failed: %v", err)

    // Return specific error
    c.JSON(http.StatusUnprocessableEntity, gin_utils.MessageResponse{
        Message: err.Error(),
    })
    return
}
```

### Error Handling with Rate Limiting (Login Example)
```go
validate, err := jsonDataValidate.RunValidate()
if err != nil {
    // Handle failed login attempt tracking
    key := fmt.Sprintf(constant.CacheKeyLoginFailedAttempt, email)
    attempt, _ := h.cacheService.Get(key)

    var num int
    if attempt != "" {
        num, _ = strconv.Atoi(attempt)
    }

    // Increment counter
    str := "1"
    if attempt != "" {
        str = strconv.Itoa(num + 1)
        h.cacheService.Del(key)
        h.cacheService.Set(key, str, time.Duration(failedWaitTimeMinutes)*time.Minute)
    } else {
        h.cacheService.Set(key, "1", time.Duration(failedWaitTimeMinutes)*time.Minute)
    }

    // Determine response code based on attempts
    statusCode := http.StatusBadRequest
    if num >= 2 {
        statusCode = http.StatusTooManyRequests
    }

    c.JSON(statusCode, gin_utils.MessageResponse{
        Message: fmt.Sprintf("%s: Show password attempt : %s of 3", err.Error(), str),
    })
    return
}
```

## Complete Real-World Examples

### Example 1: Create Registry Request
```go
// Build validation rules
roles := map_validator.BuildRoles()
    .SetRule("name", map_validator.Rules{
        Type:        reflect.String,
        Max:         map_validator.SetTotal(255),
        RegexString: constant.RegexExcludeSpecialCharSpace,
        CustomMsg: map_validator.CustomMsg{
            OnRegexString: common_utils.ToPointer("the name field should not contains special character and space"),
        },
    })
    .SetRule("description", map_validator.Rules{
        Type: reflect.String,
        Max:  map_validator.SetTotal(1000),
        Null: true,
    })
    .SetRule("purpose", map_validator.Rules{
        Type: reflect.String,
        Max:  map_validator.SetTotal(255),
    })
    .SetRule("public", map_validator.Rules{
        Type: reflect.Bool,
    })
    .SetRule("agreeSA", map_validator.Rules{
        Type: reflect.Bool,
    })
    .SetRule("agreeSoW", map_validator.Rules{
        Type: reflect.Bool,
    })
    .SetRule("ip_address", map_validator.Rules{
        Type:  reflect.String,
        IPV4:  true,
        Max:   map_validator.SetTotal(15),
    })

// Apply manipulators
roles.
    SetManipulator("description", map_validator_utils.TrimValidation).
    SetManipulator("name", map_validator_utils.TrimValidation).
    SetManipulator("purpose", map_validator_utils.TrimValidation).
    SetManipulator("ip_address", map_validator_utils.TrimValidation)

// Create and run validator
jsonDataRoles := map_validator.NewValidateBuilder().SetRules(roles)
jsonDataValidate, err := jsonDataRoles.LoadJsonHttp(c.Request)
if err != nil {
    c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
    return
}

jsonData, err := jsonDataValidate.RunValidate()
if err != nil {
    c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
    return
}

// Bind to struct
var request registry_controller.CreateRegistry
jsonData.Bind(&request)
```

### Example 2: Update User with Conditional Fields
```go
roles := map_validator.BuildRoles()
    .SetRule("id", map_validator.Rules{
        Type: reflect.Int,
    })
    .SetRule("full_name", map_validator.Rules{
        Type:        reflect.String,
        Max:         map_validator.SetTotal(100),
        Null:        true,
        RegexString: constant.RegexExcludeSpecialChar,
    })
    .SetRule("email", map_validator.Rules{
        Type:        reflect.String,
        Max:         map_validator.SetTotal(255),
        RegexString: constant.RegexEmail,
        Null:        true,
    })
    .SetRule("role", map_validator.Rules{
        Type: reflect.String,
        Null: true,
        Enum: &map_validator.EnumField[any]{
            Items: []string{"ProjectAdmin", "Developer", "Guest", "Maintainer"},
        },
    })
    .SetRule("password", map_validator.Rules{
        Type: reflect.String,
        Min:  map_validator.SetTotal(8),
        Max:  map_validator.SetTotal(64),
        Null: true,
        RequiredIf: []string{"confirm_password"},
        Unique: []string{"old_password"},
    })

// Apply manipulator for string fields
roles.SetFieldsManipulator([]string{
    "full_name",
    "email",
}, map_validator_utils.TrimValidation)
```

### Example 3: Webhook Event Configuration
```go
roles := map_validator.BuildRoles()
    .SetRule("name", map_validator.Rules{
        Type:        reflect.String,
        Max:         map_validator.SetTotal(255),
        RegexString: constant.RegexExcludeSpecialCharSpace,
        CustomMsg: map_validator.CustomMsg{
            OnRegexString: common_utils.ToPointer("Webhook name should not contains special character"),
        },
    })
    .SetRule("description", map_validator.Rules{
        Type: reflect.String,
        Max:  map_validator.SetTotal(1000),
        Null: true,
    })
    .SetRule("url", map_validator.Rules{
        Type:        reflect.String,
        Max:         map_validator.SetTotal(2048),
        RegexString: constant.RegexURL,
        CustomMsg: map_validator.CustomMsg{
            OnRegexString: common_utils.ToPointer("Invalid URL format"),
        },
    })
    .SetRule("secret_token", map_validator.Rules{
        Type: reflect.String,
        Min:  map_validator.SetTotal(16),
        Max:  map_validator.SetTotal(256),
        Null: true,
    })
    .SetRule("project_id", map_validator.Rules{
        UUID: true,
    })
    .SetRule("enabled_events", map_validator.Rules{
        Type: reflect.String,
        List: map_validator.BuildListRoles()
            .SetRule("[]", map_validator.Rules{
                Enum: &map_validator.EnumField[any]{
                    Items: []string{
                        "push",
                        "pull_request",
                        "release",
                        "tag",
                        "delete",
                    },
                },
            }),
    })
    .SetRule("config", map_validator.Rules{
        Type: reflect.String,
        Object: map_validator.BuildRoles()
            .SetRule("retry_count", map_validator.Rules{
                Type: reflect.Int,
                Min:  map_validator.SetTotal(0),
                Max:  map_validator.SetTotal(5),
                Null: true,
                IfNull: 3,
            })
            .SetRule("timeout", map_validator.Rules{
                Type: reflect.Int,
                Min:  map_validator.SetTotal(5),
                Max:  map_validator.SetTotal(300),
                Null: true,
                IfNull: 30,
            })
            .SetRule("content_type", map_validator.Rules{
                Type: reflect.String,
                Enum: &map_validator.EnumField[any]{
                    Items: []string{"json", "form"},
                },
                Null: true,
                IfNull: "json",
            }),
    })
    .SetSetting(*map_validator.BuildSetting().MakeStrict())

// Apply manipulators
roles.
    SetManipulator("name", map_validator_utils.TrimValidation).
    SetManipulator("description", map_validator_utils.TrimValidation).
    SetManipulator("url", map_validator_utils.TrimValidation)
```

## Best Practices

1. **Always use TrimValidation** for string inputs to clean data
2. **Provide clear error messages** for regex validations
3. **Use enums** for fields with limited possible values
4. **Set appropriate max lengths** for all string fields
5. **Use UUID validation** for ID fields
6. **Handle errors consistently** across all endpoints
7. **Use strict mode** when you want to reject unknown fields
8. **Validate nested objects** for complex data structures
9. **Use conditional requirements** for optional but dependent fields
10. **Log validation errors** for debugging purposes

## ⚠️ PERFORMANCE & SECURITY WARNINGS

### Performance Considerations

#### ❌ AVOID - Heavy Operations in Validation
```go
// Don't do expensive operations
.SetRule("email", map_validator.Rules{
    Type: reflect.String,
    // ❌ WRONG - No database checks, no external API calls
})
```

#### ⚠️ COMPLEX NESTED VALIDATIONS
```go
// This pattern exists but has performance impact
.SetRule("event_data", map_validator.Rules{
    Object: map_validator.BuildRoles()
        .SetRule("resources", map_validator.Rules{
            ListObject: map_validator.BuildRoles()
                .SetRule("scan_overview", map_validator.Rules{
                    Object: map_validator.BuildRoles()
                        .SetRule("summary", map_validator.Rules{
                            Object: map_validator.BuildRoles() // 5+ levels deep!
```

#### ✅ PERFORMANCE TIPS
1. Keep validation shallow (max 3 levels)
2. Use regex for simple pattern matching
3. Avoid expensive string operations on long texts
4. Use appropriate Max limits to prevent DoS

### Security Best Practices

#### ✅ INPUT SANITIZATION
```go
// Always sanitize strings
.SetFieldsManipulator([]string{"name", "description"}, map_validator_utils.TrimValidation)

// Use regex to prevent injection
.SetRule("field", map_validator.Rules{
    RegexString: constant.RegexExcludeSpecialChar,
})
```

#### ✅ ENUM VALIDATION FOR SECURITY
```go
// Prevent unauthorized values
.SetRule("role", map_validator.Rules{
    Enum: &map_validator.EnumField[any]{
        Items: []string{"ProjectAdmin", "Developer", "Guest"},
    },
})
```

#### ✅ UUID VALIDATION
```go
// Prevent ID injection
.SetRule("project_id", map_validator.Rules{UUID: true})
```

#### ❌ SECURITY ANTI-PATTERNS
```go
// Don't validate passwords beyond format
.SetRule("password", map_validator.Rules{
    Min: map_validator.SetTotal(8), // OK
    // ❌ DON'T check password strength here - do it in service layer
})
```

### Additional Rules for AI/LLM

### 7. **PERFORMANCE CONSIDERATIONS**
```go
// ✅ GOOD - Validation is lightweight and fast
.SetRule("name", map_validator.Rules{
    Type: reflect.String,
    Max:  map_validator.SetTotal(255),
})

// ❌ AVOID - Don't do expensive operations in validation
.SetRule("email", map_validator.Rules{
    Type: reflect.String,
    // Don't check database here!
})
```

### 8. **ERROR RESPONSE CONSISTENCY**
```go
// ✅ ALWAYS use this pattern for error responses
if err != nil {
    c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{
        Message: err.Error(), // Direct error message
    })
    return
}

// ❌ DON'T wrap or modify error messages
if err != nil {
    c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{
        Message: "Validation failed: " + err.Error(), // Don't add prefixes
    })
    return
}
```

### 9. **VALIDATION RULES LOCALITY**
```go
// ✅ CORRECT - All validation rules in one place
func (h *restHandler) CreateUser(c *gin.Context) {
    roles := map_validator.BuildRoles()
        .SetRule("email", emailRules)
        .SetRule("password", passwordRules)
        .SetRule("name", nameRules)

    // Don't split validation logic across multiple functions
}
```

### 10. **TYPE SAFETY**
```go
// ✅ Use specific types
.SetRule("age", map_validator.Rules{
    Type: reflect.Int,  // Not reflect.String
})

// ✅ Convert after validation
age := int(jsonData.Get("age").(float64))
```

## 🚨 ANTI-PATTERNS TO AVOID

### 1. **Architecture Violations**
```go
// ❌ WRONG - Don't use in service layer
func (s *service) CreateLabel(data Data) error {
    // Validation should be in controller layer
}

// ❌ WRONG - Don't use in repository layer
func (r *repository) Save(data Data) error {
    // Validation should be in controller layer
}

// ✅ CORRECT - In any controller layer
func (h *anyHandler) CreateLabel(c *gin.Context) {
    // Validation here in controller layer
}

// ✅ CORRECT - In REST view layer (if exists)
func (h *restRegistry) CreateLabel(c *gin.Context) {
    // Validation here in controller layer
}
```

### 2. **Inconsistent Error Handling**
```go
// ❌ WRONG - Different response formats
c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": err.Error()})
c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})

// ✅ CORRECT - Consistent format
c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
```

### 3. **Magic Numbers Without Context**
```go
// ❌ WRONG - Unclear magic numbers
.SetRule("timestamp", map_validator.Rules{
    Max: map_validator.SetTotal(17225601611111), // What is this?
})

// ✅ BETTER - Add comments for clarity
.SetRule("timestamp", map_validator.Rules{
    Max: map_validator.SetTotal(17225601611111), // Max timestamp for year 2024
})
```

### 4. **Overly Complex Validations**
```go
// ❌ AVOID - 5+ levels of nesting
.SetRule("event_data", map_validator.Rules{
    Object: map_validator.BuildRoles()
        .SetRule("level1", map_validator.Rules{
            Object: map_validator.BuildRoles()
                .SetRule("level2", map_validator.Rules{
                    Object: map_validator.BuildRoles()
                        .SetRule("level3", map_validator.Rules{
                            Object: map_validator.BuildRoles()... // Too deep!
```

### 5. **Inconsistent TrimValidation Usage**
```go
// ❌ INCONSISTENT - Some fields trimmed, others not
.SetManipulator("name", map_validator_utils.TrimValidation)
// Missing TrimValidation for "description"

// ✅ CONSISTENT - Apply to all string fields
.SetFieldsManipulator([]string{
    "name", "description", "purpose", "notes",
}, map_validator_utils.TrimValidation)
```

### 6. **Business Logic in Validation**
```go
// ❌ WRONG - Checking business rules
.SetRule("email", map_validator.Rules{
    // Don't check if email exists in database here!
})

// ✅ CORRECT - Only format validation
.SetRule("email", map_validator.Rules{
    RegexString: constant.RegexEmail,
})
// Business check in service layer
if exists {
    return errors.New("email already exists")
}
```

## SUMMARY OF RULES

### Architectural Rules
1. **Scope**: Only in controller layer for HTTP requests
2. **Purpose**: Input sanitization and format checking
3. **No sharing**: Define validation inline, no shared functions
4. **No business logic**: Keep it simple and fast
5. **Architecture**: Clear separation between validation and business logic

### Implementation Rules
6. **Consistent errors**: Direct error messages, no wrapping
7. **Performance**: Lightweight operations only
8. **Security**: Sanitize inputs, validate enums, use UUID for IDs
9. **Consistency**: Same patterns across all endpoints
10. **Simplicity**: Avoid overly complex nested validations (max 3 levels recommended)

### Advanced Features Usage
11. **RequiredIf**: For conditional field requirements
12. **IfNull**: For default values on optional fields
13. **Unique**: For field-to-field comparisons
14. **IPV4**: For IP address validation
15. **List/ListObject**: For array validation
16. **Multiple CustomMsg**: For different error types

### Data Access Patterns
17. **Bind Method**: Use for complex requests (recommended)
18. **Get Method**: Use for simple cases or direct field access
19. **Type Assertions**: Always assert to correct types
20. **Optional Fields**: Check for nil before accessing
21. **Bind Error Handling**: Always handle Bind method errors properly
22. **Log Bind Errors**: Log bind errors for debugging but return generic message to client

## Common Constants Used

```go
// Example regex constants (check constant/ package for actual values)
constant.RegexExcludeSpecialChar   // Excludes special characters
constant.RegexExcludeSpecialCharSpace  // Excludes special chars and spaces
constant.RegexEmail               // Email validation
constant.RegexURL                 // URL validation
```