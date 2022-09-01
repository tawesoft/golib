package fold_test

import (
    "io"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/v2/text/fold"
    "golang.org/x/text/transform"
)

func trans(t transform.Transformer, x string) string {
    r := transform.NewReader(strings.NewReader(x), t)
    bs, err := io.ReadAll(r)
    s := string(bs)
    if err != nil { s = "error: " + err.Error() }
    return s
}

func Test(t *testing.T) {
    type row struct {
        t transform.Transformer
        input string
        expected string
    }

    rows := []row{
        {fold.Accents,              "",             ""},        // same
        {fold.Accents,              "café",         "cafe"},    // é => e
        {fold.Accents,              "ёёёё",         "ееее"},    // ё => Cyrillic Small Letter Ie

        {fold.CanonicalDuplicates,  "",             ""},        // same
        {fold.CanonicalDuplicates,  "café",         "café"},    // same
        {fold.CanonicalDuplicates,  "aΩaé",         "aΩaé"},    // Ohm => Omega

        {fold.Dashes,               "",             ""},        // same
        {fold.Dashes,               "---",          "---"},     // same
        {fold.Dashes,               "a-b-c",        "a-b-c"},   // same
        {fold.Dashes,               "a\u2011b\u2010c", "a-b-c"},   // non-breaking hyphen, hyphen, to hyphen-minus
        {fold.Dashes,               "a⸺b⸺c",     "a-b-c"},  // to hyphen-minus

        {fold.Digit,                "",             ""},           // same
        {fold.Digit,                "abcdef",       "abcdef"},     // same
        {fold.Digit,                "0123456789",   "0123456789"}, // same
        {fold.Digit,                "٠١٢٣٤٥٦٧٨٩",   "0123456789"},
        {fold.Digit,                "۰۱۲۳۴۵۶۷۸۹",   "0123456789"},
        {fold.Digit,                "⓪①②③④⑤⑥⑦⑧⑨",   "0123456789"},
        {fold.Digit,                "⁵₅",           "55"},


        {fold.GreekLetterforms,     "",             ""},        // same
        {fold.GreekLetterforms,     "café",         "café"},    // same
        {fold.GreekLetterforms,     "ϐϑϒ",          "βθΥ"},

        {fold.HebrewAlternates,     "",             ""},        // same
        {fold.GreekLetterforms,     "café",         "café"},    // same
        {fold.HebrewAlternates,     "ﬨ",            "ת"},       // Hebrew Letter Wide Tav => Hebrew Letter Tav

        {fold.Jamo,                 "",             ""},        // same
        {fold.Jamo,                 "café",         "café"},    // same
        {fold.Jamo,                 "ㆃ",           "ᇲ"},

        {fold.Math,                 "",             ""},        // same
        {fold.Math,                 "café",         "café"},    // same
        {fold.Math,                 "𝛑",            "π"},       // Mathematical Bold Small Pi => Greek Small Letter Pi

        {fold.NoBreak,              "",             ""},        // same
        {fold.NoBreak,              "café",         "café"},    // same
        {fold.NoBreak,              "a\u00A0b",     "a b"},     // nbsp => space
        {fold.NoBreak,              "a\u202Fb",     "a b"},     // nnbsp => space
        {fold.NoBreak,              "a\u2011b",     "a\u2010b"}, // non-breaking hyphen => hyphen

        // TODO tests for fold.Positional

        {fold.Space,                "",             ""},        // Same
        {fold.Space,                "café",         "café"},    // Same
        {fold.Space,                "\t",           "\t"},      // Same - \t is control, not space
        {fold.Space,                "a\u00A0b",     "a b"},     // nbsp => space
        {fold.Space,                "a\u205Fb",     "a b"},     // Medium mathematical space
        {fold.Space,                "\u2800",       "\u2800"},  // Same - Unicode says Braille blank does not act as a space
        {fold.Space,                "\u3000",       " "},       // Ideographic space

        {fold.Small,                "",             ""},        // same
        {fold.Small,                "café",         "café"},    // same
        {fold.Small,                "f",            "f"},       // small f => regular f
    }

    for i, r := range rows {
        output := trans(r.t, r.input)
        assert.Equal(t, r.expected, output, "test %d on input %q", i, r.input)
    }

}
