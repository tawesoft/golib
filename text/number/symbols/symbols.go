// Package symbols contains the [CLDR number symbols], e.g. decimal point and
// the minus sign, for different locales.
//
// [CLDR number symbols]: https://cldr.unicode.org/translation/number-currency-formats/number-symbols
package symbols

import (
    "encoding/binary"
    "fmt"
    "strings"
)

type Symbols struct {
    // Decimal separates the integer and fractional part of a number.
    Decimal string

    // Group separates clusters of integer digits to make large numbers more
    // legible; commonly used for thousands.
    Group string

    // List separates numbers in a list intended to represent structured data such
    // as an array; must be different from the decimal value. This list separator
    // is for “non-linguistic” usage as opposed to the “linguistic” lists (e.g.
    // “Bob, Carol, and Ted”).
    List string

    // %, +, -
    PercentSign, PlusSign, MinusSign string

    // ApproximatelySign is used to denote a value that is approximate but not
    // exact.
    ApproximatelySign string

    // Exponential is used to separate mantissa and exponent values.
    Exponential string

    // SuperscriptingExponent is used in exponential notation like '×' in "1.23 ×
    // 10^4".
    SuperscriptingExponent string

    // Infinity is an infinity symbol
    Infinity string

    // PerMille is used to indicate a per-mille (1/1000th) amount.
    PerMille string

    // NaN represents "Not a number".
    NaN string

    // CurrencyDecimal is optional. If specified, then for currency
    // formatting/parsing this is used as the decimal separator instead of using
    // the regular decimal separator; otherwise, the regular decimal separator is
    // used.
    CurrencyDecimal string

    // CurrencyGroup is optional. If specified, then for currency formatting/parsing
    // this is used as the group separator instead of using the regular group
    // separator; otherwise, the regular group separator is used.
    CurrencyGroup string
}

func loadSymbol(idx int) Symbols {
    data := symbolsdata[idx * symbolsRowSize : (idx + 1) * symbolsRowSize]
    str := func(x int) string {
        x *= 3
        left := int(binary.LittleEndian.Uint16(data[x:x+2]))
        right := left + int(data[x+2])
        return stringdata[left:right]
    }
    return Symbols{
        Decimal:                str(0),
        Group:                  str(1),
        List:                   str(2),
        PercentSign:            str(3),
        PlusSign:               str(4),
        MinusSign:              str(5),
        ApproximatelySign:      str(6),
        Exponential:            str(7),
        SuperscriptingExponent: str(8),
        Infinity:               str(9),
        PerMille:               str(10),
        NaN:                    str(11),
        CurrencyDecimal:        str(12),
        CurrencyGroup:          str(13),
    }
}

func (s Symbols) update(n Symbols) Symbols {
    f := func(a string, b string) string {
        if len(b) != 0 { return b }
        return a
    }
    return Symbols{
        Decimal:                f(s.Decimal,                n.Decimal),
        Group:                  f(s.Group,                  n.Group),
        List:                   f(s.List,                   n.List),
        PercentSign:            f(s.PercentSign,            n.PercentSign),
        PlusSign:               f(s.PlusSign,               n.PlusSign),
        MinusSign:              f(s.MinusSign,              n.MinusSign),
        ApproximatelySign:      f(s.ApproximatelySign,      n.ApproximatelySign),
        Exponential:            f(s.Exponential,            n.Exponential),
        SuperscriptingExponent: f(s.SuperscriptingExponent, n.SuperscriptingExponent),
        Infinity:               f(s.Infinity,               n.Infinity),
        PerMille:               f(s.PerMille,               n.PerMille),
        NaN:                    f(s.NaN,                    n.NaN),
        CurrencyDecimal:        f(s.CurrencyDecimal,        n.CurrencyDecimal),
        CurrencyGroup:          f(s.CurrencyGroup,          n.CurrencyGroup),
    }
}

// Get_ looks up the Symbols used to format numbers in a given locale.
//
// Use empty strings to omit any argument.
//
// Get_ will be replaced by Get(language.Tag) in future versions of Go.
// Currently, this is blocked by issue [53872] (which has a fix waiting to be
// merged), so you need to specify the arguments manually.
//
// [53872]: https://github.com/golang/go/issues/53872
func Get_(language string, script string, region string, variant string, numberingSystem string) Symbols {
    if language == "" { language = "root" }
    if numberingSystem == "" { numberingSystem = "latn" }

    language = strings.ToLower(language)
    script = strings.ToLower(script)
    region = strings.ToLower(region)
    variant = strings.ToLower(variant)
    numberingSystem = strings.ToLower(numberingSystem)

    makeTag := func(l, s, r, v, n string) string {
        return fmt.Sprintf("%s-%s-%s-%s/%s", l, s, r, v, n)
    }

    var sym Symbols

    update := func(l, s, r, v, n string) {
        t := makeTag(l, s, r, v, n)
        if idx, ok := indexes[t]; ok {
            sym = sym.update(loadSymbol(idx))
        }
    }

    update("root",   "",     "",     "",      "latn")
    update("root",   "",     "",     "",      numberingSystem)
    update(language, "",     "",     "",      numberingSystem)
    update(language, "",     "",     variant, numberingSystem)
    update(language, "",     region, "",      numberingSystem)
    update(language, "",     region, variant, numberingSystem)
    update(language, script, region, "",      numberingSystem)
    update(language, script, region, variant, numberingSystem)

    return sym
}
