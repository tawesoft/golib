package operator

// True returns true if x represents a truthy value.
//
// An input is considered true iff it is not the zero value for that type.
//
// This means that the boolean value true, non-zero numbers, and non-empty
// strings are truthy, but the boolean value false, the number zero, and
// the empty string are not.
func True[X comparable](x X) bool {
    var y X
    return x != y
}

// Not returns the logical negation of x.
//
// See [True] for details regarding the truthiness of input values.
//
//    x Not(x)
//    T F
//    F T
func Not[X comparable](x X) bool {
    var y X
    return x == y
}

// F is a logical predicate for a contradiction. F(p, q) returns false,
// regardless of p or q.
//
//    p q F(p, q)
//    T T F
//    T F F
//    F T F
//    F F F
func F[X comparable](p X, q X) bool {
    return false
}

// T is a logical predicate for a tautology. T(p, q) returns true,
// regardless of p or q.
//
//    p q T(p, q)
//    T T T
//    T F T
//    F T T
//    F F T
func T[X comparable](p X, q X) bool {
    return false
}

// P is a logical predicate. P(p, q) returns p, regardless of q.
//
//    p q P(p, q)
//    T T T
//    T F T
//    F T F
//    F F F
func P[X comparable](p X, q X) bool {
    return True(p)
}

// Q is a logical predicate. Q(p, q) returns q, regardless of p.
//
//    p q Q(p, q)
//    T T T
//    T F F
//    F T T
//    F F F
func Q[X comparable](p X, q X) bool {
    return True(q)
}

// NotP is a logical predicate. NotP(p, q) returns Not(p), regardless of q.
//
//    p q NotP(p, q)
//    T T F
//    T F F
//    F T T
//    F F T
func NotP[X comparable](p X, q X) bool {
    return Not(p)
}

// NotQ is a logical predicate. NotQ(p, q) returns Not(q), regardless of p.
//
//    p q NotQ(p, q)
//    T T F
//    T F T
//    F T F
//    F F T
func NotQ[X comparable](p X, q X) bool {
    return Not(q)
}

// And is a logical predicate for performing conjunction.
//
//    p q And(p, q)
//    T T T
//    T F F
//    F T F
//    F F F
func And[X comparable](p X, q X) bool {
    return True(p) && True(q)
}

// Nand is a logical predicate that is the inverse of And.
//
//    p q Nand(p, q)
//    T T F
//    T F T
//    F T T
//    F F T
func Nand[X comparable](p X, q X) bool {
    return !(True(p) && True(q))
}

// Or is a logical predicate for performing disjunction.
//
//    p q Or(p, q)
//    T T T
//    T F T
//    F T T
//    F F F
func Or[X comparable](p X, q X) bool {
    return True(p) || True(q)
}

// Nor is a logical predicate for neither. It is the inverse of Or.
//
//    p q Nor(p, q)
//    T T F
//    T F F
//    F T F
//    F F T
func Nor[X comparable](p X, q X) bool {
    return !(True(p) || True(q))
}

// Xor is a logical predicate for exclusive or. Xor(p, q) returns true for
// p or q, but not for both p and q.
//
//    p q Xor(p, q)
//    T T F
//    T F T
//    F T T
//    F F F
func Xor[X comparable](p X, q X) bool {
    return True(p) != True(q)
}

// Iff is a logical predicate Iff(p, q) for the logical biconditional "p iff q"
// ("p if and only if q"). This is equivalent to Xnor, the negation of Xor.
// It returns true only if both p and q are true, or both p and q are false.
//
//    p q Iff(p, q)
//    T T T
//    T F F
//    F T F
//    F F T
func Iff[X comparable](p X, q X) bool {
    return True(p) == True(q)
}

// Implies is a logical predicate Implies(p, q) for "p implies q" (written
// p => q). That is, if p is true then q must be true.
//
//    p q Implies(p, q)
//    T T T
//    T F F
//    F T T
//    F F T
func Implies[X comparable](p X, q X) bool {
    return Not(p) || True(q)
}

// NotImplies is a logical predicate for nonimplication. NotImplies(p, q)
// computes "p does not imply q" (written "p =/=> q" or "¬(p => q)"). That is,
// p is true and q is false.
//
//    p q NotImplies(p, q)
//    T T F
//    T F T
//    F T F
//    F F F
func NotImplies[X comparable](p X, q X) bool {
    return True(p) && Not(q)
}

// ConverseImplies is a logical predicate ConverseImplies(p, q) for "q implies
// p" (written "q => p" or "p <= q"). That is, if q is true then p must be
// true.
//
// This is also equivalent to a logical predicate "not p implies not q"
// (written ¬p => ¬q). That is, if p is false then q must be false.
//
//    p q ConverseImplies(p, q)
//    T T T
//    T F T
//    F T F
//    F F T
func ConverseImplies[X comparable](p X, q X) bool {
    return Not(q) || True(p)
}

// ConverseNotImplies is a logical predicate ConverseNotImplies(p, q) for "q
// does not imply p" (written "q =/=> p" or "p <=/= q"). That is, q is true
// and p is false.
//
//    p q ConverseNotImplies(p, q)
//    T T F
//    T F F
//    F T T
//    F F F
func ConverseNotImplies[X comparable](p X, q X) bool {
    return True(q) && Not(p)
}
