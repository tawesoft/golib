// Package fold implements operations that map similar characters to a common
// target. These operations are called character foldings, and can be used
// to ignore certain distinctions between similar characters.
//
// Each folder implements the [transform.Transform] interface.
//
// Note that folding is NOT appropriate for secure contexts -
// see [text/secure/precis] instead.
//
// See also, for important commentary:
// - [Unicode Technical Report 30: CHARACTER FOLDINGS] (withdrawn, draft)
// - [Unicode Technical Report 25: CHARACTER FOLDINGS] (draft)
//
// [Unicode Technical Report 30: CHARACTER FOLDINGS]: http://www.unicode.org/reports/tr30/tr30-4.html
// [Unicode Technical Report 25: CHARACTER FOLDINGS]: http://www.unicode.org/L2/L2000/00261-tr25-0d1.html

package fold

import (
    "fmt"
    "unicode"

    "github.com/tawesoft/golib/v2/text/dm"
    "golang.org/x/text/cases"
    "golang.org/x/text/language"
    "golang.org/x/text/runes"
    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
)

type Folder struct {
    // hidden internal TODO optimised byte buffer versions
    str func(string) string
}

// Transform implements the [transform.Transform] interface.
func (f Folder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {

    // TODO optimise using byte buffers directly so as not to allocate memory?
    s := []byte(f.str(string(src)))

    if len(s) > len(dst) {
        return 0, len(src), transform.ErrShortDst
    }

    copy(dst, s)
    return len(s), len(src), nil
}

// Reset implements the [transform.Transform] interface.
func (f Folder) Reset() {}


var Dashes = runes.Map(func(r rune) rune {
    if unicode.Is(unicode.Pd, r) {
        return 0x002D // Hyphen-Minus
    }
    return r
})

var NoBreak = dm.New(dm.NoBreak)

var Positional = dm.New(dm.Initial, dm.Medial, dm.Final, dm.Isolated)

var Space = runes.Map(func(r rune) rune {
    if unicode.Is(unicode.Zs, r) {
        return 0x0020
    }
    return r
})


func example() {
    t := transform.Chain(cases.Lower(language.English), dm.CD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
    s, _, _ := transform.String(t, "Résumé")
    fmt.Println(s)
}
