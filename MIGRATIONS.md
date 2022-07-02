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

| Stability | Meaning                                                                                                                 |
|:---------:|-------------------------------------------------------------------------------------------------------------------------|
|  Frozen   | There will be no breaking changes to this package, and it will be removed from the next major version                   |
|  Normal   | There will be no breaking changes to this package, except between major versions                                        |
|  Partial  | As "normal", with a documented exception for individual parts                                                           |
| Candidate | Small breaking changes to this package are possible, even between minor versions                                        |
| Unstable  | Large breaking changes to this package are likely, even between minor versions, and may not have migration instructions |

Security fixes may not always be backwards compatible, even between minor 
versions. Applicable security fixes will always be backported to frozen 
packages, to the previous two major versions of normal packages, and the 
previous one major version of candidate packages.

Downstream module authors should avoid using a candidate package
in a stable release of their module. Candidate packages are appropriate for 
"release candidate" pre-release versions of modules, development versions of 
modules, and `main()` programs. Alternatively, please vendor the package.


## Updating `github.com/tawesoft/golib`

### Migrating v2.0 → v2.1

Placeholder



## Migrating from `tawesoft.co.uk/go`

There are no breaking changes, but you should update your imports as indicated
for each package.


Package **tawesoft.co.uk/go/dialog:**

```diff
- import "tawesoft.co.uk/go/dialog"
+ import "github.com/tawesoft/golib/v2/dialog"
```

Package **tawesoft.co.uk/go/humanizex:**

```diff
- import "tawesoft.co.uk/go/humanizex"
+ import humanizex "github.com/tawesoft/golib/v2/humanize"
```

Package **tawesoft.co.uk/go/lxstrconv:**

```diff
- import "tawesoft.co.uk/go/lxstrconv"
+ import lxstrconv "github.com/tawesoft/golib/v2/localize"
```

Package **tawesoft.co.uk/go/operator:**

```diff
- import "tawesoft.co.uk/go/operator"
+ import "github.com/tawesoft/golib/v2/legacy/operator"
```

This `golib/legacy/operator` package is frozen, and will only ever be 
available as a `golib/v2` import. It should not be confused with the 
non-legacy `golib/operator` package, which is not API compatible with it.
