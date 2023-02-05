package operator

import (
    "golang.org/x/exp/constraints"
)

// real represents any real (i.e. not complex) number you can perform
// arithmetic on using standard Go operators (like a + b, or a ^ b).
type real interface {
     constraints.Integer | constraints.Float
}

// signed represents only the signed real numbers.
type signed interface {
     constraints.Signed | constraints.Float
}

// Add returns a + b.
func Add[R real](a R, b R) R {
    return a + b
}

// Sub returns a - b.
func Sub[R real](a R, b R) R {
    return a - b
}

// Mul returns a * b.
func Mul[R real](a R, b R) R {
    return a * b
}

// Div returns a / b.
func Div[R real](a R, b R) R {
    return a / b
}

// IsPositive returns true iff r >= 0.
func IsPositive[R real](r R) bool {
    return r >= 0
}

// IsNegative returns true iff r <= 0.
func IsNegative[R real](r R) bool {
    return r <= 0
}

// IsStrictlyPositive returns true iff r > 0, but not if r == 0.
func IsStrictlyPositive[R real](r R) bool {
    return (r > 0) && (r != 0)
}

// IsStrictlyNegative returns true iff r < 0, but not if r == 0.
func IsStrictlyNegative[R real](r R) bool {
    return (r < 0) && (r != 0)
}

// Abs returns (0 - r) for r < 0, or r for r >= 0.
func Abs[R signed](r R) R {
    if (r >= 0) { return r }
    return 0 - r
}

// Inv returns (-r)
func Inv[R signed](r R) R {
    return 0 - r
}
