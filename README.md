# tawesoft/golib 

[![Go Reference](https://pkg.go.dev/badge/github.com/tawesoft/golib/v2.svg)](https://pkg.go.dev/github.com/tawesoft/golib/v2)
[![Coverage Status](https://coveralls.io/repos/github/tawesoft/golib/badge.svg?branch=v2)](https://coveralls.io/github/tawesoft/golib?branch=v2)

A monorepo for small Go modules maintained by
[Tawesoft®](https://www.tawesoft.co.uk). Open source ([MIT](/LICENSE.txt)).

This is modern code with generics, for go v1.19beta1+


## Packages

|   Name   |       /v2        | /v3 | Description                                             |
|:--------:|:----------------:|:---:|:--------------------------------------------------------|
|  dialog  |  [stable][101]   |  -  | simple cross-platform message boxes & file pickers      |
| digraph  | [unstable][102]  |  -  | directed graphs (including DAGs)                        |
|   drop   |       TODO       |  -  | drop process privileges and inherit handles             |
|  grace   |       TODO       |  -  | start and gracefully shutdown processes                 |
| humanize |       TODO       |  -  | locale-aware parsing & formatting of times & quantities |
|    ks    | [candidate][102] |  -  | misc helpful things                                     |
|   lazy   | [candidate][103] |  -  | composable lazy sequences                               |
|  loader  | [unstable][105]  |  -  | concurrent dependency graph solver                      |
| numbers  |  [stable][106]   |  -  | helpful things for number types                         |

[101]: https://pkg.go.dev/github.com/tawesoft/golib/v2/dialog
[102]: https://pkg.go.dev/github.com/tawesoft/golib/v2/digraph
[103]: https://pkg.go.dev/github.com/tawesoft/golib/v2/ks
[104]: https://pkg.go.dev/github.com/tawesoft/golib/v2/lazy
[105]: https://pkg.go.dev/github.com/tawesoft/golib/v2/loader
[106]: https://pkg.go.dev/github.com/tawesoft/golib/v2/numbers

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
Email [open-source@tawesoft.co.uk](mailto:open-source@tawesoft.co.uk) or visit
[tawesoft.co.uk/products/open-source-software](https://www.tawesoft.co.uk/products/open-source-software) 
to learn more.
