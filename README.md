# go-map-validator

Simple, declarative validation for Go maps and HTTP payloads with nested and list object support.

- Validate `map[string]interface{}` and `http.Request` (JSON or multipart).
- Compose rules fluently. Support nested object and list of object.
- Transform values with manipulators and plug custom extensions.
- Bind validated data back to your struct.

Examples: see the test suite: https://github.com/Rhyanz46/go-map-validator/tree/main/test

## Install

```bash
go get github.com/Rhyanz46/go-map-validator/map_validator
```

Import: `import "github.com/Rhyanz46/go-map-validator/map_validator"`

## Quick Start

```go
import "github.com/Rhyanz46/go-map-validator/map_validator"

type Login struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// Build rules once — safe to reuse across handlers and goroutines.
rules := map_validator.BuildRoles().
    SetRule("email", map_validator.Email().WithMax(100)).
    SetRule("password", map_validator.Str().Between(6, 30)).
    Done()

// HTTP handler: one-liner validate-and-bind
dto, err := map_validator.ValidateJSON[Login](httpReq, rules)
if err != nil { /* handle — ErrNoRules, ErrInvalidJsonFormat, or validation error */ }

// Or validate a plain map (no HTTP):
payload := map[string]interface{}{"email": "dev@example.com", "password": "secret123"}
extra, err := map_validator.NewValidateBuilder().SetRules(rules).Load(payload)
if err != nil { /* handle */ }
result, err := extra.RunValidate()
if err != nil { /* handle */ }
var dto2 Login
_ = result.Bind(&dto2)
```

## Features

- Validate map by keys and HTTP JSON/multipart (file upload supported).
- Types via `reflect.Kind` with nullable fields (`Null`) and default (`IfNull`).
- `Enum`, `RegexString`, `Email`, `UUID`/`UUIDToString`.
- IPv4 validators: `IPV4`, `IPV4Network` (.0 network), `IPv4OptionalPrefix` (CIDR optional).
- `Min`/`Max` for strings, numeric, and slices.
- Nested object (`Object`) and list of object (`ListObject`).
- Uniqueness across sibling fields (`Unique`).
- Conditional required: `RequiredWithout` and `RequiredIf`.
- Strict mode to reject unknown keys (`Setting{Strict:true}`).
- Custom messages for type/regex/min/max/unique/enum value.
- Manipulators to post-process values.
- Extensions lifecycle hooks.
- **Short constructors** (`Str()`, `Int()`, `Email()`, `UUID()`, `IPv4()`, `StrEnum()`, …) and chain helpers (`.WithMax()`, `.Nullable()`, `.Between()`, `.Regex()`, …) for concise rule definitions.
- **`ValidateJSON[T]`** — generic one-liner that collapses the Load → Validate → Bind pipeline for HTTP handlers.
- **Safe for shared & concurrent use** — rules no longer hold per-call state, so a single `rules` value can be declared as a package-level var and reused across handlers and goroutines.

## What's new in v0.0.41

All changes are additive on the public API — existing usage patterns keep working.

**New APIs**
- `ValidateJSON[T any](r *http.Request, rules RulesWrapper) (T, error)` — collapses Load → Validate → Bind into a single call. See [One-liner for JSON handlers](#one-liner-for-json-handlers).
- Short constructors: `Str()`, `Int()`, `Int64()`, `Float64()`, `Bool()`, `Email()`, `UUID()`, `IPv4()`, `StrEnum(items…)`, `IntEnum(items…)`, `NestedObject(w)`, `ListOfObject(w)`.
- Chain helpers on `Rules`: `.Nullable()`, `.Default(v)`, `.WithMin(n)`, `.WithMax(n)`, `.Between(min, max)`, `.Regex(p)`, `.WithMsg(cm)`, `.UniqueFrom(…)`, `.WithRequiredIf(…)`, `.WithRequiredWithout(…)`. Value-receiver, so each call returns a copy — safe to chain.
- `Done()` on `RulesWrapper` — optional chain terminator so `BuildRoles().SetRule(…).Done()` reads cleanly.
- `OnEnumValueNotMatch` field in `CustomMsg` for custom enum-mismatch messages, plus two new template variables: `${actual_value}` and `${enum_values}`.
- New error sentinel `ErrNoRules` returned by `Load` / `LoadJsonHttp` / `LoadFormHttp` when no rules are set.

**Behavior changes**
- `SetRules(empty)` no longer panics. The error surfaces as `ErrNoRules` from the subsequent `Load*` call — easier to handle uniformly.
- Enum with an unsupported element kind no longer panics. It returns a normal validation error instead.
- Rules are now safe to share across handlers and goroutines. The previous mutable-state bug inside `rulesWrapper` is gone — you can declare `rules` as a package-level variable and reuse it freely.

**Bug fixes**
- `ListObject` items now bind into the target struct with their actual field values. Before, items often bound as zero values because of a key-shadowing bug inside `validateRecursive`.

## HTTP Integration (JSON)

```go
op, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(req)
if err != nil { /* handle */ }
extra, err := op.RunValidate()
if err != nil { /* handle */ }
// Bind to your DTO
var dto MyDTO
_ = extra.Bind(&dto)
```

Note: JSON numbers decode as `float64`. The validator tolerates integer-family comparisons when rules expect an int kind.

## One-liner for JSON handlers

`ValidateJSON[T]` collapses the Load → Validate → Bind pipeline into a single generic call. Ideal for HTTP handlers where you rarely need the intermediate state.

```go
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

rules := map_validator.BuildRoles().
    SetRule("email", map_validator.Email().WithMax(255)).
    SetRule("password", map_validator.Str().Between(6, 64)).
    Done()

req, err := map_validator.ValidateJSON[LoginRequest](httpReq, rules)
if err != nil {
    // forward as-is: ErrNoRules, ErrInvalidJsonFormat, or validation error
    return err
}
// req is already validated and bound
```

Before vs after — same validation, less ceremony:

```go
// before: 4 error checks, intermediate vars
op, err := map_validator.NewValidateBuilder().SetRules(rules).LoadJsonHttp(r)
if err != nil { return err }
extra, err := op.RunValidate()
if err != nil { return err }
var dto LoginRequest
if err := extra.Bind(&dto); err != nil { return err }

// after: one call, one error check
dto, err := map_validator.ValidateJSON[LoginRequest](r, rules)
if err != nil { return err }
```

The 5-step pipeline remains the right choice when you need access to the `ExtraOperationData` returned by `RunValidate()` before bind — e.g. reading `GetFilledField()` / `GetNullField()` / `GetData()`, or running custom logic against the validated map before committing it to your struct. Extension lifecycle hooks (`BeforeLoad`, `AfterLoad`, `BeforeValidation`, `AfterValidation`) still run under `ValidateJSON` because it internally composes the same `LoadJsonHttp` + `RunValidate` calls.

Since rules no longer hold per-call mutable state, the same `rules` value can be declared as a package-level variable and shared across handlers safely — including concurrent requests.

## Nested Objects

```go
filter := map_validator.BuildRoles().
    SetRule("search", map_validator.Rules{Type: reflect.String, Null: true}).
    SetRule("organization_id", map_validator.Rules{UUID: true, Null: true}).
    Done()

parent := map_validator.BuildRoles().
    SetRule("filter", map_validator.Rules{Object: filter, Null: true}).
    SetRule("rows_per_page", map_validator.Rules{Type: reflect.Int64, Null: true}).
    SetRule("page_index", map_validator.Rules{Type: reflect.Int64, Null: true}).
    SetRule("sort", map_validator.Rules{
        Null:   true,
        IfNull: "FULL_NAME:DESC",
        Type:   reflect.String,
        Enum:   &map_validator.EnumField[any]{Items: []string{"FULL_NAME:DESC","FULL_NAME:ASC","EMAIL:ASC","EMAIL:DESC"}},
    }).
    SetSetting(map_validator.Setting{Strict: true}).
    Done()

extra, err := map_validator.NewValidateBuilder().SetRules(parent).Load(map[string]interface{}{}).RunValidate()
_ = extra; _ = err
```

## List of Objects

```go
item := map_validator.BuildRoles().
    SetRule("name", map_validator.Rules{Type: reflect.String}).
    SetRule("quantity", map_validator.Rules{Type: reflect.Int, Min: map_validator.SetTotal(1)}).
    Done()

rules := map_validator.BuildRoles().
    SetRule("goods", map_validator.Rules{ListObject: item}).
    Done()

payload := map[string]interface{}{
    "goods": []interface{}{
        map[string]interface{}{"name": "Apple", "quantity": 2},
    },
}
_, err := map_validator.NewValidateBuilder().SetRules(rules).Load(payload).RunValidate()
if err != nil { panic(err) }
```

When using `ListObject`, the input must be an array of objects. Sending a single object yields: `"field 'goods' is not valid list object"`.

## Unique and Conditional Required

```go
rules := map_validator.BuildRoles().
    SetRule("password", map_validator.Rules{Type: reflect.String, Null: true}).
    SetRule("new_password", map_validator.Rules{Type: reflect.String, Unique: []string{"password"}, Null: true}).
    SetRule("flavor", map_validator.Rules{Type: reflect.String, RequiredWithout: []string{"custom_flavor"}}).
    SetRule("custom_flavor", map_validator.Rules{Type: reflect.String, RequiredIf: []string{"flavor"}}).
    Done()
```

## Custom Messages

Supported fields in `CustomMsg`:
- `OnTypeNotMatch`, `OnRegexString`, `OnMin`, `OnMax`, `OnUnique`, `OnEnumValueNotMatch`.

Message variables:

- `${field}`: nama field yang divalidasi (key pada rules).
- `${expected_type}`: tipe yang diharapkan (hasil `reflect.Kind.String()` dari rules).
- `${actual_type}`: tipe aktual dari nilai yang diterima.
- `${actual_length}`: panjang aktual (string: jumlah rune; angka: nilai numerik yang dibandingkan; slice: jumlah elemen).
- `${expected_min_length}`: nilai/ukuran minimum yang diharapkan (`Min`).
- `${expected_max_length}`: nilai/ukuran maksimum yang diharapkan (`Max`).
- `${unique_origin}`: nama field asal pada pengecekan unik.
- `${unique_target}`: nama field target yang dibandingkan pada pengecekan unik.
- `${actual_value}`: nilai aktual yang dikirim (tersedia di `OnEnumValueNotMatch`).
- `${enum_values}`: daftar nilai enum yang diperbolehkan (tersedia di `OnEnumValueNotMatch`).

```go
rules := map_validator.BuildRoles().
  SetRule("total", map_validator.Rules{
    Type: reflect.Int,
    Min:  map_validator.SetTotal(2),
    Max:  map_validator.SetTotal(3),
    CustomMsg: map_validator.CustomMsg{
      OnMin: map_validator.SetMessage("Min is ${expected_min_length}, got ${actual_length}"),
      OnMax: map_validator.SetMessage("Max is ${expected_max_length}, got ${actual_length}"),
    },
  }).
  Done()
```

Examples:

- Type mismatch
```go
rules := map_validator.BuildRoles().
  SetRule("qty", map_validator.Rules{
    Type: reflect.Int64,
    CustomMsg: map_validator.CustomMsg{
      OnTypeNotMatch: map_validator.SetMessage("Field ${field} must be ${expected_type}, but got ${actual_type}"),
    },
  }).
  Done()
```

- Regex validation
```go
rules := map_validator.BuildRoles().
  SetRule("email", map_validator.Rules{
    RegexString: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
    CustomMsg:   map_validator.CustomMsg{OnRegexString: map_validator.SetMessage("Your ${field} is not a valid email")},
  }).
  Done()
```

- Unique across fields
```go
rules := map_validator.BuildRoles().
  SetRule("password", map_validator.Rules{Type: reflect.String}).
  SetRule("new_password", map_validator.Rules{
    Type: reflect.String, Unique: []string{"password"},
    CustomMsg: map_validator.CustomMsg{OnUnique: map_validator.SetMessage("The value of '${unique_origin}' must be different from '${unique_target}'")},
  }).
  Done()
```

- Enum value not in list
```go
rules := map_validator.BuildRoles().
  SetRule("status", map_validator.StrEnum("active", "inactive", "pending").
    WithMsg(map_validator.CustomMsg{
      OnEnumValueNotMatch: map_validator.SetMessage("'${field}' value '${actual_value}' is not allowed; expected one of ${enum_values}"),
    })).
  Done()
// → "'status' value 'banned' is not allowed; expected one of [active inactive pending]"
```

Notes:
- If a corresponding `CustomMsg` is not set, the default error message is used.
- Variables are replaced contextually at error time (e.g., `${field}` is the rule’s key).
- Currently not customizable: null errors, and specific `RequiredWithout` / `RequiredIf` messages.

## Manipulators (Post-process)

```go
rules := map_validator.BuildRoles().
  SetRule("name", map_validator.Rules{Type: reflect.String}).
  SetManipulator("name", func(v interface{}) (interface{}, error) {
    s := strings.TrimSpace(v.(string))
    return strings.ToUpper(s), nil
  }).
  Done()
```

Manipulators run after validation and before `Bind()` on the built value tree (including nested/list fields with matching keys).

## Extensions

Implement `ExtensionType` to hook into load/validate lifecycle. Example scaffold: `example_extensions/example.go`.

Use-cases:
- Normalize/transform input across many fields.
- Convert string→number for multipart forms.
- Enrich data or apply cross-cutting rules.

## Strict Mode

Set `Setting{Strict:true}` in a rules group to reject any unknown keys at that object level. Apply again for nested rules where needed.

## Notes & Caveats

- JSON numbers decode as `float64`. Integer-family kinds are tolerated on JSON input.
- `LoadFormHttp`: non-file values arrive as strings; there is no automatic string→int/float/bool parsing. Use a manipulator or extension to convert.
- Email validation is simple (checks `@` and `.`), not full RFC compliance.
- Error reporting returns the first encountered error (no multi-error aggregation with field paths yet).
- Custom messages are not yet available for: null errors and specific `RequiredWithout` / `RequiredIf` messages.
- Empty rules no longer panic. `SetRules` accepts them silently; the subsequent `Load` / `LoadJsonHttp` / `LoadFormHttp` returns `ErrNoRules` so callers can handle it uniformly.

## Roadmap

- Detailed error reporting with field paths and multi-error aggregation.
- URL params extraction helpers.
- Base64 validation.
- Multipart file size limits and image resolution checks.
- OpenAPI spec generator extension.
- Multi-validator per field (e.g., IPv4 + UUID combined).

## Community

- Updates/Discussion: [Telegram](https://t.me/addlist/Wi84VFNkvz85MWFl)
