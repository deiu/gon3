package gon3

import "testing"

// Turtle does not require whitespace between most tokens. Each case
// parses a compact document and its whitespace-separated equivalent and
// checks that both produce the same graph.
func TestCompactSyntaxWithoutWhitespace(t *testing.T) {
	const pfx = "@prefix ex: <http://ex/> .\n"
	tests := []struct {
		name    string
		compact string
		spaced  string
	}{
		{
			name:    "a keyword followed by IRIREF",
			compact: pfx + `ex:s a<http://ex/C> .`,
			spaced:  pfx + `ex:s a <http://ex/C> .`,
		},
		{
			name:    "a keyword followed by blank node property list",
			compact: pfx + `ex:s a[ ex:p "x" ] .`,
			spaced:  pfx + `ex:s a [ ex:p "x" ] .`,
		},
		{
			name:    "a keyword followed by collection",
			compact: pfx + `ex:s a( ex:C ) .`,
			spaced:  pfx + `ex:s a ( ex:C ) .`,
		},
		{
			name:    "a keyword followed by string literal",
			compact: pfx + `ex:s a"x" .`,
			spaced:  pfx + `ex:s a "x" .`,
		},
		{
			name:    "SPARQL BASE followed by IRIREF",
			compact: "BASE<http://ex/base/>\n<s> <p> <o> .",
			spaced:  "BASE <http://ex/base/>\n<s> <p> <o> .",
		},
		{
			name:    "@prefix name followed by IRIREF",
			compact: "@prefix ex:<http://ex/> .\nex:s ex:p ex:o .",
			spaced:  "@prefix ex: <http://ex/> .\nex:s ex:p ex:o .",
		},
		{
			name:    "PNAME_NS subject followed by IRIREF predicate",
			compact: pfx + `ex:<http://ex/p> ex:o .`,
			spaced:  pfx + `ex: <http://ex/p> ex:o .`,
		},
		{
			name:    "PNAME_NS inside collection",
			compact: pfx + `ex:s ex:p (ex:) .`,
			spaced:  pfx + `ex:s ex:p ( ex: ) .`,
		},
		{
			name:    "PNAME_NS inside blank node property list",
			compact: pfx + `ex:s ex:p [ex:q ex:] .`,
			spaced:  pfx + `ex:s ex:p [ ex:q ex: ] .`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustParse(t, tt.compact)
			want := mustParse(t, tt.spaced)
			if !got.IsomorphicTo(want) {
				t.Errorf("graphs differ\ncompact:\n%s\nspaced:\n%s", got, want)
			}
		})
	}
}

// A leading 'a' that continues into a prefixed name must still lex as a
// prefixed name, not as the rdf:type keyword.
func TestPrefixedNamesStartingWithA(t *testing.T) {
	tests := []struct {
		name string
		doc  string
		pred string
	}{
		{
			name: "prefix beginning with a",
			doc:  "@prefix adms: <http://ex/adms#> .\n<http://ex/s> adms:status <http://ex/o> .",
			pred: "http://ex/adms#status",
		},
		{
			name: "single-letter prefix a",
			doc:  "@prefix a: <http://ex/a#> .\n<http://ex/s> a:concept <http://ex/o> .",
			pred: "http://ex/a#concept",
		},
		{
			name: "rdf:type keyword with whitespace",
			doc:  "<http://ex/s> a <http://ex/C> .",
			pred: "http://www.w3.org/1999/02/22-rdf-syntax-ns#type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := mustParse(t, tt.doc)
			var preds []string
			for tr := range g.IterTriples() {
				preds = append(preds, tr.Predicate.RawValue())
			}
			if len(preds) != 1 || preds[0] != tt.pred {
				t.Errorf("predicates = %v, want [%s]", preds, tt.pred)
			}
		})
	}
}
