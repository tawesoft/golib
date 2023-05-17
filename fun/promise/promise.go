// Package promise implements a simple Promise type that can be used to
// represent a computation to be performed at a later stage.
//
// This composes nicely with the idea of futures (see fun/future).
package promise

import (
    "context"
    "errors"
)

var NotOk = errors.New("promised value is not ok")

// P represents a promise to calculate and return some value when Compute or
// ComputeCtx is called.
//
// Compute, or ComputeCtx, should only be called once, unless an implementation
// otherwise indicates that it is safe to do so. A promise is not safe for
// concurrent use, unless an implementation indicates otherwise.
//
// A promise may ignore the provided context if it cannot be cancelled. The
// plain Compute method computes the promise with a context that is never
// cancelled.
//
// The error return value of Compute and ComputeCtx may be an error returned by
// the computation, or, in the case of ComputeCtx, a context error such as
// [context.Cancelled].
type P[X any] interface {
    Compute() (X, error)
    ComputeCtx(ctx context.Context) (X, error)
}

// Func is the type of a function with no arguments that satisfies the promise
// interface by calling the function, ignoring any context.
type Func[X any] func() (X, error)
func (f Func[X]) Compute() (X, error) {
    return f()
}
func (f Func[X]) ComputeCtx(ctx context.Context) (X, error) {
    return f()
}

// FromFunc creates a promise to call function f, where f returns any single
// value and has no facility to indicate an error.
func FromFunc[T any](f func() T) P[T] {
    return FromResultFunc(func() (T, error) {
        return f(), nil
    })
}

// FromResultFunc creates a promise to call function f, where f returns a
// (result, error) tuple.
func FromResultFunc[X any](f func() (X, error)) P[X] {
    return Func[X](f)
}

// WrapResultFunc wraps an existing function "f() => (X, error)" so that it
// becomes "f() => P[X]", a function that returns a promise.
func WrapResultFunc[X any](f func() (X, error)) func() P[X] {
    return func() P[X] {
        return FromValueErr(f())
    }
}

// WrapResultFunc1 wraps an existing function "f(A) => (X, error)" so that it
// becomes "f(A) => P[X]", a function that returns a promise.
func WrapResultFunc1[A, X any](f func(A) (X, error)) func(A) P[X] {
    return func(a A) P[X] {
        return FromValueErr(f(a))
    }
}

// WrapResultFunc2 wraps an existing function "f(A, B) => (X, error)" so that it
// becomes "f(A, B) => P[X]", a function that returns a promise.
func WrapResultFunc2[A, B, X any](f func(A, B) (X, error)) func(A, B) P[X] {
    return func(a A, b B) P[X] {
        return FromValueErr(f(a, b))
    }
}

// WrapResultFunc3 wraps an existing function "f(A, B, C) => (X, error)" so
// that it becomes "f(A, B, C) => P[X]", a function that returns a promise.
func WrapResultFunc3[A, B, C, X any](f func(A, B, C) (X, error)) func(A, B, C) P[X] {
    return func(a A, b B, c C) P[X] {
        return FromValueErr(f(a, b, c))
    }
}

// WrapResultFunc4 wraps an existing function "f(A, B, C, D) => (X, error)" so
// that it becomes "f(A, B, C, D) => P[X]", a function that returns a promise.
func WrapResultFunc4[A, B, C, D, X any](f func(A, B, C, D) (X, error)) func(A, B, C, D) P[X] {
    return func(a A, b B, c C, d D) P[X] {
        return FromValueErr(f(a, b, c, d))
    }
}

// FromOkFunc creates a promise to call function f, where f returns a
// (value, ok) tuple. If the returned ok is false, the promise computes the
// error [NotOk].
func FromOkFunc[X any](f func() (X, bool)) P[X] {
    return Func[X](func() (result X, err error) {
        v, ok := f()
        if !ok { err = NotOk; return }
        return v, nil
    })
}

// FuncCtx is the type of a function with a context argument that satisfies the
// promise interface by calling the function with a context.
type FuncCtx[X any] func(ctx context.Context) (X, error)
func (f FuncCtx[X]) Compute() (X, error) {
    return f(context.TODO())
}
func (f FuncCtx[X]) ComputeCtx(ctx context.Context) (X, error) {
    return f(ctx)
}

// FromFuncCtx creates a promise to call function f, where f accepts a context
// and returns any single value and has no facility to indicate an error,
// other than a context error.
func FromFuncCtx[T any](f func(context.Context) T) P[T] {
    return FromResultFuncCtx(func(ctx context.Context) (T, error) {
        return f(ctx), nil
    })
}

// FromResultFuncCtx creates a promise to call function f, where f accepts a
// context and returns a (result, error) tuple.
func FromResultFuncCtx[X any](f func(context.Context) (X, error)) P[X] {
    return FuncCtx[X](f)
}

// FromOkFuncCtx creates a promise to call function f, where f accepts a
// context and returns a (value, ok) tuple. If the returned ok is false, the
// promise computes the error [NotOk].
func FromOkFuncCtx[X any](f func(ctx context.Context) (X, bool)) P[X] {
    return FuncCtx[X](func(ctx context.Context) (result X, err error) {
        v, ok := f(ctx)
        if !ok { err = NotOk; return }
        return v, nil
    })
}

type value[X any] struct {x X}
func (v value[X]) Compute() (X, error) {
    return v.x, nil
}
func (v value[X]) ComputeCtx(ctx context.Context) (X, error) {
    return v.x, nil
}

// FromValue creates a promise that simply returns the provided argument and
// a nil error when computed.
func FromValue[X any](x X) P[X] {
    return value[X]{x: x}
}

// FromValueErr creates a promise that simply returns the provided argument or,
// if the provided error was non-nil, the provided error when computed.
func FromValueErr[X any](value X, err error) P[X] {
    if err != nil {
        return FromError[X](err)
    } else {
        return FromValue(value)
    }
}

type perror[X any] struct {err error}
func (e perror[X]) Compute() (X, error) {
    var zero X
    return zero, e.err
}
func (e perror[X]) ComputeCtx(ctx context.Context) (X, error) {
    var zero X
    return zero, e.err
}

// FromError creates a promise that simply returns the provided error when
// computed.
func FromError[X any](err error) P[X] {
    return perror[X]{err: err}
}

// Chain returns a new promise to compute function f on the result of promise
// p.
func Chain[X any, Y any](p P[X], f func(X) (Y, error)) P[Y] {
    return FromResultFuncCtx[Y](func(ctx context.Context) (Y, error) {
        v, err := p.ComputeCtx(ctx)
        if err != nil {
            var zero Y
            return zero, err
        }
        return f(v)
    })
}

/*
Example:

    a := promise.FromFunc(func() (int, error) {
        return 0, nil
    })
    half := promise.Chain(a, func(x int) (float64, error) {
        return float64(x) * 0.5, nil
    })
    double := promise.Chain(half, func(x float64) (float64, error) {
        return float64(x) * 2.0, nil
    })
    inverse := promise.Chain(double, func(x float64) (float64, error) {
        if x == 0.0 {
            return 0.0, fmt.Errorf("divide by zero error")
        }
        return 1.0 / x, nil
    })

    fmt.Printf("got %f\n", must.Result(inverse.Compute()))
 */
