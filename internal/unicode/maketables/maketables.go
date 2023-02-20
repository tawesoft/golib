package main

import (
    "flag"
    "fmt"
    "path"
    "time"

    "github.com/tawesoft/golib/v2/internal/unicode/maketables/cldr41"
)

func main() {
    var (
        dataDir, destDir string
    )
    flag.StringVar(&dataDir, "data", "../DATA", "unicode data directory")
    flag.StringVar(&destDir, "dest", "../../../", "destination relative to module root")
    flag.Parse()

    dest := func(x string) string {
        return path.Join(destDir, x)
    }

    timeit := func(desc string, f func()) {
        start := time.Now()
        f()
        elapsed := time.Since(start)
        fmt.Printf("%s: %s\n", desc, elapsed)
    }

    timeit("MakeNumberingSystemRules", func() {
        cldr41.MakeNumberingSystemRules(dataDir, dest("text/number/algorithmic/rules-cldr-41.0.txt"))
    })
    timeit("MakeNumberSymbols", func() {
        cldr41.MakeNumberSymbols(dataDir, destDir)
    })
}
