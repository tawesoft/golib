package rbnf

import (
    "testing"

    "github.com/tawesoft/golib/v2/must"
)

func TestNew(t *testing.T) {
    // from CLDR 41.0 rbnf/en.xml
    must.Result(New(nil, `
        %spellout-numbering:
            -x: minus →→;
            Inf: infinity;
            NaN: not a number;
            0: =%spellout-cardinal=;
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
        %%and:
            1: and =%spellout-cardinal-verbose=;
            100: =%spellout-cardinal-verbose=;
        %%commas:
            1: and =%spellout-cardinal-verbose=;
            100: =%spellout-cardinal-verbose=;
            1000: ←%spellout-cardinal-verbose← thousand[→%%commas→];
            1000000: =%spellout-cardinal-verbose=;
        %spellout-cardinal-verbose:
            -x: minus →→;
            x.x: ←← point →→;
            Inf: infinite;
            NaN: not a number;
            0: =%spellout-numbering=;
            100: ←← hundred[→%%and→];
            1000: ←← thousand[→%%and→];
            100000/1000: ←← thousand[→%%commas→];
            1000000: ←← million[→%%commas→];
            1000000000: ←← billion[→%%commas→];
            1000000000000: ←← trillion[→%%commas→];
            1000000000000000: ←← quadrillion[→%%commas→];
            1000000000000000000: =#,##0=;
    `))

    // Output:
}
