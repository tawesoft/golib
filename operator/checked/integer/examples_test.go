package integer_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/operator/checked/integer"
)

func Example_Simple() {
    {
        result, ok := integer.Uint8.Add(250, 5)
        fmt.Printf("integer.Uint8.Add(250, 5): %d, ok?=%t\n", result, ok)
    }

    {
        result, ok := integer.Uint8.Add(250, 6)
        fmt.Printf("integer.Uint8.Add(250, 6): %d, ok?=%t\n", result, ok)
    }

    // Output:
    // integer.Uint8.Add(250, 5): 255, ok?=true
    // integer.Uint8.Add(250, 6): 0, ok?=false
}

func Example_Limits() {
    {
        const min = 0
        const max = 99
        result, ok := integer.Sub(min, max, 10, 9)
        fmt.Printf("integer.Sub(min, max, 10, 9): %d, ok?=%t\n", result, ok)
    }

    {
        limit := integer.Limits[int]{Min: 0, Max: 99}
        result, ok := limit.Sub(10, 25)
        fmt.Printf("limit.Sub(10, 25): %d, ok?=%t\n", result, ok)
    }

    // Output:
    // integer.Sub(min, max, 10, 9): 1, ok?=true
    // limit.Sub(10, 25): 0, ok?=false
}
