// Package fold implements operations that map similar characters to a common
// target. These operations are called character foldings, and can be used
// to ignore certain distinctions between similar characters.
//
// Each folder implements the [transform.Transform] interface.
//
// DISCLAIMER: these folders are based on suggested foldings that appear in
// withdrawn drafts of Unicode technical reports. They may not be complete.
// Their names come from those technical reports.
//
// WARNING: folding is NOT appropriate for secure contexts -
// see [x/text/secure/precis] instead.
//
// See, for important commentary:
// - [Unicode Technical Report 30: CHARACTER FOLDINGS] (withdrawn, draft)
// - [Unicode Technical Report 25: CHARACTER FOLDINGS] (draft)
//
// [Unicode Technical Report 30: CHARACTER FOLDINGS]: http://www.unicode.org/reports/tr30/tr30-4.html
// [Unicode Technical Report 25: CHARACTER FOLDINGS]: http://www.unicode.org/L2/L2000/00261-tr25-0d1.html
package fold

import (
    "unicode"

    "github.com/tawesoft/golib/v2/operator"
    "github.com/tawesoft/golib/v2/text/dm"
    "github.com/tawesoft/golib/v2/text/np"
    "golang.org/x/text/runes"
    "golang.org/x/text/transform"
)

// Accents is a transformer that removes accents from Latin/Greek/Cyrillic
// characters.
var Accents = accents
var accents = transform.Chain(
    dm.CD.TransformerWithFilter(func (r rune) bool {
        return unicode.In(r, unicode.Latin, unicode.Greek, unicode.Cyrillic)
    }),
    runes.Remove(runes.Predicate(func (r rune) bool {
        return unicode.Is(unicode.Mn, r)
    })),
)

// CanonicalDuplicates is a transformer that folds duplicate singletons
// (usually when the same character, for historical reasons, has two different
// code points) (e.g. Ohm => Omega)
var CanonicalDuplicates = canonicalDuplicates
var canonicalDuplicates = dm.CD.TransformerWithFilter(func (r rune) bool {
    if operator.In(r,
        0x0374, 0x037E, 0x0387, 0x1FBE,
        0x1FEF, 0x1FFD, 0x2000, 0x2001,
        0x2126, 0x212A, 0x212B,
    ) {
        return true
    }
    if (r >= 0x2329) && (r <= 0x232A) { return true }
    return false
})

// Dashes is a transformer that folds everything in Unicode class Pd ("dash
// punctuation") to hyphen-minus '-'.
var Dashes = dashes
var dashes = runes.Map(func(r rune) rune {
    if unicode.Is(unicode.Pd, r) {
        return 0x002D // Hyphen-Minus
    }
    return r
})

// Digits is a transformer that folds digits in a native language or a
// typographical context to a substitute ASCII digit. Note that this maps to
// Unicode code points for the digits '0' to '9', not to the codepoints with
// integer values 0 to 9.
var Digits = digits
var digits = runes.Map(func(r rune) rune {
    ty, value := np.Get(r)
    if ty == np.Decimal || ty == np.Digit {
        if (value.Denominator == 1) && (value.Numerator >= 0) && (value.Numerator <= 9) {
            return '0' + rune(value.Numerator)
        }
    }
    return r
})

// GreekLetterforms is a transformer that folds alternative Greek letterforms
// e.g. 'ϐ' to 'β'.
var GreekLetterforms = greekLetterforms
var greekLetterforms = dm.KD.TransformerWithFilter(func (r rune) bool {
    switch {
        case (r >= 0x03D0) && (r <= 0x03D2): return true
        case (r >= 0x03D5) && (r <= 0x03D6): return true
        case (r >= 0x03F0) && (r <= 0x03F2): return true
        case (r >= 0x03F4) && (r <= 0x03F5): return true
        default: return false
    }
})

// HebrewAlternates is a transformer that folds e.g. wide Hebrew characters
// to non-wide variants.
var HebrewAlternates = hebrewAlternates
var hebrewAlternates = dm.KD.TransformerWithFilter(func (r rune) bool {
    return (r >= 0xFB20) && (r <= 0xFB28)
})

// Jamo folding converts from the Hangul Compatibility Jamo Unicode block to
// the Hangul Jamo Unicode block.
var Jamo = jamo
var jamo = dm.KD.TransformerWithFilter(func (r rune) bool {
    return (r >= 0x3131) && (r <= 0x3183)
})

// Math folding converts font variants, excluding the HebrewAlternates.
var Math = math
var math = dm.New(dm.Font).TransformerWithFilter(func (r rune) bool {
    return (r < 0xFB20) || (r > 0xFB28)
})

// NoBreak folding converts non-breaking space and non-breaking hyphens.
var NoBreak = noBreak
var noBreak = dm.New(dm.NoBreak).Transformer()

// Positional folding performs positional forms folding including Arabic ligatures.
var Positional = positional
var positional = dm.New(dm.Initial, dm.Medial, dm.Final, dm.Isolated).Transformer()

// Space folding converts all spaces to a single 0x0020 space.
var Space = space
var space = runes.Map(func(r rune) rune {
    if unicode.Is(unicode.Zs, r) {
        return 0x0020
    }
    return r
})

// Small folding converts small variant forms into normal forms.
var Small = small
var small = dm.New(dm.Small).Transformer()
