package env

import (
	"fmt"
	"slices"
)

// OneOf returns the value of key if it is one of allowed, otherwise def. It
// returns def when key is unset. When key is set to a value outside allowed, the
// bad-default hook fires and def is returned. It never records an error.
//
// def is not required to be a member of allowed; callers are responsible for
// supplying a sensible default.
func (p *Parser) OneOf(key, def string, allowed ...string) string {
	v, ok := p.lookup(key)
	if !ok {
		return def
	}
	if slices.Contains(allowed, v) {
		return v
	}
	if p.onBadDefault != nil {
		p.onBadDefault(key, v, oneOfErr(key, v, allowed))
	}
	return def
}

// RequiredOneOf returns the value of key if it is one of allowed. If key is
// unset it records a missing error; if key is set to a value outside allowed it
// records a validation error. On failure it returns "".
func (p *Parser) RequiredOneOf(key string, allowed ...string) string {
	v, ok := p.lookup(key)
	if !ok {
		p.addErr(missingErr(key))
		return ""
	}
	if slices.Contains(allowed, v) {
		return v
	}
	p.addErr(oneOfErr(key, v, allowed))
	return ""
}

func oneOfErr(key, raw string, allowed []string) error {
	return fmt.Errorf("env %q: value %q is not one of %v", key, raw, allowed)
}

// Validated returns the value of key converted by parse, or def if key is unset
// or parse returns an error. A parse error on a present value fires the
// bad-default hook. It never records an error.
//
// parse may both convert and validate: returning a non-nil error rejects the
// value. This is the generic escape hatch for types and rules the built-in
// accessors do not cover (e.g. *url.URL, net.IP, a bounded int).
//
// It is a package-level function rather than a method because Go methods cannot
// declare their own type parameters.
func Validated[T any](p *Parser, key string, def T, parse func(string) (T, error)) T {
	return parseDefault(p, key, def, "value", parse)
}

// RequiredValidated returns the value of key converted by parse. If key is unset
// or parse returns an error, it records an error and returns the zero value of
// T.
func RequiredValidated[T any](p *Parser, key string, parse func(string) (T, error)) T {
	return parseRequired(p, key, "value", parse)
}
