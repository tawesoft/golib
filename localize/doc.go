// x-doc-short-desc: locale-aware number parsing
// x-doc-stable: yes

// Package localize is an attempt at implementing locale-aware parsing of
// numbers, integrating with golang.org/x/text.
//
// Todo:
//
//   - This is proof of concept and could be tidied up
//   - Checks for integer overflow
//   - Support different representations of negative numbers
//     e.g. `(123)` vs `-123`
//   - In cases where AcceptInteger/AcceptFloat reach a syntax error, they
//     currently underestimate how many bytes they successfully parsed when
//     the byte length of the string is not equal to the number of Unicode
//     code points in the string.
//
package localize
