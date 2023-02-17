package rbnf

import (
    "fmt"
    "math"
    "strings"
    "unicode"
    "unicode/utf8"

    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/text/number/rbnf/internal/body"
    "github.com/tawesoft/golib/v2/text/runeio"
)

func encodeTokenType(ty body.Type, sty body.SubstType) uint8 {
    return (0b00001111 & uint8(ty)) +
           (0b01110000 & (uint8(sty) << 4))
}

func decodeTokenType(x uint8) (body.Type, body.SubstType) {
    return body.Type(0b00001111 & x), body.SubstType((0b01110000 & x) >> 4)
}

func (g *Group) parse(s string) error {
    g.rulesetNames = make(map[string]int)
    if err := g.parse1(s); err != nil {
        return err
    }
    if err := g.parse2(s); err != nil {
        return err
    }
    /*if err := g.check(); err != nil {
        return err
    }*/
    return nil
}

// parse1 does a first pass through the string, capturing ruleset names
// for efficient lookup later.
func (g *Group) parse1(s string) (retError error) {
    var sb strings.Builder
    p := newStringParser(s)

    defer func() {
        if r := recover(); r != nil {
            retError = fmt.Errorf("parse error at offset %+v: %w", p.r.Offset(), r.(error))
        }
    }()

    for {
        p.ConsumeWhitespace()
        if (runeio.Must(p.r.Peek()) == '%') {
            sb.Reset()
            p.ConsumeRulesetName(&sb)
            g.addRuleset(sb.String())
            p.ConsumeSymbol(':')
            p.ConsumeWhitespace()
        }

        p.ConsumeUntilSymbol(nil, ':')
        p.ConsumeSymbol(':')
        p.ConsumeWhitespace()
        p.ConsumeOptionalSymbol('\'')
        p.ConsumeUntilSymbol(nil, ';')
        p.ConsumeSymbol(';')
        p.ConsumeWhitespace()
        if p.AtEOF() { break }
    }

    return nil
}

// parse2 does a second pass through the string, encoding rules.
func (g *Group) parse2(s string) (retError error) {
    var sb strings.Builder
    var stringData strings.Builder
    p := newStringParser(s)

    defer func() {
        if r := recover(); r != nil {
            retError = fmt.Errorf("parse error at offset %+v: %w", p.r.Offset(), r.(error))
            panic(retError)
        }
    }()

    var rs *ruleset
    var ruleDesc *desc
    for {
        p.ConsumeWhitespace()
        if (runeio.Must(p.r.Peek()) == '%') {
            sb.Reset()
            p.ConsumeRulesetName(&sb)
            rs = g.initRuleset(sb.String())
            p.ConsumeSymbol(':')
            p.ConsumeWhitespace()
        }
        if rs == nil { panic(fmt.Errorf("expected ruleset name")) }

        sb.Reset()
        p.ConsumeUntilSymbol(&sb, ':')
        {
            ruleDescriptor := sb.String()
            ruleDesc = g.addRuleDescriptor(rs, ruleDescriptor)
        }

        p.ConsumeSymbol(':')
        p.ConsumeWhitespace()
        p.ConsumeOptionalSymbol('\'')

        sb.Reset()
        p.ConsumeUntilSymbol(&sb, ';')
        ruleBody := sb.String()
        tokenizer := body.NewTokenizer(ruleBody)
        for {
            t := must.Result(tokenizer())
            if t.Type == body.TypeEOF { break }
            if ruleDesc.NumTokens >= math.MaxUint8 {
                panic(fmt.Errorf("too many tokens in rule body"))
            }
            ruleDesc.NumTokens++

            st := t.SimpleSubstType(ruleBody)
            length := 0 // or idx
            left := 0

            if t.Content.Len() > 0 {
                if st == body.SubstTypeRulesetName {
                    name := t.Content.Of(ruleBody)
                    if idx, ok := g.getRulesetIndex(name); ok {
                        length = idx
                    } else {
                        panic(fmt.Errorf("ruleset not found: %s", name))
                    }
                } else {
                    start := stringData.Len()
                    stringData.WriteString(t.Content.Of(ruleBody))
                    left = start
                    length = t.Content[1] - t.Content[0]
                }
            }

            if (t.Type > 0b00001111) || (st > 0b111) || (length < 0) ||
                (length > math.MaxUint8) || (left > math.MaxUint16) {
                panic(fmt.Errorf("rule body too large"))
            }

            g.bodies = append(g.bodies, token{
                Type: encodeTokenType(t.Type, st),
                Len:  uint8(length),
                Left: uint16(left),
            })
        }

        p.ConsumeSymbol(';')
        p.ConsumeWhitespace()
        if p.AtEOF() { break }
    }

    g.stringData = stringData.String()
    return
}

// parser parses a stream of Unicode runes
type parser struct {
    r *runeio.Reader
}

// newStringParser returns a parser with a limited peek & pushback buffer
func newStringParser(s string) parser {
    var buf [utf8.UTFMax * 2]byte
    rdr := runeio.NewReader(strings.NewReader(s))
    rdr.Buffer(buf[:], utf8.UTFMax * 2)
    return parser{rdr}
}

func (p parser) AtEOF() bool {
    return runeio.Must(p.r.Peek()) == runeio.RuneEOF
}

func (p parser) ConsumeWhitespace() {
    for unicode.IsSpace(runeio.Must(p.r.Peek())) {
        runeio.Must(p.r.Next())
    }
}

func (p parser) ConsumeSymbol(symbol rune) {
    c := runeio.Must(p.r.Next())
    if c != symbol {
        panic(fmt.Errorf("expected %c", symbol))
    }
}

func (p parser) ConsumeOptionalSymbol(symbol rune) {
    c := runeio.Must(p.r.Peek())
    if c == symbol {
        runeio.Must(p.r.Next())
    }
}

func (p parser) ConsumeRulesetName(sb *strings.Builder) {
    // must start with at least one %
    if sb != nil { sb.WriteRune('%') }
    p.ConsumeSymbol('%')
    p.ConsumeUntilSymbol(sb, ':')
}

// ConsumeUntilSymbol consumes until reaching the given symbol, raising a panic
// if EOF is reached before the symbol is found.
func (p parser) ConsumeUntilSymbol(sb *strings.Builder, symbol rune) {
    for {
        c := runeio.Must(p.r.Next())
        switch c {
            case symbol:
                p.r.Push(c)
                return
            case runeio.RuneEOF:
                panic(fmt.Errorf("expected %c", symbol))
            default:
                if sb != nil { sb.WriteRune(c) }
        }
    }
}

// ConsumeUntilSymbol2 consumes until reaching the given symbol pair, raising a
// panic if EOF is reached before the symbol pair is found.
func (p parser) ConsumeUntilSymbol2(sb *strings.Builder, a, b rune) {
    var buf [2]rune
    for {
        must.Result(p.r.PeekN(buf[:], 2))
        current, next := buf[0], buf[1]
        if (current == a) && (next == b) {
            return
        }
        if current == runeio.RuneEOF {
            panic(fmt.Errorf("expected %c%c", a, b))
        }
        if (next != runeio.RuneEOF) && (sb != nil) { sb.WriteRune(next) }
        must.Check(p.r.Skip(2))
    }
}

// ConsumeUntil2 consumes and appends to a string builder until f returns true
// for the following two runes. f should return true for EOF if allowed here.
func (p parser) ConsumeUntil2(sb *strings.Builder, f func(a, b rune) bool) {
    var buf [2]rune
    for {
        must.Result(p.r.PeekN(buf[:], 2))
        current, next := buf[0], buf[1]
        if f(current, next) {
            return
        }
        if current == runeio.RuneEOF {
            panic(fmt.Errorf("unexpected EOF"))
        }
        if sb != nil { sb.WriteRune(current) }
        must.Check(p.r.Skip(1))
    }
}
