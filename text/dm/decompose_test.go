package dm_test

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "os"
    "path"
    "runtime"
    "strconv"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/v2/ks"
    "github.com/tawesoft/golib/v2/text/dm"
    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
)

func ExampleMap() {

    input := '²'
    dt, dm := dm.Map(input)
    fmt.Printf("%c => decomposition (%s): %s\n", input, dt, string(dm))

    if dt.IsCompat() {
        fmt.Println("This is a compatibility decomposition, not a canonical one")
    } else if dt.IsCanonical() {
        fmt.Println("This is a canonical decomposition")
    } else {
        fmt.Println("There isn't a decomposition for this input")
    }

    // Output:
    // ² => decomposition (Super): 2
    // This is a compatibility decomposition, not a canonical one
}

func TestMap(t *testing.T) {
    type row struct {
        input rune
        dt    dm.Type
        dm    []rune
    }

    rows := []row{
        {'a', dm.None,      nil},
        {'ñ', dm.Canonical, []rune{0x006E, 0x0303}},

        // canonical singleton replacement
        {'Å', dm.Canonical, []rune{0x00C5}}, // Å
        // Å from above is still not a complete decomposition...
        {'Å', dm.Canonical, []rune{0x0041, 0x030A}}, // final

        // // canonical singleton replacement
        {'Ω', dm.Canonical, []rune{0x03A9}},

        {'ế', dm.Canonical, []rune{0x00EA, 0x301}}, // ê plus Combining Acute Accent,
        // ê from above is not a complete decomposition...
        {'ê', dm.Canonical, []rune{0x0065, 0x0302}}, // e plus Combining Circumflex Accent

        {'²', dm.Super,     []rune{'2'}},
        {'½', dm.Fraction,  []rune{'1', 0x2044, '2'}}, // Fraction Slash

    }

    for i, r := range rows {
        input, expectedDt, expectedDm := r.input, r.dt, r.dm
        dt, dm := dm.Map(input)
        assert.Equal(t, expectedDt, dt, "test(%d) %c dt", i, input)
        assert.Equal(t, expectedDm, dm, "test(%d) %c dm", i, input)
    }
}

func TestDecomposer_String(t *testing.T) {
    type row struct {
        dc     dm.Decomposer
        input  []rune
        output []rune
    }

    noFra := dm.Except(dm.Fraction)

    rows := []row{
        {dm.CD, []rune{'a'}, []rune{'a'}},
        {dm.CD, []rune{'ñ'}, []rune{0x006E, 0x0303}},

        // http://wiki.squeak.org/squeak/6265
        {dm.CD, []rune{0x1E0B, 0x0323}, []rune{0x0064, 0x0323, 0x0307}},

        // canonical singleton replacement
        {dm.CD, []rune{'Å'}, []rune{0x0041, 0x030A}},
        {dm.CD, []rune{'Ω'}, []rune{0x03A9}},

        {dm.CD, []rune{'ế'}, []rune{0x0065, 0x0302, 0x301}},

        {dm.KD, []rune{'²'}, []rune{'2'}},
        {dm.KD, []rune{'½'}, []rune{'1', 0x2044, '2'}}, // Fraction Slash

        // Suppress certain decompositions
        {noFra, []rune{'½'}, []rune{'½'}},
    }

    for i, r := range rows {
        s := r.dc.String(string(r.input))
        assert.Equal(t, string(r.output), s, "test(%d) %x, got %x, expected %x", i, r.input, []rune(s), r.output)
    }
}

func TestDecomposer_Transform(t *testing.T) {

    type row struct {
        input func(int) string
        expected func (int) string
        norm dm.Decomposer
    }

    rows := []row{
        { // test 0
            func(i int) string { return strings.Repeat("a", i) },
            func(i int) string { return strings.Repeat("a", i) },
            dm.CD,
        },
        { // test 1
            func(i int) string { return strings.Repeat("abcde", i) }, // 5 byte
            func(i int) string { return strings.Repeat("abcde", i) }, // 5 byte
            dm.CD,
        },
        { // test 2
            func(i int) string { return strings.Repeat("ab£d", i) }, // 5 byte
            func(i int) string { return strings.Repeat("ab£d", i) }, // 5 byte
            dm.CD,
        },
        { // test 3
            func(i int) string { return strings.Repeat("1\u20442", i) }, // 5 byte
            func(i int) string { return strings.Repeat("1\u20442", i) }, // 5 byte
            dm.KD,
        },
        { // test 4
            func(i int) string { return strings.Repeat("½", i) }, // 2 byte
            func(i int) string { return strings.Repeat("1\u20442", i) }, // 5 byte
            dm.KD,
        },
        { // test 5
            func(i int) string { return strings.Repeat("ế", i) + "a" },
            func(i int) string { return strings.Repeat("\u0065\u0302\u0301", i) + "a" },
            dm.CD,
        },
        { // test 6
            func(i int) string { return strings.Repeat("ế", i) }, // 3 byte
            func(i int) string { return strings.Repeat("\u0065\u0302\u0301", i) }, // 5 byte
            dm.CD,
        },
        { // test 7
            func(i int) string { return strings.Repeat("\u1E0B\u0323", i) + "a" }, // 5 byte
            func(i int) string { return strings.Repeat("\u0064\u0323\u0307", i) + "a" }, // 5 byte
            dm.CD,
        },
        { // test 8
            func(i int) string { return strings.Repeat("\u1E0B\u0323", i) }, // 5 byte
            func(i int) string { return strings.Repeat("\u0064\u0323\u0307", i) }, // 5 byte
            dm.CD,
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
        if j > 4 { break }
        // Test various lengths
        for _, i := range counts {
            input := r.input(i)
            expected := r.expected(i)

            // construct a new transformer so that we can hit default size buffers
            rdr := transform.NewReader(strings.NewReader(input), r.norm.Transformer())
            result, err := io.ReadAll(rdr)

            if !assert.Nil(t, err, "test %d with i=%d", j, i) { break }
            if !assert.Equal(t, expected, string(result),
                "test %d with i=%d\n" +
                "%x\n%x", j, i, expected, string(result)) { break }
        }
    }
}

func relpath(t *testing.T, file string) string {
    _, filename, _, _ := runtime.Caller(0)
    return path.Join(path.Dir(filename), file)
}

func TestUCD(t *testing.T) {
    var INPUT = relpath(t, "NormalizationTest.13.0.0.txt")

    parse := func(x string) string {
        x = strings.TrimSpace(x)
        result := make([]rune, 0)
        xs := strings.Split(x, " ")
        for _, i := range xs {
            i = strings.TrimSpace(i)
            r, err := strconv.ParseUint(i, 16, 32)
            if err != nil { panic(err) }
            result = append(result, rune(r))
        }
        return string(result)
    }

    var lineno int
    seenLineno := make(map[int]struct{})
    equal := func(a string, b string, test string, comment string) {
        if _, ok := seenLineno[lineno]; ok { return }
        if a != b {
            matchNFD := ks.IfThenElse(a == norm.NFD.String(a), "match", "differ")
            matchNFKD := ks.IfThenElse(a == norm.NFKD.String(a), "match", "differ")
            t.Errorf("test %s line %d:\n%x != %x\n%x Go's NFD (%s)\n%x Go's NFKD (%s)\n%s",
                test, lineno, []rune(a), []rune(b),
                []rune(norm.NFD.String(a)),  matchNFD,
                []rune(norm.NFKD.String(a)), matchNFKD,
                comment,
            )
        }
        seenLineno[lineno] = struct{}{}
    }

    f, err := os.Open(INPUT)
    if errors.Is(err, os.ErrNotExist) {
        t.Skipf("missing test data")
        return
    } else if err != nil { panic(err) }
    defer f.Close()

    rdr := bufio.NewReaderSize(f, 64*1024)
    for {
        lineno++
        ln, isPrefix, err := rdr.ReadLine()
        if isPrefix { ks.Never() }
        if (err != nil) && errors.Is(err, io.EOF) { break }
        if err != nil { panic(err) }
        if ln == nil { break }
        if len(ln) == 0 { continue }
        if ln[0] == '#' { continue } // # Comment
        if ln[0] == '@' { continue } // @PartN
        l := string(ln)

        // source; NFC; NFD; NFKC; NFKD; # comment
        cols := strings.SplitN(l, ";", 6)
        ks.Assert(len(cols) >= 5)
        c1 := parse(cols[0])
        c2 := parse(cols[1])
        c3 := parse(cols[2])
        c4 := parse(cols[3])
        c5 := parse(cols[4])
        comment := strings.TrimSpace(cols[5])

        // NFD
        // c3 ==  toNFD(c1) ==  toNFD(c2) ==  toNFD(c3)
        equal(c3, dm.CD.String(c1), "c3 == toNFD(c1)", comment)
        equal(c3, dm.CD.String(c2), "c3 == toNFD(c2)", comment)
        equal(c3, dm.CD.String(c3), "c3 == toNFD(c3)", comment)
        // c5 ==  toNFD(c4) ==  toNFD(c5)
        equal(c5, dm.CD.String(c4), "c5 == toNFD(c4)", comment)
        equal(c5, dm.CD.String(c5), "c5 == toNFD(c5)", comment)

        // NFKD
        // c5 == toNFKD(c1) == toNFKD(c2) == toNFKD(c3) == toNFKD(c4) == toNFKD(c5)
        equal(c5, dm.KD.String(c1), "c5 == toNFKD(c1)", comment)
        equal(c5, dm.KD.String(c2), "c5 == toNFKD(c2)", comment)
        equal(c5, dm.KD.String(c3), "c5 == toNFKD(c3)", comment)
        equal(c5, dm.KD.String(c4), "c5 == toNFKD(c4)", comment)
        equal(c5, dm.KD.String(c5), "c5 == toNFKD(c5)", comment)

        // transformer versions...
        trans := func(a string, f dm.Decomposer, b string, test string, comment string) {
            r := transform.NewReader(strings.NewReader(b), f.Transformer())
            xs, err := io.ReadAll(r)
            if !assert.Nil(t, err) { return }
            equal(a, string(xs), test, comment)
        }
        trans(c3, dm.CD, c1, "c3 == toNFD(c1) [transformer]", comment)
        trans(c3, dm.CD, c2, "c3 == toNFD(c2) [transformer]", comment)
        trans(c3, dm.CD, c3, "c3 == toNFD(c3) [transformer]", comment)
        trans(c5, dm.CD, c4, "c5 == toNFD(c4) [transformer]", comment)
        trans(c5, dm.CD, c5, "c5 == toNFD(c5) [transformer]", comment)

        trans(c5, dm.KD, c1, "c5 == toNFKD(c1) [transformer]", comment)
        trans(c5, dm.KD, c2, "c5 == toNFKD(c2) [transformer]", comment)
        trans(c5, dm.KD, c3, "c5 == toNFKD(c3) [transformer]", comment)
        trans(c5, dm.KD, c4, "c5 == toNFKD(c4) [transformer]", comment)
        trans(c5, dm.KD, c5, "c5 == toNFKD(c5) [transformer]", comment)
    }
}
