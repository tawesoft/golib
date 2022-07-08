// Package ks ("kitchen sink") implements assorted helpful things that don't
// fit anywhere else.
package ks

import (
    "fmt"
    "reflect"
    "strings"
    "unicode/utf8"

    "golang.org/x/exp/utf8string"
)

// Catch calls the input function f. If successful, Catch passes on the return
// value from f and also returns a nil error. If f panics, Catch recovers from
// the panic and returns a non-nil error.
//
// If the panic raised by f contains is of type error, the returned error
// is wrapped once.
//
// The opposite of Catch is [Must]: Catch(Must(os.Open(""))
func Catch[X any](f func() X) (x X, err error) {
    defer func() {
        if r := recover(); r != nil {
            if rErr, ok := r.(error); ok {
                err = fmt.Errorf("caught panic: %w", rErr)
            } else {
                err = fmt.Errorf("caught panic: %v", r)
            }
        }
    }()

    return f(), nil
}

// initsh casts an int to comparable. If the comparable is not an integer type,
// this will panic.
func intish[T comparable](i int) T {
    var t T
    ref := reflect.ValueOf(&t).Elem()
    ref.SetInt(int64(i))
    return t
}

// Range calls some function f(k, v) => bool over any [Rangeable]. If the
// return value of f is false, the iteration stops.
//
// This is roughly equivalent to "k, v := range(x); if !f(x) { break }".
//
// Caution: invalid key types will panic at runtime. The key type must be int
// for any type other than a map. See [Rangeable] for details. In a channel,
// the key is always zero.
func Range[K comparable, V any, R Rangeable[K, V]](
    f func(K, V) bool,
    r R,
) {
    switch ref := reflect.ValueOf(r); ref.Kind() {
        case reflect.Array:
            for i := 0; i < ref.Len(); i++ {
                k := intish[K](i)
                v := ref.Index(i).Interface().(V)
                if !f(k, v) { break }
            }
        case reflect.Chan:
            for {
                x, ok := ref.Recv()
                if !ok { break }
                v := x.Interface().(V)
                if !f(intish[K](0), v) { break }
            }
        case reflect.Map:
            iter := ref.MapRange()
            for iter.Next() {
                k, v := iter.Key().Interface().(K), iter.Value().Interface().(V)
                if !f(k, v) { break }
            }
        case reflect.Slice:
            for i := 0; i < ref.Len(); i++ {
                k := intish[K](i)
                v := ref.Index(i).Interface().(V)
                if !f(k, v) { break }
            }
        case reflect.String:
            for i := 0; i < ref.Len(); i++ {
                k := intish[K](i)
                v := ref.Index(i).Interface().(V)
                if !f(k, v) { break }
            }
    }
}

// CheckedRange calls fn(k, v) => error for each key, value in the input slice,
// but halts if an error is returned at any point. If so, it returns the key
// and value being examined at the time of the error, and the encountered
// error, or a nil error otherwise.
func CheckedRange[K comparable, V any, R Rangeable[K, V]](
    fn func(k K, v V) error,
    r R,
) (K, V, error) {
    var (k K; v V; err error)
    f := func(k2 K, v2 V) bool {
        k, v = k2, v2
        err = fn(k, v);
        return err == nil
    }
    Range(f, r)
    return k, v, err
}

// CheckedRangeValue is like [CheckedRange], but calls fn(value), not fn(key,
// value), and returns only (value, error), not (key, value, error).
func CheckedRangeValue[K comparable, V any, R Rangeable[K, V]](
    fn func(v V) error,
    r R,
) (V, error) {
    var (v V; err error)
    f := func(_ K, v2 V) bool {
        v = v2
        err = fn(v)
        return err == nil
    }
    Range(f, r)
    return v, err
}

// IfThenElse returns a value based on a boolean condition, q. Iff q is true,
// returns the ifTrue. Iff q is false, returns ifFalse. This [IfThenElse
// expression] (as distinct from If-Then-Else statements) is much like the
// ternary operator in some other languages, however it is not short-circuited
// and both arguments are evaluated.
//
// For a lazily-evaluated version, see [lazy.IfThenElse].
//
// [IfThenElse expression]: https://en.wikipedia.org/wiki/Conditional_(computer_programming)#If%E2%80%93then%E2%80%93else_expressions
func IfThenElse[X any] (
    q       bool,
    ifTrue  X,
    ifFalse X,
) X {
    if q {
        return ifTrue
    } else {
        return ifFalse
    }
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

// Must accepts a (value, err) tuple as input and panics if err != nil,
// otherwise returns value. The error raised by panic is wrapped in another
// error.
//
// For example, Must(os.Open("doesnotexist")) panics with an error like
// "unexpected error in Must[*os.File]: open doesnotexist: no such file or
// directory". Must(os.Open("filethatexists")) returns a pointer to an
// [os.File].
//
// The opposite of Must is [Catch]: Catch(Must(os.Open(""))
func Must[T any](t T, err error) T {
    if err != nil {
        panic(fmt.Errorf("unexpected error in Must[%T]: %w", t, err))
    }
    return t
}

// MustFunc accepts a function that takes an input of type X, where that
// function then returns a (value Y, err) tuple. Must then returns a function
// that panics if the returned err != nil, otherwise returns value Y. The
// returned error is wrapped in another error.
//
// For example, MustFunc(os.Open) returns a function (call this f).
// f("doesnotexist") panics with an error (like [Must]), and
// f("filethatexists") returns a pointer to an [os.File].
func MustFunc[X any, Y any](
    f func (x X) (Y, error),
) func (x X) Y {
    return func(x X) Y {
        return Must(f(x))
    }
}

// MustInit calls f, a function that returns a (value, error) tuple. If the
// error is nil, returns the value. Otherwise, returns the default value d.
// Intended to be used to initialise package-level constants.
func MustInit[K any](f func () (value K, err error), d K) K {
    if v, err := f(); err == nil {
        return v
    } else {
        return d
    }
}

// Never signifies code that should never be reached. It raises a panic when
// called.
func Never() {
    panic("this should never happen")
}

// Zero returns the zero value for any type.
func Zero[T any]() T {
    var t T
    return t
}

// Rangeable defines any type of value x where it is possible to range over
// using "for k, v := range x" or "v := range x" (in the case of a channel,
// only "v := range x" is permitted). For every Rangeable other than a map,
// K must always be int.
type Rangeable[K comparable, V any] interface {
    ~string | ~map[K]V | ~[]V | chan V
}

// WrapBlock word-wraps a whitespace-delimited string to a given number of
// columns. The column length is given in runes (Unicode code points), not
// bytes.
//
// This is a simple implementation without any configuration options, designed
// for circumstances such as quickly wrapping a single error message for
// display.
//
// Save for bug fixes, the output of this function for any given input is
// frozen and will not be changed in future. This means you can reliably test
// against the return value of this function without your tests being brittle.
//
// Caveat: Single words longer than the column length will be truncated.
//
// Caveat: all whitespace, including existing new lines, is collapsed. An input
// consisting of multiple paragraphs will be wrapped into a single word-wrapped
// paragraph.
//
// Caveat: assumes all runes in the input string represent a glyph of length
// one. Whether this is true or not depends on how the display and font treats
// different runes. For example, some runes where [Unicode.IsGraphic] returns
// false might still be displayed as a special escaped character. Some letters
// might be displayed wider than usual, even in a monospaced font.
func WrapBlock(message string, columns int) string {
    var atoms = strings.Fields(strings.TrimSpace(message))
    var sb = strings.Builder{}
    var currentLength int

    if columns <= 0 { return "" }

    for i, atom := range atoms {
        isLast := (i + 1 == len(atoms))
        atomLength := utf8.RuneCountInString(atom)

        // special case for an atom longer than a whole line
        if (currentLength == 0) && (atomLength >= columns) {
            truncated := utf8string.NewString(atom).Slice(0, columns)
            sb.WriteString(truncated)
            if !isLast { sb.WriteByte('\n') }
            currentLength = 0
            continue
        }

        // will overflow?
        if currentLength + atomLength + 1 > columns {
            sb.WriteByte('\n')
            currentLength = 0
        }

        // mid-line?
        if currentLength > 0 {
            sb.WriteByte(' ')
            currentLength += 1
        }

        sb.WriteString(atom)
        currentLength += atomLength
    }

    return sb.String()
}
