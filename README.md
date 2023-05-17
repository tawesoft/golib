# tawesoft/golib 

[![Go Reference](https://pkg.go.dev/badge/github.com/tawesoft/golib/v2.svg)](https://pkg.go.dev/github.com/tawesoft/golib/v2)
[![Coverage Status](https://coveralls.io/repos/github/tawesoft/golib/badge.svg?branch=v2)](https://coveralls.io/github/tawesoft/golib?branch=v2)

A monorepo for small Go (v1.20+) modules maintained by
[Tawesoft®](https://www.tawesoft.co.uk), with few dependencies.

```go
import "github.com/tawesoft/golib/v2/..."
```

This is free &amp; open source software made available under the
[MIT Licence](/LICENSE.txt).

Some portions, particularly portions relating to CSS processing and Unicode 
support, are additionally covered by compatible [MIT-like licences](/LICENSE-PARTS.txt).


## Packages

### General Packages

|         Name          |  Stable   |  Latest   | Description                                          |
|:---------------------:|:---------:|:---------:|:-----------------------------------------------------|
|     css/tokenizer     |     -     | [v2][c01] | CSS tokenizer for [CSS Syntax Module Level 3][css1]  |
|        dialog         | [v2][d01] |     -     | cross-platform message boxes & file pickers          |
|        digraph        |     -     | [v2][d02] | *(unstable)* directed graphs (including DAGs)        |
|         drop          |     -     |     -     | *(TODO)* drop process privileges and inherit handles |
|      fun/either       | [v2][f01] |     -     | "Either" sum type                                    |
|      fun/future       | [v2][f02] |     -     | synchronous and asynchronous future  values          |
|       fun/maybe       | [v2][f03] |     -     | "Maybe" sum type                                     |
|      fun/partial      | [v2][f04] |     -     | partial function application                         |
|      fun/promise      | [v2][f05] |     -     | store computations to be performed later             |
|      fun/result       | [v2][f06] |     -     | "Result" sum type                                    |
|      fun/slices       | [v2][f07] |     -     | higher-order functions for slices                    |
|         grace         |     -     |     -     | *(TODO)* start and gracefully shutdown processes     |
|       humanize        |     -     |     -     | *(TODO)* locale-aware numbers &amp; quantities       |
|         iter          | [v2][i01] |     -     | composable lazy iteration                            |
|          ks           |     -     | [v2][k01] | *(unstable)* "kitchen sink" of extras                |
|        loader         |     -     |     -     | *(TODO)* concurrent dependency graph solver          |
|  html/meta/opengraph  | [v2][h01] |     -     | HTML meta tags for Facebook's Open Graph protocol    |
| html/meta/twittercard | [v2][h02] |     -     | HTML meta tags for Twitter Cards                     |
|         must          | [v2][m03] |     -     | assertions                                           |
|       operator        | [v2][o01] |     -     | operators as functions                               |
|         tuple         | [v2][p01] |     -     | convert to/from tuples                               |
|         view          | [v2][v01] |     -     | dynamic views over collections                       |

**Note:** Additional v2/legacy packages exist for users migrating from
`tawesoft.co.uk/go`. See [MIGRATIONS.md](/MIGRATIONS.md).

**Note:** "Stable" packages have the
[normal stability guarantees](https://go.dev/doc/modules/version-numbers)
expected for a Go package of v2 or higher. "Latest" packages, or
"Latest *(unstable)*" packages do not. See [MIGRATIONS.md](/MIGRATIONS.md). 

### Text Packages

|          Name           |  Stable   |  Latest   | Description                                               |
|:-----------------------:|:---------:|:---------:|:----------------------------------------------------------|
|        text/ccc         |     -     | [v2][t01] | Unicode Canonical Combining Class values                  |
|         text/dm         |     -     | [v2][t02] | Unicode decomposition mappings & selective decompositions |
|      text/fallback      |     -     | [v2][t03] | Unicode Character Fallback Substitutions                  | 
|        text/fold        |     -     | [v2][t04] | Unicode text folding                                      |
|         text/np         |     -     | [v2][t05] | Unicode numeric properties                                |
| text/number/algorithmic | [v2][t07] |     -     | CLDR algorithmic (non-decimal) numbering systems          |
|   text/number/plurals   | [v2][t08] |     -     | CLDR plural rules with a simple interface                 |
|    text/number/rbnf     |     -     | [v2][t09] | CLDR Rule-Based Number Formats                            |
|   text/number/symbols   |     -     | [v2][t10] | CLDR locale-appropriate Number Symbols                    |


**Note:** "Stable" packages have the
[normal stability guarantees](https://go.dev/doc/modules/version-numbers)
expected for a Go package of v2 or higher. "Latest" packages, or
"Latest *(unstable)*" packages do not. See [MIGRATIONS.md](/MIGRATIONS.md). 

[css1]: https://www.w3.org/TR/css-syntax-3/
[c01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/css/tokenizer
[d01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/dialog
[d02]: https://pkg.go.dev/github.com/tawesoft/golib/v2/digraph
[f01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/fun/either
[f02]: https://pkg.go.dev/github.com/tawesoft/golib/v2/fun/future
[f03]: https://pkg.go.dev/github.com/tawesoft/golib/v2/fun/maybe
[f04]: https://pkg.go.dev/github.com/tawesoft/golib/v2/fun/partial
[f05]: https://pkg.go.dev/github.com/tawesoft/golib/v2/fun/promise
[f06]: https://pkg.go.dev/github.com/tawesoft/golib/v2/fun/result
[f07]: https://pkg.go.dev/github.com/tawesoft/golib/v2/fun/slices
[i01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/iter
[k01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/ks
[h01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/meta/opengraph
[h02]: https://pkg.go.dev/github.com/tawesoft/golib/v2/meta/twittercard
[m03]: https://pkg.go.dev/github.com/tawesoft/golib/v2/must
[o01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/operator
[p01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/tuple
[t01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/ccc
[t02]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/dm
[t03]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/fallback
[t04]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/fold
[t05]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/np
[t06]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/runeio
[t07]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/number/algorithmic
[t08]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/number/plurals
[t09]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/number/rbnf
[t10]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/number/symbols
[v01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/view


## Support

### Free and Community Support

Use the [tawesoft/golib issue tracker](), powered by GitHub issues.

### Commercial Support

Open source software from Tawesoft® is backed by commercial support options.
Email [opensource@tawesoft.co.uk](mailto:opensource@tawesoft.co.uk) or visit
[tawesoft.co.uk/products/open-source-software](https://www.tawesoft.co.uk/products/open-source-software) 
to learn more.
