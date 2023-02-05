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
    var c rune

    r.Next()
    assert.Equal(t, runeio.Offset{1, 1, 0}, r.Offset())
    r.PeekN(dest[:], 4)
    assert.Equal(t, string(dest[:]), "éllo")
    r.Peek()
    r.Push('界')
    r.Skip(1)
    c = runeio.Must(r.Next())
    assert.Equal(t, 'é', c)
    assert.Equal(t, c, r.Last())
    assert.Equal(t, runeio.Offset{3, 2, 0}, r.Offset())
    r.Push('界')
    r.Peek();
    r.Skip(1)
    r.Skip(3) // llo
    assert.Equal(t, runeio.Offset{6, 5, 0}, r.Offset())
    r.Push('界')
    r.Peek(); r.Peek(); r.Peek(); r.Skip(1)
    c = runeio.Must(r.Next())
    assert.Equal(t, '\n', c)
    assert.Equal(t, c, r.Last())
    assert.Equal(t, runeio.Offset{7, 0, 1}, r.Offset())
    r.Next()
    assert.Equal(t, runeio.Offset{8, 1, 1}, r.Offset())
    r.Push('界')
    r.Next()
    assert.Equal(t, runeio.Offset{8, 1, 1}, r.Offset())
    c = runeio.Must(r.Next())
    assert.Equal(t, 'o', c)
    assert.Equal(t, runeio.Offset{9, 2, 1}, r.Offset())
    r.Skip(3) // rld
    assert.Equal(t, runeio.Offset{12, 5, 1}, r.Offset())
}

func TestOffsetEof(t *testing.T) {
}
