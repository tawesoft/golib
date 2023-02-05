package tokenizer

import (
    "github.com/tawesoft/golib/v2/css/tokenizer/token"
    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/text/runeio"
)

// position0 returns a [token.Position] from an offset, with a length of
// zero.
func position0(start runeio.Offset) token.Position {
    return token.Position{
        Byte: start.Byte,
        Rune: start.Rune,
        Line: start.Line,
        End:  start.Byte,
    }
}

// position1 returns a [token.Position] from an offset with a length of
// one byte.
func position1(start runeio.Offset) token.Position {
    return token.Position{
        Byte: start.Byte,
        Rune: start.Rune,
        Line: start.Line,
        End:  start.Byte + 1,
    }
}

// position returns a [token.Position] from a (start, end) offset pair.
func position(start runeio.Offset, end runeio.Offset) token.Position {
    if end.Byte < start.Byte {
        must.Never("token end before token start")
    }
    return token.Position{
        Byte: start.Byte,
        Rune: start.Rune,
        Line: start.Line,
        End:  end.Byte - start.Byte,
    }
}

func runeIsWhitespace(x rune) bool {
    // Note that U+000D CARRIAGE RETURN and U+000C FORM FEED are not included in
    // this definition, as they are converted to U+000A LINE FEED during
    // preprocessing.
    return x == '\n' || x == '\t' || x == ' '
}

func runeIsSurrogate(x rune) bool {
    return (x >= 0xD800) && (x <= 0xDFFF)
}

func runeIsHexDigit(x rune) bool {
    return ((x >= '0') && (x <= '9')) ||
           ((x >= 'a') && (x <= 'f')) ||
           ((x >= 'A') && (x <= 'F'))
}

func runeIsDigit(x rune) bool {
    return (x >= '0') && (x <= '9')
}

func runeIsLetter(x rune) bool {
    return ((x >= 'a') && (x <= 'z')) ||
           ((x >= 'A') && (x <= 'Z'))
}

func runeIsNonAscii(x rune) bool {
    return (x >= 0x80) && (x != runeio.RuneEOF)
}

func runeIsIdentStartCodepoint(x rune) bool {
    return runeIsLetter(x) || runeIsNonAscii(x) || (x == '_')
}

func runeIsIdentCodepoint(x rune) bool {
    return runeIsIdentStartCodepoint(x) || runeIsDigit(x) || (x == '-')
}

func runeIsNonPrintable(x rune) bool {
    // A code point between U+0000 NULL and U+0008 BACKSPACE inclusive, or U+000B
    // LINE TABULATION, or a code point between U+000E SHIFT OUT and U+001F
    // INFORMATION SEPARATOR ONE inclusive, or U+007F DELETE.
    if (x >= 0x00) && (x <= 0x08) { return true }
    if (x == 0x0B) { return true }
    if (x >= 0x0E) && (x <= 0x1F) { return true }
    if (x == 0x7F) { return true }
    return false
}

func isValidEscape(a rune, b rune) bool {
    return (a == '\\') && (b != '\n')
}

func isStartOfIdentSequence(a rune, b rune, c rune) bool {
    // https://www.w3.org/TR/css-syntax-3/#would-start-an-identifier
    if runeIsIdentStartCodepoint(a) { return true }
    switch a {
        case '-': // U+002D HYPHEN-MINUS
            // If the second code point is an ident-start code point or a
            // U+002D HYPHEN-MINUS, or the second and third code points are a
            // valid escape, return true.
            if runeIsIdentStartCodepoint(b) || (b == '-') { return true }
            if isValidEscape(b, c) { return true }
            return false
        case '\\': // U+005C REVERSE SOLIDUS (\)
            return isValidEscape(a, b)
        default:
            return false
    }
}

func isStartOfNumber(a rune, b rune, c rune) bool {
    // https://www.w3.org/TR/css-syntax-3/#starts-with-a-number
    switch {
        case a == '+':
            fallthrough
        case a == '-':
            // If the second code point is a digit, return true.
            // Otherwise, if the second code point is a U+002E FULL STOP (.)
            // and the third code point is a digit, return true.
            if runeIsDigit(b) { return true }
            if (b == '.') && runeIsDigit(c) { return true }
            return false
        case a == '.':
            return runeIsDigit(b)
        case runeIsDigit(a):
            return true
        default:
            return false
    }
}
