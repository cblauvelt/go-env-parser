package env

import "strings"

// DefaultSeparator is the delimiter used by CSV and IntSlice when no explicit
// separator is supplied.
const DefaultSeparator = ","

// CSV returns the value of key split on DefaultSeparator, with each element
// trimmed of surrounding whitespace and empty elements dropped. It returns def
// if key is unset. If key is set but yields no non-empty elements, an empty
// (non-nil) slice is returned rather than def — a set value always overrides the
// default. It never records an error.
func (p *Parser) CSV(key string, def []string) []string {
	return p.CSVSep(key, DefaultSeparator, def)
}

// CSVSep behaves like CSV but splits on the supplied separator.
func (p *Parser) CSVSep(key, sep string, def []string) []string {
	v, ok := p.lookup(key)
	if !ok {
		return def
	}
	return splitTrim(v, sep)
}

// RequiredCSV returns the value of key split like CSV. If key is unset it
// records a missing error and returns nil. A key that is set but yields no
// non-empty elements is not an error and returns an empty slice.
func (p *Parser) RequiredCSV(key string) []string {
	return p.RequiredCSVSep(key, DefaultSeparator)
}

// RequiredCSVSep behaves like RequiredCSV but splits on the supplied separator.
func (p *Parser) RequiredCSVSep(key, sep string) []string {
	v, ok := p.lookup(key)
	if !ok {
		p.addErr(missingErr(key))
		return nil
	}
	return splitTrim(v, sep)
}

// IntSlice returns the value of key split on DefaultSeparator with each element
// parsed as a base-10 int. It returns def if key is unset. A malformed element
// fires the bad-default hook and causes def to be returned. It never records an
// error.
func (p *Parser) IntSlice(key string, def []int) []int {
	v, ok := p.lookup(key)
	if !ok {
		return def
	}
	out, err := parseIntSlice(v, DefaultSeparator)
	if err != nil {
		if p.onBadDefault != nil {
			p.onBadDefault(key, v, parseErr(key, v, "int slice", err))
		}
		return def
	}
	return out
}

// RequiredIntSlice returns the value of key split on DefaultSeparator with each
// element parsed as a base-10 int. If key is unset, or any element fails to
// parse, it records an error and returns nil.
func (p *Parser) RequiredIntSlice(key string) []int {
	v, ok := p.lookup(key)
	if !ok {
		p.addErr(missingErr(key))
		return nil
	}
	out, err := parseIntSlice(v, DefaultSeparator)
	if err != nil {
		p.addErr(parseErr(key, v, "int slice", err))
		return nil
	}
	return out
}

// splitTrim splits s on sep, trims whitespace from each element, and drops empty
// elements. It always returns a non-nil slice when s is non-empty after
// splitting; an input that yields no elements returns an empty, non-nil slice.
func splitTrim(s, sep string) []string {
	out := []string{}
	for _, part := range strings.Split(s, sep) {
		if v := strings.TrimSpace(part); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// parseIntSlice splits s with splitTrim and parses every element as a base-10
// int, returning the first parse error encountered.
func parseIntSlice(s, sep string) ([]int, error) {
	parts := splitTrim(s, sep)
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		n, err := parseInt(part)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}
