package color_test

import (
    "testing"

    "github.com/tawesoft/golib/v2/css/color"
)

// TODO fuzz round trip testing

func TestParseColor(t *testing.T) {
    type row struct {
        input string
        specified string
        computed string
        ok bool
    }
    rows := []row{
        // clamping
        {
            input:     "rgb(512 -64 32)",
            specified: "rgb(512, -64, 32)",
            computed:  "rgb(255, 0, 32)",
            ok: true,
        },

        // rgb plus alpha => rgba legacy representation
        {
            input:     "rgb(128 64 32 / 0.5)",
            specified: "rgba(128, 64, 32, 0.5)",
            computed:  "rgba(128, 64, 32, 0.5)",
            ok: true,
        },

        // legacy representation parsing
        {
            input:     "rgb(128, 64, 32, 50%)",
            specified: "rgba(128, 64, 32, 0.5)",
            computed:  "rgba(128, 64, 32, 0.5)",
            ok: true,
        },

        // percentages
        {
            input:     "rgb(57.28% 42.14% 51.46% / 25%)",
            specified: "rgba(146.064, 107.457, 131.223, 0.25)",
            computed:  "rgba(146.064, 107.457, 131.223, 0.25)",
            ok: true,
        },

        // strip trailing zeroes
        {
            input:     "rgb(128 0064 010.03500)",
            specified: "rgb(128, 64, 10.035)",
            computed:  "rgb(128, 64, 10.035)",
            ok: true,
        },

        // "none" keyword including case folding
        {
            input:     "RGB(none NoNe none / none)",
            specified: "rgb(0, 0, 0)",
            computed:  "rgb(0, 0, 0)",
            ok: true,
        },

        // hexadecimal representation
        {
            input:     "#FA7", // 3
            specified: "#ffaa77",
            computed:  "rgb(255, 170, 119)",
            ok: true,
        },
        {
            input:     "#FA73", // 4
            specified: "#ffaa7733",
            computed:  "rgba(255, 170, 119, 0.2)",
            ok: true,
        },
        {
            input:     "#Fba57a", // 6
            specified: "#fba57a",
            computed:  "rgb(251, 165, 122)",
            ok: true,
        },
        {
            input:     "#Fba57a33", // 8
            specified: "#fba57a33",
            computed:  "rgba(251, 165, 122, 0.2)",
            ok: true,
        },
        {
            input:     "#B6002F", // with zeros
            specified: "#b6002f",
            computed:  "rgb(182, 0, 47)",
            ok: true,
        },
    }

    for _, r := range rows {
        specified, err := color.ParseColorString(r.input)
        computed := specified.Norm()

        if r.ok != (err == nil) {
            t.Errorf("expected ok=%v for input %s but got (%v, %v)",
                r.ok, r.input, specified, err)
            continue
        }
        if specified.String() != r.specified {
            t.Errorf("expected specified %s but got %s on input %q",
                r.specified, specified, r.input)
        }
        if computed.String() != r.computed {
            t.Errorf("expected computed %s but got %s on input %q",
                r.computed, computed, r.input)
        }
    }
}
