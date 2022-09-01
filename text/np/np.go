// Package np provides a way to query the Numeric Properties of a Unicode
// code point. This allows, for example, parsing the value of digits and
// numerals in other languages.
package np

import (
    "github.com/tawesoft/golib/v2/must"
)

// Type returns the type of a numeral.
//
//  - Decimal is a numeral in a decimal-radix number system, such as the ASCII
//    digits 0-9, Devanagari digits, Arabic digits, etc.
//  - Digit is like Decimal, but in some typographic context, e.g. superscript,
//    a number in a circle, etc.
//  - Numeric is a numeral that has a value, but does not appear in a decimal
//    system. For example, fractions, Tamil numbers, or Roman numerals.
type Type int

const (
    None    = 0
    Decimal = 1 // decimal-radix numeral e.g. 0-9
    Digit   = 2 // typographic context e.g. superscript
    Numeric = 3 // non-decimal e.g. Roman numerals
)

// Get returns, for the given codepoint, its Type and Value. If not found, Type
// is None.
func Get(x rune) (Type, Fraction) {
    s := getspan(x)
    if s.length == 0 { return 0, Fraction{} }

    offset := int(x - s.codepoint)
    must.True(offset >= 0)
    must.True(offset < s.length)

    return Type(s.nt), Fraction{
        Numerator:   s.nv.Numerator + int64(offset),
        Denominator: s.nv.Denominator,
    }
}
