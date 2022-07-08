package localize_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/ks"
    "github.com/tawesoft/golib/v2/localize"
    "golang.org/x/text/language"
    "golang.org/x/text/language/display"
    "golang.org/x/text/message"
    "golang.org/x/text/number"
)

// Demonstrates converting a numbers to and from a strings in a given locale,
// with example outputs for a variety of locales.
func Example_ManyNumbers() {
    const input = 123456789.012

    langs := []language.Tag{
        language.Arabic,
        language.BritishEnglish,
        language.Dutch,
        language.French,
        language.Malayalam,
        language.Bengali,
    }

    // language name as a string
    namer := display.Tags(language.English)

    for _, t := range langs {
        printer := message.NewPrinter(t)
        localized := printer.Sprintf("%.4f", number.Decimal(input))

        parser := localize.NewDecimalFormat(t)
        result := ks.Must(parser.ParseFloat(localized))

        fmt.Printf("Language: %s\nPrints %T %.4f as %q\nParses %q back to %T %.4f\n\n",
            namer.Name(t), input, input, localized, localized, result, result)
    }

    // Output:
    // Language: Arabic
    // Prints float64 123456789.0120 as "١٢٣٬٤٥٦٬٧٨٩٫٠١٢٠"
    // Parses "١٢٣٬٤٥٦٬٧٨٩٫٠١٢٠" back to float64 123456789.0120
    //
    // Language: British English
    // Prints float64 123456789.0120 as "123,456,789.0120"
    // Parses "123,456,789.0120" back to float64 123456789.0120
    //
    // Language: Dutch
    // Prints float64 123456789.0120 as "123.456.789,0120"
    // Parses "123.456.789,0120" back to float64 123456789.0120
    //
    // Language: French
    // Prints float64 123456789.0120 as "123\u00a0456\u00a0789,0120"
    // Parses "123\u00a0456\u00a0789,0120" back to float64 123456789.0120
    //
    // Language: Malayalam
    // Prints float64 123456789.0120 as "12,34,56,789.0120"
    // Parses "12,34,56,789.0120" back to float64 123456789.0120
    //
    // Language: Bangla
    // Prints float64 123456789.0120 as "১২,৩৪,৫৬,৭৮৯.০১২০"
    // Parses "১২,৩৪,৫৬,৭৮৯.০১২০" back to float64 123456789.0120
}
