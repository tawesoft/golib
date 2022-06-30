// Command doc builds static documentation using the new Go 1.19 documentation
// features.
package main

import (
    "fmt"
    "go/ast"
    "go/doc"
    "go/parser"
    "go/token"
    "math"
    "os"
    "path/filepath"

    cli "github.com/jawher/mow.cli"
)

// TODO generate HTML like https://cs.opensource.google/go/x/tools/+/master:godoc/linkify.go

func warnf(format string, args ... any) {
    fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func genTarget(target string, files []string) error {

    fset := token.NewFileSet()
    astFiles := make([]*ast.File, 0)
    sources := make(map[string][]byte)

    for _, f := range files {
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

    pkg, err := doc.NewFromFiles(fset, astFiles, "tawesoft.co.uk/go/"+target)
    if err != nil {
        return fmt.Errorf("")
    }

    // fmt.Printf(">>> %s: %+v\n", target, pkg)
    p := pkg.Parser()
    pr := pkg.Printer()
    doc := p.Parse(pkg.Doc)
    os.Stdout.Write(pr.HTML(doc))

    if len(pkg.Examples) > 0 {
        // file := fset.File(pkg.Examples[0].Play.Pos()).Name()
        //fmt.Println(string(sources[file]))
    }

    return nil
}


func walkTargets(targets []string, f func (t string) error) error {
    for _, t := range targets {
        err := f(t)
        if err != nil {
            return fmt.Errorf("error processing target %s: %w", t, err)
        }
    }
    return nil
}

func gen(targets []string) error {

    return walkTargets(targets, func (t string) error {

        matches, err := filepath.Glob(filepath.FromSlash(t + "/*.go"))
        if err != nil {
            return fmt.Errorf("error globbing target %s: %w", t, err)
        }

        if len(matches) == 0 {
            warnf("warning: no files for target %s", t)
        }

        err = genTarget(t, matches)
        if err != nil {
            return fmt.Errorf("error generating target %s: %w", t, err)
        }

        return nil
    })
}


func main() {

    app := cli.App("doc", "generate documentation")

    app.Command("gen", "generate docs for named folders", func(cmd *cli.Cmd) {
        cmd.Spec = "TARGETS..."
        targets := cmd.StringsArg("TARGETS", nil, "target names e.g. foo foo/bar")

        cmd.Action = func() {
            err := gen(*targets)
            if err != nil {
                fmt.Fprintf(os.Stderr, "fatal error: %w\n", err)
                cli.Exit(1)
            }
        }
    })

    app.Run(os.Args)

}
