# gon3: agent guide

A guide for AI-assisted development of this repo, whatever harness is reading
it. gon3 is a Go library that parses Turtle (and N-Triples) into an RDF graph
and checks whether two graphs are isomorphic. It has one dependency,
`github.com/rychipman/easylex`, and no binary. `README.md` is the user-facing
documentation of the public API.

## Repo map

| Path | What |
|---|---|
| `token.go` | Token types the lexer emits |
| `matchers.go` | Character-class matchers for the Turtle grammar |
| `statefns.go` | Lexer state functions (one per lexing state, on easylex) |
| `lexer.go` | The `lexer` interface the parser reads tokens from, plus unused mock lexers |
| `parser.go` | `NewParser` and `Parse`: tokens to a `*Graph` |
| `rdf.go` | Data model: `Term`, `IRI`, `BlankNode`, `Literal`, `Triple`, `Graph` |
| `isomorphism.go` | Canonical labelling (Carroll 2001) behind `IsomorphicTo` |
| `ttl_test.go` | Runs the W3C Turtle suite: positive, positive without isomorphism, negative |
| `rdf_test.go` | Unit tests for the data model helpers |
| `tests/turtle/` | The W3C Turtle test suite, under its own `LICENSE` |
| `.github/workflows/test.yml` | CI |

## Landmines

- **`tests/turtle/` is third-party data.** Never edit a fixture to make a test
  pass. It is the W3C suite, licensed separately from the code.
- **Suite cases are listed by hand.** `ttl_test.go` holds three name lists
  (`positiveParserTests`, `positiveParserTestsNoIso`, `negativeParserTests`).
  A file in `tests/turtle/` runs only if its name is in a list. A commented-out
  name is turned off: it failed, it hangs, or it is off on purpose (read the
  comment above it). Some now pass. Run a case with a timeout
  (`go test -timeout 30s`) before you uncomment it, because a hanging case
  stops the whole run.
- **A `Parser` holds the state of one document.** Code and tests create a new
  parser for each document.
- **`lexer.go` holds unused mock lexers in a non-test file.** No test calls
  them, but they ship in the package. Do not add more code to this file.
- **The `go.mod` floor is a promise.** `go.mod` and `README.md` both state the
  minimum Go version, and CI tests on it. Do not use a newer language feature
  or standard library API without raising all three together.

## Verification

All three must pass before anything is called done:

    gofmt -l .        # must print nothing
    go vet ./...
    go test -race ./...

CI (`.github/workflows/test.yml`) runs build, vet and test on push to `master`
and on every pull request, on the `go.mod` minimum and on the latest stable Go.
To check the minimum locally: `GOTOOLCHAIN=go<version>.0 go test -race ./...`.

## Team agreements

- **Commits** carry no AI co-author trailers. Commit only when asked.
- **One change at a time**: no refactor folded into a fix or a feature.
- **Comments say why**, in full sentences.
- **Bugs get a test first**: for a non-trivial parser bug, add the failing
  case (a suite name uncommented, or a new unit test), then fix it.

## Keep these docs truthful and brief

If your change makes a line here or in `README.md` false, fix it in the same
change. A stale line is worse than a missing one, because the next agent will
believe it. Search the whole file for the old claim, not only the section you
are editing. Add a line only where it would have saved you time.

## Harness-specific config

This file is the source of truth for how to work here. Root `CLAUDE.md` is one
line (`@AGENTS.md`) so Claude Code reads this file and nothing drifts.
