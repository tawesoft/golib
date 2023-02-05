package partial_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/fun/maybe"
    "github.com/tawesoft/golib/v2/fun/partial"
)

func Example() {
    // The formula for a line can be given by "y = mx + c", where m is the
    // gradient, and c is the offset where the line crosses the x-axis.
    line := func(x int, m int, c int) int { // solves for y
        return (m * x) + c // y = mx + c
    }

    // That's the general formula. Let's say we know c and m, and instead want
    // a single function y = f(x) that we can plug in values for x and get
    // a result for y.

    // For example, y = 3x + 5 where f(n) is f((3*n) + 5).

    // Given c, we can partially apply the function so y = f(mx) where
    // f(n) = n+5. To put it another way, c is now bound, and we only need m
    // and x as inputs now.
    lineC := partial.Right3(line)(5) // f(x, m, c) => f(x, m)

    // Given m, we can partially apply the function so y = f(g(x)) where g(n)
    // = m * n. To put it another way, m is now also bound, and we only need x
    // as an input now.
    lineMC := partial.Right2(lineC)(2) // f(x, m) => f(x)

    // Now we have a function y = 2x + 5.

    // Inputting x = 3...
    fmt.Println(lineMC(3)) // y = (2 * 3) + 5 = 11

    // Output:
    // 11
}

func Example_Maybe() {
    // divides two numbers, while checking for divide by zero.
    divide := func(x int, y int) (value int, ok bool) {
        if y == 0 { return 0, false }
        return x / y, true
    }

    // we need a single return value, so convert the function to one that
    // returns a single [maybe.M].
    maybeDivide := maybe.WrapFunc2(divide)

    // bind y to the divide function, and also convert it back from a function
    // that returns (value int, ok bool) instead of a [maybe.M].
    divideByTwo := maybe.UnwrapFunc(partial.Right2(maybeDivide)(2))
    divideByZero := maybe.UnwrapFunc(partial.Right2(maybeDivide)(0))

    {
        result, ok := divideByTwo(10)
        fmt.Printf("divideByTwo(10) = %d, %t\n", result, ok)
    }
    {
        result, ok := divideByZero(10)
        fmt.Printf("divideByZero(10) = %d, %t\n", result, ok)
    }

    // Output:
    // divideByTwo(10) = 5, true
    // divideByZero(10) = 0, false
}
