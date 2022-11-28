// Package token defines CSS tokens produced by a tokenizer.
package token

import (
    "fmt"
    "unicode/utf8"
)

type Type string
const (
    TypeWhitespace         = Type("whitespace-token")
    TypeEOF                = Type("EOF-token")
    TypeString             = Type("string-token")
    TypeBadString          = Type("bad-string-token")
    TypeDelim              = Type("delim-token")
    TypeComma              = Type("comma-token")
    TypeHash               = Type("hash-token")
    TypeLeftParen          = Type("(-token")
    TypeRightParen         = Type(")-token")
    TypeNumber             = Type("number-token")
    TypeDimension          = Type("dimension-token")
    TypePercentage         = Type("percentage-token")
    TypeCDC                = Type("CDC-token")
    TypeIdent              = Type("ident-token")
    TypeFunction           = Type("function-token")
    TypeUrl                = Type("url-token")
    TypeBadUrl             = Type("bad-url-token")
    TypeColon              = Type("colon-token")
    TypeSemicolon          = Type("semicolon-token")
    TypeCDO                = Type("CDO-token")
    TypeAtKeyword          = Type("at-keyword-token")
    TypeLeftSquareBracket  = Type("[-token")
    TypeRightSquareBracket = Type("]-token")
    TypeLeftCurlyBracket   = Type("{-token")
    TypeRightCurlyBracket  = Type("}-token")
)

type HashType string
const (
    HashTypeID           = HashType("id")
    HashTypeUnrestricted = HashType("unrestricted")
)

type NumberType string
const (
    NumberTypeInteger = NumberType("integer")
    NumberTypeNumber  = NumberType("number")
)

type Token struct {
    _type Type

    // repr preserves details such as whether .009 was written as .009 or 9e-3, and
    // whether a character was written literally or as a CSS escape. Only used by
    // <number-token>, <dimension-token>, <percentage-token>, ...
    repr string

    // Value is used by <ident-token>, <function-token>, <at-keyword-token>,
    // <hash-token>, <string-token>, and <url-token>.
    stringValue string

    unit string // used by <dimension-token>.

    // Hash tokens have a type flag set to either "id" or "unrestricted".
    // The type flag defaults to "unrestricted" if not otherwise set.
    hashType HashType

    // <number-token> and <dimension-token> additionally have a type flag set
    // to either "integer" or "number".
    // NOTE: spec is a bit ambiguous, but assume this also applies to
    // percentage-tokens.
    numberType NumberType

    delim rune // used by <delim-token>

    // numberValue is used by <number-token>, <dimension-token>,
    // <percentage-token>.
    numberValue float64
}

func (t Token) Is(x Type) bool {
    return t._type == x
}

func (t Token) Type() Type {
    return t._type
}

func (t Token) String() string {
    switch t._type {
        case TypeString:    fallthrough
        case TypeAtKeyword: fallthrough
        case TypeUrl:       fallthrough
        case TypeFunction:  fallthrough
        case TypeIdent:
            return fmt.Sprintf("<%s>{value: %q}", t._type, t.stringValue)
        case TypeDelim:
            return fmt.Sprintf("<%s>{delim: %q}", t._type, t.delim)
        case TypeHash:
            return fmt.Sprintf("<%s>{type: %q, value: %q}", t._type, t.hashType, t.stringValue)
        case TypeNumber:
            fallthrough
        case TypePercentage:
            return fmt.Sprintf("<%s>{type: %q, value: %f, repr: %q}", t._type, t.numberType, t.numberValue, t.repr)
        case TypeDimension:
            return fmt.Sprintf("<%s>{type: %q, value: %f, unit: %q, repr: %q}", t._type, t.numberType, t.numberValue, t.unit, t.repr)
        default:
            return fmt.Sprintf("<%s>", t._type)
    }
}

// Equals returns true iff the given two tokens are of the same type and of
// the same value (and other applicable details, such as hash type, number
// type, dimension unit, etc).
//
// In the case of <number-token>, <percentage-token>, and <dimension-token>,
// the tokens are also only considered equal if their underlying representation
// (i.e. the result of the [Token.Repr] method) is exactly equal.
func Equals(a Token, b Token) bool {
    if a._type != b._type { return false }
    switch a._type {
        case TypeHash:
            if a.hashType != b.hashType { return false }
            fallthrough
        case TypeString:    fallthrough
        case TypeAtKeyword: fallthrough
        case TypeUrl:       fallthrough
        case TypeFunction:  fallthrough
        case TypeIdent:
            return a.stringValue == b.stringValue

        case TypeDelim:
            return a.delim == b.delim

        case TypeDimension:
            if a.unit != b.unit { return false }
            fallthrough
        case TypeNumber:
            fallthrough
        case TypePercentage:
            if a.numberType != b.numberType { return false }
            if a.repr != b.repr { return false }
            return true

        default:
            return true
    }
    return false
}

// Repr returns the original representation of a token (or part of a token).
// This preserves details such as whether .009 was written as .009 or 9e-3 in
// a number. This is only valid for certain types. If called on an unsupported
// type, this function returns an empty string.
func (t Token) Repr() string {
    switch t._type {
        case TypeNumber:
        case TypePercentage:
        case TypeDimension:
            return t.repr
    }
    return ""
}

// StringValue returns the string value of a <ident-token>, <function-token>,
// <at-keyword-token>, <hash-token>, <string-token>, or <url-token>, or the
// empty string if the token is not one of these types.
func (t Token) StringValue() string {
    switch t._type {
        case TypeHash:      fallthrough
        case TypeString:    fallthrough
        case TypeAtKeyword: fallthrough
        case TypeUrl:       fallthrough
        case TypeFunction:  fallthrough
        case TypeIdent:
            return t.stringValue
    }
    return ""
}

// NumericValue returns the numeric value of a <number-token>,
// <percentage-token>, or <dimension-token>, as well as the number type
// (NumberTypeInteger or NumberTypeNumber). If the token is not one of these
// types, returns (0, NumberType(""))
func (t Token) NumericValue() (float64, NumberType) {
    switch t._type {
        case TypeNumber:     fallthrough
        case TypePercentage: fallthrough
        case TypeDimension:
            return t.numberValue, t.numberType
    }
    return 0, ""
}

// IsNumeric returns true if a token is a <number-token>,
// <percentage-token>, or <dimension-token>.
func (t Token) IsNumeric() bool {
    switch t._type {
        case TypeNumber:     fallthrough
        case TypePercentage: fallthrough
        case TypeDimension:
            return true
    }
    return false
}

// Unit returns the unit of a <dimension-token>. If the token is not a dimension
// token, returns "".
func (t Token) Unit() string {
    if t._type == TypeDimension {
        return t.unit
    } else {
        return ""
    }
}

// HashType returns the hash type of a <hash-token>. If the token is not a
// hash token, returns HashType("")
func (t Token) HashType() HashType {
    if t._type == TypeHash {
        return t.hashType
    } else {
        return ""
    }
}

// Delim returns the delimeter of a <delim-token>. If the token is not a
// delim token, returns utf.RuneError.
func (t Token) Delim() rune {
    if t._type == TypeDelim {
        return t.delim
    } else {
        return utf8.RuneError
    }
}

func EOF()                Token { return Token{_type: TypeEOF} }
func CDC()                Token { return Token{_type: TypeCDC} }
func CDO()                Token { return Token{_type: TypeCDO} }
func Colon()              Token { return Token{_type: TypeColon} }
func Comma()              Token { return Token{_type: TypeComma} }
func BadUrl()             Token { return Token{_type: TypeBadUrl} }
func BadString()          Token { return Token{_type: TypeBadString} }
func Semicolon()          Token { return Token{_type: TypeSemicolon} }
func Whitespace()         Token { return Token{_type: TypeWhitespace} }

func LeftParen()          Token { return Token{_type: TypeLeftParen} }
func RightParen()         Token { return Token{_type: TypeRightParen} }
func LeftSquareBracket()  Token { return Token{_type: TypeLeftSquareBracket} }
func RightSquareBracket() Token { return Token{_type: TypeRightSquareBracket} }
func LeftCurlyBracket()   Token { return Token{_type: TypeLeftCurlyBracket} }
func RightCurlyBracket()  Token { return Token{_type: TypeRightCurlyBracket} }

func String(s string) Token {
    return Token{
        _type:        TypeString,
        stringValue: s,
    }
}

func Delim(x rune) Token {
    return Token{
        _type:  TypeDelim,
        delim: x,
    }
}

func Hash(t HashType, s string) Token {
    return Token{
        _type:        TypeHash,
        stringValue: s,
        hashType:    t,
    }
}

func Number(nt NumberType, repr string, value float64) Token {
    return Token{
        _type:        TypeNumber,
        repr:        repr,
        numberValue: value,
        numberType:  nt,
    }
}

func Percentage(nt NumberType, repr string, value float64) Token {
    return Token{
        _type:        TypePercentage,
        repr:        repr,
        numberValue: value,
        numberType:  nt,
    }
}

func Dimension(nt NumberType, repr string, value float64, unit string) Token {
    return Token{
        _type:        TypeDimension,
        repr:        repr,
        numberValue: value,
        numberType:  nt,
        unit:        unit,
    }
}

func Ident(s string) Token {
    return Token{
        _type:        TypeIdent,
        stringValue: s,
    }
}

func Function(s string) Token {
    return Token{
        _type:        TypeFunction,
        stringValue: s,
    }
}

func Url(s string) Token {
    return Token{
        _type:        TypeUrl,
        stringValue: s,
    }
}

func AtKeyword(s string) Token {
    return Token{
        _type:        TypeAtKeyword,
        stringValue: s,
    }
}
