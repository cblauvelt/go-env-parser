package env_test

import (
	"strings"
	"testing"

	env "github.com/cblauvelt/go-env-parser"
)

func TestAllOrNone(t *testing.T) {
	tests := []struct {
		name    string
		vars    env.MapLookuper
		keys    []string
		wantErr bool
		// substrings that must appear in the error, if any.
		contains []string
	}{
		{
			name:    "all set",
			vars:    env.MapLookuper{"A": "1", "B": "2"},
			keys:    []string{"A", "B"},
			wantErr: false,
		},
		{
			name:    "all unset",
			vars:    env.MapLookuper{},
			keys:    []string{"A", "B"},
			wantErr: false,
		},
		{
			name:     "one of two set",
			vars:     env.MapLookuper{"A": "1"},
			keys:     []string{"A", "B"},
			wantErr:  true,
			contains: []string{"A", "B", "set: [A]"},
		},
		{
			name:    "no keys",
			vars:    env.MapLookuper{},
			keys:    nil,
			wantErr: false,
		},
		{
			name:    "single key set",
			vars:    env.MapLookuper{"A": "1"},
			keys:    []string{"A"},
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := env.From(tc.vars)
			p.AllOrNone(tc.keys...)
			err := p.Err()
			if tc.wantErr && err == nil {
				t.Fatal("Err = nil, want error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("Err = %v, want nil", err)
			}
			for _, want := range tc.contains {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("Err = %q, want it to contain %q", err, want)
				}
			}
		})
	}
}

func TestAllOrNone_EmptyAsUnset(t *testing.T) {
	// B is set-but-empty, so with WithEmptyAsUnset only A counts as set.
	p := env.From(env.MapLookuper{"A": "1", "B": ""}, env.WithEmptyAsUnset())
	p.AllOrNone("A", "B")
	if err := p.Err(); err == nil {
		t.Fatal("Err = nil, want error because B is empty and treated as unset")
	}
}

func TestMutuallyExclusive(t *testing.T) {
	tests := []struct {
		name    string
		vars    env.MapLookuper
		keys    []string
		wantErr bool
	}{
		{
			name:    "none set",
			vars:    env.MapLookuper{},
			keys:    []string{"A", "B"},
			wantErr: false,
		},
		{
			name:    "one set",
			vars:    env.MapLookuper{"A": "1"},
			keys:    []string{"A", "B"},
			wantErr: false,
		},
		{
			name:    "two set",
			vars:    env.MapLookuper{"A": "1", "B": "2"},
			keys:    []string{"A", "B"},
			wantErr: true,
		},
		{
			name:    "no keys",
			vars:    env.MapLookuper{},
			keys:    nil,
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := env.From(tc.vars)
			p.MutuallyExclusive(tc.keys...)
			err := p.Err()
			if tc.wantErr && err == nil {
				t.Fatal("Err = nil, want error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("Err = %v, want nil", err)
			}
		})
	}
}

func TestMutuallyExclusive_NamesConflictingKeys(t *testing.T) {
	p := env.From(env.MapLookuper{"A": "x", "B": "y"})
	p.MutuallyExclusive("A", "B")
	err := p.Err()
	if err == nil {
		t.Fatal("Err = nil, want error")
	}
	for _, want := range []string{"A", "B"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Err = %q, want it to contain %q", err, want)
		}
	}
}

func TestMutuallyExclusive_EmptyAsUnset(t *testing.T) {
	// Only A is truly set; B is empty and treated as unset, so no conflict.
	p := env.From(env.MapLookuper{"A": "1", "B": ""}, env.WithEmptyAsUnset())
	p.MutuallyExclusive("A", "B")
	if err := p.Err(); err != nil {
		t.Fatalf("Err = %v, want nil", err)
	}
}

func TestRequiredWith(t *testing.T) {
	tests := []struct {
		name       string
		vars       env.MapLookuper
		key        string
		dependents []string
		wantErrs   int
	}{
		{
			name:       "key unset skips dependents",
			vars:       env.MapLookuper{},
			key:        "FEATURE",
			dependents: []string{"A", "B"},
			wantErrs:   0,
		},
		{
			name:       "key set all dependents present",
			vars:       env.MapLookuper{"FEATURE": "on", "A": "1", "B": "2"},
			key:        "FEATURE",
			dependents: []string{"A", "B"},
			wantErrs:   0,
		},
		{
			name:       "key set two dependents missing",
			vars:       env.MapLookuper{"FEATURE": "on"},
			key:        "FEATURE",
			dependents: []string{"A", "B"},
			wantErrs:   2,
		},
		{
			name:       "key set one dependent missing",
			vars:       env.MapLookuper{"FEATURE": "on", "A": "1"},
			key:        "FEATURE",
			dependents: []string{"A", "B"},
			wantErrs:   1,
		},
		{
			name:       "no dependents",
			vars:       env.MapLookuper{"FEATURE": "on"},
			key:        "FEATURE",
			dependents: nil,
			wantErrs:   0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := env.From(tc.vars)
			p.RequiredWith(tc.key, tc.dependents...)
			err := p.Err()
			got := 0
			if err != nil {
				got = strings.Count(err.Error(), "\n") + 1
			}
			if got != tc.wantErrs {
				t.Fatalf("got %d errors (%v), want %d", got, err, tc.wantErrs)
			}
		})
	}
}

func TestRequiredWith_EmptyAsUnset(t *testing.T) {
	// FEATURE set, dependent A is empty and treated as unset -> one error.
	p := env.From(env.MapLookuper{"FEATURE": "on", "A": ""}, env.WithEmptyAsUnset())
	p.RequiredWith("FEATURE", "A")
	err := p.Err()
	if err == nil {
		t.Fatal("Err = nil, want error because A is empty and treated as unset")
	}
	if !strings.Contains(err.Error(), `"A"`) {
		t.Errorf("Err = %q, want it to name A", err)
	}
}

// TestConstraints_JoinInCallOrder verifies constraint errors accumulate into
// Err() alongside accessor errors, in call order.
func TestConstraints_JoinInCallOrder(t *testing.T) {
	p := env.From(env.MapLookuper{"A": "1"})
	p.RequiredString("HOST") // accessor error first
	p.AllOrNone("A", "B")    // constraint error second
	err := p.Err()
	if err == nil {
		t.Fatal("Err = nil, want two errors")
	}
	msg := err.Error()
	hostIdx := strings.Index(msg, "HOST")
	constraintIdx := strings.Index(msg, "must be all set or all unset")
	if hostIdx == -1 || constraintIdx == -1 {
		t.Fatalf("Err = %q, want both accessor and constraint errors", msg)
	}
	if hostIdx > constraintIdx {
		t.Errorf("errors out of call order: %q", msg)
	}
}
