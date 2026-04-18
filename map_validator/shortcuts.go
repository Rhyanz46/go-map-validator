package map_validator

import "reflect"

// Shortcut constructors and chain helpers to reduce boilerplate at call sites.
//
// Example:
//
//	rules := map_validator.BuildRoles().
//	    SetRule("email", map_validator.Email().WithMax(255)).
//	    SetRule("name", map_validator.Str().Between(1, 100)).
//	    SetRule("role", map_validator.StrEnum("admin", "guest").Nullable()).
//	    Done()
//
// These return Rules by value so chain calls stay composable without aliasing.

// --- Type constructors ---

func Str() Rules     { return Rules{Type: reflect.String} }
func Int() Rules     { return Rules{Type: reflect.Int} }
func Int64() Rules   { return Rules{Type: reflect.Int64} }
func Float64() Rules { return Rules{Type: reflect.Float64} }
func Bool() Rules    { return Rules{Type: reflect.Bool} }
func Email() Rules   { return Rules{Email: true} }
func UUID() Rules    { return Rules{UUID: true} }
func IPv4() Rules    { return Rules{IPV4: true} }

// --- Enum constructors (base-typed to avoid reflect gymnastics at call site) ---

func StrEnum(items ...string) Rules {
	return Rules{Type: reflect.String, Enum: &EnumField[any]{Items: items}}
}

func IntEnum(items ...int) Rules {
	return Rules{Type: reflect.Int, Enum: &EnumField[any]{Items: items}}
}

// --- Nesting shortcuts ---

func NestedObject(w RulesWrapper) Rules { return Rules{Object: w} }
func ListOfObject(w RulesWrapper) Rules { return Rules{ListObject: w} }

// --- Chain helpers on Rules ---
//
// Methods use a value receiver so each call returns a modified copy; the
// original Rules value is never mutated. Use `With*` prefix wherever the
// method name would collide with an existing struct field.

func (r Rules) Nullable() Rules             { r.Null = true; return r }
func (r Rules) Default(v interface{}) Rules { r.IfNull = v; return r }
func (r Rules) WithMin(n int64) Rules       { r.Min = SetTotal(n); return r }
func (r Rules) WithMax(n int64) Rules       { r.Max = SetTotal(n); return r }
func (r Rules) Between(min, max int64) Rules {
	r.Min = SetTotal(min)
	r.Max = SetTotal(max)
	return r
}
func (r Rules) Regex(pattern string) Rules        { r.RegexString = pattern; return r }
func (r Rules) WithMsg(m CustomMsg) Rules         { r.CustomMsg = m; return r }
func (r Rules) UniqueFrom(fields ...string) Rules { r.Unique = fields; return r }
func (r Rules) WithRequiredIf(fields ...string) Rules {
	r.RequiredIf = fields
	return r
}
func (r Rules) WithRequiredWithout(fields ...string) Rules {
	r.RequiredWithout = fields
	return r
}
