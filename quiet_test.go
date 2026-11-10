package gon3

import (
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout runs f and returns everything it wrote to os.Stdout.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()
	f()
	w.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func mustParse(t *testing.T, doc string) *Graph {
	t.Helper()
	g, err := NewParser("http://example.org/").Parse(strings.NewReader(doc))
	if err != nil {
		t.Fatalf("parse %q: %v", doc, err)
	}
	return g
}

func TestIsomorphismWritesNothingToStdout(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{
			name: "isomorphic blank nodes",
			a:    `<http://ex/s> <http://ex/p> [ <http://ex/q> "x" ] .`,
			b:    `<http://ex/s> <http://ex/p> _:b1 . _:b1 <http://ex/q> "x" .`,
			want: true,
		},
		{
			name: "different sizes",
			a:    `<http://ex/s> <http://ex/p> "x" .`,
			b:    `<http://ex/s> <http://ex/p> "x", "y" .`,
			want: false,
		},
		{
			name: "same size, different IRIs",
			a:    `<http://ex/s> <http://ex/p> <http://ex/o1> .`,
			b:    `<http://ex/s> <http://ex/p> <http://ex/o2> .`,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ga, gb := mustParse(t, tt.a), mustParse(t, tt.b)
			var got bool
			out := captureStdout(t, func() {
				got = ga.IsomorphicTo(gb)
				_ = ga.NodesSorted()
			})
			if got != tt.want {
				t.Errorf("IsomorphicTo = %v, want %v", got, tt.want)
			}
			if out != "" {
				t.Errorf("library wrote to stdout:\n%s", out)
			}
		})
	}
}
