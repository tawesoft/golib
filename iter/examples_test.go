package iter_test

import (
    "fmt"
    "strconv"
    "strings"

    lazy "github.com/tawesoft/golib/v2/iter"
)

func ExampleFromMap() {
    type Person struct {
        name string
        age int
    }
    type ID string

    // a map of people
    people := lazy.FromMap(map[ID]Person{
        "ATUR001": {name: "Alice Turing", age: 23},
        "GHOP001": {name: "George Hopper", age: 60},
        "FKAH001": {name: "Freddy Kahlo", age: 29},
    })

    // this filter function returns true for people under thirty
    underThirty := func(kv lazy.Pair[ID, Person]) bool {
        return kv.Value.age < 30
    }

    // apply the filter and finally generate a dict
    peopleUnderThirty := lazy.ToMap(nil, lazy.Filter(underThirty, people))

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

func ExampleFunc() {
    // generate an infinite sequence of integers with a function
    integers := func() lazy.It[int] { // or func() func() (int, bool) {
        i := 0
        return func() (int, bool) {
            result := i
            i = i + 1
            return result, true
        }
    }

    integerGenerator := lazy.Func(integers())
    firstFour := lazy.Take(4, integerGenerator)

    fmt.Printf("First four integers are: %v\n", lazy.ToSlice(firstFour))

    // Output:
    // First four integers are: [0 1 2 3]
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

func ExampleToMap() {
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
    personToTuple := func (person Person) lazy.Pair[ID, Person] {
        return lazy.Pair[ID, Person]{person.id, person}
    }

    // apply the function over all people (lazily...)
    peopleTuples := lazy.Map(personToTuple, people)

    // finally generate a dict
    peopleByID := lazy.ToMap(nil, peopleTuples)

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

func ExampleWalk() {
    var sb strings.Builder

    strings := lazy.FromSlice([]string{"one", "two", "three"})

    lazy.Walk(func(x string) { sb.WriteString(x) }, strings)

    fmt.Println(sb.String())

    // Output:
    // onetwothree
}
