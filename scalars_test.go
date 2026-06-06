package env

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name string
		src  MapLookuper
		opts []Option
		def  string
		want string
	}{
		{name: "present", src: MapLookuper{"K": "v"}, def: "d", want: "v"},
		{name: "absent uses default", src: MapLookuper{}, def: "d", want: "d"},
		{name: "empty is present by default", src: MapLookuper{"K": ""}, def: "d", want: ""},
		{name: "empty as unset uses default", src: MapLookuper{"K": ""}, opts: []Option{WithEmptyAsUnset()}, def: "d", want: "d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := From(tt.src, tt.opts...)
			if got := p.String("K", tt.def); got != tt.want {
				t.Errorf("String = %q, want %q", got, tt.want)
			}
			if err := p.Err(); err != nil {
				t.Errorf("Err = %v, want nil", err)
			}
		})
	}
}

func TestRequiredString(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		p := From(MapLookuper{"K": "v"})
		if got := p.RequiredString("K"); got != "v" {
			t.Errorf("RequiredString = %q, want v", got)
		}
		if err := p.Err(); err != nil {
			t.Errorf("Err = %v, want nil", err)
		}
	})
	t.Run("absent records error", func(t *testing.T) {
		p := From(MapLookuper{})
		if got := p.RequiredString("K"); got != "" {
			t.Errorf("RequiredString = %q, want empty", got)
		}
		err := p.Err()
		if err == nil || !strings.Contains(err.Error(), `"K"`) {
			t.Errorf("Err = %v, want error mentioning K", err)
		}
	})
	t.Run("empty as unset records error", func(t *testing.T) {
		p := From(MapLookuper{"K": ""}, WithEmptyAsUnset())
		p.RequiredString("K")
		if p.Err() == nil {
			t.Error("Err = nil, want missing error")
		}
	})
}

func TestInt(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		p := From(MapLookuper{"N": "42"})
		if got := p.Int("N", 7); got != 42 {
			t.Errorf("Int = %d, want 42", got)
		}
	})
	t.Run("absent uses default", func(t *testing.T) {
		p := From(MapLookuper{})
		if got := p.Int("N", 7); got != 7 {
			t.Errorf("Int = %d, want 7", got)
		}
	})
	t.Run("malformed uses default and fires hook", func(t *testing.T) {
		var gotKey, gotRaw string
		p := From(MapLookuper{"N": "abc"}, WithBadDefaultHook(func(key, raw string, _ error) {
			gotKey, gotRaw = key, raw
		}))
		if got := p.Int("N", 7); got != 7 {
			t.Errorf("Int = %d, want 7", got)
		}
		if gotKey != "N" || gotRaw != "abc" {
			t.Errorf("hook got (%q,%q), want (N,abc)", gotKey, gotRaw)
		}
		if err := p.Err(); err != nil {
			t.Errorf("Err = %v, want nil (defaulted accessors never record)", err)
		}
	})
}

func TestRequiredInt(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		p := From(MapLookuper{"N": "42"})
		if got := p.RequiredInt("N"); got != 42 {
			t.Errorf("RequiredInt = %d, want 42", got)
		}
		if err := p.Err(); err != nil {
			t.Errorf("Err = %v, want nil", err)
		}
	})
	t.Run("absent records missing error", func(t *testing.T) {
		p := From(MapLookuper{})
		p.RequiredInt("N")
		if p.Err() == nil {
			t.Error("Err = nil, want missing error")
		}
	})
	t.Run("malformed records parse error", func(t *testing.T) {
		p := From(MapLookuper{"N": "abc"})
		if got := p.RequiredInt("N"); got != 0 {
			t.Errorf("RequiredInt = %d, want 0", got)
		}
		err := p.Err()
		if err == nil || !strings.Contains(err.Error(), "invalid int") {
			t.Errorf("Err = %v, want invalid int error", err)
		}
	})
}

// TestErr_CollectsAll verifies the batch behavior: multiple failures surface in
// a single Err() rather than stopping at the first.
func TestErr_CollectsAll(t *testing.T) {
	p := From(MapLookuper{"BAD": "abc"})
	p.RequiredString("MISSING_A")
	p.RequiredInt("BAD")
	p.RequiredString("MISSING_B")

	err := p.Err()
	if err == nil {
		t.Fatal("Err = nil, want aggregated error")
	}
	msg := err.Error()
	for _, want := range []string{"MISSING_A", "BAD", "MISSING_B"} {
		if !strings.Contains(msg, want) {
			t.Errorf("Err missing %q in: %v", want, msg)
		}
	}
}
