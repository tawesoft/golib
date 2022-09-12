package result_test

import (
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/tawesoft/golib/v2/fun/maybe"
    "github.com/tawesoft/golib/v2/fun/partial"
    "github.com/tawesoft/golib/v2/fun/result"
    "github.com/tawesoft/golib/v2/fun/slices"
)

func ExampleR() {
    resultOpen := result.WrapFunc(os.Open)
    toReader := result.Map(func (x *os.File) io.Reader { return x })

    closer := func(f *os.File) error {
        fmt.Printf("Closing %s\n", f.Name())
        return f.Close()
    }

    paths := []string{
        "testdata/example1.txt",
        "testdata/example2.txt",
        "testdata/example3.txt",
    }

    handles := slices.Map(resultOpen, paths)
    defer slices.Map(result.Map(closer), handles)

    readers := slices.Map(toReader, handles)

    resultRead := result.FlatMap(result.WrapFunc(io.ReadAll))
    toString := result.Map(func(x []byte) string { return string(x) })
    resultReadString := result.Compose(resultRead, toString)

    contents := slices.Map(resultReadString, readers)

    for _, x := range contents {
        if x.Failed() {
            fmt.Printf("Could not read from file: %v\n", x.Error)
            continue
        }
        fmt.Println(strings.TrimSpace(x.Value))
    }

    // Output:
    // This is the first file!
    // Could not read from file: open testdata/example2.txt: no such file or directory
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
