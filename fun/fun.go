// SPDX-License-Identifier: MIT
// x-doc-short-desc: higher-order functions for various data structures
// x-doc-stable: candidate

// Package fun implements some common higher-order functions (like Map, Filter,
// Reduce) on slices. For lazy evaluation, see [golib/lazy] instead.
//
// Unlike many existing generic functional libraries for Go, our map function
// can also map to different types.
//
// To avoid confusion, in this package the Go map data structure is referred to
// as a "dict". "Map" in this package refers to the higher-order function "map"
// used in functional programming.
//
// This is a placeholder package. All functionality is implemented by packages
// in the subdirectories.
package fun

import (
    _ "github.com/tawesoft/golib/v2/fun/slice"
)
