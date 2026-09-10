// Package files parses ttybus `files` payloads.
//
// The easy form is an absolute path. Optional `:line` or `:line:col`.
// When context grows, publishers may send a single-line JSON object.
package files

import (
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"
)

// Event is a selected file (and optional cursor) on the bus.
type Event struct {
	Path   string `json:"path"`
	Line   int    `json:"line,omitempty"`
	Col    int    `json:"col,omitempty"`
	Root   string `json:"root,omitempty"`
	Origin string `json:"origin,omitempty"`
	Raw    string `json:"-"`
}

// Parse a `files` payload. Empty input is not an error; Path stays empty.
func Parse(payload string) (Event, error) {
	payload = strings.TrimSpace(payload)
	ev := Event{Raw: payload}
	if payload == "" {
		return ev, nil
	}
	if strings.HasPrefix(payload, "{") {
		if err := json.Unmarshal([]byte(payload), &ev); err != nil {
			return Event{Raw: payload}, err
		}
		ev.Raw = payload
		ev.Path = cleanPath(ev.Path)
		return ev, nil
	}
	path, line, col := splitLocation(payload)
	ev.Path = cleanPath(path)
	ev.Line = line
	ev.Col = col
	return ev, nil
}

func cleanPath(p string) string {
	p = strings.TrimSpace(p)
	p = strings.TrimPrefix(p, "file://")
	if p == "" {
		return ""
	}
	return filepath.Clean(p)
}

// splitLocation pulls :line[:col] off a path only when the prefix exists
// as a path-shaped string. We do not stat here — publishers may point at
// a file that was just deleted.
func splitLocation(s string) (path string, line, col int) {
	s = strings.TrimPrefix(s, "file://")
	path = s
	// Walk from the right: optional col, then line.
	a, aok := cutLastNumber(s)
	if !aok {
		return s, 0, 0
	}
	b, bok := cutLastNumber(a)
	if bok {
		// path:line:col
		col = mustInt(s[len(a)+1:])
		line = mustInt(a[len(b)+1:])
		return b, line, col
	}
	// path:line
	line = mustInt(s[len(a)+1:])
	return a, line, 0
}

func cutLastNumber(s string) (prefix string, ok bool) {
	i := strings.LastIndexByte(s, ':')
	if i <= 0 {
		return s, false
	}
	n := s[i+1:]
	if n == "" {
		return s, false
	}
	for _, r := range n {
		if r < '0' || r > '9' {
			return s, false
		}
	}
	return s[:i], true
}

func mustInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
