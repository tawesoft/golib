// Package either implements a simple generic "Either" type that can represent
// exactly one value out of two options.
package either

import (
    "fmt"
)

// E represents a type that can hold either a value "a" of type A or a
// value "b" of type B.
type E[A any, B any] struct {
    a A
    b B
    index byte // 'a' or 'b'
}

// Pack returns an E that contains eotjer a value "a" of type A (if index ==
// 'a'), a value "b" of type B (if index == 'b'), or panics if index is not 'a'
// or 'b'.
func Pack[A any, B any](a A, b B, index byte) E[A, B] {
    if index == 'a' {
        return E[A, B]{a: a, index: index}
    } else if index == 'b' {
        return E[A, B]{b: b, index: index}
    } else {
        panic(fmt.Errorf("either: Pack[%T, %T](): invalid index value", a, b))
    }
}

// Unpack returns the components of an E. The last return value is a
// discriminator with the value 'a' or 'b' representing the existence of the
// value "a" of type A or the value "b" of type B.
func (e E[A, B]) Unpack() (A, B, byte) {
    return e.a, e.b, e.index
}

// A returns a new E that holds a value "a" of Type A.
func A[A any, B any](a A) E[A, B] {
    return E[A, B]{
        a: a,
        index: 'a',
    }
}

// B returns a new E that holds a value "b" of type B.
func B[A any, B any](b B) E[A, B] {
    return E[A, B]{
        b: b,
        index: 'b',
    }
}

// A returns the value "a" of type A and true if the E contains that value,
// or the zero value and false otherwise.
func (e E[A, B]) A() (result A, ok bool) {
    if e.index == 'a' {
        result = e.a
        ok = true
    }
    return
}

// B returns the value "b" of type B and true if the E contains that value,
// or the zero value and false otherwise.
func (e E[A, B]) B() (result B, ok bool) {
    if e.index == 'b' {
        result = e.b
        ok = true
    }
    return
}
