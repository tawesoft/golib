// Command gen.go generates a mapping from the first character of a
// decomposition to an ordered list of codepoints whose decomposition starts
// with that first character.
//
// The use-case is in Enumerating Equivalent Strings in Unicode
package main

import (
    "archive/zip"
    "bufio"
    "bytes"
    "encoding/binary"
    "encoding/xml"
    "fmt"
    "io"
    "os"
    "sort"
    "strconv"
    "strings"

    "github.com/tawesoft/golib/v2/internal/legacy"
    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/operator"
)

type Char struct {
    codepoint rune
    coderange [2]rune // First and Last codepoint
    decompositionType string
    decompositionMapping string
    Offset int // set later, index into table
}

func (c Char) Codepoint() rune {
    return c.codepoint
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
            case "dt":
                c.decompositionType = a.Value
            case "dm":
                c.decompositionMapping = a.Value
        }
    }
    return c
}

// DecompositionMappings returns the code points, excluding itself, of a
// decomposition mapping
func (c Char) DecompositionMappings() []rune {
    var xs []rune
    for _, d := range strings.Split(c.decompositionMapping, " ") {
        if d == "#" { continue }
        xs = append(xs, ParseCodepoint(d))
    }
    return xs
}

type DT int8
const (
    DTNone DT =  0
    DTCan  DT =  1 // Canonical
    DTCom  DT =  2 // Compatibility
)

func (c Char) DecompositionType() DT {
    switch c.decompositionType {
        case "none": return DTNone
        case "can":  return DTCan
        default:     return DTCom
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
                            if c.DecompositionType() == 0 { break }
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

    dmToCodepoints := make(map[rune][]rune)

    for _, c := range chars {
        dms := c.DecompositionMappings()
        if c.DecompositionType() != DTCan { continue }
        if len(dms) == 0 { continue }
        first := dms[0]

        dmToCodepoints[first] = append(dmToCodepoints[first], c.Codepoint())
    }

    var sortedKeys []rune
    for k, _ := range dmToCodepoints {
        sortedKeys = append(sortedKeys, k)
    }
    sort.Slice(sortedKeys, func(i, j int) bool {
        return sortedKeys[i] < sortedKeys[j]
    })



    {
        // dstarts = decomposition starters
        dest := must.Result(os.Create("../../../../text/fallback/dstarts.bin"))
        defer dest.Close()

        must.Check(binary.Write(dest, binary.LittleEndian, int32(len(sortedKeys))))

        // decomposition starters index
        dstartsEncode(dest, sortedKeys, dmToCodepoints)
    }
}

func dstartsEncode(dest io.Writer, keys []rune, mapping map[rune][]rune) {
    var offset int
    var bb bytes.Buffer

    encode := func(xs []rune) []byte {
        bb.Reset()
        for _, x := range xs {
            bb.WriteRune(x)
        }
        return bb.Bytes()
    }

    for _, k := range keys {
        values := mapping[k]
        must.True(sort.SliceIsSorted(values, func(i int, j int) bool {
            return values[i] < values[j]
        }))

        // 21 bit codepoint
        // 3 bits padding
        // 16 bit index
        // = 40 bits
        b := uint32(k)
        o := uint32(offset)
        lm := len(encode(mapping[k]))

        dest.Write([]byte{
            uint8((b      ) & 0xFF), //  8
            uint8((b >>  8) & 0xFF), // 16
            uint8((b >> 16) & 0x1F), // 21, 3 bits padding = 24
            uint8((o      ) & 0xFF), // 32
            uint8((o >>  8) & 0xFF), // 40
        })

        offset += lm
    }

    // it's possible to infer the length by subtracting from the offset of
    // the next entry...
    for _, k := range keys {
        dest.Write(encode(mapping[k]))
    }

    must.True(offset < 65536)
    fmt.Printf("max offset %d\n", offset)
    fmt.Printf("last key: 0x%x\n", keys[len(keys)-1])
}
