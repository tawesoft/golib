# tawesoft/golib 

[![Go Reference](https://pkg.go.dev/badge/github.com/tawesoft/golib/v2.svg)](https://pkg.go.dev/github.com/tawesoft/golib/v2)
[![Coverage Status](https://coveralls.io/repos/github/tawesoft/golib/badge.svg?branch=v2)](https://coveralls.io/github/tawesoft/golib?branch=v2)

A monorepo for small Go modules maintained by
[Tawesoft®](https://www.tawesoft.co.uk). Open source ([MIT](/LICENSE.txt)).

This is modern code with generics, for go v1.19

## Packages

|       Name       |       /v2        | /v3 | Description                                               |
|:----------------:|:----------------:|:---:|:----------------------------------------------------------|
|      dialog      |  [stable][d01]   |  -  | cross-platform message boxes & file pickers               |
|     digraph      | [unstable][d02]  |  -  | directed graphs (including DAGs)                          |
|       drop       |       TODO       |  -  | drop process privileges and inherit handles               |
|    fun/maybe     | [candidate][f01] |  -  | implements a "Maybe" sum type                             |
|   fun/partial    |  [stable][f02]   |  -  | partial function application                              |
|    fun/result    | [candidate][f03] |  -  | implements a "Result" sum type                            |
|    fun/slices    | [candidate][f04] |  -  | higher-order functions (map, flatmap) for slices          |
|      grace       |       TODO       |  -  | start and gracefully shutdown processes                   |
|     humanize     |       TODO       |  -  | locale-aware parsing & formatting of times & quantities   |
|       iter       | [candidate][i01] |  -  | composable lazy iteration over sequences                  |
|        ks        | [unstable][k01]  |  -  | "kitchen sink" of misc helpful things                     |
|      loader      |       TODO       |  -  | concurrent dependency graph solver                        |
|  meta/opengraph  | [candidate][m01] |  -  | render Open Graph protocol HTML meta tags                 |
| meta/twittercard | [candidate][m02] |  -  | render Twitter Card HTML meta tags                        |
|       must       | [candidate][m03] |  -  | assertions                                                |
|     numbers      | [candidate][n01] |  -  | helpful things for number types                           |
|     text/ccc     | [candidate][t01] |  -  | Unicode Canonical Combining Class values                  |
|     text/dm      | [candidate][t02] |  -  | Unicode decomposition mappings & selective decompositions |
|  text/fallback   | [candidate][t03] |  -  | Unicode Character Fallback Substitutions                  | 
|    text/fold     | [candidate][t04] |  -  | selectively merge distinctions in Unicode text            |
|     text/np      | [candidate][t05] |  -  | Unicode numeric properties                                |
|       view       | [candidate][v01] |  -  | dynamic views over collections                            |

[d01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/dialog
[d02]: https://pkg.go.dev/github.com/tawesoft/golib/v2/digraph
[f01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/fun/maybe
[f02]: https://pkg.go.dev/github.com/tawesoft/golib/v2/fun/partial
[f03]: https://pkg.go.dev/github.com/tawesoft/golib/v2/fun/result
[f04]: https://pkg.go.dev/github.com/tawesoft/golib/v2/fun/slices
[i01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/iter
[k01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/ks
[m01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/meta/opengraph
[m02]: https://pkg.go.dev/github.com/tawesoft/golib/v2/meta/twittercard
[m03]: https://pkg.go.dev/github.com/tawesoft/golib/v2/must
[n01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/numbers
[t01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/ccc
[t02]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/dm
[t03]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/fallback
[t04]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/fold
[t05]: https://pkg.go.dev/github.com/tawesoft/golib/v2/text/np
[v01]: https://pkg.go.dev/github.com/tawesoft/golib/v2/view

**Note:** Additional v2/legacy packages exist for users migrating from
`tawesoft.co.uk/go`. See [MIGRATIONS.md](/MIGRATIONS.md).

**Note:** Packages have the
[normal stability guarantees](https://go.dev/doc/modules/version-numbers)
expected for a Go package of v2 or higher, except where marked 
"partial", "candidate", or "unstable". See 
[MIGRATIONS.md](/MIGRATIONS.md) 
for the meaning of other terms. 

## Support

### Free and Community Support

Use the [tawesoft/golib issue tracker](), powered by GitHub issues.

### Commercial Support

Open source software from Tawesoft® is backed by commercial support options.
Email [opensource@tawesoft.co.uk](mailto:opensource@tawesoft.co.uk) or visit
[tawesoft.co.uk/products/open-source-software](https://www.tawesoft.co.uk/products/open-source-software) 
to learn more.
