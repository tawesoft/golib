package ldml_test

import (
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/tawesoft/golib/v2/internal/unicode/ldml"
    "github.com/tawesoft/golib/v2/must"
)

func ExampleParseRbnf() {
    r := must.Result(os.Open("testdata/en.xml"))
    defer r.Close()
    gdoc := must.Result(ldml.Parse(must.Result(io.ReadAll(r))))
    must.Equal(gdoc.Type, ldml.DocTypeLdml)
    doc := gdoc.Ldml

    fmt.Printf("Doc %q %q %q %q:\n", doc.Language, doc.Script, doc.Region, doc.Variant)

    for _, group := range doc.RbnfRulesetGroupings {
        fmt.Printf("....Group: %q\n", group.Type)
        for _, set := range group.Rulesets {
            if set.Type != "spellout-cardinal-verbose" {
                fmt.Printf("........%%%%%s: (omitted)\n", set.Type)
                continue
            }
            if strings.EqualFold(set.Access, "private") {
                fmt.Printf("........%%%%%s:\n", set.Type)
            } else {
                fmt.Printf("........%%%s:\n", set.Type)
            }
            for _, rule := range set.Rules {
                fmt.Printf("............%s\n", rule.IcuStyle())
            }
        }
    }

    // Output:
    // Doc "en" "Zzzz" "ZZ" "example":
    // ....Group: "SpelloutRules"
    // ........%%and: (omitted)
    // ........%%commas: (omitted)
    // ........%spellout-cardinal-verbose:
    // ............-x: minus →→;
    // ............x.x: ←← point →→;
    // ............Inf: infinite;
    // ............NaN: not a number;
    // ............0: =%spellout-numbering=;
    // ............100: ←← hundred[→%%and→];
    // ............1000: ←← thousand[→%%and→];
    // ............100000/1000: ←← thousand[→%%commas→];
    // ............1000000: ←← million[→%%commas→];
    // ............1000000000: ←← billion[→%%commas→];
    // ............1000000000000: ←← trillion[→%%commas→];
    // ............1000000000000000: ←← quadrillion[→%%commas→];
    // ............1000000000000000000: =#,##0=;
    // ....Group: "OrdinalRules"
    // ........%%digits-ordinal: (omitted)
}

func ExampleParsePlurals() {
    r := must.Result(os.Open("testdata/plurals.xml"))
    defer r.Close()
    gdoc := must.Result(ldml.Parse(must.Result(io.ReadAll(r))))
    must.Equal(gdoc.Type, ldml.DocTypeSupplemental)
    doc := gdoc.Supplemental


    fmt.Printf("Plurals type=%q\n\n", doc.Plurals.Type)
    for _, rules := range doc.Plurals.Rules {
        fmt.Printf("[%s]\n", rules.Locales)
        for _, rule := range rules.Rules {
            fmt.Printf("> %s:\n>     %s\n", rule.Count, rule.Content)
        }
        fmt.Printf("\n")
    }

    // Output:
    // Plurals type="cardinal"
    //
    // [bm bo dz hnj id ig ii in ja jbo jv jw kde kea km ko lkt lo ms my nqo osa root sah ses sg su th to tpi vi wo yo yue zh]
    // > other:
    // >      @integer 0~15, 100, 1000, 10000, 100000, 1000000, … @decimal 0.0~1.5, 10.0, 100.0, 1000.0, 10000.0, 100000.0, 1000000.0, …
    //
    // [am as bn doi fa gu hi kn pcm zu]
    // > one:
    // >     i = 0 or n = 1 @integer 0, 1 @decimal 0.0~1.0, 0.00~0.04
    // > other:
    // >      @integer 2~17, 100, 1000, 10000, 100000, 1000000, … @decimal 1.1~2.6, 10.0, 100.0, 1000.0, 10000.0, 100000.0, 1000000.0, …
}
