package localize_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/legacy/localize"
    "golang.org/x/text/language"
    "golang.org/x/text/message"
)

// Demonstrates converting a numbers to and from a strings in a given locale
func Example_SimpleNumber() {

    const million = 1_000_000
    english := language.BritishEnglish
    dutch := language.Dutch

    message.NewPrinter(english).Printf("A million in English locale is '%d'\n", million)
    message.NewPrinter(dutch).Printf("But a million in Dutch locale is '%d'\n", million)

    pEng := localize.NewDecimalFormat(english)
    i, err := pEng.ParseInt("1,234")
    if err != nil { panic("error parsing a number in English locale:"+err.Error()) }
    fmt.Printf("1,234 parsed in English locale %T %d\n", i, i)

    pDutch := localize.NewDecimalFormat(dutch)
    f, err := pDutch.ParseFloat("1,234")
    if err != nil { panic("error parsing a number in Dutch locale: "+err.Error()) }
    fmt.Printf("But 1,234 parsed in Dutch locale is %T %f\n", f, f)

    _, err = pDutch.ParseInt("1,234")
    if err != nil {
        fmt.Println("In fact, parsing as an integer in the Dutch locale gives an error!")
    }

    // Output:
    // A million in English locale is '1,000,000'
    // But a million in Dutch locale is '1.000.000'
    // 1,234 parsed in English locale int64 1234
    // But 1,234 parsed in Dutch locale is float64 1.234000
    // In fact, parsing as an integer in the Dutch locale gives an error!

}
