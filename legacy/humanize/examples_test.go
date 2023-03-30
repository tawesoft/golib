package humanize_test

import (
    "fmt"
    "time"

    humanizex "github.com/tawesoft/golib/v2/legacy/humanize"
    "golang.org/x/text/language"
)

func Example_simple() {
    mustInt64 := func(v int64, err error) int64 {
        if err != nil { panic(err) }
        return v
    }

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

    // Output:
    // 1.5 KiB
    // 1,5 KiB
    // ১.৫ KiB
    // 1536
}

func Example_customFactors() {
    factors := humanizex.Factors{
        Factors:    []humanizex.Factor{
            {1,                         humanizex.Unit{"millicenton", "millicenton"}, humanizex.FactorModeReplace},
            {60,                        humanizex.Unit{"centon",      "centon"},      humanizex.FactorModeReplace},
            {60 * 60,                   humanizex.Unit{"centar",      "centar"},      humanizex.FactorModeReplace},
            {24 * 60 * 60,              humanizex.Unit{"cycle",       "cycle"},       humanizex.FactorModeReplace},
            {7 * 24 * 60 * 60,          humanizex.Unit{"secton",      "secton"},      humanizex.FactorModeReplace},
            {28 * 24 * 60 * 60,         humanizex.Unit{"sectar",      "sectar"},      humanizex.FactorModeReplace},
            {365 * 24 * 60 * 60,        humanizex.Unit{"yahren",      "yahren"},      humanizex.FactorModeReplace},
            {100 * 365 * 24 * 60 * 60,  humanizex.Unit{"centauron",   "centauron"},   humanizex.FactorModeReplace},
        },
        Components: 2,
    }

    h := humanizex.NewHumanizer(language.English)

    est := float64((2 * 365 * 24 * 60 * 60) + 1)

    fmt.Printf("Hey, I'll be with you in %s. Watch out for toasters!\n",
        h.Format(est, humanizex.Unit{"millicenton", "millicenton"}, factors).Utf8)

    // Output:
    // Hey, I'll be with you in 2 yahren 1 millicenton. Watch out for toasters!
}

func Example_customDurations() {
    plural := func (x float64) string {
        if x > 0.99 && x < 1.01 { return "" }
        return "s"
    }

    duration := (2 * time.Hour) + (20 * time.Second)

    // prints "Basic time: 2 h 20 s"
    fmt.Printf("Basic time: %s\n", humanizex.NewHumanizer(language.English).FormatDuration(duration))

    // Get the raw format parts
    parts := humanizex.FormatParts(
        duration.Seconds(),
        humanizex.CommonUnits.Second,
        humanizex.CommonFactors.Time,
    )

    // prints "Nice time: 2 hours and 20 seconds ago"
    fmt.Printf("Nice time: ")
    if (len(parts) == 1) && (parts[0].Unit.Utf8 == "s") {
        fmt.Printf("just now\n")
    } else {
        for i, part := range parts {
            fmt.Printf("%d", int(part.Magnitude + 0.5))

            if part.Unit.Utf8 == "y" {
                fmt.Printf(" year%s", plural(part.Magnitude))
            } else if part.Unit.Utf8 == "d" {
                fmt.Printf(" day%s", plural(part.Magnitude))
            } else if part.Unit.Utf8 == "h" {
                fmt.Printf(" hour%s", plural(part.Magnitude))
            } else if part.Unit.Utf8 == "min" {
                fmt.Printf(" minute%s", plural(part.Magnitude))
            } else if part.Unit.Utf8 == "s" {
                fmt.Printf(" second%s", plural(part.Magnitude))
            }

            if i + 1 < len(parts) {
                fmt.Printf(" and ")
            } else {
                fmt.Printf(" ago\n")
            }
        }
    }

    // Output:
    // Basic time: 2 h 20 s
    // Nice time: 2 hours and 20 seconds ago
}
