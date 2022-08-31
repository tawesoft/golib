package iter_test

import (
    "fmt"

    lazy "github.com/tawesoft/golib/v2/iter"
)

// FibFunc lazily generates a sequence of fibonacci numbers.
func FibFunc() lazy.It[int] {
    // if you prefer,
    //     func FibFunc() func() func() (int, bool) { ...
    pre1 := 0
    pre2 := 1
    return func() (int, bool) {
        result := pre1 + pre2
        pre1 = pre2
        pre2 = result
        return result, true
    }
}

func ExampleFibonacci() {
    // return a new generator of fibonacci numbers
    isOdd := func(x int) bool { return x % 2 != 0 }

    sum := lazy.Reducer[int]{
        Identity: 0,
        Reduce: func (a int, b int) int { return a + b },
    }

    // create four iterators
    fib := lazy.Tee[int](4, lazy.Func(FibFunc()))

    fmt.Printf("First ten Fibonacci numbers are:\n    %v\n",
        lazy.ToSlice(
            lazy.Enumerate(
                lazy.Take(10,
                    fib[0]))))

    fmt.Printf("First five odd Fibonacci numbers are:\n    %v\n",
        lazy.ToSlice(
            lazy.Enumerate(
                lazy.Take(5,
                    lazy.Filter(isOdd,
                        fib[1])))))

    fmt.Printf("Sum of the first 10 Fibonacci numbers is: %d\n",
        lazy.Reduce(sum,
            lazy.Take(10,
                fib[2])))

    average := func(n int, xs lazy.It[int]) float64 {
        total := lazy.Reduce(sum, lazy.Take(n, xs))
        return float64(total) / float64(n)
    }

    fmt.Printf("Average of the first 5 Fibonacci numbers is: %.2f\n",
        average(5, fib[3]))

    // or...
    // fmt.Printf("Average of the first 5 Fibonacci numbers is: %.2f\n",
    //    lazy.Join(lazy.AverageJoiner[int](), lazy.Take(5, fib[3])))

    // Output:
    // First ten Fibonacci numbers are:
    //     [{0 1} {1 2} {2 3} {3 5} {4 8} {5 13} {6 21} {7 34} {8 55} {9 89}]
    // First five odd Fibonacci numbers are:
    //     [{0 1} {1 3} {2 5} {3 13} {4 21}]
    // Sum of the first 10 Fibonacci numbers is: 231
    // Average of the first 5 Fibonacci numbers is: 3.80
}
