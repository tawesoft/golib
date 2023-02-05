package color

import (
    "fmt"
    "strings"

    "github.com/tawesoft/golib/v2/css/tokenizer"
    "github.com/tawesoft/golib/v2/css/tokenizer/token"
    "github.com/tawesoft/golib/v2/fun/maybe"
)

// Tokenizer produces CSS tokens - it is implemented, for example, by a
// [tokenizer.Tokenizer].
type Tokenizer interface {
    Next() token.Token
}

type errSyntax struct {
    err error
    at token.Position
}

func (e errSyntax) Unwrap() error {
    return e.err
}

func (e errSyntax) Error() string {
    return fmt.Sprintf("error at %+v: %s", e.at, e.err.Error())
}

const maxFunctionArgs = 7 // e.g. color("RGB" 1 1 1 / 1) or rgba(1 , 2 , 3 , 4)
var (
    ErrTooManyFunctionArguments  = fmt.Errorf("too many function arguments")
    ErrSyntax                    = fmt.Errorf("invalid color syntax")
    ErrUnexpectedEOF             = fmt.Errorf("unexpected end of file")
    ErrUnexpectedTrailing        = fmt.Errorf("unexpected trailing input")
    ErrNotSupportedNamedOrSystem = fmt.Errorf("named and system colors not supported")
    ErrUnrecognisedFunction      = fmt.Errorf("unrecognised function")
    ErrInvalidArguments          = fmt.Errorf("invalid function arguments")
    ErrInvalidHex                = fmt.Errorf("invalid hexadecimal color")
)

func nextExcept(tokenizer Tokenizer, exclude ... token.Type) token.Token {
    for {
        skip:
        t := tokenizer.Next()
        for _, ex := range exclude {
            if t.Is(ex) {
                goto skip
            }
            return t
        }
    }
}

func nextExceptWS(tokenizer Tokenizer) token.Token {
    return nextExcept(tokenizer, token.TypeWhitespace)
}

// consumeSimpleFunctionArgs consumes a function, but does not allow recursion
// into blocks or functions in the arguments.
func consumeSimpleFunctionArgs(tokenizer Tokenizer) (args []token.Token, err error) {
    for {
        t := nextExceptWS(tokenizer)
        switch {
            case t.Is(token.TypeRightParen):
                return args, nil
            case t.Is(token.TypeEOF):
                err = ErrUnexpectedEOF
                return args, err
            case t.Is(token.TypeNumber):     fallthrough
            case t.Is(token.TypeDimension):  fallthrough
            case t.Is(token.TypePercentage): fallthrough
            case t.Is(token.TypeDelim):      fallthrough
            case t.Is(token.TypeString):     fallthrough
            case t.Is(token.TypeComma):      fallthrough
            case t.Is(token.TypeWhitespace): fallthrough
            case t.Is(token.TypeIdent):
                if len(args) + 1 > maxFunctionArgs {
                    return args, errSyntax{
                        err: ErrTooManyFunctionArguments,
                        at:  t.Position(),
                    }
                }
                args = append(args, t)
        }
    }
}

// ParseColorString parses a color value from a string containing a color in
// CSS syntax.
func ParseColorString(s string) (Color, error) {
    k := tokenizer.New(strings.NewReader(s))
    color, err := ParseColor(k)
    if err == nil {
        t := nextExceptWS(k)
        if !t.Is(token.TypeEOF) {
            err = errSyntax{ErrUnexpectedTrailing, t.Position()}
        }
    }
    if len(k.Errors()) > 0 {
        return color, fmt.Errorf("parse errors: %+v\n", k.Errors())
    }
    return color, err
}

// ParseColor parses a color value from CSS tokens.
//
// Note: the caller should ensure that, after a valid color is returned, the
// tokenizer produces token.EOF() (possibly preceded by whitespace) if
// necessary when parsing an entire input as a color.
func ParseColor(tokenizer Tokenizer) (Color, error) {
    zero := Color{}
    t := nextExcept(tokenizer, token.TypeWhitespace)
    if t.Is(token.TypeHash) { // #RRGGBB format
        return parseColorFromHexadecimalString(t)
    } else if t.Is(token.TypeFunction) {
        args, err := consumeSimpleFunctionArgs(tokenizer)
        if err != nil { return zero, err }
        return parseColorFromFunction(t, args)
    } else if t.Is(token.TypeIdent) { // e.g. "red"
        return zero, errSyntax{ErrNotSupportedNamedOrSystem, t.Position()}
    } else {
        return zero, errSyntax{ErrSyntax, t.Position()}
    }
}

func parseColorFromHexadecimalString(t token.Token) (color Color, err error) {
    x := t.StringValue()
    digit := func(x byte) uint8 {
        if x >= '0' && x <= '9' {
            return x - '0'
        } else if x >= 'a' && x <= 'f' {
            return x - 'a' + 10
        } else if x >= 'A' && x <= 'F' {
            return x - 'A' + 10
        } else {
            err = errSyntax{ErrInvalidHex, t.Position()}
            return 0
        }
    }

    // scaleHex2 turns "AB" into 0xAB
    scaleHex2 := func(a byte, b byte) uint8{
        u := digit(a)
        v := digit(b)
        return ((u * 16.0) + v)
    }

    // scaleHex1 turns 'A' into 0xAA
    scaleHex1 := func(a byte) uint8 {
        return scaleHex2(a, a)
    }

    var r, g, b, a uint8
    a = 255
    switch len(x) {
        case 4:
            a = scaleHex1(x[3])
            fallthrough
        case 3:
            r = scaleHex1(x[0])
            g = scaleHex1(x[1])
            b = scaleHex1(x[2])
        case 8:
            a = scaleHex2(x[6], x[7])
            fallthrough
        case 6:
            r = scaleHex2(x[0], x[1])
            g = scaleHex2(x[2], x[3])
            b = scaleHex2(x[4], x[5])
        default:
            err = errSyntax{ErrInvalidHex, t.Position()}
    }
    if err != nil { return Color{}, err }
    return Hexadecimal(r, g, b, a), nil
}

func parseColorFromFunction(f token.Token, args []token.Token) (Color, error) {
    zero := Color{}
    name := f.StringValue()
    switch {
        case strings.EqualFold(name, "rgb"): fallthrough
        case strings.EqualFold(name, "rgba"):
            return parseRGBFromFunction(f, args)
    default:
        return zero, errSyntax{ErrUnrecognisedFunction, f.Position()}
    }
    return zero, nil
}

func step(args []token.Token) (next token.Token, rest []token.Token) {
    if len(args) == 0 { return token.EOF(), nil }
    return args[0], args[1:]
}

func acceptEither(
    t token.Token,
    acceptors ... func(t token.Token) (maybe.M[float64], bool),
) (maybe.M[float64], bool) {
    for _, acceptor := range acceptors {
        value, ok := acceptor(t)
        if ok { return value, ok }
    }
    return maybe.Nothing[float64](), false
}

func numericAcceptor(_type token.Type, scale float64) func(t token.Token) (maybe.M[float64], bool) {
    return func(t token.Token) (maybe.M[float64], bool) {
        if t.Is(_type) {
            _, nv := t.NumericValue()
            nv *= scale
            return maybe.Some(nv), true
        } else {
            return maybe.Nothing[float64](), false
        }
    }
}

var acceptPercentage = numericAcceptor(token.TypePercentage, 0.01)

func acceptNone(t token.Token) (maybe.M[float64], bool) {
    ok := t.Is(token.TypeIdent) && (strings.EqualFold(t.StringValue(), "none"))
    return maybe.Nothing[float64](), ok
}

func parseRGBFromFunction(f token.Token, args []token.Token) (Color, error) {
    zero := Color{}
    var r,g,b,a maybe.M[float64]
    acceptNumber := numericAcceptor(token.TypeNumber, 1.0 / 255.0)
    acceptRawNumber := numericAcceptor(token.TypeNumber, 1.0)

    modern := func(acceptor func(t token.Token) (maybe.M[float64], bool)) bool {
        var ok bool
        if !strings.EqualFold(f.StringValue(), "rgb") { return false }
        rest := args

        t, rest := step(rest)
        r, ok = acceptEither(t, acceptor, acceptNone)
        if !ok { return false }

        t, rest = step(rest)
        g, ok = acceptEither(t, acceptor, acceptNone)
        if !ok { return false }

        t, rest = step(rest)
        b, ok = acceptEither(t, acceptor, acceptNone)
        if !ok { return false }

        t, rest = step(rest)
        if t.Is(token.TypeEOF) {
            a = maybe.Some(1.0)
            return true
        }
        if !(t.Is(token.TypeDelim) && (t.Delim() == '/')) { return false }

        t, rest = step(rest)
        a, ok = acceptEither(t, acceptPercentage, acceptRawNumber, acceptNone)
        if !ok { return false }

        t, rest = step(rest)
        return t.Is(token.TypeEOF)
    }

    // RGB( [<percentage> | none]{3} [ / [<alpha-value> | none] ]? )
    if modern(acceptPercentage) { return RGB(r, g, b, a), nil }
    // RGB( [<number> | none]{3} [ / [<alpha-value> | none] ]? )
    if modern(acceptNumber) { return RGB(r, g, b, a), nil }

    legacy := func(acceptor func(t token.Token) (maybe.M[float64], bool)) bool {
        var ok bool
        rest := args

        t, rest := step(rest)
        r, ok = acceptEither(t, acceptor, acceptNone)
        if !ok { return false }

        t, rest = step(rest)
        if !t.Is(token.TypeComma) { return false }

        t, rest = step(rest)
        g, ok = acceptEither(t, acceptor, acceptNone)
        if !ok { return false }

        t, rest = step(rest)
        if !t.Is(token.TypeComma) { return false }

        t, rest = step(rest)
        b, ok = acceptEither(t, acceptor, acceptNone)
        if !ok { return false }

        t, rest = step(rest)
        if t.Is(token.TypeEOF) {
            a = maybe.Some(1.0)
            return true
        }
        if !t.Is(token.TypeComma) { return false }

        t, rest = step(rest)
        a, ok = acceptEither(t, acceptPercentage, acceptRawNumber)
        if !ok { return false }

        t, rest = step(rest)
        return t.Is(token.TypeEOF)
    }

    // rgba?( <percentage>#{3} , <alpha-value>? )
    if legacy(acceptPercentage) { return RGB(r, g, b, a), nil }
    // rgba?( <number>#{3} , <alpha-value>? )
    if legacy(acceptNumber) { return RGB(r, g, b, a), nil }

    return zero, errSyntax{ErrInvalidArguments, f.Position()}
}
