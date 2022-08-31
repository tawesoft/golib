// Package fallback implements [Unicode Character Fallback Substitutions] using
// the Unicode CLDR 41.0 supplemental data file characters.xml, and an
// algorithm for enumerating every canonically equivalent string.
//
// This can be useful for robustly parsing Unicode strings where for practical
// reasons (e.g. missing keyboard keys, missing font support) certain
// fallbacks have been used, or for picking a sensible default when certain
// Unicode strings cannot be displayed (e.g. missing font support).
//
// Note that care must be taken not to change the meaning of a text - for
// example, superscript two '²', will have a (last resort) Character Fallback
// Substitution to the digit '2' via NKFC normalisation, but these have
// entirely different meanings. Similarly, the string "1½" changes meaning if
// naively converted to "11/2". The Unicode Character Fallback Substitutions
// rules as implemented in this package would produce "1 1/2", but this doesn't
// help for superscript two.
//
// See the (withdrawn draft) Unicode Technical Report [30: CHARACTER FOLDINGS],
// as well as the earlier draft Unicode Technical Report [25: CHARACTER FOLDINGS], for commentary.
//
// [Unicode Character Fallback Substitutions]: https://unicode-org.github.io/cldr-staging/charts/41/supplemental/character_fallback_substitutions.html
// [30: CHARACTER FOLDINGS]: http://www.unicode.org/reports/tr30/tr30-4.html
// [25: CHARACTER FOLDINGS]: http://www.unicode.org/L2/L2000/00261-tr25-0d1.html
package fallback

import (
    "bytes"
    "fmt"
    "strings"
    "unicode/utf8"

    lazy "github.com/tawesoft/golib/v2/iter"
    "github.com/tawesoft/golib/v2/ks"
    "github.com/tawesoft/golib/v2/numbers"
    "github.com/tawesoft/golib/v2/text/ccc"
    "golang.org/x/text/unicode/norm"
)

/*
func New() {
}
*/

// Is returns true iff the provided string is a possible fallback string
// produced by Unicode Character Fallback Substitution rules applied to the
// input rune. Neither argument is required to be normalised on input.
//
// For example,
//
//   Is('㎦', "㎞³") // true
//   Is('㎦', "km³") // true
//   Is('㎦', "km3") // true
func Is(r rune, s string) bool {
    // 1. toNFC(value)
    // 2. other canonically equivalent sequences, if there are any
    // Two sequences are canonically equivalent if
    //   * toNFD(a) == toNFD(b), or
    //   * toNFC(a) == toNFC(b)

    q := string(r)
    if !norm.NFC.IsNormalString(q) { q = norm.NFC.String(q) }
    if !norm.NFC.IsNormalString(s) { s = norm.NFC.String(s) }
    if q == s { return true }

    // 3. the explicit substitutes value from characters.xml (in order)
    for _, u := range getsubs(r) {
        if (s == u) { return true }
    }

    // 4. toNFKC(value)
    if norm.NFKC.String(q) == norm.NFKC.String(s) { return true }

    return false
}

// Subs returns a complete list of strings that can be used as fallbacks for
// the input rune, in order of priority, according to the Unicode Character
// Fallback Substitutions rules.
func Subs(x rune) []string {
    nfc, nfkc := norm.NFC, norm.NFKC
    s := string(x)
    snfc := s // NDC normal
    xs := make([]string, 0)
    seen := make(map[string]struct{})

    see := func(x string) {
        if _, exists := seen[x]; !exists {
            seen[x] = struct{}{}
            xs = append(xs, x)
        }
    }
    see(s)

    // 1. toNFC(value)
    if !nfc.IsNormalString(s) {
        snfc = nfc.String(s)
        see(snfc)
    }

    // 2. other canonically equivalent sequences, if there are any
    it, err := Equivalent(snfc)
    if err == nil {
        for {
            if i, ok := it(); ok {
                see(i)
            } else {
                break
            }
        }
    }

    // 3. the explicit substitutes value from characters.xml (in order)
    for _, u := range getsubs(x) {
        see(u)
    }

    // 4. toNFKC(value)
    if !nfkc.IsNormalString(s) {
        sn := nfkc.String(s)
        see(sn)
    }

    return xs
}

// Equivalent is a [lazy.It] that produces all strings canonically-equivalent
// to the input. Note that this is very expensive for large strings. Note also
// that this does not include any Unicode Character Fallback Substitutions.
//
// This is a clean-room implementation of Mark Davies' algorithm described at
// https://unicode.org/notes/tn5/#Enumerating_Equivalent_Strings
func Equivalent(in string) (lazy.It[string], error) {
    // we use this as a delimiter, so it must not appear in the input
    if strings.ContainsRune(in, utf8.RuneError) { return nil, fmt.Errorf("encoding error") }

    // 1. Transform the input string into its NFD form.
    in = norm.NFD.String(in)

    // 2. Partition the string into segments, with each starter character in the
    // string at the beginning of a segment.
    segs := segments(in) // FFFD delimited list of strings

    // 3. For each segment enumerate canonically equivalent forms.
    // -- calls equivalent for each segment
    var segVariants []string // a list of FFFD delimited lists of strings
    segVariants = lazy.AppendToSlice(segVariants,
        lazy.Map(equivalent,
            lazy.CutString(segs, 0xFFFD)))

    // 4. Enumerate the combinations of all forms of all segments.
    return combinations(segVariants)
}

// combinations produces every possible combination of segment variants.
//
// e.g. for two segments with two variants each:
//
//     (Segment 1, Variant 1) + (Segment 2, Variant 1),
//     (Segment 1, Variant 2) + (Segment 2, Variant 1),
//     (Segment 1, Variant 1) + (Segment 2, Variant 2)
//     (Segment 1, Variant 2) + (Segment 2, Variant 2)
func combinations(segVariants []string) (lazy.It[string], error) {

    // for n slices a, b, c... containing a_z, b_z, c_z substrings,
    // nCombinations = a_z * b_z * c_z

    nSegments := len(segVariants)
    nVariants := make([]int, nSegments)
    for i, v := range segVariants {
        nVariants[i] = strings.Count(v, "\uFFFD") + 1
    }
    var nCombinations = 1
    for i := 0; i < nSegments; i++ {
        n, ok := numbers.Int.CheckedMul(nCombinations, nVariants[i])
        if !ok { return nil, fmt.Errorf("too many combinations") }
        nCombinations = n
    }

    // we can map each combination to a variant in each segment:
    // for each i in [0, Combinations)
    //   * let idx = i
    //   for each segment:
    //   * idx mod first_z gives the variant of the first segment
    //   * idx := idx / first_z, pop the first segment
    //
    // The outer loop is hoisted out to the caller because this is an iterator

    {
        i := 0
        sb := &strings.Builder{}
        return func() (string, bool) {
            if i == nCombinations {
                sb = nil
                return "", false
            }
            sb.Reset()

            idx := i
            for j := 0; j < nSegments; j++ {
                // 1st segment...
                first := idx % nVariants[j]
                // remaining segments
                idx = idx / nVariants[j]

                sb.WriteString(stringGet(segVariants[j], first))
            }

            i++
            return sb.String(), true
        }, nil
    }
}

// stringGet gets the n'th 0xFFFD delimited string
func stringGet (x string, n int) string {
    offset := 0
    seen := 0
    start := 0
    z := utf8.RuneLen(0xFFFD)
    for _, r := range x {
        rZ := utf8.RuneLen(r)

        if (r == 0xFFFD) && (rZ == z) {
            seen++
            if seen > n {
                return x[start:offset]
            }
            start = offset+z
        }

        offset += rZ
    }
    return x[start:]
}

// equivalent generates all strings canonically equivalent to the input
// segment. The return value is a string list delimited by utf8.RuneError.
func equivalent(in string) string {
    var sb strings.Builder
    equivalent_r(&sb, in)
    return sb.String()
}

func equivalent_r(sb *strings.Builder, in string) {
    if len(in) == 0 { return }

    if sb.Len() > 0 { sb.WriteRune(0xFFFD) }
    sb.WriteString(in)

    // a: Use the set of characters whose decomposition begins with the
    // segment's starter.
    r, _ := utf8.DecodeRuneInString(in)
    ks.Assert(r != utf8.RuneError)
    ds := dstarts(r)

    // b: For each character in this set
    var bb bytes.Buffer
    for _, c := range string(ds) {

        // i: Get the character's decomposition.
        dcom := norm.NFD.String(string(c))

        // ii: If the decomposition contains characters that are not in the
        // segment, then skip this character.
        for _, d := range dcom {
            if !strings.ContainsRune(in, d) {
                goto skipChar
            }
        }
        goto noskipchar
        skipChar: continue
        noskipchar:

        // iii: If the decomposition contains a character that is blocked
        // in the segment (preceded by a combining mark with the same
        // combining class), then also skip this character.
        // TODO

        // iv: Otherwise, start building a new string with this character.
        bb.Reset()
        bb.WriteRune(c)
        cZ := utf8.RuneLen(c)

        // v: Append all characters from the input segment that are not in
        // this character's decomposition in canonical order.
        for _, d := range in {
            if !strings.ContainsRune(dcom, d) {
                bb.WriteRune(d)
            }
        }
        // canonical order
        bytes := bb.Bytes()
        ccc.Reorder(bytes[cZ:])

        // vi: Add this string to the set of canonical equivalents for the
        // current segment.
        if len(bytes) > 0 {
            sb.WriteRune(0xFFFD)
            sb.Write(bytes)
        }

        // vii: Recurse: Treat all but the initial character of this new
        // string as a segment and add to the set for the current segment
        // all combinations of the initial character and the equivalent
        // strings of the rest.
        if len(bytes) > cZ {
            equivalent_r(sb, string(bytes[cZ:]))
        }
    }
}

// segments breaks an NFD normalized string into segments that start with a
// starter character. The return value is a string list delimited by
// utf8.RuneError.
func segments(in string) string {
    var segs strings.Builder

    // preallocate once, len(in) * 2 is the maximum even if every rune
    // forms a new segment
    segs.Grow(len(in) * 2)

    for _, c := range in {
        cls := ccc.Of(c)

        if cls == 0 {
            // new segment
            if segs.Len() > 0 { segs.WriteRune(utf8.RuneError) }
        }
        segs.WriteRune(c)
    }

    return segs.String()
}
