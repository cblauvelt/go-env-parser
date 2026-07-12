package env

import "fmt"

// AllOrNone records one error unless the named keys are all set or all unset. It
// is the constraint for values that only make sense together, such as a TLS
// certificate and its private key. The error names every key and lists which of
// them are set. It respects WithEmptyAsUnset: a set-but-empty variable counts as
// unset. Zero or one key always satisfies the constraint and records nothing.
func (p *Parser) AllOrNone(keys ...string) {
	set := p.setKeys(keys)
	if len(set) == 0 || len(set) == len(keys) {
		return
	}
	p.addErr(fmt.Errorf("env %v must be all set or all unset; set: %v", keys, set))
}

// MutuallyExclusive records one error when more than one of the named keys is
// set. It is the constraint for values that name competing sources of the same
// thing, such as two alternative binding stores. The error names every key and
// lists which of them are set. It respects WithEmptyAsUnset: a set-but-empty
// variable counts as unset. Zero or one key never conflicts and records nothing.
func (p *Parser) MutuallyExclusive(keys ...string) {
	set := p.setKeys(keys)
	if len(set) <= 1 {
		return
	}
	p.addErr(fmt.Errorf("env %v are mutually exclusive but multiple are set: %v", keys, set))
}

// RequiredWith records one error per unset dependent when key is set. It is the
// constraint for values that become required once a feature key turns the
// feature on. When key is unset the dependents are not checked and nothing is
// recorded. It respects WithEmptyAsUnset: a set-but-empty key or dependent
// counts as unset.
func (p *Parser) RequiredWith(key string, dependents ...string) {
	if _, ok := p.lookup(key); !ok {
		return
	}
	for _, dep := range dependents {
		if _, ok := p.lookup(dep); !ok {
			p.addErr(fmt.Errorf("env %q is set, so %q is required but not set", key, dep))
		}
	}
}

// setKeys returns the subset of keys that are present, honoring WithEmptyAsUnset
// through lookup. Order follows keys so error messages are deterministic.
func (p *Parser) setKeys(keys []string) []string {
	set := make([]string, 0, len(keys))
	for _, k := range keys {
		if _, ok := p.lookup(k); ok {
			set = append(set, k)
		}
	}
	return set
}
