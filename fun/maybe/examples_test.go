package maybe_test

import (
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/tawesoft/golib/v2/fun/maybe"
    "github.com/tawesoft/golib/v2/fun/partial"
    "github.com/tawesoft/golib/v2/fun/slices"
)

func ExampleM() {
    maybeOpen := func(x string) maybe.M[*os.File] {
        f, err := os.Open(x)
        if err != nil { return maybe.Nothing[*os.File]() }
        return maybe.Some(f)
    }

    closer := func(f *os.File) error {
        fmt.Printf("Closing %s\n", f.Name())
        return f.Close()
    }

    paths := []string{
        "testdata/example1.txt",
        "testdata/example2.txt",
        "testdata/example3.txt",
    }

    handles := slices.Map(maybeOpen, paths)
    defer slices.Map(maybe.Map(closer), handles)

    maybeRead := func(x *os.File) maybe.M[string] {
        content, err := io.ReadAll(x)
        if err != nil { return maybe.Nothing[string]() }
        return maybe.Some(string(content))
    }

    contents := slices.Map(maybe.FlatMap(maybeRead), handles)

    for i, x := range contents {
        if x.Ok {
            fmt.Println(strings.TrimSpace(x.Value))
        } else {
            fmt.Printf("Could not read from file %d\n", i + 1)
        }
    }

    // Output:
    // This is the first file!
    // Could not read from file 2
    // This is the third file!
    // Closing testdata/example1.txt
    // Closing testdata/example3.txt
}

func ExampleApplicator() {
    friend := func(a string, b string) string {
        return fmt.Sprintf("%s is friends with %s", a, b)
    }
    friendPartial := partial.Left2(friend) // f(a) => f(b) => string
    friendM := maybe.Map(friendPartial)

    fmt.Printf("func friend: %T\n", friend)
    fmt.Printf("func friendPartial: %T\n", friendPartial)
    fmt.Printf("func friendM: %T\n", friendM)

    alice := maybe.Some("Alice")
    bob := maybe.Some("Bob")
    charlie := maybe.Some("Charlie")
    nobody := maybe.Nothing[string]()

    friendsWithAlice  := maybe.FlatMap(maybe.Applicator(friendM(alice)))
    friendsWithBob    := maybe.FlatMap(maybe.Applicator(friendM(bob)))
    friendsWithNobody := maybe.FlatMap(maybe.Applicator(friendM(nobody)))

    fmt.Printf("func friendsWithAlice: %T\n\n", friendsWithAlice)

    fmt.Println(friendsWithAlice(bob).Must())
    fmt.Println(friendsWithBob(charlie).Must())
    friendsWithNobody(alice).MustNot()
    friendsWithAlice(nobody).MustNot()

    // Output:
    // func friend: func(string, string) string
    // func friendPartial: func(string) func(string) string
    // func friendM: func(maybe.M[string]) maybe.M[func(string) string]
    // func friendsWithAlice: func(maybe.M[string]) maybe.M[string]
    //
    // Alice is friends with Bob
    // Bob is friends with Charlie
}
