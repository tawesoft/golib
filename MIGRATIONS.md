# Migrating from Earlier Versions

This document describes any breaking changes between versions,
and how to transition.

## Versioning

Packages in this project generally have
[normal stability guarantees](https://go.dev/doc/modules/version-numbers) (for
a project of v2+).

However, we further mark each individual package with a stability rating to
manage expectations. Packages marked as "candidate" or "unstable" do not
have the normal stability guarantee. Instead, they have the guarantees
described in the following table.

| Stability | Meaning                                                                                                                  |
|:---------:|--------------------------------------------------------------------------------------------------------------------------|
|  Legacy   | There will be no changes to this package, and it may be removed in the next major version                                |
|  Normal   | Breaking changes to this package are possible, but only between major versions                                           |
| Candidate | Small breaking changes to this package are possible, even between minor versions                                         |
| Unstable  | Large breaking changes to this package are likely, even between minor versions, and may not have migration instructions  |

Security fixes may not always be backwards compatible, even between minor 
versions.

Downstream module authors should avoid using a candidate package
in a stable release of their module. Candidate packages are appropriate for 
"release candidate" pre-release versions of modules, development versions of 
modules, and `main()` programs.


## Updating `github.com/tawesoft/golib`

### Migrating v2.0 → v2.1

Placeholder


### Migrating from `tawesoft.co.uk/go` → v2.0

#### Requirements

* Update to go v1.19 or better

#### Update imports

* `tawesoft.co.uk/go/dialog` → `github.
  com/tawesoft/golib/v2/dialog`
 
* `tawesoft.co.uk/go/lxstrconv` → `github.com/tawesoft/golib/v2/lxstrconv`
