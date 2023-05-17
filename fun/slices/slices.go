// Package slices provides generic higher-order functions over slices of values.
package slices

// FromArgs returns the slice of the variadic arguments list.
func FromArgs[X any](xs ... X) []X {
    return xs
}

// ForEach applies the void function "f(x)" to each value X of the input slice.
func ForEach[X any](f func(X), xs []X) {
    for _, x := range xs {
        f(x)
    }
}

// Filter applies function "f : X => bool" to each value X of the input slice,
// and returns a new slice containing only each X for which f(X) is true.
//
// If f(X) is false for every item, this function may return a nil slice.
func Filter[X any](f func(X) bool, xs []X) []X {
    if xs == nil { return nil }
    var result []X
    for _, x := range xs {
        if !f(x) { continue }
        result = append(result, x)
    }
    return result
}

// Map applies function "f : X => Y" to each value of the input slice, and
// returns a new slice of each output in sequence.
//
// If the input slice is a nil slice, the return value is also a nil slice.
//
// For a lazy version, see the iter package.
func Map[X any, Y any](f func (X) Y, xs []X) []Y {
    if xs == nil { return nil }
    result := make([]Y, 0, len(xs))
    for _, x := range xs {
        result = append(result, f(x))
    }
    return result
}

// FlatMap applies function "f : X => []Y" to each value of the input slice,
// and returns a new slice of each output in sequence. The output sequence
// is "flattened" so that each slice returned by f is concatenated into a
// single slice.
//
// If the input slice is a nil slice, the return value is also a nil slice.
// If f returns a nil slice, it is not included in the output.
//
// For a lazy version, see the iter package.
func FlatMap[X any, Y any](f func (X) []Y, xs []X) []Y {
    if xs == nil { return nil }
    result := make([]Y, 0)
    for _, x := range xs {
        ys := f(x)
        for _, y := range ys {
            result = append(result, y)
        }
    }
    return result
}

// Reduce applies function f to each element of the slice (in increasing order
// of element index) with the previous return value from f as the first
// argument and the element as the second argument (for the first call to f,
// the supplied initial value is used instead of a return value), returning the
// final value returned by f.
//
// If the input slice is empty, the result is the initial value.
func Reduce[X any](initial X, f func (X, X) X, xs []X) X {
    if len(xs) == 0 { return initial }

    result := initial
    for _, x := range xs {
        result = f(result, x)
    }

    return result
}

// Reducer constructs a partially applied Reduce function with the arguments
// "initial" and "f" already bound.
func Reducer[X any](initial X, f func(X, X) X) func(xs []X) X {
    return func(xs []X) X {
        return Reduce[X](initial, f, xs)
    }
}
