# Migrating from Earlier Versions

This document describes any breaking changes between versions,
and how to transition.

## Versioning

Packages in this project generally have
[normal stability guarantees](https://go.dev/doc/modules/version-numbers) (for
a project of v2+) if they are listed under "stable" in the README table.

However, packages listed under "latest" do not have the normal stability
guarantees. Instead, they have the guarantees described in the following table.

|      Stability      | Meaning                                                                                                                          |
|:-------------------:|----------------------------------------------------------------------------------------------------------------------------------|
|       Stable        | There will be no breaking changes to this package, except between major versions.                                                |
|       Latest        | Small breaking changes to this package are possible, even between minor versions. They will usually have migration instructions. |
| Latest *(unstable)* | Large breaking changes to this package are likely, even between minor versions, and may not have migration instructions.         |

Security fixes may cause breaking changes if necessary, even across 
minor version changes.


## Updating `github.com/tawesoft/golib`

### Migrating v2.0 → v2.7

Fewer breaking changes are expected after this point.

* `meta` packages have been moved to `html/meta`
* `fun/maybe`, `fun/result`, and `view` packages have been redone
* `numbers` package removed with much implemented by the new `operator` package
* some functions removed from `fun/result`, `iter`, `ks` packages
* some functions removed from `ks` to `operator`

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

### (Optional) migrate operator → **goblib/v2/operator:**

Package `goblib/v2/operator` provides similar features to the old
`tawesoft.co.uk/go/operator` package. It uses generics and will be maintained
going forward. It's fairly easy to migrate existing code to use the new golib 
package.

#### Update imports:

```diff
- import "tawesoft.co.uk/go/operator"
+ import "github.com/tawesoft/golib/v2/operator"
+ import "github.com/tawesoft/golib/v2/operator/checked/integer"
```

#### Rewrite checked code

Note that checked functions now return (value, ok), not (value, error), so 
the return value logic is inverted:

```diff
- if sum, err := operator.IntChecked.Binary.Add(a, b); err == nil { ...
+ if sum, ok := integer.Int32.Add(a, b); ok { ...
```

Or, generically:

```diff
func foo[N constraints.Integer](a N, b N) N {
    limits := integer.Limits[N]()
    sum, ok := limits.Add(a, b)
    ...
}
```

#### Rewrite non-checked code

Remaining code is trivial to update, and types are often automatically 
detected by the Go compiler.

```diff
- operator.Int.Binary.Add(2, 3) // 5
+ operator.Add(2, 5) // 5

- reduce(operator.Int.Binary.Add, []int{1, 2, 3})
+ reduce(operator.Add[int], []int{1, 2, 3})
```
