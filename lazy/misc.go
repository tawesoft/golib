package lazy

// IfThenElse returns a lazily-evaluated value based on a boolean condition, q.
// Iff q is true, returns the return value of ifTrue(). Iff q is false, returns
// the return value of ifFalse(). This [IfThenElse expression] (as distinct from
// If-Then-Else statements) is much like the ternary operator in some other
// languages.
//
// For a non-lazy version, see [ks.IfThenElse].
//
// [IfThenElse expression]: https://en.wikipedia.org/wiki/Conditional_(computer_programming)#If%E2%80%93then%E2%80%93else_expressions
func IfThenElse[X any] (
    q       bool,
    ifTrue  func() X,
    ifFalse func() X,
) X {
    if q {
        return ifTrue()
    } else {
        return ifFalse()
    }
}
