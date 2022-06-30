package digraph

import (
    "fmt"
    "math"

    "golang.org/x/exp/constraints"
)

type Number interface {
    constraints.Integer | constraints.Float
}

const negativeInfiniteEdgeCount = math.MinInt32
const positiveInfiniteEdgeCount = math.MaxInt32

func isInfiniteDistance(d DistanceT) bool {
    return (d == positiveInfiniteEdgeCount) || (d == negativeInfiniteEdgeCount)
}

// infinities stores positive and negative infinity for type of number
type infinities[N Number] struct {
    positive N
    negative N
}

// sum returns `a + b`, where a or b may be (possibly signed) infinite.
// `a + Inf == a + Inf == Inf` where `a != -Inf`
// and `a + (-Inf) == (-Inf) + a == -Inf` where `a != Inf`
func (i infinities[N]) sum(a N, b N) N {
    posInf := i.positive
    negInf := i.negative

    if ((a == posInf) || (b == posInf)) && ((a != negInf) && (b != negInf)) {
        return posInf
    }
    if ((a == negInf) || (b == negInf)) && ((a != posInf) && (b != posInf)) {
        return negInf
    }
    if (a == posInf) || (b == posInf) || (a == negInf) || (b == negInf) {
        panic(fmt.Sprintf("undefined operation on %v + %v (inf: +%v, -%v)", a, b, posInf, negInf))
    }

    return a + b
}

// MaxWeights returns a maximum value for a weight of a given type. Useful for
// e.g. the initial value in a `min` reducer function. Not to be confused with
// the [InfWeight] function.
func MaxWeight[N Number]() N {
    var n N
    switch x := any(&n).(type) {
        case   *int:   *x = (     int (math.MaxInt  ));
        case   *int8:  *x = (     int8(math.MaxInt8 ));
        case   *int16: *x = (    int16(math.MaxInt16));
        case   *int32: *x = (    int32(math.MaxInt32));
        case   *int64: *x = (    int64(math.MaxInt64));
        case  *uint:   *x = (   uint (math.MaxUint  ));
        case  *uint8:  *x = (   uint8(math.MaxUint8 ));
        case  *uint16: *x = (  uint16(math.MaxUint16));
        case  *uint32: *x = (  uint32(math.MaxUint32));
        case  *uint64: *x = (  uint64(math.MaxUint64));
        case *float32: *x = (float32(math.MaxFloat32));
        case *float64: *x = (float64(math.MaxFloat64));
        default:
            panic("generic max not implemented for this Number type")
    }
    return n
}

// InfWeight returns an "infinite" value for a weight of a given type, with a
// given sign (i.e. positive or negative infinity). For signed integer values,
// this is the largest and smallest possible values. For unsigned integer
// values, this is the largest possible value minus one (for positive infinity)
// and the largest possible value  (for negative infinity). For floating point
// values, this is the appropriate floating-point infinity from [math.Inf].
func InfWeight[N Number, S Number](sign S) N {
    var n N

    if sign >= 0 {
        switch x := any(&n).(type) {
            case   *int:   *x = (     int (math.MaxInt  ));
            case   *int8:  *x = (     int8(math.MaxInt8 ));
            case   *int16: *x = (    int16(math.MaxInt16));
            case   *int32: *x = (    int32(math.MaxInt32));
            case   *int64: *x = (    int64(math.MaxInt64));
            case  *uint:   *x = (   uint (math.MaxUint   - 1));
            case  *uint8:  *x = (   uint8(math.MaxUint8  - 1));
            case  *uint16: *x = (  uint16(math.MaxUint16 - 1));
            case  *uint32: *x = (  uint32(math.MaxUint32 - 1));
            case  *uint64: *x = (  uint64(math.MaxUint64 - 1));
            case *float32: *x = (float32(math.Inf(1)));
            case *float64: *x = (float64(math.Inf(1)));
            default:
                panic("generic max not implemented for this Number type")
        }
    } else {
    switch x := any(&n).(type) {
        case   *int:   *x = (     int (math.MinInt  ));
        case   *int8:  *x = (     int8(math.MinInt8 ));
        case   *int16: *x = (    int16(math.MinInt16));
        case   *int32: *x = (    int32(math.MinInt32));
        case   *int64: *x = (    int64(math.MinInt64));
        case  *uint:   *x = (   uint (math.MaxUint  ));
        case  *uint8:  *x = (   uint8(math.MaxUint8 ));
        case  *uint16: *x = (  uint16(math.MaxUint16));
        case  *uint32: *x = (  uint32(math.MaxUint32));
        case  *uint64: *x = (  uint64(math.MaxUint64));
        case *float32: *x = (float32(math.Inf(-1)));
        case *float64: *x = (float64(math.Inf(-1)));
        default:
            panic("generic max not implemented for this Number type")
    }
    }
    return n
}

// zero returns the zero value of any type.
func zero[T any]() T {
    var t T
    return t
}

// zeroPtr returns a nil-valued pointer of type T
func zeroPtr[T any]() *T {
    var t *T
    return t
}

// growCap grows a slice of elements of type T to a certain size and capacity
// and returns the buffer
func growCap[T any](buf []T, size int, capacity int) []T {
    if cap(buf) < capacity {
        return make([]T, size, capacity)
    } else {
        return buf[:size]
    }
}

// clear sets every element in a slice of type []T to the zero value of the
// type.
func clear[T any](buf []T) {
    z := zero[T]()
    for i := 0; i < len(buf); i++ {
        buf[i] = z
    }
}

// reverse swaps the order of a slice in place
func reverse[T any](a []T) {
    for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
        a[left], a[right] = a[right], a[left]
    }
}
