package fate

import (
	"reflect"
	"testing"
)

func TestDefaultFilterOptionIsBounded(t *testing.T) {
	o := DefaultFilterOption()
	if o.MinStroke != 3 || o.MaxStroke != 18 {
		t.Fatalf("stroke range = [%d,%d], want [3,18]", o.MinStroke, o.MaxStroke)
	}
	if o.MinStroke <= 0 || o.MaxStroke <= 0 || o.MinStroke > o.MaxStroke {
		t.Fatalf("stroke range is not bounded: %+v", o)
	}
	if !o.CharacterFilter || !o.RegularFilter || !o.DaYanFilter || !o.WuXingFilter {
		t.Fatalf("bounded default must enable candidate-pruning filters: %+v", o)
	}
	if DefaultFilter() == nil {
		t.Fatal("DefaultFilter() is nil; the Session path depends on it")
	}
}

func TestResolveFilterOption(t *testing.T) {
	// The nil policy is the regression guard: before the fix, nil resolved to
	// the permissive shared DefaultFilter and never converged.
	if got := resolveFilterOption(nil); !reflect.DeepEqual(got, DefaultFilterOption()) {
		t.Fatalf("nil Filter resolved to %+v, want DefaultFilterOption()", got)
	}

	custom := DefaultFilterOption().WithAvoidCharacters("病", "死")
	if got := resolveFilterOption(&custom); !reflect.DeepEqual(got, custom) {
		t.Fatalf("explicit Filter was not preserved: got %+v, want %+v", got, custom)
	}
}
