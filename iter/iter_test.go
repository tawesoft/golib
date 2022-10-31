package iter_test

import (
    "fmt"
    "math"
    "sort"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    lazy "github.com/tawesoft/golib/v2/iter"
    "golang.org/x/exp/maps"
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

func TestCounter(t *testing.T) {
    min := math.MinInt
    max := math.MaxInt

    type row struct {
        expected []int
        inputStart int
        inputStep  int
        take       int
    }

    rows := []row{
        { // row 0
            expected: []int{2, 4, 6, 8, 10},
            inputStart: 2,
            inputStep:  2,
            take:       5,
        },
        { // row 1
            expected: []int{max - 1, max + 0},
            inputStart: max - 1,
            inputStep:  1,
            take:       5,
        },
        { // row 2
            expected: []int{min + 1, min + 0},
            inputStart: min + 1,
            inputStep:  -1,
            take:       5,
        },
    }

    for i, r := range rows {
        expected := r.expected
        input := lazy.Counter(r.inputStart, r.inputStep)
        actual := lazy.ToSlice(lazy.Take(r.take, input))
        assert.Equalf(t, expected, actual, "row %d", i)
    }
}

func TestCutString(t *testing.T) {
    type row struct {
        sep rune
        input string
        expected []string
    }

    rows := []row{
        {'|', "",            []string{""}},
        {'|', "a",           []string{"a"}},
        {'|', "foo",         []string{"foo"}},
        {'|', "a|b|c",       []string{"a", "b", "c"}},
        {'|', "abc|def|xyz", []string{"abc", "def", "xyz"}},
        {'|', "¹|²|³",       []string{"¹", "²", "³"}},
        {'|',
            string([]rune{0x0041, 0x030a, '|', 0x064, 0x0307, 0x0327}),
            []string{
                string([]rune{0x0041, 0x030a}),
                string([]rune{0x0064, 0x0307, 0x0327}),
            },
        },
        {'\uFFFD',
            string([]rune{0x0041, 0x030a, 0xFFFD, 0x064, 0x0307, 0x0327}),
            []string{
                string([]rune{0x0041, 0x030a}),
                string([]rune{0x0064, 0x0307, 0x0327}),
            },
        },
    }

    for _, r := range rows {
        result := lazy.ToSlice(lazy.CutString(r.input, r.sep))
        assert.Equal(t, r.expected, result)
    }
}

func _testCutString(t *testing.T) {
    type row struct {
        sep string
        input string
        expected []string
    }

    rows := []row{
        {"|", "",            []string{""}},
        {"|", "|",           []string{"", ""}},
        {"|", "a",           []string{"a"}},
        {"|", "foo",         []string{"foo"}},
        {"|", "a|b|c",       []string{"a", "b", "c"}},
        {"|", "||",          []string{"", "", ""}},
        {"|", "abc|def|xyz", []string{"abc", "def", "xyz"}},
        {"|", "¹|²|³",       []string{"¹", "²", "³"}},
        {"|",
            string([]rune{0x0041, 0x030a, '|', 0x064, 0x0307, 0x0327}),
            []string{
                string([]rune{0x0041, 0x030a}),
                string([]rune{0x0064, 0x0307, 0x0327}),
            },
        },
        {"\uFFFD",
            string([]rune{0x0041, 0x030a, 0xFFFD, 0x064, 0x0307, 0x0327}),
            []string{
                string([]rune{0x0041, 0x030a}),
                string([]rune{0x0064, 0x0307, 0x0327}),
            },
        },
    }

    for _, r := range rows {
        result := lazy.ToSlice(lazy.CutStringStr(r.input, r.sep))
        assert.Equal(t, r.expected, result)
    }
}

func TestEnumerate(t *testing.T) {
    abc := lazy.FromSlice([]rune("abc"))

    expected := []lazy.Pair[int, rune]{
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

    expected := []lazy.Pair[string, string]{
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
    it := lazy.FromString("abc")
    z := rune(0)
    x, ok := it(); assert.Equal(t, 'a', x); assert.Equal(t,  true, ok)
    x, ok  = it(); assert.Equal(t, 'b', x); assert.Equal(t,  true, ok)
    x, ok  = it(); assert.Equal(t, 'c', x); assert.Equal(t,  true, ok)
    x, ok  = it(); assert.Equal(t,  z,  x); assert.Equal(t, false, ok)
    x, ok  = it(); assert.Equal(t,  z,  x); assert.Equal(t, false, ok)
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

func TestJoin_string(t *testing.T) {
    j := lazy.StringJoiner(", ") // can be reused

    assert.Equal(t, "a, b, c", lazy.Join(j, lazy.FromSlice([]string{"a", "b", "c"})))
    assert.Equal(t, "a, b, c", lazy.Join(j, lazy.FromSlice([]string{"a", "b", "c"}))) // reuse j
    assert.Equal(t, "a",       lazy.Join(j, lazy.FromSlice([]string{"a"})))
    assert.Equal(t, "",        lazy.Join(j, lazy.FromSlice([]string{""})))
    assert.Equal(t, "",        lazy.Join(j, lazy.FromSlice([]string(nil))))
}

func TestJoin_average(t *testing.T) {
    j := lazy.AverageJoiner[int]() // can be reused

    epsilon := 0.01

    {
        avg := lazy.Join(j, lazy.FromSlice([]int{2, 4, 6, 8}))
        assert.InDelta(t, 5.0, avg, epsilon)
    }
    {
        avg := lazy.Join(j, lazy.FromSlice([]int{2}))
        assert.InDelta(t, 2.0, avg, epsilon)
    }
    {
        avg := lazy.Join(j, lazy.FromSlice([]int{}))
        assert.InDelta(t, 0.0, avg, epsilon)
    }
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

func TestPairwiseEnd(t *testing.T) {
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

func TestReduce(t *testing.T) {
    mul := func(a int, b int) int { return a * b }
    multiplier := lazy.Reducer[int]{Reduce: mul, Identity: 1}
    {
        x := lazy.Reduce(multiplier, lazy.FromSlice([]int{1, 2, 3, 4}))
        // 1 * 2 * 3 * 4 = 24
        assert.Equal(t, 24, x)
    }
    {
        x := lazy.Reduce(multiplier, lazy.FromSlice([]int{4}))
        assert.Equal(t, 4, x)
    }
    {
        x := lazy.Reduce(multiplier, lazy.FromSlice([]int{}))
        assert.Equal(t, 1, x)
    }
    {
        x := lazy.Reduce(multiplier, lazy.FromSlice([]int(nil)))
        assert.Equal(t, 1, x)
    }
}

func TestRepeat(t *testing.T) {
    {
        input := lazy.Repeat(3, [2]int{1, 2})
        expected := [][2]int{{1, 2}, {1, 2}, {1, 2}}
        assert.Equal(t, expected, lazy.ToSlice(input))
    }
    {
        input := lazy.Repeat(5, "foo")
        expected := []string{"foo", "foo", "foo", "foo", "foo"}
        assert.Equal(t, expected, lazy.ToSlice(input))
    }
    {
        input := lazy.Repeat(1, "foo")
        expected := []string{"foo"}
        assert.Equal(t, expected, lazy.ToSlice(input))
    }
    {
        input := lazy.Repeat(0, "foo")
        expected := []string{}
        assert.Equal(t, expected, lazy.ToSlice(input))
    }
}

func TestTake(t *testing.T) {
    {
        input := lazy.FromSlice([]int{1, 2, 3, 4, 5, 6})
        expected := []int{1, 2, 3}
        taken := lazy.Take(3, input)
        assert.Equal(t, expected, lazy.ToSlice(taken))
    }
    {
        input := lazy.FromSlice([]int{1, 2, 3})
        expected := []int{1, 2, 3}
        taken := lazy.Take(5, input)
        assert.Equal(t, expected, lazy.ToSlice(taken))
    }
    {
        input := lazy.FromSlice([]int{})
        expected := []int{}
        taken := lazy.Take(3, input)
        assert.Equal(t, expected, lazy.ToSlice(taken))
    }
    {
        input := lazy.FromSlice([]int(nil))
        expected := []int{}
        taken := lazy.Take(3, input)
        assert.Equal(t, expected, lazy.ToSlice(taken))
    }
}

func TestTee(t *testing.T) {
    abc := lazy.FromSlice([]rune("abc"))

    gs := lazy.Tee(3, abc)
    // gs[0] produces 'a', 'b', 'c'
    // gs[1] produces 'a', 'b', 'c'
    // gs[2] produces 'a', 'b', 'c'

    type row struct{
        generator    lazy.It[rune]
        expectedRune rune
        expectedOk   bool
    }

    // the iterators must produce their sequence, even if iterated out-of-order
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

func TestToMap(t *testing.T) {
    f := func() lazy.It[lazy.Pair[string, int]] {
        i := 0
        return func() (lazy.Pair[string, int], bool) {
            i++
            switch i {
                case 1: return lazy.Pair[string, int]{"foo", 1}, true
                case 2: return lazy.Pair[string, int]{"bar", 2}, true
                case 3: return lazy.Pair[string, int]{"baz", 3}, true
                case 4: return lazy.Pair[string, int]{"baz", 4}, true
                default: return lazy.Pair[string, int]{}, false
            }
        }
    }

    // nil chooser means that on key collision, last added wins by default
    choose := (func(string, int, int) int)(nil)

    m := lazy.ToMap(choose, f())

    expected := map[string]int{
        "foo": 1,
        "bar": 2,
        "baz": 4,
    }
    assert.True(t, maps.Equal(expected, m))
}

func TestToMap_empty(t *testing.T) {
    f := func() (lazy.Pair[string, int], bool) {
        return lazy.Pair[string, int]{}, false
    }
    assert.Equal(t, map[string]int{}, lazy.ToMap(nil, f))
}

func TestToSlice(t *testing.T) {
    f := func() lazy.It[int] {
        i := 0
        return func() (int, bool) {
            i++
            switch i {
                case 1: return 1, true
                case 2: return 2, true
                case 3: return 3, true
                case 4: return 4, true
                default: return 0, false
            }
        }
    }

    s := lazy.ToSlice(f())

    expected := []int{1, 2, 3, 4}
    assert.Equal(t, expected, s)
}

func TestToSlice_empty(t *testing.T) {
    f := func() (int, bool) {
        return 0, false
    }
    assert.Equal(t, []int{}, lazy.ToSlice(f))
}

func TestToString(t *testing.T) {
    f := func() lazy.It[rune] {
        i := 0
        return func() (rune, bool) {
            i++
            switch i {
                case 1: return 'a', true
                case 2: return 'b', true
                case 3: return 'c', true
                default: return 0, false
            }
        }
    }

    assert.Equal(t, "abc", lazy.ToString(f()))
}

func TestToString_empty(t *testing.T) {
    f := func() (rune, bool) {
        return 0, false
    }
    assert.Equal(t, "", lazy.ToString(f))
}

func TestWalk_stringBuilder(t *testing.T) {
    var sb strings.Builder
    strings := lazy.FromSlice([]string{"one", "two", "three"})
    write := func(x string) { sb.WriteString(x) }
    lazy.Walk(write, strings)
    assert.Equal(t, "onetwothree", sb.String())
}

func TestWalkFinal_stringBuilder(t *testing.T) {
    var sb strings.Builder
    strings := lazy.FromSlice([]string{"one", "two", "three"})
    write := func(x string, final bool) {
        sb.WriteString(x)
        if !final { sb.WriteString(", ") }
    }
    lazy.WalkFinal(write, strings)
    assert.Equal(t, "one, two, three", sb.String())
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
