package parser_test

import (
    "fmt"
    "strings"

    "github.com/tawesoft/golib/v2/css/internal/parser"
    "github.com/tawesoft/golib/v2/css/tokenizer/token"
)

func ExampleParser() {
    str := `rgb(128, 64, 64)`
    p := parser.New(strings.NewReader(str))
    k := p.Tokenizer()

    for {
        cv := p.ConsumeComponentValue(k)
        if cv.IsPreservedToken(token.EOF()) { break }
        fmt.Println(cv)
    }

    /*
    if len(t.Errors()) > 0 {
        fmt.Printf("%v\n", t.Errors())
    }
    */

    // Output:
    // <ComponentValue/Function{Name: "rgb", Value: [<ComponentValue:<number-token>{type: "integer", value: 128.000000, repr: "128"}> <ComponentValue:<comma-token>> <ComponentValue:<whitespace-token>> <ComponentValue:<number-token>{type: "integer", value: 64.000000, repr: "64"}> <ComponentValue:<comma-token>> <ComponentValue:<whitespace-token>> <ComponentValue:<number-token>{type: "integer", value: 64.000000, repr: "64"}>]}>
}
