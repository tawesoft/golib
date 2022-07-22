// Command gen.go generates decomposition mappings from Unicode data.
//
// Unlike the package /x/text/unicode/norm, this also includes the
// [compatibility mapping tags] extracted from the
// [XML decomposition properties].
//
// [compatibility mapping tags]: https://unicode.org/reports/tr44/#Formatting_Tags_Table
// [XML decomposition properties]: https://www.unicode.org/reports/tr42/#d1e2932
package main

import (
    "archive/zip"
    "bufio"
    "compress/gzip"
    "encoding/xml"
    "fmt"
    "io"
    "os"
    "sort"
    "strconv"
    "strings"

    "github.com/tawesoft/golib/v2/ks"
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

// Pack encodes a unicode character's codepoint (21 bits), decomposition
// type (dt) (5 bits), length of the decomposition mapping (dml) (5 bits),
// an index dmi that points into table decomposition table dms (16 bits),
// for a total of 47 bits, which fits in 6 bytes with 1 bit of padding.
func (c Char) Pack() uint64 {
    cp := c.codepoint
    dt := c.DecompositionType()
    dmi := c.Offset
    dml := len(c.DecompositionMappings())

    ks.Assert(dt     <=     0x1F)
    ks.Assert(dml    <=     0x1F)
    ks.Assert(dmi    <=   0xFFFF)
    ks.Assert(cp     <= 0x1FFFFF)

    var x uint64
    x |= (uint64(cp)  & 0x1FFFFF)
    x |= (uint64(dt)  &     0x1F) << 21
    x |= (uint64(dml) &     0x1F) << (21 + 5)
    x |= (uint64(dmi) &   0xFFFF) << (21 + 5 + 5)
    return x
}

func ParseCodepoint(x string) rune {
    return rune(ks.Must(strconv.ParseInt(x, 16, 32)))
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

// DecompositionMappingsLiteral formats DecompositionMappings so that
// the output looks like "0x1, 0x2, 0x3"
func (c Char) DecompositionMappingsLiteral() string {
    var sb strings.Builder
    rs := c.DecompositionMappings()
    for i, r := range rs {
        sb.WriteString(fmt.Sprintf("0x%X", r))
        if i + 1 < len(rs) {
            sb.WriteRune(',')
        }
    }
    return sb.String()
}


type DT int8
const (
    DTNone DT =  0
    DTCan  DT =  1 // Canonical
    DTCom  DT =  2 // Otherwise unspecified compatibility character
    DTEnc  DT =  3 // Encircled form
    DTFin  DT =  4 // Final presentation form (Arabic)
    DTFont DT =  5 // Font variant (for example, a blackletter form)
    DTFra  DT =  6 // Vulgar fraction form
    DTInit DT =  7 // Initial presentation form (Arabic)
    DTIso  DT =  8 // Isolated presentation form (Arabic)
    DTMed  DT =  9 // Medial presentation form (Arabic)
    DTNar  DT = 10 // Narrow (or hankaku) compatibility character
    DTNb   DT = 11 // No-break version of a space or hyphen
    DTSml  DT = 12 // Small variant form (CNS compatibility)
    DTSqr  DT = 13 // CJK squared font variant
    DTSub  DT = 14 // Subscript form
    DTSup  DT = 15 // Superscript form
    DTVert DT = 16 // Vertical layout presentation form
    DTWide DT = 17 // Wide (or zenkaku) compatibility character
) // fits in a 5 bit number

func (c Char) DecompositionType() DT {
/*
    XML stores this as

    attribute dt { "can"  | "com" | "enc" | "fin"  | "font" | "fra"
                 | "init" | "iso" | "med" | "nar"  | "nb"   | "sml"
                 | "sqr"  | "sub" | "sup" | "vert" | "wide" | "none"}?
*/
    switch c.decompositionType {
        case "none": return DTNone
        case "can":  return DTCan
        case "com":  return DTCom
        case "enc":  return DTEnc
        case "fin":  return DTFin
        case "font": return DTFont
        case "fra":  return DTFra
        case "init": return DTInit
        case "iso":  return DTIso
        case "med":  return DTMed
        case "nar":  return DTNar
        case "nb":   return DTNb
        case "sml":  return DTSml
        case "sqr":  return DTSqr
        case "sup":  return DTSup
        case "sub":  return DTSub
        case "vert": return DTVert
        case "wide": return DTWide
        default:
            fmt.Println(c.decompositionType)
            ks.Never()
    }

    return DTNone
}

func main() {
    zr := ks.Must(zip.OpenReader("../../DATA/ucd.nounihan.grouped.13.0.0.zip"))

    opener := func(name string) func() (io.ReadCloser, error) {
        return func() (io.ReadCloser, error) {
            return zr.Open(name)
        }
    }

    chars := make([]Char, 0)

    ks.Check(ks.WithCloser(opener("ucd.nounihan.grouped.xml"), func(f io.ReadCloser) error {
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
                            if inRepertoire { ks.Never() }
                            inRepertoire = true
                        case "group":
                            if !inRepertoire { break }
                            if inGroup { ks.Never() }
                            inGroup = true
                            group = CharFromAttrs(ty.Attr, ks.Zero[Char]())
                        case "char":
                            if !inRepertoire { break }
                            if !inGroup { ks.Never() }
                            c := CharFromAttrs(ty.Attr, group)
                            if c.DecompositionType() == 0 { break }
                            if c.IsRange() { ks.Never() }
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

    ks.Assert(sort.SliceIsSorted(chars, func(i int, j int) bool {
        return chars[i].codepoint < chars[j].codepoint
    }))

    offset := 0
    for i := 0; i < len(chars); i++ {
        chars[i].Offset = offset
        offset += len(chars[i].DecompositionMappings())
    }

    {
        dest := ks.Must(os.Create("../../../../text/dm/dti.bin"))
        defer dest.Close()
        dtiEncode(dest, chars)
    }
    {
        dest := ks.Must(os.Create("../../../../text/dm/dms.bin.gz"))
        defer dest.Close()

        zdest := gzip.NewWriter(dest)
        defer zdest.Close()

        dmsEncode(zdest, chars)
    }
}

func dtiEncode(dest io.Writer, chars []Char) {
    for _, c := range chars {
        p := c.Pack() // 48 bit
        b0 :=       (p) & 0xFF
        b1 := (p >>  8) & 0xFF
        b2 := (p >> 16) & 0xFF
        b3 := (p >> 24) & 0xFF
        b4 := (p >> 32) & 0xFF
        b5 := (p >> 40) & 0xFF
        ks.Must(dest.Write([]byte{
            uint8(b0),
            uint8(b1),
            uint8(b2),
            uint8(b3),
            uint8(b4),
            uint8(b5),
        }))
    }
}

func dmsEncode(dest io.Writer, chars []Char) {
    for _, c := range chars {
        for _, m := range c.DecompositionMappings() {
            x := uint32(m)
            b0 :=       (x) & 0xFF
            b1 := (x >>  8) & 0xFF
            b2 := (x >> 16) & 0xFF
            ks.Must(dest.Write([]byte{
                uint8(b0),
                uint8(b1),
                uint8(b2),
            }))
        }
    }
}
