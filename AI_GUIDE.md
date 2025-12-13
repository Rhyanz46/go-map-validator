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

## ⚠️ IMPORTANT RULES FOR AI/LLM

### 1. **USAGE SCOPE RESTRICTIONS**
- **❌ NEVER** use `map_validator` in:
  - Service layer (usecase)
  - Repository layer
  - Model/Entity structs
  - Utility functions
  - Shared validation logic
  - Database operations
- **✅ ONLY** use `map_validator` in:
  - Controller layer (any folder containing HTTP handlers)
  - HTTP request validation

### 1.1 **FINDING THE CONTROLLER LAYER**
```
Common controller layer locations:
- /app/v2/views/rest/          # REST API handlers
- /app/v1/views/rest/          # REST API handlers
- /app/controllers/           # Traditional controllers
- /app/controller/           # Traditional controllers
- /app/handlers/             # HTTP request handlers
- /internal/http/            # Internal HTTP handlers
- /api/handlers/             # API handlers
```

**IMPORTANT**: If you cannot identify the controller layer in the project, **ASK THE USER** before proceeding with map_validator implementation. Example questions:
- "Where are the HTTP request handlers located in this project?"
- "Which folder contains the controller/REST handlers?"
- "What is the project structure for HTTP handlers?"

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
1. [Basic Validation Pattern](#basic-validation-pattern)
2. [Field Types and Rules](#field-types-and-rules)
3. [Manipulators](#manipulators)
4. [Custom Messages](#custom-messages)
5. [Complex Validations](#complex-validations)
6. [Settings](#settings)
7. [Error Handling](#error-handling)

## Basic Validation Pattern

Pattern standar yang digunakan di seluruh project:

```go
// 1. Build validation rules
roles := map_validator.BuildRoles()
    .SetRule("field_name", map_validator.Rules{...})
    .SetRule("another_field", map_validator.Rules{...})

// 2. Create validator
jsonDataRoles := map_validator.NewValidateBuilder().SetRules(roles)

// 3. Load from HTTP request
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

### Using common_utils.ToPointer
```go
CustomMsg: map_validator.CustomMsg{
    OnRegexString: common_utils.ToPointer("the name field should not contains special character and space"),
}
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
    OnMin:   map_validator.SetMessage("Value is too small"),
    OnMax:   map_validator.SetMessage("Value exceeds maximum allowed"),
    OnRegex: map_validator.SetMessage("Invalid format"),
}
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

1. **Scope**: Only in controller layer for HTTP requests
2. **Purpose**: Input sanitization and format checking
3. **No sharing**: Define validation inline, no shared functions
4. **No business logic**: Keep it simple and fast
5. **Consistent errors**: Direct error messages, no wrapping
6. **Performance**: Lightweight operations only
7. **Architecture**: Clear separation between validation and business logic
8. **Security**: Sanitize inputs, validate enums, use UUID for IDs
9. **Consistency**: Same patterns across all endpoints
10. **Simplicity**: Avoid overly complex nested validations

## Common Constants Used

```go
// Example regex constants (check constant/ package for actual values)
constant.RegexExcludeSpecialChar   // Excludes special characters
constant.RegexExcludeSpecialCharSpace  // Excludes special chars and spaces
constant.RegexEmail               // Email validation
constant.RegexURL                 // URL validation
```