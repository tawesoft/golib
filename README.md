# tawesoft/golib 

[![Go Reference](https://pkg.go.dev/badge/github.com/tawesoft/golib/v2.svg)](https://pkg.go.dev/github.com/tawesoft/golib/v2)
[![Coverage Status](https://coveralls.io/repos/github/tawesoft/golib/badge.svg?branch=v2)](https://coveralls.io/github/tawesoft/golib?branch=v2)

A monorepo for small Go modules maintained by
[Tawesoft®](https://www.tawesoft.co.uk). Open source ([MIT](/LICENSE.txt)).

This is modern code with generics, for go v1.19beta1+


## Packages


|      Name       |   v2 Status    | v3 Status | Description                                 |
|:---------------:|:--------------:|:---------:|:--------------------------------------------|
|     dialog      |     stable     |     -     | simple cross-platform message boxes         |
|     digraph     | [unstable][1]  |     -     | directed graphs (including DAGs)            |
|      drop       |   candidate    |     -     | drop process privileges and inherit handles |
|      grace      |   candidate    |     -     | start and gracefully shutdown processes     |
|    humanize     |     stable     | unstable  | locale-aware natural number formatting      |
|       ks        | [candidate][2] |     -     | misc helpful things                         |
|      lazy       | [candidate][3] |     -     | lazy evaluation                             |
| legacy/operator |  [frozen][4]   |     -     | operators as functions                      |
|     loader      | [unstable][5]  |     -     | concurrent dependency graph solver          |
|    localize     |  [stable][6]   | unstable  | locale-aware number parsing                 |
|     numbers     |  [stable][7]   |     -     | helpful things for number types             |


[1]: https://pkg.go.dev/github.com/tawesoft/golib/v2/digraph
[2]: https://pkg.go.dev/github.com/tawesoft/golib/v2/ks
[3]: https://pkg.go.dev/github.com/tawesoft/golib/v2/lazy
[4]: https://pkg.go.dev/github.com/tawesoft/golib/v2/legacy/operator
[5]: https://pkg.go.dev/github.com/tawesoft/golib/v2/loader
[6]: https://pkg.go.dev/github.com/tawesoft/golib/v2/localize
[7]: https://pkg.go.dev/github.com/tawesoft/golib/v2/numbers

**Note:** "Frozen" or "normal" stability packages have the
[normal stability guarantees](https://go.dev/doc/modules/version-numbers)
expected for a Go package of v2 or higher. See [MIGRATIONS.md](/MIGRATIONS.md) 
for the meaning of other terms. 

## Support

### Free and Community Support

Use the [tawesoft/golib issue tracker](), powered by GitHub issues.

### Commercial Support

Open source software from Tawesoft® is backed by commercial support options.
Email [open-source@tawesoft.co.uk](mailto:open-source@tawesoft.co.uk) or visit
[tawesoft.co.uk/products/open-source-software](https://www.tawesoft.co.uk/products/open-source-software) 
to learn more.
