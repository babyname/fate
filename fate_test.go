package fate

import (
	"testing"

	"github.com/babyname/fate/config"
	_ "github.com/sqlite3ent/sqlite3"
)

func TestNew(t *testing.T) {
	cfg := config.DefaultConfig()
	got, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got == nil {
		t.Fatal("New() returned nil")
	}

	s := got.NewSession()
	if s == nil {
		t.Fatal("NewSession() returned nil")
	}

	s2 := got.NewSessionWithFilter(NewFilter(FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: CharacterFilterTypeDefault,
		MinStroke:           3,
		MaxStroke:           18,
		SexFilter:           true,
	}))
	if s2 == nil {
		t.Fatal("NewSessionWithFilter() returned nil")
	}
}

func TestFilterOption(t *testing.T) {
	fo := FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: CharacterFilterTypeChs,
		MinStroke:           5,
		MaxStroke:           20,
		RegularFilter:       true,
		DaYanFilter:         true,
		WuXingFilter:        true,
		SexFilter:           false,
	}
	f := NewFilter(fo)
	if f == nil {
		t.Fatal("NewFilter() returned nil")
	}
	if f.FilterType() != CharacterFilterTypeChs {
		t.Errorf("FilterType() = %v, want %v", f.FilterType(), CharacterFilterTypeChs)
	}
}

func TestDefaultFilter(t *testing.T) {
	f := DefaultFilter()
	if f == nil {
		t.Fatal("DefaultFilter() returned nil")
	}
	if f.FilterType() != CharacterFilterTypeDefault {
		t.Errorf("FilterType() = %v, want %v", f.FilterType(), CharacterFilterTypeDefault)
	}
}
