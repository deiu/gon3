# gon3

gon3 is a [Turtle](https://www.w3.org/TR/turtle/) parser for Go. It reads a Turtle (or N-Triples) document and returns an RDF graph of triples. It can also check whether two graphs are isomorphic (the same graph, up to blank node labels).

## Install

```sh
go get github.com/deiu/gon3
```

You need Go 1.23 or later.

## Usage

```go
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/deiu/gon3"
)

func main() {
	doc := `
@prefix foaf: <http://xmlns.com/foaf/0.1/> .

<#me> a foaf:Person ;
    foaf:name "Alice"@en ;
    foaf:knows [ foaf:name "Bob" ] .
`
	// The base IRI resolves relative IRIs such as <#me>.
	p := gon3.NewParser("http://example.org/profile")
	graph, err := p.Parse(strings.NewReader(doc))
	if err != nil {
		log.Fatal(err)
	}

	for t := range graph.IterTriples() {
		if t.Predicate.RawValue() == "http://xmlns.com/foaf/0.1/name" {
			fmt.Println(t.Subject.RawValue(), "is named", t.Object.RawValue())
		}
	}
}
```

Output:

```
http://example.org/profile#me is named Alice
a0 is named Bob
```

A parser holds the state of one document. Create a new parser for each document that you parse.

### Compare graphs

```go
same := graphA.IsomorphicTo(graphB)
```

`IsomorphicTo` uses the canonical labelling algorithm from Jeremy Carroll, [Matching RDF Graphs](http://www.hpl.hp.com/techreports/2001/HPL-2001-293.pdf) (HP Labs, 2001).

## Data model

The types follow [RDF 1.1 Concepts](https://www.w3.org/TR/rdf11-concepts/).

| Type | Description |
| --- | --- |
| `Graph` | A set of triples. Read them with `IterTriples()`. |
| `Triple` | A `Subject`, a `Predicate` and an `Object`, each a `Term`. |
| `Term` | An interface implemented by `*IRI`, `*BlankNode` and `*Literal`. |
| `IRI` | An absolute IRI. |
| `BlankNode` | A node with a local label only. |
| `Literal` | A `LexicalForm` with an optional `DatatypeIRI` or `LanguageTag`. |

Every `Term` has three methods:

- `String()` returns the term in N-Triples style, for example `<http://example.org/a>`.
- `RawValue()` returns the plain value: the IRI, the blank node label, or the literal's lexical form.
- `Equals(other)` compares two terms.

## Status

gon3 is not complete. Its tests use the [W3C Turtle test suite](https://www.w3.org/2013/TurtleTests/), which is in `tests/turtle`. Some of those tests are turned off because they fail or hang. The known gaps are:

- Some IRI resolution cases.
- Some escape sequences in strings and IRIs.
- Some invalid documents that should fail but do not.
- SPARQL-style `PREFIX` written in mixed case, such as `PreFIX`.

See the commented-out entries in `ttl_test.go` for the full list.

## Test

```sh
go test ./...
```

## Credits

Ryan Chipman wrote the original gon3. The lexer uses [easylex](https://github.com/rychipman/easylex).

## License

Apache License 2.0. See [LICENSE](LICENSE).

The W3C test files in `tests/turtle` have their own license. See [tests/turtle/LICENSE](tests/turtle/LICENSE).
