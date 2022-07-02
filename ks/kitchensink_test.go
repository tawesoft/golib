package ks_test

// CONTRIBUTORS: keep tests in alphabetical order, but with examples grouped
// first.

import (
    "errors"
    "fmt"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/v2/ks"
)

func ExampleZero() {

    type thing struct {
        number int
        phrase string
    }

    fmt.Printf("The zero value is %+v\n", ks.Zero[thing]())
    fmt.Printf("The zero value is %+v\n", ks.Zero[int32]())

    // Output:
    // The zero value is {number:0 phrase:}
    // The zero value is 0
}

func TestCatch(t *testing.T) {
    {
        failingFunctionWithString := func() string {
            panic("oops")
            return "something"
        }

        x, err := ks.Catch[string](failingFunctionWithString)
        assert.Equal(t, "", x)
        assert.Error(t, err)
        assert.Equal(t, nil, errors.Unwrap(err))
    }

    {
        failingFunctionWithError := func() string {
            panic(fmt.Errorf("oops"))
            return "something"
        }

        x, err := ks.Catch[string](failingFunctionWithError)
        assert.Equal(t, "", x)
        assert.Error(t, err)
        wrapped := errors.Unwrap(err)
        assert.Error(t, wrapped)
        assert.Equal(t, "oops", wrapped.Error())
    }
}

func TestMust(t *testing.T) {
    successfulFunction := func() (string, error) {
        return "success", nil
    }

    assert.NotPanics(t, func() {
        ks.Must(successfulFunction())
    })

    failingFunction := func() (string, error) {
        return "", fmt.Errorf("oops")
    }

    assert.Panics(t, func() {
        ks.Must(failingFunction())
    })
}

func TestMustFunc(t *testing.T) {
    assert.Panics(t, func() {
        f := ks.MustFunc(os.Open)("doesnotexist")
        f.Close()
    })

    successfulFunction := func(input int) (string, error) {
        return fmt.Sprintf("%d\n", input), nil
    }

    f := ks.MustFunc(successfulFunction)
    assert.Equal(t, "3", f(3))
}
