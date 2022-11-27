package tokenizer_test

// TODO test cases from
// https://chromium.googlesource.com/chromium/src/+/22eeef8fc52576bf54a81b39555030eea9629d35/third_party/blink/renderer/core/css/parser/css_tokenizer_test.cc

import (
    "fmt"
    "strings"

    "github.com/tawesoft/golib/v2/css/tokenizer"
)

func ExampleTokenizer() {
    str := `
/* hello *//* world */
/* don't mind the comments */

"a hundred pounds is \A3 100"

#something {
    background: rgb(128, 64, 64);
}

+100%

+125.45cm
`
    t := tokenizer.New(strings.NewReader(str))

    for {
        tok, ok := t.Next()
        if !ok { break }
        fmt.Println(tok)
    }

    if len(t.Errors()) > 0 {
        fmt.Printf("%v\n", t.Errors())
    }

    // Output:
    // whitespace
}
