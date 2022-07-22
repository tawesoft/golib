package human_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/human"
    "golang.org/x/text/language"
)

func ExampleParseFloat() {
    tag := language.Dutch
    dp := human.NewDecimalParser(tag)
    v, _ := human.ParseFloat(dp, "1.234.567,89")
    fmt.Printf("%.2f", v)

    // Output:
    // 1234567.89
}

func ExampleParseInt() {
    tag := language.Dutch
    dp := human.NewDecimalParser(tag)
    v, _ := human.ParseInt(dp, "1.234.567")
    fmt.Printf("%d", v)

    // Output:
    // 1234567
}

func ExampleHumanizer_ParseDecimalFloat() {
    h := human.New(language.Dutch)
    v, _ := h.ParseDecimalFloat("1.234.567.89")
    fmt.Printf("%.2f\n", v)

    // Output:
    // 1234.567.89
}

func ExampleHumanizer_ParseDecimalInt() {
    h := human.New(language.Dutch)
    v, _ := h.ParseDecimalInt("1.234.567")
    fmt.Printf("%d", v)

    // Output:
    // 1234567
}
