// Package descriptor parses a RNBF rule descriptor.
//
// See: https://unicode-org.github.io/icu-docs/apidoc/released/icu4c/classicu_1_1RuleBasedNumberFormat.html#details
package descriptor

import (
    "fmt"
    "strconv"
    "strings"
    "unicode"

    "golang.org/x/text/runes"
)

type Type uint8
const (
    TypeDefault = Type(iota)
    TypeBaseValue
    TypeBaseValueAndRadix
    TypeNegativeNumber
    TypeProperFraction
    TypeImproperFraction
    TypeInfinity
    TypeNaN
)

type Descriptor struct {
    Type Type
    Base int64
    Divisor int64
}

// decimalFold is a transformer for decimal integers that filters out
// spaces, period, and commas.
var decimalFold = runes.Remove(runes.Predicate(func(r rune) bool {
        if (r >= '0') && (r <= '9') { return false } // quick case
        if (r == '.') || (r == ',') { return true }
        return unicode.IsSpace(r)
}))

// isAsciiDigitString returns true iff the only characters in x are the ASCII
// digits '0' to '9', and x contains at least one digit.
func isAsciiDigitString(x string) bool {
    if len(x) == 0 { return false }
    for _, c := range x {
        if (c < '0') || (c > '9') { return false }
    }
    return true
}

// Parse parses a rule descriptor in the following format, or panics:
//
// * bv and rad are the names of tokens formatted as decimal numbers expressed
//   using ASCII digits with spaces, period, and commas ignored.
// * bv specifies the rule's base value. The rule's divisor is the highest power
//   of 10 less than or equal to the base value.
// * bv/rad: The rule's divisor is the highest power of rad less than or equal to
//   the base value.
// * -x: The rule is a negative-number rule.
// * x.x: The rule is an improper fraction rule.
// * 0.x: The rule is a proper fraction rule.
// * x.0: The rule is a default rule.
// * Inf: The rule for infinity.
// * NaN: The rule for an IEEE 754 NaN (not a number).
func Parse(s string) Descriptor {
    const ferr = "invalid rule descriptor syntax %q"

    switch s {
        case "-x":  { return Descriptor{Type: TypeNegativeNumber} }
        case "x.x": { return Descriptor{Type: TypeImproperFraction} }
        case "0.x": { return Descriptor{Type: TypeProperFraction} }
        case "x.0": { return Descriptor{Type: TypeDefault} }
        case "Inf": { return Descriptor{Type: TypeInfinity} }
        case "NaN": { return Descriptor{Type: TypeNaN} }
    }

    idx := strings.IndexRune(s, '/')
    if idx > 0 {
        left := decimalFold.String(s[:idx])
        right := decimalFold.String(s[idx+1:])

        if (!isAsciiDigitString(left)) || (!isAsciiDigitString(right)) {
            panic(fmt.Errorf(ferr, s))
        }

        base,  baseErr  := strconv.ParseInt(left, 10, 64)
        radix, radixErr := strconv.ParseInt(right, 10, 64)
        if (baseErr != nil) || (radixErr != nil) {
            panic(fmt.Errorf("rule descriptor range error while parsing %q", s))
        }

        return Descriptor{
            Base:    base,
            Divisor: radix,
            Type:    TypeBaseValueAndRadix,
        }
    } else if idx < 0 {
        s = decimalFold.String(s)

        if !isAsciiDigitString(s) {
            panic(fmt.Errorf(ferr, s))
        }

        base, baseErr  := strconv.ParseInt(s, 10, 64)
        if baseErr != nil {
            panic(fmt.Errorf("rule descriptor range error while parsing %q", s))
        }

        return Descriptor{
            Base: base,
            Type: TypeBaseValue,
        }
    } else {
        panic(fmt.Errorf(ferr, s))
    }
}
