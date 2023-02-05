package parser

import (
    "strings"

    "github.com/tawesoft/golib/v2/css/tokenizer/token"
    "github.com/tawesoft/golib/v2/must"
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

func mirror(x token.Token) token.Token {
    switch {
        case x.Is(token.TypeLeftParen):         return token.RightParen()
        case x.Is(token.TypeLeftCurlyBracket):  return token.RightCurlyBracket()
        case x.Is(token.TypeLeftSquareBracket): return token.RightSquareBracket()
    }
    must.Never("invalid mirror")
    return token.Token{}
}

// TokenIsBlockStart returns true for a <(-token>, <{-token> or <[-token>.
func TokenIsBlockStart(t token.Token) bool {
    return t.Is(token.TypeLeftParen) ||
        t.Is(token.TypeLeftCurlyBracket) ||
        t.Is(token.TypeLeftSquareBracket)
}

// TokenIsNDimension returns true for a token that is a <n-dimension>: a
// <dimension-token> with its type flag set to "integer", and a unit that is an
// ASCII case-insensitive match for "n".
func TokenIsNDimension(t token.Token) bool {
    nt, _ := t.NumericValue()
    return true &&
        t.Is(token.TypeDimension) &&
        nt == token.NumberTypeInteger &&
        strings.EqualFold(t.Unit(), "n")
}

// TokenIsNdashdigitDimension returns true for a token that is a
// <ndashdigit-dimension>: a <dimension-token> with its type flag set to
// "integer", and a unit that is an ASCII case-insensitive match for "n-*",
// where "*" is a series of one or more digits.
func TokenIsNdashdigitDimension(t token.Token) bool {
    nt, _ := t.NumericValue()
    unit := t.Unit()
    return true &&
        t.Is(token.TypeDimension) &&
        nt == token.NumberTypeInteger &&
        len(unit) >= 3 &&
        strings.EqualFold(unit[0:2], "n-") &&
        eachRune(unit[2:], runeIsDigit)
}

// TokenIsNdashdigitIdent returns true for a token that is a
// <ndashdigit-ident>: an <ident-token> whose value is an ASCII
// case-insensitive match for "n-*", where "*" is a series of one or more
// digits.
func TokenIsNdashdigitIdent(t token.Token) bool {
    value := t.StringValue()
    return true &&
        t.Is(token.TypeIdent) &&
        len(value) >= 3 &&
        strings.EqualFold(value[0:2], "n-") &&
        eachRune(value[2:], runeIsDigit)
}

// TokenIsDashndashdigitIdent returns true for a token that is a
// <dashndashdigit-ident>: an <ident-token> whose value is an ASCII
// case-insensitive match for "-n-*", where "*" is a series of one or more
// digits.
func TokenIsDashndashdigitIdent(t token.Token) bool {
    value := t.StringValue()
    return true &&
        t.Is(token.TypeIdent) &&
        len(value) >= 4 &&
        strings.EqualFold(value[0:3], "-n-") &&
        eachRune(value[3:], runeIsDigit)
}

// TokenIsInteger returns true for a token that is an <integer>: a
// <number-token> with its type flag set to "integer".
func TokenIsInteger(t token.Token) bool {
    nt, _ := t.NumericValue()
    return true &&
        t.Is(token.TypeNumber) &&
        nt == token.NumberTypeInteger
}

// TokenIsSignedInteger returns true for a token that is a <signed-integer>: a
// <number-token> with its type flag set to "integer", and whose representation
// starts with "+" or "-".
func TokenIsSignedInteger(t token.Token) bool {
    nt, _ := t.NumericValue()
    repr := t.Repr()
    return true &&
        t.Is(token.TypeNumber) &&
        nt == token.NumberTypeInteger &&
        len(repr) > 0 &&
        (repr[0] == '+') || (repr[0] == '-')
}

// TokenIsSignlessInteger returns true for a token that is a
// <signless-integer>: a <number-token> with its type flag set to "integer",
// and whose representation starts with a digit.
func TokenIsSignlessInteger(t token.Token) bool {
    nt, _ := t.NumericValue()
    repr := t.Repr()
    return true &&
        t.Is(token.TypeNumber) &&
        nt == token.NumberTypeInteger &&
        len(repr) > 0 && runeIsDigit(rune(repr[0]))
}
