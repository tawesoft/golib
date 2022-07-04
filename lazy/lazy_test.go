package lazy_test

import (
    "fmt"
    "sort"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/v2/lazy"
)

// CONTRIBUTORS: keep tests in alphabetical order, but with examples grouped
// first.

func TestAll(t *testing.T) {
    isEven := func(x int) bool { return x % 2 == 0 }
    assert.True (t, lazy.All(isEven, lazy.FromSlice([]int{})))
    assert.True (t, lazy.All(isEven, lazy.FromSlice([]int{2})))
    assert.True (t, lazy.All(isEven, lazy.FromSlice([]int{2, 4, 6})))
    assert.False(t, lazy.All(isEven, lazy.FromSlice([]int{2, 3, 6})))
}

func TestAny(t *testing.T) {
    isEven := func(x int) bool { return x % 2 == 0 }
    assert.False(t, lazy.Any(isEven, lazy.FromSlice([]int{})))
    assert.True (t, lazy.Any(isEven, lazy.FromSlice([]int{2})))
    assert.False(t, lazy.Any(isEven, lazy.FromSlice([]int{3})))
    assert.True (t, lazy.Any(isEven, lazy.FromSlice([]int{2, 4, 6})))
    assert.True (t, lazy.Any(isEven, lazy.FromSlice([]int{2, 3, 6})))
    assert.False(t, lazy.Any(isEven, lazy.FromSlice([]int{1, 3, 5})))
}

func TestAppendToSlice(t *testing.T) {
    {
        xs := []int{1, 2, 3}
        ys := lazy.FromSlice([]int{4, 5, 6})
        expected := []int{1, 2, 3, 4, 5, 6}
        xs = lazy.AppendToSlice(xs, ys)
        assert.Equal(t, expected, xs)
    }
    {
        xs := []int{}
        ys := lazy.FromSlice([]int{4, 5, 6})
        expected := []int{4, 5, 6}
        xs = lazy.AppendToSlice(xs, ys)
        assert.Equal(t, expected, xs)
    }
    {
        xs := []int{1, 2, 3}
        ys := lazy.FromSlice([]int{})
        expected := []int{1, 2, 3}
        xs = lazy.AppendToSlice(xs, ys)
        assert.Equal(t, expected, xs)
    }
    {
        xs := []int{}
        ys := lazy.FromSlice([]int{})
        expected := []int{}
        xs = lazy.AppendToSlice(xs, ys)
        assert.Equal(t, expected, xs)
    }
    {
        var xs []int = nil
        ys := lazy.FromSlice([]int{4, 5, 6})
        expected := []int{4, 5, 6}
        xs = lazy.AppendToSlice(xs, ys)
        assert.Equal(t, expected, xs)
    }
    {
        // appending nothing to a nil slice should leave it as a nil slice
        var xs []int = nil
        var nilInts []int = nil
        ys := lazy.FromSlice(nilInts)
        var expected []int = nil
        xs = lazy.AppendToSlice(xs, ys)
        assert.Equal(t, expected, xs)
    }
}

func TestCat(t *testing.T) {
    abc := lazy.FromSlice([]rune("abc"))
    def := lazy.FromSlice([]rune("def"))
    xyz := lazy.FromSlice[rune](nil) // empty slice
    abcdef := lazy.Cat(abc, def, xyz)
    assert.Equal(t, []rune("abcdef"), lazy.ToSlice(abcdef))
}

func TestCheck(t *testing.T) {
    errorIfOdd := func (x int) error {
        if x % 2 != 0 {
            return fmt.Errorf("got odd number %d", x)
        }
        return nil
    }

    type row struct {
        input []int
        expectedValue int
        expectedError error
    }
    rows := []row{
        {
            input: []int{},
            expectedValue: 0,
            expectedError: nil,
        },
        {
            input: []int{2, 4, 8, 10},
                expectedValue: 0,
            expectedError: nil,
        },
        {
            input: []int{2, 4, 7, 10},
            expectedValue: 7,
            expectedError: fmt.Errorf("got odd number 7"),
        },
    }

    for i, r := range rows {
        x, err := lazy.Check(errorIfOdd, lazy.FromSlice(r.input))
        assert.Equal(t,   x, r.expectedValue, "test %i expected value", i)
        assert.Equal(t, err, r.expectedError, "test %i expected error", i)
    }
}

func TestEnumerate(t *testing.T) {
    abc := lazy.FromSlice([]rune("abc"))

    expected := []lazy.Item[int, rune]{
        {0, 'a'},
        {1, 'b'},
        {2, 'c'},
    }

    result := lazy.ToSlice(lazy.Enumerate(abc))

    assert.Equal(t, expected, result)
}

func TestFilter(t *testing.T) {
    isOdd := func(x int) bool { return x % 2 == 1 }

    assert.Equal(t, []int{1, 3},
        lazy.ToSlice(lazy.Filter(isOdd, lazy.FromSlice([]int{1, 2, 3}))))

    assert.Equal(t, []int{},
        lazy.ToSlice(lazy.Filter(isOdd, lazy.FromSlice([]int{}))))
}

func TestFinal(t *testing.T) {
    {
        abc := lazy.FromSlice([]int{1, 2, 3})

        expected := []lazy.FinalValue[int]{
            {1, false},
            {2, false},
            {3,  true},
        }

        assert.Equal(t, expected, lazy.ToSlice(lazy.Final(abc)))
    }
    {
        abc := lazy.FromSlice([]int{1})
        expected := []lazy.FinalValue[int]{{1, true}}
        assert.Equal(t, expected, lazy.ToSlice(lazy.Final(abc)))
    }
    {
        abc := lazy.FromSlice([]int{})
        expected := []lazy.FinalValue[int]{}
        assert.Equal(t, expected, lazy.ToSlice(lazy.Final(abc)))
    }
    {
        abc := lazy.FromSlice([]int(nil))
        expected := []lazy.FinalValue[int]{}
        assert.Equal(t, expected, lazy.ToSlice(lazy.Final(abc)))
    }
}

func TestFromMap(t *testing.T) {
    original := map[string]string{
        "cat": "meow",
        "dog": "woof",
        "cow": "moo",
    }
    kvs := lazy.ToSlice(lazy.FromMap(original))

    sort.Slice(kvs, func(i int, j int) bool {
        return kvs[i].Key < kvs[j].Key
    })

    expected := []lazy.Item[string, string]{
        {Key: "cat", Value: "meow"},
        {Key: "cow", Value: "moo"},
        {Key: "dog", Value: "woof"},
    }

    assert.Equal(t, expected, kvs)
}

func TestFromSlice(t *testing.T) {
    it := lazy.FromSlice([]int{1, 2, 3})

    x, ok := it(); assert.Equal(t, 1, x); assert.Equal(t,  true, ok)
    x, ok  = it(); assert.Equal(t, 2, x); assert.Equal(t,  true, ok)
    x, ok  = it(); assert.Equal(t, 3, x); assert.Equal(t,  true, ok)
    x, ok  = it(); assert.Equal(t, 0, x); assert.Equal(t, false, ok)
    x, ok  = it(); assert.Equal(t, 0, x); assert.Equal(t, false, ok)

    it = lazy.FromSlice([]int(nil))
    x, ok  = it(); assert.Equal(t, 0, x); assert.Equal(t, false, ok)
}

func TestFromString(t *testing.T) {
    abc := lazy.FromString("abc")
    assert.Equal(t, []rune{'a', 'b', 'c'}, lazy.ToSlice(abc))
}

func TestFunc(t *testing.T) {
    it := func () lazy.It[int] {
        var i int
        return func() (int, bool) {
            if i == 3 { return 0, false }
            i++
            return i, true
        }
    }

    f := lazy.Func(it())

    x, ok := f(); assert.Equal(t, 1, x); assert.Equal(t,  true, ok)
    x, ok  = f(); assert.Equal(t, 2, x); assert.Equal(t,  true, ok)
    x, ok  = f(); assert.Equal(t, 3, x); assert.Equal(t,  true, ok)
    x, ok  = f(); assert.Equal(t, 0, x); assert.Equal(t, false, ok)
    x, ok  = f(); assert.Equal(t, 0, x); assert.Equal(t, false, ok)
}

func TestInsertToMap(t *testing.T) {
    {
        base := map[string]string{
            "cat": "meow",
            "dog": "woof",
            "cow": "moo",
        }
        extras := lazy.FromMap(map[string]string{
            "sheep": "baa",
            "dog": "bark",
        })

        lazy.InsertToMap(base, nil, extras)

        expected := map[string]string{
            "cat": "meow",
            "dog": "bark",
            "cow": "moo",
            "sheep": "baa",
        }

        assert.Equal(t, expected, base)
    }

    {
        base := map[string]string{
            "cat": "meow",
            "dog": "woof",
            "cow": "moo",
        }
        extras := lazy.FromMap(map[string]string{
            "sheep": "baa",
            "dog": "bark",
        })

        choose := func(key string, original string, new string) string {
            return original
        }

        lazy.InsertToMap(base, choose, extras)

        expected := map[string]string{
            "cat": "meow",
            "dog": "woof",
            "cow": "moo",
            "sheep": "baa",
        }

        assert.Equal(t, expected, base)
    }
}

func TestJoin(t *testing.T) {
    sum := func(a int, b int) int { return a + b }

    assert.Equal(t, 0, lazy.Join(sum, lazy.FromSlice([]int(nil))))
    assert.Equal(t, 0, lazy.Join(sum, lazy.FromSlice([]int{})))
    assert.Equal(t, 5, lazy.Join(sum, lazy.FromSlice([]int{5})))
    assert.Equal(t, 6, lazy.Join(sum, lazy.FromSlice([]int{1, 2, 3})))

}

func TestMap(t *testing.T) {
    {
        double := func(a int) int { return a + a }
        xs := lazy.FromSlice([]int{1, 2, 3})
        assert.Equal(t, []int{2, 4, 6}, lazy.ToSlice(lazy.Map(double, xs)))
    }
    {
        toString := func(a int) string { return fmt.Sprintf("%d", a) }
        xs := lazy.FromSlice([]int{1, 2, 3})
        assert.Equal(t, []string{"1", "2", "3"}, lazy.ToSlice(lazy.Map(toString, xs)))
    }
    {
        double := func(a int) int { return a + a }
        xs := lazy.FromSlice([]int{})
        assert.Equal(t, []int{}, lazy.ToSlice(lazy.Map(double, xs)))
    }
    {
        double := func(a int) int { return a + a }
        xs := lazy.FromSlice([]int(nil))
        assert.Equal(t, []int{}, lazy.ToSlice(lazy.Map(double, xs)))
    }
}

func TestPairwise(t *testing.T) {
    {
        abcd := lazy.FromSlice([]rune("abcd"))

        expected := [][2]rune{
            {'a', 'b'},
            {'b', 'c'},
            {'c', 'd'},
        }

        result := lazy.ToSlice(lazy.Pairwise(abcd))

        assert.Equal(t, expected, result)
    }
    {
        in := []int(nil)
        expected := [][2]int{}
        result := lazy.ToSlice(lazy.Pairwise(lazy.FromSlice(in)))
        assert.Equal(t, expected, result)
    }
    {
        in := []int{}
        expected := [][2]int{}
        result := lazy.ToSlice(lazy.Pairwise(lazy.FromSlice(in)))
        assert.Equal(t, expected, result)
    }
    {
        in := []int{1}
        expected := [][2]int{}
        result := lazy.ToSlice(lazy.Pairwise(lazy.FromSlice(in)))
        assert.Equal(t, expected, result)
    }
}

func TestPairwiseFill(t *testing.T) {
    {
        abcd := lazy.FromSlice([]rune("abcd"))

        expected := [][2]rune{
            {'a', 'b'},
            {'b', 'c'},
            {'c', 'd'},
            {'d', 0xFFFF},
        }

        result := lazy.ToSlice(lazy.PairwiseEnd(0xFFFF, abcd))

        assert.Equal(t, expected, result)
    }
    {
        in := []int(nil)
        expected := [][2]int{}
        result := lazy.ToSlice(lazy.PairwiseEnd(999, lazy.FromSlice(in)))
        assert.Equal(t, expected, result)
    }
    {
        in := []int{}
        expected := [][2]int{}
        result := lazy.ToSlice(lazy.PairwiseEnd(999, lazy.FromSlice(in)))
        assert.Equal(t, expected, result)
    }
    {
        in := []int{1}
        expected := [][2]int{
            {1, 999},
        }
        result := lazy.ToSlice(lazy.PairwiseEnd(999, lazy.FromSlice(in)))
        assert.Equal(t, expected, result)
    }
}

// TODO tests from here (alphabetical order)

func TestTee(t *testing.T) {
    abc := lazy.FromSlice([]rune("abc"))

    gs := lazy.Tee(3, abc)

    type row struct{
        generator    lazy.It[rune]
        expectedRune rune
        expectedOk   bool
    }
    rows := []row{
        {gs[0], 'a',  true},
        {gs[1], 'a',  true},
        {gs[2], 'a',  true},
        {gs[0], 'b',  true},
        {gs[0], 'c',  true},
        {gs[0],  0,  false},
        {gs[0],  0,  false},
        {gs[0],  0,  false},
        {gs[0],  0,  false},
        {gs[1], 'b',  true},
        {gs[2], 'b',  true},
        {gs[2], 'c',  true},
        {gs[2],  0,  false},
        {gs[1], 'c',  true},
        {gs[1],  0,  false},
    }

    for i, row := range rows {
        r, ok := row.generator()
        assert.Equal(t, row.expectedOk, ok, "test %d", i)
        assert.Equal(t, string(row.expectedRune), string(r), "test %d", i)
    }
}

func TestToString(t *testing.T) {
    abc := lazy.FromSlice([]rune{'a', 'b', 'c'})
    assert.Equal(t, "abc", lazy.ToString(abc))
}

func TestWalk_stringBuilder(t *testing.T) {
    var sb strings.Builder
    strings := lazy.FromSlice([]string{"one", "two", "three"})
    write := func(x string) { sb.WriteString(x) }
    lazy.Walk(write, strings)
    assert.Equal(t, "onetwothree", sb.String())
}

func TestZip(t *testing.T) {
    {
        a := lazy.FromSlice([]int{  1,   2,   3})
        b := lazy.FromSlice([]int{ 10,  20,  30})
        c := lazy.FromSlice([]int{100, 200, 300, 400})
        expected := [][]int{{1, 10, 100}, {2, 20, 200}, {3, 30, 300}}
        assert.Equal(t, expected, lazy.ToSlice(lazy.Zip(a, b, c)))
    }
    {
        a := lazy.FromSlice([]int{  1,   2,   3})
        expected := [][]int{{1}, {2}, {3}}
        assert.Equal(t, expected, lazy.ToSlice(lazy.Zip(a)))
    }
    {
        a := lazy.FromSlice([]int{  1,   2,   3})
        b := lazy.FromSlice([]int{})
        expected := [][]int{}
        assert.Equal(t, expected, lazy.ToSlice(lazy.Zip(a, b)))
    }
    {
        a := lazy.FromSlice([]int(nil))
        b := lazy.FromSlice([]int(nil))
        expected := [][]int{}
        assert.Equal(t, expected, lazy.ToSlice(lazy.Zip(a, b)))
    }
    {
        a := lazy.FromSlice([]int(nil))
        expected := [][]int{}
        assert.Equal(t, expected, lazy.ToSlice(lazy.Zip(a)))
    }
    {
        assert.Equal(t, [][]any{}, lazy.ToSlice(lazy.Zip[any]()))
    }
}

func TestZip_withStrings(t *testing.T) {
    abc := lazy.FromString("abc")
    def := lazy.FromString("def")
    wxyz := lazy.FromString("wxyz")

    expectedStrings := []string{
        "adw", "bex", "cfy",
    }

    runes2string := func (x []rune) string {
        return string(x)
    }

    zippedToStrings := lazy.ToSlice(lazy.Map(runes2string, lazy.Zip(abc, def, wxyz)))
    assert.Equal(t, expectedStrings, zippedToStrings)
}

func TestZipFlat(t *testing.T) {
    abc := lazy.FromString("abc")
    def := lazy.FromString("def")
    wxyz := lazy.FromString("wxyz")

    result := lazy.ToString(lazy.ZipFlat(abc, def, wxyz))

    assert.Equal(t, "adwbexcfy", result)
}
