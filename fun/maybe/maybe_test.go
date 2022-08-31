package maybe_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/v2/fun/maybe"
)

func TestMaybe(t *testing.T) {
    divide := func(x int, y int) maybe.Maybe[int] {
        if y == 0 { return maybe.Nothing[int]() }
        return maybe.Some(x / y)
    }

    double := func(x int) int {
        return x * 2
    }

    assert.Equal(t, 4, maybe.Map(divide(16, 8), double).Must())
    assert.Equal(t, 4, maybe.FlatMap(divide(16, 8), maybe.Lift(double)).Must())
    maybe.Map(divide(8, 0), double).MustNot()

    assert.Equal(t, 4, maybe.FlatMap(divide(16, 2), func(x int) maybe.Maybe[int] {
        return divide(x, 2)
    }).Must())

    maybe.FlatMap(divide(16, 8), func(x int) maybe.Maybe[int] {
        return divide(x, 0)
    }).MustNot()

    incrementer := maybe.Some(func(x int) int {
        return x + 1
    })
    assert.Equal(t, 6, maybe.Apply(maybe.Some(5), incrementer).Must())
}
