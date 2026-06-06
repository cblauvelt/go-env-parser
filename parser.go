package env

import (
	"errors"
	"fmt"
)

// Parser reads environment variables through a Lookuper and accumulates any
// errors encountered. Required accessors that fail append an error; defaulted
// accessors never append an error but report malformed values through the
// onBadDefault hook before falling back to the supplied default.
//
// The zero value is not usable; construct a Parser with New or From.
//
// A Parser is not safe for concurrent use. Typical usage builds a config struct
// from a single goroutine at startup and then checks Err once.
type Parser struct {
	src          Lookuper
	emptyAsUnset bool
	onBadDefault func(key, raw string, err error)
	errs         []error
}

// Option configures a Parser.
type Option func(*Parser)

// WithEmptyAsUnset makes a set-but-empty variable (KEY="") be treated as if the
// variable were absent: required accessors report it missing, and defaulted
// accessors return their default. This restores the semantics of callers that
// consider an empty value meaningless.
func WithEmptyAsUnset() Option {
	return func(p *Parser) { p.emptyAsUnset = true }
}

// WithBadDefaultHook registers a callback invoked when a defaulted accessor
// finds a present value that fails to parse. The accessor still returns its
// default; the hook lets callers log the typo instead of silently swallowing
// it. fn is called with the key, the raw value, and the parse error.
func WithBadDefaultHook(fn func(key, raw string, err error)) Option {
	return func(p *Parser) { p.onBadDefault = fn }
}

// New returns a Parser that reads from the OS environment.
func New(opts ...Option) *Parser {
	return From(OSLookuper{}, opts...)
}

// From returns a Parser that reads from src. A nil src defaults to OSLookuper.
func From(src Lookuper, opts ...Option) *Parser {
	if src == nil {
		src = OSLookuper{}
	}
	p := &Parser{src: src}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Err returns the accumulated errors joined into a single error, or nil if no
// required accessor has failed. Call it once after reading all variables.
func (p *Parser) Err() error {
	return errors.Join(p.errs...)
}

// lookup resolves key honoring the emptyAsUnset option. It returns the value and
// whether it should be considered present.
func (p *Parser) lookup(key string) (string, bool) {
	v, ok := p.src.Lookup(key)
	if !ok {
		return "", false
	}
	if p.emptyAsUnset && v == "" {
		return "", false
	}
	return v, true
}

// addErr appends err to the accumulated set.
func (p *Parser) addErr(err error) {
	p.errs = append(p.errs, err)
}

// missingErr builds the canonical "required and unset" error for key.
func missingErr(key string) error {
	return fmt.Errorf("env %q is required but not set", key)
}

// parseErr builds the canonical "set but unparseable" error for key.
func parseErr(key, raw, typ string, cause error) error {
	return fmt.Errorf("env %q: invalid %s value %q: %w", key, typ, raw, cause)
}
