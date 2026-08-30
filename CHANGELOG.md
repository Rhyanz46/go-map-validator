# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to incremental patch versioning (`v0.0.x`).

## [v0.0.47]

A bug-fix release. No public API additions. Four non-breaking fixes: two
correctness fixes (integer coercion from maps, `uint64` message overflow), one
defensive guard, and dead-code cleanup.

### Fixed

- **Integer rules now coerce sibling integer kinds from `Load(map)`.** `Int()` / `Int64()` rejected values such as `int64`, `int32`, `uint`, or `uint64` coming from plain Go maps, reporting `"should be 'int'"` even though the value was a valid integer. This was inconsistent with JSON/form sources (which decode numbers as `float64`) and with `IntEnum`. Integer kinds now cross-coerce from map sources; **floats are still rejected** for integer rules from maps (integer-only coercion — no float→int loosening).
- **`uint64` values in `${actual_length}` no longer wrap negative.** A huge `uint64` above `math.MaxInt64` shown through a custom Min/Max message displayed as a negative number (e.g. `-9223372036854775709`). The message now reports the correct unsigned value.
- **`GetData()` no longer panics on a nil-data state.** Calling `GetData()` on a zero-value `ExtraOperationData` dereferenced a nil pointer; it now returns an empty map, consistent with `GetFilledField()` / `GetNullField()`.

### Internal

- New helpers: `isIntegerKind` (integer kinds without floats) and `integerCoercion` (the shared JSON-vs-map coercion decision).
- `MessageMeta.ActualLength` widened from `*int64` to `*interface{}` so Min/Max messages can carry either signed or unsigned lengths.
- Removed dead code: `isEqualInt` / `isEqualInt64` / `isEqualFloat64`, the unused `convertValue` function, and an unreachable boolean branch.
- Added regression tests in `test/bugfix_v0047_test.go` (5 tests, one table-driven).

### Migration notes

- No code changes required for existing callers.
- `Int()` / `Int64()` rules now accept integer-kind values (e.g. `int64`, `uint`) from `Load(map)` that previously errored — this only makes previously-rejected payloads pass, so it is non-breaking.

## [v0.0.46]

A bug-fix release. No public API additions. Ten defects fixed across two
themes: rules that were silently not enforced, and panics/incorrect results
on edge-case inputs. One behavioral correction (nullable `Bool` null-field
reporting) changes `GetFilledField()` output for a previously-misclassified
case — see Migration notes.

### Fixed

- **Forms no longer reject flag-based rules.** `LoadFormHttp` gated on `rule.Type` being `String`/`Bool`/integer-family, so `Email()`, `UUID()`, `IPv4()`, `.Regex()`, and enum-only rules — whose `Type` is the zero value because they validate string *content* — always failed with `ErrUnsupportType`. They are now accepted and validated as strings.
- **Float values no longer bypass `Max` via truncation.** Min/Max comparisons converted floats with `int64(v)`, so `3.9` passed `WithMax(3)` — every value in `(max, max+1)` was wrongly accepted. Floats are now compared as `float64`.
- **`uint64` values above `math.MaxInt64` no longer bypass Min/Max.** `int64(v.Uint())` wrapped huge values to negative, so a gigantic `uint64` passed `WithMax(100)` and could wrongly fail a small `WithMin`. Unsigned values are now compared in `uint64` space.
- **Manipulators declared on `ListObject` item rules now run.** Each list item was validated on a detached chain that the root-level `RunManipulator()` traversal never reached, so item manipulators were silently skipped.
- **`Unique` across sibling fields inside `ListObject` items is now enforced.** Same root cause as the manipulator skip — the unique checker never visited item chains.
- **`ListRules.Unique` is now enforced.** The field was settable via `BuildListRoles().SetListRule(ListRules{Unique: true})` but had no reader, so duplicate list elements were never rejected. Duplicates are detected with `reflect.DeepEqual` and honor `CustomMsg.OnUnique`.
- **Named string/int types no longer panic.** Hard assertions like `data.(string)` panicked on values such as `type ID string` even though their `reflect.Kind` matched the rule. Affected paths: `Regex`, `Email`, string Min/Max, `StrEnum`/`IntEnum`, and primitive list element checks. All extraction is now kind-safe (`UUID`/`IPv4` were already safe).
- **Integer enums now coerce from any source, not just JSON.** `IntEnum(1, 2, 3)` previously rejected `int64(2)` from `Load(map)` because cross-type integer coercion only applied to HTTP JSON/form sources. Integer-family enum values are now compared numerically regardless of concrete type or source.
- **Enum type-mismatch error no longer prints `list<nil>`.** The mismatch path passed `nil` instead of the enum values into the message builder.
- **A rule carrying both `Object` and `ListObject` no longer panics.** The conflicting configuration now returns a validation error instead of a type-assertion panic.

### Behavior change

- **A missing nullable `Bool` is now reported as a null field.** Previously `Bool().Nullable()` returned `false` (non-nil) for an absent field, so `GetFilledField()` included it and `GetNullField()` did not — inconsistent with every other type. It now returns `nil`, so the field appears in `GetNullField()`.

### Internal

- New helpers: `asString` (kind-safe string extraction), `numericAsFloat64`, `isUintKind`.
- The enum validation block (~158 lines of per-kind duplication with a JSON-only coercion path) was replaced by a single unified numeric comparison (~60 lines).
- `ListObject` item validation now runs `RunManipulator()` and `RunUniqueChecker()` locally on each item chain.
- Added regression tests in `test/bugfix_v0046_test.go` (6 tests) and `test/bugfix2_v0046_test.go` (4 tests, one table-driven with 7 sub-cases). Full suite: 167 tests green, race-detector clean; no existing test required modification.

### Migration notes

- No code changes required for existing callers.
- Payloads that previously *passed incorrectly* are now *correctly rejected*: floats above `Max`, huge `uint64` values outside Min/Max, duplicate elements in `ListRules.Unique` lists, and duplicate values across `UniqueFrom` siblings inside `ListObject` items.
- Payloads that previously *failed incorrectly* now pass: flag-based rules in forms, named string/int types, and integer-family enum values of a different concrete type (e.g. `int64` against `IntEnum`).
- **Action may be needed** if you read `GetFilledField()` for nullable `Bool` fields: absent bools move from the filled list to the null list. Adjust any logic that depended on the old classification.

## [v0.0.45]

A bug-fix release. No public API additions; two behavior corrections are
non-breaking (they fix paths that previously errored or panicked).

### Fixed

- **`.Default(v)` now works without an explicit `.Nullable()`.** Previously `Str().Default("guest")` still errored with `"we need 'x' field"` when the field was absent, so the default was never applied. `Default` now implies the field is optional. Non-breaking: present fields behave identically, and the previously-erroring case was unusable.
- **Integer / float / bool fields from HTTP forms now validate.** `LoadFormHttp` stored every value as a string but type coercion only ran for JSON sources, so `Int`, `Float64`, and `Bool` form fields always failed with `"should be 'int'/'bool'"` — effectively only `String` worked. Form values are now normalized (numbers → `float64`, bools → `bool`) and run through the same coercion path as JSON. `Float32`/`Float64` and the full integer family are now accepted form types. Invalid input (e.g. `"abc"` for an `Int`) still yields a clear type error.
- **Unique check no longer panics on non-comparable values.** Comparing sibling values with `==` panicked when a field held a slice/map (e.g. via `Any()`). The comparison now uses `reflect.DeepEqual`; equal slices are still flagged as violations.
- **`${unique_origin}` / `${unique_target}` template guard corrected.** `buildMessage` checked `meta.Field` for nil but dereferenced `meta.UniqueOrigin` / `meta.UniqueTarget` — a latent nil-deref. Now guards the correct pointers.

### Internal

- New `coercesNumbers(dataFrom)` helper centralizes the JSON-like-source check and includes multipart/urlencoded forms.
- New `coerceFormValue(value, kind)` helper for form value normalization.
- Added `test/form_test.go` (first-ever `LoadFormHttp` coverage) and regression tests: `TestDefaultImpliesNullable`, `TestUniqueNonComparableValues`.

### Migration notes

- No code changes required for existing callers.
- If you declared `Int`/`Bool`/`Float` rules for form handlers and worked around them by treating everything as `String`, you can now declare the real types.

## [v0.0.44]

A bug-fix release. No public API changes.

### Fixed

- **Negative integers no longer bypass `Min` / `Max` (`WithMin` / `WithMax` / `Between`).** `extractInteger` previously matched digits with the regex `[0-9]+`, which silently dropped the leading sign and concatenated separate digit groups. As a result `-5` was treated as `5`: it passed `WithMin(1)` and was wrongly rejected by `WithMax(3)`. The helper now reads the numeric value directly via reflection (int / uint / float families), correctly preserving sign and magnitude — including large floats that render in scientific notation.

### Internal

- Removed the now-unused `removeAfter` string helper and the `fmt.Sprintf` + `e+` stripping at the integer `Min`/`Max` call sites.
- Added regression tests: `TestNegativeIntegerFailsWithMin`, `TestNegativeIntegerPassesWithMax`, `TestNegativeIntegerBetween`.

### Migration notes

- No code changes required. Callers that (incorrectly) relied on negative integers passing a positive `Min` will now see those values correctly rejected.

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
