package slice

import (
    "fmt"
)

// multiplier(m) returns the function f(x) => x * m
func multiplier(m int) (func (int) int) {
    return func (x int) int {
        return x * m
    }
}

// isOdd(x) returns true iff x is odd
func isOdd(x int) bool {
    return x % 2 != 0
}

// lessThan(m) returns the function f(x) => true iff x < m
func lessThan(m int) (func (int) bool) {
    return func (x int) bool {
        return x < m
    }
}

// toString(x) returns the string "%d", inserting x as "%d"
func toString(x int) string {
    return fmt.Sprintf("(%d)", x)
}

// join(x) returns a reducer that joins a slice of strings together
func join(sep string) SliceReducer[string] {
    return SliceReducer[string]{
        Identity: "",
        Reduce: func(a string, b string) string {
            if a == "" { return b }
            return a + sep + b
        },
    }
}

func Example() {
    xs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

    xs = FilterSliceInPlace(lessThan(8), xs)
    fmt.Printf("%+v\n",
        ReduceSlice(join(", "),
            MapSlice(toString,
                FilterSlice(isOdd,
                    MapSlice(multiplier(3), xs)))))

    // Output:
    // (3), (9), (15), (21)
}
