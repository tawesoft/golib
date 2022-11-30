package runeio_test

import (
    "strings"
    "testing"
    "unicode/utf8"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/v2/text/runeio"
)

func TestPeekN(t *testing.T) {
    var buf [6]rune
    r := runeio.NewReader(strings.NewReader("hello"))
    r.Buffer(nil, utf8.UTFMax * 6)
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
    r.Buffer(nil, utf8.UTFMax * 1)

    assert.Equal(t, 'h', runeio.Must(r.Next()))
    assert.Equal(t, 'e', runeio.Must(r.Next()))
    r.Push('x')
    assert.Equal(t, 'x', runeio.Must(r.Next()))
    assert.Equal(t, 'l', runeio.Must(r.Next()))
}

func TestPeek(t *testing.T) {
    r := runeio.NewReader(strings.NewReader("hello"))
    r.Buffer(nil, utf8.UTFMax * 1)

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

func TestOffset(t *testing.T) {
    r := runeio.NewReader(strings.NewReader("héllo\nworld"))
    r.Buffer(nil, utf8.UTFMax * 4)
    var dest [4]rune

    r.Next();  assert.Equal(t, runeio.Offset{1, 1, 0}, r.Offset())
    r.PeekN(dest[:], 4)
    assert.Equal(t, string(dest[:]), "éllo")
    r.Next();  assert.Equal(t, runeio.Offset{3, 2, 0}, r.Offset())
    r.Push('x');
    r.Peek();
    r.Skip(5); assert.Equal(t, runeio.Offset{7, 0, 1}, r.Offset())
    r.Peek(); r.Peek(); r.Peek()
    r.Next();  assert.Equal(t, runeio.Offset{8, 1, 1}, r.Offset())
    r.Push('x');
    r.Next();  assert.Equal(t, runeio.Offset{8, 1, 1}, r.Offset())
}
