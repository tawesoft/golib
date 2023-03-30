// Package iter implements iteration over sequences, including lazy evaluation.
//
// To avoid confusion, between the higher-order function "map" and the Go map
// data structure, the former will always be referred to as a "map function"
// and the latter will always be referred to as a "map collection".
package iter

import (
    "strings"
    "unicode/utf8"

    "github.com/tawesoft/golib/v2/operator"
    "github.com/tawesoft/golib/v2/operator/checked"
    "golang.org/x/exp/constraints"
    "golang.org/x/exp/maps"
)

// Pair is any Key, Value pair. Type K is any type that would be suitable as a
// KeyType in a map collection.
type Pair[K comparable, V any] struct {
    Key   K
    Value V
}

// It is an iterator, defined as a function that produces a (possibly infinite)
// sequence of (value, true) tuples through successive calls. If the sequence
// has ended then returned tuple is always (undefined, false).
//
// In this package, an iterator returning a sequence of (value, true) tuples is
// said to be "producing" values. An iterator returning (undefined, false) is
// said to be "exhausted". The sequence of values produced by an iterator is
// said to go from "left to right" where the leftmost value is the first one
// returned, and the rightmost value is the last one returned. Something that
// causes an iterator to produce values is called a "consumer". The number of
// values produced is the length of an iterator. Iterators may have infinite
// length.
//
// For example, if it := It.FromSlice([]int{1, 2, 3}), then it is an iterator
// of length three that produces the sequence 1, 2, 3. The leftmost value is 1.
// The rightmost value is 3. The first call to it() returns (1, true) and is
// said to have produced the value 1, the second call to it() returns (2, true)
// and is said to have produced the value 2, the third call to it() returns (3,
// true) and is said to have produced the value 3, the fourth call to it()
// returns (0, false) and the iterator is said to be exhausted, and successive
// calls to it() while exhausted continue to always return (0, false). A second
// iterator, e.g. it2 := It.Take(5, it1), is an example of an iterator that
// consumes a production of the input iterator (in this case, it1) whenever it2
// produces values.
type It[X any] func()(X, bool)

// Next implements the proposed Iter interface.
func (it It[X]) Next()(X, bool) {
    return it()
}

// All calls function f for each value x produced by an iterator until either f
// returns false, or the iterator is exhausted. All returns true iff f(x) was
// true for every x. If the iterator produced no values (was empty), then
// All still returns true.
func All[X any](
    f func(X) bool,
    it It[X],
) bool {
    for {
        x, ok := it()
        if !ok { break }
        if !f(x) { return false }
    }
    return true
}

// Any calls function f for each value x produced by an iterator until either
// f(x) returns true, or the iterator is exhausted. It returns true iff f(x)
// was true for at least one x.
func Any[X any](
    f func(X) bool,
    it It[X],
) bool {
    for {
        x, ok := it()
        if !ok { break }
        if f(x) { return true }
    }
    return false
}

// AppendToSlice appends every value produced by an iterator to dest, and
// returns the modified dest (like the builtin append).
func AppendToSlice[X any](
    dest []X,
    xs It[X],
) []X {
    for {
        x, ok := xs()
        if !ok { break }
        dest = append(dest, x)
    }
    return dest
}

// Cat (for "concatenate") returns an iterator that merges several input
// iterators, consuming an input iterator in its entirety to produce values
// before moving on to the next input iterator. The input iterators should not
// be used anywhere else once provided to this function.
//
// For example, given an iterator abc that produces the letters 'a', 'b', 'c',
// and an iterator def that produces the letters 'd', 'e', 'f', Cat(abc, def)
// will return an iterator that consumes abc and def to produce the letters
// 'a', 'b', 'c', 'd', 'e', 'f'.
//
// For an iterator that produces values from each input in lockstep, consuming
// one value from each iterator in turn, see [Zip] or [ZipFlat].
//
// Some libraries call this function "chain" instead.
func Cat[X any](its ... It[X]) It[X] {
    zero := operator.Zero[X]()
    return func() (X, bool) {
        for {
            if len(its) == 0 { return zero, false }
            if x, ok := its[0](); ok {
                return x, true
            } else {
                its = its[1:]
            }
        }
    }
}

// Check calls function f for each value produced by an iterator. It halts when
// a non-nil error is returned by f, and immediately returns the value being
// examined at the time and the error. Otherwise, returns a zero value and a
// nil error.
//
// See also [Walk], which is similar but does not check for errors.
func Check[X any](
    f func(X) error,
    it It[X],
) (X, error) {
    zero := operator.Zero[X]()
    for {
        x, ok := it()
        if !ok { break }
        if err := f(x); err != nil { return x, err }
    }
    return zero, nil
}

// Count calls function f for each value x produced by an iterator. It returns
// a 3-tuple containing a count of the number of times f(x) returned true,
// a count of the number of times f(x) returned false, and the total number
// of times f(x) was called. If f is nil, acts as if f(x) => true for all x.
func Count[X any](
    f func(X) bool,
    it It[X],
) (numTrue, numFalse, total int) {
    for {
        x, ok := it()
        if !ok { break }
        if (f == nil) || (f(x)) {
            numTrue++
        } else {
            numFalse++
        }
        total++
    }
    return
}

// Counter returns an iterator that produces a series of integers, starting at
// the given number, and increasing by step each time. It terminates at the
// maximum representable value for the number type.
func Counter[I constraints.Integer](start I, step I) It[I] {
    limit := checked.GetLimits[I]()
    done := false
    return func() (I, bool) {
        if done { return 0, false }
        if x, ok := limit.Add(start, step); ok {
            result := start
            start = x
            return result, true
        } else {
            done = true
            return start, true
        }
    }
}

// Cut returns an iterator that splits the input slice on an element
// TODO

// CutString is like [Cut], but operates on a string, producing strings
// delimited by a separator rune.
//
// For example, CutString("a|b|c", "|") returns an iterator that produces
// the strings "a", "b", "c".
//
// Unlike some Go functions, CutString does not also treat invalid Unicode or
// invalid Utf8 byte sequences as a valid delimiter when the separator is
// utf8.RuneError. That is, the utf8.RuneError delimiter only matches bytes
// that literally encode the Unicode value of utf8.RuneError.
//
// Note that in many situations, it is probably faster, simpler, clearer, and
// more idiomatic to use a [bufio.Scanner]. Prefer CutString only if you are
// combining the result iterator with other higher-order functions in this
// package.
func CutString(in string, sep rune) It[string] {
    done := false
    z := utf8.RuneLen(sep)
    return func() (string, bool) {
        if done { return "", false }
        if idx := stringsIndexRune(in, sep); idx < 0 {
            done = true
            return in, true
        } else {
            left := in[0:idx]
            right := in[idx+z:]
            in = right
            return left, true
        }
    }
}

// stringsIndexRune reimplements strings.IndexRune because that function can't
// capture utf8.RuneError.
func stringsIndexRune(in string, sep rune) int {
    offset := 0
    z := utf8.RuneLen(sep)
    for _, r := range in {
        rZ := utf8.RuneLen(r)

        if (r == sep) && (rZ == z) {
            return offset
        }

        offset += rZ
    }
    return -1
}

// CutStringStr is like [CutString], but the delimiter is a string, not just
// a single rune.
func CutStringStr(in string, sep string) It[string] {
    done := false
    return func() (string, bool) {
        if done { return "", false }
        if left, right, found := strings.Cut(in, sep); !found {
            done = true
            return in, true
        } else {
            in = right
            return left, true
        }
    }
}

// Empty returns an iterator that is typed, but empty.
func Empty[X any]() It[X] {
    zero := operator.Zero[X]()
    return func () (X, bool) {
        return zero, false
    }
}

// Enumerate produces a [Pair] for each value produced by the input iterator,
// where Pair.Key is an integer that starts at zero and increases by one with
// each produced value. The input iterator should not be used anywhere else
// once provided to this function.
//
// For example, for an iterator abc that produces the runes 'a', 'b', 'c',
// Enumerate(abc) produces the values Pair[0, 'a'], Pair[1, 'b'], Pair[2, 'c'].
func Enumerate[X any, Y Pair[int, X]](
    it It[X],
) It[Y] {
    var n int
    return func() (Y, bool) {
        if x, ok := it(); ok {
            y := Y{Key: n, Value: x}
            n++
            return y, true
        }
        return Y{}, false
    }
}

// Exhaust consumes an iterator and discards the produced values. As a special
// case, a nil iterator safely returns immediately.
func Exhaust[X any](it It[X]) {
    if it == nil { return }
    for {
        if _, ok := it(); !ok { break }
    }
}

// Filter returns an iterator that consumes an input iterator and only
// produces those values where the provided filter function f returns true.
//
// As a special case, if f is nil, it is treated as the function f(x) => true.
//
// For example:
//
//     function isOdd := func(x int) bool { return x % 2 == 1 }
//     Filter(isOdd, FromSlice([]int{1, 2, 3})) // produces 1 and 3
//
func Filter[X any](
    f func(X) bool,
    it It[X],
) It[X] {
    zero := operator.Zero[X]()

    if f == nil {
        f = func(x X) bool { return true }
    }

    return func() (X, bool) {
        for {
            x, ok := it()
            if !ok { return zero, false }
            if !f(x) { continue }
            return x, true
        }
    }
}

// Final produces a [FinalValue] for each value produced by the input iterator.
// FinalValue.IsFinal is true iff value is the last value produced before the
// input iterator would be exhausted. If the input iterator was exhausted at
// the time of input, then the returned iterator is empty. The input iterators
// should not be used anywhere else once provided to this function.
//
// For example, Final(FromSlice([]int{1, 2, 3}) produces the values
// (FinalValue{1, false}, true), (FinalValue{2, false}, true),
// (FinalValue{3, true}, true), (FinalValue{}, false).
func Final[X any](it It[X]) It[FinalValue[X]] {
    zero := FinalValue[X]{}
    done := false

    last, ok := it()
    if !ok {
        return func() (FinalValue[X], bool) { return zero, false }
    }

    return func() (FinalValue[X], bool) {
        if done { return zero, false }

        x, ok := it()
        if !ok {
            done = true
            return FinalValue[X]{last, true}, true
        }

        result := last
        last = x

        return FinalValue[X]{result, false}, true
    }
}

// FinalValue is a value returned by the [Final] function. Iff IsFinal.IsFinal
// is true, Value is the last value produced and the iterator is now exhausted.
type FinalValue[X any] struct {
    Value X
    IsFinal bool
}

// FromMap returns an iterator that produces each (key, value) pair from the
// input [builtin.Map] (of Go type map[X]Y, not to be confused with the higher
// order function [Map]) as an Pair. Do not modify the underlying map's keys
// until no longer using the returned iterator.
func FromMap[X comparable, Y any](kvs map[X]Y) It[Pair[X, Y]] {
    rest := maps.Keys(kvs)
    zero := operator.Zero[Pair[X, Y]]()

    return func() (Pair[X, Y], bool) {
        if len(rest) == 0 {
            return zero, false
        } else {
            key := rest[0]
            rest = rest[1:]
            return Pair[X, Y]{
                Key: key,
                Value: kvs[key],
            }, true
        }
    }
}

// FromSlice returns an iterator that produces each item in the input slice.
// Do not modify the underlying slice or backing array until no longer using
// the returned iterator.
func FromSlice[X any](xs []X) It[X] {
    rest := xs
    zero := operator.Zero[X]()
    return func() (X, bool) {
        if len(rest) == 0 {
            return zero, false
        } else {
            x := rest[0]
            rest = rest[1:]
            return x, true
        }
    }
}

// FromString returns an iterator that produces the runes of string s.
func FromString(s string) It[rune] {
    return FromSlice([]rune(s))
}

// Func performs a type cast that returns an iterator (type [It]) from any
// function meeting the iterator interface.
func Func[X any](f func() (X, bool)) It[X] {
    return f
}

// InsertToMap modifies a [builtin.Map] (of Go type map[X]Y, not to be confused
// with the higher order function [Map]). For each Pair produced by the input
// iterator, Pair.Key is used as a map key and Pair.Value is used as the
// matching value.
//
// If a key already exists in the destination map (either originally, or as a
// result of the input iterator producing two values with the same key), then
// the choose function is called with the key, the original value, and the new
// value. The return value of that function is used as the value to keep. If
// the choose function is nil, the new value in the case of conflicts is always
// the one most recently produced by the iterator.
func InsertToMap[X comparable, Y any](
    dest map[X]Y,
    choose func(key X, original Y, new Y) Y,
    kvs It[Pair[X,Y]],
) {
    for {
        kv, ok := kvs()
        if !ok { break }
        if original, exists := dest[kv.Key]; exists && (choose != nil) {
            dest[kv.Key] = choose(kv.Key, original, kv.Value)
        } else {
            dest[kv.Key] = kv.Value
        }
    }
}

func Keys[X comparable, Y any](xs It[Pair[X, Y]]) It[X] {
    return Map(func(i Pair[X, Y]) X { return i.Key }, xs)
}

func Values[X comparable, Y any](xs It[Pair[X, Y]]) It[Y] {
    return Map(func(i Pair[X, Y]) Y { return i.Value }, xs)
}

// Map returns an iterator that consumes an input iterator (of type X) and
// produces values (of type Y) for each input value, according to some mapping
// function f(x X) => y Y.
//
// For example:
//
//     double := func(i int) int { return i * 2 }
//     Map(double, FromSlice([int]{1, 2, 3})) // produces 2, 3, 6.
//
// For example, changing the type of the result:
//
//     stringify := func(i int) string { return fmt.Sprintf("%d", i) }
//     Map(stringify, FromSlice([int]{1, 2, 3})) // produces "1", "2", "3".
//
func Map[X any, Y any](
    f func(X) Y,
    it It[X],
) It[Y] {
    zero := operator.Zero[Y]()

    return func () (Y, bool) {
        x, ok := it()
        if !ok { return zero, false }
        return f(x), true
    }
}

// Pairwise returns an iterator that produces overlapping pairs of values
// produced by the input. The input iterator should not be used anywhere else
// once provided to this function.
//
// For example, for an iterator abc that produces the runes 'a', 'b', 'c',
// Pairwise(abc) produces the pairs [2]rune{'a', b'} and [2]rune{'b', 'c'}.
func Pairwise[X any, Y [2]X](
    it It[X],
) It[Y] {
    gs := Tee(2, it)
    g0, g1 := gs[0], gs[1]

    g1() // discard in order to advance the second iterator

    return func() (Y, bool) {
        a, ok := g0()
        if !ok { return [2]X{}, false }

        b, ok := g1()
        if !ok { return [2]X{}, false }

        return [2]X{a, b}, true
    }
}

// PairwiseEnd is like [Pairwise], however a final result pair [2]X{value,
// lastValue} is produced for the last value.
//
// For example, for an iterator abc that produces the runes 'a', 'b', 'c',
// PairwiseFill(0, abc) produces the pairs [2]rune{'a', b'}, [2]rune{'b', 'c'}
// and [2]rune{'c', 0}.
func PairwiseEnd[X any, Y [2]X](
    lastValue X,
    g It[X],
) It[Y] {
    return Pairwise[X, Y](Cat(g, Repeat(1, lastValue)))
}

// Reduce applies function f to each element of the sequence (in order) with
// the previous return value from f as the first argument and the element as
// the second argument (for the first call to f, the supplied initial value is
// used instead of a return value), returning the final value returned by f.
//
// If the input iterator is empty, the result is the initial value.
func Reduce[X any](
    initial X,
    f func(X, X) X,
    it It[X],
) X {
    v := initial
    for {
        x, ok := it()
        if !ok { return v }
        v = f(v, x)
    }
}

// Repeat produces the value x, with n repetitions. If n is negative, it
// produces x with infinite repetitions.
//
// For example, Repeat(3, "foo") produces the values "foo", "foo", "foo".
func Repeat[X any](
    n int,
    x X,
) It[X] {
    return func() (X, bool) {
        if n == 0 {
            return operator.Zero[X](), false
        }
        if n > 0 { n-- }
        return x, true
    }
}

// Take returns an iterator that produces only (up to) the first n items of
// the input iterator.
func Take[X any](
    n int,
    xs It[X],
) It[X] {
    zero := operator.Zero[X]()
    return func() (X, bool) {
        if n == 0 { return zero, false }
        x, ok := xs()
        n--
        return x, ok
    }
}

// Tee returns a slice of n iterators that each, individually, produce the
// same values otherwise produced by the input iterator. This can be thought
// of at "copying" an iterator. The input iterators should not be used anywhere
// else once provided to this function.
//
// For example, given an iterator abc that produces the letters "a", "b", "c",
// Tee will return n iterators. Each returned iterator will, independently,
// produce the letters "a", "b", "c".
//
// Where the returned iterators produces their inputs at different speeds
// (their consumers are "out of step"), this requires growing amounts of
// auxiliary storage.
func Tee[X any](
    n int,
    g It[X],
) []It[X] {

    // TODO/PERF this could be more efficient.

    gs := make([]It[X], n)
    queues := make([][]X, n)

    for i := 0; i < n; i++ {
        queues[i] = make([]X, 0) // FIFO

        gs[i] = func() It[X] {
            queue := &queues[i]

            return (func() (X, bool) {

                // empty queue? produce more values
                if len(*queue) == 0 {
                    if x, ok := g(); ok {
                        // send to all queues
                        for j := 0; j < n; j++ {
                            queues[j] = append(queues[j], x)
                        }
                    } else {
                        // exhausted
                        return x, ok
                    }
                }

                // pop from front of queue
                item := (*queue)[0]
                *queue = (*queue)[1:]
                return item, true
            })
        }()
    }

    return gs
}

// ToMap returns a new [builtin.map] (of Go type map[X]Y, not to be confused
// with the higher order function [Map]). For each Pair produced by the input
// iterator, Pair.Key is used as a map key and Pair.Value is used as the
// matching value. If the iterator produces two Items with the same Pair.Key,
// only the last value is kept.
//
// If a key already exists in the destination map (as a result of the input
// iterator producing two values with the same key), then the choose function
// is called with the key, the original value, and the new value. The return
// value of that function is used as the value to keep. If the choose function
// is nil, the new value in the case of conflicts is always the one most
// recently produced by the iterator.
func ToMap[X comparable, Y any](
    choose func(key X, original Y, new Y) Y,
    kvs It[Pair[X,Y]],
) map[X]Y {
    result := make(map[X]Y)
    InsertToMap(result, choose, kvs)
    return result
}

// ToSlice returns a slice of every value produced by an iterator. If the input
// iterator is exhausted, returns an empty slice (not nil).
func ToSlice[X any](xs It[X]) []X {
    result := make([]X, 0)
    return AppendToSlice(result, xs)
}

// ToString returns a string from an iterator that produces runes.
func ToString(it It[rune]) string {
    return string(ToSlice[rune](it))
}

// Walk calls a visitor function f for each value produced by an iterator.
//
// See also [Check], which is like Walk but aborts on error, and [WalkFinal], which
// is like Walk but the last item produced by the input iterator is detected.
func Walk[X any](
    f func (X),
    it It[X],
) {
    for {
        x, ok := it()
        if !ok { break }
        f(x)
    }
}

// WalkFinal is like [Walk], but the second argument to the visitor function is
// true if x is the last item to be produced by an iterator.
func WalkFinal[X any](
    f func(X, bool),
    it It[X],
) {
    final := Final(it)
    for {
        x, ok := final()
        if !ok { break }
        f(x.Value, x.IsFinal)
    }
}

// Zip returns an iterator that produces slices of the results of each input
// iterator, in lockstep, terminating when any input is exhausted. The input
// iterators should not be used anywhere else once provided to this function.
//
// For example, for an iterator abc that produces the runes 'a', 'b', 'c', and
// an input iterator wxyz that produces the runes 'w', 'x', 'y', 'z', Zip(abc,
// wxyz) produces the values []rune{'a', 'w'}, []rune{'b', 'x'}, []rune{'c',
// 'y'} before becoming exhausted.
//
// If zipping multiple different types together, you will need to use
// iterators of type It[any].
func Zip[X any, Y []X](
    its ... It[X],
) It[Y] {
    n := len(its)

    if n == 0 { return func() (Y, bool) {
        return []X{}, false }
    }

    return func() (Y, bool) {
        result := make([]X, n)

        for i := 0; i < n; i++ {
            x, ok := its[i]()
            if !ok { return []X{}, false }
            result[i] = x
        }

        return result, true
    }
}

// ZipFlat returns an iterator that produces the results of each input
// iterator, in turn, terminating when any input is exhausted. The input
// iterators should not be used anywhere else once provided to this function.
//
// For example, for an iterator abc that produces the runes 'a', 'b', 'c', and
// an input iterator wxyz that produces the runes 'w', 'x', 'y', 'z',
// ZipFlat(abc, wxyz) produces the runes 'a', 'w', 'b', 'x', 'c', 'y' before
// becoming exhausted.
//
// If zipping multiple different types together, you will need to use
// iterators of type It[any].
//
// Some libraries call this function "round-robin" instead.
func ZipFlat[X any](
    its ... It[X],
) It[X] {
    zero := operator.Zero[X]()
    n := len(its)
    i := 0

    if n == 0 { return func() (X, bool) {
        return zero, false }
    }

    return func() (X, bool) {
        x, ok := its[i]()
        if !ok { return zero, false }

        i++
        if i == n { i = 0 }

        return x, true
    }
}
