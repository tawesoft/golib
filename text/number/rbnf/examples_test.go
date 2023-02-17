package rbnf_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/text/number/rbnf"
)

func Example_spelloutCardinal() {
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

// Example using custom time factors from the Battlestar Galactica 1978 TV
// series.
func Example_fictional() {
    /*
            168000: ←← secton [ →→];
            672000: ←← quatron [ →→];
            80640000: =%%spellout-cardinal-verbose= yahren [ >>];
    161280000/161280000: ←%%spellout-cardinal-verbose← yahrens[ →→→];
     */
    g := must.Result(rbnf.New(nil, `
        %%s:
            0: s;
            1: ;
            2: s;
        %%es:
            0: es;
            1: ;
            2: es;
        %%timecomma:
            0: =%time=;
            1: , =%time=;
        %%microns:
            0: =%%spellout-cardinal= microns;
            1: =%%spellout-cardinal= micron;
            2: =%%spellout-cardinal= microns;
        %%hyphen-microns:
            0: ' microns;
            1: -=%%spellout-cardinal= micron;
            2: -=%%spellout-cardinal= microns;
        %time:
            -x: minus →→;
            0: =%%microns=;
            1: =%%microns=;
            2: =%%microns=;
            20: twenty→%%hyphen-microns→;
            30: thirty→%%hyphen-microns→;
            40: forty→%%hyphen-microns→;
            50: fifty→%%hyphen-microns→;
            60: sixty→%%hyphen-microns→;
            70: seventy→%%hyphen-microns→;
            80: eighty→%%hyphen-microns→;
            90: ninety→%%hyphen-microns→;
            100: ←%%spellout-cardinal← centon[←%%s←→%%timecomma→];
            6000/6000: ←%%spellout-cardinal← centar[←%%es←→%%timecomma→];
            144000/144000: ←%%spellout-cardinal← cycle[←%%s←→%%timecomma→];
            1008000/1008000: ←%%spellout-cardinal← secton[←%%s←→%%timecomma→];
            4032000/4032000: ←%%spellout-cardinal← quatron[←%%s←→%%timecomma→];
            48384000/48384000: ←%%spellout-cardinal-verbose← yahren[←%%s←→%%timecomma→];
        %%spellout-cardinal:
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
        %%spellout-cardinal-verbose:
            0: =%%spellout-numbering=;
            100: ←← hundred[→%%and→];
            1000: ←← thousand[→%%and→];
            100000/1000: ←← thousand[→%%commas→];
            1000000: ←← million[→%%commas→];
            1000000000: ←← billion[→%%commas→];
            1000000000000: ←← trillion[→%%commas→];
            1000000000000000: ←← quadrillion[→%%commas→];
            1000000000000000000: =#,##0=;
        %%spellout-numbering:
            0: =%%spellout-cardinal=;
        %%and:
            1: ' and =%%spellout-cardinal-verbose=;
            100: ' =%%spellout-cardinal-verbose=;
        %%commas:
            1:' and =%%spellout-cardinal-verbose=;
            100: ' =%%spellout-cardinal-verbose=;
            1000: ' ←%%spellout-cardinal-verbose← thousand[→%%commas→];
            1000000: ' =%%spellout-cardinal-verbose=;
    `))

    type microns int64

    printTime := func(v microns) {
        fmt.Printf("printTime(microns(%d)): %s\n", v,
            must.Result(g.FormatInteger("%time", int64(v))))
    }

    const centon  = 100           // in microns, ~= 1 minute.
    const centar  = 60 * centon   // ~= 1 hour, plural "centares"
    const cycle   = 24 * centar   // ~= 1 day
    const secton  = 7 * cycle    // ~= 1 week
    const quatron = 4 * secton   // ~= 1 month
    const yahren  = 12 * quatron // ~= 1 year

    printTime(microns(0))
    printTime(microns(1))
    printTime(microns(5))
    printTime(microns(1*centar))
    printTime(microns(2*centar))
    printTime(microns((1*centon) + 95))
    printTime(microns((2*centar)+ (5*centon) + 1))
    printTime(microns((1*cycle) + (1*centar) + 5))
    printTime(microns(1*secton))
    printTime(microns(1*quatron))
    printTime(microns((3*quatron) + (2*secton)))
    printTime(microns(1*yahren))
    printTime(microns(2*yahren))
    printTime(microns(150*yahren))
    printTime(microns((101*yahren) + (6*quatron) + (3*secton) + (4*cycle) + (2*centar) + 50))

    // Output:
    // printTime(microns(0)): zero microns
    // printTime(microns(1)): one micron
    // printTime(microns(5)): five microns
    // printTime(microns(1)): one centar
    // printTime(microns(2)): one centares
    // printTime(microns(195)): one centon, ninety-five microns
    // printTime(microns(12501)): two centares, five centons, one micron
    // printTime(microns(150005)): one cycle, one centar, five microns
    // printTime(microns(1008000)): one secton
    // printTime(microns(4032000)): one quatron
    // printTime(microns(14112000)): three quatrons, two secton
    // printTime(microns(48384000)): one yahren
    // printTime(microns(96768000)): two yahren
    // printTime(microns(7257600000)): one hundred and fifty yahren
    // printTime(microns(4914588050)): one hundred and one yahrens, six quatrons, three sectons, four cycles, two centares, fifty microns
}
