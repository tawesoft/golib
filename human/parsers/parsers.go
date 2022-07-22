// Package parsers implements parsers
package parsers

import (
    "github.com/tawesoft/golib/v2/human"
    "golang.org/x/text/language"
)

// MatchDecimalParser picks a parser for strings that encoded decimal numbers
// based on the specified locale.
func MatchDecimalParser(t language.Tag) human.NumberParser {

    // arabic script matches all of
    //   - fa persian/farsi
    //   - ar arabic
    //   - sd Sindhi
    //   - ur Urdu
    script, _ := t.Script()
    if arabic, _ := language.Arabic.Script(); arabic == script {
        return EasternArabicDecimalParser
    }

    // european style...

    // Reference: Unicode ICU 71.1

    // fallback
    return human.NewDecimalParser(t)
}

// EasternArabicDecimalParser returns a parser for strings that encode decimal
// numbers in multiple types of Eastern Arabic scripts, including
// Arabic-derived scripts such as Persian, Sindhi, and Urdu.
var EasternArabicDecimalParser = &human.DecimalParser{
    // Reference: https://www.unicode.org/charts/PDF/U0600.pdf (v14.0)

    GroupSeparators: []rune{
        0x002C, // , comma
        0x0027, // ' apostrophe
        0x2019, // ’ right single quotation mark
        0x0603, // ، ARABIC COMMA
        0x066C, // ٬ ARABIC THOUSANDS SEPARATOR
        0x2E32, // ⸲ turned comma
        0x2E41, // ⹁ reversed comma
    },
    DecimalSeparators: []rune{
        0x066B, // ٫ ARABIC DECIMAL SEPARATOR
    },
    DecimalDigits: [10][]rune{
        // Arabic-Indic digits ("Arabic proper").
        // Starting with 0x0660 ٠ ARABIC-INDIC DIGIT ZERO

        // Eastern Arabic-Indic digits used with Arabic-script languages of
        // Iran, Pakistan, and India (Persian, Sindhi, Urdu, etc.).
        // Starting with 0x06F0 ۰ EXTENDED ARABIC-INDIC DIGIT ZERO
        {0x0660, 0x06F0},
        {0x0661, 0x06F1},
        {0x0662, 0x06F2},
        {0x0663, 0x06F3},
        {0x0664, 0x06F4},
        {0x0665, 0x06F5},
        {0x0666, 0x06F6},
        {0x0667, 0x06F7},
        {0x0668, 0x06F8},
        {0x0669, 0x06F9},
    },
}

// EuropeanDecimalParser returns a parser for strings that encode decimal
// numbers in the styles common in continental Europe, in Russia, and im many
// ex-territories of European colonial empires, where the decimal separator is
// a comma, and the group separator (separating groups of three digits) is
// often a dot, space, or apostrophe.
var EuropeanDecimalParser = &human.DecimalParser{
    GroupSeparators: []rune{
        // Whitespace is skipped, but common here in SI style.
        // U+0020 SPACE
        // U+2009 THIN SPACE
        // U+202F NARROW NO-BREAK SPACE
        0x0027, // ' APOSTROPHE
        0x002E, // . FULL STOP,
        0x02D9, // ˙ DOT ABOVE
    },
    DecimalSeparators: []rune{
        0x002C, // , comma
        0x00B7, // · MIDDLE DOT
        0x22C5, // ⋅ DOT OPERATOR
        0x2E33, // ⸳ RAISED DOT
        0x2396, // ⎖ DECIMAL SEPARATOR KEY SYMBOL
    },
    DecimalDigits:     [10][]rune{},
}
