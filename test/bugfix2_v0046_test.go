package test

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

type namedString string
type namedInt int

// Bug 1: named string/int types panicked on hard type assertions
// even though their reflect.Kind matches the declared rule.
func TestNamedTypesNoPanic(t *testing.T) {
	cases := []struct {
		name  string
		rules map_validator.RulesWrapper
		data  map[string]interface{}
	}{
		{"regex", buildRules1("code", map_validator.Str().Regex("^[a-z]+$")), map[string]interface{}{"code": namedString("abc")}},
		{"email", buildRules1("email", map_validator.Email()), map[string]interface{}{"email": namedString("a@b.co")}},
		{"str-max", buildRules1("name", map_validator.Str().WithMax(10)), map[string]interface{}{"name": namedString("abc")}},
		{"str-min", buildRules1("name", map_validator.Str().WithMin(2)), map[string]interface{}{"name": namedString("abc")}},
		{"str-enum", buildRules1("role", map_validator.StrEnum("admin", "guest")), map[string]interface{}{"role": namedString("admin")}},
		{"int-enum", buildRules1("level", map_validator.IntEnum(1, 2, 3)), map[string]interface{}{"level": namedInt(2)}},
		{"list-str-max", buildRules1("tags", map_validator.List(map_validator.Str().WithMax(10))), map[string]interface{}{"tags": []interface{}{namedString("abc")}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panicked: %v", r)
				}
			}()
			op, err := map_validator.NewValidateBuilder().SetRules(tc.rules).Load(tc.data)
			if err != nil {
				t.Fatalf("load: %s", err)
			}
			if _, err := op.RunValidate(); err != nil {
				t.Fatalf("expected named type to validate cleanly: %s", err)
			}
		})
	}
}

func buildRules1(field string, rule map_validator.Rules) map_validator.RulesWrapper {
	return map_validator.BuildRoles().SetRule(field, rule).Done()
}

// Bug 2: uint64 values above MaxInt64 wrapped to negative in Min/Max checks.
func TestUint64NoOverflow(t *testing.T) {
	big := uint64(math.MaxInt64) + 100

	rules := map_validator.BuildRoles().SetRule("big", map_validator.Rules{
		Type: reflect.Uint64,
		Max:  map_validator.SetTotal(100),
	}).Done()
	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{"big": big})
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	if _, err := op.RunValidate(); err == nil {
		t.Fatal("expected huge uint64 to fail Max=100")
	}

	// huge value must still pass when within bounds
	rules2 := map_validator.BuildRoles().SetRule("big", map_validator.Rules{
		Type: reflect.Uint64,
		Min:  map_validator.SetTotal(100),
	}).Done()
	op2, err := map_validator.NewValidateBuilder().SetRules(rules2).Load(map[string]interface{}{"big": big})
	if err != nil {
		t.Fatalf("load2: %s", err)
	}
	if _, err := op2.RunValidate(); err != nil {
		t.Fatalf("expected huge uint64 to pass Min=100: %s", err)
	}
}

// Bug 3: IntEnum rejected int64 values from Load(map) and printed "list<nil>".
func TestIntEnumCoercionFromMap(t *testing.T) {
	rules := buildRules1("level", map_validator.IntEnum(1, 2, 3))
	op, err := map_validator.NewValidateBuilder().SetRules(rules).
		Load(map[string]interface{}{"level": int64(2)})
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	if _, err := op.RunValidate(); err != nil {
		t.Fatalf("expected int64(2) to match IntEnum(1,2,3): %s", err)
	}

	// out-of-range value must fail with a message that lists the enum values
	op2, err := map_validator.NewValidateBuilder().SetRules(rules).
		Load(map[string]interface{}{"level": int64(9)})
	if err != nil {
		t.Fatalf("load2: %s", err)
	}
	if _, err := op2.RunValidate(); err == nil {
		t.Fatal("expected 9 to be rejected")
	} else if strings.Contains(err.Error(), "<nil>") {
		t.Fatalf("error message contains <nil>: %s", err)
	}
}

// Bug 4: a missing nullable Bool was reported as a filled field.
func TestNullableBoolGoesToNullField(t *testing.T) {
	rules := buildRules1("active", map_validator.Bool().Nullable())
	op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(map[string]interface{}{})
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	extra, err := op.RunValidate()
	if err != nil {
		t.Fatalf("validate: %s", err)
	}
	if len(extra.GetFilledField()) != 0 {
		t.Fatalf("expected no filled fields, got %v", extra.GetFilledField())
	}
	nullFields := extra.GetNullField()
	if len(nullFields) != 1 || nullFields[0] != "active" {
		t.Fatalf("expected null field [active], got %v", nullFields)
	}
}
