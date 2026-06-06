package env_test

import (
	"os"
	"testing"

	env "github.com/cblauvelt/go-env-parser"
)

// TestNew_UsesOSEnvironment verifies that New reads from the actual process
// environment, and that OSLookuper.Lookup works correctly. Uses the real env
// rather than MapLookuper to exercise both code paths.
func TestNew_UsesOSEnvironment(t *testing.T) {
	const key = "GO_ENV_PARSER_TEST_KEY"
	os.Setenv(key, "hello")
	defer os.Unsetenv(key)

	p := env.New()
	if got := p.String(key, "default"); got != "hello" {
		t.Errorf("String = %q, want hello", got)
	}
}

func TestNew_AbsentKey(t *testing.T) {
	const key = "GO_ENV_PARSER_DEFINITELY_NOT_SET_12345"
	os.Unsetenv(key)

	p := env.New()
	if got := p.String(key, "fallback"); got != "fallback" {
		t.Errorf("String = %q, want fallback", got)
	}
}

// TestFrom_NilSrcDefaultsToOS verifies that From(nil) falls back to OSLookuper.
func TestFrom_NilSrcDefaultsToOS(t *testing.T) {
	const key = "GO_ENV_PARSER_NIL_SRC_TEST"
	os.Setenv(key, "present")
	defer os.Unsetenv(key)

	p := env.From(nil)
	if got := p.String(key, "default"); got != "present" {
		t.Errorf("String = %q, want present", got)
	}
}
