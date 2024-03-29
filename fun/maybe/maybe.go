// Package maybe implements a M{value, ok} "sum type" that has a value only
// when ok is true.
//
// Note that in many cases, it is more idiomatic for a function to return a
// naked (value, ok). Use [WrapFunc] to convert such a function to return
// a M maybe type.
package maybe

import (
    "fmt"
)

// M is a (value, ok) "sum type" that has a value only when ok is
// true.
type M[V any] struct {
    Value V
    Ok bool
}

// New returns a M. It is syntax sugar for M{value, ok}. If ok is
// a known constant, use [Some] or [Nothing] instead.
func New[V any](value V, ok bool) M[V] {
    if !ok { return Nothing[V]() }
    return Some(value)
}

// Unpack returns a plain (value, ok) tuple from a M.
func (m M[V]) Unpack() (V, bool) {
    return m.Value, m.Ok
}

// Must returns a M's value. If the M is not ok, panics.
func (m M[V]) Must() V {
    if !m.Ok {
        panic(fmt.Sprintf("maybe.M[%T].Must called, but value missing.", m))
    }
    return m.Value
}

// MustNot panics if the M is ok.
func (m M[V]) MustNot() {
    if m.Ok {
        panic(fmt.Sprintf("maybe.M[%T].MustNot called, but value present.", m))
    }
}

// Nothing returns a (typed) M that has no value.
func Nothing[V any]() M[V] {
    return M[V]{}
}

// Some returns a M that contains a value.
func Some[V any](value V) M[V] {
    return M[V]{
        Value: value,
        Ok:    true,
    }
}

// Or returns M.value if ok, otherwise returns the provided argument instead.
func (m M[V]) Or(v V) V {
    if m.Ok { return m.Value }
    return v
}

// Map turns function "f: X => Y" into "f: M(X) => M[Y]".
func Map[X any, Y any](
    f func(x X) Y,
) func(x2 M[X]) M[Y] {
    return func(x2 M[X]) M[Y] {
        if (!x2.Ok) { return Nothing[Y]() }
        return Some(f(x2.Value))
    }
}

// FlatMap turns function "f: X => M[Y]" into "f: M(X) => M[Y]".
func FlatMap[X any, Y any](
    f func(x X) M[Y],
) func(x2 M[X]) M[Y] {
    return func(x2 M[X]) M[Y] {
        if (!x2.Ok) { return Nothing[Y]() }
        return f(x2.Value)
    }
}

// Collect takes a slice of Maybe[X] and returns a slice of []X and true iff
// every element contains a value.
func Collect[X any](
    xs []M[X],
) ([]X, bool) {
    if len(xs) == 0 { return nil, true }
    for _, x := range xs {
        if !x.Ok { return nil, false }
    }
    result := make([]X, 0, len(xs))
    for _, x := range xs {
        result = append(result, x.Value)
    }
    return result, true
}

// Applicator turns function "M[f: X => Y]" into "f: X => M[Y]".
func Applicator[X any, Y any](
    f M[func(x X) Y],
) func(x X) M[Y] {
    if !f.Ok { return func(x X) M[Y] { return Nothing[Y]() } }
    return func(x X) M[Y] { return Some(f.Value(x)) }
}

// WrapFunc converts a function of the form "f: X => (Y, ok bool)" to the form
// "f: X => M[X].
func WrapFunc[X any, Y any](
    f func(x X) (Y, bool),
) func(x X) M[Y] {
    return func(x X) M[Y] {
        return New(f(x))
    }
}

// WrapFunc2 converts a function of the form "f(A, B) => (C, ok bool)" to
// the form "f(A, B) => M[C].
func WrapFunc2[A any, B any, C any](
    f func(a A, b B) (C, bool),
) func(a A, b B) M[C] {
    return func(a A, b B) M[C] {
        return New(f(a, b))
    }
}

// WrapFunc3 converts a function of the form "f(A, B, C) => (D, ok bool)" to
// the form "f(A, B, C) => M[D].
func WrapFunc3[A any, B any, C any, D any](
    f func(a A, b B, c C) (D, bool),
) func(a A, b B, c C) M[D] {
    return func(a A, b B, c C) M[D] {
        return New(f(a, b, c))
    }
}

// WrapFunc4 converts a function of the form "f(A, B, C, D) => (E, ok bool)" to
// the form "f(A, B, C, D) => M[E].
func WrapFunc4[A any, B any, C any, D any, E any](
    f func(a A, b B, c C, d D) (E, bool),
) func(a A, b B, c C, d D) M[E] {
    return func(a A, b B, c C, d D) M[E] {
        return New(f(a, b, c, d))
    }
}

// UnwrapFunc converts a function of the form "f: X => M[Y]" to the
// form "f: X => ([Y], ok bool)".
func UnwrapFunc[X any, Y any](
    f func(x X) M[Y],
) func(x X) (Y, bool) {
    return func(x X) (Y, bool) {
        return f(x).Unpack()
    }
}

// Lift converts a function of the form "f: X => Y" to the form "f: X => M[Y]"
// where M[Y] == Some(y).
func Lift[X any, Y any](
    f func(x X) Y,
) func(x X) M[Y] {
    return func(x X) M[Y] {
        return Some(f(x))
    }
}

// Compose takes two functions of the form "xy: M[X] => M[Y]" and
// "yz: M[Y] => M[Z]" and returns a function "xz(M[X]) => M[Z]".
func Compose[X any, Y any, Z any](
    xy func(M[X]) M[Y],
    yz func(M[Y]) M[Z],
) func(M[X]) M[Z] {
    return func(x M[X]) M[Z] {
        return yz(xy(x))
    }
}
