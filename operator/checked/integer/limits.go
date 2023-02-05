package integer

import (
    "math"

    "github.com/tawesoft/golib/v2/must"
    "golang.org/x/exp/constraints"
)

// Limits provides a convenient way to fill the min and max arguments to the
// checked operator functions. The inequality Min <= Max must be satisfied.
type Limits[I constraints.Integer] struct {
    Min I
    Max I
}

// GetLimits returns a filled-in [Limit] for the given integer type.
func GetLimits[I constraints.Integer]() Limits[I] {
    var n Limits[I]
    switch x := any(&n).(type) {
        case *Limits[int]:    *x = Int
        case *Limits[int8]:   *x = Int8
        case *Limits[int16]:  *x = Int16
        case *Limits[int32]:  *x = Int32
        case *Limits[int64]:  *x = Int64
        case *Limits[uint]:   *x = Uint
        case *Limits[uint8]:  *x = Uint8
        case *Limits[uint16]: *x = Uint16
        case *Limits[uint32]: *x = Uint32
        case *Limits[uint64]: *x = Uint64
        default:
            must.Never("Limits are not defined for type %T", n)
    }
    return n
}

// Filled-in [Limits] about different integer types with minimum and maximum
// set to the largest range supported by the limit.
var (
    Int   = Limits[int]  {math.MinInt,   math.MaxInt}
    Int8  = Limits[int8] {math.MinInt8,  math.MaxInt8}
    Int16 = Limits[int16]{math.MinInt16, math.MaxInt16}
    Int32 = Limits[int32]{math.MinInt32, math.MaxInt32}
    Int64 = Limits[int64]{math.MinInt64, math.MaxInt64}

    Uint   = Limits[uint]  {0, math.MaxUint}
    Uint8  = Limits[uint8] {0, math.MaxUint8}
    Uint16 = Limits[uint16]{0, math.MaxUint16}
    Uint32 = Limits[uint32]{0, math.MaxUint32}
    Uint64 = Limits[uint64]{0, math.MaxUint64}
)

// Add calls [checked.Add] with min and max filled in with the associated
// [Limits] values.
func (l Limits[I]) Add(a I, b I) (I, bool) {
    return Add(l.Min, l.Max, a, b)
}

// Sub calls [checked.Sub] with min and max filled in with the associated
// [Limits] values.
func (l Limits[I]) Sub(a I, b I) (I, bool) {
    return Sub(l.Min, l.Max, a, b)
}

// Mul calls [checked.Mul] with min and max filled in with the associated
// [Limits] values.
func (l Limits[I]) Mul(a I, b I) (I, bool) {
    return Mul(l.Min, l.Max, a, b)
}

// Abs calls [checked.Abs] with min and max filled in with the associated
// [Limits] values.
func (l Limits[I]) Abs(i I) (I, bool) {
    return Abs(l.Min, l.Max, i)
}

// Inv calls [checked.Inv] with min and max filled in with the associated
// [Limits] values.
func (l Limits[I]) Inv(i I) (I, bool) {
    return Inv(l.Min, l.Max, i)
}
