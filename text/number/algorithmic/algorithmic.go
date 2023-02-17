// Package algorithmic implements parsing and formatting of non-decimal
// number systems, such as roman numerals and traditional tamil numbers.
package algorithmic

import (
    _ "embed"

    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/text/number/rbnf"
)

//go:embed rules-cldr-41.0.txt
var rules string

var group = must.Result(rbnf.New(nil, rules))

func Format(name string, number int64) (string, error) {
    return group.FormatInteger(name, number)
}
