// x-doc-short-desc: lazy evaluation of sequences
// x-doc-stable: candidate

// Package lazy implements lazy evaluation (strictly typed with Go generics).
//
// To avoid confusion, between the higher-order function "map" and the Go map
// data structure, the latter is referred to as a [builtin.map] in this
// package.
package lazy

import (
    "github.com/tawesoft/golib/v2/ks"
    "golang.org/x/exp/maps"
)

// Item is any Key, Value pair. Type K is any type that would be suitable as a
// KeyType in a Go [builtin.map].
type Item[K comparable, X any] struct {
    Key   K
    Value X
}

// It is a function anything that lazily produces a (possibly infinite) sequence of
// (value, ok) tuples through successive calls to Next(). If the sequence has
// ended then in the returned tuple, value is undefined and ok is false.
//
// For example, if i := It.FromSlice([]int{1, 2, 3}), then i.Next()
// returns (1, true), then (2, true), then (3, true), then (0, false),
// and successive calls to i.Next() continue to return (0, false).
type Iterator[X any] interface {
    Next() (X, bool)
}

// IteratorFunc is a type for a function that can be used as an [It].
// Typing a function this way is optional, but sometimes clearer.
type IteratorFunc[X any] func()(X, bool)
func (f IteratorFunc[X]) Next() (X, bool) { return f() }

// Func returns an iterator from a function (or an [IteratorFunc]).
func Func[X any](f func() (X, bool)) Iterator[X] {
    return IteratorFunc[X](f)
}

// ToSlice returns a slice of every value produced by a generator.
func ToSlice[X any](xs Iterator[X]) []X {
    result := make([]X, 0)
    return AppendToSlice(result, xs)
}

// AppendToSlice appends every value produced by a generator to dest, and
// returns dest (like [builtin.Append]).
func AppendToSlice[X any](dest []X, xs Iterator[X]) []X {
    for {
        x, ok := xs.Next()
        if !ok { break }
        dest = append(dest, x)
    }
    return dest
}

// ToDict returns a dict (of Go type map[X]Y) for each Item produced by a
// generator
func ToDict[X comparable, Y any](kvs Iterator[Item[X,Y]]) map[X]Y {
    result := make(map[X]Y)
    for {
        kv, ok := kvs.Next()
        if !ok { break }
        result[kv.Key] = kv.Value
    }
    return result
}

// DictInsert modifies a dict (of Go type map[X]Y) to insert each Item
// produced by a generator. If a key already exists in the destination dict,
// the given function is called with the original value and the new value, and
// the return value is used at the new valur. If the given function is nil, the
// new value is always the one created by the generator.
func DictInsert[X comparable, Y any](
    dest map[X]Y,
    choose func(key X, original Y, new Y) Y,
    kvs Iterator[Item[X,Y]],
) map[X]Y {
    for {
        kv, ok := kvs.Next()
        if !ok { break }
        if original, exists := dest[kv.Key]; exists && (choose != nil) {
            dest[kv.Key] = choose(kv.Key, original, kv.Value)
        } else {
            dest[kv.Key] = kv.Value
        }
    }
    return dest
}

type taker[X any] struct {
    g Iterator[X]
    n int
}

func (g *taker[X]) Next() (X, bool) {
    if g.n == 0 { return ks.Zero[X](), false }
    x, ok := g.g.Next()
    g.n--
    return x, ok
}

// TakeN returns a generator that only generates the first n items of the
// input generator.
func TakeN[X any](n int, xs Iterator[X]) Iterator[X] {
    return &taker[X]{
        g: xs,
        n: n,
    }
}

type slice[X any] struct {
    rest []X
}

func (g *slice[X]) Next() (X, bool) {
    if len(g.rest) == 0 {
        return ks.Zero[X](), false
    } else {
        x := g.rest[0]
        g.rest = g.rest[1:]
        return x, true
    }
}

// FromSlice returns a generator that generates each item in the input slice.
// Do not modify the underlying slice or backing array until no longer using
// the generator.
func FromSlice[X any](xs []X) Iterator[X] {
    return &slice[X]{xs}
}

type dict[X comparable, Y any] struct {
    dict map[X]Y
    keys []X
}

func (g *dict[X, Y]) Next() (Item[X, Y], bool) {
    if len(g.keys) == 0 {
        return ks.Zero[Item[X, Y]](), false
    } else {
        key := g.keys[0]
        g.keys = g.keys[1:]
        return Item[X, Y]{
            Key: key,
            Value: g.dict[key],
        }, true
    }
}

// FromDict returns a generator that generates each (key, value) pair from the
// input dict as an Item. Do not modify the underlying dict's keys until no
// longer using the generator.
func FromDict[X comparable, Y any](kvs map[X]Y) Iterator[Item[X, Y]] {
    return &dict[X, Y]{
        dict: kvs,
        keys: maps.Keys(kvs),
    }
}

type filter[X any] struct {
    g Iterator[X]
    f func(X) bool
}

func (g *filter[X]) Next() (X, bool) {
    for {
        x, ok := g.g.Next()
        if !ok { return x, ok }

        if g.f(x) {
            return x, true
        }
    }
}

// Filter returns a generator that only generates values where the provided
// filter function f returns true.
func Filter[X any](
    f func(X) bool,
    g Iterator[X],
) Iterator[X] {
    return &filter[X]{
        g: g,
        f: f,
    }
}

type mapper[X any, Y any] struct {
    g Iterator[X]
    f func(X) Y
}

func (g *mapper[X, Y]) Next() (Y, bool) {
    for {
        x, ok := g.g.Next()
        if !ok { return ks.Zero[Y](), ok }

        y := g.f(x)
        return y, true
    }
}

// Map returns a generator that generates values of type Y for each input value
// of type X.
func Map[X any, Y any](
    f func(X) Y,
    g Iterator[X],
) Iterator[Y] {
    return &mapper[X, Y]{
        g: g,
        f: f,
    }
}

// Reducer is a function and an identity value (see [Reduce]) for reducing
// a sequence of values into a single value.
type Reducer[X any] struct {
    Reduce func(a X, b X) X
    Identity X
}

// Reduce generates the result of applying a reduce function from left to
// right pairwise on the elements of a generator. The first argument to the
// first invocation of the reduction function is the identity value where
// reduce(identity, x) == x.
func Reduce[X any](
    reducer Reducer[X],
    g Iterator[X],
) X  {
    v := reducer.Identity
    for {
        x, ok := g.Next()
        if !ok { break }

        v = reducer.Reduce(v, x)
    }
    return v
}

// Walk calls function f for each value generated by a generator.
func Walk[X any](
    f func (X),
    g Iterator[X],
) {
    for {
        x, ok := g.Next()
        if !ok { break }
        f(x)
    }
}

// Check calls function f for each value of type X generated by a generator. If
// f returns nil, Check continues normally through each element. Otherwise, f
// returns a function with no arguments that returns a value of type Y when
// called (e.g. an error message), and Check terminates and returns the
// function returned by f.
func Check[X any, Y any](
    f func(X) (func () Y),
    g Iterator[X],
) (func () Y) {
    for {
        x, ok := g.Next()
        if !ok { break }
        if y := f(x); y != nil { return y }
    }
    return nil
}

// Any calls function f for each value x generated by a generator until either
// f(x) returns true, or the generator is exhausted. It returns true iff f(x)
// was true for at least one x.
func Any[X any](
    f func(X) bool,
    g Iterator[X],
) bool {
    for {
        x, ok := g.Next()
        if !ok { break; }
        if f(x) { return true }
    }
    return false
}

// All calls function f for each value x generated by a generator until either f
// returns false, or the generator is exhausted. It returns true iff f(x) was
// true for every x. In the case of an empty generator, always returns true.
func All[X any](
    f func(X) bool,
    g Iterator[X],
) bool {
    for {
        x, ok := g.Next()
        if !ok { break; }
        if f(x) { return true }
    }
    return true
}

// Cat returns a generator that generates values from each generator until
// exhausted before moving on to the next generator. The input generators
// should not be used anywhere else once provided to this function.
//
// For example, given a generator abc that generates the letters "a", "b", "c",
// and a generator def that generates the letters "d", "e", "f", Cat(abc, def)
// will return a generator that generates the letters "a", "b", "c", "d", "e",
// "f".
func Cat[X any](
    gs ... Iterator[X],
) Iterator[X] {
    g := gs[:]
    return Function[X](func() (X, bool){
        for {
            if len(g) == 0 { break }
            if x, ok := g[0].Next(); ok {
                return x, true
            } else {
                g = g[1:]
            }
        }
        return ks.Zero[X](), false
    })
}

// Tee returns a slice of n generators that each, individually, produce the
// sames values otherwise produced by the input generator. You can think of
// this as "copying" a generator. The input generators should not be used
// anywhere else once provided to this function.
//
// Where the returned generators generate their inputs at different speeds (are
// "out of step"), this requires growing amounts of auxiliary storage.
//
// For example, given a generator abc that generates the letters "a", "b", "c",
// Tee will return n generators that also, independently, generate the letters
// "a", "b", "c".
func Tee[X any](
    n int,
    g Iterator[X],
) []Iterator[X] {

    // TODO/PERF this could be more efficient.

    gs := make([]Iterator[X], n)
    queues := make([][]X, n)

    for i := 0; i < n; i++ {
        queues[i] = make([]X, 0) // FIFO

        gs[i] = func() Iterator[X] {
            queue := &queues[i]

            return Function[X](func() (X, bool) {

                // empty queue? generate more values
                if len(*queue) == 0 {
                    if x, ok := g.Next(); ok {
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

// Zip returns a generator that generates slices of the results of each input
// generator, in lockstep. The input generators should not be used anywhere
// else once provided to this function.
//
// If the input generators are of different lengths, Zip terminates once
// any input generator is exhausted.
//
// For example, for a generator abc that generates the runes 'a', 'b', 'c',
// and a generator wxyz that generates the runes 'w', 'x', 'y', 'z',
// Zip(abc, wxyz) generates the values []rune{'a', 'w'}, []rune{'b', 'x'},
// []rune{'c', 'y'}.
//
// If zipping multiple different types together, you will need to use
// generators of type It[any].
func Zip[X any, Y []X](
    gs ... Iterator[X],
) Iterator[Y] {
    n := len(gs)
    result := make([]X, n)

    return Function[Y](func() (Y, bool) {
        for i := 0; i < n; i++ {
            x, ok := gs[i].Next()
            if !ok { return result, false }
            result[i] = x
        }
        return result, true
    })
}

// Pairwise returns a generator that generates overlapping pairs of values
// generated by the input generator. The input generator should not be used
// anywhere else once provided to this function.
//
// For example, for a generator abc that generates the runes 'a', 'b', 'c',
// Pairwise(abc) generates the pairs [2]rune{'a', b'} and [2]rune{'b', 'c'}.
func Pairwise[X any, Y [2]X](
    g Iterator[X],
) Iterator[Y] {
    gs := Tee(2, g)
    g0, g1 := gs[0], gs[1]

    g1.Next() // discard to advance second stream

    return Function[Y](func() (Y, bool) {
        a, ok := g0.Next()
        if !ok { return [2]X{}, false }

        b, ok := g1.Next()
        if !ok { return [2]X{}, false }

        return [2]X{a, b}, true
    })
}

// PairwiseFill is like [Pairwise], however a final result pair [2]X{value,
// fillValue} is generated for the last value generated by g.
//
// For example, for a generator abc that generates the runes 'a', 'b', 'c',
// PairwiseFill(0, abc) generates the pairs [2]rune{'a', b'}, [2]rune{'b', 'c'}
// and [2]rune{'c', 0}.
func PairwiseFill[X any, Y [2]X](
    fillValue X,
    g Iterator[X],
) Iterator[Y] {
    return Pairwise[X, Y](Cat(g, RepeatN(1, fillValue)))
}

// RepeatN generates the value x, with n repetitions. If n is negative,
// generates x with infinite repetitions.
//
// For example, RepeatN(3, "foo") generates the values "foo", "foo", "foo".
func RepeatN[X any](
    n int,
    x X,
) Iterator[X] {
    return Function[X](func() (X, bool) {
        if n == 0 {
            return ks.Zero[X](), false
        }
        if n > 0 { n-- }
        return x, true
    })
}

// Enumerate generates an Item tuple (index, value) for each value generated by
// the input generator, where index is an integer that starts at zero and
// increases by one with each generated value. The input generator should not
// be used anywhere else once provided to this function.
//
// For example, for a generator abc that generates the runes 'a', 'b', 'c',
// Enumerate(abc) generates the values Item[0, 'a'], Item[1, 'b'], Item[2,
// 'c'].
func Enumerate[X any, Y Item[int, X]](
    g Iterator[X],
) Iterator[Y] {
    var n int
    return Function[Y](func() (Y, bool) {
        if x, ok := g.Next(); ok {
            y := Y{Key: n, Value: x}
            n++
            return y, true
        }
        return Y{}, false
    })
}
