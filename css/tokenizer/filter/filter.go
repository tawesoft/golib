// Package filter implements a [transform.Transformer] that performs the Unicode
// code point filtering preprocessing step defined in
// [CSS Syntax Module Level 3, section 3.3]:
//
// "To filter code points from a stream of (unfiltered) code points input:
//
// "Replace any U+000D CARRIAGE RETURN (CR) code points, U+000C FORM FEED (FF)
// code points, or pairs of U+000D CARRIAGE RETURN (CR) followed by U+000A LINE
// FEED (LF) in input by a single U+000A LINE FEED (LF) code point.
//
// "Replace any U+0000 NULL or surrogate code points in input with U+FFFD REPLACEMENT
// CHARACTER (�)."
//
// [CSS Syntax Module Level 3, section 3.3]: https://www.w3.org/TR/css-syntax-3/#input-preprocessing
package filter

import (
    "fmt"
    "unicode/utf8"

    "golang.org/x/text/transform"
)

var DecodeError = fmt.Errorf("error decoding input (not valid UTF8)")

type filter struct {
    last rune
}

func Transformer() transform.Transformer {
    return &filter{}
}

func (t *filter) Reset() {
    t.last = 0
}

func emit(r rune, size int, dst []byte, nDst *int, err *error) {
    if *nDst + size > len(dst) {
        *err = transform.ErrShortDst
    } else {
        utf8.EncodeRune(dst[*nDst:], r)
        (*nDst)+= size
    }
    return
}

func (t *filter) Transform(dst, src []byte, atEOF bool) (nDst int, nSrc int, err error) {
    for nSrc < len(src) {
        r, size := utf8.DecodeRune(src[nSrc:])

        if !atEOF && !utf8.FullRune(src[nSrc:]) {
            err = transform.ErrShortSrc
            break
        } else if atEOF && (size == 0) {
            break
        } else if r == utf8.RuneError {
            err = DecodeError
            break
        }

        if (r != '\n') && (t.last == '\r') {
            emit('\n', 1, dst, &nDst, &err)
            if err != nil { break }
            t.last = 0
        }

        if r != '\r' {
            t.last = 0
        }

        if r == '\r' {
            t.last = '\r'
        } else if r == '\f' {
            emit('\n', 1, dst, &nDst, &err)
            if err != nil { break }
        } else if (r == 0) || !utf8.ValidRune(r) {
            emit(utf8.RuneError, utf8.RuneLen(utf8.RuneError), dst, &nDst, &err)
        } else {
            emit(r, size, dst, &nDst, &err)
            if err != nil { break }
        }

        nSrc += size
    }

    // trailing carriage return
    if atEOF && (t.last == '\r') {
        emit('\n', 1, dst, &nDst, &err)
    }

    return
}
