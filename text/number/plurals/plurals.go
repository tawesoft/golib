// Package plurals is an easy-to-use wrapper around the /x/text/feature/plural
// package.
//
// Plural rules controls how, for a given locale, plurals are counted in forms
// such as "1st", "2nd", "3rd", or "1 cat", "2 cats", etc.
package plurals

import (
    "strconv"
    "strings"

    "golang.org/x/text/feature/plural"
    "golang.org/x/text/language"
)

type Form byte
const (
    Other Form = iota
    Zero
    One
    Two
    Few
    Many
)

// Rules implements plural forms of a number for a given locale.
//
// Numbers represented by a string must be formatted with the ASCII digits
// '0' to '9', and optionally with a single decimal point encoded as the
// ASCII character '.'.
//
// If an argument is of the wrong format, or out of range of an int64, the
// result will be Other.
type Rules interface {
    // Ordinal is used to determine e.g. 0 days, 1 day, 2 days. For example,
    // for locale "en", Ordinal("2") == Many.
    Ordinal(number string) Form

    // Cardinal is used to determine e.g. 1st, 2nd, 3rd. For example,
    // for locale "en", Cardinal("3") == Few.
    Cardinal(number string) Form
}

// New returns a value implementing the Rules interface for that locale.
//
// e.g. plurals.New(language.MustParse("en"))
func New(locale language.Tag) Rules {
    return plurals{
        locale: locale,
    }
}

type plurals struct {
    locale language.Tag
}

func (p plurals) Ordinal(number string) Form {
    return form(p.locale, plural.Ordinal, number)
}

func (p plurals) Cardinal(number string) Form {
    return form(p.locale, plural.Cardinal, number)
}

func form(locale language.Tag, rules *plural.Rules, number string) Form {
    i, v, w, f, t, ok := operands(number)
    if !ok { return Other }
    return Form(rules.MatchPlural(locale, i, v, w, f, t))
}

// operands calculates the plural operands as described by
// https://unicode.org/reports/tr35/tr35-numbers.html#table-plural-operand-meanings
func operands(number string) (i, v, w, f, t int, ok bool) {
    number = strings.TrimLeft(number, "0") // remove leading zeros
    wrap := int64(10_000_000) // /x/text says this its ok for arguments to be mod this
    if len(number) == 0 { ok = true; return }

    onlyDigits := func(x string) bool {
        for _, c := range x {
            if (c < '0') || (c > '9') { return false }
        }
        return true
    }

    /*
        Plural Operand Meanings for N (*ignoring exponent notation)
        n: the absolute value of N.*
        i: the integer digits of N.*
        v: the number of visible fraction digits in N, with trailing zeros.*
        w: the number of visible fraction digits in N, without trailing zeros.*
        f: the visible fraction digits in N, with trailing zeros, expressed as an integer.*
        t: the visible fraction digits in N, without trailing zeros, expressed as an integer.*
    */

    point := strings.IndexRune(number, '.')
    if point == -1 { // integer
        if !onlyDigits(number) { ok = false; return }
        if len(number) == 0    { ok = false; return }

        if x, err := strconv.ParseInt(number, 10, 64); err == nil {
            i = int(x % wrap)
        } else {
            ok = false; return
        }
    } else { // has decimal point
        left, right := number[0:point], number[point+1:]
        if !onlyDigits(left)           { ok = false; return }
        if !onlyDigits(right)          { ok = false; return }
        if len(left) + len(right) == 0 { ok = false; return } // just "."!?

        if len(left) != 0 {
            if x, err := strconv.ParseInt(left, 10, 64); err == nil {
                i = int(x % wrap)
            } else {
                ok = false; return
            }
        }

        if len(right) != 0 {
            if x, err := strconv.ParseInt(right, 10, 64); err == nil {
                f = int(x % wrap)
                v = len(right)
            } else {
                ok = false; return
            }
        }

        rightTrimmed := strings.TrimRight(right, "0")
        if len(rightTrimmed) != 0 {
            if x, err := strconv.ParseInt(rightTrimmed, 10, 64); err == nil {
                t = int(x % wrap)
                w = len(rightTrimmed)
            } else {
                ok = false; return
            }
        }
    }

    ok = true; return
}
