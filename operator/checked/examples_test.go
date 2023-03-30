package checked_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/operator/checked"
)

func ExampleSimple() {
    {
        result, ok := checked.Uint8.Add(250, 5)
        fmt.Printf("checked.Uint8.Add(250, 5): %d, ok?=%t\n", result, ok)
    }

    {
        result, ok := checked.Uint8.Add(250, 6)
        fmt.Printf("checked.Uint8.Add(250, 6): %d, ok?=%t\n", result, ok)
    }

    // Output:
    // checked.Uint8.Add(250, 5): 255, ok?=true
    // checked.Uint8.Add(250, 6): 0, ok?=false
}

func ExampleLimits() {
    {
        const min = 0
        const max = 99
        result, ok := checked.Sub(min, max, 10, 9)
        fmt.Printf("checked.Sub(min, max, 10, 9): %d, ok?=%t\n", result, ok)
    }

    {
        limit := checked.Limits[int]{Min: 0, Max: 99}
        result, ok := limit.Sub(10, 25)
        fmt.Printf("limit.Sub(10, 25): %d, ok?=%t\n", result, ok)
    }

    // Output:
    // checked.Sub(min, max, 10, 9): 1, ok?=true
    // limit.Sub(10, 25): 0, ok?=false
}
