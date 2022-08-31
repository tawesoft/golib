// Package view provides customisable abstractions over collections. Changes to
// an underlying collection are reflected in its views.
package view

import (
    "github.com/tawesoft/golib/v2/iter"
    "github.com/tawesoft/golib/v2/ks"
)

// Key is a shorthand to create a new [Keyer].
func Key[A comparable, B comparable](AtoB func(A) B, BtoA func(B) A) Keyer[A, B] {
    return Keyer[A, B]{
        To: AtoB,
        From: BtoA,
    }
}

// Value is a shorthand to create a new [Valuer].
func Value[A any, B any](AtoB func(A) B, BtoA func(B) A) Valuer[A, B] {
    return Valuer[A, B]{
        To: AtoB,
        From: BtoA,
    }
}

// Keyer defines a mapping between comparable types A and B.
type Keyer[A comparable, B comparable] struct {
    To func(A) B
    From func(B) A
}

// Valuer defines a mapping between any types A and B.
type Valuer[A any, B any] struct {
    To func(A) B
    From func(B) A
}

type Iter[T any] interface {
    Next() (T, bool)
}

type Pair[K comparable, V any] struct {
    Key K
    Value V
}

type View[K comparable, V any] interface {
    Get(K) (V, bool)
    Set(K, V)
    Delete(K)
    Iter() Iter[Pair[K, V]]
}

type view[K comparable, V any] struct {
    get func(K) (V, bool)
    set func(K, V)
    delete func(K)
    iter func() Iter[Pair[K, V]]
}

// Mapper defines how a collection is mapped to and from a view. Call the View
// method to get a new view based on this map.
type Mapper[K comparable, V any, ToK comparable, ToV any] struct {
    // Filterer defines a function that "hides" the given value from the view
    // when accessed through Get or Iter. It is implemented in terms of the
    // types of the underlying collection, and should not implement the keyer
    // or valuer logic.
    Filterer func(K, V) bool

    // Getter defines a function that accesses a given value from the
    // underlying collection. It is implemented in terms of the types of the
    // underlying collection, and should not implement the filtering, keyer or
    // valuer logic.
    Getter func(K) (V, bool)

    // Setter defines a function that accesses a given value from the
    // underlying collection. It is implemented in terms of the types of the
    // underlying collection, and should not implement the filtering, keyer or
    // valuer logic. Note that a Setter ignores filtering.
    Setter func(K, V)

    // Deleter defines a function that deletes a value in the underlying
    // collection. It is implemented in terms of the types of the
    // underlying collection, and should not implement the filtering, keyer or
    // valuer logic. Note that a Deleter ignores filtering.
    Deleter func(K)

    // Keyer defines a mapping to and from the keys in the underlying
    // collection to the view. In the case of a list, the keys will be simple
    // integers.
    Keyer Keyer[K, ToK]

    // Valuer defines a mapping to and from the values in the underlying
    // collection to the view.
    Valuer Valuer[V, ToV]

    // Iterer defines a function that returns a new iterator over the
    // underlying collection. It is implemented in terms of the types of the
    // underlying collection, and should not implement the filtering, keyer or
    // valuer logic.
    Iterer func() Iter[Pair[K, V]]
}

// FromMap returns a View from a Go map type. See [Mapper] for details on
// the function arguments.
func FromMap[K comparable, V any, ToK comparable, ToV any](
    m map[K]V,
    filterer func(K, V) bool,
    Keyer Keyer[K, ToK],
    Valuer Valuer[V, ToV],
) View[ToK, ToV] {
    return Mapper[K, V, ToK, ToV]{
        Filterer: filterer,
        Getter: func(k K) (V, bool) {
            v, ok := m[k]
            return v, ok
        },
        Setter: func(k K, v V) {
            m[k] = v
        },
        Deleter: nil,
        Keyer: Keyer,
        Valuer: Valuer,
        Iterer: func() Iter[Pair[K, V]] {
            return IterPair(iter.FromMap(m))
        },
    }.View()
}

// FromSlice returns a View from a Go list type. See [Mapper] for details on
// the function arguments.
func FromSlice[V any, ToV any](
    s []V,
    filterer func(int, V) bool,
    Valuer Valuer[V, ToV],
) View[int, ToV] {
    return Mapper[int, V, int, ToV]{
        Filterer: filterer,
        Getter: func(k int) (V, bool) {
            ok := (k < len(s)) && (k >= 0)
            if !ok { return ks.Zero[V](), false }
            return s[k], true
        },
        Setter: func(k int, v V) {
            s[k] = v
        },
        Deleter: nil,
        Keyer: Key[int, int](ks.Identity[int], ks.Identity[int]),
        Valuer: Valuer,
        Iterer: func() Iter[Pair[int, V]] {
            return IterPair(iter.Enumerate(iter.FromSlice(s)))
        },
    }.View()
}

// IterPair returns an [Iter] of [Pair] elements from an [iter.It] of
// [iter.Item] elements.
func IterPair[K comparable, V any](it iter.It[iter.Item[K, V]]) Iter[Pair[K, V]] {
    return iter.Map(func (i iter.Item[K, V]) Pair[K, V] {
        return Pair[K, V]{i.Key, i.Value}
    }, it)
}

func (m Mapper[K, V, ToK, ToV]) get(key ToK) (ToV, bool) {
    k := m.Keyer.From(key)
    v, ok := m.Getter(k)
    if ok && m.Filterer(k, v) {
        return m.Valuer.To(v), true
    }
    return ks.Zero[ToV](), false
}

func (m Mapper[K, V, ToK, ToV]) set(k ToK, v ToV) {
    m.Setter(m.Keyer.From(k), m.Valuer.From(v))
}

func (m Mapper[K, V, ToK, ToV]) delete(k ToK) {
    m.Deleter(m.Keyer.From(k))
}

func (m Mapper[K, V, ToK, ToV]) iter() Iter[Pair[ToK, ToV]] {
    mapper := func(i Pair[K, V]) Pair[ToK, ToV] {
        return Pair[ToK, ToV]{m.Keyer.To(i.Key), m.Valuer.To(i.Value)}
    }

    filterer := func(i Pair[K, V]) bool {
        return m.Filterer(i.Key, i.Value)
    }

    return iter.Map[Pair[K, V], Pair[ToK, ToV]](mapper,
        iter.Filter[Pair[K, V]](filterer,
            iter.Func(m.Iterer().Next),
        ),
    )
}

func (m Mapper[K, V, ToK, ToV]) View() View[ToK, ToV] {
    return view[ToK, ToV]{
        get: m.get,
        set: m.set,
        delete: m.delete,
        iter: m.iter,
    }
}

func (w view[K, V]) Get(k K) (V, bool) {
    return w.get(k)
}

func (w view[K, V]) Set(k K, v V) {
    w.set(k, v)
}

func (w view[K, V]) Delete(k K) {
    w.delete(k)
}

func (w view[K, V]) Iter() Iter[Pair[K, V]] {
    return w.iter()
}
