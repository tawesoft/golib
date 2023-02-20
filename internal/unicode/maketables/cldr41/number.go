package cldr41

import (
    "compress/gzip"
    "encoding/binary"
    "fmt"
    "io"
    "math"
    "os"
    "path"
    "sort"
    "strings"
    "text/template"

    "github.com/tawesoft/golib/v2/internal/unicode/ldml"
    "github.com/tawesoft/golib/v2/iter"
    "github.com/tawesoft/golib/v2/must"
)

func MakeNumberingSystemRules(basedir string, dest string) {
    wr := must.Result(os.Create(dest))
    defer wr.Close()

    r := must.Result(os.Open(path.Join(basedir, "cldr-41.0/common/rbnf/root.xml")))
    defer r.Close()

    gdoc := must.Result(ldml.Parse(must.Result(io.ReadAll(r))))
    must.Equal(gdoc.Type, ldml.DocTypeLdml)
    doc := gdoc.Ldml

    for _, group := range doc.RbnfRulesetGroupings {
        if group.Type != "NumberingSystemRules" { continue }
        for _, set := range group.Rulesets {
            if strings.EqualFold(set.Access, "private") {
                must.Result(fmt.Fprintf(wr, "%%%%%s:\n", set.Type))
            } else {
                must.Result(fmt.Fprintf(wr, "%%%s:\n", set.Type))
            }
            for _, rule := range set.Rules {
                must.Result(fmt.Fprintf(wr, "\t%s\n", rule.IcuStyle()))
            }
            must.Result(fmt.Fprintf(wr, "\n"))
        }
    }
}

func MakeNumberSymbols(basedir string, destdir string) {
    type Slice struct {
        left, right uint16
    }

    type Symbols struct {
        Decimal                 Slice
        Group                   Slice
        List                    Slice
        PercentSign             Slice
        PlusSign                Slice
        MinusSign               Slice
        ApproximatelySign       Slice
        Exponential             Slice
        SuperscriptingExponent  Slice
        PerMille                Slice
        Infinity                Slice
        NaN                     Slice
        CurrencyDecimal         Slice
        CurrencyGroup           Slice
    }

    dir := must.Result(os.Open(path.Join(basedir, "cldr-41.0/common/main/")))
    defer dir.Close()
    names := must.Result(dir.Readdirnames(0))

    walk := func(f func(d ldml.Ldml)) {
        for _, name := range names {
            func() {
                r := must.Result(os.Open(path.Join(basedir, "cldr-41.0/common/main/"+name)))
                defer r.Close()
                gdoc := must.Result(ldml.Parse(must.Result(io.ReadAll(r))))
                must.Equal(gdoc.Type, ldml.DocTypeLdml)
                doc := gdoc.Ldml
                f(doc)
            }()
        }
    }

    // pass 1: unique strings

    // stores a slice{left, right} where left is the cumulative total of
    // character count of preceding sorted strings and right is left + length
    // of the current sorted string.
    seen := make(map[string]Slice)

    // stores language-script-region-variant/numbersystem,
    // and increments by one each time
    keyedIndexes := make(map[string]int)

    skip := func(s ldml.Symbols) bool {
        if strings.HasSuffix(s.Alias.Path, "[@numberSystem='latn']") { return true }
        if ((s.Decimal == "") &&
            (s.Group == "") &&
            (s.List == "") &&
            (s.PercentSign == "") &&
            (s.PlusSign == "") &&
            (s.MinusSign == "") &&
            (s.ApproximatelySign == "") &&
            (s.Exponential == "") &&
            (s.SuperscriptingExponent == "") &&
            (s.Infinity == "") &&
            (s.PerMille == "") &&
            (s.NaN == "") &&
            (s.CurrencyGroup == "") &&
            (s.CurrencyDecimal == "")) { return true}
        return false
    }

    walk(func(doc ldml.Ldml) {
        for _, symbols := range doc.Numbers.Symbols {
            if skip(symbols) { continue }

            if symbols.Alias.Path != "" {
                fmt.Printf("got unexpected alias: %q", symbols.Alias.Path)
            }

            markSeen := func(x string) { if len(x) > 0 { seen[x] = Slice{} } }

            markSeen(symbols.Decimal)
            markSeen(symbols.Group)
            markSeen(symbols.List)
            markSeen(symbols.PercentSign)
            markSeen(symbols.PlusSign)
            markSeen(symbols.MinusSign)
            markSeen(symbols.ApproximatelySign)
            markSeen(symbols.Exponential)
            markSeen(symbols.SuperscriptingExponent)
            markSeen(symbols.Infinity)
            markSeen(symbols.PerMille)
            markSeen(symbols.NaN)
            markSeen(symbols.CurrencyGroup)
            markSeen(symbols.CurrencyDecimal)
        }
    })

    // sort and count offsets for efficient lookup
    sortedSeen := iter.ToSlice(iter.Keys(iter.FromMap(seen)))
    sort.Strings(sortedSeen)
    offset := 0
    for _, s := range sortedSeen {
        must.True(len(s) + offset <= math.MaxUint16, "offset out of range")
        seen[s] = Slice{uint16(offset), uint16(offset + len(s))}
        offset += len(s)
    }

    // pass 2: collect data
    i := 0
    symbolData := make([]Symbols, 0)
    seenSymbols := make(map[Symbols]int)

    walk(func(doc ldml.Ldml) {
        for _, symbols := range doc.Numbers.Symbols {
            if skip(symbols) { continue }
            get := func(x string) Slice { return seen[x] }

            key := fmt.Sprintf("%s-%s-%s-%s/%s",
                strings.ToLower(doc.Language.String()),
                strings.ToLower(doc.Script.String()),
                strings.ToLower(doc.Region.String()),
                strings.ToLower(doc.Variant.String()),
                strings.ToLower(symbols.NumberSystem))

            if _, exists := keyedIndexes[key]; exists {
                must.Never("duplicate key %q", key)
            }

            s := Symbols{
                    Decimal:                get(symbols.Decimal),
                    Group:                  get(symbols.Group),
                    List:                   get(symbols.List),
                    PercentSign:            get(symbols.PercentSign),
                    PlusSign:               get(symbols.PlusSign),
                    MinusSign:              get(symbols.MinusSign),
                    ApproximatelySign:      get(symbols.ApproximatelySign),
                    Exponential:            get(symbols.Exponential),
                    SuperscriptingExponent: get(symbols.SuperscriptingExponent),
                    Infinity:               get(symbols.Infinity),
                    PerMille:               get(symbols.PerMille),
                    NaN:                    get(symbols.NaN),
                    CurrencyDecimal:        get(symbols.CurrencyDecimal),
                    CurrencyGroup:          get(symbols.CurrencyGroup),
            }

            if id, exists := seenSymbols[s]; exists {
                keyedIndexes[key] = id
            } else {
                symbolData = append(symbolData, s)
                seenSymbols[s] = i
                keyedIndexes[key] = i
                i++
            }
        }
    })

    sortedKeys := iter.ToSlice(iter.Keys(iter.FromMap(keyedIndexes)))
    sort.Strings(sortedKeys)

    encodedSymbolData := make([]uint8, 0)
    putSlice := func(s Slice) {
        var buf [3]byte
        length := int(s.right) - int(s.left)
        must.True(length < math.MaxUint8)
        binary.LittleEndian.PutUint16(buf[0:2], s.left)
        buf[2] = uint8(length)
        encodedSymbolData = append(encodedSymbolData, buf[0], buf[1], buf[2])
    }
    for _, s := range symbolData {
        putSlice(s.Decimal)
        putSlice(s.Group)
        putSlice(s.List)
        putSlice(s.PercentSign)
        putSlice(s.PlusSign)
        putSlice(s.MinusSign)
        putSlice(s.ApproximatelySign)
        putSlice(s.Exponential)
        putSlice(s.SuperscriptingExponent)
        putSlice(s.Infinity)
        putSlice(s.PerMille)
        putSlice(s.NaN)
        putSlice(s.CurrencyDecimal)
        putSlice(s.CurrencyGroup)
    }
    rowSize := 14*3
    must.True(rowSize * len(symbolData) == len(encodedSymbolData))

    data := struct {
        SortedKeys []string
        KeyedIndexes map[string]int
        String string
        RowSize int
    }{
        SortedKeys: sortedKeys,
        KeyedIndexes: keyedIndexes,
        String: fmt.Sprintf("%q", strings.Join(sortedSeen, "")),
        RowSize: rowSize,
    }

    fm := map[string]any{
        "strlen": func(x string) int { return len(x) },
    }

    {
        wr := must.Result(os.Create(path.Join(destdir, "text/number/symbols/tables.go")))
        defer wr.Close()
        must.Check(must.Result(template.New("").Funcs(fm).Parse(strings.TrimSpace(`
// Code generated by internal/unicode/gen.sh - DO NOT EDIT.

package symbols

import (
    "bytes"
    "compress/gzip"
    _ "embed"
    "io"
    "strconv"

    "github.com/tawesoft/golib/v2/must"
)

type slice struct{left, right uint16}

// {{strlen .String}} bytes
var stringdata = {{ .String }}

//go:embed indexes.gz
var indexesGz []byte

//go:embed symbols.gz
var symbolsGz []byte

func ungz(src []byte) []byte {
    rdr := must.Result(gzip.NewReader(bytes.NewReader(src)))
    defer rdr.Close()
    return must.Result(io.ReadAll(rdr))
}

var indexes = func() map[string]int {
    v := make(map[string]int)
    data := ungz(indexesGz)
    for {
        idx1 := bytes.IndexByte(data, byte(':'))
        if idx1 < 0 { break }

        idx2 := bytes.IndexByte(data, byte(';'))
        if idx2 < 0 { break }

        key := string(data[0:idx1])
        value := string(data[idx1+1:idx2])
        v[key] = must.Result(strconv.Atoi(value))
        data = data[idx2+1:]
    }
    return v
}()

var symbolsdata = ungz(symbolsGz)
const symbolsRowSize = {{.RowSize}}

    `))).Execute(wr, data))
    }

    {
        wr := must.Result(os.Create(path.Join(destdir, "text/number/symbols/indexes.gz")))
        defer wr.Close()
        gz := gzip.NewWriter(wr)
        defer gz.Close()
        must.Check(must.Result(template.New("").Funcs(fm).Parse(strings.TrimSpace(`
{{range $k := .SortedKeys }}{{$k}}:{{index $.KeyedIndexes $k}};{{end}}
        `))).Execute(gz, data))
    }

    {
        wr := must.Result(os.Create(path.Join(destdir, "text/number/symbols/symbols.gz")))
        defer wr.Close()
        gz := gzip.NewWriter(wr)
        defer gz.Close()
        must.Result(gz.Write(encodedSymbolData))
    }
}
