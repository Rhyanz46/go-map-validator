# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to incremental patch versioning (`v0.0.x`).

## [v0.0.43]

All changes are additive on the public API — existing usage patterns keep
working without modification. This release responds to feedback about silent
data-loss when struct fields lack a corresponding rule, by adding two opt-in
escape hatches.

### Added

- **`List(elem Rules) Rules`** — primitive list shortcut. Element constraints (`Min`, `Max`, `Regex`, `Enum`, `UUID`, `Email`, `IPv4`) inherit from the inner rule. Container size constraints chain after `List(...)` via `.WithMin` / `.WithMax` / `.Between`.
  ```go
  SetRule("tags",   List(Str().WithMax(64)))           // each tag ≤ 64 chars
  SetRule("ids",    List(UUID()).WithMin(1))           // ≥ 1 valid UUID
  SetRule("emails", List(Email()))                     // each item is valid email
  SetRule("colors", List(StrEnum("red","blue","green")))
  SetRule("scores", List(Int().Between(0,100)).WithMax(10))
  ```

- **`Any() Rules`** — passthrough escape hatch. The field must be present (use `.Nullable()` to make optional) but its value is not validated and is preserved verbatim through `Bind()`. For heterogeneous metadata, raw config, third-party payloads.
  ```go
  SetRule("metadata", Any())              // required, any value
  SetRule("settings", Any().Nullable())   // optional, any value
  ```

- **`Rules.Any` field** — supports the `Any()` short-circuit in validate logic.

### Improved

- Element-level `CustomMsg` (e.g. `OnMin` / `OnMax`) inside primitive lists now propagates correctly. Previously `List(Str().WithMax(3).WithMsg(CustomMsg{OnMax: ...}))` fell through to the default error message; now the custom message fires with `${field}`, `${actual_length}`, `${expected_max_length}` template variables.

### Documentation

- New "Whitelist Binding & Escape Hatches" section in `README.md` explicitly documents that fields without rules are stripped at `Bind()`. This is intentional (mass-assignment protection) but was previously undocumented — a common AI-generated code footgun.
- `AI_GUIDE.md` adds two new MUST rules (slot 8 and 9 in BEST PRACTICES CHECKLIST) covering `List` for slices and `Any` for heterogeneous fields. New "🛡 Whitelist Binding (Important)" subsection with diagnostic flow for "data hilang setelah ValidateJSON" reports.
- `llms.txt` and `llms-full.txt` updated: `List` and `Any` added to constructor catalog. `llms-full.txt` adds dedicated "Whitelist Binding (CRITICAL — common AI mistake)" section.

### Migration notes

- No code changes required for existing callers. All additions are opt-in.
- If you have struct fields like `Tags []string` that were silently dropped before, declare them now with `SetRule("tags", List(Str()))` (or `List(Str().WithMax(N))`).
- If you have heterogeneous fields (e.g. `Metadata map[string]interface{}`) that were dropped, declare with `SetRule("metadata", Any())` or `Any().Nullable()`.

## [v0.0.41]

All changes are additive on the public API — existing usage patterns keep
working without modification. Two behavior shifts (panics replaced by
errors) are non-breaking for callers that handle errors normally.

### Added

- **`ValidateJSON[T any](r *http.Request, rules RulesWrapper) (T, error)`** — one-shot generic helper that collapses the `Load → Validate → Bind` pipeline into a single call. Extension lifecycle hooks continue to run internally.
- **Short constructors**: `Str()`, `Int()`, `Int64()`, `Float64()`, `Bool()`, `Email()`, `UUID()`, `IPv4()`, `StrEnum(items…)`, `IntEnum(items…)`, `NestedObject(w)`, `ListOfObject(w)`.
- **Chain helpers on `Rules`** (value-receiver, each returns a new copy): `.Nullable()`, `.Default(v)`, `.WithMin(n)`, `.WithMax(n)`, `.Between(min, max)`, `.Regex(p)`, `.WithMsg(cm)`, `.UniqueFrom(…)`, `.WithRequiredIf(…)`, `.WithRequiredWithout(…)`.
- **`Done()` on `RulesWrapper`** — optional chain terminator so `BuildRoles().SetRule(…).Done()` reads cleanly.
- **`OnEnumValueNotMatch`** field in `CustomMsg` for custom enum-mismatch messages, plus two new template variables: `${actual_value}` and `${enum_values}`.
- **`ErrNoRules`** sentinel error — returned by `Load` / `LoadJsonHttp` / `LoadFormHttp` when no rules are set.
- GoDoc and runnable `ExampleValidateJSON` function.
- Dedicated `test/validate_json_test.go` with 13 scenarios, including a concurrent-execution regression test (20 goroutines, race-detector clean).

### Changed

- **`SetRules(empty)` no longer panics.** The error now surfaces as `ErrNoRules` from the subsequent `Load*` call — easier to handle uniformly across HTTP handlers.
- **Enum with an unsupported element kind no longer panics.** A normal validation error is returned instead.
- **Rules are now safe to share across handlers and goroutines.** The previous mutable per-call state inside `rulesWrapper` has been moved to a per-invocation `wrapperRunState`. A single `rules` value can be declared as a package-level variable and reused — including concurrent requests. The `AI_GUIDE` "no shared rules, inline only" guideline is retired accordingly.
- README and `AI_GUIDE.md` updated: Quick Start and preferred-style sections now showcase `ValidateJSON[T]` and short constructors; legacy 5-step pipeline remains documented as the advanced path.

### Fixed

- **`ListObject` items now bind with real values.** Previously, items in `ListObject` arrays bound as zero-value structs due to a key-shadowing bug inside `validateRecursive` (local `chainKey` shadowed the global `chainKey` constant used by `ToMap`). Renamed the local to `nodeKey`; `tmpChain` now correctly uses the global root key.

### Internal

- Removed state methods from `RulesWrapper` interface and `rulesWrapper` struct (`getFilledField`, `setFilledField`, `appendFilledField`, `getNullFields`, `setNullFields`, `appendNullFields`, `getRequiredWithout`, `setRequiredWithout`, `getRequiredIf`, `setRequiredIf`, `getUniqueValues`, `setUniqueValues`). These were never callable from outside the package.
- Added `*.log` to `.gitignore`.

### Migration notes

- No code changes required for existing callers.
- If your code wraps `SetRules` in a panic-recovery, note that the panic path is gone; handle `ErrNoRules` from `Load*` instead.
- If you relied on enum validation panicking for unsupported element kinds (unlikely), the error is now returned like any other validation error.
