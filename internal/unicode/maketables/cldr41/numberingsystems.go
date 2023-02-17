package cldr41

import (
    "fmt"
    "io"
    "os"
    "path"
    "strings"

    "github.com/tawesoft/golib/v2/internal/unicode/ldml"
    "github.com/tawesoft/golib/v2/must"
)

func MakeNumberingSystemRules(srcdir string, dest string) {
    wr := must.Result(os.Create(dest))
    defer wr.Close()
    makeNumberingSystemRules(srcdir, wr)
}

func makeNumberingSystemRules(basedir string, wr io.Writer) {
    r := must.Result(os.Open(path.Join(basedir, "cldr-41.0/common/rbnf/root.xml")))
    defer r.Close()
    gdoc := must.Result(ldml.Parse(must.Result(io.ReadAll(r))))
    must.Equal(gdoc.Type, ldml.DocTypeLdml)
    doc := gdoc.Ldml

    for _, group := range doc.RbnfRulesetGroupings {
        if group.Type != "NumberingSystemRules" { continue }
        for _, set := range group.Rulesets {
            if strings.EqualFold(set.Access, "private") {
                must.Result(fmt.Fprintf(wr, "%%%%%s:\n", set.Type))
            } else {
                must.Result(fmt.Fprintf(wr, "%%%s:\n", set.Type))
            }
            for _, rule := range set.Rules {
                must.Result(fmt.Fprintf(wr, "\t%s\n", rule.IcuStyle()))
            }
            must.Result(fmt.Fprintf(wr, "\n"))
        }
    }
}
