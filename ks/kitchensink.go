// SPDX-License-Identifier: MIT-0
// x-doc-short-desc: misc helpful things
// x-doc-stable: no

// Package ks ("kitchen sink") implements assorted helpful things that don't
// fit anywhere else.
package ks

// Number defines anything you can perform arithmetic with using standard Go
// operators (like a + b, or a ^ b). A package should use it to define its own
// number type (e.g. type Number = ks.Number) rather than use the type directly
// from here in its exported interface.
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
