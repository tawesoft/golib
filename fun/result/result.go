// Package result implements a Result{value, error} "sum type" that has a value
// only when error is nil.
//
// Note that in many cases, it is more idiomatic for a function to return a
// naked (value, error). Use [WrapFunc] to convert such a function to return
// a Result type.
//
// For examples, see the sibling "maybe" package.
package result

// Result is a (value, error) "sum type" that has a value only when error is
// nil.
type Result[V any] struct {
    Value V
    Error error
}

// Ok returns true if the Result's error is nil.
func (r Result[V]) Ok() bool {
    return r.Error == nil
}

// New returns a Result. It is syntax sugar for Result{value, error}. If error
// is a known constant, use [Some] or [Error] instead.
func New[V any](value V, err error) Result[V] {
    if err != nil { return Error[V](err) }
    return Some(value)
}

// Result returns a plain (value, error) tuple from a Result.
func (r Result[V]) Unpack() (V, error) {
    return r.Value, r.Error
}

// Else returns the Results's value (if error is nil), otherwise returns the
// provided argument instead.
func (r Result[V]) Else(v V) V {
    if r.Error == nil { return r.Value }
    return v
}

// Match examines a Result and returns true if it is not an error and if the
// provided predicate returns true for the value.
func (r Result[V]) Match(f func(v V) bool) bool {
    if r.Error != nil { return false }
    return f(r.Value)
}

// Filter examines a Result. If error is nil and f(value) returns true, returns
// Some(value), otherwise returns Error(error).
func (r Result[V]) Filter(f func(v V) bool) Result[V] {
    if r.Error != nil { return Error[V](r.Error) }
    if f(r.Value) { return r }
    return Error[V](r.Error)
}

// Must returns a Result's value. If the Result is an error, panics with the
// error.
func (r Result[V]) Must() V {
    if r.Error != nil {
        panic(r.Error)
    }
    return r.Value
}

// MustError panics if the Result is not an error.
func (r Result[V]) MustError() {
    if r.Error != nil {
        panic(r.Error)
    }
}

// Error returns a Result type that represents an error.
func Error[V any](err error) Result[V] {
    return Result[V]{Error: err}
}

// Some (a.k.a. "Just") returns a Result that contains a value with a nil error.
func Some[V any](value V) Result[V] {
    return Result[V]{
        Value: value,
        Error: nil,
    }
}

// WrapFunc converts a function of the form f(x) => (value, error) to the form
// f(x) => Result{value, error}.
func WrapFunc[A any, B any](
    f func(a A) (B, error),
) func(a A) Result[B] {
    return func(a A) Result[B] {
        return New(f(a))
    }
}

// UnwrapFunc converts a function of the form f(x) => Result{value, error} to
// the form f(x) => (value, error)
func UnwrapFunc[A any, B any](
    f func(a A) Result[B],
) func(a A) (B, error) {
    return func(a A) (B, error) {
        return f(a).Unpack()
    }
}

// Apply applies Maybe(f(x)) => y. That is, function f is itself a Result. If
// either a or f are errors, returns an Error. Otherwise, returns
// Some(f(a.Value)).
//
// If both a and f are errors, returns either error.
func Apply[A any, B any](
    a Result[A],
    f Result[func(a A) B],
) Result[B] {
    if a.Error != nil { return Error[B](a.Error) }
    if f.Error != nil { return Error[B](f.Error) }
    return Some(f.Value(a.Value))
}

// Map applies a function f(x) => y. If a is an error, returns an error.
// Otherwise, returns Some(f(a.Value)).
//
// For a function that applies f(x) => Maybe(y) instead, see [Then].
func Map[A any, B any](
    a Result[A],
    f func(a A) B,
) Result[B] {
    if a.Error != nil { return Error[B](a.Error) }
    return Some(f(a.Value))
}

// Then (a.k.a. "flatMap") applies a function f(x) => Maybe(y). If a is
// Nothing, returns Nothing, otherwise returns f(a.Value).
//
// For a function that applies f(x) => y instead, see [Map].
func Then[A any, B any](
    a Result[A],
    f func(a A) Result[B],
) Result[B] {
    if a.Error != nil { return Error[B](a.Error) }
    return f(a.Value)
}

// Lift converts a function of the form f(a) => b to the form f(a) => Result(b).
func Lift[A any, B any](
    f func(a A) B,
) func(a A) Result[B] {
    return func(a A) Result[B] {
        return Some(f(a))
    }
}
