package future

import (
    "context"

    "github.com/tawesoft/golib/v2/fun/maybe"
    "github.com/tawesoft/golib/v2/fun/promise"
    "github.com/tawesoft/golib/v2/fun/result"
    "github.com/tawesoft/golib/v2/fun/slices"
    "github.com/tawesoft/golib/v2/operator"
)

type sync[T any] struct {
    promise promise.P[T]
    result maybe.M[result.R[T]]
}

// NewSync creates a new future to be run synchronously based on a given
// promise to compute a value.
//
// It's methods are not safe for concurrent access.
func NewSync[T any](promise promise.P[T]) F[T] {
    return &sync[T]{promise: promise,}
}

// NewSyncs is like [NewSync], but accepts a slice of promises and returns a
// slice of futures.
func NewSyncs[T any](xs []promise.P[T]) []F[T] {
    return slices.Map(NewSync[T], xs)
}

func (f *sync[T]) Collect() (T, error) {
    return f.CollectCtx(context.TODO())
}

func (f *sync[T]) CollectCtx(ctx context.Context) (T, error) {
    if !f.result.Ok {
        r, err := f.promise.ComputeCtx(ctx)
        f.result = maybe.Some(result.New(r, err))
    }
    return f.result.Value.Unpack()
}

func (f *sync[T]) Peek() (T, error) {
    if !f.result.Ok {
        return operator.Zero[T](), NotReady
    }
    return f.result.Value.Unpack()
}

func (f *sync[T]) Stop() {
    f.result = maybe.Some(result.Error[T](context.Canceled))
}
