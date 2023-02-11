package rbnf

import (
    "fmt"
    "strconv"
    "strings"
    "unicode"
    "unicode/utf8"

    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/operator"
    "github.com/tawesoft/golib/v2/text/runeio"
    "golang.org/x/text/runes"
)

// decimalFold is a transformer for decimal integers that filters out
// spaces, period, and commas.
var decimalFold = runes.Remove(runes.Predicate(func(r rune) bool {
        if (r >= '0') && (r <= '9') { return false } // quick case
        if (r == '.') || (r == ',') { return true }
        return unicode.IsSpace(r)
}))

// isAsciiDigitString returns true iff the only characters in x are the ASCII
// digits '0' to '9', and x contains at least one digit.
func isAsciiDigitString(x string) bool {
    if len(x) == 0 { return false }
    for _, c := range x {
        if (c < '0') || (c > '9') { return false }
    }
    return true
}

// parseRuleDescriptor parses a rule descriptor in the following format, or
// panics:
//
// * bv and rad are the names of tokens formatted as decimal numbers expressed
//   using ASCII digits with spaces, period, and commas ignored.
// * bv specifies the rule's base value. The rule's divisor is the highest power
//   of 10 less than or equal to the base value.
// * bv/rad: The rule's divisor is the highest power of rad less than or equal to
//   the base value.
// * -x: The rule is a negative-number rule.
// * x.x: The rule is an improper fraction rule.
// * 0.x: The rule is a proper fraction rule.
// * x.0: The rule is a default rule.
// * Inf: The rule for infinity.
// * NaN: The rule for an IEEE 754 NaN (not a number).
func parseRuleDescriptor(s string) rule {
    const ferr = "invalid rule descriptor syntax %q"

    switch s {
        case "-x":  { return rule{Type: ruleTypeNegativeNumber} }
        case "x.x": { return rule{Type: ruleTypeImproperFraction} }
        case "0.x": { return rule{Type: ruleTypeProperFraction} }
        case "x.0": { return rule{Type: ruleTypeDefault} }
        case "Inf": { return rule{Type: ruleTypeInfinity} }
        case "NaN": { return rule{Type: ruleTypeNaN} }
    }

    idx := strings.IndexRune(s, '/')
    if idx > 0 {
        left := decimalFold.String(s[:idx])
        right := decimalFold.String(s[idx+1:])

        if (!isAsciiDigitString(left)) || (!isAsciiDigitString(right)) {
            panic(fmt.Errorf(ferr, s))
        }

        base,  baseErr  := strconv.ParseInt(left, 10, 64)
        radix, radixErr := strconv.ParseInt(right, 10, 64)
        if (baseErr != nil) || (radixErr != nil) {
            panic(fmt.Errorf("rule descriptor range error while parsing %q", s))
        }

        return rule{
            Base:    base,
            Divisor: radix,
            Type:    ruleTypeBaseValueAndRadix,
        }
    } else if idx < 0 {
        s = decimalFold.String(s)

        if !isAsciiDigitString(s) {
            panic(fmt.Errorf(ferr, s))
        }

        base, baseErr  := strconv.ParseInt(s, 10, 64)
        if baseErr != nil {
            panic(fmt.Errorf("rule descriptor range error while parsing %q", s))
        }

        return rule{
            Base: base,
            Type: ruleTypeBaseValue,
        }
    } else {
        panic(fmt.Errorf(ferr, s))
    }
}

// parseBracketed reads and stores a bracketed term such as "→foo→", with
// a limit on the number of occurrences, and stores the result and updates the
// number of seen occurrences.
func parseBracketed(p parser, sb *strings.Builder, current, next, terminator rune, seen *int, max int, dest *string) {
    if (current == next) && (next == '→') {} // TODO optimise for this case
    if (current == next) && (next == '←') {} // TODO optimise for this case
    if *seen > max {
        panic(fmt.Errorf("too many %c%c parts in rule body (seen %d/%d)", current, terminator, *seen, max))
    }
    must.Check(p.r.Skip(1))
    sb.Reset()
    p.ConsumeUntilSymbol(sb, terminator)
    *dest = sb.String()
    must.Check(p.r.Skip(1))
    (*seen)++
}

func parseRuleBody(r rule, s string) rule {
    // zero, one, or two substitution tokens
    // zero or one text in brackets
    // zero or one literal
    var seenSubs, seenOptional, seenLiteral int
    var sb strings.Builder

    p := newStringParser(s)

    defer func() {
        if r := recover(); r != nil {
            panic(fmt.Errorf("rule body parse error while parsing: %q: %w", s, r.(error)))
        }
    }()

    outer:
    for {
        var buf [2]rune
        must.Result(p.r.PeekN(buf[:], 2))
        current, next := buf[0], buf[1]

        switch current {
            case runeio.RuneEOF: break outer
            case '→': fallthrough
            case '←': fallthrough
            case '=': parseBracketed(p, &sb, current, next, current, &seenSubs, 2, &r.Subs[seenSubs])
            case '[': parseBracketed(p, &sb, current, next, ']', &seenOptional, 1, &r.Optional)
            default: {
                if (current == '$') && (next == '(') {
                    if seenOptional > 0 {
                        panic("unexpected optional in rule body")
                    }
                    must.Check(p.r.Skip(2))
                    sb.Reset()
                    p.ConsumeUntilSymbol2(&sb, ')', '$')
                    r.Literal = sb.String()
                    must.Check(p.r.Skip(2))
                } else {
                    if seenLiteral > 0 {
                        panic("unexpected literal in rule body")
                    }
                    sb.Reset()
                    p.ConsumeUntil2(&sb, func(a, b rune) bool {
                        return operator.In(a, '→', '←', '=', '[') ||
                            ((a == '$') && (b == '(')) ||
                            (a == runeio.RuneEOF)
                    })
                    r.Literal = sb.String()
                }
            }
        }
    }

    return r
}

func (g *Group) parseGroups(s string) (retError error) {
    var sb strings.Builder
    p := newStringParser(s)

    defer func() {
        if r := recover(); r != nil {
            retError = fmt.Errorf("parse error at offset %+v: %w", p.r.Offset(), r.(error))
        }
    }()

    var rsName string
    for {
        p.ConsumeWhitespace()
        if (runeio.Must(p.r.Peek()) == '%') || (rsName == "") {
            sb.Reset()
            p.ConsumeRulesetName(&sb)
            rsName = sb.String()

            if _, exists := g.rulesets[rsName]; exists {
                panic(fmt.Errorf("duplicate ruleset name %q", rsName))
            }
            g.rulesets[rsName] = ruleset{}

            p.ConsumeSymbol(':')
            p.ConsumeWhitespace()
        }

        sb.Reset()
        p.ConsumeUntilSymbol(&sb, ':')
        descriptor := sb.String()
        r := parseRuleDescriptor(descriptor)

        p.ConsumeSymbol(':')
        p.ConsumeWhitespace()
        p.ConsumeOptionalSymbol('\'')

        sb.Reset()
        p.ConsumeUntilSymbol(&sb, ';')
        body := sb.String()
        r = parseRuleBody(r, body)

        p.ConsumeSymbol(';')
        p.ConsumeWhitespace()

        g.rulesets[rsName] = append(g.rulesets[rsName], r)

        if p.AtEOF() { break }
    }

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
    sb.WriteRune('%')
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
                sb.WriteRune(c)
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
        if next != runeio.RuneEOF { sb.WriteRune(next) }
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
        sb.WriteRune(current)
        must.Check(p.r.Skip(1))
    }
}
