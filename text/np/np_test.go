package np_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/v2/text/np"
)

func TestGet(t *testing.T) {
    type row struct {
        codepoint rune
        t np.Type
        v np.Fraction
    }

    rows := []row{
        // Not numberals
        {'a', np.None, np.Fraction{}},
        {'X', np.None, np.Fraction{}},
        {'0'-1, np.None, np.Fraction{}},
        {'9'+1, np.None, np.Fraction{}},

        // ASCII Latin
        {'0', np.Decimal, np.Fraction{0, 1}},
        {'5', np.Decimal, np.Fraction{5, 1}},
        {'9', np.Decimal, np.Fraction{9, 1}},

        // Other Decimal
        {'६', np.Decimal, np.Fraction{6, 1}}, // Devanagari
        {'೬', np.Decimal, np.Fraction{6, 1}}, // Kannada
        {'٤', np.Decimal, np.Fraction{4, 1}}, // Arabic-Indic
        {'۴', np.Decimal, np.Fraction{4, 1}}, // Extended Arabic-Indic

        // Typographic
        {'⑥', np.Digit, np.Fraction{6, 1}},
        {'⑨', np.Digit, np.Fraction{9, 1}},

        // Fractions
        {'¾', np.Numeric, np.Fraction{3, 4}},
        {'⅐', np.Numeric, np.Fraction{1, 7}},
        {'↉', np.Numeric, np.Fraction{0, 1}}, // Unicode treats this as 0/1, not 0/3

        // Roman numerals
        {'Ⅰ', np.Numeric, np.Fraction{1, 1}},
        {'Ⅱ', np.Numeric, np.Fraction{2, 1}},
        {'Ⅹ', np.Numeric, np.Fraction{10, 1}},
        {'Ⅻ', np.Numeric, np.Fraction{12, 1}},
        {'Ⅽ', np.Numeric, np.Fraction{100, 1}},

        // Tibet
        {'༭', np.Numeric, np.Fraction{7, 2}}, // TIBETAN DIGIT HALF FOUR
        {'༳', np.Numeric, np.Fraction{-1, 2}}, // TIBETAN DIGIT HALF ZERO

        // Tamil
        {'௰', np.Numeric, np.Fraction{10, 1}},
        {'௱', np.Numeric, np.Fraction{100, 1}},
        {'௲', np.Numeric, np.Fraction{1000, 1}},

        // Other Numeral
        {'𝍪', np.Numeric, np.Fraction{20, 1}}, // Counting Rods
        {'𑁜', np.Numeric, np.Fraction{20, 1}}, // Brahmi
        {'𒑡', np.Numeric, np.Fraction{1, 6}}, // Cuneiform
        {'𐳿', np.Numeric, np.Fraction{1000, 1}}, // Old Hungarian
        {'𖭡', np.Numeric, np.Fraction{1000000000000, 1}}, // Pahawh Hmong
    }

    for i, r := range rows {
        ty, v := np.Get(r.codepoint)
        assert.Equal(t, r.t, ty, "type for test %d", i)
        assert.Equal(t, r.v, v, "value for test %d", i)
    }
}

/*
func Foo() {
    for _, s := range spans {
        fmt.Printf("%c %x: %d, %d/%d (%d)\n",
            s.codepoint, s.codepoint, s.nt, s.nv.Numerator, s.nv.Denominator, s.length)
    }
}
 */
