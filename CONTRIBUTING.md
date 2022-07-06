# Contributing

Open source is a gift to the commons and all contributions - bug reports,
documentation, code, etc. - are welcomed.

## Issues and bug reports

Open an issue on the
[tawesoft/golib issue tracker](https://github.com/tawesoft/golib/issues).

Please make the issue title start with the package relative
path including version (or "all" for everything or misc chores), followed 
by a colon. For example:

```
v2/lazy: error when doing something
v2/all: help getting started
```

Please also label the issue "bug", "enchancement", "question", "help wanted", 
etc.


## Contributing code

### Pull requests

In the git commit message, make the first line (the subject) start with the
package relative path (or "all" for everything or chores), followed by a colon. 
For example:

```
all: update README
lazy: fix something
legacy/foo: fix something else
```

### Copyright

All code contributions **must** be made available under the exact same 
terms as the [LICENSE.txt](/LICENSE.txt) file (and cannot be accepted 
otherwise).

Contributors will have their name and email added to the CONTRIBUTORS file
at the start of each minor version release.

If you are contributing in the course of your employment, and your employer
is the copyright holder of your contribution (this is usually the case when
you are contributing code at work),  then please also add your employer's 
legal name, the country where that company is registered, and the year(s) 
of the contribution(s) to the AUTHORS file (please preserve alphabetical
order).

### Style guide

* Unlike most Go projects, we exclusively use spaces for indentation.
* Use British English (en-GB) for code comments and documentation, US 
  English (en-US) for identifiers and file names.
