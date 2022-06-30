package ks_test

import (
    "fmt"

    "lib/kitchensink"
)

func ExampleZero() {

    type thing struct {
        number int
        phrase string
    }

    fmt.Printf("The zero value is %+v\n", kitchensink.Zero[thing]())
    fmt.Printf("The zero value is %+v\n", kitchensink.Zero[int32]())

    // Output:
    // The zero value is {number:0 phrase:}
    // The zero value is 0
}
