package numbers_test

import (
    "math"
    "testing"

    "github.com/tawesoft/golib/v2/numbers"
)

func TestCheckedAddReals(t *testing.T) {
    min := math.MinInt
    max := math.MaxInt
    yes := 1
    no := 0

    rows := [][4]int{
        //     a,       b, expected, ok?
        {      1,       2,       3, yes},

        {max - 1,       1, max    , yes},
        {max    ,       1,       0,  no},
        {      1, max    ,       0,  no},
        {      1, max - 1, max    , yes},

        {max - 1,       1, max    , yes},
        {max    ,       1,       0,  no},
        {      1, max    ,       0,  no},
        {      1, max - 1, max    , yes},

    }

    for i, r := range rows {
        a, b, expectedSum := r[0], r[1], r[2]
        expectedOk := r[3] == yes
        actual, ok := numbers.CheckedAddReals(min, max, a, b)

        if ok == expectedOk {
            if actual != expectedSum {
                t.Errorf("%d (%v): expected sum %d but got %d", i, r, expectedSum, actual)
            }
        } else {
            t.Errorf("%d (%v): expected ok %t but got %t", i, r, expectedOk, ok)
        }
    }
}

func TestCheckedMulReals(t *testing.T) {
    min := math.MinInt16
    max := math.MaxInt16
    half := (math.MaxInt16/2) + 10 // half plus a bit
    mini := (math.MaxInt16/2) - 10 // half minus a bit
    yes := 1
    no := 0

    rows := [][4]int{
        //     a,       b, expected, ok?
        { 3,     4,        12, yes}, //  0

        { 0,   max,         0, yes}, //  1
        { 0,   min,         0, yes}, //  2
        { 1,   max,       max, yes}, //  3
        { 1,   min,       min, yes}, //  4
        {-1,   max,      -max, yes}, //  5
        {-1,   min,         0,  no}, //  6

        { 2,  half,         0,  no}, //  7
        { 2, -half,         0,  no}, //  8
        {-2,  half,         0,  no}, //  9
        {-2, -half,         0,  no}, // 10

        { 2,  mini,  2 * mini, yes}, // 11
        { 2, -mini, -2 * mini, yes}, // 12
        {-2,  mini, -2 * mini, yes}, // 13
        {-2, -mini,  2 * mini, yes}, // 14
    }

    for i, r := range rows {
        for j := 0; j < 2; j++ {
            a, b, expectedSum := r[0], r[1], r[2]

            if j == 1 {
                a, b = b, a // try reversed arguments, too
            }

            expectedOk := r[3] == yes
            actual, ok := numbers.CheckedMulReals(min, max, a, b)

            if ok == expectedOk {
                if actual != expectedSum {
                    t.Errorf("%d, %d (%v): expected result %d but got %d", i, j, r, expectedSum, actual)
                }
            } else {
                t.Errorf("%d, %d (%v): expected ok %t but got %t", i, j, r, expectedOk, ok)
            }
        }
    }
}
