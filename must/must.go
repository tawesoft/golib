// Package must implements
package must

import (
    "fmt"
)

// Result accepts a (value, err) tuple as input and panics if err != nil,
// otherwise returns value. The error raised by panic is wrapped in another
// error.
//
// For example, must.Result(os.Open("doesnotexist")) may panic with an error
// like "unexpected error in must.Result[*os.File]: open doesnotexist: no such
// file or directory". On success, returns *os.File
func Result[T any](t T, err error) T {
    if err != nil {
        panic(fmt.Errorf("error in must.Result[%T]: %w", t, err))
    }
    return t
}

// Ok accepts a (value, ok) tuple as input and panics if ok is false, otherwise
// returns value.
func Ok[T any](t T, ok bool) T {
    if ok == false {
        panic(fmt.Errorf("error in must.Ok[%T]: not ok", t))
    }
    return t
}

// Check panics if the error is not nil. Otherwise, it returns a nil error (so
// that it is convenient to chain).
func Check(err error) error {
    if err != nil {
        panic(fmt.Errorf("must.Check: unexpected error: %w", err))
    }
    return nil
}

// CatchFunc takes a function f() => x that may panic, and instead returns a
// function f() => (x, error).
func CatchFunc[X any](f func() X) func() (x X, err error) {
    return func() (x X, err error) {
        defer func() {
            if r := recover(); r != nil {
                if rErr, ok := r.(error); ok {
                    err = fmt.Errorf("must.CatchFunc: caught panic: %w", rErr)
                } else {
                    err = fmt.Errorf("must.CatchFunc: caught panic: %v", r)
                }
            }
        }()

        return f(), nil
    }
}

// Func takes a function f() => (x, error), and returns a function f() => x
// that may panic in the event of error.
func Func[X any](
    f func () (X, error),
) func () X {
    return func() X {
        return Result(f())
    }
}
