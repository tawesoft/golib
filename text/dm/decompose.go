// Package dm provides a way to query Unicode decomposition mappings and
// perform a custom compatibility decomposition using
// [compatibility mapping tags].
//
// This is slower than the optimised NFD and NFKD versions in
// [text/unicode/norm], so this package is only appropriate in situations where
// a custom decomposition is required.
//
// See [Unicode Normalization Forms] and [Character Decomposition Mapping].
//
// [compatibility mapping tags]: https://unicode.org/reports/tr44/#Formatting_Tags_Table
// [Unicode Normalization Forms]: https://unicode.org/reports/tr15/
// [Character Decomposition Mapping]: https://www.unicode.org/reports/tr44/#Character_Decomposition_Mappings
// [Stream-Safe Text Format]: https://unicode.org/reports/tr15/#Stream_Safe_Text_Format
package dm

import (
    "fmt"
    "sort"
    "unicode/utf8"

    "github.com/tawesoft/golib/v2/ks"
    "github.com/tawesoft/golib/v2/text/ccc"
    "golang.org/x/text/transform"
)

// Type is the compatibility formatting tag that controls decomposition
// mapping. The exact integer value is arbitrary and has no meaning.
type Type int

// IsCompat returns true if Type is any of the compatability mapping types
// i.e. is not Canonical, and is not None.
func (t Type) IsCompat() bool {
    return t >= Compat
}

// IsCanonical returns true if Type is a canonical mapping.
func (t Type) IsCanonical() bool {
    return t == Canonical
}

// note that the constants below MUST match those in internal/unicode/gen_test_decompose.go

const (
    None      Type =  0
    Canonical Type =  1 // Canonical
    Compat    Type =  2 // Otherwise unspecified compatibility character
    Encircled Type =  3 // Encircled form
    Final     Type =  4 // Final presentation form (Arabic)
    Font      Type =  5 // Font variant (for example, a blackletter form)
    Fraction  Type =  6 // Vulgar fraction form
    Initial   Type =  7 // Initial presentation form (Arabic)
    Isolated  Type =  8 // Isolated presentation form (Arabic)
    Medial    Type =  9 // Medial presentation form (Arabic)
    Narrow    Type = 10 // Narrow (or hankaku) compatibility character
    NoBreak   Type = 11 // No-break version of a space or hyphen
    Small     Type = 12 // Small variant form (CNS compatibility)
    Square    Type = 13 // CJK squared font variant
    Sub       Type = 14 // Subscript form
    Super     Type = 15 // Superscript form
    Vertical  Type = 16 // Vertical layout presentation form
    Wide      Type = 17 // Wide (or zenkaku) compatibility character
)

func (t Type) String() string {
    switch t {
        case None:      return "None"
        case Canonical: return "Canonical"
        case Compat:    return "Compat"
        case Encircled: return "Encircled"
        case Final:     return "Final"
        case Font:      return "Font"
        case Fraction:  return "Fraction"
        case Initial:   return "Initial"
        case Isolated:  return "Isolated"
        case Medial:    return "Medial"
        case Narrow:    return "Narrow"
        case NoBreak:   return "NoBreak"
        case Small:     return "Small"
        case Square:    return "Square"
        case Sub:       return "Sub"
        case Super:     return "Super"
        case Vertical:  return "Vertical"
        case Wide:      return "Wide"
    }
    ks.Never()
    return ""
}

// Map returns the decomposition mapping type and mappings for a single input
// rune. If there are no decomposition mappings, the returned type is None
// and the returned mappings are undefined.
//
// Note that, while the Unicode data files also have a default mapping of a
// character to itself, these are not counted (None is returned instead).
//
// Note also that this is a single mapping, not a full decomposition. For that,
// use dm.CD.String for a full canonical decomposition, dm.KD.String for a
// full compatibility decomposition, or [New] to define a custom decomposition
// and call the [Decomposer.String] method on it.
//
// For historical Unicode reasons, the longest compatibility mapping is 18
// characters long. Compatibility mappings are guaranteed to be no longer than
// 18 characters, although most consist of just a few characters.
func Map(r rune) (Type, []rune) {
    n := len(dtis)
    i := sort.Search(n, func(i int) bool {
        return r <= dtis[i].codepoint
    })

    if (i == n) || (dtis[i].codepoint != r) {
        return None, nil
    }

    dt := dtis[i].dt
    if dt == None {
        return None, nil
    }

    dmi := dtis[i].dmi
    dml := dtis[i].dml
    ms := dms[dmi:dmi+dml]

    return dt, ms
}

// Decomposer performs a full recursive decomposition of an input then
// applies the canonical reordering algorithm.
type Decomposer uint64

// CD is a Decomposer that performs a canonical decomposition
var CD  = Decomposer(1 << Canonical)

// KD is a Decomposer that performs a compatibility decomposition
var KD = Decomposer(0xFFFFFFFF)

// New returns a new Decomposer that performs a decomposition, but only
// for certain decomposition types.
func New(types ... Type) Decomposer {
    var d Decomposer
    return d.Extend(types...)
}

// Extend returns a canonical decomposer, extended with the compatibility
// mapping types given here, to create a new compatibility decomposer.
func Extend(types ... Type) Decomposer {
    return CD.Extend(types...)
}

// Except returns a compatability decomposer, except for the compatibility
// mapping types given here.
func Except(types ... Type) Decomposer {
    return KD.Except(types...)
}

// Extend returns a new Decomposer that performs a decomposition on the same
// types as its parent, in addition to those given here.
//
// For example:
//   decompose.CD.Extend(decompose.Super, decompose.Sub)
func (d Decomposer) Extend(types ... Type) Decomposer {
    for _, t := range types {
        d = Decomposer(uint64(d) | (uint64(1) << t))
    }
    return d
}

// Except returns a new Decomposer that performs a decomposition on the same
// types as its parent, except those given here.
//
// For example:
//   decompose.KD.Except(decompose.Super, decompose.Sub)
func (d Decomposer) Except(types ... Type) Decomposer {
    for _, t := range types {
        d = Decomposer(uint64(d) & (^(uint64(1) << t)))
    }
    return d
}

// Map returns the decomposition type and decomposition mapping for an input
// rune, provided that the decomposition type is one that the Decomposer
// supports. If the decomposition mapping is not supported, the unsupported
// type is returned and the returned mappings are nil. If there are no
// decomposition mappings, the returned type is None and the returned mappings
// are nil.
//
// Note that, while the Unicode data files also have a default mapping of a
// character to itself, these are not counted here (None is returned instead).
//
// Note also that this is a single mapping, not a full decomposition. For that,
// call [Decomposer.String] method, or [Decomposer.Rune] for a single rune.
//
// For historical Unicode reasons, the longest compatibility mapping is 18
// characters long. Compatibility mappings are guaranteed to be no longer than
// 18 characters, although most consist of just a few characters.
func (d Decomposer) Map(r rune) (Type, []rune) {
    dt, dm := Map(r)
    mask := uint64(1) << dt
    if (uint64(d) & mask) != mask { return dt, nil }
    return dt, dm
}

// String returns the full decomposition of s, but only applies the
// decomposition mappings that match the types registered with the Decomposer
// with [New], [Decomposer.Extend], or [Decomposer.Except].
func (d Decomposer) String(s string) (string, error) {
    dest := make([]rune, 0)
    d.flatten_r(&dest, []rune(s))
    xs := []byte(string(dest))
    err := ccc.Reorder(xs)
    if err != nil { return "", err }
    return string(xs), nil
}

func (d Decomposer) flatten_r(dest *[]rune, xs []rune) {
    // TODO Hangul rules (Go doesn't do these yet either)

    for _, x := range xs {
        dt, dm := d.Map(x)

        if dt == Canonical {
            // "A canonical mapping may also consist of a pair of characters,
            // but is never longer than two characters. When a canonical
            // mapping consists of a pair of characters, the first character
            // may itself be a character with a decomposition mapping, but the
            // second character never has a decomposition mapping."
            // TODO optimize for this case
        }

        if len(dm) == 0 {
            *dest = append(*dest, x)
        } else {
            d.flatten_r(dest, dm)
        }
    }
}

// Transformer returns an object implementing the [transform.Transform]
// interface that applies the decomposition specified by the decomposer across
// its input. It outputs the decomposed result.
//
// The returned transformer is stateless, so may be used concurrently.
func (d Decomposer) Transformer() transform.Transformer {
    return transform.Chain(
        mappingTransformer{d, nil},
        ccc.Transformer,
    )
}

// TransformerWithFilter is like [Transformer], however for each input rune x
// where filter(x) returns false, the decomposition process is skipped and that
// rune is output normally.
func (d Decomposer) TransformerWithFilter(filter func (x rune) bool) transform.Transformer {
    return transform.Chain(
        mappingTransformer{d, filter},
        ccc.Transformer,
    )
}

// mappingTransformer applies a mapping, but doesn't reorder the input
type mappingTransformer struct {
    d Decomposer
    filter func(x rune) bool
}

func (m mappingTransformer) Reset() {}
func (m mappingTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
    d := m.d
    for {
        r, rZ := utf8.DecodeRune(src[nSrc:])
        if r == utf8.RuneError {
            if (rZ == 0) && (atEOF) { return nDst, nSrc, nil }
            if atEOF { return nDst, nSrc, fmt.Errorf("invalid utf8 sequence") }
            return nDst, nSrc, transform.ErrShortSrc
        }

        if (m.filter != nil) && !m.filter(r) {
            if cap(dst) - nDst < rZ {
                return nDst, nSrc, transform.ErrShortDst
            }
            nDst += utf8.EncodeRune(dst[nDst:], r)
            nSrc += rZ
            continue
        }

        dt, dm := d.Map(r)

        if dt == Canonical {
            // "A canonical mapping may also consist of a pair of characters,
            // but is never longer than two characters. When a canonical
            // mapping consists of a pair of characters, the first character
            // may itself be a character with a decomposition mapping, but the
            // second character never has a decomposition mapping."
            // TODO optimize for this case
        }

        if len(dm) == 0 {
            // no mapping, so copy as-is
            if cap(dst) - nDst < rZ {
                return nDst, nSrc, transform.ErrShortDst
            }
            nDst += utf8.EncodeRune(dst[nDst:], r)
            nSrc += rZ
        } else {
            // Compatibility mappings are guaranteed to be no longer than
            // 18 characters.
            if cap(dst) - nDst < 18 * 4 { // max
                return nDst, nSrc, transform.ErrShortDst
            }
            buf := [18]rune{}
            bufs := buf[0:0]
            d.flatten_r(&bufs, dm)
            for i := 0; i < len(bufs); i++ {
                nDst += utf8.EncodeRune(dst[nDst:], bufs[i])
            }
            nSrc += rZ
        }
    }
}
