package env

import (
	"errors"
	"time"
)

// errNotPositive is the cause reported when a value parses cleanly but is not
// greater than zero. parseErr wraps it with the key, raw value, and type name.
var errNotPositive = errors.New("must be greater than zero")

// positiveNumeric is the set of types the Positive accessors cover. time.Duration
// has underlying type int64, so ~int64 admits it.
type positiveNumeric interface {
	~int | ~int64 | ~float64
}

// positive adapts a parse function to reject values that are not greater than
// zero, so a non-positive value follows the same path as an unparseable one:
// recorded by the Required accessors, routed to the bad-default hook by the
// defaulted ones.
func positive[T positiveNumeric](parse func(string) (T, error)) func(string) (T, error) {
	return func(s string) (T, error) {
		v, err := parse(s)
		if err != nil {
			var zero T
			return zero, err
		}
		if v <= 0 {
			var zero T
			return zero, errNotPositive
		}
		return v, nil
	}
}

// PositiveInt returns the value of key parsed as a base-10 int greater than
// zero, or def if key is unset, malformed, or not positive. A rejected present
// value fires the bad-default hook.
func (p *Parser) PositiveInt(key string, def int) int {
	return parseDefault(p, key, def, "positive int", positive(parseInt))
}

// RequiredPositiveInt returns the value of key parsed as a base-10 int greater
// than zero. If key is unset, malformed, or not positive it records an error and
// returns 0.
func (p *Parser) RequiredPositiveInt(key string) int {
	return parseRequired(p, key, "positive int", positive(parseInt))
}

// PositiveInt64 returns the value of key parsed as a base-10 int64 greater than
// zero, or def if key is unset, malformed, or not positive. A rejected present
// value fires the bad-default hook.
func (p *Parser) PositiveInt64(key string, def int64) int64 {
	return parseDefault(p, key, def, "positive int64", positive(parseInt64))
}

// RequiredPositiveInt64 returns the value of key parsed as a base-10 int64
// greater than zero. If key is unset, malformed, or not positive it records an
// error and returns 0.
func (p *Parser) RequiredPositiveInt64(key string) int64 {
	return parseRequired(p, key, "positive int64", positive(parseInt64))
}

// PositiveFloat64 returns the value of key parsed as a float64 greater than
// zero, or def if key is unset, malformed, or not positive. A rejected present
// value fires the bad-default hook.
func (p *Parser) PositiveFloat64(key string, def float64) float64 {
	return parseDefault(p, key, def, "positive float64", positive(parseFloat64))
}

// RequiredPositiveFloat64 returns the value of key parsed as a float64 greater
// than zero. If key is unset, malformed, or not positive it records an error and
// returns 0.
func (p *Parser) RequiredPositiveFloat64(key string) float64 {
	return parseRequired(p, key, "positive float64", positive(parseFloat64))
}

// PositiveDuration returns the value of key parsed as a time.Duration greater
// than zero, or def if key is unset, malformed, or not positive. A rejected
// present value fires the bad-default hook.
func (p *Parser) PositiveDuration(key string, def time.Duration) time.Duration {
	return parseDefault(p, key, def, "positive duration", positive(time.ParseDuration))
}

// RequiredPositiveDuration returns the value of key parsed as a time.Duration
// greater than zero. If key is unset, malformed, or not positive it records an
// error and returns 0.
func (p *Parser) RequiredPositiveDuration(key string) time.Duration {
	return parseRequired(p, key, "positive duration", positive(time.ParseDuration))
}
