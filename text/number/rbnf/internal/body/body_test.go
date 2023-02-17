package body_test

import (
    "fmt"
    "io"
    "os"
    "path"
    "runtime"
    "testing"

    "github.com/tawesoft/golib/v2/internal/unicode/ldml"
    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/text/number/rbnf/internal/body"
)

func TestTokenizer_Next(t *testing.T) {
    type row struct {
        input string
        ok bool
    }
    rows := []row{
        {input: "",             ok: true},
        {input: "==",           ok: false}, // empty not permitted between ==
        {input: "=%literal=",   ok: true},
        {input: "=#,##0=",      ok: true},
        {input: "=invalid=",    ok: false},
        {input: "←←",           ok: true},  // may be empty
        {input: "←←←←←←",       ok: true},
        {input: "←←←←←←←←",     ok: false}, // too many
        {input: "[",            ok: false}, // open bracket
        {input: "]",            ok: false}, // unexpected bracket
        {input: "[]",           ok: true},
        {input: "X←←[X←←]←←X",  ok: true},  // mixing literals
        {input: "[][]",         ok: false}, // too many brackets
        {input: "←←[←←]←←",     ok: true},
        {input: "←←[←←←←]←←",   ok: false}, // too many

        {input: "$(ordinal,one{A}other{B})$",  ok: true},
        {input: "$(cardinal,one{A}other{B})$", ok: true},
        {input: "$(invalid,one{A}other{B})$",  ok: false},
        {input: "$()$", ok: false},
        {input: "$()$", ok: false},

        // sample real rules from CLDR
        {input: "=#,##0=$(ordinal,one{st}two{nd}few{rd}other{th})$", ok: true},
        {input: "←%spellout-numbering← quadrillion→%%th→;", ok: true},
        {input: "twen→%%tieth→;", ok: true},
    }
    for _, test := range rows {
        try := func() error {
            tokenizer := body.NewTokenizer(test.input)
            for {
                tok, err := tokenizer()
                if err != nil { return err }
                if tok.Type == body.TypeEOF {
                    return nil
                }
            }
        }
        err := try()
        if (err != nil) && test.ok {
            t.Errorf("tokenize %q: expected nil error but got: %v", test.input, err)
        } else if (err == nil) && (!test.ok) {
            t.Errorf("tokenize %q: expected error but got unexpected success", test.input)
        }
    }
}

func cldrpath(v string, p string) string {
    _, b, _, _ := runtime.Caller(0)
    return path.Join(
        path.Dir(b),
        "../../../../../internal/unicode/DATA/cldr-"+v+"/common/",
        p,
    )
}

func TestTokenizer_NextFromCLDR(t *testing.T) {
    type row struct {
        source string
        input string
    }
    cldrvers := []string{"41.0"}
    var passes, total int

    for _, cldrver := range cldrvers {
        rbnfdir := cldrpath(cldrver, "rbnf")
        fp, err := os.Open(rbnfdir)
        if err != nil {
            t.Skipf("skip: could not read %q", rbnfdir)
        }
        sources := must.Result(fp.Readdirnames(0))
        t.Logf("sources: %v", sources)

        for _, source := range sources {
            source = cldrpath(cldrver, "rbnf/"+source)
            rows := []row{}

            fp, err := os.Open(source)
            if err != nil {
                t.Logf("skip %q: error: %v", source, err)
                continue
            }

            gdoc := must.Result(ldml.Parse(must.Result(io.ReadAll(fp))))
            if gdoc.Type != ldml.DocTypeLdml {
                t.Logf("skip %q: not a ldml file", source)
                continue
            }
            doc := gdoc.Ldml

            for _, group := range doc.RbnfRulesetGroupings {
                for _, set := range group.Rulesets {
                    if set.Type == "lenient-parse" { continue }
                    for _, rule := range set.Rules {
                        rows = append(rows, row{
                            source: fmt.Sprintf("%s-%s-%s-%s:%s/%s:%s",
                                doc.Language, doc.Script, doc.Region, doc.Variant,
                                set.Type, rule.Value, rule.Radix),
                                input: rule.Content,
                        })
                    }
                }
            }

            for _, test := range rows {
                try := func() error {
                    tokenizer := body.NewTokenizer(test.input)
                    for {
                        tok, err := tokenizer()
                        if err != nil { return err }
                        if tok.Type == body.TypeEOF {
                            return nil
                        }
                    }
                }
                err := try()
                if (err != nil) {
                    t.Errorf("tokenize %q error (from %q): %v", test.input, test.source, err)
                } else {
                    passes++
                }
                total++
            }
        }
    }

    t.Logf("Pass: %d/%d", passes, total)
}
