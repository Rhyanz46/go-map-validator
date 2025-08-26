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
// Validate a map and bind to struct
payload := map[string]interface{}{"email": "dev@example.com", "password": "secret123"}

rules := map_validator.BuildRoles().
    SetRule("email", map_validator.Rules{Type: reflect.String, Email: true, Max: map_validator.SetTotal(100)}).
    SetRule("password", map_validator.Rules{Type: reflect.String, Min: map_validator.SetTotal(6), Max: map_validator.SetTotal(30)}).
    Done()

op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(payload)
if err != nil { panic(err) }
extra, err := op.RunValidate()
if err != nil { panic(err) }

type Login struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
var dto Login
if err := extra.Bind(&dto); err != nil { panic(err) }
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
- Custom messages for type/regex/min/max/unique.
- Manipulators to post-process values.
- Extensions lifecycle hooks.

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
- `OnTypeNotMatch`, `OnRegexString`, `OnMin`, `OnMax`, `OnUnique`.

Message variables:

- `${field}`: nama field yang divalidasi (key pada rules).
- `${expected_type}`: tipe yang diharapkan (hasil `reflect.Kind.String()` dari rules).
- `${actual_type}`: tipe aktual dari nilai yang diterima.
- `${actual_length}`: panjang aktual (string: jumlah rune; angka: nilai numerik yang dibandingkan; slice: jumlah elemen).
- `${expected_min_length}`: nilai/ukuran minimum yang diharapkan (`Min`).
- `${expected_max_length}`: nilai/ukuran maksimum yang diharapkan (`Max`).
- `${unique_origin}`: nama field asal pada pengecekan unik.
- `${unique_target}`: nama field target yang dibandingkan pada pengecekan unik.

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

Notes:
- If a corresponding `CustomMsg` is not set, the default error message is used.
- Variables are replaced contextually at error time (e.g., `${field}` is the rule’s key).
- Currently not customizable: enum mismatch, null errors, and specific RequiredWithout/If messages.

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
- Custom messages are not yet available for: enum mismatch, null errors, and RequiredWithout/If specific messages.

## Roadmap

- Detailed error reporting with field paths and multi-error aggregation.
- URL params extraction helpers.
- Base64 validation.
- Multipart file size limits and image resolution checks.
- OpenAPI spec generator extension.
- Multi-validator per field (e.g., IPv4 + UUID combined).

## Community

- Updates/Discussion: [Telegram](https://t.me/addlist/Wi84VFNkvz85MWFl)
