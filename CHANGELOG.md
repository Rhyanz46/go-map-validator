# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to incremental patch versioning (`v0.0.x`).

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
