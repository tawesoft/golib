package symbols_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/text/number/symbols"
)

func ExampleGet_() {
    print := func(lang, region, numbersystem string, desc string) {
        s := symbols.Get_(lang, "", region, "", numbersystem)

        fmt.Printf("For %s:\n", desc)
        fmt.Printf("> Decimal point: %q\n", s.Decimal)
        fmt.Printf("> Approximately: %q\n", s.ApproximatelySign)
        fmt.Printf("> Infinity: %q\n", s.Infinity)
    }

    print("en", "gb", "",     "en-gb")
    print("fr", "",   "",     "fr")
    print("en", "",   "arab", "en-u-nu-arab")

    // Output:
    // For en-gb:
    // > Decimal point: "."
    // > Approximately: "~"
    // > Infinity: "∞"
    // For fr:
    // > Decimal point: ","
    // > Approximately: "≃"
    // > Infinity: "∞"
    // For en-u-nu-arab:
    // > Decimal point: "٫"
    // > Approximately: "~"
    // > Infinity: "∞"
}
