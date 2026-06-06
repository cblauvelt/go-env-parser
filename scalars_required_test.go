package env

import (
	"strings"
	"testing"
)

// TestRequiredNumericTypes covers the Required* variants for every numeric type
// that was only exercised via the defaulted path in TestNumericTypes.
func TestRequiredNumericTypes(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		set       bool
		call      func(p *Parser) any
		wantValue any
		wantErr   string // non-empty means we expect an error containing this substring
	}{
		// int32
		{name: "RequiredInt32 present", raw: "100", set: true, wantValue: int32(100),
			call: func(p *Parser) any { return p.RequiredInt32("K") }},
		{name: "RequiredInt32 absent", set: false, wantValue: int32(0), wantErr: "required",
			call: func(p *Parser) any { return p.RequiredInt32("K") }},
		{name: "RequiredInt32 overflow", raw: "9999999999", set: true, wantValue: int32(0), wantErr: "invalid int32",
			call: func(p *Parser) any { return p.RequiredInt32("K") }},
		// int64
		{name: "RequiredInt64 present", raw: "9999999999", set: true, wantValue: int64(9999999999),
			call: func(p *Parser) any { return p.RequiredInt64("K") }},
		{name: "RequiredInt64 absent", set: false, wantValue: int64(0), wantErr: "required",
			call: func(p *Parser) any { return p.RequiredInt64("K") }},
		{name: "RequiredInt64 malformed", raw: "abc", set: true, wantValue: int64(0), wantErr: "invalid int64",
			call: func(p *Parser) any { return p.RequiredInt64("K") }},
		// uint
		{name: "RequiredUint present", raw: "42", set: true, wantValue: uint(42),
			call: func(p *Parser) any { return p.RequiredUint("K") }},
		{name: "RequiredUint absent", set: false, wantValue: uint(0), wantErr: "required",
			call: func(p *Parser) any { return p.RequiredUint("K") }},
		{name: "RequiredUint negative", raw: "-1", set: true, wantValue: uint(0), wantErr: "invalid uint",
			call: func(p *Parser) any { return p.RequiredUint("K") }},
		// uint64
		{name: "RequiredUint64 present", raw: "42", set: true, wantValue: uint64(42),
			call: func(p *Parser) any { return p.RequiredUint64("K") }},
		{name: "RequiredUint64 absent", set: false, wantValue: uint64(0), wantErr: "required",
			call: func(p *Parser) any { return p.RequiredUint64("K") }},
		{name: "RequiredUint64 malformed", raw: "x", set: true, wantValue: uint64(0), wantErr: "invalid uint64",
			call: func(p *Parser) any { return p.RequiredUint64("K") }},
		// float64
		{name: "RequiredFloat64 present", raw: "2.718", set: true, wantValue: 2.718,
			call: func(p *Parser) any { return p.RequiredFloat64("K") }},
		{name: "RequiredFloat64 absent", set: false, wantValue: float64(0), wantErr: "required",
			call: func(p *Parser) any { return p.RequiredFloat64("K") }},
		{name: "RequiredFloat64 malformed", raw: "inf-ish", set: true, wantValue: float64(0), wantErr: "invalid float64",
			call: func(p *Parser) any { return p.RequiredFloat64("K") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := MapLookuper{}
			if tt.set {
				src["K"] = tt.raw
			}
			p := From(src)
			got := tt.call(p)
			if got != tt.wantValue {
				t.Errorf("got %v (%T), want %v (%T)", got, got, tt.wantValue, tt.wantValue)
			}
			err := p.Err()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("Err = %v, want nil", err)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Err = %v, want error containing %q", err, tt.wantErr)
				}
			}
		})
	}
}

// TestTimeLayout_DefaultedVariant exercises the non-required TimeLayout path.
func TestTimeLayout_DefaultedVariant(t *testing.T) {
	p := From(MapLookuper{"D": "2026-06-06"})
	want := "2026-06-06"
	got := p.TimeLayout("D", "2006-01-02", zeroTime)
	if got.Format("2006-01-02") != want {
		t.Errorf("TimeLayout = %v, want %s", got, want)
	}
	if err := p.Err(); err != nil {
		t.Errorf("Err = %v, want nil", err)
	}
}

func TestTimeLayout_DefaultedVariant_Absent(t *testing.T) {
	p := From(MapLookuper{})
	got := p.TimeLayout("D", "2006-01-02", zeroTime)
	if !got.Equal(zeroTime) {
		t.Errorf("TimeLayout absent = %v, want zero", got)
	}
}
