package gon3

import (
	"strings"
	"testing"
)

// TestLocalNameConsecutiveEscapes checks that the lexer accepts a local name
// with more than one PN_LOCAL_ESC or PLX separated by name characters.
func TestLocalNameConsecutiveEscapes(t *testing.T) {
	for _, local := range []string{`a\/b\/c`, `a\?b\#c`, `a%20b\/c`, `a\/b%20c`} {
		doc := "@prefix ex: <http://ex/> .\nex:" + local + " <http://ex/p> <http://ex/o> .\n"
		if _, err := NewParser("http://ex/").Parse(strings.NewReader(doc)); err != nil {
			t.Errorf("Parse(%q) error: %v", doc, err)
		}
	}
}

// TestLocalNameEscapes checks that a prefixed name's IRI is the namespace IRI
// followed by the local name with each PN_LOCAL_ESC backslash removed, as
// Turtle 1.1 section 6.3 requires, and that %-escapes are kept as written.
func TestLocalNameEscapes(t *testing.T) {
	tests := []struct {
		name, local, want string
	}{
		{"escaped slash", `a\/b`, "http://ex/a/b"},
		{"consecutive escaped slashes", `a\/b\/c`, "http://ex/a/b/c"},
		{"adjacent escapes", `a\/\/b`, "http://ex/a//b"},
		{"escape first", `\/a`, "http://ex//a"},
		{"escaped hyphen", `a\-b`, "http://ex/a-b"},
		{"escape then percent", `a\/%20b`, "http://ex/a/%20b"},
		{"percent then escape", `a%20\/b`, "http://ex/a%20/b"},
		{"escaped query and fragment", `a\?b\#c`, "http://ex/a?b#c"},
		{"colon in local name", `a:b`, "http://ex/a:b"},
		{"trailing colon", `s:`, "http://ex/s:"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := "@prefix ex: <http://ex/> .\nex:" + tt.local + " <http://ex/p> <http://ex/o> .\n"
			g, err := NewParser("http://ex/").Parse(strings.NewReader(doc))
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", doc, err)
			}
			var got []string
			for tr := range g.IterTriples() {
				got = append(got, tr.Subject.(*IRI).url)
			}
			if len(got) != 1 || got[0] != tt.want {
				t.Errorf("subject IRIs = %q, want [%q]", got, tt.want)
			}
		})
	}
}

// TestLocalNameAppendsToNamespace checks that the local name is appended to
// the namespace IRI as a string rather than resolved against it as a
// relative reference.
func TestLocalNameAppendsToNamespace(t *testing.T) {
	tests := []struct {
		name, ns, local, want string
	}{
		{"hash namespace", "http://ex/vocab#", "Thing", "http://ex/vocab#Thing"},
		{"no trailing separator", "http://ex/vocab", "Thing", "http://ex/vocabThing"},
		{"escaped slash after hash", "http://ex/vocab#", `a\/b`, "http://ex/vocab#a/b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := "@prefix ex: <" + tt.ns + "> .\nex:" + tt.local + " <http://ex/p> <http://ex/o> .\n"
			g, err := NewParser("http://ex/").Parse(strings.NewReader(doc))
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", doc, err)
			}
			for tr := range g.IterTriples() {
				if got := tr.Subject.(*IRI).url; got != tt.want {
					t.Errorf("subject = %q, want %q", got, tt.want)
				}
			}
		})
	}
}
