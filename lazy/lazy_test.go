package lazy_test

import (
    "fmt"
    "strconv"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/lazy"
)

func ExampleFunction() {
    // generate an infinite sequence of integers with a function
    integers := func() (func() (int, bool)) {
        i := 0
        return func() (int, bool) {
            result := i
            i = i + 1
            return result, true
        }
    }

    integerGenerator := lazy.Function(integers())
    firstFour := lazy.TakeN(4, integerGenerator)

    fmt.Printf("First four integers are: %v\n", lazy.ToSlice(firstFour))

    // Output:
    // First four integers are: [0 1 2 3]
}

func ExampleToDict() {
    type Person struct {
        name string
        age int
    }
    type ID string

    // a map of people
    people := lazy.FromDict(map[ID]Person{
        "ATUR001": {name: "Alice Turing", age: 23},
        "GHOP001": {name: "George Hopper", age: 60},
        "FKAH001": {name: "Freddy Kahlo", age: 29},
    })

    // this filter function returns true for people under thirty
    underThirty := func(kv lazy.Item[ID, Person]) bool {
        return kv.Value.age < 30
    }

    // apply the filter and finally generate a dict
    peopleUnderThirty := lazy.ToDict(lazy.Filter(underThirty, people))

    // printer function
    p := func(lookup map[ID]Person, id ID) {
        if person, ok := lookup[id]; ok {
            fmt.Printf("%s: %+v\n", id, person)
        } else {
            fmt.Printf("%s: NOT FOUND\n", id)
        }
    }

    p(peopleUnderThirty, "ATUR001")
    p(peopleUnderThirty, "GHOP001") // missing!
    p(peopleUnderThirty, "FKAH001")

    // Output:
    // ATUR001: {name:Alice Turing age:23}
    // GHOP001: NOT FOUND
    // FKAH001: {name:Freddy Kahlo age:29}
}

func ExampleMap() {
    numbersAsStrings := lazy.FromSlice([]string{"1", "2", "3", "4"})

    // atoi returns the integer x from the string "x"
    atoi := func (s string) int {
        i, _ := strconv.Atoi(s)
        return i
    }

    doubler := func (i int) int {
        return i * 2
    }

    fmt.Printf("%v\n", lazy.ToSlice(
        lazy.Map[int, int](doubler,     // =>  [2 4 6 8]
            lazy.Map[string, int](atoi, // => [1 2 3 4]
                numbersAsStrings))))    // => ["1" "2" "3" "4"]

    // Output:
    // [2 4 6 8]
}

func ExampleMap_dict() {
    type ID string
    type Person struct {
        id   ID
        name string
        age  int
    }

    // given a list of people, we want a map (id -> person)
    people := lazy.FromSlice([]Person{
        {id: "ATUR001", name: "Alice Turing",  age: 23},
        {id: "GHOP001", name: "George Hopper", age: 60},
        {id: "FKAH001", name: "Freddy Kahlo",  age: 29},
    })

    // for a person input, this function returns (id, person)
    personToTuple := func (person Person) lazy.Item[ID, Person] {
        return lazy.Item[ID, Person]{person.id, person}
    }

    // apply the function over all people (lazily...)
    peopleTuples := lazy.Map(personToTuple, people)

    // finally generate a dict
    peopleByID := lazy.ToDict(peopleTuples)

    // printer function
    p := func(lookup map[ID]Person, id ID) {
        if person, ok := lookup[id]; ok {
            fmt.Printf("%s: %+v\n", id, person)
        } else {
            fmt.Printf("%s: NOT FOUND\n", id)
        }
    }

    p(peopleByID, "ATUR001")

    // Output:
    // ATUR001: {id:ATUR001 name:Alice Turing age:23}
}

func TestCat(t *testing.T) {
    abc := lazy.FromSlice([]rune("abc"))
    def := lazy.FromSlice([]rune("def"))
    xyz := lazy.FromSlice[rune](nil) // empty slice
    abcdef := lazy.Cat(abc, def, xyz)
    assert.Equal(t, []rune("abcdef"), lazy.ToSlice(abcdef))
}

func TestTee(t *testing.T) {
    abc := lazy.FromSlice([]rune("abc"))

    gs := lazy.Tee(3, abc)

    type row struct{
        generator    lazy.Generator[rune]
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
        r, ok := row.generator.Next()
        assert.Equal(t, row.expectedOk, ok, "test %d", i)
        assert.Equal(t, string(row.expectedRune), string(r), "test %d", i)
    }
}

func TestZip(t *testing.T) {
    abc := lazy.FromSlice([]rune("abc"))
    def := lazy.FromSlice([]rune("def"))
    wxyz := lazy.FromSlice([]rune("wxyz"))

    expected := []string{
        "adw", "bex", "cfy",
    }

    runesToString := func (xs []rune) string { return string(xs) }

    result := lazy.ToSlice(lazy.Map(runesToString, lazy.Zip(abc, def, wxyz)))

    assert.Equal(t, expected, result)
}

func TestPairwise(t *testing.T) {
    abcd := lazy.FromSlice([]rune("abcd"))

    expected := [][2]rune{
        {'a', 'b'},
        {'b', 'c'},
        {'c', 'd'},
    }

    result := lazy.ToSlice(lazy.Pairwise(abcd))

    assert.Equal(t, expected, result)
}


func TestPairwiseFill(t *testing.T) {
    abcd := lazy.FromSlice([]rune("abcd"))

    expected := [][2]rune{
        {'a', 'b'},
        {'b', 'c'},
        {'c', 'd'},
        {'d', 0xFFFF},
    }

    result := lazy.ToSlice(lazy.PairwiseFill(0xFFFF, abcd))

    assert.Equal(t, expected, result)
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
