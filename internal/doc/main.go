// Command doc builds static documentation using the new Go 1.19 documentation
// features.
package main

import (
    "fmt"
    "go/ast"
    "go/doc"
    "go/doc/comment"
    "go/parser"
    "go/token"
    "math"
    "os"
    "path/filepath"
    "strings"

    cli "github.com/jawher/mow.cli"
    "github.com/tawesoft/golib/v2/fun/slice"
    "github.com/tawesoft/golib/v2/ks"
)

// TODO generate HTML like https://cs.opensource.google/go/x/tools/+/master:godoc/linkify.go

func genTarget(t string) error {
    matches, err := filepath.Glob(filepath.FromSlash(t + "/*.go"))
    if err != nil {
        return fmt.Errorf("glob error: %w", err)
    }
    if len(matches) == 0 {
        return fmt.Errorf("no files")
    }

    fset := token.NewFileSet()
    astFiles := make([]*ast.File, 0)
    sources := make(map[string][]byte)

    for _, f := range matches {
        stat, err := os.Stat(f)
        if err != nil {
            return fmt.Errorf("error statting file %q: %w", f, err)
        }

        size := stat.Size()
        if size > math.MaxInt32 {
            return fmt.Errorf("file too large: %q", f)
        }

        src, err := os.ReadFile(f)
        sources[f] = src
        if err != nil {
            return fmt.Errorf("error reading file %q: %w", f, err)
        }

        pf, err := parser.ParseFile(fset, f, src, parser.ParseComments)
        if err != nil {
            return fmt.Errorf("error parsing file %q: %w", f, err)
        }

        astFiles = append(astFiles, pf)
    }

    pkg, err := doc.NewFromFiles(fset, astFiles, "tawesoft.co.uk/go/"+t)
    if err != nil {
        return fmt.Errorf("")
    }

    // fmt.Printf(">>> %s: %+v\n", target, pkg)
    p := pkg.Parser()
    pr := pkg.Printer()
        pr.DocLinkBaseURL = "https://pkg.go.dev"


    pr.DocLinkURL = func(link *comment.DocLink) string {
        fmt.Printf("> %+v\n", link)
        if link.ImportPath == "builtin" && link.Name == "" {
            link.Name = string(link.Text[0].(comment.Plain))
            link.Name = strings.TrimPrefix(link.Name, "builtin.")
            fmt.Printf(">>> %+v\n", link)
        }
        return link.DefaultURL(pr.DocLinkBaseURL)
    }

    // oldLookupPackage := p.LookupPackage
    p.LookupPackage = func(name string) (importPath string, ok bool) {
        fmt.Println(name)

        if (name == "builtin") { // e.g. builtin.IntegerType
            return name, true
        } else if strings.HasPrefix(name, "builtin.") { // e.g. builtin.map
            return "builtin", true
        } else if name == "comparable" { // or any lowercase builtin identifier...
            // name is actual builtin
            return "builtin", true
        }
        return "", false
    }

    /*
    pr.DocLinkURL = func(link *comment.DocLink) string {
        panic(fmt.Sprintf("got %+v", link.DefaultURL("https://example.net/")))
    }
     */

    doc := p.Parse(pkg.Doc)

    fmt.Printf("%s\n", pr.Markdown(doc))

    return nil
}

func gen(targets []string) error {
    if t, err := slice.CheckedWalk(genTarget, targets); err != nil {
        return fmt.Errorf("error generating target %s", t)
    }
    return nil
}

func main() {

    app := cli.App("doc", "generate documentation")

    app.Command("gen", "generate docs for named folders", func(cmd *cli.Cmd) {
        cmd.Spec = "TARGETS..."
        targets := cmd.StringsArg("TARGETS", nil, "target names e.g. foo foo/bar")

        cmd.Action = func() {
            err, perr := ks.Catch[error](func() error { return gen(*targets) })
            if perr != nil { err = perr }
            if err != nil {
                fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
                cli.Exit(1)
            }
        }
    })

    app.Run(os.Args)
}
