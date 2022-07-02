// Package pkgexample is an example
//
// [encoding/json.Marshal] is a doc link, and so is [rsc.io/quote.Glass].
//
// So is [sort] and [sort.Strings], but not [sort.foo]
//
// So is [builtin.Integer] or even [builtin.comparable] or [comparable]
package pkgexample

import (
    "encoding/json"
)

func foobar() {
    _, _ = json.Marshal(nil)
}
