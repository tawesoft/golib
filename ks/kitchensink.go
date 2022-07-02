// x-doc-short-desc: misc helpful things
// x-doc-stable: no

// Package ks ("kitchen sink") implements assorted helpful things that don't
// fit anywhere else.
package ks

// CONTRIBUTORS: keep definitions in alphabetical order.

import (
    "fmt"
)

// Catch calls the input function f. If successful, Catch passes on the return
// value from f and also returns a nil error. If f panics, Catch recovers from
// the panic and returns a non-nil error.
//
// If the panic raised by f contains is of type error, the returned error
// is wrapped once.
func Catch[X any](f func() X) (x X, err error) {
    defer func() {
        if r := recover(); r != nil {
            if rErr, ok := r.(error); ok {
                err = fmt.Errorf("caught panic: %w", rErr)
            } else {
                err = fmt.Errorf("caught panic: %v", r)
            }
        }
    }()

    return f(), nil
}

// IfThenElse returns a value based on a boolean condition, q. Iff q is true,
// returns the ifTrue. Iff q is false, returns ifFalse. This [IfThenElse
// expression] (as distinct from If-Then-Else statements) is much like the
// ternary operator in some other languages, however it is not short-circuited
// and both arguments are evaluated.
//
// For a lazily-evaluated version, see [lazy.IfThenElse].
//
// [IfThenElse expression]: https://en.wikipedia.org/wiki/Conditional_(computer_programming)#If%E2%80%93then%E2%80%93else_expressions
func IfThenElse[X any] (
    q       bool,
    ifTrue  X,
    ifFalse X,
) X {
    if q {
        return ifTrue
    } else {
        return ifFalse
    }
}

// Item is any Key, Value pair. Type K is any type that would be suitable as a
// KeyType in a Go [builtin.map].
//
// A downstream package should use this to define its own number type (e.g.
// type Item[K comparable, V any] = ks.Item[K, V]) rather than use the type
// directly from here in its exported interface.
type Item[K comparable, V any] struct {
    Key   K
    Value V
}

// Must accepts a (value, err) tuple as input and panics if err != nil,
// otherwise returns value. The error raised by panic is wrapped in another
// error.
//
// For example, Must(os.Open("doesnotexist")) panics with an error like
// "unexpected error in Must[*os.File]: open doesnotexist: no such file or
// directory". Must(os.Open("filethatexists")) returns a pointer to an
// [os.File].
func Must[T any](t T, err error) T {
    if err != nil {
        panic(fmt.Errorf("unexpected error in Must[%T]: %w", t, err))
    }
    return t
}

// MustFunc accepts a function that takes an input of type X, where that
// function then returns a (value Y, err) tuple. Must then returns a function
// that panics if the returned err != nil, otherwise returns value Y. The
// returned error is wrapped in another error.
//
// For example, MustFunc(os.Open) returns a function (call this f).
// f("doesnotexist") panics with an error (like [Must]), and
// f("filethatexists") returns a pointer to an [os.File].
func MustFunc[X any, Y any](
    f func (x X) (Y, error),
) func (x X) Y {
    return func(x X) Y {
        return Must(f(x))
    }
}

// Number defines anything you can perform arithmetic with using standard Go
// operators (like a + b, or a ^ b).
//
// A downstream package should use this to define its own number type (e.g.
// type Number = ks.Number) rather than use the type directly from here in its
// exported interface.
type Number interface {
     ~int8 |  ~int16 |  ~int32 |  ~int64 |
    ~uint8 | ~uint16 | ~uint32 | ~uint64 |
                      ~float32 | float64 |
                              ~complex64 | ~complex128
}

// Zero returns the zero value for any type.
func Zero[T any]() T {
    var t T
    return t
}
