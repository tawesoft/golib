package ks_test

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

func TestCatchFunc(t *testing.T) {
    {
        failingFunctionWithString := func() string {
            panic("oops")
            return "something"
        }

        x, err := ks.CatchFunc[string](failingFunctionWithString)()
        assert.Equal(t, "", x)
        assert.Error(t, err)
        assert.Equal(t, nil, errors.Unwrap(err))
    }

    {
        failingFunctionWithError := func() string {
            panic(fmt.Errorf("oops"))
            return "something"
        }

        x, err := ks.CatchFunc[string](failingFunctionWithError)()
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
        return fmt.Sprintf("%d", input), nil
    }

    f := ks.MustFunc(successfulFunction)
    assert.Equal(t, "3", f(3))
}

func TestWordWrap(t *testing.T) {
    var tests = [][3]interface{}{
        {"",                        "",                           1},
        {"a",                       "",                          -1},
        {"a",                       "",                           0},
        {"a",                       "a",                          1},
        {"a",                       "a",                          2},
        {"  a  ",                   "a",                          1},
        {"helloworld",              "hello",                      5},
        {"hello\nworld",            "he\nwo",                     2},
        {"hello world",             "he\nwo",                     2},
        {"hello world",             "hello\nworld",               5},
        {"hello world",             "hello\nworld",              10},
        {"hello world",             "hello world",               11},
        {"hello world",             "hello world",               12},
        {"hello\nworld",            "hello world",               12},
        {"hello      world",        "hello world",               12},
        {"    hello    world   ",   "hello world",               12},
        {"a b c d e f g h i",       "a\nb\nc\nd\ne\nf\ng\nh\ni",  1},
        {"a b c d e f g h i",       "a b c\nd e f\ng h i",        5},
        {"a b c d e f g h i",       "a b c\nd e f\ng h i",        6},
        {"a b c d e f g h i",       "a b c d e\nf g h i",         9},
        {"a b c d e f g h i",       "a b c d e\nf g h i",        10},
        {`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce a
        tortor sagittis, elementum velit id, scelerisque erat. Sed mollis odio
        molestie dui venenatis condimentum. Donec massa ligula, auctor rutrum
        interdum a, faucibus sed sapien. Vivamus neque massa, porttitor vel
        nulla eu, gravida egestas massa. Aliquam interdum pellentesque elit.
        Quisque vestibulum, libero condimentum venenatis commodo, erat lectus
        convallis libero, at pellentesque nibh enim vel risus. Duis elit mi,
        lacinia ut ex vitae, ullamcorper tempus ex. Lorem ipsum dolor sit amet,
        consectetur adipiscing elit. Fusce eu elit molestie, tempor nulla
        vehicula, tempor nulla. Maecenas pellentesque, lectus non accumsan
        pharetra, neque justo dignissim dolor, sit amet luctus mi leo ut dui.`,
        `Lorem ipsum dolor sit amet,
consectetur adipiscing elit.
Fusce a tortor sagittis,
elementum velit id,
scelerisque erat. Sed mollis
odio molestie dui venenatis
condimentum. Donec massa
ligula, auctor rutrum interdum
a, faucibus sed sapien.
Vivamus neque massa, porttitor
vel nulla eu, gravida egestas
massa. Aliquam interdum
pellentesque elit. Quisque
vestibulum, libero condimentum
venenatis commodo, erat lectus
convallis libero, at
pellentesque nibh enim vel
risus. Duis elit mi, lacinia
ut ex vitae, ullamcorper
tempus ex. Lorem ipsum dolor
sit amet, consectetur
adipiscing elit. Fusce eu elit
molestie, tempor nulla
vehicula, tempor nulla.
Maecenas pellentesque, lectus
non accumsan pharetra, neque
justo dignissim dolor, sit
amet luctus mi leo ut dui.`,
        30},
    }

    for index, test := range tests {
        var original, expected, length = test[0], test[1], test[2]
        var result = ks.WrapBlock(original.(string), length.(int))
        if result != expected.(string) {
            t.Errorf("Test %d failed: wrap(\"%s\", %d), got \"%s\" but wanted \"%s\"\n",
                index, original.(string), length.(int), result, expected.(string))
        }
    }
}
