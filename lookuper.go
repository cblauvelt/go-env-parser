// Package env provides explicit, batch-validated parsing of environment
// variables into typed Go values.
//
// Unlike struct-tag binding libraries, every variable is read with an explicit
// function call, keeping the mapping between env vars and fields easy to read
// and grep. A Parser accumulates errors as it goes so a single Err() call
// reports every missing or malformed variable in one pass, rather than failing
// on the first one.
//
// Values are read through a Lookuper, which decouples parsing from the process
// environment. The default source is the OS environment; tests and layered
// configuration can inject a MapLookuper instead, avoiding global state such as
// os.Setenv / t.Setenv.
package env

import "os"

// Lookuper resolves an environment variable key to its raw string value and
// reports whether the key was present at all. Distinguishing "present but
// empty" from "absent" is intentional: some callers treat KEY="" as a
// meaningful explicit value, others as unset (see WithEmptyAsUnset).
type Lookuper interface {
	// Lookup returns the value for key and whether key was set. A set key with
	// an empty value returns ("", true).
	Lookup(key string) (value string, ok bool)
}

// OSLookuper reads from the process environment via os.LookupEnv. It is the
// default source used by New.
type OSLookuper struct{}

// Lookup implements Lookuper against the OS environment.
func (OSLookuper) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

// MapLookuper reads from an in-memory map. A nil or missing entry is reported as
// unset. It is intended for tests and for composing layered configuration
// without touching the real environment.
type MapLookuper map[string]string

// Lookup implements Lookuper against the backing map.
func (m MapLookuper) Lookup(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

// Compile-time assertions that the built-in sources satisfy Lookuper.
var (
	_ Lookuper = OSLookuper{}
	_ Lookuper = MapLookuper(nil)
)
