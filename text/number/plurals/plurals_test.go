package plurals

import (
    "fmt"
    "testing"

    "golang.org/x/text/feature/plural"
    "golang.org/x/text/language"
)

func ExampleCardinal() {
    englishRules := New(language.MustParse("en"))
    fish := func(i int) string {
        var format string
        form := englishRules.Cardinal(fmt.Sprintf("%d", i))
        switch form {
            case One: format = "%d fish is in your cart; do you want to buy it?"
            case Other: format = "%d fish are in your cart; do you want to buy them?"
        }
        return fmt.Sprintf(format, i)
    }
    fmt.Println(fish(0))
    fmt.Println(fish(1))
    fmt.Println(fish(2))

    welshRules := New(language.MustParse("cy"))
    catsAndDogs := func(cats int, dogs int) string {
        catMutations := map[Form]string{
            Zero: "cathod",
            One: "gath",
            Two: "gath",
            Few: "cath",
            Many: "chath",
            Other: "cath",
        }
        dogMutations := map[Form]string{
            Zero: "cŵn",
            One: "ci",
            Two: "gi",
            Few: "chi",
            Many: "chi",
            Other: "ci",
        }

        catForm := welshRules.Cardinal(fmt.Sprintf("%d", cats))
        dogForm := welshRules.Cardinal(fmt.Sprintf("%d", dogs))
        return fmt.Sprintf("Mae %d %s, %d %s.",
            dogs, dogMutations[dogForm],
            cats, catMutations[catForm],
        )
    }

    // example from https://cldr.unicode.org/index/cldr-spec/plural-rules
    fmt.Println(catsAndDogs(0, 0))
    fmt.Println(catsAndDogs(1, 1))
    fmt.Println(catsAndDogs(2, 2))
    fmt.Println(catsAndDogs(3, 3))
    fmt.Println(catsAndDogs(4, 4))
    fmt.Println(catsAndDogs(6, 6))

    // Output:
    // 0 fish are in your cart; do you want to buy them?
    // 1 fish is in your cart; do you want to buy it?
    // 2 fish are in your cart; do you want to buy them?
    // Mae 0 cŵn, 0 cathod.
    // Mae 1 ci, 1 gath.
    // Mae 2 gi, 2 gath.
    // Mae 3 chi, 3 cath.
    // Mae 4 ci, 4 cath.
    // Mae 6 chi, 6 chath.
}

func TestForm(t *testing.T) {
    /*
    From CLDR supplemental/ordinals.xml:
        <pluralRules locales="en">
            <pluralRule count="one">n % 10 = 1 and n % 100 != 11 @integer 1, 21, 31, 41, 51, 61, 71, 81, 101, 1001, …</pluralRule>
            <pluralRule count="two">n % 10 = 2 and n % 100 != 12 @integer 2, 22, 32, 42, 52, 62, 72, 82, 102, 1002, …</pluralRule>
            <pluralRule count="few">n % 10 = 3 and n % 100 != 13 @integer 3, 23, 33, 43, 53, 63, 73, 83, 103, 1003, …</pluralRule>
            <pluralRule count="other"> @integer 0, 4~18, 100, 1000, 10000, 100000, 1000000, …</pluralRule>
        </pluralRules>
        <pluralRules locales="cy">
            <pluralRule count="zero">n = 0,7,8,9 @integer 0, 7~9</pluralRule>
            <pluralRule count="one">n = 1 @integer 1</pluralRule>
            <pluralRule count="two">n = 2 @integer 2</pluralRule>
            <pluralRule count="few">n = 3,4 @integer 3, 4</pluralRule>
            <pluralRule count="many">n = 5,6 @integer 5, 6</pluralRule>
            <pluralRule count="other"> @integer 10~25, 100, 1000, 10000, 100000, 1000000, …</pluralRule>
        </pluralRules>

    From rbnf/en.xml:
        <rbnfrule value="0">=#,##0=$(ordinal,one{st}two{nd}few{rd}other{th})$;</rbnfrule>

    From CLDR supplemental/plurals.xml:
        <pluralRules locales="ast ca de en et fi fy gl ia io ji lij nl sc scn sv sw ur yi">
            <pluralRule count="one">i = 1 and v = 0 @integer 1</pluralRule>
            <pluralRule count="other"> @integer 0, 2~16, 100, 1000, 10000, 100000, 1000000, … @decimal 0.0~1.5, 10.0, 100.0, 1000.0, 10000.0, 100000.0, 1000000.0, …</pluralRule>
        </pluralRules>
    */
    rules := map[string]*plural.Rules{
        "plural.Ordinal": plural.Ordinal,
        "plural.Cardinal": plural.Cardinal,
    }
    type row struct {
        locale string
        rules string
        input string
        form Form
    }
    rows := []row{
        // locale, input, form
        {"en", "plural.Ordinal", "0", Other},
        {"en", "plural.Ordinal", "1", One},
        {"en", "plural.Ordinal", "2", Two},
        {"en", "plural.Ordinal", "3", Few},
        {"en", "plural.Ordinal", "4", Other},

        {"cy", "plural.Ordinal", "0", Zero},
        {"cy", "plural.Ordinal", "8", Zero},
        {"cy", "plural.Ordinal", "2", Two},
        {"cy", "plural.Ordinal", "3", Few},
        {"cy", "plural.Ordinal", "6", Many},
        {"cy", "plural.Ordinal", "12", Other},

        {"en", "plural.Cardinal", "0", Other},
        {"en", "plural.Cardinal", "1", One},
        {"en", "plural.Cardinal", "2", Other},
    }

    for _, test := range rows {
        f := form(language.MustParse(test.locale), rules[test.rules], test.input)
        if f != test.form {
            t.Errorf("test form(%q, %q, %q) failed: got form %d not %d",
                test.locale, test.rules, test.input, f, test.form)
        }
    }
}

func TestOperands(t_ *testing.T) {
    type row struct {
        input string
        i, v, w, f, t int
        ok bool
    }
    type result struct {
        i, v, w, f, t int
        ok bool
    }

    rows := []row{
        //                i  v  w    f   t  ok
        {"5.",            5, 0, 0,   0,  0, true},
        {".5",            0, 1, 1,   5,  5, true},
        // from https://unicode.org/reports/tr35/tr35-numbers.html#table-plural-operand-examples
        {"1",             1, 0, 0,   0,  0, true},
        {"1.0",           1, 1, 0,   0,  0, true},
        {"1.00",          1, 2, 0,   0,  0, true},
        {"1.3",           1, 1, 1,   3,  3, true},
        {"1.03",          1, 2, 2,   3,  3, true},
        {"1.230",         1, 3, 2, 230, 23, true},
        {"1200000", 1200000, 0, 0,   0,  0, true},
    }

    for _, test := range rows {
        i, v, w, f, t, ok := operands(test.input)
        got := result{i, v, w, f, t, ok}
        wanted := result{test.i, test.v, test.w, test.f, test.t, test.ok}

        if (got.ok != test.ok) {
            t_.Errorf("test %+v failed: expected %t but got %+v", test, test.ok, got)
        } else if (wanted != got) {
            t_.Errorf("test %+v failed: got %+v", test, got)
        }
    }
}
