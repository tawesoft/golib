package algorithmic_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/text/number/algorithmic"
)

func ExampleFormat() {
    print := func(system string, n int64) {
        if s, err := algorithmic.Format(system, n); err == nil {
            fmt.Printf("algorithmic.Format(%q, %d): %q\n", system, n, s)
        } else {
            fmt.Printf("algorithmic.Format(%q, %d): error: %v\n", system, n, err)
        }
    }

    print("%roman-upper", 2023)

    // Output:
    // algorithmic.Format("%roman-upper", 2023): "MMXXIII"
}
