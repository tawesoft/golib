// Package operator implements builtin language operators, such as "=="
// (equals) or "+" (addition), as functions that can be passed to higher order
// functions.
//
// See also the [github.com/tawesoft/golib/v2/operator/checked/integer] package
// which implements operators that are robust against integer overflow.
//
// See also the [github.com/tawesoft/golib/v2/operator/reflect] package
// which implements operators that require reflection.
package operator

import (
    "golang.org/x/exp/constraints"
)

// Zero returns the zero value for any type.
func Zero[T any]() T {
    var t T
    return t
}

// Ternary returns a if q is true, or b if q is false.
func Ternary[X any](q bool, a X, b X) X {
    if q { return a } else { return b }
}

// In returns true if x equals any of the following arguments.
func In[X comparable](x X, xs ... X) bool {
    for _, i := range xs {
        if x == i { return true }
    }
    return false
}

// Identity implements the function f(x) => x.
func Identity[X any](x X) X {
    return x
}

// IsZero returns true iff x is equal to the zero value of its type.
func IsZero[X comparable](x X) bool {
    var y X
    return x == y
}

// IsNonZero returns true iff x is not equal to the zero value of its type.
func IsNonZero[X comparable](x X) bool {
    var y X
    return x != y
}

// Equal returns a == b.
func Equal[C comparable](a C, b C) bool {
    return a == b
}

// NotEqual returns a != b.
func NotEqual[C comparable](a C, b C) bool {
    return a != b
}

// Cmp returns integer 1, 0, or -1 depending on whether a is greater than,
// equal to, or less than b.
func Cmp[O constraints.Ordered](a O, b O) int {
    switch {
        case (a == b): return  0
        case (a < b):  return -1
        default:       return +1
    }
}

// LT returns a < b.
func LT[O constraints.Ordered](a O, b O) bool {
    return a < b
}

// LTE returns a <= b.
func LTE[O constraints.Ordered](a O, b O) bool {
    return a <= b
}

// GT returns a > b.
func GT[O constraints.Ordered](a O, b O) bool {
    return a > b
}

// GTE returns a >= b.
func GTE[O constraints.Ordered](a O, b O) bool {
    return a >= b
}
