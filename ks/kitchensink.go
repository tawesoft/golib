// Package ks ("kitchen sink") implements assorted helpful things that don't
// fit anywhere else.
package ks

import (
    "fmt"
    "io"
    "reflect"
    "strings"
    "testing"
    "time"
    "unicode/utf8"

    "github.com/tawesoft/golib/v2/must"
    "golang.org/x/exp/utf8string"
)

// Assert panics if the value is not true. Optionally, follow with a
// printf-style format string and arguments.
func Assert(q bool, args ... interface{}) {
    if q { return }

    if len(args) == 0 {
        panic(fmt.Errorf("assertion error"))
    } else {
        panic(fmt.Errorf("assertion error: " + args[0].(string), args[1:]...))
    }
}

func Cast[X any, Y any](x X) Y {
    ref := reflect.ValueOf(&x).Elem()
    return ref.Interface().(Y)
}

/*
// CatchFunc returns a function that wraps input function f. When called, if f
// is successful, CatchFunc passes on the return value from f and also returns
// a nil error. If f panics, CatchFunc recovers from the panic and returns a
// non-nil error.
//
// If the panic raised by f contains is of type error, the returned error
// is wrapped once.
//
// The opposite of CatchFunc is [MustFunc] - e.g.
//   CatchFunc(MustFunc(os.Open))("example.txt")
func CatchFunc[X any](f func() X) func() (x X, err error) {
    return func() (x X, err error) {
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
}
*/

// FirstNonZero returns the first argument that isn't equal to the zero value
// of T, or otherwise the zero value.
func FirstNonZero[T comparable](args ... T) T {
    var zero T
    for _, a := range args {
        if a == zero { continue }
        return a
    }
    return zero
}

// Identity implements the function f(x) => x.
func Identity[X any](x X) X {
    return x
}

// IfThenElse returns a value based on a boolean condition, q. Iff q is true,
// returns the ifTrue. Iff q is false, returns ifFalse. This [IfThenElse
// expression] (as distinct from If-Then-Else statements) is much like the
// ternary operator in some other languages, however it is not short-circuited
// and both arguments are evaluated.
//
// For a lazily-evaluated version, see [IfThenElseLazy].
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

// IfThenElseLazy returns a lazily-evaluated value based on a boolean
// condition, q. Iff q is true, returns the return value of ifTrue(). Iff q is
// false, returns the return value of ifFalse(). This [IfThenElse expression]
// (as distinct from If-Then-Else statements) is much like the ternary operator
// in some other languages.
//
// For a non-lazy version, see [IfThenElse].
//
// [IfThenElse expression]: https://en.wikipedia.org/wiki/Conditional_(computer_programming)#If%E2%80%93then%E2%80%93else_expressions
func IfThenElseLazy[X any] (
    q       bool,
    ifTrue  func() X,
    ifFalse func() X,
) X {
    if q {
        return ifTrue()
    } else {
        return ifFalse()
    }
}

// In returns true if x equals any of the following arguments.
func In[X comparable](x X, xs ... X) bool {
    for _, i := range xs {
        if x == i { return true }
    }
    return false
}

/*
// Must accepts a (value, err) tuple as input and panics if err != nil,
// otherwise returns value. The error raised by panic is wrapped in another
// error.
//
// For example, Must(os.Open("doesnotexist")) panics with an error like
// "unexpected error in Must[*os.File]: open doesnotexist: no such file or
// directory". Must(os.Open("filethatexists")) returns a pointer to an
// [os.File].
func Must[T any](t T, err error) T {
    if err != nil {
        panic(fmt.Errorf("unexpected error in Must[%T]: %w", t, err))
    }
    return t
}

// MustOk accepts a (value, ok bool) tuple as input and panics if ok is false,
// otherwise returns value.
func MustOk[T any](t T, ok bool) T {
    if !ok {
        panic(fmt.Errorf("unexpected error in MustOk[%T]: not ok", t))
    }
    return t
}
*/

// MustFunc accepts a function f(x) => (y, err) and returns a function
// g(x) => y that panics if f(x) returned an error, otherwise returns y.
//
// The opposite of MustFunc is CatchFunc - e.g.
//   CatchFunc(MustFunc(os.Open))("example.txt")
/*
func MustFunc[X any, Y any](
    f func (x X) (Y, error),
) func (x X) Y {
    return func(x X) Y {
        return Must(f(x))
    }
}
*/

// MustInit calls f, a function that returns a (value, error) tuple. If the
// error is nil, returns the value. Otherwise, returns the default value d.
// Intended to be used to initialise package-level constants.
/*
func MustInit[K any](f func () (value K, err error), d K) K {
    if v, err := f(); err == nil {
        return v
    } else {
        return d
    }
}
TODO replace with Else
*/


// Never signifies code that should never be reached. It raises a panic when
// called.
//
// Deprecated: use [must.Never] instead.
func Never(args ... interface{}) {
    if len(args) > 0 {
        panic(fmt.Errorf("this should never happen: " + args[0].(string), args[1:]...))
    }
    panic(fmt.Errorf("this should never happen"))
}

// initsh casts an int to comparable. If the comparable is not an integer type,
// this will panic.
func intish[T any](i int) T {
    var t T
    ref := reflect.ValueOf(&t).Elem()
    ref.SetInt(int64(i))
    return t
}

// Range calls some function f(k, v) => err over any [Rangeable] of (K, V)s. If
// the return value of f is not nil, the iteration stops immediately, and
// returns (k, err) for the given k. Otherwise, returns (zero, nil).
//
// Caution: invalid key types will panic at runtime. The key type must be int
// for any type other than a map. See [Rangeable] for details. In a channel,
// the key is always zero.
func Range[K comparable, V any, R Rangeable[K, V]](
    f func(K, V) error,
    r R,
) (K, error) {
    switch ref := reflect.ValueOf(r); ref.Kind() {
        case reflect.Array:
            for i := 0; i < ref.Len(); i++ {
                k := intish[K](i)
                v := ref.Index(i).Interface().(V)
                if err := f(k, v); err != nil { return k, err }
            }
        case reflect.Chan:
            for {
                x, ok := ref.Recv()
                if !ok { break }
                k := intish[K](0)
                v := x.Interface().(V)
                if err := f(k, v); err != nil { return k, err }
            }
        case reflect.Map:
            iter := ref.MapRange()
            for iter.Next() {
                k, v := iter.Key().Interface().(K), iter.Value().Interface().(V)
                if err := f(k, v); err != nil { return k, err }
            }
        case reflect.Slice:
            for i := 0; i < ref.Len(); i++ {
                k := intish[K](i)
                v := ref.Index(i).Interface().(V)
                if err := f(k, v); err != nil { return k, err }
            }
        case reflect.String:
            runes := 0
            for i := 0; i < ref.Len(); i++ {
                k := intish[K](runes)
                runes++
                v := ref.Index(i).Interface().(byte)

                if utf8.FullRune([]byte{v}) {
                    if err := f(k, intish[V](int(v))); err != nil { return k, err }
                    continue
                }

                if utf8.RuneStart(v) {
                    buf := [4]byte{v, 0, 0, 0}
                    for j := 0; j < 3; j++ {
                        //if i + j + 1 >= ref.Len() { break }
                        v = ref.Index(i+j+1).Interface().(byte)
                        buf[1+j] = v
                        if utf8.FullRune(buf[0:2+j]) {
                            n, _ := utf8.DecodeRune(buf[0:2+j])
                            i += j + 1
                            if err := f(k, intish[V](int(n))); err != nil { return k, err }
                            break
                        }
                    }
                } else {
                    return Zero[K](), fmt.Errorf("Unicode error at byte %d %x", i, v)
                }
            }
    }
    return Zero[K](), nil
}

// Rangeable defines any type of value x where it is possible to range over
// using "for k, v := range x" or "v := range x" (in the case of a channel,
// only "v := range x" is permitted). For every Rangeable other than a map,
// K must always be int.
type Rangeable[K comparable, V any] interface {
    ~string | ~map[K]V | ~[]V | chan V
}

// Reserve grows a slice to fit at least size extra elements. Like the builtin
// append, it may return an updated slice.
func Reserve[T any](xs []T, size int) []T {
    // https://github.com/golang/go/wiki/SliceTricks#extend-capacity
    if cap(xs) - len(xs) < size {
        return append(make([]T, 0, len(xs) + size), xs...)
    }
    return xs
}

// TestCompletes executes f (in a goroutine), and blocks until either f returns,
// or the provided duration has elapsed. In the latter case, calls t.Errorf to
// fail the test. Provide optional format string and arguments to add
// context to the test error message.
func TestCompletes(t *testing.T, duration time.Duration, f func(), args ... interface{}) {
    done := make(chan struct{}, 1)
    timeout := time.After(duration)
    go func() {
        f()
        done <- struct{}{}
    }()

    select {
        case <-done: // OK
        case <-timeout:
            if len(args) > 0 {
                t.Errorf("test timed out after "+duration.String()+": " + args[0].(string), args[1:]...)
            } else {
                t.Errorf("test timed out after %s", duration.String())
            }
    }
}

// Zero returns the zero value for any type.
func Zero[T any]() T {
    var t T
    return t
}

// WithCloser.
//
// Deprecated: use [with.Closer].
func WithCloser[T io.Closer](opener func() (T, error), do func(v T) error) error {
    var zero T

    f, err := opener()
    if err != nil { return fmt.Errorf("WithCloser[%T] open error: %w", zero, err) }

    doer := must.CatchFunc(func() error { return do(f) })
    err, panicErr := doer()
    if err != nil {
        err = fmt.Errorf("WithCloser[%T] error: %w", zero, err)
    } else if panicErr != nil {
        err = fmt.Errorf("WithCloser[%T] error: panic: %w", zero, panicErr)
    }

    errClose := f.Close()
    if errClose != nil {
        err = fmt.Errorf("WithCloser[%T] close error: %v; %w", zero, errClose, err)
    }

    return err
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
