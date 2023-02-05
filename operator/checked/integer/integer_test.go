package integer_test

import (
    "math"
    "testing"

    "github.com/tawesoft/golib/v2/operator/checked/integer"
)

func FuzzAdd_Int32(f *testing.F) {
    type row struct {
        a int32
        b int32
        min int32
        max int32
    }

    const halfmax = math.MaxInt32/2
    const halfmin = math.MinInt32/2

    rows := []row{
        {-1, -1, math.MinInt32, math.MaxInt32},
        {0, 0, math.MinInt32, math.MaxInt32},
        {1, 1, math.MinInt32, math.MaxInt32},
        {-1, -1, math.MinInt32, math.MaxInt32},
        {0, 0, 0, math.MaxInt32},
        {halfmax, halfmin, 0, math.MaxInt32},
        {halfmax, halfmin + 4, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32-1},
    }

    for _, tc := range rows {
        f.Add(tc.min, tc.max, tc.a, tc.b)
    }

    f.Fuzz(func(t *testing.T, min, max, a, b int32) {
        if ((min > max) ||
            (a < min) ||
            (a > max) ||
            (b < min) ||
            (b > max)) {
            t.SkipNow()
        }

        result, resultOk := integer.Add(min, max, a, b)
        expectedResult := int64(a) + int64(b)
        ok := (expectedResult <= int64(max)) && (expectedResult >= int64(min))

        if ok != resultOk {
            t.Errorf("integer.Add(%v, %v, %v, %v): got %t, expected %t",
                min, max, a, b,
                resultOk, ok)
        } else if ok && (int64(result) != expectedResult) {
            t.Errorf("integer.Add(%v, %v, %v, %v): got %v, %t, expected %v",
                min, max, a, b,
                result, resultOk, expectedResult)
        }
    })
}

func FuzzAdd_Uint32(f *testing.F) {
    type row struct {
        a uint32
        b uint32
        min uint32
        max uint32
    }

    const halfmax = math.MaxUint32/2

    rows := []row{
        {0, 0, 0, math.MaxInt32},
        {0, 1, 0, math.MaxInt32},
        {1, 0, 0, math.MaxInt32},
        {1, 1, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32-1},
    }

    for _, tc := range rows {
        f.Add(tc.min, tc.max, tc.a, tc.b)
    }

    f.Fuzz(func(t *testing.T, min, max, a, b uint32) {
        if ((min > max) ||
            (a < min) ||
            (a > max) ||
            (b < min) ||
            (b > max)) {
            t.SkipNow()
        }

        result, resultOk := integer.Add(min, max, a, b)
        expectedResult := int64(a) + int64(b)
        ok := (expectedResult <= int64(max)) && (expectedResult >= int64(min))

        if ok != resultOk {
            t.Errorf("integer.Add(%v, %v, %v, %v): got %t, expected %t",
                min, max, a, b,
                resultOk, ok)
        } else if ok && (int64(result) != expectedResult) {
            t.Errorf("integer.Add(%v, %v, %v, %v): got %v, %t, expected %v",
                min, max, a, b,
                result, resultOk, expectedResult)
        }
    })
}

func FuzzSub_Int32(f *testing.F) {
    type row struct {
        a int32
        b int32
        min int32
        max int32
    }

    const halfmax = math.MaxInt32/2
    const halfmin = math.MinInt32/2

    rows := []row{
        {-1, -1, math.MinInt32, math.MaxInt32},
        {0, 0, math.MinInt32, math.MaxInt32},
        {1, 1, math.MinInt32, math.MaxInt32},
        {-1, -1, math.MinInt32, math.MaxInt32},
        {0, 0, 0, math.MaxInt32},
        {halfmax, halfmin, 0, math.MaxInt32},
        {halfmax, halfmin + 4, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32-1},
    }

    for _, tc := range rows {
        f.Add(tc.min, tc.max, tc.a, tc.b)
    }

    f.Fuzz(func(t *testing.T, min, max, a, b int32) {
        if ((min > max) ||
            (a < min) ||
            (a > max) ||
            (b < min) ||
            (b > max)) {
            t.SkipNow()
        }

        result, resultOk := integer.Sub(min, max, a, b)
        expectedResult := int64(a) - int64(b)
        ok := (expectedResult <= int64(max)) && (expectedResult >= int64(min))

        if ok != resultOk {
            t.Errorf("integer.Sub(%v, %v, %v, %v): got %t, expected %t",
                min, max, a, b,
                resultOk, ok)
        } else if ok && (int64(result) != expectedResult) {
            t.Errorf("integer.Sub(%v, %v, %v, %v): got %v, %t, expected %v",
                min, max, a, b,
                result, resultOk, expectedResult)
        }
    })
}

func FuzzSub_Uint32(f *testing.F) {
    type row struct {
        a uint32
        b uint32
        min uint32
        max uint32
    }

    const halfmax = math.MaxInt32/2

    rows := []row{
        {0, 0, 0, math.MaxInt32},
        {0, 1, 0, math.MaxInt32},
        {1, 0, 0, math.MaxInt32},
        {1, 1, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32},
        {halfmax, halfmax + 1, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32-1},
    }

    for _, tc := range rows {
        f.Add(tc.min, tc.max, tc.a, tc.b)
    }

    f.Fuzz(func(t *testing.T, min, max, a, b uint32) {
        if ((min > max) ||
            (a < min) ||
            (a > max) ||
            (b < min) ||
            (b > max)) {
            t.SkipNow()
        }

        result, resultOk := integer.Sub(min, max, a, b)
        expectedResult := int64(a) - int64(b)
        ok := (expectedResult <= int64(max)) && (expectedResult >= int64(min))

        if ok != resultOk {
            t.Errorf("integer.Sub(%v, %v, %v, %v): got %t, expected %t",
                min, max, a, b,
                resultOk, ok)
        } else if ok && (int64(result) != expectedResult) {
            t.Errorf("integer.Sub(%v, %v, %v, %v): got %v, %t, expected %v",
                min, max, a, b,
                result, resultOk, expectedResult)
        }
    })
}

func FuzzMul_Int32(f *testing.F) {
    type row struct {
        a int32
        b int32
        min int32
        max int32
    }

    const halfmax = math.MaxInt32/2
    const halfmin = math.MinInt32/2

    rows := []row{
        {-1, -1, math.MinInt32, math.MaxInt32},
        {0, 0, math.MinInt32, math.MaxInt32},
        {1, 1, math.MinInt32, math.MaxInt32},
        {-1, -1, math.MinInt32, math.MaxInt32},
        {0, 0, 0, math.MaxInt32},
        {halfmax, halfmin, 0, math.MaxInt32},
        {halfmax, halfmin + 4, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32-1},
    }

    for _, tc := range rows {
        f.Add(tc.min, tc.max, tc.a, tc.b)
    }

    f.Fuzz(func(t *testing.T, min, max, a, b int32) {
        if ((min > max) ||
            (a < min) ||
            (a > max) ||
            (b < min) ||
            (b > max)) {
            t.SkipNow()
        }

        result, resultOk := integer.Mul(min, max, a, b)
        expectedResult := int64(a) * int64(b)
        ok := (expectedResult <= int64(max)) && (expectedResult >= int64(min))

        if ok != resultOk {
            t.Errorf("integer.Mul(%v, %v, %v, %v): got %t, expected %t",
                min, max, a, b,
                resultOk, ok)
        } else if ok && (int64(result) != expectedResult) {
            t.Errorf("integer.Mul(%v, %v, %v, %v): got %v, %t, expected %v",
                min, max, a, b,
                result, resultOk, expectedResult)
        }
    })
}

func FuzzMul_Uint32(f *testing.F) {
    type row struct {
        a uint32
        b uint32
        min uint32
        max uint32
    }

    const halfmax = math.MaxInt32/2

    rows := []row{
        {0, 0, 0, math.MaxInt32},
        {0, 1, 0, math.MaxInt32},
        {1, 0, 0, math.MaxInt32},
        {1, 1, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32},
        {halfmax, halfmax + 1, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32},
        {halfmax, halfmax, 0, math.MaxInt32-1},
    }

    for _, tc := range rows {
        f.Add(tc.min, tc.max, tc.a, tc.b)
    }

    f.Fuzz(func(t *testing.T, min, max, a, b uint32) {
        if ((min > max) ||
            (a < min) ||
            (a > max) ||
            (b < min) ||
            (b > max)) {
            t.SkipNow()
        }

        result, resultOk := integer.Mul(min, max, a, b)
        expectedResult := int64(a) * int64(b)
        ok := (expectedResult <= int64(max)) && (expectedResult >= int64(min))

        if ok != resultOk {
            t.Errorf("integer.Mul(%v, %v, %v, %v): got %t, expected %t",
                min, max, a, b,
                resultOk, ok)
        } else if ok && (int64(result) != expectedResult) {
            t.Errorf("integer.Mul(%v, %v, %v, %v): got %v, %t, expected %v",
                min, max, a, b,
                result, resultOk, expectedResult)
        }
    })
}
