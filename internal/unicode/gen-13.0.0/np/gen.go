// Command gen.go generates Numeric Properties information
package main

import (
    "archive/zip"
    "bufio"
    "encoding/xml"
    "fmt"
    "io"
    "os"
    "sort"
    "strconv"
    "strings"

    "github.com/tawesoft/golib/v2/internal/legacy"
    "github.com/tawesoft/golib/v2/iter"
    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/operator"
)

type Char struct {
    codepoint rune
    coderange [2]rune // First and Last codepoint
    Nt string
    Nv string

    // Span is length of a span where Nv increments by one and Nt is equal.
    span int
}

func ParseCodepoint(x string) rune {
    return rune(must.Result(strconv.ParseInt(x, 16, 32)))
}

func (c Char) IsRange() bool {
    return (c.codepoint == 0) && ((c.coderange[0] != 0) || (c.coderange[1] != 0))
}

func CharFromAttrs(attr []xml.Attr, parent Char) Char {
    c := parent

    for _, a := range attr {
        switch a.Name.Local {
            case "cp":
                c.codepoint = ParseCodepoint(a.Value)
            case "first-cp":
                c.coderange[0] = ParseCodepoint(a.Value)
            case "last-cp":
                c.coderange[1] = ParseCodepoint(a.Value)
            case "nt":
                c.Nt = a.Value
            case "nv":
                c.Nv = a.Value
        }
    }
    return c
}

type NT int8
const (
    NTNone NT = 0
    NTDe   NT = 1 // Decimal
    NTDi   NT = 2 // Decimal, but in typographic context
    NTNu   NT = 3 // Numeric, but not decimal
)

func (c Char) NumericType() NT {
/*
    XML stores this as
        attribute nt { "None" | "De" | "Di" | "Nu" }??
*/
    switch c.Nt {
        case "None": return NTNone
        case "De":   return NTDe
        case "Di":   return NTDi
        case "Nu":   return NTNu
    }

    must.Never()
    return NTNone
}

type Fraction struct {
    Negative bool
    N int64 // Numerator
    D int64 // Denominator
}

func (f Fraction) IsInteger() bool {
   return f.D == 1
}

func (c Char) NumericValue() Fraction {
    neg := false
    left, right, found := strings.Cut(c.Nv, "-")
    if found {
        neg = true
    } else {
        right = left
    }

    left, right, found = strings.Cut(right, "/")
    if found {
        return Fraction{
            Negative: neg,
            N: must.Result(strconv.ParseInt(left, 10, 64)),
            D: must.Result(strconv.ParseInt(right, 10, 64)),
        }
    } else {
        return Fraction{
            Negative: neg,
            N: must.Result(strconv.ParseInt(left, 10, 64)),
            D: 1,
        }
    }
}

func main() {
    zr := must.Result(zip.OpenReader("../../DATA/ucd.nounihan.grouped.13.0.0.zip"))

    opener := func(name string) func() (io.ReadCloser, error) {
        return func() (io.ReadCloser, error) {
            return zr.Open(name)
        }
    }

    chars := make([]Char, 0)

    must.Check(legacy.WithCloser(opener("ucd.nounihan.grouped.xml"), func(f io.ReadCloser) error {
        d := xml.NewDecoder(bufio.NewReaderSize(f, 64 * 1024))
        var group Char
        var inRepertoire, inGroup bool

        for {
            tok, err := d.Token()
            if err == io.EOF {
                break
            } else if err != nil {
                panic(fmt.Errorf("xml decode error: %w", err))
            }

            switch ty := tok.(type) {
                case xml.StartElement:
                    switch ty.Name.Local {
                        case "repertoire":
                            if inRepertoire { must.Never() }
                            inRepertoire = true
                        case "group":
                            if !inRepertoire { break }
                            if inGroup { must.Never() }
                            inGroup = true
                            group = CharFromAttrs(ty.Attr, operator.Zero[Char]())
                        case "char":
                            if !inRepertoire { break }
                            if !inGroup { must.Never() }
                            c := CharFromAttrs(ty.Attr, group)
                            if c.codepoint == 0 { break }
                            if c.NumericType() == NTNone { break }
                            if c.IsRange() { must.Never() }
                            chars = append(chars, c)
                    }
                case xml.EndElement:
                    switch ty.Name.Local {
                        case "repertoire":
                            inRepertoire = false
                        case "group":
                            inGroup = false
                        case "char":
                            break
                    }
            }
        }

        return nil
    }))

    must.True(sort.SliceIsSorted(chars, func(i int, j int) bool {
        return chars[i].codepoint < chars[j].codepoint
    }))

    // Many are ranges of 0-9 that we can compact into "spans"
    for i := 0; i < len(chars); i++ {
        c := countSpan(chars, i)
        chars[i].span = c
        i += c - 1
    }

    /*
    for _, v := range chars {
        if v.span == 0 {
            // fmt.Printf("%x, %c: skip\n", v.codepoint, v.codepoint)
        } else {
            fmt.Printf("%x, %c, %s, %s [%d]\n", v.codepoint, v.codepoint, v.Nt, v.Nv, v.span)
        }
    }
    */

    filterInSpan := func (c Char) bool {
        return c.span > 0
    }
    nBefore := len(chars)
    chars = iter.ToSlice(iter.Filter(filterInSpan, iter.FromSlice(chars)))
    nAfter := len(chars)

    for _, v := range chars {
        fmt.Printf("%x, %c, %s, %s [%d]\n", v.codepoint, v.codepoint, v.Nt, v.Nv, v.span)
    }

    var maxN, maxD, maxPrefix int64
    var maxS, maxTrailing int
    for _, v := range chars {
        f := v.NumericValue()
        if f.N > maxN { maxN = f.N }
        if f.D > maxD { maxD = f.D }
        if v.span > maxS { maxS = v.span }
        p, t := split(f.N)
        if p > maxPrefix { maxPrefix = p }
        if t > maxTrailing { maxTrailing = t }
    }

    fmt.Printf("%d codepoints, packed as %d spans\n", nBefore, nAfter)
    fmt.Printf("max value n/d: %d/%d (max prefix %d, max trailing %d)\n", maxN, maxD, maxPrefix, maxTrailing)
    fmt.Printf("max span length: %d\n", maxS)

    {
        dest := must.Result(os.Create("../../../../text/np/np.bin"))
        defer dest.Close()
        nvEncode(dest, chars)
    }
}

// split a number into a prefix and a number of trailing zeros
func split(x int64) (int64, int) {
    if x == 0 { return 0, 0 }
    s := fmt.Sprintf("%d", x)
    t := strings.TrimRight(s, "0")
    return must.Result(strconv.ParseInt(t, 10, 64)), len(s) - len(t)
}

// Pack encodes a unicode character's codepoint (21 bits), numeric type nt
// (2 bits), and faction n/d with n encoded as int12 prefix + int6 number of
// trailing zeros and d encoded as (int12), and span encoded as int 8, and
// a bit to indicate negative fraction, with 2 bits padding.
func (c Char) Pack() uint64 {
    cp := c.codepoint
    nt := c.NumericType()
    nv := c.NumericValue()

    prefix, trailing := split(nv.N)

    must.True(nt       >=        1)
    must.True(nt       <=        3)
    must.True(nv.N     >=        0)
    must.True(prefix   <=   0x0FFF) // 2^12
    must.True(trailing <=     0x3F) // 2^ 6
    must.True(nv.D     <=   0x0FFF) // 2^12
    must.True(nv.D     >=        1)
    must.True(c.span   >=        1)
    must.True(c.span   <=     0xFF)

    var x uint64
    x |= (uint64(cp)  &      0x1FFFFF)
    x |= (uint64(nt)  &             3) << 21
    x |= (uint64(prefix) &     0x0FFF) << (21 + 2)
    x |= (uint64(trailing) &     0x3F) << (21 + 2 + 12)
    x |= (uint64(nv.D) &       0x0FFF) << (21 + 2 + 12 + 6)
    x |= (uint64(c.span) &       0xFF) << (21 + 2 + 12 + 6 + 12)

    if nv.Negative {
        x |= 1 << (21 + 2 + 12 + 6 + 12 + 8)
    }

    return x
}

func nvEncode(dest io.Writer, chars []Char) {
    for _, c := range chars {
        p := c.Pack()
        b0 :=       (p) & 0xFF
        b1 := (p >>  8) & 0xFF
        b2 := (p >> 16) & 0xFF
        b3 := (p >> 24) & 0xFF
        b4 := (p >> 32) & 0xFF
        b5 := (p >> 40) & 0xFF
        b6 := (p >> 48) & 0xFF
        b7 := (p >> 56) & 0xFF
        must.Result(dest.Write([]byte{
            uint8(b0),
            uint8(b1),
            uint8(b2),
            uint8(b3),
            uint8(b4),
            uint8(b5),
            uint8(b6),
            uint8(b7),
        }))
    }
}

func countSpan(chars []Char, idx int) int {
    c := chars[idx]
    i := 0

    for {
        i++

        if idx + i >= len(chars) { break }
        d := chars[idx + i]
        if c.Nt != d.Nt { break }
        cv, dv := c.NumericValue(), d.NumericValue()
        if !cv.IsInteger() { break }
        if !dv.IsInteger() { break }
        if cv.Negative { break }
        if dv.Negative { break }
        if dv.N != cv.N + int64(i) { break }
    }

    return i
}
