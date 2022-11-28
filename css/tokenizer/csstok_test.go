package tokenizer_test

// NOTE: some tests, marked "from chromium.googlesource.com", are based on
// test cases given at:
// https://chromium.googlesource.com/chromium/src/+/22eeef8fc52576bf54a81b39555030eea9629d35/third_party/blink/renderer/core/css/parser/css_tokenizer_test.cc

import (
    "errors"
    "fmt"
    "math"
    "strings"
    "testing"
    "unicode"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/v2/css/tokenizer"
    "github.com/tawesoft/golib/v2/css/tokenizer/token"
)

func ExampleTokenizer() {
    str := `/* example */
#something[rel~="external"] {
    background-color: rgb(128, 64, 64);
}`
    t := tokenizer.New(strings.NewReader(str))

    for {
        tok, ok := t.NextExcept(token.TypeWhitespace, token.TypeEOF)
        if !ok { break }
        fmt.Println(tok)
    }

    if len(t.Errors()) > 0 {
        fmt.Printf("%v\n", t.Errors())
    }

    // Output:
    // <hash-token>{type: "id", value: "something"}
    // <[-token>
    // <ident-token>{value: "rel"}
    // <delim-token>{delim: '~'}
    // <delim-token>{delim: '='}
    // <string-token>{value: "external"}
    // <]-token>
    // <{-token>
    // <ident-token>{value: "background-color"}
    // <colon-token>
    // <function-token>{value: "rgb"}
    // <number-token>{type: "integer", value: 128.000000, repr: "128"}
    // <comma-token>
    // <number-token>{type: "integer", value: 64.000000, repr: "64"}
    // <comma-token>
    // <number-token>{type: "integer", value: 64.000000, repr: "64"}
    // <)-token>
    // <semicolon-token>
    // <}-token>
}

func roughlyEqual(a float64, b float64) bool {
    epsilon := math.Nextafter(1, 2) - 1
    return math.Abs(a - b) < epsilon
}

func equal(t *testing.T, expected token.Token, actual token.Token) bool {
    if !token.Equals(expected, actual) {
        return false
    }

    if expected.IsNumeric() {
        ev, _ := expected.NumericValue()
        av, _ := actual.NumericValue()
        if !roughlyEqual(ev, av) {
            return false
        }
    }

    return true
}

func Test_Tester(t *testing.T) {
    assert.True(t, roughlyEqual(2.0, 6.0/3.0))
    assert.False(t, roughlyEqual(2.1, 6.0/3.0))
}

func testWithErrCheck(t *testing.T, css string, errCheck func(errors []error) bool, tokens ... token.Token) {
    p := tokenizer.New(strings.NewReader(css))
    func() {
        seen := make([]token.Token, 0)
        err := func(msg string) {
            t.Errorf("%s\n    input: %q\n    expected: %v\n    seen: %v", msg, css, tokens, seen)
        }

        for _, k := range tokens {
            n, ok := p.Next()
            if !ok { break }

            seen = append(seen, n)
            if !equal(t, k, n) {
                err("parse error")
                return
            }
        }

        if len(seen) != len(tokens) {
            err("unexpected tokenizer termination")
            return
        }
        n, ok := p.Next()
        if !ok {
            err("unexpected tokenizer termination")
            return
        }
        if !n.Is(token.TypeEOF) {
            seen = append(seen, n)
            err("expected EOF")
        }
    }()

    if (len(p.Errors()) > 0) && !errCheck(p.Errors()) {
        t.Errorf("parser error:\n    input: %q\n    errors: %v", css, p.Errors())
    }
}

func test(t *testing.T, css string, tokens ... token.Token) {
    testWithErrCheck(t, css, func([]error) bool { return false }, tokens...)
}

func TestTokenizer_SingleCharacterTokens(t *testing.T) {
    // Tests from chromium.googlesource.com
    test(t, "(",  token.LeftParen())
    test(t, ")",  token.RightParen())
    test(t, "[",  token.LeftSquareBracket())
    test(t, "]",  token.RightSquareBracket())
    test(t, ",",  token.Comma())
    test(t, ":",  token.Colon())
    test(t, ";",  token.Semicolon())
    test(t, ")[", token.RightParen(), token.LeftSquareBracket())
    test(t, "[)", token.LeftSquareBracket(), token.RightParen())
    test(t, "{}", token.LeftCurlyBracket(), token.RightCurlyBracket())
    test(t, ",,", token.Comma(), token.Comma())
}

func TestTokenizer_MultiCharacterTokens(t *testing.T) {
    // Tests from chromium.googlesource.com
    // -- but not all, as they tokenize at a slightly different level
    // instead of at the CSS spec level.
    test(t, "<!--",  token.CDO())
    test(t, "<!---", token.CDO(), token.Delim('-'))
    test(t, "-->",   token.CDC())
}

func TestTokenizer_DelimeterTokens(t *testing.T) {
    // Tests from chromium.googlesource.com
    test(t, "^", token.Delim('^'))
    test(t, "*", token.Delim('*'))
    test(t, "%", token.Delim('%'))
    test(t, "~", token.Delim('~'))
    test(t, "&", token.Delim('&'))
    test(t, "|", token.Delim('|'))
    test(t, "\x7F", token.Delim(0x7F))
    test(t, "\x01", token.Delim(0x01))
    test(t, "~-", token.Delim('~'), token.Delim('-'))
    test(t, "^|", token.Delim('^'), token.Delim('|'))
    test(t, "$~", token.Delim('$'), token.Delim('~'))
    test(t, "*^", token.Delim('*'), token.Delim('^'))
}

func TestTokenizer_WhitespaceTokens(t *testing.T) {
    // Tests from chromium.googlesource.com
    test(t, "   ", token.Whitespace())
    test(t, "\n\rS", token.Whitespace(), token.Ident("S"))
    test(t, "   *", token.Whitespace(), token.Delim('*'))
    test(t, "\r\n\f\t2", token.Whitespace(), token.Number(token.NumberTypeInteger, "2", 2.0))
}

func TestTokenizer_Escapes(t *testing.T) {
    // Tests from chromium.googlesource.com
    replacement := string([]rune{0xFFFD});
    test(t, "hel\\6Co",         token.Ident("hello"))
    test(t, "\\26 B",           token.Ident("&B"))
    test(t, "'hel\\6c o'",      token.String("hello"))
    test(t, "'spac\\65\r\ns'",  token.String("spaces"))
    test(t, "spac\\65\r\ns",    token.Ident("spaces"))
    test(t, "spac\\65\n\rs",    token.Ident("space"), token.Whitespace(), token.Ident("s"))
    test(t, "sp\\61\tc\\65\fs", token.Ident("spaces"))
    test(t, "hel\\6c  o",       token.Ident("hell"), token.Whitespace(), token.Ident("o"))
    test(t, "test\\D799",       token.Ident("test\uD799"))
    test(t, "\\E000",           token.Ident("\uE000"))
    test(t, "te\\s\\t",         token.Ident("test"))
    test(t, "spaces\\ in\\\tident", token.Ident("spaces in\tident"))
    test(t, "\\.\\,\\:\\!",     token.Ident(".,:!"))
    test(t, "null\\\000",       token.Ident("null" + replacement))
    test(t, "null\\\000\000",   token.Ident("null" + replacement + replacement))
    test(t, "null\\0",          token.Ident("null" + replacement))
    test(t, "null\\0000",       token.Ident("null" + replacement))
    test(t, "large\\110000",    token.Ident("large" + replacement))
    test(t, "large\\23456a",    token.Ident("large" + replacement))
    test(t, "surrogate\\D800",  token.Ident("surrogate" + replacement))
    test(t, "surrogate\\0DABC", token.Ident("surrogate" + replacement))
    test(t, "\\00DFFFsurrogate",token.Ident(replacement + "surrogate"))
    test(t, "\\10fFfF",         token.Ident(string([]rune{unicode.MaxRune})))
    test(t, "\\10fFfF0",        token.Ident(string([]rune{unicode.MaxRune}) + "0"))
    test(t, "\\10000000",       token.Ident(string([]rune{0x100000}) + "00"))
    test(t, "eof\\",            token.Ident("eof" + replacement))

    // the following tests tokenize successfully, but only because they
    // recover from a syntax error. Assert (exactly!) that error is reported.
    check := func(errs []error) bool {
        return (len(errs) == 1) && (errors.Is(errs[0], tokenizer.ErrUnexpectedInput))
    }
    testWithErrCheck(t, "test\\\n", check, token.Ident("test"), token.Delim('\\'), token.Whitespace())
    testWithErrCheck(t, "\\\r",     check, token.Delim('\\'), token.Whitespace())
    testWithErrCheck(t, "\\\f",     check, token.Delim('\\'), token.Whitespace())
    testWithErrCheck(t, "\\\r\n",   check, token.Delim('\\'), token.Whitespace())
}

func TestTokenizer_IdentToken(t *testing.T) {
    // Tests from chromium.googlesource.com
    test(t, "simple-ident",     token.Ident("simple-ident"))
    test(t, "testing123",       token.Ident("testing123"))
    test(t, "hello!",           token.Ident("hello"), token.Delim('!'))
    test(t, "world\005",        token.Ident("world"), token.Delim('\005'))
    test(t, "_under score",     token.Ident("_under"), token.Whitespace(), token.Ident("score"))
    test(t, "-_underscore",     token.Ident("-_underscore"))
    test(t, "-text",            token.Ident("-text"))
    test(t, "-\\6d",            token.Ident("-m"))
    test(t, "--abc",            token.Ident("--abc"))
    fmt.Println(">>> --")
    test(t, "--",               token.Ident("--"))
    fmt.Println(">>> --11")
    test(t, "--11",             token.Ident("--11"))
    fmt.Println(">>> ---")
    test(t, "---",              token.Ident("---"))
    test(t, "\u2003",           token.Ident(string([]rune{0x2003}))) // em-space
    test(t, "\u00A0",           token.Ident(string([]rune{0x00A0}))) // non-breaking space
    test(t, "\u1234",           token.Ident(string([]rune{0x1234})))
    test(t, "\U00012345",       token.Ident(string([]rune{0x12345})))
    test(t, "\000",             token.Ident(string([]rune{0xFFFD})))
    test(t, "ab\000c",          token.Ident("ab" + string([]rune{0xFFFD}) + "c"))
    test(t, "ab\000c",          token.Ident("ab" + string([]rune{0xFFFD}) + "c"))
}
