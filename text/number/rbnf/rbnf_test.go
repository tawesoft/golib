package rbnf

import (
    "strings"
    "testing"

    "github.com/tawesoft/golib/v2/must"
)

func TestParseRuleDescriptor(t *testing.T) {
    type row struct {
        input string
        r rule
    }
    rows := []row{
        {input: "100000",      r: rule{Type: ruleTypeBaseValue,         Base: 100_000, Divisor: 0}},
        {input: "100000/1000", r: rule{Type: ruleTypeBaseValueAndRadix, Base: 100_000, Divisor: 1000}},
    }
    for _, test := range rows {
        got := parseRuleDescriptor(test.input)
        if (got != test.r) {
            t.Errorf("pareRuleDescriptor(%q): got %+v but expected %+v",
                test.input, got, test.r)
        }
    }
}

func TestParseBracketed(t *testing.T) {
    type row struct {
        input string
        inner string
        start rune
        end rune
    }
    rows := []row{
        {"→→", "", '→', '→'},
        {"←←", "", '←', '←'},
        {"→a→", "a", '→', '→'},
        {"←b←", "b", '←', '←'},
        {"[abc]", "abc", '[', ']'},
    }
    parse := func(r row) {
        var seen int
        var dest string
        var sb strings.Builder
        p := newStringParser(r.input)
        next := []rune(r.input)[1]
        parseBracketed(p, &sb, r.start, next, r.end, &seen, 1, &dest)
        if seen != 1 {
            t.Errorf("parseBracketed(%q): expected seen=1 but seen=%d", r.input, seen)
        }
        if dest != r.inner {
            t.Errorf("parseBracketed(%q): expected inner=%q but got %q", r.input, r.inner, dest)
        }
    }
    for _, test := range rows {
        parse(test)
    }
}

func TestParseRuleBody(t *testing.T) {
    type row struct {
        input string
        r rule
    }
    rows := []row{
        {input: "→→", r: rule{}},
        {input: "→a→", r: rule{Subs: [3]string{"a"}}},
        {input: "→a→→b→→c→", r: rule{Subs: [3]string{"a", "b", "c"}}},
        {input: "←← point →→", r: rule{Literal: " point "}},
    }
    for _, test := range rows {
        got := parseRuleBody(rule{}, test.input)
        if (got != test.r) {
            t.Errorf("pareRuleDescriptor(%q): got %+v but expected %+v",
                test.input, got, test.r)
        }
    }
}

func TestNew(t *testing.T) {
    g := must.Result(New(nil, `
        %spellout-cardinal-verbose:
            -x: minus →→;
            x.x: ←← point →→;
            Inf: infinite;
            NaN: not a number;
            0: =%spellout-numbering=;
            100: ←← hundred[→%%and→];
            1000: ←← thousand[→%%and→];
            100000/1000: ←← thousand[→%%commas→];
            1000000: ←← million[→%%commas→];
            1000000000: ←← billion[→%%commas→];
            1000000000000: ←← trillion[→%%commas→];
            1000000000000000: ←← quadrillion[→%%commas→];
            1000000000000000000: =#,##0=;
    `))
    t.Errorf("%+v", g)
}
