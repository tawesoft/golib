# Migrating from Earlier Versions

This document describes any breaking changes between versions,
and how to transition.

## Versioning

Packages in this project generally have
[normal stability guarantees](https://go.dev/doc/modules/version-numbers) (for
a project of v2+).

However, we further mark each individual package with a stability rating to
manage expectations. Packages marked as "partial", "candidate" or "unstable" 
do not have the normal stability guarantees. Instead, they have the guarantees
described in the following table.

| Stability | Meaning                                                                                                                  |
|:---------:|--------------------------------------------------------------------------------------------------------------------------|
|  Frozen   | There will be no breaking changes to this package. There will be no new versions.                                        |
|  Normal   | There will be no breaking changes to this package, except between major versions.                                        |
|  Partial  | As "normal", with a documented exception for individual parts.                                                           |
| Candidate | Small breaking changes to this package are possible, even between minor versions.                                        |
| Unstable  | Large breaking changes to this package are likely, even between minor versions, and may not have migration instructions. |

Security fixes may not always be backwards compatible, even between minor 
versions.

Downstream module authors should avoid using a candidate package
in a stable release of their module, particularly in their public API. 
Candidate packages are appropriate for "release candidate" pre-release 
versions of modules, development versions of modules, and `main()` programs. 
Alternatively, please vendor the package as a temporary measure.


## Updating `github.com/tawesoft/golib`

### Migrating v2.0 → v2.1

Placeholder



## Migrating from `tawesoft.co.uk/go`

There are no breaking changes, but you should update your imports as indicated
below for each package.


### Package **tawesoft.co.uk/go/dialog:**

#### Update imports:

```diff
- import "tawesoft.co.uk/go/dialog"
+ import "github.com/tawesoft/golib/v2/dialog"
```

#### (Optional) update code

Although the new package is API-compatible with the old one, check out the
[new features](https://pkg.go.dev/github.com/tawesoft/golib/v2/dialog) too.

```diff
- dialog.Alert(...) // deprecated but still works
+ err := dialog.Raise(...) // preferred
+ err := dialog.MessageBox{...}.Raise() // more options
```

### Package **tawesoft.co.uk/go/humanizex:**

#### Update imports:

```diff
- import "tawesoft.co.uk/go/humanizex"
+ import humanizex "github.com/tawesoft/golib/v2/legacy/humanize"
```

This `golib/v2/legacy/humanize` package is frozen and will not appear 
in v3. It will always be available at the v2 import path.


### Package **tawesoft.co.uk/go/lxstrconv:**

#### Update imports:

```diff
- import "tawesoft.co.uk/go/lxstrconv"
+ import lxstrconv "github.com/tawesoft/golib/v2/legacy/localize"
```

This `golib/v2/legacy/localize` package is frozen and will not appear 
in v3. It will always be available at the v2 import path.


### Package **tawesoft.co.uk/go/operator:**

#### Update imports:

```diff
- import "tawesoft.co.uk/go/operator"
+ import "github.com/tawesoft/golib/v2/legacy/operator"
```

This `golib/v2/legacy/operator` package is frozen and will not appear 
in v3. It will always be available at the v2 import path.

### (Optional) migrate operator → **goblib/v2/number:**

Package `goblib/v2/number` provides similar features to the old
`tawesoft.co.uk/go/operator` package. It uses generics and will be maintained
going forward. It's fairly easy to migrate existing code to use the new golib 
package.

#### Update imports:

```diff
- import "tawesoft.co.uk/go/operator"
+ import "github.com/tawesoft/golib/v2/numbers"
```

#### Rewrite checked code

Note that checked functions now return (value, ok), not (value, error), so 
the return value logic is inverted:

```diff
- if sum, err := operator.IntChecked.Binary.Add(a, b); err == nil { ...
+ if sum, ok := numbers.Int.CheckedAdd(a, b); ok { ...

- if sum, err := operator.IntChecked.Nary.AddN(a, b, c, d); err == nil { ...
+ if sum, ok := numbers.Int.CheckedAddN(a, b, c, d); ok { ...
```

Or, generically:

```diff
func foo[N numbers.Real](a N, b N) N {
    min := numbers.Min[N]()
    max := numbers.Min[N]() 
    sum, ok := numbers.CheckedAddReals(min, max, a, b)
    ...
}
```

#### Rewrite non-checked code

Remaining code is trivial to update, and types are automatically detected by 
the Go compiler.

```diff
- operator.Int.Binary.Add(2, 3) // 5
+ numbers.Add(2, 5) // 5

- reduce(operator.Int.Binary.Add, []int{1, 2, 3})
+ reduce(numbers.Add, []int{1, 2, 3})

- sum := operator.Int.Nary.Add([]int{1, 2, 3}...)
+ sum := numbers.AddN([]int{1, 2, 3}...)
```
