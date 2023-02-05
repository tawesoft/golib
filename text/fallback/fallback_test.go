package fallback_test

import (
    "fmt"
    "strings"

    lazy "github.com/tawesoft/golib/v2/iter"
    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/text/fallback"
    "golang.org/x/text/unicode/runenames"
)

func ExampleIs() {
    type row struct {
        input rune
        alternative string
    }

    rows := []row{
        {'㎦', "㎞³"},
        {'㎦', "km³"},
        {'㎦', "km3"},
        {'㎦', "foo"},
        {'²',  "2"},
        {'½',  "1⁄2"}, // 0x2044
        {'½',  " 1/2"}, // 0x002F
    }

    for _, r := range rows {
        q := fallback.Is(r.input, r.alternative)
        fmt.Printf("Is %s a fallback for %c? %t\n",
            r.alternative, r.input, q)
    }

    // output:
    // Is ㎞³ a fallback for ㎦? true
    // Is km³ a fallback for ㎦? true
    // Is km3 a fallback for ㎦? true
    // Is foo a fallback for ㎦? false
    // Is 2 a fallback for ²? true
    // Is 1⁄2 a fallback for ½? true
    // Is  1/2 a fallback for ½? true
}

func ExampleSubs() {
    rows := []rune{
        '㎦', '²', '½',
    }

    for _, r := range rows {
        fmt.Printf("=== %c ===\n", r)
        for _, s := range fallback.Subs(r) {
            fmt.Println(s)
        }
    }

    // output:
    // === ㎦ ===
    // ㎦
    // km³
    // km3
    // === ² ===
    // ²
    // 2
    // === ½ ===
    // ½
    //  1/2
    // 1⁄2
}

func ExampleEquivalent() {
    input := "\u0041\u030A\u0064\u0307\u0327"
    fmt.Printf("Input: %s %x %x\n", input, []rune(input), []byte(input))
    eq := must.Result(fallback.Equivalent(input))

    lazy.Walk(func (x string) {
        fmt.Printf("%s: %x = %s\n", x, []rune(x),
            lazy.Join(lazy.StringJoiner(", "),
                lazy.Map(strings.ToLower,
                    lazy.Map[rune, string](runenames.Name,
                        lazy.FromString(x)))))
    }, eq)

    /* (Found in the Unicode ICU as a test case)

    Results for: {LATIN CAPITAL LETTER A WITH RING ABOVE}{LATIN SMALL LETTER D}{COMBINING DOT ABOVE}{COMBINING CEDILLA}

    1: \u0041\u030A\u0064\u0307\u0327
     = {LATIN CAPITAL LETTER A}{COMBINING RING ABOVE}{LATIN SMALL LETTER D}{COMBINING DOT ABOVE}{COMBINING CEDILLA}
    2: \u0041\u030A\u0064\u0327\u0307
     = {LATIN CAPITAL LETTER A}{COMBINING RING ABOVE}{LATIN SMALL LETTER D}{COMBINING CEDILLA}{COMBINING DOT ABOVE}
    3: \u0041\u030A\u1E0B\u0327
     = {LATIN CAPITAL LETTER A}{COMBINING RING ABOVE}{LATIN SMALL LETTER D WITH DOT ABOVE}{COMBINING CEDILLA}
    4: \u0041\u030A\u1E11\u0307
     = {LATIN CAPITAL LETTER A}{COMBINING RING ABOVE}{LATIN SMALL LETTER D WITH CEDILLA}{COMBINING DOT ABOVE}
    5: \u00C5\u0064\u0307\u0327
     = {LATIN CAPITAL LETTER A WITH RING ABOVE}{LATIN SMALL LETTER D}{COMBINING DOT ABOVE}{COMBINING CEDILLA}
    6: \u00C5\u0064\u0327\u0307
     = {LATIN CAPITAL LETTER A WITH RING ABOVE}{LATIN SMALL LETTER D}{COMBINING CEDILLA}{COMBINING DOT ABOVE}
    7: \u00C5\u1E0B\u0327
     = {LATIN CAPITAL LETTER A WITH RING ABOVE}{LATIN SMALL LETTER D WITH DOT ABOVE}{COMBINING CEDILLA}
    8: \u00C5\u1E11\u0307
     = {LATIN CAPITAL LETTER A WITH RING ABOVE}{LATIN SMALL LETTER D WITH CEDILLA}{COMBINING DOT ABOVE}
    9: \u212B\u0064\u0307\u0327
     = {ANGSTROM SIGN}{LATIN SMALL LETTER D}{COMBINING DOT ABOVE}{COMBINING CEDILLA}
    10: \u212B\u0064\u0327\u0307
     = {ANGSTROM SIGN}{LATIN SMALL LETTER D}{COMBINING CEDILLA}{COMBINING DOT ABOVE}
    11: \u212B\u1E0B\u0327
     = {ANGSTROM SIGN}{LATIN SMALL LETTER D WITH DOT ABOVE}{COMBINING CEDILLA}
    12: \u212B\u1E11\u0307
     = {ANGSTROM SIGN}{LATIN SMALL LETTER D WITH CEDILLA}{COMBINING DOT ABOVE}

     */

    // TODO for some reason our implementation is missing the two variants with an Angstrom Sign.
    //   This is probably due to Go's Unicode version being older than the example.
    //   Revisit once new Unicode versions land

    // output:
    // Input: Åḑ̇ [41 30a 64 307 327] 41cc8a64cc87cca7
    // Åḑ̇: [41 30a 64 327 307] = latin capital letter a, combining ring above, latin small letter d, combining cedilla, combining dot above
    // Åḑ̇: [c5 64 327 307] = latin capital letter a with ring above, latin small letter d, combining cedilla, combining dot above
    // Åḑ̇: [41 30a 1e0b 327] = latin capital letter a, combining ring above, latin small letter d with dot above, combining cedilla
    // Åḑ̇: [c5 1e0b 327] = latin capital letter a with ring above, latin small letter d with dot above, combining cedilla
    // Å̧: [41 30a 327] = latin capital letter a, combining ring above, combining cedilla
    // Å̧: [c5 327] = latin capital letter a with ring above, combining cedilla
    // Åḑ̇: [41 30a 1e11 307] = latin capital letter a, combining ring above, latin small letter d with cedilla, combining dot above
    // Åḑ̇: [c5 1e11 307] = latin capital letter a with ring above, latin small letter d with cedilla, combining dot above
    // Å̇: [41 30a 307] = latin capital letter a, combining ring above, combining dot above
    // Å̇: [c5 307] = latin capital letter a with ring above, combining dot above
}
