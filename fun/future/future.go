// Package future implements "Futures", which represent a placeholder handle
// for a value that may not yet be ready, but is (eventually) computed by a
// promise.
package future

import (
    "context"
    "errors"

    "github.com/tawesoft/golib/v2/fun/promise"
)

// F represents a "future", a placeholder handle backed by a promise to
// compute a value at some later point.
//
// It is up to the implementation of the interface if the promise will run
// synchronously or asynchronously, or if it is safe to collect the future
// using concurrent code. See [NewSync] and [NewAsync] for synchronous and
// concurrent implementations, respectively.
//
// Collect and CollectCtx return the value or error returned by the computed
// promise, blocking if necessary. A promise is computed exactly once. In
// synchronous code, this happens at the time of the first call to Collect or
// CollectCtx. The cached value or error is returned to each subsequent call to
// Collect or CollectCtx.
//
// Peek is like Collect, but if the value is not yet available, returns
// immediately with the error [NotReady]. In the case of [Sync], Peek always
// returns with [NotReady] if the value has yet to be computed (i.e. by a call
// to Collect or CollectCtx).
//
// The error return value of Collect, Peek, and CollectCtx may include context
// errors such as [context.Cancelled].
//
// The "Ctx" method variants support cancellation only so far as the underlying
// promise supports being cancelled.
//
// Stop indicates that the future is no longer needed and its resources can be
// released. In the case of asynchronous code, this includes terminating the
// backing goroutine so that future calls to Collect, CollectCtx, and Peek
// return [context.Cancelled]. Futures created with [NewSync] do not need to be
// closed, but it is not an error to do so. An asynchronous future may not
// actually stop until the underlying promise has finished computing or
// accepted a cancellation signal. It is not an error to stop a future multiple
// times.
type F[T any] interface {
    Collect() (T, error)
    CollectCtx(ctx context.Context) (T, error)
    Peek() (T, error)
    Stop()
}

// NotReady is the error returned by the Peek methods on a future when the
// computed value or error is not yet available.
var NotReady = errors.New("future.NotReady")

// ForEach applies function f(x) to every future that is both computed
// and whose promise did not return an error i.e. for every successful Peek.
func ForEach[T any](f func(x T), xs []F[T]) {
    for _, x := range xs {
        value, err := x.Peek()
        if err != nil { continue}
        f(value)
    }
}

// CollectAll returns a promise to compute the slice of the values of the input
// futures, stopping at the first error.
func CollectAll[T any](xs []F[T]) promise.P[[]T] {
    return CollectAllCtx(context.TODO(), xs)
}

// CollectAllCtx is like [CollectAll] but uses the given context while
// computing the promises.
func CollectAllCtx[T any](ctx context.Context, xs []F[T]) promise.P[[]T] {
    return promise.FromResultFunc(func () ([]T, error) {
        values := make([]T, 0, len(xs))
        for _, x := range xs {
            value, err := x.CollectCtx(ctx)
            if err != nil {
                return nil, err
            }
            values = append(values, value)
        }
        return values, nil
    })
}
