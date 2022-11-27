// Package runeio implements a mechanism to read a stream of Unicode code
// points (runes) from an io.Reader, with an internal buffer to push code
// points back to the front of the stream to allow limited peeking and rewind.
package runeio

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "unicode/utf8"
)

// Reader has a useful zero value, but doesn't have a pushback buffer by
// default.
type Reader struct {
    br *bufio.Reader
    buf []byte
    bufmax int
    offset Offset
    last rune
}

type Offset struct {
    Byte int64 // current steam offset in bytes (0-indexed)
    Rune int64 // current line offset in runes (0-indexed)
    Line int64 // current line (0-indexed)
}

const RuneEOF = rune(0xFFFFFF)

func NewReader(rd io.Reader) *Reader {
    return NewReaderBuffered(bufio.NewReader(rd))
}

func NewReaderBuffered(rd *bufio.Reader) *Reader {
    return &Reader{
        br: rd,
    }
}

func (r *Reader) Offset() Offset {
    return r.offset
}

// Last returns the rune most recently returned by [Next]. If [Next] has not
// yet been called, it panics.
func (r *Reader) Last() rune {
    if r.offset.Byte == 0 { panic(fmt.Errorf("runeio invalid reader Last")) }
    return r.last
}

// Buffer sets the initial buffer to use when pushing runes back onto the stack
// with [Next] and the maximum capacity (in bytes) for that stack. If the
// maximum is less than or equal to the buffer capacity, then it will never
// grow. (TODO right now it always grows)
//
// Passing a nil buffer will allocate a new buffer.
//
// Buffer panics if it is called after reading has started.
func (r *Reader) Buffer(buf []byte, max int) {
    if buf == nil { buf = make([]byte, 0, max) }
    r.buf = buf
    if max < cap(buf) { max = 0 }
    r.bufmax = max
}

// Must accepts a (rune, error) pair and always returns a rune value or raises
// a panic. If the error is nil, returns the input rune as normal. If the error
// is io.EOF, returns a special RuneEOF value. Otherwise, if the error is not
// nil, panics, wrapping the error.
func Must(x rune, err error) rune {
    if err == nil {
        return x
    } else if errors.Is(err, io.EOF) {
        return RuneEOF
    } else {
        panic(fmt.Errorf("runeio read error: %w", err))
    }
}

// MustFunc wraps a function that returns a (rune, error) pair with [Must] and
// returns a new function that always returns a rune, RuneEOF, or panics.
func MustFunc(f func() (rune, error)) func() rune {
    return func() rune {
        return Must(f())
    }
}

// Push adds a rune to an internal stack that runes will be read from before
// the input stream is read from.
//
// If a maximum buffer capacity has been set, a panic will be raised if pushing
// a rune would exceed that capacity.
func (r *Reader) Push(x rune) {
    size := utf8.RuneLen(x)
    if size < 0 { return }
    if len(r.buf) + size > r.bufmax {
        panic(fmt.Errorf("runeio pushback buffer overflow"))
    }

    r.buf = utf8.AppendRune(r.buf, x)
}

func (r *Reader) next() (rune, int, error) {
    if len(r.buf) > 0 {
        x, size := utf8.DecodeLastRune(r.buf)
        r.buf = r.buf[0 : len(r.buf) - size]
        return x, size, nil
    }

    return r.br.ReadRune()
}

// Next reads a single rune from the input stream or, if runes have been pushed
// back, reads and removes the most recently pushed rune from the stack. Error
// may be io.EOF or a read error.
func (r *Reader) Next() (rune, error) {
    x, size, err := r.next()

    r.offset.Byte += int64(size)
    if x == '\n' {
        r.offset.Rune = 0
        r.offset.Line++
    } else {
        r.offset.Rune++
    }

    r.last = x
    return x, err
}

// Peek returns what the next call to [Next] would return, without advancing
// the input stream.
func (r *Reader) Peek() (rune, error) {
    x, _, err := r.next()
    if err == nil {
        r.Push(x)
    }
    return x, err
}

// PeekN stores in dest what the next n calls to [Next] would return, without
// advancing the input stream. It returns the number of elements read, ending
// early in the event of EOF. The first n elements of dest are set to the
// special value RuneEOF unless updated with a successfully peeked value. The
// pushback buffer must be able to handle at least n elements, in addition to
// its existing contents.
func (r *Reader) PeekN(dest []rune, n int) (int, error) {

    for i := 0; i < n; i++ {
        dest[i] = RuneEOF
    }

    var numRead int
    for i := 0; i < n; i++ {
        x, _, err := r.next()
        if errors.Is(err, io.EOF) {
            break
        } else if err != nil {
            return 0, err
        }
        dest[numRead] = x
        numRead++
    }

    for i := 0; i < numRead; i++ {
        r.Push(dest[numRead-i-1])
    }

    return numRead, nil
}

func (r *Reader) Skip(n int) error {
    for i := 0; i < n; i++ {
        _, err := r.Next()
        if err == io.EOF {
            break
        } else if err != nil {
            return err
        }
    }
    return nil
}
