package test

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/Rhyanz46/go-map-validator/map_validator"
)

// Bug A: integer rules (Int, Int64, ...) rejected sibling integer kinds from
// Load(map). JSON/form already coerce numbers (decoded as float64); map data
// carries exact kinds (int64 from a DB timestamp, etc.) and was rejected with
// "should be 'int'" even though the value is a valid integer.
func TestIntegerCoercionFromMap(t *testing.T) {
	cases := []struct {
		name string
		rule map_validator.Rules
		data map[string]interface{}
	}{
		{"int64-for-int", map_validator.Int(), map[string]interface{}{"v": int64(30)}},
		{"int32-for-int", map_validator.Int(), map[string]interface{}{"v": int32(30)}},
		{"uint-for-int", map_validator.Int(), map[string]interface{}{"v": uint(30)}},
		{"uint64-for-int", map_validator.Int(), map[string]interface{}{"v": uint64(30)}},
		{"int-for-int64", map_validator.Int64(), map[string]interface{}{"v": int(30)}},
		{"uint64-for-int64", map_validator.Int64(), map[string]interface{}{"v": uint64(30)}},
		{"uint-for-int64", map_validator.Int64(), map[string]interface{}{"v": uint(30)}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rules := buildRules1("v", tc.rule)
			op, err := map_validator.NewValidateBuilder().SetRules(rules).Load(tc.data)
			if err != nil {
				t.Fatalf("load: %s", err)
			}
			if _, err := op.RunValidate(); err != nil {
				t.Fatalf("expected %s to coerce, got: %s", tc.name, err)
			}
		})
	}
}

// List(Int()) elements should also cross-coerce integer kinds from map
// (the list pre-check had the same JSON-only gate).
func TestListIntegerCoercionFromMap(t *testing.T) {
	rules := buildRules1("nums", map_validator.List(map_validator.Int()))
	op, err := map_validator.NewValidateBuilder().SetRules(rules).
		Load(map[string]interface{}{"nums": []int64{1, 2, 3}})
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	if _, err := op.RunValidate(); err != nil {
		t.Fatalf("expected []int64 to coerce for List(Int()) from map: %s", err)
	}
}

// Guard: integer coercion must stay integer-only. A float for an Int() rule
// from map must still be rejected (this is the breaking-change boundary).
func TestIntFromMapStillRejectsFloat(t *testing.T) {
	cases := []struct {
		name string
		rule map_validator.Rules
		val  interface{}
	}{
		{"float64-for-int", map_validator.Int(), 3.9},
		{"float32-for-int", map_validator.Int(), float32(3)},
		{"float64-for-int64", map_validator.Int64(), 3.9},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rules := buildRules1("v", tc.rule)
			op, err := map_validator.NewValidateBuilder().SetRules(rules).
				Load(map[string]interface{}{"v": tc.val})
			if err != nil {
				t.Fatalf("load: %s", err)
			}
			if _, err := op.RunValidate(); err == nil {
				t.Fatalf("expected %v to be rejected for an integer rule from map", tc.val)
			}
		})
	}
}

// Bug B: the ${actual_length} in a custom Min/Max message wrapped huge uint64
// values (above math.MaxInt64) to a negative number.
func TestUint64MessageNoOverflow(t *testing.T) {
	big := uint64(math.MaxInt64) + 100
	want := fmt.Sprintf("%d", big)

	rules := map_validator.BuildRoles().SetRule("big", map_validator.Rules{
		Type: reflect.Uint64,
		Max:  map_validator.SetTotal(100),
		CustomMsg: map_validator.CustomMsg{
			OnMax: map_validator.SetMessage("got ${actual_length}"),
		},
	}).Done()
	op, err := map_validator.NewValidateBuilder().SetRules(rules).
		Load(map[string]interface{}{"big": big})
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	_, err = op.RunValidate()
	if err == nil {
		t.Fatal("expected huge uint64 to fail Max=100")
	}
	if strings.Contains(err.Error(), "-") {
		t.Fatalf("actual_length wrapped negative: %s", err)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("expected message to contain %q, got: %s", want, err)
	}
}

// Bug E: GetData() on a zero-value ExtraOperationData (nil data) panicked on
// dereference, unlike GetFilledField/GetNullField which guard against empty.
func TestGetDataNilSafe(t *testing.T) {
	d := &map_validator.ExtraOperationData{}
	got := d.GetData()
	if got == nil {
		t.Fatal("expected GetData on a nil-data state to return an empty map, not nil")
	}
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}
