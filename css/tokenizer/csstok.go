// Package tokenizer performs the tokenization step defined in
// [CSS Syntax Module Level 3] (part 4).
//
// The main elements of this package are the [New] function, which returns a
// new [Tokenizer], and that Tokenizer's [Tokenizer.Next] method.
//
// This package also exposes several low-level "Consume" functions, which
// implement specific algorithms in the CSS specification.
//
// Note that all "Consume" functions may panic on I/O error. The
// [Tokenizer.Next] method catches these panics.
//
// Note that all "Consume" functions operate on a steam of filtered code points
// (see https://www.w3.org/TR/css-syntax-3/#input-preprocessing). This is
// handled by a [New] Tokenizer.
//
// [CSS Syntax Module Level 3]: https://www.w3.org/TR/css-syntax-3/
//
// Portions Copyright © 2022 W3C® (MIT, ERCIM, Keio, Beihang)
package tokenizer

import (
    "bufio"
    "fmt"
    "io"
    "math"
    "strconv"
    "strings"
    "unicode"
    "unicode/utf8"

    "github.com/tawesoft/golib/v2/css/tokenizer/filter"
    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/text/runeio"
    "golang.org/x/text/transform"
)

const MaxLookahead = 7 // normally up to 3, but 7 used in urange

var (
    ErrUnexpectedEOF = fmt.Errorf("unexpected end of file")
    ErrUnexpectedLinebreak = fmt.Errorf("unexpected line break")
    ErrUnexpectedInput = fmt.Errorf("unexpected input")
    ErrBadUrl = fmt.Errorf("invalid URL syntax")
)

type Tokenizer struct {
    rdr *runeio.Reader
    errors []error
    eof bool
}

func reader(r io.Reader) *runeio.Reader {
    br := bufio.NewReader(r)
    rdr := runeio.NewReader(transform.NewReader(br, filter.Transformer()))
    rdr.Buffer(nil, utf8.UTFMax * MaxLookahead)
    return rdr
}

func New(r io.Reader) Tokenizer {
    return Tokenizer{
        rdr: reader(r),
    }
}

// Errors reports parse errors.
func (z *Tokenizer) Errors() []error {
    return z.errors
}

type readError struct {
    err error
    offset runeio.Offset
}

func (e readError) Error() string {
    return fmt.Sprintf("parse error at %+v: %s", e.offset, e.err)
}

func (e readError) Unwrap() error {
    return e.err
}

func (z *Tokenizer) error(err error) {
    rerr := readError{
        err: err,
        offset: z.rdr.Offset(),
    }
    z.errors = append(z.errors, rerr)
}

type TokenType string
const (
    TokenTypeWhitespace         = TokenType("whitespace-token")
    TokenTypeEOF                = TokenType("EOF-token")
    TokenTypeString             = TokenType("string-token")
    TokenTypeBadString          = TokenType("bad-string-token")
    TokenTypeDelim              = TokenType("delim-token")
    TokenTypeHash               = TokenType("hash-token")
    TokenTypeLeftParen          = TokenType("(-token")
    TokenTypeRightParen         = TokenType(")-token")
    TokenTypeNumber             = TokenType("number-token")
    TokenTypeDimension          = TokenType("dimension-token")
    TokenTypePercentage         = TokenType("percentage-token")
    TokenTypeCDC                = TokenType("CDC-token")
    TokenTypeIdent              = TokenType("ident-token")
    TokenTypeFunction           = TokenType("function-token")
    TokenTypeUrl                = TokenType("url-token")
    TokenTypeBadUrl             = TokenType("bad-url-token")
    TokenTypeColon              = TokenType("colon-token")
    TokenTypeSemicolon          = TokenType("semicolon-token")
    TokenTypeCDO                = TokenType("CDO-token")
    TokenTypeAtKeyword          = TokenType("at-keyword-token")
    TokenTypeLeftSquareBracket  = TokenType("[-token")
    TokenTypeRightSquareBracket = TokenType("]-token")
    TokenTypeLeftCurlyBracket  = TokenType("{-token")
    TokenTypeRightCurlyBracket = TokenType("}-token")
)

type HashType string
const (
    HashTokenTypeID           = HashType("id")
    HashTokenTypeUnrestricted = HashType("unrestricted")
)

type NumberType string
const (
    NumberTypeInteger = NumberType("integer")
    NumberTypeNumber  = NumberType("number")
)

type Token struct {
    Type TokenType

    // The lexeme as it appears in the input stream. This preserves details
    // such as whether .009 was written as .009 or 9e-3, and whether a
    // character was written literally or as a CSS escape. Only used by
    // <number-token>, <dimension-token>, <percentage-token>, ...
    repr string

    // Value is used by <ident-token>, <function-token>, <at-keyword-token>,
    // <hash-token>, <string-token>, and <url-token>.
    stringValue string

    unit string // used by <dimension-token>.

    // Hash tokens have a type flag set to either "id" or "unrestricted".
    // The type flag defaults to "unrestricted" if not otherwise set.
    hashType HashType

    // <number-token> and <dimension-token> additionally have a type flag set
    // to either "integer" or "number".
    // NOTE: spec is a bit ambiguous, but assume this also applies to
    // percentage-tokens.
    numberType NumberType

    delim rune // used by <delim-token>

    // numberValue is used by <number-token>, <dimension-token>,
    // <percentage-token>.
    numberValue float64
}

func (t Token) String() string {
    switch t.Type {
        case TokenTypeString:    fallthrough
        case TokenTypeAtKeyword: fallthrough
        case TokenTypeUrl:       fallthrough
        case TokenTypeFunction:  fallthrough
        case TokenTypeIdent:
            return fmt.Sprintf("<%s>{value: %q}", t.Type, t.stringValue)
        case TokenTypeDelim:
            return fmt.Sprintf("<%s>{delim: %q}", t.Type, t.delim)
        case TokenTypeHash:
            return fmt.Sprintf("<%s>{type: %q, value: %q}", t.Type, t.hashType, t.stringValue)
        case TokenTypeNumber:
            fallthrough
        case TokenTypePercentage:
            return fmt.Sprintf("<%s>{type: %q, value: %f, repr: %q}", t.Type, t.numberType, t.numberValue, t.repr)
        case TokenTypeDimension:
            return fmt.Sprintf("<%s>{type: %q, value: %f, unit: %q, repr: %q}", t.Type, t.numberType, t.numberValue, t.unit, t.repr)
        default:
            return fmt.Sprintf("<%s>", t.Type)
    }
}

// Next returns the next token from the input stream while the second return
// value is true.
func (z *Tokenizer) Next() (result Token, ok bool) {
    if z.eof { return Token{Type: TokenTypeEOF}, false }

    defer func() {
        if r := recover(); r != nil {
            z.error(r.(error))
            ok = false
        }
    }()

    err := ConsumeComments(z.rdr)
    if err != nil { z.error(err) } // recovers

    c := runeio.Must(z.rdr.Next())
    switch {
        case runeIsWhitespace(c):
            return ConsumeWhitespace(z.rdr), true
        case c == '"': // U+0022 QUOTATION MARK (")
            fallthrough
        case c == '\'': // U+0027 APOSTROPHE (')
            t, err := ConsumeString(z.rdr, c)
            if err != nil { z.error(err) }
            return t, true
        case c == '#': // U+0023 NUMBER SIGN (#)
            // If the next input code point is an ident code point or the next
            // two input code points are a valid escape, then:
            var xs [3]rune
            must.Result(z.rdr.PeekN(xs[:], 3))
            if runeIsIdentCodepoint(xs[0]) || isValidEscape(xs[0], xs[1]) {
                // Create a <hash-token>.
                hashTokenType := HashTokenTypeUnrestricted
                // If the next 3 input code points would start an ident
                // sequence, set the <hash-token>’s type flag to "id".
                if isStartOfIdentSequence(xs[0], xs[1], xs[2]) {
                    hashTokenType = HashTokenTypeID
                }
                // Consume an ident sequence, and set the <hash-token>’s value
                // to the returned string. Return the <hash-token>.
                ident := ConsumeIdentSequence(z.rdr)
                return NewHashToken(hashTokenType, ident), true
            } else {
                // Otherwise, return a <delim-token> with its value set to the
                // current input code point.
                return NewDelimToken(c), true
            }
        case c == '(': // U+0028 LEFT PARENTHESIS (()
            return Token{Type:TokenTypeLeftParen}, true
        case c == ')': // U+0029 RIGHT PARENTHESIS ())
            return Token{Type:TokenTypeRightParen}, true
        case c == '+': // U+002B PLUS SIGN (+)
            // If the input stream starts with a number...
            var xs[3]rune
            must.Result(z.rdr.PeekN(xs[:], 3))
            if isStartOfNumber(xs[0], xs[1], xs[2]) {
                // reconsume the current input code point, consume a numeric
                // token, and return it.
                z.rdr.Push(c)
                return ConsumeNumericToken(z.rdr), true
            } else {
                // Otherwise, return a <delim-token> with its value set to the
                // current input code point.
                return NewDelimToken(c), true
            }
        case c == ',': // U+002C COMMA (,)
            return NewDelimToken(c), true
        case c == '-': // U+002D HYPHEN-MINUS (-)
            // If the input stream starts with a number...
            var xs[3]rune
            must.Result(z.rdr.PeekN(xs[:], 3))
            if isStartOfNumber(xs[0], xs[1], xs[2]) {
                // reconsume the current input code point, consume a numeric
                // token, and return it.
                z.rdr.Push(c)
                return ConsumeNumericToken(z.rdr), true
                // Otherwise, if the next 2 input code points are
                // U+002D HYPHEN-MINUS U+003E GREATER-THAN SIGN (->)...
            } else if (xs[0] == '-') && (xs[1] == '>') {
                // consume them and return a <CDC-token>.
                z.rdr.Skip(2)
                return Token{Type: TokenTypeCDC}, true
                // Otherwise, if the input stream starts with an ident sequence,
            } else if isStartOfIdentSequence(xs[0], xs[1], xs[2]) {
                // reconsume the current input code point, consume an
                // ident-like token, and return it.
                z.rdr.Push(c)
                t, err := ConsumeIdentLikeToken(z.rdr)
                if err != nil { z.error(err) }
                return t, true
            } else {
                // Otherwise, return a <delim-token> with its value set to the
                // current input code point.
                return NewDelimToken(c), true
            }
        case c == '.': // U+002E FULL STOP (.)
            // If the input stream starts with a number...
            var xs[3]rune
            must.Result(z.rdr.PeekN(xs[:], 3))
            if isStartOfNumber(xs[0], xs[1], xs[2]) {
                // reconsume the current input code point,
                z.rdr.Push(c)
                // consume a numeric token, and return it.
                return ConsumeNumericToken(z.rdr), true
                // Otherwise...
            } else {
                // return a <delim-token> with its value set to the current
                // input code point.
                return NewDelimToken(c), true
            }
        case c == ':': // U+003A COLON (:)
            return Token{Type: TokenTypeColon}, true
        case c == ';': // U+003B SEMICOLON (;)
            return Token{Type: TokenTypeSemicolon}, true
        case c == '<': // U+003C LESS-THAN SIGN (<)
            // If the next 3 input code points are
            // U+0021 EXCLAMATION MARK
            // U+002D HYPHEN-MINUS
            // U+002D HYPHEN-MINUS (!--)...
            var xs[3]rune
            must.Result(z.rdr.PeekN(xs[:], 3))
            if (xs[0] == '!') && (xs[1] == '-') && (xs[2] == '-') {
                // consume them and return a <CDO-token>.
                z.rdr.Skip(3)
                return Token{Type: TokenTypeCDO}, true
                // Otherwise...
            } else {
                // return a <delim-token> with its value set
                // to the current input code point.
                return NewDelimToken(c), true
            }
        case c == '@': // U+0040 COMMERCIAL AT (@)
            // If the next 3 input code points would start an ident sequence...
            var xs[3]rune
            must.Result(z.rdr.PeekN(xs[:], 3))
            if isStartOfIdentSequence(xs[0], xs[1], xs[2]) {
                // consume an ident sequence, create an <at-keyword-token> with
                // its value set to the returned value, and return it.
                return NewAtKeywordToken(ConsumeIdentSequence(z.rdr)), true
                // Otherwise...
            } else {
                // return a <delim-token> with its value set
                // to the current input code point.
                return NewDelimToken(c), true
            }
        case c == '[': // U+005B LEFT SQUARE BRACKET ([)
            return Token{Type:TokenTypeLeftSquareBracket}, true
        case c == '\\': // U+005C REVERSE SOLIDUS (\)
            // If the input stream starts with a valid escape...
            p := runeio.Must(z.rdr.Peek())
            if isValidEscape(c, p) {
                // reconsume the current input code point,
                z.rdr.Push(c)
                // consume an ident-like token, and return it.
                t, err := ConsumeIdentLikeToken(z.rdr)
                if err != nil { z.error(err) }
                return t, true
                // Otherwise...
            } else {
                // this is a parse error.
                z.error(ErrUnexpectedInput)
                // Return a <delim-token> with its value set to
                // the current input code point.
                return NewDelimToken(c), true
            }
        case c == ']': // U+005D RIGHT SQUARE BRACKET (])
            return Token{Type: TokenTypeRightSquareBracket}, true
        case c == '{': // U+007B LEFT CURLY BRACKET ({)
            return Token{Type: TokenTypeLeftCurlyBracket}, true
        case c == '}': // U+007D RIGHT CURLY BRACKET (})
            return Token{Type: TokenTypeRightCurlyBracket}, true
        case runeIsDigit(c):
            // Reconsume the current input code point,
            z.rdr.Push(c)
            // consume a numeric token, and return it.
            return ConsumeNumericToken(z.rdr), true
        case runeIsIdentStartCodepoint(c):
            // Reconsume the current input code point,
            z.rdr.Push(c)
            // consume an ident-like token, and return it.
            t, err := ConsumeIdentLikeToken(z.rdr)
            if err != nil { z.error(err) }
            return t, true
        case c == runeio.RuneEOF:
            z.eof = true
            return Token{Type: TokenTypeEOF}, true
        default: // anything else
            return NewDelimToken(c), true
    }
}

// ConsumeComments consumes zero or more CSS comments.
func ConsumeComments(rdr *runeio.Reader) error {
    for {
        var xs [2]rune

        // If the next two input code points...
        must.Result(rdr.PeekN(xs[:], 2))

        // are U+002F SOLIDUS (/) followed by a U+002A ASTERISK (*)...
        if xs[0] == '/' && xs[1] == '*' {
            // consume them...
            must.Check(rdr.Skip(2))

            x := rune(0)
            for {
                y := runeio.Must(rdr.Next())

                // ... and all following code points up to and including
                // the first U+002A ASTERISK (*) followed by a U+002F SOLIDUS (/),
                if x == '*' && y == '/' {
                    break
                }

                // or up to an EOF code point
                // (this is a parse error)
                if y == runeio.RuneEOF {
                    return ErrUnexpectedEOF
                }

                x = y
            }
        } else {
            return nil
        }
        // and repeat
    }
    return nil
}

// ConsumeWhitespace consumes as much whitespace as possible and returns a
// <whitespace-token>.
func ConsumeWhitespace(rdr *runeio.Reader) Token {
    for runeIsWhitespace(runeio.Must(rdr.Peek())) {
        runeio.Must(rdr.Next())
    }

    return Token{Type: TokenTypeWhitespace}
}

// ConsumeString consumes a string token. It is assumed that the character that
// opens a string (if any) has already been consumed. Returns either a
// <string-token> or a <bad-string-token>. Endpoint specifies the codepoint
// that terminates the string (e.g. a double or single quotation mark).
func ConsumeString(rdr *runeio.Reader, endpoint rune) (t Token, err error) {
    // https://www.w3.org/TR/css-syntax-3/#consume-string-token
    var sb strings.Builder

    for {
        c := runeio.Must(rdr.Next())
        switch c {
            case endpoint:
                return NewStringToken(sb.String()), nil
            case runeio.RuneEOF:
                return NewStringToken(sb.String()), ErrUnexpectedEOF
            case '\n':
                rdr.Push(c)
                return NewBadStringToken(sb.String()), ErrUnexpectedLinebreak
            case '\\': // U+005C REVERSE SOLIDUS (\)
                n := runeio.Must(rdr.Peek())
                if n == runeio.RuneEOF { continue }
                if n == '\n' { rdr.Skip(1); continue }
                sb.WriteRune(ConsumeEscapedCodepoint(rdr))
            default:
                sb.WriteRune(c)
        }
    }
}

// ConsumeEscapedCodepoint consumes an escaped code point. It assumes that
// the U+005C REVERSE SOLIDUS (\) has already been consumed and that the next
// input code point has already been verified to be part of a valid escape.
func ConsumeEscapedCodepoint(rdr *runeio.Reader) rune {
    // Consume the next input code point.
    c := runeio.Must(rdr.Next())
    if !runeIsHexDigit(c) { return c }

    // Consume as many hex digits as possible, but no more than 5. Note that this
    // means 1-6 hex digits have been consumed in total.
    var buf [6]rune
    var size int

    buf[0] = c
    size++
    for i := 1; i < 6; i++ {
        c = runeio.Must(rdr.Peek())
        if !runeIsHexDigit(c) { break }
        rdr.Skip(1)

        buf[i] = c
        size++
    }

    // If the next input code point is whitespace, consume it as well.
    p := runeio.Must(rdr.Peek())
    if runeIsWhitespace(p) { rdr.Skip(1) }

    // Interpret the hex digits as a hexadecimal number. If this number is zero, or
    // is for a surrogate, or is greater than the maximum allowed code point,
    // return U+FFFD REPLACEMENT CHARACTER (�). Otherwise, return the code point
    // with that value.
    n, err := strconv.ParseInt(string(buf[0:size]), 16, 64)
    if (err != nil) || (n == 0) || (n > unicode.MaxRune) || runeIsSurrogate(rune(n)) {
        return unicode.ReplacementChar
    } else {
        return rune(n)
    }
}

// ConsumeIdentSequence consumes an ident sequence from a stream of code
// points. It returns a string containing the largest name that can be formed
// from adjacent code points in the stream, starting from the first.
//
// Note: This algorithm does not do the verification of the first few code
// points that are necessary to ensure the returned code points would
// constitute an <ident-token>. If that is the intended use, ensure that the
// stream starts with an ident sequence before calling this algorithm.
func ConsumeIdentSequence(rdr *runeio.Reader) string {
    // https://www.w3.org/TR/css-syntax-3/#consume-name
    var sb strings.Builder

    for {
        c := runeio.Must(rdr.Next())

        if runeIsIdentCodepoint(c) {
            sb.WriteRune(c)
            continue
        }

        n := runeio.Must(rdr.Peek())
        if isValidEscape(c, n) {
            sb.WriteRune(ConsumeEscapedCodepoint(rdr))
            continue
        }

        rdr.Push(c)
        return sb.String()
    }
}

// ConsumeNumericToken consumes a numeric token from a stream of code points.
// It returns either a <number-token>, <percentage-token>, or
// <dimension-token>.
func ConsumeNumericToken(rdr *runeio.Reader) Token {
    // Consume a number and let number be the result.
    nt, repr, value := ConsumeNumber(rdr)

    // If the next 3 input code points would start an ident sequence:
    var xs [3]rune
    must.Result(rdr.PeekN(xs[:], 3))
    if isStartOfIdentSequence(xs[0], xs[1], xs[2]) {
        // Create a <dimension-token> with the same value and type flag as
        // number, and a unit set initially to the empty string.
        // Consume an ident sequence. Set the <dimension-token>’s unit to the
        // returned value. Return the <dimension-token>.
        unit := ConsumeIdentSequence(rdr)
        return NewDimensionToken(nt, repr, value, unit)
    }

    // Otherwise, if the next input code point is U+0025
    // PERCENTAGE SIGN (%), consume it. Create a <percentage-token> with the same
    // value as number, and return it.
    if xs[0] == '%' {
        rdr.Skip(1)
        return NewPercentageToken(nt, repr, value)
    }

    // Otherwise, create a <number-token> with the same value and type flag as
    // number, and return it.
    return NewNumberToken(nt, repr, value)
}

// ConsumeNumber consumes a number from a stream of code points. It returns a
// representation, a numeric value, and a type which is either "integer" or
// "number".
//
// The representation is the token lexeme as it appears in the input stream.
// This preserves details // such as whether .009 was written as .009 or 9e-3.
//
// Note: This algorithm does not do the verification of the first few code
// points that are necessary to ensure a number can be obtained from the
// stream. Ensure that the stream starts with a number before calling this
// algorithm.
func ConsumeNumber(rdr *runeio.Reader) (nt NumberType, repr string, value float64) {
    // https://www.w3.org/TR/css-syntax-3/#consume-number

    // Initially set type to "integer". Let repr be the empty string.
    nt = NumberTypeInteger
    var sb strings.Builder // repr string builder

    // If the next input code point is U+002B PLUS SIGN (+) or
    // U+002D HYPHEN-MINUS (-), consume it and append it to repr.
    n := runeio.Must(rdr.Peek())
    if (n == '+') || (n == '-') {
        rdr.Skip(1)
        sb.WriteRune(n)
    }

    // While the next input code point is a digit, consume it and append it to
    // repr.
    consumeAndAppendWhile(rdr, &sb, runeIsDigit)

    // If the next 2 input code points are U+002E FULL STOP (.) followed by a
    // digit, then:
    var xs [2]rune
    must.Result(rdr.PeekN(xs[:], 2))
    if (xs[0] == '.') && runeIsDigit(xs[1]) {
        // Consume them.
        // Append them to repr.
        sb.WriteRune(runeio.Must(rdr.Next()))
        sb.WriteRune(runeio.Must(rdr.Next()))

        // Set type to "number".
        nt = NumberTypeNumber

        // While the next input code point is a digit, consume it and append it
        // to repr.
        consumeAndAppendWhile(rdr, &sb, runeIsDigit)
    }

    // If the next 2 or 3 input code points are
    // U+0045 LATIN CAPITAL LETTER E (E) or U+0065 LATIN SMALL LETTER E (e),
    // optionally followed by U+002D HYPHEN-MINUS (-) or U+002B PLUS SIGN (+),
    // followed by a digit, then:
    var eNotation int
    var e[3]rune
    must.Result(rdr.PeekN(e[:], 3))
    if (e[0] == 'E') || (e[0] == 'e') {
        if (e[1] == '-') || (e[1] == '+') {
            if runeIsDigit(e[2]) {
                eNotation = 3
            }
        } else if runeIsDigit(e[1]) {
            eNotation = 2
        }
    }
    if eNotation > 0 {
        // Consume them.
        // Append them to repr.
        for i := 0; i < eNotation; i++ {
            sb.WriteRune(runeio.Must(rdr.Next()))
        }

        // Set type to "number".
        nt = NumberTypeNumber

        // While the next input code point is a digit, consume it and append it
        // to repr.
        consumeAndAppendWhile(rdr, &sb, runeIsDigit)
    }

    // Convert repr to a number, and set the value to the returned value.
    repr = sb.String()
    value = StringToNumber(repr)
    return
}

// StringToNumber describes how to convert a string to a number according to
// the CSS specification.
//
// Note: This algorithm does not do any verification to ensure that the string
// contains only a number. Ensure that the string contains only a valid CSS
// number before calling this algorithm.
func StringToNumber(x string) float64 {
    digits := func(s string) (string, string) {
        n := 0
        for _, c := range s {
            if !runeIsDigit(c) { break }
            n++
        }
        result := s[0:n]
        s = s[n:]
        return s, result
    }

    // Divide the string into seven components, in order from left to right:

    // A sign: a single U+002B PLUS SIGN (+) or U+002D HYPHEN-MINUS (-),
    // or the empty string.
    var sign byte
    if len(x) > 0 {
        if x[0] == '+' {
            sign = x[0]
            x = x[1:]
        } else if x[0] == '-' {
            sign = x[1]
            x = x[1:]
        }
    }

    // An integer part: zero or more digits.
    var integer string
    x, integer = digits(x)

    // A decimal point: a single U+002E FULL STOP (.), or the empty string.
    if (len(x) > 0) && (x[0] == '.') {
        x = x[1:]
    }

    // A fractional part: zero or more digits
    var frac string
    x, frac = digits(x)

    var expsign byte

    // An exponent indicator: a single U+0045 LATIN CAPITAL LETTER E (E) or
    // U+0065 LATIN SMALL LETTER E (e), or the empty string.
    if len(x) > 0 {
        if (x[0] == 'E') || (x[0] == 'e') {
            x = x[1:]

            // An exponent sign: a single U+002B PLUS SIGN (+) or
            // U+002D HYPHEN-MINUS (-), or the empty string.
            if len(x) > 0 {
                if (x[0] == '+') || (x[0] == '-') {
                    expsign = x[0]
                    x = x[1:]
                }
            }
        }
    }

    // An exponent: zero or more digits.
    var exponent string
    x, exponent = digits(x)
    if len(x) > 0 { panic(fmt.Errorf("StringToNumber: unexpected trailing bytes at end of number")) }

    var s, i, f, d, t, e float64

    // Let s be the number -1 if the sign is U+002D HYPHEN-MINUS (-);
    // otherwise, let s be the number 1.
    if sign == '-' { s = -1 } else { s = 1 }

    // Integer part: If there is at least one digit, let i be the number formed
    // by interpreting the digits as a base-10 integer;
    // otherwise, let i be the number 0.
    if len(integer) > 0 {
        n, err := strconv.ParseInt(integer, 10, 64)
        if err == strconv.ErrRange { // ok, n is largest representable integer
        } else if err != nil {
            panic(fmt.Errorf("StringToNumber: invalid integer component"))
        }
        i = float64(n)
    }

    // Fractional part: If there is at least one digit, let f be the number
    // formed by interpreting the digits as a base-10 integer and d be the
    // number of digits; otherwise, let f and d be the number 0.
    if len(frac) > 0 {
        n, err := strconv.ParseInt(integer, 10, 64)
        if err == strconv.ErrRange { // ok, n is largest representable integer
        } else if err != nil {
            panic(fmt.Errorf("StringToNumber: invalid fractional component"))
        }
        f = float64(n)
        d = float64(len(frac))
    }

    // Let t be the number -1 if the exponent sign is U+002D HYPHEN-MINUS (-);
    // otherwise, let t be the number 1.
    if expsign == '-' { t = -1 } else { t = 1 }

    // Exponent: If there is at least one digit, let e be the number formed by
    // interpreting the digits as a base-10 integer;
    // otherwise, let e be the number 0.
    if len(exponent) > 0 {
        n, err := strconv.ParseInt(integer, 10, 64)
        if err == strconv.ErrRange { // ok, n is largest representable integer
        } else if err != nil {
            panic(fmt.Errorf("StringToNumber: invalid fractional component"))
        }
        e = float64(n)
    }

    // Return the number s·(i + f·10^(-d))·10^(te).
    part := i + (f * math.Pow(10, -d))
    return s * part * math.Pow(10, t * e)
}

// ConsumeIdentLikeToken consumes an ident-like token from a stream of code
// points. It returns an <ident-token>, <function-token>, <url-token>, or
// <bad-url-token>.
func ConsumeIdentLikeToken(rdr *runeio.Reader) (Token, error) {
    // Consume an ident sequence, and let string be the result.
    ident := ConsumeIdentSequence(rdr)

    // If string’s value is an ASCII case-insensitive match for "url", and the next
    // input code point is U+0028 LEFT PARENTHESIS (()...
    if strings.EqualFold(ident, "url") && '(' == runeio.Must(rdr.Peek()) {
        // consume it.
        rdr.Skip(1)

        // While the next two input code points are whitespace, consume the
        // next input code point.
        for {
            var xs [2]rune
            must.Result(rdr.PeekN(xs[:], 2))
            if runeIsWhitespace(xs[0]) && runeIsWhitespace(xs[1]) {
                rdr.Skip(1)
                continue
            }
            break
        }

        // If the next one or two input code points are
        // U+0022 QUOTATION MARK ("), U+0027 APOSTROPHE ('),
        // or whitespace followed by
        // U+0022 QUOTATION MARK (") or U+0027 APOSTROPHE ('),
        // then
        var xs [2]rune
        must.Result(rdr.PeekN(xs[:], 2))
        var isFuncToken bool

        if (xs[0] == '"') || (xs[0] == '\'') { isFuncToken = true }
        if runeIsWhitespace(xs[0]) && ((xs[1] == '"') || (xs[1] == '\'')) {
            isFuncToken = true
        }

        if isFuncToken {
            // create a <function-token>
            // with its value set to string and return it.
            return NewFunctionToken(ident), nil
        } else {
            // Otherwise, consume a url token, and return it.
            return ConsumeUrlToken(rdr)
        }

        // Otherwise, if the next input code point is
        // U+0028 LEFT PARENTHESIS (()...
    } else if runeio.Must(rdr.Peek()) == '(' {
        // consume it.
        rdr.Skip(1)
        // Create a <function-token> with its value set to string and return it.
        return NewFunctionToken(ident), nil
        // Otherwise...
    } else {
        // create an <ident-token> with its value set to string and return it.
        return NewIdentToken(ident), nil
    }
}

// ConsumeUrlToken describes how to consume a url token from a stream of code
// points. It returns either a <url-token> or a <bad-url-token>.
//
// Note: This algorithm assumes that the initial "url(" has already been
// consumed. This algorithm also assumes that it’s being called to consume an
// "unquoted" value, like url(foo). A quoted value, like url("foo"), is parsed
// as a <function-token>. ConsumeIdentLikeToken automatically handles
// this distinction; this algorithm shouldn’t be called directly otherwise.
func ConsumeUrlToken(rdr *runeio.Reader) (Token, error) {
    // Initially create a <url-token> with its value set to the empty string.
    var sb strings.Builder

    // Consume as much whitespace as possible.
    ConsumeWhitespace(rdr)

    // Repeatedly consume the next input code point from the stream:
    for {
        c := runeio.Must(rdr.Next())
        switch {
            case c == ')': // U+0029 RIGHT PARENTHESIS ())
                return NewUrlToken(sb.String()), nil
            case c == runeio.RuneEOF:
                // This is a parse error. Return the <url-token>.
                return NewUrlToken(sb.String()), ErrUnexpectedEOF
            case runeIsWhitespace(c):
                // // Consume as much whitespace as possible.
                ConsumeWhitespace(rdr)
                // If the next input code point is
                // U+0029 RIGHT PARENTHESIS ()) or EOF...
                p := runeio.Must(rdr.Peek())
                // consume it and return the <url-token>
                // (if EOF was encountered, this is a parse error);
                if (p == ')') {
                    rdr.Skip(1)
                    return NewUrlToken(sb.String()), nil
                } else if (p == runeio.RuneEOF) {
                    return NewUrlToken(sb.String()), ErrUnexpectedEOF
                    // otherwise
                } else {
                    // consume the remnants of a bad url,
                    // create a <bad-url-token>, and return it.
                    ConsumeBadUrl(rdr)
                    return Token{Type: TokenTypeBadUrl}, ErrBadUrl
                }
            case c == '"':  fallthrough
            case c == '\'': fallthrough
            case c == '(':  fallthrough
            case isNonPrintable(c):
                // This is a parse error. Consume the remnants of a bad url,
                // create a <bad-url-token>, and return it.
                    ConsumeBadUrl(rdr)
                    return Token{Type: TokenTypeBadUrl}, ErrBadUrl
            case c == '\\': // U+005C REVERSE SOLIDUS (\)
                // If the stream starts with a valid escape...
                p := runeio.Must(rdr.Peek())
                if isValidEscape(c, p) {
                    // consume an escaped code point and append the returned
                    // code point to the <url-token>’s value.
                    sb.WriteRune(ConsumeEscapedCodepoint(rdr))
                    // Otherwise...
                } else {
                    // this is a parse error.
                    // Consume the remnants of a bad url,
                    // create a <bad-url-token>, and return it.
                    ConsumeBadUrl(rdr)
                    return Token{Type: TokenTypeBadUrl}, ErrBadUrl
                }
            default: // anything else
                // Append the current input code point to the <url-token>’s value.
                sb.WriteRune(c)
        }
    }
}

// ConsumeBadUrl consumes the remnants of a bad url from a stream of code
// points, "cleaning up" after the tokenizer realizes that it’s in the middle
// of a <bad-url-token> rather than a <url-token>. It returns nothing; its sole
// use is to consume enough of the input stream to reach a recovery point where
// normal tokenizing can resume.
func ConsumeBadUrl(rdr *runeio.Reader) {
    for {
        // Repeatedly consume the next input code point from the stream:
        c := runeio.Must(rdr.Next())
        p := runeio.Must(rdr.Peek())
        if c == ')' || c == runeio.RuneEOF {
            return
        } else if isValidEscape(c, p) {
            ConsumeEscapedCodepoint(rdr)
        }
    }
}

func consumeAndAppendWhile(rdr *runeio.Reader, builder *strings.Builder, pred func(x rune) bool) {
        for pred(runeio.Must(rdr.Peek())) {
            builder.WriteRune(runeio.Must(rdr.Next()))
        }
}

func NewStringToken(s string) Token {
    return Token{
        Type:        TokenTypeString,
        stringValue: s,
    }
}

func NewBadStringToken(s string) Token {
    return Token{
        Type:        TokenTypeBadString,
        stringValue: s,
    }
}

func NewDelimToken(x rune) Token {
    return Token{
        Type:  TokenTypeDelim,
        delim: x,
    }
}

func NewHashToken(t HashType, s string) Token {
    return Token{
        Type:        TokenTypeHash,
        stringValue: s,
        hashType:    t,
    }
}

func NewNumberToken(nt NumberType, repr string, value float64) Token {
    return Token{
        Type:        TokenTypeNumber,
        repr:        repr,
        numberValue: value,
        numberType:  nt,
    }
}

func NewPercentageToken(nt NumberType, repr string, value float64) Token {
    return Token{
        Type:        TokenTypePercentage,
        repr:        repr,
        numberValue: value,
        numberType:  nt,
    }
}

func NewDimensionToken(nt NumberType, repr string, value float64, unit string) Token {
    return Token{
        Type:        TokenTypeDimension,
        repr:        repr,
        numberValue: value,
        numberType:  nt,
        unit:        unit,
    }
}

func NewIdentToken(s string) Token {
    return Token{
        Type:        TokenTypeIdent,
        stringValue: s,
    }
}

func NewFunctionToken(s string) Token {
    return Token{
        Type:        TokenTypeFunction,
        stringValue: s,
    }
}

func NewUrlToken(s string) Token {
    return Token{
        Type:        TokenTypeUrl,
        stringValue: s,
    }
}

func NewAtKeywordToken(s string) Token {
    return Token{
        Type:        TokenTypeAtKeyword,
        stringValue: s,
    }
}