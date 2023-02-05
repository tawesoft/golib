package tokenizer_test

// NOTE: some tests, marked "from chromium.googlesource.com", are based on
// test cases given at:
// https://chromium.googlesource.com/chromium/src/+/22eeef8fc52576bf54a81b39555030eea9629d35/third_party/blink/renderer/core/css/parser/css_tokenizer_test.cc
// (see LICENSE-PARTS.txt).
//
// Note that tests have some differences where specified, and additional
// error detection has been added to some tests. Additionally, we have
// different number token type representations.

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
        tok := t.NextExcept(token.TypeWhitespace)
        if tok.Is(token.TypeEOF) { break }
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
    // based on python "isClose"
    // https://peps.python.org/pep-0485/#proposed-implementation
    reltol := 1e-12
    abstol := 0.0
    return math.Abs(a-b) <= math.Max(reltol * math.Max(math.Abs(a), math.Abs(b)), abstol)
}

func equal(expected token.Token, actual token.Token) bool {
    if !token.Equals(expected, actual) {
        return false
    }

    if expected.IsNumeric() {
        _, ev := expected.NumericValue()
        _, av := actual.NumericValue()
        if !roughlyEqual(ev, av) {
            return false
        }
    }

    return true
}

func Test_Tester(t *testing.T) {
    // asserts the tests actually test what they claim to test!
    assert.True(t, roughlyEqual(2.0, 6.0/3.0))
    assert.True(t, roughlyEqual(1.2, 1.2))
    assert.True(t, roughlyEqual(12.34e2, 12.34e2))
    assert.True(t, roughlyEqual(12.34e10, 12.34e10))
    assert.False(t, roughlyEqual(2.1, 2.0))
    assert.False(t, roughlyEqual(12.34e6+1, 12.34e6))
    assert.False(t, roughlyEqual(math.MaxInt32 - 1, math.MaxInt32))
    assert.True(t,  roughlyEqual(math.MaxInt32, math.MaxInt32))
    assert.False(t,  roughlyEqual(float64(math.MaxInt32) + 1.0, math.MaxInt32))
    assert.True(t, equal(
        token.Number(token.NumberTypeInteger, "123", 123),
        token.Number(token.NumberTypeInteger, "123", 123),
    ))
    assert.False(t, equal(
        token.Number(token.NumberTypeInteger, "123", 123),
        token.Number(token.NumberTypeInteger, "124", 124),
    ))
    assert.True(t, equal(
        token.Number(token.NumberTypeNumber, "123.456", 123.456),
        token.Number(token.NumberTypeNumber, "123.456", 123.456),
    ))
    assert.False(t, equal(
        token.Number(token.NumberTypeNumber, "123.4", 123.4),
        token.Number(token.NumberTypeNumber, "123.456", 123.456),
    ))
    assert.False(t, equal(
        token.Number(token.NumberTypeInteger, "123", 123),
        token.Number(token.NumberTypeNumber,  "123", 123),
    ))
    assert.True(t, equal(
        token.Number(token.NumberTypeInteger, "12E-1", 1.2),
        token.Number(token.NumberTypeInteger, "12E-1", 1.2),
    ))
}

func testWithErrCheck(t *testing.T, css string, errCheck func(errors []error) bool, tokens ... token.Token) {
    p := tokenizer.New(strings.NewReader(css))
    func() {
        seen := make([]token.Token, 0)
        err := func(msg string) {
            t.Errorf("%s\n    input: %q\n    expected: %v\n    seen: %v", msg, css, tokens, seen)
        }

        for _, k := range tokens {
            n := p.Next()
            if n.Is(token.TypeEOF) { break }

            seen = append(seen, n)
            if !equal(k, n) {
                err("parse error")
                return
            }
        }

        if len(seen) != len(tokens) {
            err("unexpected tokenizer termination")
            return
        }
        n := p.Next()
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

func TestTokenizer_Delimiters(t *testing.T) {
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

func TestTokenizer_Whitespace(t *testing.T) {
    // Tests from chromium.googlesource.com
    test(t, "   ", token.Whitespace())
    test(t, "\n\rS", token.Whitespace(), token.Ident("S"))
    test(t, "   *", token.Whitespace(), token.Delim('*'))
    test(t, "\r\n\f\t2", token.Whitespace(),
        token.Number(token.NumberTypeInteger, "2", 2.0))
}

func TestTokenizer_Escapes(t *testing.T) {
    // Tests from chromium.googlesource.com
    replacement := string([]rune{0xFFFD})
    test(t, "hel\\6Co",         token.Ident("hello"))
    test(t, "\\26 B",           token.Ident("&B"))
    test(t, "'hel\\6c o'",      token.String("hello"))
    test(t, "'spac\\65\r\ns'",  token.String("spaces"))
    test(t, "spac\\65\r\ns",    token.Ident("spaces"))
    test(t, "spac\\65\n\rs",    token.Ident("space"),
        token.Whitespace(), token.Ident("s"))
    test(t, "sp\\61\tc\\65\fs", token.Ident("spaces"))
    test(t, "hel\\6c  o",       token.Ident("hell"),
        token.Whitespace(), token.Ident("o"))
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
    testWithErrCheck(t, "test\\\n", check,
        token.Ident("test"), token.Delim('\\'), token.Whitespace())
    testWithErrCheck(t, "\\\r",     check, token.Delim('\\'), token.Whitespace())
    testWithErrCheck(t, "\\\f",     check, token.Delim('\\'), token.Whitespace())
    testWithErrCheck(t, "\\\r\n",   check, token.Delim('\\'), token.Whitespace())
}

func TestTokenizer_Idents(t *testing.T) {
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
    test(t, "--",               token.Ident("--"))
    test(t, "--11",             token.Ident("--11"))
    test(t, "---",              token.Ident("---"))
    test(t, "\u2003",           token.Ident(string([]rune{0x2003}))) // em-space
    test(t, "\u00A0",           token.Ident(string([]rune{0x00A0}))) // non-breaking space
    test(t, "\u1234",           token.Ident(string([]rune{0x1234})))
    test(t, "\U00012345",       token.Ident(string([]rune{0x12345})))
    test(t, "\000",             token.Ident(string([]rune{0xFFFD})))
    test(t, "ab\000c",          token.Ident("ab" + string([]rune{0xFFFD}) + "c"))
    test(t, "ab\000c",          token.Ident("ab" + string([]rune{0xFFFD}) + "c"))
}

func TestTokenizer_Functions(t *testing.T) {
    // Tests from chromium.googlesource.com
    test(t, "scale(2)",         token.Function("scale"),
        token.Number(token.NumberTypeInteger, "2", 2),
        token.RightParen())
    test(t, "foo-bar\\ baz(",   token.Function("foo-bar baz"))
    test(t, "fun\\(ction(",     token.Function("fun(ction"))
    test(t, "-foo(", token.Function("-foo"))
    test(t, "url(\"foo.gif\"",  token.Function("url"), token.String("foo.gif"))
    test(t, "foo(  'bar.gif'",  token.Function("foo"), token.Whitespace(), token.String("bar.gif"))
    // unlike the chromium test, we don't drop the whitespace after Function("url") here:
    test(t, "url(  'bar.gif'",  token.Function("url"), token.Whitespace(), token.String("bar.gif"))
}

func TestTokenizer_AtKeywords(t *testing.T) {
    // Tests from chromium.googlesource.com
    test(t, "@at-keyword",      token.AtKeyword("at-keyword"))
    test(t, "@testing123",      token.AtKeyword("testing123"))
    test(t, "@hello!",          token.AtKeyword("hello"), token.Delim('!'))
    test(t, "@-text",           token.AtKeyword("-text"))
    test(t, "@--abc",           token.AtKeyword("--abc"))
    test(t, "@--",              token.AtKeyword("--"))
    test(t, "@--11",            token.AtKeyword("--11"))
    test(t, "@---",             token.AtKeyword("---"))
    test(t, "@\\ ",             token.AtKeyword(" "))
    test(t, "@-\\ ",            token.AtKeyword("- "))
    test(t, "@@",               token.Delim('@'), token.Delim('@'))
    test(t, "@2",               token.Delim('@'), token.Number(token.NumberTypeInteger,  "2",  2))
    test(t, "@-1",              token.Delim('@'), token.Number(token.NumberTypeInteger, "-1", -1))
}

func TestTokenizer_Urls(t *testing.T) {
    // Tests from chromium.googlesource.com
    test(t, "url(foo.gif)",                         token.Url("foo.gif"))
    test(t, "urL(https://example.com/cats.png)",    token.Url("https://example.com/cats.png"))
    test(t, "uRl(what-a.strange^URL~this\\ is!)",   token.Url("what-a.strange^URL~this is!"))
    test(t, "uRL(123#test)",                        token.Url("123#test"))
    test(t, "Url(escapes\\ \\\"\\'\\)\\()",         token.Url("escapes \"')("))
    test(t, "UrL(   whitespace   )",                token.Url("whitespace"))
    test(t, "url(not/*a*/comment)",                 token.Url("not/*a*/comment"))
    test(t, "urL()",                                token.Url(""))

    // the following tests tokenize successfully, but only because they
    // recover from a specific error. Assert (exactly!) that error is reported.
    checkEof := func(errs []error) bool {
        return (len(errs) == 1) && (errors.Is(errs[0], tokenizer.ErrUnexpectedEOF))
    }
    checkSyn := func(errs []error) bool {
        return (len(errs) == 1) && (errors.Is(errs[0], tokenizer.ErrBadUrl))
    }
    testWithErrCheck(t, "URL(eof",                  checkEof, token.Url("eof"))
    testWithErrCheck(t, "URl( whitespace-eof ",     checkEof, token.Url("whitespace-eof"))
    testWithErrCheck(t, "uRl(white space),",        checkSyn, token.BadUrl(), token.Comma())
    testWithErrCheck(t, "Url(b(ad),",               checkSyn, token.BadUrl(), token.Comma())
    testWithErrCheck(t, "uRl(ba'd):",               checkSyn, token.BadUrl(), token.Colon())
    testWithErrCheck(t, "urL(b\"ad):",              checkSyn, token.BadUrl(), token.Colon())
    testWithErrCheck(t, "uRl(b\"ad):",              checkSyn, token.BadUrl(), token.Colon())
    testWithErrCheck(t, "Url(b\\\rad):",            checkSyn, token.BadUrl(), token.Colon())
    testWithErrCheck(t, "url(b\\\nad):",            checkSyn, token.BadUrl(), token.Colon())
    testWithErrCheck(t, "url(/*'bad')*/",           checkSyn, token.BadUrl(),
        token.Delim('*'), token.Delim('/'))
    testWithErrCheck(t, "url(ba'd\\\\))",           checkSyn, token.BadUrl(), token.RightParen())
}

func TestTokenizer_Strings(t *testing.T) {
    // Tests from chromium.googlesource.com
    test(t, "'text'",           token.String("text"))
    test(t, "\"text\"",         token.String("text"))
    test(t, "'testing, 123!'",  token.String("testing, 123!"))
    test(t, "'es\\'ca\\\"pe'",  token.String("es'ca\"pe"))
    test(t, "'\"quotes\"'",     token.String("\"quotes\""))
    test(t, "\"'quotes'\"",     token.String("'quotes'"))
    test(t, "'text\005\t\013'", token.String("text\005\t\013"))
    test(t, "'esca\\\nped'",    token.String("escaped"))
    test(t, "\"esc\\\faped\"",  token.String("escaped"))
    test(t, "'new\\\rline'",    token.String("newline"))
    test(t, "\"new\\\r\nline\"",token.String("newline"))
    test(t, "'\000'",           token.String(string([]rune{0xFFFD})))
    test(t, "'hel\000lo'",      token.String("hel" + string([]rune{0xFFFD}) + "lo"))
    test(t, "'h\\065l\000lo'",  token.String("hel" + string([]rune{0xFFFD}) + "lo"))

    // the following tests tokenize successfully, but only because they
    // recover from a specific error. Assert (exactly!) that error is reported.
    checkEof := func(errs []error) bool {
        return (len(errs) == 1) && (errors.Is(errs[0], tokenizer.ErrUnexpectedEOF))
    }
    checkLbr := func(errs []error) bool {
        return (len(errs) == 1) && (errors.Is(errs[0], tokenizer.ErrUnexpectedLinebreak))
    }
    testWithErrCheck(t, "\"mismatch'",      checkEof, token.String("mismatch'"))
    testWithErrCheck(t, "\"end on eof",     checkEof, token.String("end on eof"))
    testWithErrCheck(t, "'bad\nstring",     checkLbr, token.BadString(), token.Whitespace(), token.Ident("string"))
    testWithErrCheck(t, "'bad\rstring",     checkLbr, token.BadString(), token.Whitespace(), token.Ident("string"))
    testWithErrCheck(t, "'bad\r\nstring",   checkLbr, token.BadString(), token.Whitespace(), token.Ident("string"))
    testWithErrCheck(t, "'bad\fstring",     checkLbr, token.BadString(), token.Whitespace(), token.Ident("string"))
}

func TestTokenizer_HashTokens(t *testing.T) {
    // Tests from chromium.googlesource.com
    test(t, "#id-selector",     token.Hash(token.HashTypeID, "id-selector"))
    test(t, "#FF7700",          token.Hash(token.HashTypeID, "FF7700"))
    test(t, "#3377FF",          token.Hash(token.HashTypeUnrestricted, "3377FF"))
    test(t, "#\\ ",             token.Hash(token.HashTypeID, " "))
    test(t, "# ",               token.Delim('#'), token.Whitespace())
    test(t, "#!",               token.Delim('#'), token.Delim('!'))

    // the following tests tokenize successfully, but only because they
    // recover from a specific error. Assert (exactly!) that error is reported.
    checkErr := func(errs []error) bool {
        return (len(errs) == 1) && (errors.Is(errs[0], tokenizer.ErrUnexpectedInput))
    }
    testWithErrCheck(t, "#\\\n",    checkErr, token.Delim('#'), token.Delim('\\'), token.Whitespace())
    testWithErrCheck(t, "#\\\r\n",  checkErr, token.Delim('#'), token.Delim('\\'), token.Whitespace())
}

func TestTokenizer_Numbers(t *testing.T) {
    // Tests based on chromium.googlesource.com
    test(t, "10",           token.Number(token.NumberTypeInteger, "10",     10))
    test(t, "12.0",         token.Number(token.NumberTypeNumber,  "12.0",   12))
    test(t, "+45.6",        token.Number(token.NumberTypeNumber,  "+45.6",  45.6))
    test(t, "-7",           token.Number(token.NumberTypeInteger, "-7",     -7))
    test(t, "010",          token.Number(token.NumberTypeInteger, "010",    10))
    test(t, "10e0",         token.Number(token.NumberTypeNumber,  "10e0",   10e0))
    test(t, "12e3",         token.Number(token.NumberTypeNumber,  "12e3",   12000))
    test(t, "3e+1",         token.Number(token.NumberTypeNumber,  "3e+1",   30))
    test(t, "12E-1",        token.Number(token.NumberTypeNumber, "12E-1",   1.2))
    test(t, ".7",           token.Number(token.NumberTypeNumber, ".7",      0.7))
    test(t, "-.3",          token.Number(token.NumberTypeNumber, "-.3",     -0.3))
    test(t, "+637.54e-2",   token.Number(token.NumberTypeNumber, "+637.54e-2", 6.3754))
    test(t, "-12.34E+2",    token.Number(token.NumberTypeNumber, "-12.34E+2", -1234))
    test(t, "+ 5",          token.Delim('+'), token.Whitespace(),
        token.Number(token.NumberTypeInteger, "5", 5))
    test(t, "-+12",         token.Delim('-'),
        token.Number(token.NumberTypeInteger, "+12", 12))
    test(t, "+-21",         token.Delim('+'),
        token.Number(token.NumberTypeInteger, "-21", -21))
    test(t, "++22",         token.Delim('+'),
        token.Number(token.NumberTypeInteger, "+22", 22))
    test(t, "13.",          token.Number(token.NumberTypeInteger, "13", 13),
        token.Delim('.'))
    test(t, "1.e2",         token.Number(token.NumberTypeInteger, "1", 1),
        token.Delim('.'),
        token.Ident("e2"))
    test(t, "2e3.5",        token.Number(token.NumberTypeNumber, "2e3", 2000),
        token.Number(token.NumberTypeNumber, ".5", 0.5))
    test(t, "2e3.",         token.Number(token.NumberTypeNumber, "2e3", 2000),
        token.Delim('.'))

    // Unlike chromium, we clamp to Min/MaxInt32
    test(t, "1000000000000000000000000",
        token.Number(token.NumberTypeInteger, "1000000000000000000000000", math.MaxInt32))
}

func TestTokenizer_Dimensions(t *testing.T) {
    test(t, "10px",         token.Dimension(token.NumberTypeInteger, "10",    10,    "px"))
    test(t, "12.0em",       token.Dimension(token.NumberTypeNumber,  "12.0",  12,    "em"))
    test(t, "-12.0em",      token.Dimension(token.NumberTypeNumber,  "-12.0", -12.0, "em"))
    test(t, "+45.6__qem",   token.Dimension(token.NumberTypeNumber,  "+45.6", 45.6,  "__qem"))
    test(t, "5e",           token.Dimension(token.NumberTypeInteger, "5",     5,     "e"))
    test(t, "5px-2px",      token.Dimension(token.NumberTypeInteger, "5",     5,     "px-2px"))
    test(t, "5e-",          token.Dimension(token.NumberTypeInteger, "5",     5,     "e-"))
    test(t, "5\\ ",         token.Dimension(token.NumberTypeInteger, "5",     5,     " "))
    test(t, "40\\70\\78",   token.Dimension(token.NumberTypeInteger, "40",    40,    "px"))
    test(t, "4e3e2",        token.Dimension(token.NumberTypeNumber,  "4e3",   4000,  "e2"))
    test(t, "0x10px",       token.Dimension(token.NumberTypeInteger, "0",     0,     "x10px"))
    test(t, "4unit ",       token.Dimension(token.NumberTypeInteger, "4",     4,     "unit"),
        token.Whitespace())
    test(t, "5e+",          token.Dimension(token.NumberTypeInteger, "5",     5,     "e"),
        token.Delim('+'))
    test(t, "2e.5",         token.Dimension(token.NumberTypeInteger, "2",     2,     "e"),
        token.Number(token.NumberTypeNumber, ".5", 0.5))
    test(t, "2e+.5",        token.Dimension(token.NumberTypeInteger, "2",     2,     "e"),
        token.Number(token.NumberTypeNumber, "+.5", 0.5))
}

func TestTokenizer_Percentages(t *testing.T) {
    test(t, "10%",      token.Percentage(token.NumberTypeInteger, "10",     10))
    test(t, "+12.0%",   token.Percentage(token.NumberTypeNumber,  "+12.0",  12))
    test(t, "-48.99%",  token.Percentage(token.NumberTypeNumber,  "-48.99", -48.99))
    test(t, "6e-1%",    token.Percentage(token.NumberTypeNumber,  "6e-1",   0.6))
    test(t, "5%%",      token.Percentage(token.NumberTypeInteger, "5",      5),
        token.Delim('%'))
}

// func TestTokenizer_UnicodeRange(t *testing.T)
// not implemented because CSS now defines urange in terms of other CSS tokens,
// and urange is detected at the parser step instead.

func TestTokenizer_Comments(t *testing.T) {
    test(t, "/*comment*/a",       token.Ident("a"))
    test(t, "/**\\2f**//",        token.Delim('/'))
    test(t, "/**y*a*y**/ ",       token.Whitespace())
    test(t, ",/* \n :) \n */)",   token.Comma(), token.RightParen())
    test(t, ":/*/*/",             token.Colon())
    test(t, "/**/*",              token.Delim('*'))

    // the following tests tokenize successfully, but only because they
    // recover from a specific error. Assert (exactly!) that error is reported.
    checkEof := func(errs []error) bool {
        return (len(errs) == 1) && (errors.Is(errs[0], tokenizer.ErrUnexpectedEOF))
    }
    testWithErrCheck(t, ";/******", checkEof, token.Semicolon())
}
