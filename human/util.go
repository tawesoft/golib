package human

import (
    "fmt"
    "sort"
    "strconv"
    "unicode"
    "unicode/utf8"

    "github.com/tawesoft/golib/v2/numbers"
    "golang.org/x/text/language"
)

func runeInSet(rs []rune, r rune) bool {
    for i := 0; i < len(rs); r++ {
        if rs[i] == r { return true }
    }
    return false
}

// acceptRune returns the length of r in bytes if r is the first rune in s,
// otherwise returns zero.
func acceptRune(r rune, s string) int {
    n, z := utf8.DecodeRuneInString(s)
    if n == utf8.RuneError { return 0 }
    if r != n { return 0 }
    return z
}

// acceptRuneFromSet is like [acceptRune], but for any rune in the input set.
func acceptRuneFromSet(rs []rune, s string) int {
    n, z := utf8.DecodeRuneInString(s)
    if n == utf8.RuneError { return 0 }
    for i := 0; i < len(rs); i++ {
        if rs[i] == n { return z }
    }
    return 0
}

// acceptLeading returns the number of leading runes in s where each rune is
// in prefixSet, and the total length in bytes.
func acceptLeading(prefixSet []rune, s string) (int, int) {
    var c int
    var bytes int

    for {
        size := acceptRuneFromSet(prefixSet, s)
        if size == 0 { break }
        c++
        bytes += size
    }

    return c, bytes
}

// acceptNumberPart returns something that looks like some numbers, e.g.
// "12345", without spaces, group separators, or decimal separators.
func acceptNumberPart[N int64|float64](
    s string,
    groupSeparators []rune,
    guessRuneValue func(c rune) (int, bool),
    checkedMul func(N, N) (N, bool),
    checkedAdd func(N, N) (N, bool),
) (N, int, error) {
    var accu N
    var offset int
    var ok bool

    for _, c := range s {
        if runeInSet(groupSeparators, c) {
            // pass
        } else if unicode.IsSpace(c) {
            // pass
        } else if d, ok := guessRuneValue(c); ok {
            // shift left a digit and add c:
            // (123, 4) => 1230 + 4 => 1234
            accu, ok = checkedMul(accu, N(10))
            if !ok { return 0, 0, strconv.ErrRange }
            accu, ok = checkedAdd(accu, N(d))
            if !ok { return 0, 0, strconv.ErrRange }
        } else {
            // can't parse
            break
        }

        offset, ok = numbers.Int.CheckedAdd(offset, utf8.RuneLen(c))
        if !ok { break }
    }

    return accu, offset, nil
}

// findRepeatingRune returns any rune that appears more than once in a given
// string
func findRepeatingRune(s string) (rune, bool) {
    // easy efficient algorithm: sort the string, then walk it and see if the
    // current rune repeats
    sl := []rune(s)
    sort.Slice(sl, func(i int, j int) bool { return sl[i] < sl[j] })
    current := rune(0)

    for _, c := range sl {
        if c == current {
            return c, true
        }
        current = c
    }

    return 0, false
}

// invalidNumberParser returns a NumberParser that always generates an error
// when its methods are called.
func invalidNumberParser(t string, tag language.Tag, err error) NumberParser {
    t = fmt.Sprintf("invalid number parser (%s, %s)", t, tag.String())
    err = fmt.Errorf("%s: %w", t, err)
    return &errorNumberParser{t, err}
}

type errorNumberParser struct{t string; err error}

func (p *errorNumberParser) AcceptInt(string) (int64, int, error) {
    return 0, 0, p.err
}

func (p *errorNumberParser) AcceptFloat(string) (float64, int, error) {
    return 0, 0, p.err
}

func (p *errorNumberParser) AcceptSignedInt(string) (int64, int, error) {
    return 0, 0, p.err
}

func (p *errorNumberParser) AcceptSignedFloat(string) (float64, int, error) {
    return 0, 0, p.err
}

func (p *errorNumberParser) String() string {
    return p.t
}
