package rbnf_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/text/number/rbnf"
)

func ExampleGroup_FormatInteger() {
    g := must.Result(rbnf.New(nil, `
        %spellout-cardinal:
            -x: minus →→;
            x.x: ←← point →→;
            Inf: infinite;
            NaN: not a number;
            0: zero;
            1: one;
            2: two;
            3: three;
            4: four;
            5: five;
            6: six;
            7: seven;
            8: eight;
            9: nine;
            10: ten;
            11: eleven;
            12: twelve;
            13: thirteen;
            14: fourteen;
            15: fifteen;
            16: sixteen;
            17: seventeen;
            18: eighteen;
            19: nineteen;
            20: twenty[-→→];
            30: thirty[-→→];
            40: forty[-→→];
            50: fifty[-→→];
            60: sixty[-→→];
            70: seventy[-→→];
            80: eighty[-→→];
            90: ninety[-→→];
            100: ←← hundred[ →→];
            1000: ←← thousand[ →→];
            1000000: ←← million[ →→];
            1000000000: ←← billion[ →→];
            1000000000000: ←← trillion[ →→];
            1000000000000000: ←← quadrillion[ →→];
            1000000000000000000: =#,##0=;
    `))

    spellout := func(x int64) {
        fmt.Printf("spellout(%d): %s\n", x,
            must.Result(g.FormatInteger("%spellout-cardinal", x)))
    }

    spellout(0)
    spellout(1)
    spellout(2)
    spellout(-5)
    spellout(25)
    spellout(-325)

    // Output:
    // spellout(0): zero
    // spellout(1): one
    // spellout(2): two
    // spellout(-5): minus five
    // spellout(25): twenty-five
    // spellout(-325): minus three hundred twenty-five
}
