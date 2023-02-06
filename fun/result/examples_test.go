package result_test

import (
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/tawesoft/golib/v2/fun/result"
    "github.com/tawesoft/golib/v2/fun/slices"
)

func ExampleResult() {
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
        if !x.Success() {
            err := x.Error
            if os.IsNotExist(err) {
                // hack to ensure error has the same text on every OS
                err = fmt.Errorf("no such file or directory")
            }
            fmt.Printf("Could not read from file: %v\n", err)
            continue
        }
        fmt.Println(strings.TrimSpace(x.Value))
    }

    // Output:
    // This is the first file!
    // Could not read from file: no such file or directory
    // This is the third file!
    // Closing testdata/example1.txt
    // Closing testdata/example3.txt
}
