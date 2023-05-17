package partial_test

import (
    "fmt"
    "math"

    "github.com/tawesoft/golib/v2/fun/maybe"
    "github.com/tawesoft/golib/v2/fun/partial"
)

func Example_Line() {
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

    // we create a function, divider, with one argument, that can be used to
    // construct new functions that divide by a constant factor.
    divider := partial.Right2(maybeDivide)

    // bind a constant factor to the divide function, and also convert it back
    // to a function that returns (value int, ok bool) instead of a [maybe.M].
    divideByTwo := maybe.UnwrapFunc(divider(2))
    divideByZero := maybe.UnwrapFunc(divider(0))

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

func Example_All() {
    // Pythagoras theorem for calculating the hypotenuse of a triangle:
    // a squared + b squared = c squared.
    hyp := func(a float64, b float64) float64 {
        return math.Sqrt((a*a) + (b*b))
    }

    // Suppose we want a function with no arguments that returns the answer
    // for a specific triangle. Imagine it's an expensive computation, we've
    // got lots of computations like this, and they're ones we might want to
    // run asynchronously across multiple workers. This is a simple "promise"
    // construct.
    //
    // Normally we could define it like this:
    hyp_2_3_verbose := func() float64 { return hyp(2, 3) }

    // But we can use "partial.All*" functions to do this for us (here, 2,
    // for the two arguments).
    hyp_2_3_terse := partial.All2(hyp)(2, 3)

    fmt.Printf("hyp_2_3_verbose = %.3f\n", hyp_2_3_verbose())
    fmt.Printf("hyp_2_3_terse = %.3f\n", hyp_2_3_terse())

    // Output:
    // hyp_2_3_verbose = 3.606
    // hyp_2_3_terse = 3.606
}
