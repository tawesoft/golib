package tuple_test

import (
    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/tuple"
)

func Example() {
    // types inferred implicitly
    a := tuple.ToT2(5, true)
    b := tuple.ToT2(7, true)
    must.Not(a == b)

    // explicitly typed
    c := tuple.ToT2[int32, bool](5, true)
    d := tuple.ToT2[int32, bool](7, true)
    must.Not(c == d)

    // Output:
}
