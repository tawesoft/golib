// x-doc-short-desc: lazy evaluation
// x-doc-stable: candidate

// Package lazy implements lazy evaluation (strictly typed with Go generics).
//
// To avoid confusion, between the higher-order function "map" and the Go map
// data structure, the latter is referred to as a [builtin.map] in this
// package.
package lazy

// CONTRIBUTORS: keep definitions in alphabetical order, but with types
// grouped first.

import (
    "github.com/tawesoft/golib/v2/ks"
    "golang.org/x/exp/maps"
)

// FinalValue is a value returned by the [Final] function. Iff IsFinal.IsFinal
// is true, Value is the last value produced and the iterator is now exhausted.
type FinalValue[X any] struct {
    Value X
    IsFinal bool
}

// Item is any Key, Value pair. Type K is any type that would be suitable as a
// KeyType in a Go [builtin.map].
type Item[K comparable, V any] ks.Item[K, V]

// It is an iterator, defined as a function that lazily produces a (possibly
// infinite) sequence of (value, true) tuples through successive calls. If the
// sequence has ended then returned tuple is always (undefined, false).
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
// iterator, e.g. it2 := It.TakeN(5, it1), is an example of an iterator that
// consumes a production of the input iterator (in this case, it1) whenever it2
// produces values.
type It[X any] func()(X, bool)

// Reducer is a function and an identity value for reducing a sequence of
// values into a single value by calling v = Reduce(v, x) for each x of an
// input sequence from left to right.
//
// The first argument to the first invocation of the Reduce function is always
// the provided Identity value which is defined so that Reduce(Identity, x)
// always returns x.
//
// Simple example Reducers:
//
//     sum := Reducer{Reduce: func(a int, b int) { return a + b }, Identity: 0}
//     mul := Reducer{Reduce: func(a int, b int) { return a * b }, Identity: 1}
//
type Reducer[X any] struct {
    Reduce func(a X, b X) X
    Identity X
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
        if !ok { break; }
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
        if !ok { break; }
        if f(x) { return true }
    }
    return false
}

// AppendToSlice appends every value produced by an iterator to dest, and
// returns the modified dest (like [builtin.append]).
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
    zero := ks.Zero[X]()
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
func Check[X any](
    f func(X) error,
    it It[X],
) (X, error) {
    zero := ks.Zero[X]()
    for {
        x, ok := it()
        if !ok { break }
        if err := f(x); err != nil { return x, err }
    }
    return zero, nil
}

// Enumerate produces an [Item] for each value produced by the input iterator,
// where Item.Key is an integer that starts at zero and increases by one with
// each produced value. The input iterator should not be used anywhere else
// once provided to this function.
//
// For example, for an iterator abc that produces the runes 'a', 'b', 'c',
// Enumerate(abc) produces the values Item[0, 'a'], Item[1, 'b'], Item[2, 'c'].
func Enumerate[X any, Y Item[int, X]](
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

// Filter returns an iterator that consumes an input iterator and only
// produces those values where the provided filter function f returns true.
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
    zero := ks.Zero[X]()

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
// For example, Final(FromSlice([]int{1, 2, 3}) produces the values FinalValue{1, false},
// FinalValue{2, false}, FinalValue{3, true}
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

// FromMap returns an iterator that produces each (key, value) pair from the
// input [builtin.Map] (of Go type map[X]Y, not to be confused with the higher
// order function [Map]) as an Item. Do not modify the underlying map's keys
// until no longer using the returned iterator.
func FromMap[X comparable, Y any](kvs map[X]Y) It[Item[X, Y]] {
    rest := maps.Keys(kvs)
    zero := ks.Zero[Item[X, Y]]()

    return func() (Item[X, Y], bool) {
        if len(rest) == 0 {
            return zero, false
        } else {
            key := rest[0]
            rest = rest[1:]
            return Item[X, Y]{
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
    zero := ks.Zero[X]()
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
// with the higher order function [Map]). For each Item produced by the input
// iterator, Item.Key is used as a map key and Item.Value is used as the
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
    kvs It[Item[X,Y]],
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

// Join is similar to [Reduce], but without using an identity value. It applies
// an operation between successive productions of an iterator.
//
// An empty iterator joins to the zero value, an iterator of length one joins
// to the single value produced, and an iterator of length greater than one
// joins v = f(v, x) for each value x it produces after the first, where v is
// the first value.
//
// For example,
//
//     sum := func(a int, b int) int { return a + b }
//     Join(sum, FromSlice([]int{}) // returns 0
//     Join(sum, FromSlice([]int{123}) // returns 123
//     Join(sum, FromSlice([]int{1, 2, 3})
//     // calculates ((1 + 2) + 3) and returns 6
//
// Note, this is a terrible way to join strings!
func Join[X any](
    f func(a X, b X) X,
    it It[X],
) X {
    zero := ks.Zero[X]()

    v, ok := it()
    if !ok { return zero }

    for {
        x, ok := it()
        if !ok { return v }
        v = f(v, x)
    }
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
    zero := ks.Zero[Y]()

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

// Reduce returns the result of applying a [Reducer] to the elements produced
// by an iterator.
//
// For example,
//
//     sum := func(a int, b int) { return a + b }
//     summer := Reducer{Reduce: sum, Identity: 0}
//     Reduce(summer, FromSlice([]int{1, 2, 3})
//     // calculates (((0 + 1) + 2) + 3) and returns 6
//
func Reduce[X any](
    reducer Reducer[X],
    it It[X],
) X  {
    v := reducer.Identity
    for {
        x, ok := it()
        if !ok { return v }
        v = reducer.Reduce(v, x)
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
            return ks.Zero[X](), false
        }
        if n > 0 { n-- }
        return x, true
    }
}

// TakeN returns an iterator that produces only (up to) the first n items of
// the input iterator.
func TakeN[X any](
    n int,
    xs It[X],
) It[X] {
    zero := ks.Zero[X]()
    return func() (X, bool) {
        if n == 0 { return zero, false }
        x, ok := xs()
        n--
        return x, ok
    }
}

// Tee returns a slice of n iterators that each, individually, produce the
// sames values otherwise produced by the input iterator. This can be thought
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
// with the higher order function [Map]). For each Item produced by the input
// iterator, Item.Key is used as a map key and Item.Value is used as the
// matching value. If the iterator produces two Items with the same Item.Key,
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
    kvs It[Item[X,Y]],
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

// Walk calls function f for each value produced by an iterator.
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
// an input iterator wxyz that produces the runes 'w', 'x', 'y', 'z', Zip(abc,
// wxyz) produces the runes 'a', 'w', 'b', 'x', 'c', 'y' before becoming
// exhausted.
//
// If zipping multiple different types together, you will need to use
// iterators of type It[any].
//
// Some libraries call this function "round-robin" instead.
func ZipFlat[X any](
    its ... It[X],
) It[X] {
    zero := ks.Zero[X]()
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
