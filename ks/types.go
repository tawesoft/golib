package ks

// Pair is any Key, Value pair. Type K is any type that would be suitable as a
// KeyType in a Go [builtin.map].
type Pair[K comparable, V any] struct {
    Key   K
    Value V
}

// Item is any Key, Value pair. Type K is any type that would be suitable as a
// KeyType in a Go [builtin.map].
//
// A downstream package should use this to define its own number type (e.g.
// type Item[K comparable, V any] ks.Item[K, V]) rather than use the type
// directly from here in its exported interface.
type Item[K comparable, V any] struct {
    Key   K
    Value V
}

// Result is a (value, error) "sum type" that has a value only when Error is
// nil.
//
// Note that in many cases, it is more idiomatic for a function to return a
// naked (value, error). The Result type is more useful in iterators.
type Result[V any] struct {
    Value V
    Error error
}

// ResultFunc accepts a function that returns a naked (value, error) and
// returns a function that returns a Result{value, error} instead.
func ResultFunc[V any](f func() (V, error)) func() Result[V] {
    return func() Result[V] {
        v, err := f()
        return Result[V]{v, err}
    }
}

// Maybe is a (value, ok) "sum type" that has a value only when Ok is
// true.
//
// Note that in many cases, it is more idiomatic for a function to return a
// naked (value, ok). The Maybe type is more useful in iterators.
type Maybe[V any] struct {
    Value V
    Ok bool
}

// MaybeFunc accepts a function that returns a naked (value, exists) and
// returns a function that returns a Maybe{value, exists} instead.
func MaybeFunc[V any](f func() (V, bool)) func() Maybe[V] {
    return func() Maybe[V] {
        v, ok := f()
        return Maybe[V]{v, ok}
    }
}
