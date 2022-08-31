package view_test

import (
    "fmt"
    "strings"

    "github.com/tawesoft/golib/v2/view"
)

func Example() {

    // original map of lowercase animals to sounds
    quacks := map[string]string{
        "cat": "meow",
        "dog": "woof",
        "cow": "moo",
        "duck": "quack",
        "duckling": "quack",
    }

    // define a dynamic mapping that only shows things that quack and also
    // shows names in title-case and emotes in uppercase with an exclamation
    // mark...
    onlyQuackers := func(k string, v string) bool {
        return v == "quack"
    }
    keyer := view.Key(strings.Title, strings.ToLower)
    valuer := view.Valuer[string, string]{
        To: func(x string) string {
            return strings.ToUpper(x) + "!"
        },
        From: func (x string) string {
            return strings.ToLower(x)[0:len(x)-1]
        },
    }
    titleCaseOnlyQuackers := view.FromMap(quacks, onlyQuackers, keyer, valuer)

    // Update the view, which also updates the original
    titleCaseOnlyQuackers.Set("Duck impersonator", "QUACK!")

    // Update the original, which also updates the view
    quacks["another duckling"] = "quack"

    // Iterate over the view
    fmt.Printf("View:\n")
    it := titleCaseOnlyQuackers.Iter()
    for {
        quacker, ok := it.Next()
        if !ok { break }
        fmt.Printf("%s: %s\n", quacker.Key, quacker.Value)
    }

    fmt.Printf("\nOriginal (modified):\n")
    for k, v := range quacks {
        fmt.Printf("%s: %s\n", k, v)
    }

    /*
    Prints something like:

    View:
    Duckling: QUACK!
    Duck Impersonator: QUACK!
    Another Duckling: QUACK!
    Duck: QUACK!

    Original (modified):
    duck: quack
    duckling: quack
    duck impersonator: quack
    another duckling: quack
    cat: meow
    dog: woof
    cow: moo

     */
}
