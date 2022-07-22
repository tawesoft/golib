package ccc_test

import (
    "bytes"
    "io"
    "strings"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/v2/ks"
    "github.com/tawesoft/golib/v2/text/ccc"
    "golang.org/x/text/transform"
)

func TestOf(t *testing.T) {
    type row struct {
        codepoint rune
        ccc ccc.CCC
    }

    rows := []row{
        // tests at range boundaries
        {0x0299,   0}, // Latin Letter Small Capital B, Not Reordered
        {0x0300, 230}, // Combining Grave Accent, Above
        {0x0314, 230}, // Combining Reversed Comma Above, Above
        {0x0315, 232}, // Combining Comma Above Right, Above Right
        {0x0334,   1}, // Combining Tilde Overlay, Overlay

        // test single-codepoint range
        {0x0591, 220}, // Hebrew Accent Etnahta, Below
        {0x1E94A,  7}, // Adlam Nukta, Nukta (this is also the last entry)

        // test values with none recorded
        {'a',     0}, // Not Reordered (start of ranges)
        {0x1E900, 0}, // Adlam Capital Letter Alif, Not Reordered (between ranges)
        {0x1E94B, 0},  // Adlam Nasalization Mark, Not Reordered (end of ranges)
    }

    for i, r := range rows {
        ccc := ccc.Of(r.codepoint)
        assert.Equal(t, r.ccc, ccc, "for test %d of codepoint %x", i, r.codepoint)
    }
}

func TestReorder(t *testing.T) {

    // This is tested more thoroughly by the tests in text/dm

    type row struct {
        input []rune
        output []rune
    }

    rows := []row{
        {
            []rune{0x0064, 0x0307, 0x0323},
            []rune{0x0064, 0x0323, 0x0307},
        },
        {
            []rune{0x0064, 0x0064, 0x0064, 0x0307, 0x0307, 0x0307, 0x0307, 0x0323, 0x0064,},
            []rune{0x0064, 0x0064, 0x0064, 0x0323, 0x0307, 0x0307, 0x0307, 0x0307, 0x0064,},
        },
        {
            []rune{0x0064, 0x0064, 0x0064, 0x0307, 0x0307, 0x0307, 0x0307, 0x0323},
            []rune{0x0064, 0x0064, 0x0064, 0x0323, 0x0307, 0x0307, 0x0307, 0x0307},
        },
    }

    for i, r := range rows {
        {
            input := make([]rune, len(r.input))
            copy(input, r.input)
            ccc.ReorderRunes(input)
            assert.Equal(t, r.output, input, "for test %d", i)
        }
        {
            input := []byte(string(r.input))
            output := []byte(string(r.output))
            ccc.Reorder(input)
            assert.Equal(t, output, input, "for test %d", i)
        }
    }
}

func TestReorder_MaliciousInput(t *testing.T) {
    // Tests that you can't DoS Reorder with malicious input by ensuring it
    // completes in a reasonable time.

    var inBytes []byte
    inBytes = append(inBytes, []byte("\u0064")...)
    inBytes = append(inBytes, bytes.Repeat([]byte("\u0307"), 100)...)
    inBytes = append(inBytes, []byte("\u0323")...)
    inRunes := []rune(string(inBytes))

    ks.TestCompletes(t, 1 * time.Second, func() {
        var inBytesCopy []byte
        inBytesCopy = append(inBytesCopy, inBytes...)
        assert.Equal(t, ccc.ErrMaxNonStarters, ccc.ReorderRunes(inRunes))
        assert.Equal(t, ccc.ErrMaxNonStarters, ccc.Reorder(inBytes))
        rdr := transform.NewReader(strings.NewReader(string(inBytesCopy)), ccc.Transformer)
        _, err := io.ReadAll(rdr)
        assert.Equal(t, ccc.ErrMaxNonStarters, err)
    })
}

func TestTransform(t *testing.T) {

    type row struct {
        input func(int) string // already normalized
        expected func (int) string
    }

    rows := []row{
        // tests 0 to 3 have no runes with a combining class,
        // so should pass through unchanged

        { // test 0
            func(i int) string { return strings.Repeat("a", i) },
            func(i int) string { return strings.Repeat("a", i) },
        },
        { // test 1
            func(i int) string { return strings.Repeat("abcde", i) }, // 5 byte
            func(i int) string { return strings.Repeat("abcde", i) }, // 5 byte
        },
        { // test 2
            func(i int) string { return strings.Repeat("ab£d", i) }, // 5 byte
            func(i int) string { return strings.Repeat("ab£d", i) }, // 5 byte
        },
        { // test 3
            func(i int) string { return strings.Repeat("1\u20442", i) }, // 5 byte
            func(i int) string { return strings.Repeat("1\u20442", i) }, // 5 byte
        },
        { // test 4
            func(i int) string { return strings.Repeat("\u0064\u0307\u0323", i) + "a" }, // 5 byte
            func(i int) string { return strings.Repeat("\u0064\u0323\u0307", i) + "a" }, // 5 byte
        },
        { // test 5
            func(i int) string { return strings.Repeat("\u0064\u0307\u0323", i) }, // 5 byte
            func(i int) string { return strings.Repeat("\u0064\u0323\u0307", i) }, // 5 byte
        },
    }

    counts := []int{
        0, 1, 2, 3, 4, 5,
        63, 64, 65,
        70, 71, 72,
        127, 128, 129,
        255, 256, 257,
        511, 512, 513,
        584, 585, 586, // ~= 4096/7
        681, 682, 684, // ~= 4096/6
        819, 820, 821, // ~= 4096/5
        1023, 1024, 1025,
        1364, 1365, 1367, // ~= 4096/3
        1750, 1751, 1752, // ~= 4096 * (3/7)
        2047, 2048, 2049,
        4095, 4096, 4097,
        4125, 4126, 4127, 4128, 4129, // ~4096 + 30-32 non-starters maximums
        8191, 8192, 8193,
    }

    // Test various inputs
    for j, r := range rows {
        // Test various lengths
        for _, i := range counts {
            // fmt.Printf("Test %d at count %d\n", j, i)

            input := r.input(i)
            expected := r.expected(i)

            // construct a new transformer so that we can hit default size buffers
            rdr := transform.NewReader(strings.NewReader(input), ccc.Transformer)
            result, err := io.ReadAll(rdr)

            if !assert.Nil(t, err, "test %d with i=%d", j, i) { break }
            if !assert.Equal(t, expected, string(result),
                "test %d with i=%d\n" +
                "%x\n%x", j, i, expected, string(result)) { break }
        }
    }
}
