package env

import (
	"strconv"
	"time"
)

// String returns the value of key, or def if key is unset (or empty when
// WithEmptyAsUnset is in effect). It never records an error.
func (p *Parser) String(key, def string) string {
	if v, ok := p.lookup(key); ok {
		return v
	}
	return def
}

// RequiredString returns the value of key. If key is unset it records a missing
// error and returns "".
func (p *Parser) RequiredString(key string) string {
	if v, ok := p.lookup(key); ok {
		return v
	}
	p.addErr(missingErr(key))
	return ""
}

// Int returns the value of key parsed as a base-10 int, or def if key is unset
// or malformed. A malformed value invokes the bad-default hook (if set) and is
// otherwise silently replaced by def.
func (p *Parser) Int(key string, def int) int {
	return parseDefault(p, key, def, "int", parseInt)
}

// RequiredInt returns the value of key parsed as a base-10 int. If key is unset
// or malformed it records an error and returns 0.
func (p *Parser) RequiredInt(key string) int {
	return parseRequired(p, key, "int", parseInt)
}

// parseInt is the base-10 int parse used by the int accessors.
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// Int32 returns the value of key parsed as a base-10 int32, or def if key is
// unset or malformed (including overflow). A malformed value fires the
// bad-default hook.
func (p *Parser) Int32(key string, def int32) int32 {
	return parseDefault(p, key, def, "int32", parseInt32)
}

// RequiredInt32 returns the value of key parsed as a base-10 int32. If key is
// unset or malformed (including overflow) it records an error and returns 0.
func (p *Parser) RequiredInt32(key string) int32 {
	return parseRequired(p, key, "int32", parseInt32)
}

func parseInt32(s string) (int32, error) {
	n, err := strconv.ParseInt(s, 10, 32)
	return int32(n), err
}

// Int64 returns the value of key parsed as a base-10 int64, or def if key is
// unset or malformed. A malformed value fires the bad-default hook.
func (p *Parser) Int64(key string, def int64) int64 {
	return parseDefault(p, key, def, "int64", parseInt64)
}

// RequiredInt64 returns the value of key parsed as a base-10 int64. If key is
// unset or malformed it records an error and returns 0.
func (p *Parser) RequiredInt64(key string) int64 {
	return parseRequired(p, key, "int64", parseInt64)
}

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// Uint returns the value of key parsed as a base-10 unsigned int, or def if key
// is unset or malformed (including negative or overflow). A malformed value
// fires the bad-default hook.
func (p *Parser) Uint(key string, def uint) uint {
	return parseDefault(p, key, def, "uint", parseUint)
}

// RequiredUint returns the value of key parsed as a base-10 unsigned int. If key
// is unset or malformed it records an error and returns 0.
func (p *Parser) RequiredUint(key string) uint {
	return parseRequired(p, key, "uint", parseUint)
}

func parseUint(s string) (uint, error) {
	n, err := strconv.ParseUint(s, 10, strconv.IntSize)
	return uint(n), err
}

// Uint64 returns the value of key parsed as a base-10 uint64, or def if key is
// unset or malformed. A malformed value fires the bad-default hook.
func (p *Parser) Uint64(key string, def uint64) uint64 {
	return parseDefault(p, key, def, "uint64", parseUint64)
}

// RequiredUint64 returns the value of key parsed as a base-10 uint64. If key is
// unset or malformed it records an error and returns 0.
func (p *Parser) RequiredUint64(key string) uint64 {
	return parseRequired(p, key, "uint64", parseUint64)
}

func parseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// Float64 returns the value of key parsed as a float64, or def if key is unset
// or malformed. A malformed value fires the bad-default hook.
func (p *Parser) Float64(key string, def float64) float64 {
	return parseDefault(p, key, def, "float64", parseFloat64)
}

// RequiredFloat64 returns the value of key parsed as a float64. If key is unset
// or malformed it records an error and returns 0.
func (p *Parser) RequiredFloat64(key string) float64 {
	return parseRequired(p, key, "float64", parseFloat64)
}

func parseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// Bool returns the value of key parsed as a bool, or def if key is unset or
// malformed. A malformed value fires the bad-default hook. Accepted values are
// those of strconv.ParseBool: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false,
// False.
func (p *Parser) Bool(key string, def bool) bool {
	return parseDefault(p, key, def, "bool", strconv.ParseBool)
}

// RequiredBool returns the value of key parsed as a bool. If key is unset or
// malformed it records an error and returns false.
func (p *Parser) RequiredBool(key string) bool {
	return parseRequired(p, key, "bool", strconv.ParseBool)
}

// Duration returns the value of key parsed as a time.Duration (e.g. "300ms",
// "1.5h", "2h45m"), or def if key is unset or malformed. A malformed value
// fires the bad-default hook.
func (p *Parser) Duration(key string, def time.Duration) time.Duration {
	return parseDefault(p, key, def, "duration", time.ParseDuration)
}

// RequiredDuration returns the value of key parsed as a time.Duration. If key is
// unset or malformed it records an error and returns 0.
func (p *Parser) RequiredDuration(key string) time.Duration {
	return parseRequired(p, key, "duration", time.ParseDuration)
}

// Time returns the value of key parsed as a time.Time using RFC 3339, or def if
// key is unset or malformed. A malformed value fires the bad-default hook. Use
// TimeLayout for a non-RFC3339 layout.
func (p *Parser) Time(key string, def time.Time) time.Time {
	return parseDefault(p, key, def, "time", parseTime(time.RFC3339))
}

// RequiredTime returns the value of key parsed as a time.Time using RFC 3339. If
// key is unset or malformed it records an error and returns the zero time.
func (p *Parser) RequiredTime(key string) time.Time {
	return parseRequired(p, key, "time", parseTime(time.RFC3339))
}

// TimeLayout returns the value of key parsed as a time.Time using the supplied
// reference layout (see the time package), or def if key is unset or malformed.
func (p *Parser) TimeLayout(key, layout string, def time.Time) time.Time {
	return parseDefault(p, key, def, "time", parseTime(layout))
}

// RequiredTimeLayout returns the value of key parsed as a time.Time using the
// supplied reference layout. If key is unset or malformed it records an error
// and returns the zero time.
func (p *Parser) RequiredTimeLayout(key, layout string) time.Time {
	return parseRequired(p, key, "time", parseTime(layout))
}

func parseTime(layout string) func(string) (time.Time, error) {
	return func(s string) (time.Time, error) {
		return time.Parse(layout, s)
	}
}

// parseRequired is the shared implementation for Required* scalar accessors: it
// looks up key, parses it, and records a missing or parse error on failure.
func parseRequired[T any](p *Parser, key, typ string, parse func(string) (T, error)) T {
	var zero T
	v, ok := p.lookup(key)
	if !ok {
		p.addErr(missingErr(key))
		return zero
	}
	out, err := parse(v)
	if err != nil {
		p.addErr(parseErr(key, v, typ, err))
		return zero
	}
	return out
}

// parseDefault is the shared implementation for defaulted scalar accessors: it
// returns def when key is unset, and when key is present but malformed it
// notifies the bad-default hook and returns def. It never records an error.
func parseDefault[T any](p *Parser, key string, def T, typ string, parse func(string) (T, error)) T {
	v, ok := p.lookup(key)
	if !ok {
		return def
	}
	out, err := parse(v)
	if err != nil {
		if p.onBadDefault != nil {
			p.onBadDefault(key, v, parseErr(key, v, typ, err))
		}
		return def
	}
	return out
}
