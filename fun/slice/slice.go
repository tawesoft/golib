// SPDX-License-Identifier: MIT
// x-doc-short-desc: higher-order functions for slices
// x-doc-stable: candidate

// Package slice implements generically-typed higher order functions, like
// [Map], [Filter], [Reduce], including some in-place variants, for slices. For
// lazy evaluation, see [golib/lazy] instead.
//
// To avoid confusion, in this package the Go map data structure is referred to
// as a "dict". "Map" in this package refers to the higher-order function "map"
// used in functional programming.
//
// Unlike many existing generic functional libraries for Go, our map function
// can also map to different types.
package slice

import (
    "github.com/tawesoft/golib/ks"
)

// Walk calls fn(i, x) for each pointer to an element in the input slice xs,
// where i is the index 0 <= 0 < len(xs) This may safely be used to modify the
// elements of the slice in place (but not add or remove from a slice).
func Walk[X any](
    fn func(i int, x *X),
    xs []X,
) {
    for i := 0; i < len(xs); i++ {
        fn(i, &xs[i])
    }
}

// Map returns a new slice made up of the results of calling a function on each
// element of the input slice in sequence.
func Map[X any, Y any](
    fn func(x X) Y,
    xs []X,
) []Y {
    ys := make([]Y, len(xs))

    for i := 0; i < len(xs); i++ {
        ys[i] = fn(xs[i])
    }

    return ys
}

// MapInPlace modifies a slice so that each element is replaced with the result of
// calling a function on it.
func MapInPlace[X any](
    fn func(x X) X,
    xs []X,
) {
    for i := 0; i < len(xs); i++ {
        xs[i] = fn(xs[i])
    }
}

// Filter returns a new slice made up of only the members of the input slice
// where the provided filter function returns true.
func Filter[X any](
    filter func(x X) (bool),
    xs []X,
) []X {
    ys := make([]X, len(xs))
    idx := 0

    for i := 0; i < len(xs); i++ {
        if !filter(xs[i]) { continue }
        ys[idx] = xs[i]
        idx++
    }

    return ys[:idx]
}

// FilterInPlace returns the input slice, but excluding any elements where the
// provided filter function returns true.
func FilterInPlace[X any](
    filter func(x X) (bool),
    xs []X,
) []X {
    idxOut := 0
    z := len(xs)

    for i := 0; i < z; i++ {
        if !filter(xs[i]) { continue }
        xs[idxOut] = xs[i]
        idxOut++
    }

    // zero any pointers so GC doesn't have dangling references
    for j := idxOut; j < z; j++ {
        xs[j] = ks.Zero[X]()
    }

    return xs[:idxOut]
}

type Reducer[X any] struct {
    Reduce func(a X, b X) X
    Identity X
}

// Reduce returns the result of applying a reduce function from left to
// right pairwise on the elements of a slice. The first argument to the first
// invocation of the reduction function is the identity value where
// reduce(identity, x) == x.
func Reduce[X any](
    reducer Reducer[X],
    xs []X,
) X {
    last := reducer.Identity
    for i := 0; i < len(xs); i++ {
        last = reducer.Reduce(last, xs[i])
    }
    return last
}

// ReduceRight is like [ReduceS], but the operations are applied right-to-left
// and the first argument to the first invocation of the reduction function is
// the identity value where reduce(x, identity) == x.
func ReduceRight[X any](
    reducer Reducer[X],
    xs []X,
) X {
    last := reducer.Identity
    for i := len(xs) - 1; i >= 0; i-- {
        last = reducer.Reduce(xs[i], last)
    }
    return last
}
