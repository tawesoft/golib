// Package view provides customisable abstractions over collections. Changes to
// an underlying collection are reflected in its views, and vice-versa.
package view

import (
    "github.com/tawesoft/golib/v2/iter"
    "github.com/tawesoft/golib/v2/ks"
)
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

// View is an abstraction over a collection. In the case of a slice, the key
// is an int.
type View[K comparable, V any] interface {
    Get(K) (V, bool)
    Set(K, V)
    Delete(K)
    Iter() iter.It[iter.Pair[K, V]]
}

// view is a concrete implementation of View.
type view[K comparable, V any] struct {
    get func(K) (V, bool)
    set func(K, V)
    delete func(K)
    iter func() iter.It[iter.Pair[K, V]]
}

// mapper defines how a collection is mapped to and from a view.
//
// Call the [mapper.View] method to construct a view from a mapper.
type mapper[K comparable, V any, ToK comparable, ToV any] struct {
    // Filterer defines a function that "hides" the given value from the view
    // when accessed through Get or Iter. It is implemented in terms of the
    // types of the underlying collection, and should not implement the keyer
    // or valuer logic.
    Filterer func(iter.Pair[K, V]) bool

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
    Iterer func() iter.It[iter.Pair[K, V]]
}

// FromMap returns a View from a Go map collection type.
//
// Filterer defines a function that "hides" the given value from the view
// when accessed through Get or Iter. It is implemented in terms of the
// types of the underlying collection, and should not implement the keyer
// or valuer logic.
//
// Keyer defines a mapping to and from the keys in the underlying
// collection to the view.
//
// Valuer defines a mapping to and from the values in the underlying
// collection to the view.
func FromMap[K comparable, V any, ToK comparable, ToV any](
    m map[K]V,
    filterer func(K, V) bool,
    keyer Keyer[K, ToK],
    valuer Valuer[V, ToV],
) View[ToK, ToV] {
    return mapper[K, V, ToK, ToV]{
        Filterer: func(x iter.Pair[K, V]) bool {
            if filterer == nil { return true }
            return filterer(x.Key, x.Value)
        },
        Getter: func(k K) (V, bool) {
            v, ok := m[k]
            return v, ok
        },
        Setter: func(k K, v V) {
            m[k] = v
        },
        Deleter: nil,
        Keyer: keyer,
        Valuer: valuer,
        Iterer: func() iter.It[iter.Pair[K, V]] {
            return iter.FromMap(m)
        },
    }.View()
}

// FromSlice returns a View from a Go list collection type.
//
// Filterer defines a function that "hides" the given value from the view
// when accessed through Get or Iter. It is implemented in terms of the
// types of the underlying collection, and should not implement the keyer
// or valuer logic.
//
// Valuer defines a mapping to and from the values in the underlying
// collection to the view.
func FromSlice[V any, ToV any](
    s []V,
    filterer func(int, V) bool,
    valuer Valuer[V, ToV],
) View[int, ToV] {
    return mapper[int, V, int, ToV]{
        Filterer: func(x iter.Pair[int, V]) bool {
            if filterer == nil { return true }
            return filterer(x.Key, x.Value)
        },
        Getter: func(k int) (V, bool) {
            ok := (k < len(s)) && (k >= 0)
            if !ok { return ks.Zero[V](), false }
            return s[k], true
        },
        Setter: func(k int, v V) {
            s[k] = v
        },
        Deleter: nil,
        Keyer: Keyer[int, int]{To: ks.Identity[int], From: ks.Identity[int]},
        Valuer: valuer,
        Iterer: func() iter.It[iter.Pair[int, V]] {
            return iter.Enumerate(iter.FromSlice(s))
        },
    }.View()
}

func (m mapper[K, V, ToK, ToV]) get(key ToK) (ToV, bool) {
    k := m.Keyer.From(key)
    v, ok := m.Getter(k)
    if ok && m.Filterer(iter.Pair[K, V]{k, v}) {
        return m.Valuer.To(v), true
    }
    return ks.Zero[ToV](), false
}

func (m mapper[K, V, ToK, ToV]) set(k ToK, v ToV) {
    m.Setter(m.Keyer.From(k), m.Valuer.From(v))
}

func (m mapper[K, V, ToK, ToV]) delete(k ToK) {
    m.Deleter(m.Keyer.From(k))
}

func (m mapper[K, V, ToK, ToV]) iter() iter.It[iter.Pair[ToK, ToV]] {
    mapper := func(i iter.Pair[K, V]) iter.Pair[ToK, ToV] {
        return iter.Pair[ToK, ToV]{m.Keyer.To(i.Key), m.Valuer.To(i.Value)}
    }

    return iter.Map[iter.Pair[K, V], iter.Pair[ToK, ToV]](mapper,
        iter.Filter[iter.Pair[K, V]](m.Filterer,
            m.Iterer(),
        ),
    )
}

// View returns a new view from an instantiated [mapper].
func (m mapper[K, V, ToK, ToV]) View() View[ToK, ToV] {
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

func (w view[K, V]) Iter() iter.It[iter.Pair[K, V]] {
    return w.iter()
}
