package runeio_test

import (
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/v2/text/runeio"
)

func TestPeekN(t *testing.T) {
    var buf [6]rune
    r := runeio.NewReader(strings.NewReader("hello"))
    r.Buffer(nil, 6)
    n, err := r.PeekN(buf[:], 6)

    assert.Nil(t, err)
    assert.Equal(t, 5, n)
    assert.Equal(t, 'h', buf[0])
    assert.Equal(t, 'e', buf[1])
    assert.Equal(t, 'l', buf[2])
    assert.Equal(t, 'l', buf[3])
    assert.Equal(t, 'o', buf[4])
    assert.Equal(t, runeio.RuneEOF, buf[5])
}

func TestPush(t *testing.T) {
    r := runeio.NewReader(strings.NewReader("hello"))
    r.Buffer(nil, 2)

    assert.Equal(t, 'h', runeio.Must(r.Next()))
    assert.Equal(t, 'e', runeio.Must(r.Next()))
    r.Push('x')
    assert.Equal(t, 'x', runeio.Must(r.Next()))
    assert.Equal(t, 'l', runeio.Must(r.Next()))
}

func TestPeek(t *testing.T) {
    r := runeio.NewReader(strings.NewReader("hello"))
    r.Buffer(nil, 2)

    assert.Equal(t, 'h', runeio.Must(r.Peek()))
    assert.Equal(t, 'h', runeio.Must(r.Peek()))
    assert.Equal(t, 'h', runeio.Must(r.Next()))
    r.Push('x')
    assert.Equal(t, 'x', runeio.Must(r.Peek()))
    assert.Equal(t, 'x', runeio.Must(r.Peek()))
    assert.Equal(t, 'x', runeio.Must(r.Next()))
    assert.Equal(t, 'e', runeio.Must(r.Peek()))
    assert.Equal(t, 'e', runeio.Must(r.Peek()))
    assert.Equal(t, 'e', runeio.Must(r.Next()))
}
