package numbers

// CheckedAddReals returns (a + b, true) iff the sum would not overflow (including
// negative overflow if adding a negative number).
//
// The arguments min and max should be the greatest magnitude negative (zero
// for unsigned) and the greatest magnitude maximum numbers representable by
// the type - see the generic functions [Min] and [Max].
//
// If you know the type of N in advance, you can use [RealInfo.CheckedAdd],
// which populates min and max for you.
func CheckedAddReals[N Real](min N, max N, a N, b N) (N, bool) {
    if (b > 0) && (a > (max - b)) { return 0, false }
    if (b < 0) && (a < (min - b)) { return 0, false }
    return a + b, true
}

// CheckedAddRealsN is like [CheckedAddReals], but returns the sum of any
// number of arguments, not just two.
func CheckedAddRealsN[N Real](min N, max N, xs ... N) (N, bool) {
    var total N
    for i := 0; i < len(xs); i++ {
        if x, ok := CheckedAddReals(min, max, total, xs[0]); ok {
            total = x
        } else {
            return 0, false
        }
    }
    return total, true
}

// CheckedSubReals returns (a - b, true) iff a - b would not overflow (see
// [CheckedAddReals] for notes).
func CheckedSubReals[N Real](min N, max N, a N, b N) (N, bool) {
    if (b < 0) && (a > (max + b)) { return 0, false }
    if (b > 0) && (a < (min + b)) { return 0, false }
    return a - b, true
}

// CheckedAdd returns (a + b, true) iff the sum would not overflow (including
// negative overflow if adding a negative number).
func (t RealInfo[N]) CheckedAdd(a N, b N) (N, bool) {
    return CheckedAddReals(t.Min, t.Max, a, b)
}

// CheckedAddN is like [RealType.CheckedAdd], but returns the sum of any
// number of arguments, not just two.
func (t RealInfo[N]) CheckedAddN(xs ... N) (N, bool) {
    return CheckedAddRealsN(t.Min, t.Max, xs...)
}

// CheckedSub returns (a - b, true) iff a - b would not overflow.
func (t RealInfo[N]) CheckedSub(a N, b N) (N, bool) {
    return CheckedAddReals(t.Min, t.Max, a, b)
}
