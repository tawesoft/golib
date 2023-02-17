package main

import (
    "flag"
    "path"

    "github.com/tawesoft/golib/v2/internal/unicode/maketables/cldr41"
)

func main() {
    var (
        dataDir, destDir string
    )
    flag.StringVar(&dataDir, "data", "../DATA", "unicode data directory")
    flag.StringVar(&destDir, "dest", "../../../", "destination relative to module root")
    flag.Parse()

    cldr41.MakeNumberingSystemRules(dataDir, path.Join(destDir, "text/number/algorithmic/rules-cldr-41.0.txt"))
}
