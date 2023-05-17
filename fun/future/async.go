package future

import (
    "context"

    "github.com/tawesoft/golib/v2/fun/partial"
    "github.com/tawesoft/golib/v2/fun/promise"
    "github.com/tawesoft/golib/v2/fun/result"
    "github.com/tawesoft/golib/v2/fun/slices"
)

type async[T any] struct {
    context context.Context
    cancel context.CancelFunc
    channel chan result.R[T]
}

func start[T any](ctx context.Context, promise promise.P[T], channel chan result.R[T]) {
    value := result.New(promise.ComputeCtx(ctx))
    for {
        select {
            case <- ctx.Done(): return
            default: channel <- value
        }
    }
    close(channel)
}

// NewAsync creates a new future from a promise, and begins computing that
// promise asynchronously in a new goroutine.
func NewAsync[T any](ctx context.Context, promise promise.P[T]) F[T] {
    ctxWithCancel, cancel := context.WithCancel(ctx)
    f := async[T]{
        context: ctxWithCancel,
        cancel: cancel,
        channel: make(chan result.R[T]),
    }
    go start(ctxWithCancel, promise, f.channel)
    return f
}

// NewAsyncs is like [NewAsync], but accepts a slice of promises and returns a
// slice of futures.
func NewAsyncs[T any](ctx context.Context, xs []promise.P[T]) []F[T] {
    return slices.Map(
        partial.Left2(NewAsync[T])(ctx),
        xs,
    )
}

func (f async[T]) Collect() (result T, err error) {
    return f.CollectCtx(context.TODO())
}

func (f async[T]) CollectCtx(ctx context.Context) (result T, err error) {
    select {
        case <- ctx.Done():
            err = ctx.Err()
            return
        case <- f.context.Done():
            err = f.context.Err()
            return
        case r := <- f.channel:
            result, err = r.Unpack()
            return
    }
}

func (f async[T]) Stop() {
    f.cancel()
}

func (f async[T]) Peek() (result T, err error) {
    select {
        case <- f.context.Done():
            err = f.context.Err()
            return
        case r := <- f.channel:
            result, err = r.Unpack()
            return
        default:
            err = NotReady
            return
    }
}
