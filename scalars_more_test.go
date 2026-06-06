package env

import (
	"strings"
	"testing"
	"time"
)

// TestNumericTypes exercises the present/absent/malformed paths for every
// numeric defaulted accessor through a single table.
func TestNumericTypes(t *testing.T) {
	tests := []struct {
		name      string
		raw       string // value for key "K"; "" with set=false means absent
		set       bool
		check     func(t *testing.T, p *Parser) (got any)
		wantValue any
	}{
		{name: "int32 present", raw: "100", set: true, wantValue: int32(100),
			check: func(t *testing.T, p *Parser) any { return p.Int32("K", 7) }},
		{name: "int32 overflow uses default", raw: "9999999999", set: true, wantValue: int32(7),
			check: func(t *testing.T, p *Parser) any { return p.Int32("K", 7) }},
		{name: "int32 absent uses default", set: false, wantValue: int32(7),
			check: func(t *testing.T, p *Parser) any { return p.Int32("K", 7) }},
		{name: "int64 present", raw: "9999999999", set: true, wantValue: int64(9999999999),
			check: func(t *testing.T, p *Parser) any { return p.Int64("K", 7) }},
		{name: "uint present", raw: "42", set: true, wantValue: uint(42),
			check: func(t *testing.T, p *Parser) any { return p.Uint("K", 7) }},
		{name: "uint negative uses default", raw: "-1", set: true, wantValue: uint(7),
			check: func(t *testing.T, p *Parser) any { return p.Uint("K", 7) }},
		{name: "uint64 present", raw: "42", set: true, wantValue: uint64(42),
			check: func(t *testing.T, p *Parser) any { return p.Uint64("K", 7) }},
		{name: "float64 present", raw: "3.14", set: true, wantValue: 3.14,
			check: func(t *testing.T, p *Parser) any { return p.Float64("K", 1.0) }},
		{name: "float64 malformed uses default", raw: "pi", set: true, wantValue: 1.0,
			check: func(t *testing.T, p *Parser) any { return p.Float64("K", 1.0) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := MapLookuper{}
			if tt.set {
				src["K"] = tt.raw
			}
			p := From(src)
			if got := tt.check(t, p); got != tt.wantValue {
				t.Errorf("got %v (%T), want %v (%T)", got, got, tt.wantValue, tt.wantValue)
			}
			if err := p.Err(); err != nil {
				t.Errorf("Err = %v, want nil", err)
			}
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		raw  string
		set  bool
		def  bool
		want bool
	}{
		{raw: "true", set: true, def: false, want: true},
		{raw: "0", set: true, def: true, want: false},
		{raw: "T", set: true, def: false, want: true},
		{raw: "nope", set: true, def: true, want: true}, // malformed -> default
		{set: false, def: true, want: true},             // absent -> default
	}
	for _, tt := range tests {
		src := MapLookuper{}
		if tt.set {
			src["K"] = tt.raw
		}
		p := From(src)
		if got := p.Bool("K", tt.def); got != tt.want {
			t.Errorf("Bool(%q,set=%v,def=%v) = %v, want %v", tt.raw, tt.set, tt.def, got, tt.want)
		}
	}
}

func TestRequiredBool_Malformed(t *testing.T) {
	p := From(MapLookuper{"K": "yes-please"})
	p.RequiredBool("K")
	err := p.Err()
	if err == nil || !strings.Contains(err.Error(), "invalid bool") {
		t.Errorf("Err = %v, want invalid bool error", err)
	}
}

func TestDuration(t *testing.T) {
	p := From(MapLookuper{"D": "1h30m"})
	if got := p.Duration("D", time.Second); got != 90*time.Minute {
		t.Errorf("Duration = %v, want 1h30m", got)
	}

	p = From(MapLookuper{})
	if got := p.Duration("D", 5*time.Second); got != 5*time.Second {
		t.Errorf("Duration default = %v, want 5s", got)
	}

	p = From(MapLookuper{"D": "abc"})
	p.RequiredDuration("D")
	if p.Err() == nil {
		t.Error("RequiredDuration: Err = nil, want parse error")
	}
}

func TestTime(t *testing.T) {
	want := time.Date(2026, 6, 6, 12, 0, 0, 0, time.UTC)

	p := From(MapLookuper{"T": "2026-06-06T12:00:00Z"})
	if got := p.RequiredTime("T"); !got.Equal(want) {
		t.Errorf("RequiredTime = %v, want %v", got, want)
	}
	if err := p.Err(); err != nil {
		t.Errorf("Err = %v, want nil", err)
	}

	// Custom layout.
	p = From(MapLookuper{"T": "2026-06-06"})
	got := p.RequiredTimeLayout("T", "2006-01-02")
	if !got.Equal(time.Date(2026, 6, 6, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("RequiredTimeLayout = %v", got)
	}

	// Malformed default path.
	def := time.Unix(0, 0).UTC()
	p = From(MapLookuper{"T": "not-a-time"})
	if got := p.Time("T", def); !got.Equal(def) {
		t.Errorf("Time malformed = %v, want default %v", got, def)
	}
	if err := p.Err(); err != nil {
		t.Errorf("Time defaulted: Err = %v, want nil", err)
	}
}
