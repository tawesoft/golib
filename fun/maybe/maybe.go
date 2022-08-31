// Package maybe implements a Maybe{value, ok} "sum type" that has a value only
// when ok is true.
//
// Note that in many cases, it is more idiomatic for a function to return a
// naked (value, ok). Use [WrapFunc] to convert such a function to return
// a Maybe type.
package maybe

import (
    "fmt"
)

// Maybe is a (value, ok) "sum type" that has a value only when ok is
// true.
type Maybe[V any] struct {
    Value V
    Ok bool
}

// New returns a Maybe. It is syntax sugar for Maybe{value, ok}. If ok is
// a known constant, use [Some] or [Nothing] instead.
func New[V any](value V, ok bool) Maybe[V] {
    if !ok { return Nothing[V]() }
    return Some(value)
}

// Unpack returns a plain (value, ok) tuple from a Maybe.
func (m Maybe[V]) Unpack() (V, bool) {
    return m.Value, m.Ok
}

// Else returns the Maybe's value (if ok), otherwise returns the provided
// argument instead.
func (m Maybe[V]) Else(v V) V {
    if m.Ok { return m.Value }
    return v
}

// Filter examines a Maybe. If ok and f(value) returns true, returns
// Some(value), otherwise returns Nothing.
func (m Maybe[V]) Filter(f func(v V) bool) Maybe[V] {
    if !m.Ok { return Nothing[V]() }
    if f(m.Value) { return m }
    return Nothing[V]()
}

// Match examines a Maybe and returns true if ok and if the provided predicate
// returns true for the value.
func (m Maybe[V]) Match(f func(v V) bool) bool {
    if !m.Ok { return false }
    return f(m.Value)
}

// Must returns a Maybe's value. If the Maybe is not ok, panics.
func (m Maybe[V]) Must() V {
    if !m.Ok {
        panic(fmt.Sprintf("Maybe[%T].Must called on missing value.", m))
    }
    return m.Value
}

// MustNot panics if the Maybe is ok.
func (m Maybe[V]) MustNot() {
    if m.Ok {
        panic(fmt.Sprintf("Maybe[%T].MustNot called but value present.", m))
    }
}

// Nothing returns a (typed) Maybe that has no value.
func Nothing[V any]() Maybe[V] {
    return Maybe[V]{}
}

// Some (a.k.a. "Just") returns a Maybe that contains a value.
func Some[V any](value V) Maybe[V] {
    return Maybe[V]{
        Value: value,
        Ok:    true,
    }
}

// WrapFunc converts a function of the form f(x) => (value, ok) to the form
// f(x) => Maybe{value, ok}.
func WrapFunc[A any, B any](
    f func(a A) (B, bool),
) func(a A) Maybe[B] {
    return func(a A) Maybe[B] {
        return New(f(a))
    }
}

// UnwrapFunc converts a function of the form f(x) => Maybe{value, ok} to the
// form f(x) => (value, ok)
func UnwrapFunc[A any, B any](
    f func(a A) Maybe[B],
) func(a A) (B, bool) {
    return func(a A) (B, bool) {
        return f(a).Unpack()
    }
}

// Apply applies a function f(x) => y, but function f is itself a Maybe. If
// either a or f are Nothing, returns Nothing. Otherwise, returns
// Some(f(a.Value)).
//
// This is useful for working with partial function application.
//
//     doubler := partial.Left2(mul, 2)
//     maybe.Apply(Some(x), maybe.Map(Some(y), doubler))
//
// For a function that applies f(x) => Maybe(y) instead, see [FlatApply].
func Apply[A any, B any](
    a Maybe[A],
    f Maybe[func(a A) B],
) Maybe[B] {
    if !a.Ok { return Nothing[B]() }
    if !f.Ok { return Nothing[B]() }
    return Some(f.Value(a.Value))
}

// Map applies a function f(x) => y. If a is Nothing, returns Nothing.
// Otherwise, returns Some(f(a.Value)).
//
// For a function that applies f(x) => Maybe(y) instead, see [FlatMap].
func Map[A any, B any](
    a Maybe[A],
    f func(a A) B,
) Maybe[B] {
    if !a.Ok { return Nothing[B]() }
    return Some(f(a.Value))
}

// FlatMap applies a function f(x) => Maybe(y). If a is Nothing, returns
// Nothing, otherwise returns f(a.Value).
//
// This is called "flat" because applying function f with the ordinary map
// would give Maybe(Maybe(y)).
//
// For a function that applies f(x) => y instead, see [Map].
func FlatMap[A any, B any](
    a Maybe[A],
    f func(a A) Maybe[B],
) Maybe[B] {
    if !a.Ok { return Nothing[B]() }
    return f(a.Value)
}

// FlatApply applies Maybe(f(x)) => Maybe(y). That is, function f is itself a
// Maybe. If either a or f are Nothing, returns Nothing. Otherwise, returns
// f(a.Value).
//
// This is useful for working with partial function application.)
//
// This is called "flat" because applying function f with the ordinary apply
// would give Maybe(Maybe(y)).
//
// For a function that applies f(x) => y instead, see [Apply].
func FlatApply[A any, B any](
    a Maybe[A],
    f Maybe[func(a A) Maybe[B]],
) Maybe[B] {
    if !a.Ok { return Nothing[B]() }
    if !f.Ok { return Nothing[B]() }
    return f.Value(a.Value)
}

// Lift converts a function of the form f(a) => b to the form f(a) => Maybe(b).
func Lift[A any, B any](
    f func(a A) B,
) func(a A) Maybe[B] {
    return func(a A) Maybe[B] {
        return Some(f(a))
    }
}
