// Slices provides generic higher-order functions over slices of values.
package slices

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
