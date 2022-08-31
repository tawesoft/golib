package partial_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/fun/partial"
)

func Example() {
    // The formula for a line can be given by "y = mx + c", where m is the
    // gradient, and c is the offset where the line crosses the x axis.
    line := func(x int, m int, c int) int { // solves for y
        return (m * x) + c // y = mx + c
    }

    // That's the general formula. Let's say we know c and m, and instead want
    // a single function y = f(x) that we can plug in values for x and get
    // a result for y.

    // For example, y = 3x + 5 where f(n) is f((3*n) + 5).

    // Given c, we can partially apply the function so y = f(mx) where
    // f(n) = n+5. To put it another way, c is now bound and we only need m
    // and x as inputs now.
    lineC := partial.Right3(line)(5)

    // Given m, we can partially apply the function so y = f(g(x)) where
    // g(n) = m * n. To put it another way, m is now bound and we only need
    // x as an input now.
    lineMC := partial.Right2(lineC)(2)

    // Now we have a function y = 2x + 5.

    // Inputting x = 3...
    fmt.Println(lineMC(3)) // y = (2 * 3) + 5 = 11

    // output:
    // 11
}
