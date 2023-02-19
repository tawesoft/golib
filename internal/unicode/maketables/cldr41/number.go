package cldr41

import (
    "fmt"
    "io"
    "math"
    "os"
    "path"
    "reflect"
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

func MakeNumberSymbols(basedir string, dest string) {
    wr := must.Result(os.Create(dest))
    defer wr.Close()

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

    data := struct {
        SymbolsData []Symbols
        SortedKeys []string
        KeyedIndexes map[string]int
        String string
    }{
        SymbolsData: symbolData,
        SortedKeys: sortedKeys,
        KeyedIndexes: keyedIndexes,
        String: fmt.Sprintf("%q", strings.Join(sortedSeen, "")),
    }

    getField := func(s Symbols, field string) Slice {
        r := reflect.ValueOf(s)
        f := reflect.Indirect(r).FieldByName(field)
        return f.Interface().(Slice)
    }

    fm := map[string]any{
        "Field": func(name string, data Symbols) string {
            slice := getField(data, name)
            if slice.right == 0 {
                return ""
            }
            return fmt.Sprintf("%s: slice{%d, %d},", name, slice.left, slice.right)
        },
    }

    must.Check(must.Result(template.New("").Funcs(fm).Parse(strings.TrimSpace(`
// Code generated by internal/unicode/gen.sh - DO NOT EDIT.

package symbols

type slice struct{left, right uint16}
type symbols struct {
    Decimal                 slice
    Group                   slice
    List                    slice
    PercentSign             slice
    PlusSign                slice
    MinusSign               slice
    ApproximatelySign       slice
    Exponential             slice
    SuperscriptingExponent  slice
    Infinity                slice
    PerMille                slice
    NaN                     slice
    CurrencyDecimal         slice
    CurrencyGroup           slice
}

var stringdata = {{ .String }}

var indexes = map[string]int{
    {{range $k := .SortedKeys }}"{{$k}}": {{index $.KeyedIndexes $k}},
    {{end}}
}

var symbolsdata = []symbols{
    {{range $k := .SymbolsData }} {
        {{- Field "Decimal" $k -}}
        {{- Field "Group" $k -}}
        {{- Field "List" $k -}}
        {{- Field "PercentSign" $k -}}
        {{- Field "PlusSign" $k -}}
        {{- Field "MinusSign" $k -}}
        {{- Field "ApproximatelySign" $k -}}
        {{- Field "Exponential" $k -}}
        {{- Field "SuperscriptingExponent" $k -}}
        {{- Field "Infinity" $k -}}
        {{- Field "PerMille" $k -}}
        {{- Field "NaN" $k -}}
        {{- Field "CurrencyDecimal" $k -}}
        {{- Field "CurrencyGroup" $k -}} },
    {{end}}
}

    `))).Execute(wr, data))


}
