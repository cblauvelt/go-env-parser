package env

import (
	"slices"
	"strings"
	"testing"
)

func TestRequiredCSVSep(t *testing.T) {
	p := From(MapLookuper{"K": "a:b:c"})
	if got := p.RequiredCSVSep("K", ":"); !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Errorf("RequiredCSVSep = %#v, want [a b c]", got)
	}
	if err := p.Err(); err != nil {
		t.Errorf("Err = %v, want nil", err)
	}
}

func TestIntSlice_Absent(t *testing.T) {
	def := []int{1, 2, 3}
	p := From(MapLookuper{})
	if got := p.IntSlice("K", def); !slices.Equal(got, def) {
		t.Errorf("IntSlice absent = %#v, want default", got)
	}
	if err := p.Err(); err != nil {
		t.Errorf("Err = %v, want nil", err)
	}
}

func TestRequiredIntSlice_Present(t *testing.T) {
	p := From(MapLookuper{"K": "1, 2 ,3"})
	if got := p.RequiredIntSlice("K"); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("RequiredIntSlice = %#v, want [1 2 3]", got)
	}
	if err := p.Err(); err != nil {
		t.Errorf("Err = %v, want nil", err)
	}
}

func TestRequiredIntSlice_Absent(t *testing.T) {
	p := From(MapLookuper{})
	if got := p.RequiredIntSlice("K"); got != nil {
		t.Errorf("RequiredIntSlice = %#v, want nil", got)
	}
	if err := p.Err(); err == nil || !strings.Contains(err.Error(), "required") {
		t.Errorf("Err = %v, want missing error", err)
	}
}

func TestRequiredIntSlice_Malformed(t *testing.T) {
	p := From(MapLookuper{"K": "1,x,3"})
	if got := p.RequiredIntSlice("K"); got != nil {
		t.Errorf("RequiredIntSlice malformed = %#v, want nil", got)
	}
	if err := p.Err(); err == nil {
		t.Error("Err = nil, want parse error")
	}
}

func TestRequiredCSV_SetButEmpty(t *testing.T) {
	p := From(MapLookuper{"K": " , , "})
	got := p.RequiredCSV("K")
	// Set key with no non-empty elements is not an error; returns empty slice.
	if got == nil {
		t.Error("RequiredCSV = nil, want empty (non-nil) slice")
	}
	if len(got) != 0 {
		t.Errorf("RequiredCSV = %#v, want empty slice", got)
	}
	if err := p.Err(); err != nil {
		t.Errorf("Err = %v, want nil", err)
	}
}
