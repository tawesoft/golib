// Example formatting and parsing Byte quantities in various locales
package main

import (
    "fmt"

    humanizex "github.com/tawesoft/golib/v2/legacy/humanize"
    "golang.org/x/text/language"
)

func mustInt64(v int64, err error) int64 {
    if err != nil { panic(err) }
    return v
}

func main() {
    hEnglish := humanizex.NewHumanizer(language.English)
    hDanish  := humanizex.NewHumanizer(language.Danish)
    hBengali := humanizex.NewHumanizer(language.Bengali)

    // prints 1.5 KiB
    fmt.Println(hEnglish.FormatBytesIEC(1024 + 512))

    // prints 1,5 KiB
    fmt.Println(hDanish.FormatBytesIEC(1024 + 512))

    // prints ১.৫ KiB
    fmt.Println(hBengali.FormatBytesIEC(1024 + 512))

    // prints 1536
    fmt.Println(mustInt64(hEnglish.ParseBytesIEC("1.5 KiB")))
}
