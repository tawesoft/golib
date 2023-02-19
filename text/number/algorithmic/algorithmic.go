// Package algorithmic implements parsing and formatting of common non-decimal
// number systems, such as roman numerals and traditional tamil numbers.
//
// These correspond to the rules in CLDR 41.0 at "rbnf/root.xml" (and does not
// include every algorithmic numbering system mentioned in
// "numberingSystems.xml").
package algorithmic

import (
    _ "embed"
    "sort"
    "strings"

    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/text/number/rbnf"
)

//go:embed rules-cldr-41.0.txt
var rules string

var group = must.Result(rbnf.New(nil, rules))

// RulesetNames is a slice of all algorithmic rulesets implemented by this package.
var RulesetNames = func() []string {
    xs := group.RulesetNames()
    sort.Strings(xs)
    return xs
}()

// Formatter returns a new [rbnf.Formatter] for a given algorithmic ruleset
// such as "roman-upper".
func Formatter(name string) (rbnf.Formatter, bool) {
    if !strings.HasPrefix(name, "%") { name = "%" + name }
    return group.Formatter(name)
}

// Format formats a number using a given algorithmic ruleset such as
// "roman-upper".
func Format(name string, number int64) (string, error) {
    if !strings.HasPrefix(name, "%") { name = "%" + name }
    return group.FormatInteger(name, number)
}
