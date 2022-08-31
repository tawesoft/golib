package maybe_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/fun/maybe"
    "github.com/tawesoft/golib/v2/fun/partial"
)

func ExampleMaybe_FlatMap() {
    type Data struct {
        title string
        body string
    }

    getPost := func(id string) (s string, ok bool) {
        return "# Hello world\n---------\n\nWelcome to my website.", true
    }
    getPostM := maybe.WrapFunc(getPost)

    parseData := func(data string) (d Data, ok bool) {
        return Data{
            title: "Hello world",
            body:  "Welcome to my website.",
        }, true
    }
    parseDataM := maybe.WrapFunc(parseData)

    getTitle := func(data Data) string {
        return "Hello world!"
    }

    post := getPostM("20220825-hello-world")
    data := maybe.FlatMap(post, parseDataM)
    title := maybe.Map(data, getTitle)

    fmt.Println(title)

    // Output:
    // {Hello world! true}
}

func ExampleMaybe_Apply() {
    friend := func(a string, b string) string {
        return fmt.Sprintf("%s is friends with %s", a, b)
    }

    friendPartial := partial.Left2(friend)

    alice := maybe.Some("Alice")
    bob := maybe.Some("Bob")
    charlie := maybe.Some("Charlie")
    nobody := maybe.Nothing[string]()

    friendsWithAlice  := maybe.Map(alice,  friendPartial)
    friendsWithBob    := maybe.Map(bob,    friendPartial)
    friendsWithNobody := maybe.Map(nobody, friendPartial)

    fmt.Println(maybe.Apply(bob,     friendsWithAlice))
    fmt.Println(maybe.Apply(charlie, friendsWithBob))
    fmt.Println(maybe.Apply(alice,   friendsWithNobody))
    fmt.Println(maybe.Apply(nobody,  friendsWithAlice))

    // Output:
    // {Alice is friends with Bob true}
    // {Bob is friends with Charlie true}
    // { false}
    // { false}
}

func ExampleMaybe_FlatApply() {
    divide := func(a int, b int) maybe.Maybe[int] {
        if b == 0 { return maybe.Nothing[int]() }
        return maybe.Some(a / b)
    }

    dividePartial := partial.Right2(divide)

    type Row struct {
        Initial int
        DivideBy int
        Iterations int
    }

    rows := []Row{
        {32, 2, 3},  // 32 => 16 => 8 => 4
        {32, 0, 10}, // divide by zero
    }

    for _, row := range rows {
        divider := maybe.Map(maybe.Some(row.DivideBy), dividePartial)
        x := maybe.Some(row.Initial)
        for i := 0; i < row.Iterations; i ++ {
            x = maybe.FlatApply(x, divider)
        }
        fmt.Println(x)
    }

    // Output:
    // {4 true}
    // {0 false}
}
