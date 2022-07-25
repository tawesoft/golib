package fallback

import (
    "fmt"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/v2/ks"
    "github.com/tawesoft/golib/v2/lazy"
    "golang.org/x/text/unicode/norm"
)

func TestSegments(t *testing.T) {
    type row struct {
        input string
        expected string
    }

    rows := []row{
        {"",        ""},
        {"a",       "a"},
        {"aa",      "a�a"},
        {"aaa",     "a�a�a"},
        {"é",       "e\u0301"},
        {"éé",      "e\u0301�e\u0301"},
        {"\u0041\u030A\u0064\u0307\u0327", "\u0041\u030A�\u0064\u0307\u0327"},
    }

    for _, r := range rows {
        input := norm.NFD.String(r.input)
        expected := norm.NFD.String(r.expected)
        output := segments(input)
        assert.Equal(t, expected, output,
            "segs(%v) => %v, wanted %v",
            []rune(input), []rune(output), []rune(expected))
    }
}

func TestStringGet(t *testing.T) {
    assert.Equal(t, "a", stringGet("a�b�c", 0))
    assert.Equal(t, "b", stringGet("a�b�c", 1))
    assert.Equal(t, "c", stringGet("a�b�c", 2))
}

func Example_combinations() {
    input := []string{
        "a�b�c",
        "d�e�f",
        "w",
        "x�y�z",
    }

    it := ks.Must(combinations(input))
    lazy.Walk(func (x string) {
        fmt.Println(x)
    }, it)

    // output:
    // adwx
    // bdwx
    // cdwx
    // aewx
    // bewx
    // cewx
    // afwx
    // bfwx
    // cfwx
    // adwy
    // bdwy
    // cdwy
    // aewy
    // bewy
    // cewy
    // afwy
    // bfwy
    // cfwy
    // adwz
    // bdwz
    // cdwz
    // aewz
    // bewz
    // cewz
    // afwz
    // bfwz
    // cfwz
}

func Example_dstarts() {
    a := dstarts('a')
    _ = dstarts(0x2A600) // last item

    for _, r := range(string(a)) {
        fmt.Printf("%c\n", r)
    }

    // output:
    // à
    // á
    // â
    // ã
    // ä
    // å
    // ā
    // ă
    // ą
    // ǎ
    // ȁ
    // ȃ
    // ȧ
    // ḁ
    // ạ
    // ả
}
