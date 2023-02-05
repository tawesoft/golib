// Command gen.go generates Canonical Combining Class information
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

    "github.com/tawesoft/golib/v2/internal/legacy"
    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/operator"
)
type Char struct {
    codepoint rune
    coderange [2]rune // First and Last codepoint
    ccc uint8
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
            case "ccc":
                v := must.Result(strconv.ParseUint(a.Value, 10, 8))
                c.ccc = uint8(v)
        }
    }
    return c
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
                            if c.IsRange() {
                                if c.ccc != 0 {
                                    // not in this version
                                    must.Never()
                                }
                            } else {
                                chars = append(chars, c)
                            }
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

    ranges := cccRanges(chars)
    var size int
    for _, r := range ranges {
        if r.ccc == 0 { continue }
        z := int(r.end) - int(r.start)
        if z > size { size = z }
        fmt.Printf("ccc: range %+v\n", r)
    }
    fmt.Printf("ccc: max range size: %d\n", size) // gives 256

    must.True(sort.SliceIsSorted(ranges, func(i int, j int) bool {
        return ranges[i].start < ranges[j].start
    }))

    // check each codepoint has a range (ranges contains ranges with ccc=0
    // for this reason)
    for _, c := range chars {
        var seenInAnyRange bool
        for _, r := range ranges {
            if (c.codepoint >= r.start) && (c.codepoint < r.end) {
                seenInAnyRange = true
                break
            }
        }
        if !seenInAnyRange {
            fmt.Printf("Codepoint 0x%x is not in any range!\n", c.codepoint)
            must.Never()
        }
    }

    {
        dest := must.Result(os.Create("../../../../text/ccc/ccc.bin"))
        cccEncode(dest, ranges)
        defer dest.Close()
    }
}

// rng is a range of entries sharing the same ccc
type rng struct {
    ccc uint8
    start rune
    end rune
}

// Pack encodes a range as a character's codepoint (21 bits) (the start),
// length of the range (currently max 256) (11 bits), and ccc (8 bits),
// for a total of 40 bits, which fits in 5 bytes
func (r rng) Pack() uint64 {
    cp  := r.start
    rl  := int(r.end) - int(r.start)
    ccc := r.ccc

    must.True(rl     <=    0x7FF)
    must.True(cp     <= 0x1FFFFF)

    var x uint64
    x |= (uint64(cp)  & 0x1FFFFF)
    x |= (uint64(rl)  &    0x7FF) << 21
    x |= (uint64(ccc) &     0xFF) << (21 + 11)
    return x
}

func cccRanges(chars []Char) []rng {
    ranges := make([]rng, 0)

    lastSeenCCC := int(-1)
    lastIdx := rune(0)

    boundary := func(c Char) {
        ccc := c.ccc
        if lastSeenCCC >= 0 {
            next := c.codepoint
            r := rng{
                ccc:   uint8(lastSeenCCC),
                start: lastIdx,
                end:   next,
            }
            ranges = append(ranges, r)
            lastIdx = next
        }
        lastSeenCCC = int(ccc)
    }

    for _, c := range chars {
        if lastSeenCCC != int(c.ccc) {
            boundary(c)
        }
    }

    // finish
    last := chars[len(chars)-1]
    ranges = append(ranges, rng{
        ccc: last.ccc,
        start: lastIdx,
        end: 0x10FFFF,
    })

    return ranges
}

func cccEncode(dest io.Writer, ranges []rng) {
    for _, r := range ranges {
        if r.ccc == 0 { continue }

        var x = r.Pack()
        b0 :=       (x) & 0xFF
        b1 := (x >>  8) & 0xFF
        b2 := (x >> 16) & 0xFF
        b3 := (x >> 24) & 0xFF
        b4 := (x >> 32) & 0xFF
        must.Result(dest.Write([]byte{
            uint8(b0),
            uint8(b1),
            uint8(b2),
            uint8(b3),
            uint8(b4),
        }))
    }
}
