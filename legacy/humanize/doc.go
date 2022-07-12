// Package humanize is an elegant, general-purpose, extensible, modular,
// locale-aware way to format and parse numbers and quantities - like distances,
// bytes, and time - in a human-readable way ideal for config files and as a
// building-block for fully translated ergonomic user interfaces.
//
// This is the frozen version of the package previously at
// `tawesoft.co.uk/go/humanizex`. See [migration instructions].
//
// [migration instructions]: https://github.com/tawesoft/golib/blob/v2/MIGRATIONS.md#package-tawesoftcoukgooperator
//
// ## Alternative - What about dustin's go-humanize?
//
// dustin's go-humanize (https://github.com/dustin/go-humanize) is 3.9 to 4.5
// times faster formatting and 2 times faster parsing, if this is a bottleneck
// for you. It's also quite mature, so is probably very well tested by now. If
// you're only targeting the English language it also has more handy "out of
// the box" features.
//
// On the other hand, tawesoft's humanizex is more general purpose and has
// better localisation support.
package humanize // import "tawesoft.co.uk/go/humanizex"
