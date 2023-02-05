package operator_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/iter"
    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/operator"
)

func ExampleAdd() {
    // generate a sequence 1, 2, 3, ... 99, 100. This is produced lazily.
    sequence := iter.Take(100, iter.Counter[int](1, 1))

    // reduce applies a function to each element of the sequence. We want
    // addition ("+"), but we need this as a function, so we use operator.Add.
    // Here, [int] is needed to specify which type of the generic function
    // we need. This should match the type of the sequence (in this case, int).
    result := iter.Reduce(0, operator.Add[int], sequence)

    fmt.Printf("sum of numbers from 1 to 100: %d\n", result)

    // Note that the above is given as an example. A better way to sum the
    // numbers from 1 to n is to use Gauss's method or proof by induction and
    // immediately calculate (n+1) * (n/2).
    must.Equal(result, (100+1)*(100/2))

    // Output:
    // sum of numbers from 1 to 100: 5050
}
