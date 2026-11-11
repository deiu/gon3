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
