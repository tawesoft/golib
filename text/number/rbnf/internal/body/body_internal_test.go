package body

import (
    "testing"

    "github.com/tawesoft/golib/v2/must"
)

func Test_consumeDirect(t *testing.T) {
    type row struct {
        fn string
        startOffset int
        input string
        expectedValue string
        expectedType Type
        expectedOffset int
    }
    rows := []row{
        {"Literal", len("X"), "X",           "",          TypeLiteral,  0},
        {"Literal", len("X"), "X£",          "£",         TypeLiteral,  0},
        {"Literal", len("X"), "Xmillion",    "million",   TypeLiteral,  0},
        {"Literal", len("X"), "Xmillion $",  "million $", TypeLiteral,  0},
        {"Literal", len("X"), "Xmillion $(", "million ",  TypeLiteral, -2},
        {"Literal", len("X"), "Xmillion =",  "million ",  TypeLiteral, -1},

        {"SimpleSubstitution", len("X→"), "X→→ X",     "",   TypeSubstRightArrow, -2},
        {"SimpleSubstitution", len("X→"), "X→foo→ X", "foo", TypeSubstRightArrow, -2},

        {"PluralSubstitution", len("$("), "$(ordinal,foo)$ X",  "foo", TypeSubstPluralOrdinal, -2},
        {"PluralSubstitution", len("$("), "$(cardinal,foo)$ X", "foo", TypeSubstPluralCardinal,-2},
    }
    fns := map[string]func(string, int) (Token, int){
        "Literal": consumeLiteral,
        "SimpleSubstitution": func(x string, n int) (Token, int) {
            return consumeSimpleSubstitution(x, n, '→')
        },
        "PluralSubstitution": consumePluralSubstitution,
    }
    for _, test := range rows {
        f, ok := fns[test.fn]; must.True(ok, "no function %q", test.fn)
        got, offset := f(test.input, test.startOffset)
        gotValue := test.input[got.Content[0]:got.Content[1]]

        if (gotValue != test.expectedValue) ||
            (got.Type != test.expectedType) ||
            (len(test.input) + test.expectedOffset != offset) {
            t.Errorf("consume%s(%q[%d:]): wanted %d %q +%d but got %d %q +%d",
                test.fn, test.input, test.startOffset,
                test.expectedType, test.expectedValue,
                len(test.input) + test.expectedOffset,
                got.Type, gotValue, offset)
        }
    }
}

func Test_tokenizerDirect(t *testing.T) {
    type row struct {
        input string
        expected []Token
        substTypes []SubstType
    }
    rows := []row{
        {"", []Token{{Type: TypeEOF}}, []SubstType{SubstTypeNone}},
        {
            "a[b =%c=→→]d←#e←==", // note 3 bytes per '→' or '←'
            []Token{
                {
                    Type:       TypeLiteral,
                    Content:    Slice{0, 1}, // "a"
                },
                {
                    Type:       TypeOptionalStart,
                },
                {
                    Type:       TypeLiteral,
                    Content:    Slice{2, 4}, // "b "
                },
                {
                    Type:       TypeSubstEqualsSign,
                    Content:    Slice{5, 7}, // "%c"
                },
                {
                    Type:       TypeSubstRightArrow,
                    Content:    Slice{11, 11}, // ""
                },
                {
                    Type:       TypeOptionalEnd,
                },
                {
                    Type:       TypeLiteral,
                    Content:    Slice{15, 16}, // "d"
                },
                {
                    Type:       TypeSubstLeftArrow,
                    Content:    Slice{19, 21}, // "#e"
                },
                {
                    Type:       TypeSubstEqualsSign,
                    Content:    Slice{25, 25}, // ""
                },
                {
                    Type:       TypeEOF,
                },
            },
            []SubstType{
                SubstTypeNone,
                SubstTypeNone,
                SubstTypeNone,
                SubstTypeRulesetName,
                SubstTypeEmpty,
                SubstTypeNone,
                SubstTypeNone,
                SubstTypeDecimalFormat,
                SubstTypeInvalid,
                SubstTypeNone,
            },
        },
    }

    for _, test := range rows {
        tokenizer := tokenizer{str: test.input}
        for i := 0; i < len(test.expected); i++ {
            tok := tokenizer.next()

            if tok != test.expected[i] {
                t.Errorf("tokenizing %q: token %d: got %+v but expected %+v (%q)",
                    test.input, i, tok, test.expected[i],
                    test.input[test.expected[i].Content[0]:test.expected[i].Content[1]])
            }

            if tok.SimpleSubstType(test.input) != test.substTypes[i] {
                t.Errorf("tokenizing %q: token %d: got SubstType %d but expected %d",
                    test.input, i, tok.SimpleSubstType(test.input), test.substTypes[i])
            }

            if (tok.Type == TypeEOF) {
                if (i + 1) < len(test.expected) {
                    t.Errorf("tokenizing %q ended early after %d/%d tokens",
                        test.input, i + 1, len(test.expected))
                }
                break
            }
        }
    }
}
