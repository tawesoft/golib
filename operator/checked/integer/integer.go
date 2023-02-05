// Package integer (operator/checked/integer) implements operations on integers
// that are robust in the event of integer overflow.
package integer

import (
    "golang.org/x/exp/constraints"
)

// Add returns (a + b, true) iff the result lies between min and max inclusive,
// otherwise returns (0, false). This calculation is robust in the event of
// integer overflow.
//
// The input arguments must satisfy the inequalities: `min <= a <= max` and
// `min <= b <= max`.
func Add[N constraints.Integer](min N, max N, a N, b N) (N, bool) {
    if (b > 0) && (a > (max - b)) { return 0, false }
    if (b < 0) && (a < (min - b)) { return 0, false }
    return a + b, true
}

// Sub returns (a + b, true) iff the result lies between min and max inclusive,
// otherwise returns (0, false). This calculation is robust in the event of
// integer overflow.
//
// The input arguments must satisfy the inequalities: `min <= a <= max` and
// `min <= b <= max`.
func Sub[N constraints.Integer](min N, max N, a N, b N) (N, bool) {
    if (b < 0) && (a > (max + b)) { return 0, false }
    if (b > 0) && (a < (min + b)) { return 0, false }
    return a - b, true
}

// Mul returns (a * b, true) iff the result lies between min and max inclusive,
// otherwise returns (0, false). This calculation is robust in the event of
// integer overflow.
//
// The input arguments must satisfy the inequalities: `min <= a <= max` and
// `min <= b <= max`.
func Mul[N constraints.Integer](min N, max N, a N, b N) (N, bool) {
    if (a == 0) || (b == 0) { return 0, true }

    x := a * b
    if (x < min) || (x > max) { return 0, false }
    if (a != x/b) { return 0, false }
    return x, true
}

// Abs returns (-i, true) for i < 0, or (i, true) for i >= 0 iff the result lies
// between min and max inclusive. Otherwise returns (0, false).
//
// The input arguments must satisfy the inequality `min <= i <= max`.
func Abs[N constraints.Integer](min N, max N, i N) (N, bool) {
    if (i >= 0) { return i, true }
    return Sub(min, max, 0, i)
}

// Inv returns (-i) iff the result lies between min and max inclusive.
// Otherwise returns (0, false).
//
// The input arguments must satisfy the inequality `min <= i <= max`.
func Inv[N constraints.Integer](min N, max N, i N) (N, bool) {
    return Sub(min, max, 0, i)
}
