package gon3

import (
	"fmt"
)

type tokenType int

const (
	// tokens expressed as literal strings in http://www.w3.org/TR/turtle/#sec-grammar-grammar
	tokenAtPrefix = iota
	tokenAtBase
	tokenEndStatement
	tokenA
	tokenPredicateListSeparator
	tokenObjectListSeparator
	tokenStartBlankNodePropertyList
	tokenEndBlankNodePropertyList
	tokenStartCollection
	tokenEndCollection
	tokenLiteralDatatypeTag // TODO: rename
	tokenTrue
	tokenFalse

	// terminal tokens from http://www.w3.org/TR/turtle/#terminals
	tokenIRIRef
	tokenPNameNS
	tokenPNameLN
	tokenBlankNodeLabel
	tokenLangTag
	tokenInteger
	tokenDecimal
	tokenDouble
	tokenExponent
	tokenStringLiteralQuote
	tokenStringLiteralSingleQuote
	tokenStringLiteralLongQuote
	tokenStringLiteralLongSingleQuote
	tokenUChar
	tokenEChar
	tokenWhitespace
	tokenAnon
	tokenPNCharsBase
	tokenPNCharsU
	tokenPNChars
	tokenPNPrefix
	tokenPNLocal
	tokenPLX
	tokenPercent
	tokenHex
	tokenPNLocalEsc

	// TODO: use tokenError?
	// TODO: use tokenEOF?
)

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	switch t.typ {
	case tokenError:
		return t.val
	case tokenEOF:
		return "EOF"
	}
	if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}
