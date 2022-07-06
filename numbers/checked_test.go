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

        // TODO test min
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
