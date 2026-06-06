package env

import (
	"net/url"
	"slices"
	"strings"
	"testing"
)

func TestCSV(t *testing.T) {
	tests := []struct {
		name string
		src  MapLookuper
		def  []string
		want []string
	}{
		{name: "splits and trims", src: MapLookuper{"K": "a, b ,c"}, def: nil, want: []string{"a", "b", "c"}},
		{name: "drops empty elements", src: MapLookuper{"K": "a,,b,"}, def: nil, want: []string{"a", "b"}},
		{name: "absent uses default", src: MapLookuper{}, def: []string{"*"}, want: []string{"*"}},
		{name: "set-but-blank overrides default", src: MapLookuper{"K": " , "}, def: []string{"*"}, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := From(tt.src)
			got := p.CSV("K", tt.def)
			if !slices.Equal(got, tt.want) {
				t.Errorf("CSV = %#v, want %#v", got, tt.want)
			}
			if err := p.Err(); err != nil {
				t.Errorf("Err = %v, want nil", err)
			}
		})
	}
}

func TestCSVSep_And_Required(t *testing.T) {
	p := From(MapLookuper{"K": "a:b:c"})
	if got := p.CSVSep("K", ":", nil); !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Errorf("CSVSep = %#v", got)
	}

	p = From(MapLookuper{})
	if got := p.RequiredCSV("K"); got != nil {
		t.Errorf("RequiredCSV = %#v, want nil", got)
	}
	if p.Err() == nil {
		t.Error("RequiredCSV: Err = nil, want missing error")
	}
}

func TestIntSlice(t *testing.T) {
	p := From(MapLookuper{"K": "1, 2 ,3"})
	if got := p.IntSlice("K", nil); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("IntSlice = %#v, want [1 2 3]", got)
	}

	// Malformed element: defaulted path returns def, fires hook, no recorded err.
	var fired bool
	p = From(MapLookuper{"K": "1,x,3"}, WithBadDefaultHook(func(string, string, error) { fired = true }))
	if got := p.IntSlice("K", []int{9}); !slices.Equal(got, []int{9}) {
		t.Errorf("IntSlice malformed = %#v, want [9]", got)
	}
	if !fired {
		t.Error("bad-default hook did not fire")
	}
	if err := p.Err(); err != nil {
		t.Errorf("Err = %v, want nil", err)
	}

	// Malformed element: required path records an error.
	p = From(MapLookuper{"K": "1,x,3"})
	p.RequiredIntSlice("K")
	if err := p.Err(); err == nil || !strings.Contains(err.Error(), "int slice") {
		t.Errorf("RequiredIntSlice: Err = %v, want int slice error", err)
	}
}

func TestOneOf(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}

	p := From(MapLookuper{"LOG": "warn"})
	if got := p.OneOf("LOG", "info", levels...); got != "warn" {
		t.Errorf("OneOf = %q, want warn", got)
	}

	p = From(MapLookuper{})
	if got := p.OneOf("LOG", "info", levels...); got != "info" {
		t.Errorf("OneOf absent = %q, want info", got)
	}

	// Invalid value: defaulted path returns def and fires hook.
	var fired bool
	p = From(MapLookuper{"LOG": "trace"}, WithBadDefaultHook(func(string, string, error) { fired = true }))
	if got := p.OneOf("LOG", "info", levels...); got != "info" {
		t.Errorf("OneOf invalid = %q, want info", got)
	}
	if !fired {
		t.Error("bad-default hook did not fire for invalid OneOf")
	}
	if err := p.Err(); err != nil {
		t.Errorf("Err = %v, want nil", err)
	}
}

func TestRequiredOneOf(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}

	p := From(MapLookuper{"LOG": "warn"})
	if got := p.RequiredOneOf("LOG", levels...); got != "warn" {
		t.Errorf("RequiredOneOf valid = %q, want warn", got)
	}
	if err := p.Err(); err != nil {
		t.Errorf("Err = %v, want nil", err)
	}

	p = From(MapLookuper{"LOG": "trace"})
	if got := p.RequiredOneOf("LOG", levels...); got != "" {
		t.Errorf("RequiredOneOf = %q, want empty", got)
	}
	err := p.Err()
	if err == nil || !strings.Contains(err.Error(), "not one of") {
		t.Errorf("Err = %v, want 'not one of' error", err)
	}

	p = From(MapLookuper{})
	p.RequiredOneOf("LOG", levels...)
	if p.Err() == nil {
		t.Error("RequiredOneOf absent: Err = nil, want missing error")
	}
}

func TestValidated(t *testing.T) {
	parseURL := func(s string) (*url.URL, error) {
		u, err := url.Parse(s)
		if err != nil {
			return nil, err
		}
		if u.Scheme == "" {
			return nil, &url.Error{Op: "parse", URL: s, Err: errMissingScheme}
		}
		return u, nil
	}

	p := From(MapLookuper{"U": "https://example.com/path"})
	got := RequiredValidated(p, "U", parseURL)
	if err := p.Err(); err != nil {
		t.Fatalf("Err = %v, want nil", err)
	}
	if got.Host != "example.com" {
		t.Errorf("host = %q, want example.com", got.Host)
	}

	// Validation failure (no scheme) on required path records an error.
	p = From(MapLookuper{"U": "example.com"})
	RequiredValidated(p, "U", parseURL)
	if p.Err() == nil {
		t.Error("RequiredValidated: Err = nil, want validation error")
	}

	// Defaulted path returns the supplied default on failure.
	def := &url.URL{Scheme: "https", Host: "default.local"}
	p = From(MapLookuper{"U": "::not-a-url"})
	if got := Validated(p, "U", def, parseURL); got != def {
		t.Errorf("Validated = %v, want default", got)
	}
	if err := p.Err(); err != nil {
		t.Errorf("Err = %v, want nil", err)
	}
}

var errMissingScheme = errStr("missing scheme")

type errStr string

func (e errStr) Error() string { return string(e) }
