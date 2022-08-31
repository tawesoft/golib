// Package numbers implements assorted things that operate on numbers, such
// as generic access to limits, addition checked for integer overflow, and
// functions implementing builtin operators like addition so that they can
// be passed to higher order functions.
package numbers

// NOTE: the constraints package is experimental, but it's easy to avoid using
// it, so we don't.

import (
    "math"

    "github.com/tawesoft/golib/v2/ks"
)

// Complex represents any complex number
type Complex interface {
    ~complex64 | ~complex128
}

// Number represents anything you can perform arithmetic on using standard Go
// operators (like a + b, or a ^ b).
//
// See also [Real], which doesn't include the complex numbers.
type Number interface {
     ~int8 |  ~int16 |  ~int32 |   ~int64 |  ~int |
    ~uint8 | ~uint16 | ~uint32 |  ~uint64 | ~uint |
                      ~float32 | ~float64 |
                               ~complex64 | ~complex128
}

// Real represents any real (i.e. not complex) number you can perform arithmetic
// on using standard Go operators (like a + b, or a ^ b).
type Real interface {
     ~int8 |  ~int16 |  ~int32 |   ~int64 |  ~int |
    ~uint8 | ~uint16 | ~uint32 |  ~uint64 | ~uint |
                      ~float32 | ~float64
}

// Smallest returns the smallest representable real number greater than zero.
// For all integer types, this is always 1. For float types, this is
// math.SmallestNonzeroFloat32 or math.SmallestNonzeroFloat64.
func Smallest[N Real]() N {
    var n N
    switch x := any(&n).(type) {
        case      *int:   *x = 1
        case      *int8:  *x = 1
        case      *int16: *x = 1
        case      *int32: *x = 1
        case      *int64: *x = 1
        case     *uint:   *x = 1
        case     *uint8:  *x = 1
        case     *uint16: *x = 1
        case     *uint32: *x = 1
        case     *uint64: *x = 1
        case    *float32: *x = float32(math.SmallestNonzeroFloat32)
        case    *float64: *x = float64(math.SmallestNonzeroFloat64)
        default: ks.Never()
    }
    return n
}

// Max returns the maximum representable number of type N.
//
// For integers, min <= 0 < max, and min != -max.
//
// For floating point numbers, -inf < min < 0 < epsilon < max < inf, and min ==
// -max.
//
// e.g. Max[uint8]() // returns 255
func Max[N Number]() N {
    var n N
    switch x := any(&n).(type) {
        case      *int:   *x =    int (math.MaxInt  )
        case      *int8:  *x =    int8(math.MaxInt8 )
        case      *int16: *x =   int16(math.MaxInt16)
        case      *int32: *x =   int32(math.MaxInt32)
        case      *int64: *x =   int64(math.MaxInt64)
        case     *uint:   *x =   uint (math.MaxUint  )
        case     *uint8:  *x =   uint8(math.MaxUint8 )
        case     *uint16: *x =  uint16(math.MaxUint16)
        case     *uint32: *x =  uint32(math.MaxUint32)
        case     *uint64: *x =  uint64(math.MaxUint64)
        case    *float32: *x = float32(math.MaxFloat32)
        case    *float64: *x = float64(math.MaxFloat64)
        case  *complex64: *x = complex(math.MaxFloat32, math.MaxFloat32)
        case *complex128: *x = complex(math.MaxFloat64, math.MaxFloat64)
        default: ks.Never()
    }
    return n
}

// Min returns the minimum representable number of type N. By minimum, this
// means the negative number with the greatest magnitude.
//
// For integers, min <= 0 < max, and min != -max.
//
// For floating point numbers, -inf < min < 0 < epsilon < max < inf, and min ==
// -max.
//
// e.g. Max[uint8]() // returns 255
func Min[N Number]() N {
    var n N
    switch x := any(&n).(type) {
        case      *int:   *x =  int (math.MinInt  )
        case      *int8:  *x =  int8(math.MinInt8 )
        case      *int16: *x = int16(math.MinInt16)
        case      *int32: *x = int32(math.MinInt32)
        case      *int64: *x = int64(math.MinInt64)
        case     *uint:   *x = 0
        case     *uint8:  *x = 0
        case     *uint16: *x = 0
        case     *uint32: *x = 0
        case     *uint64: *x = 0
        case    *float32: *x = float32(-math.MaxFloat32)
        case    *float64: *x = float64(-math.MaxFloat64)
        case  *complex64: *x = complex(-math.MaxFloat32, -math.MaxFloat32)
        case *complex128: *x = complex(-math.MaxFloat64, -math.MaxFloat64)
        default: ks.Never()
    }
    return n
}

// RealInfo stores filled-in information about a [Real] number type.
type RealInfo[N Real] struct {
    Min N
    Max      N
    Smallest N
    Signed   bool
}

// Filled-in [RealInfo] information about different [Real] number types.
var (
    Int   = nInt
    Int8  = nInt8
    Int16 = nInt16
    Int32 = nInt32
    Int64 = nInt64

    Uint   = nUint
    Uint8  = nUint8
    Uint16 = nUint16
    Uint32 = nUint32
    Uint64 = nUint64

    Float32 = nFloat32
    Float64 = nFloat64
)

var (
    nInt =        RealInfo[int]{
        Min:      Min[int](),
        Max:      Max[int](),
        Smallest: Smallest[int](),
        Signed:   true,
    }
    nInt8 =       RealInfo[int8]{
        Min:      Min[int8](),
        Max:      Max[int8](),
        Smallest: Smallest[int8](),
        Signed:   true,
    }
    nInt16 =     RealInfo[int16]{
        Min:      Min[int16](),
        Max:      Max[int16](),
        Smallest: Smallest[int16](),
        Signed:   true,
    }
    nInt32 =     RealInfo[int32]{
        Min:      Min[int32](),
        Max:      Max[int32](),
        Smallest: Smallest[int32](),
        Signed:   true,
    }
    nInt64 =     RealInfo[int64]{
        Min:      Min[int64](),
        Max:      Max[int64](),
        Smallest: Smallest[int64](),
        Signed:   true,
    }

    nUint =       RealInfo[uint]{
        Min:      Min[uint](),
        Max:      Max[uint](),
        Smallest: Smallest[uint](),
    }
    nUint8 =      RealInfo[uint8]{
        Min:      Min[uint8](),
        Max:      Max[uint8](),
        Smallest: Smallest[uint8](),
    }
    nUint16 =    RealInfo[uint16]{
        Min:      Min[uint16](),
        Max:      Max[uint16](),
        Smallest: Smallest[uint16](),
    }
    nUint32 =    RealInfo[uint32]{
        Min:      Min[uint32](),
        Max:      Max[uint32](),
        Smallest: Smallest[uint32](),
    }
    nUint64 =    RealInfo[uint64]{
        Min:      Min[uint64](),
        Max:      Max[uint64](),
        Smallest: Smallest[uint64](),
    }

    nFloat32 =  RealInfo[float32]{
        Min:      Min[float32](),
        Max:      Max[float32](),
        Smallest: Smallest[float32](),
        Signed:   true,
    }
    nFloat64 =  RealInfo[float64]{
        Min:      Min[float64](),
        Max:      Max[float64](),
        Smallest: Smallest[float64](),
        Signed:   true,
    }
)
