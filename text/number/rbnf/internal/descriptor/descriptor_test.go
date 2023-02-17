package descriptor

import (
    "testing"
)

func TestParseRuleDescriptor(t *testing.T) {
    type row struct {
        input string
        r Descriptor
    }
    rows := []row{
        {input: "100000",      r: Descriptor{Type: TypeBaseValue,         Base: 100_000, Divisor: 0}},
        {input: "100000/1000", r: Descriptor{Type: TypeBaseValueAndRadix, Base: 100_000, Divisor: 1000}},
    }
    for _, test := range rows {
        got := Parse(test.input)
        if (got != test.r) {
            t.Errorf("parse(%q): got %+v but expected %+v",
                test.input, got, test.r)
        }
    }
}
