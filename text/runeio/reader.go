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
    rdr *bufio.Reader
    buf []byte
    bufmax int
    offset Offset
    last rune

    // pushedRunes counts caller "manually" pushed back runes, to avoid
    // incrementing the offset
    pushedRunes int
}

// Offset describes an offset into a stream at the time of a call to
// [Reader.Next].
type Offset struct {
    Byte int64 // current steam offset in bytes (0-indexed)
    Rune int64 // current line offset in runes (0-indexed)
    Line int64 // current line (0-indexed)
}

const RuneEOF = rune(0xFFFFFF)

func NewReader(rd io.Reader) *Reader {
    return &Reader{
        rdr: bufio.NewReaderSize(rd, utf8.MaxRune * 64),
    }
}

func (r *Reader) Offset() Offset {
    return r.offset
}

// Last returns the rune most recently returned by [Next]. If [Next] has not
// yet been called, it panics.
func (r *Reader) Last() rune {
    if r.offset.Byte == 0 { panic(fmt.Errorf("runeio error: call to Last before Next")) }
    return r.last
}

// Buffer sets the initial buffer to use when pushing runes back onto the stack
// with [Next]. The capacity is in bytes. If buf is not nil, Cap(buf) must be
// >= the capacity argument. If the Reader already has an existing buffer that
// is not empty, the existing buffer is copied to the new buffer (if there is
// not enough space, this will panic - call [Reader.Clear] first if necessary).
// If buf is nil, then a new buffer is allocated with the given capacity. If
// the capacity is zero, then the reader has no pushback buffer and the
// provided buf is ignored.
//
// For example, to be able to push back at least n Unicode codepoints (runes),
// pass a buffer with capacity utf8.UTFMax * n, or 4n.
func (r *Reader) Buffer(buf []byte, capacity int) {
    fitErr := fmt.Errorf("runeio error: existing buffer does not fit in new buffer")
    if capacity == 0 {
        if len(r.buf) > 0 {
            // can't copy existing buffer
            panic(fitErr)
        }
        r.buf = nil
        r.bufmax = 0
        return
    }
    if buf == nil {
        buf = make([]byte, 0, capacity)
    }
    if capacity < cap(buf) {
        panic(fmt.Errorf("runeio error: cap(buffer) < capacity"))
    }
    buf = buf[0:0]

    // copy existing buffer
    if capacity < len(r.buf) {
        panic(fmt.Errorf("runeio error: cap(buffer) < length of existing buffer"))
    }
    for i := 0; i < len(r.buf); i++ {
        buf[i] = r.buf[i]
    }

    r.buf = buf
    r.bufmax = capacity
}

// Clear clears the reader's pushback buffer (if one exists).
func (r *Reader) Clear() {
    if r.buf != nil {
        r.buf = r.buf[0:0]
    }
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

func (r *Reader) push(x rune) {
    size := utf8.RuneLen(x)
    if size < 0 {
        if x == RuneEOF { return }
        x = utf8.RuneError
        size = 3
    }
    if len(r.buf) + size > r.bufmax {
        panic(fmt.Errorf("runeio pushback buffer overflow"))
    }

    r.buf = utf8.AppendRune(r.buf, x)
}

// Push adds a rune to an internal stack that runes will be read from before
// the input stream is read from.
//
// If a maximum buffer capacity has been set, a panic will be raised if pushing
// a rune would exceed that capacity.
func (r *Reader) Push(x rune) {
    r.push(x)
    r.pushedRunes++
}

func (r *Reader) next() (rune, int, error) {
    if len(r.buf) > 0 {
        x, size := utf8.DecodeLastRune(r.buf)
        r.buf = r.buf[0 : len(r.buf) - size]
        return x, size, nil
    }

    return r.rdr.ReadRune()
}

// Next reads a single rune from the input stream or, if runes have been pushed
// back, reads and removes the most recently pushed rune from the stack. Error
// may be io.EOF or a read error.
func (r *Reader) Next() (rune, error) {
    x, size, err := r.next()

    if r.pushedRunes == 0 {
        r.offset.Byte += int64(size)
        if x == '\n' {
            r.offset.Rune = 0
            r.offset.Line++
        } else {
            r.offset.Rune++
        }
    } else {
        r.pushedRunes--
    }

    r.last = x
    return x, err
}

// Peek returns what the next call to [Reader.Next] would return, without
// advancing the input stream. This requires room on the pushback buffer.
func (r *Reader) Peek() (rune, error) {
    x, _, err := r.next()
    if err == nil {
        r.push(x)
    }
    return x, err
}

// PeekN stores in dest what the next n calls to [Reader.Next] would return,
// without advancing the input stream. This requires room on the pushback
// buffer. It returns the number of elements read, ending early in the event of
// EOF. The first n elements of dest are set to the special value RuneEOF
// unless updated with a successfully peeked value. The pushback buffer must be
// able to handle at least n elements, in addition to its existing contents.
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
        r.push(dest[numRead-i-1])
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
