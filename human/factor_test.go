package human_test

import (
    "testing"

    "github.com/tawesoft/golib/v2/human"
)

func TestFactors_Min(t *testing.T) {
    factors := human.Factors{
        Factors:    []human.Factor{
            { 0.10, human.Unit{"small", "small"},   0},
            { 1.00, human.Unit{"normal", "normal"}, 0},
            {10.00, human.Unit{"big",    "big"},    0},
        },
    }

    type test struct {
        value float64
        expectedFactorIndex int
    }

    tests := []test{
        { 0.01,  0},
        { 0.10,  0},
        { 0.50,  0},
        { 1.00,  1},
        { 1.50,  1},
        {10.00,  2},
        {11.00,  2},
    }

    for _, test := range tests {
        idx := factors.Min(test.value)
        if idx != test.expectedFactorIndex {
            t.Errorf("factors.bracket(%f): got idx %d but expected %d",
                test.value, idx, test.expectedFactorIndex)
        }
    }
}
