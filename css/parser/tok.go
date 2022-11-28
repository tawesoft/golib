package parser

import (
    "strings"

    "github.com/tawesoft/golib/v2/css/tokenizer/token"
)

// Helpful functions for parsing tokens into the <an+b> type
// https://www.w3.org/TR/css-syntax-3/#the-anb-type

// TODO maybe pull tokenizer util functions into internal or tokutil packge?
func runeIsDigit(x rune) bool {
    return (x >= '0') && (x <= '9')
}

func eachRune(xs string, pred func(x rune) bool) bool {
    for _, x := range xs {
        if !pred(x) { return false }
    }
    return true
}

// IsNDimensionToken returns true for a token that is a <n-dimension>: a
// <dimension-token> with its type flag set to "integer", and a unit that is an
// ASCII case-insensitive match for "n".
func IsNDimensionToken(t token.Token) bool {
    nt, _ := t.NumericValue()
    return true &&
        t.Is(token.TypeDimension) &&
        nt == token.NumberTypeInteger &&
        strings.EqualFold(t.Unit(), "n")
}

// IsNdashdigitDimensionToken returns true for a token that is a
// <ndashdigit-dimension>: a <dimension-token> with its type flag set to
// "integer", and a unit that is an ASCII case-insensitive match for "n-*",
// where "*" is a series of one or more digits.
func IsNdashdigitDimensionToken(t token.Token) bool {
    nt, _ := t.NumericValue()
    unit := t.Unit()
    return true &&
        t.Is(token.TypeDimension) &&
        nt == token.NumberTypeInteger &&
        len(unit) >= 3 &&
        strings.EqualFold(unit[0:2], "n-") &&
        eachRune(unit[2:], runeIsDigit)
}

// IsNdashdigitIdentToken returns true for a token that is a
// <ndashdigit-ident>: an <ident-token> whose value is an ASCII
// case-insensitive match for "n-*", where "*" is a series of one or more
// digits.
func IsNdashdigitIdentToken(t token.Token) bool {
    value := t.StringValue()
    return true &&
        t.Is(token.TypeIdent) &&
        len(value) >= 3 &&
        strings.EqualFold(value[0:2], "n-") &&
        eachRune(value[2:], runeIsDigit)
}

// IsDashndashdigitIdentToken returns true for a token that is a
// <dashndashdigit-ident>: an <ident-token> whose value is an ASCII
// case-insensitive match for "-n-*", where "*" is a series of one or more
// digits.
func IsDashndashdigitIdentToken(t token.Token) bool {
    value := t.StringValue()
    return true &&
        t.Is(token.TypeIdent) &&
        len(value) >= 4 &&
        strings.EqualFold(value[0:3], "-n-") &&
        eachRune(value[3:], runeIsDigit)
}

// IsIntegerToken returns true for a token that is an <integer>: a
// <number-token> with its type flag set to "integer".
func IsIntegerToken(t token.Token) bool {
    nt, _ := t.NumericValue()
    return true &&
        t.Is(token.TypeNumber) &&
        nt == token.NumberTypeInteger
}

// IsSignedIntegerToken returns true for a token that is a <signed-integer>: a
// <number-token> with its type flag set to "integer", and whose representation
// starts with "+" or "-".
func IsSignedIntegerToken(t token.Token) bool {
    nt, _ := t.NumericValue()
    repr := t.Repr()
    return true &&
        t.Is(token.TypeNumber) &&
        nt == token.NumberTypeInteger &&
        len(repr) > 0 &&
        (repr[0] == '+') || (repr[0] == '-')
}

// IsSignlessIntegerToken returns true fora token that is a <signless-integer>
// is a <number-token> with its type flag set to "integer", and whose
// representation starts with a digit.
func IsSignlessIntegerToken(t token.Token) bool {
    nt, _ := t.NumericValue()
    repr := t.Repr()
    return true &&
        t.Is(token.TypeNumber) &&
        nt == token.NumberTypeInteger &&
        len(repr) > 0 && runeIsDigit(rune(repr[0]))
}
