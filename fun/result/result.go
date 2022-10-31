// Package result implements a R{value, error} "sum type" that has a value
// only when error is nil.
//
// Note that in many cases, it is more idiomatic for a function to return a
// naked (value, error). Use [WrapFunc] to convert such a function to return
// a Result type.
//
// For examples, see the sibling "maybe" package.
package result

import (
    "fmt"
)
// R is a (value, error) "sum type" that has a value only when error is
// nil.
type R[V any] struct {
    Value V
    Error error
}

// Success returns true if the R is not an error.
func (r R[V]) Success() bool {
    return r.Error == nil
}

// Failed returns true if R has an error.
func (r R[V]) Failed() bool {
    return r.Error != nil
}

// New returns a R. It is syntax sugar for R{value, error}. If error is
// a known constant, use [Some] or [Error] instead.
func New[V any](value V, err error) R[V] {
    if err != nil { return Error[V](err) }
    return Some(value)
}

// Unpack returns a plain (value, error) tuple from a R.
func (r R[V]) Unpack() (V, error) {
    return r.Value, r.Error
}

// Must returns a R's value. If the R is an error, panics.
func (r R[V]) Must() V {
    if r.Error != nil {
        panic(fmt.Sprintf("result.R[%T].Must called, but is error.", r))
    }
    return r.Value
}

// MustError returns a R's error. Panics if the R is not an error.
func (r R[V]) MustError() error {
    if r.Error != nil {
        panic(fmt.Sprintf("result.R[%T].MustError called, but is not an error.", r))
    }
    return r.Error
}

// If calls function f on a R's value, but only if R is not an error.
func (r R[V]) If(f func(v V)) {
    if r.Error != nil {
        f(r.Value)
    }
}

// IfNot calls function f on a R's error, but only if R is an error.
func (r R[V]) IfNot(f func(err error)) {
    if r.Error != nil {
        f(r.Error)
    }
}

// Error returns a R that is an error.
func Error[V any](err error) R[V] {
    return R[V]{
        Error: err,
    }
}

// Some returns a R that contains a value and is not an error.
func Some[V any](value V) R[V] {
    return R[V]{
        Value: value,
        Error: nil,
    }
}

// Else returns R.value if not an error, otherwise returns the provided
// argument instead.
func (r R[V]) Else(v V) V {
    if r.Error == nil { return r.Value }
    return v
}

// Map turns function "f: X => Y" into "f: R(X) => R[Y]".
func Map[X any, Y any](
    f func(x X) Y,
) func(x2 R[X]) R[Y] {
    return func(x2 R[X]) R[Y] {
        if x2.Error != nil { return Error[Y](x2.Error) }
        return Some(f(x2.Value))
    }
}

// FlatMap turns function "f: X => R[Y]" into "f: R(X) => R[Y]".
func FlatMap[X any, Y any](
    f func(x X) R[Y],
) func(x2 R[X]) R[Y] {
    return func(x2 R[X]) R[Y] {
        if x2.Error != nil { return Error[Y](x2.Error) }
        return f(x2.Value)
    }
}

// Applicator turns function "R[f]: X => Y" into "f: X => R[Y]".
func Applicator[X any, Y any](
    f R[func(x X) Y],
) func(x X) R[Y] {
    if f.Error != nil { return func(x X) R[Y] { return Error[Y](f.Error) } }
    return func(x X) R[Y] { return Some(f.Value(x)) }
}

// WrapFunc converts a function of the form "f: X => (Y, error)" to the form
// "f: X => R[X].
func WrapFunc[X any, Y any](
    f func(x X) (Y, error),
) func(x X) R[Y] {
    return func(x X) R[Y] {
        return New(f(x))
    }
}

// UnwrapFunc converts a function of the form "f: X => R[Y]" to the
// form "f: X => ([Y], error)".
func UnwrapFunc[X any, Y any](
    f func(x X) R[Y],
) func(x X) (Y, error) {
    return func(x X) (Y, error) {
        return f(x).Unpack()
    }
}

// Lift converts a function of the form "f: X => Y" to the form "f: X => R[Y]"
// where R[Y] == Some(y).
func Lift[X any, Y any](
    f func(x X) Y,
) func(x X) R[Y] {
    return func(x X) R[Y] {
        return Some(f(x))
    }
}

// Compose takes two functions of the form "xy: R[X] => R[Y]" and
// "yz: R[Y] => R[Z]" and returns a function "xz(R[X]) => R[Z]".
func Compose[X any, Y any, Z any](
    xy func(R[X]) R[Y],
    yz func(R[Y]) R[Z],
) func(R[X]) R[Z] {
    return func(x R[X]) R[Z] {
        return yz(xy(x))
    }
}
