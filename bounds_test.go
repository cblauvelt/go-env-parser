package env

import (
	"strings"
	"testing"
	"time"
)

func TestPositiveInt(t *testing.T) {
	tests := []struct {
		name string
		src  MapLookuper
		def  int
		want int
	}{
		{name: "positive", src: MapLookuper{"K": "7"}, def: 3, want: 7},
		{name: "absent uses default", src: MapLookuper{}, def: 3, want: 3},
		{name: "zero uses default", src: MapLookuper{"K": "0"}, def: 3, want: 3},
		{name: "negative uses default", src: MapLookuper{"K": "-1"}, def: 3, want: 3},
		{name: "malformed uses default", src: MapLookuper{"K": "many"}, def: 3, want: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := From(tt.src)
			if got := p.PositiveInt("K", tt.def); got != tt.want {
				t.Errorf("PositiveInt = %d, want %d", got, tt.want)
			}
			if err := p.Err(); err != nil {
				t.Errorf("Err = %v, want nil", err)
			}
		})
	}
}

func TestRequiredPositiveInt(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		p := From(MapLookuper{"K": "7"})
		if got := p.RequiredPositiveInt("K"); got != 7 {
			t.Errorf("RequiredPositiveInt = %d, want 7", got)
		}
		if err := p.Err(); err != nil {
			t.Errorf("Err = %v, want nil", err)
		}
	})
	t.Run("zero records error", func(t *testing.T) {
		p := From(MapLookuper{"K": "0"})
		if got := p.RequiredPositiveInt("K"); got != 0 {
			t.Errorf("RequiredPositiveInt = %d, want 0", got)
		}
		err := p.Err()
		if err == nil {
			t.Fatal("Err = nil, want a parse error")
		}
		if !strings.Contains(err.Error(), "greater than zero") {
			t.Errorf("Err = %v, want it to say the value must be greater than zero", err)
		}
	})
	t.Run("absent records error", func(t *testing.T) {
		p := From(MapLookuper{})
		if got := p.RequiredPositiveInt("K"); got != 0 {
			t.Errorf("RequiredPositiveInt = %d, want 0", got)
		}
		if p.Err() == nil {
			t.Error("Err = nil, want a missing error")
		}
	})
}

func TestPositiveInt64(t *testing.T) {
	tests := []struct {
		name string
		src  MapLookuper
		def  int64
		want int64
	}{
		{name: "positive", src: MapLookuper{"K": "65536"}, def: 1024, want: 65536},
		{name: "absent uses default", src: MapLookuper{}, def: 1024, want: 1024},
		{name: "zero uses default", src: MapLookuper{"K": "0"}, def: 1024, want: 1024},
		{name: "negative uses default", src: MapLookuper{"K": "-9"}, def: 1024, want: 1024},
		{name: "overflow uses default", src: MapLookuper{"K": "99999999999999999999"}, def: 1024, want: 1024},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := From(tt.src)
			if got := p.PositiveInt64("K", tt.def); got != tt.want {
				t.Errorf("PositiveInt64 = %d, want %d", got, tt.want)
			}
			if err := p.Err(); err != nil {
				t.Errorf("Err = %v, want nil", err)
			}
		})
	}
}

func TestRequiredPositiveInt64(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		p := From(MapLookuper{"K": "65536"})
		if got := p.RequiredPositiveInt64("K"); got != 65536 {
			t.Errorf("RequiredPositiveInt64 = %d, want 65536", got)
		}
		if err := p.Err(); err != nil {
			t.Errorf("Err = %v, want nil", err)
		}
	})
	t.Run("negative records error", func(t *testing.T) {
		p := From(MapLookuper{"K": "-1"})
		if got := p.RequiredPositiveInt64("K"); got != 0 {
			t.Errorf("RequiredPositiveInt64 = %d, want 0", got)
		}
		if p.Err() == nil {
			t.Error("Err = nil, want a parse error")
		}
	})
}

func TestPositiveFloat64(t *testing.T) {
	tests := []struct {
		name string
		src  MapLookuper
		def  float64
		want float64
	}{
		{name: "positive", src: MapLookuper{"K": "5.5"}, def: 1, want: 5.5},
		{name: "absent uses default", src: MapLookuper{}, def: 1, want: 1},
		{name: "zero uses default", src: MapLookuper{"K": "0"}, def: 1, want: 1},
		{name: "negative uses default", src: MapLookuper{"K": "-0.5"}, def: 1, want: 1},
		{name: "malformed uses default", src: MapLookuper{"K": "fast"}, def: 1, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := From(tt.src)
			if got := p.PositiveFloat64("K", tt.def); got != tt.want {
				t.Errorf("PositiveFloat64 = %v, want %v", got, tt.want)
			}
			if err := p.Err(); err != nil {
				t.Errorf("Err = %v, want nil", err)
			}
		})
	}
}

func TestRequiredPositiveFloat64(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		p := From(MapLookuper{"K": "5.5"})
		if got := p.RequiredPositiveFloat64("K"); got != 5.5 {
			t.Errorf("RequiredPositiveFloat64 = %v, want 5.5", got)
		}
		if err := p.Err(); err != nil {
			t.Errorf("Err = %v, want nil", err)
		}
	})
	t.Run("zero records error", func(t *testing.T) {
		p := From(MapLookuper{"K": "0"})
		if got := p.RequiredPositiveFloat64("K"); got != 0 {
			t.Errorf("RequiredPositiveFloat64 = %v, want 0", got)
		}
		if p.Err() == nil {
			t.Error("Err = nil, want a parse error")
		}
	})
}

func TestPositiveDuration(t *testing.T) {
	tests := []struct {
		name string
		src  MapLookuper
		def  time.Duration
		want time.Duration
	}{
		{name: "positive", src: MapLookuper{"K": "45m"}, def: time.Minute, want: 45 * time.Minute},
		{name: "absent uses default", src: MapLookuper{}, def: time.Minute, want: time.Minute},
		{name: "zero uses default", src: MapLookuper{"K": "0s"}, def: time.Minute, want: time.Minute},
		{name: "negative uses default", src: MapLookuper{"K": "-5m"}, def: time.Minute, want: time.Minute},
		{name: "malformed uses default", src: MapLookuper{"K": "soon"}, def: time.Minute, want: time.Minute},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := From(tt.src)
			if got := p.PositiveDuration("K", tt.def); got != tt.want {
				t.Errorf("PositiveDuration = %s, want %s", got, tt.want)
			}
			if err := p.Err(); err != nil {
				t.Errorf("Err = %v, want nil", err)
			}
		})
	}
}

func TestRequiredPositiveDuration(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		p := From(MapLookuper{"K": "15m"})
		if got := p.RequiredPositiveDuration("K"); got != 15*time.Minute {
			t.Errorf("RequiredPositiveDuration = %s, want 15m", got)
		}
		if err := p.Err(); err != nil {
			t.Errorf("Err = %v, want nil", err)
		}
	})
	t.Run("negative records error", func(t *testing.T) {
		p := From(MapLookuper{"K": "-5m"})
		if got := p.RequiredPositiveDuration("K"); got != 0 {
			t.Errorf("RequiredPositiveDuration = %s, want 0", got)
		}
		err := p.Err()
		if err == nil {
			t.Fatal("Err = nil, want a parse error")
		}
		if !strings.Contains(err.Error(), "positive duration") {
			t.Errorf("Err = %v, want it to name the expected type", err)
		}
	})
	t.Run("absent records error", func(t *testing.T) {
		p := From(MapLookuper{})
		if got := p.RequiredPositiveDuration("K"); got != 0 {
			t.Errorf("RequiredPositiveDuration = %s, want 0", got)
		}
		if p.Err() == nil {
			t.Error("Err = nil, want a missing error")
		}
	})
}

// A non-positive value is malformed, so the defaulted accessors route it to the
// bad-default hook exactly as they do an unparseable one.
func TestPositiveFiresBadDefaultHook(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		read func(*Parser)
	}{
		{name: "int", raw: "0", read: func(p *Parser) { p.PositiveInt("K", 3) }},
		{name: "int64", raw: "-1", read: func(p *Parser) { p.PositiveInt64("K", 3) }},
		{name: "float64", raw: "-0.5", read: func(p *Parser) { p.PositiveFloat64("K", 3) }},
		{name: "duration", raw: "0s", read: func(p *Parser) { p.PositiveDuration("K", time.Minute) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotKey, gotRaw string
			var gotErr error
			p := From(MapLookuper{"K": tt.raw}, WithBadDefaultHook(func(key, raw string, err error) {
				gotKey, gotRaw, gotErr = key, raw, err
			}))
			tt.read(p)

			if gotKey != "K" || gotRaw != tt.raw {
				t.Errorf("hook got (%q, %q), want (K, %q)", gotKey, gotRaw, tt.raw)
			}
			if gotErr == nil {
				t.Fatal("hook error = nil, want a parse error")
			}
			if !strings.Contains(gotErr.Error(), "greater than zero") {
				t.Errorf("hook error = %v, want it to say the value must be greater than zero", gotErr)
			}
			if err := p.Err(); err != nil {
				t.Errorf("Err = %v, want nil; a defaulted accessor never records", err)
			}
		})
	}
}
