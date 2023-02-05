// Package view provides customisable abstractions over collections. Changes to
// an underlying collection are reflected in its views, and vice-versa.
package view

import (
    "github.com/tawesoft/golib/v2/iter"
    "github.com/tawesoft/golib/v2/operator"
)

// View is an abstraction over a collection. In the case of a slice, the key
// is an int representing the index.
type View[K comparable, V any] interface {
    Get(K) (V, bool)
    Set(K, V)
    Delete(K)
    Iter() iter.It[iter.Pair[K, V]]
}

// Viewer describes a View over a Go collection type that maps to a new
// collection of mapped keys and values.
//
// This is a generic, low-level definition for implementors of custom
// collection types. Most users will want to call the Bind method on the
// simpler [Map] or [Slice] constructor types instead.
type Viewer[K comparable, V any, ToK comparable, ToV any] struct {
    // Filterer defines a function that controls, given an element (key, value)
    // pair from the underlying collection, whether that element is mapped
    // to the new collection in Get or Iter methods on a View. Elements only
    // appear if the filter function returns true. If omitted, defaults to a
    // function that always returns true.
    Filterer func(K, V) bool

    // Getter defines a function that accesses a given value from the underlying
    // collection. It is implemented in terms of the types of the underlying
    // collection, and should not implement the filtering itself.
    Getter func(K) (V, bool)

    // Setter defines a function that sets a given value in the underlying
    // collection. It is implemented in terms of the types of the underlying
    // collection, and ignores filtering.
    Setter func(K, V)

    // Deleter defines a function that deletes a given value in the underlying
    // collection. It is implemented in terms of the types of the underlying
    // collection, and ignores filtering.
    Deleter func(K)

    // ToKey defines a function that maps K to ToK.
    ToKey func(K) ToK

    // FromKey defines the inverse of ToKey.
    FromKey func(ToK) K

    // ToValue defines a function that maps V to ToV.
    ToValue func(V) ToV

    // FromValue defines the inverse of ToValue.
    // Note that FromValue may be omitted if
    FromValue func(ToV) V

    // Iterer defines a function that returns a new iterator over the underlying
    // collection. It is implemented in terms of the types of the underlying
    // collection, and should not implement the filtering itself.
    Iterer func() iter.It[iter.Pair[K, V]]
}

// Map describes a way to create a View from a Go map collection type that
// maps to a collection of mapped keys and values.
//
// The types K and V define element types in an original collection.
// The types ToK and ToV define element types in a new collection formed by
// mapping K to ToK and mapping v to ToV.
type Map[K comparable, V any, ToK comparable, ToV any] struct {
    // Filterer defines a function that controls, given an element (key, value)
    // pair from the underlying collection, whether that element is mapped
    // to the new collection in Get or Iter methods on a View. Elements only
    // appear if the filter function returns true. If omitted, defaults to a
    // function that always returns true.
    Filterer func(K, V) bool

    // ToKey defines a function that maps K to ToK.
    ToKey func(K) ToK

    // FromKey defines the inverse of ToKey.
    FromKey func(ToK) K

    // ToValue defines a function that maps V to ToV.
    ToValue func(V) ToV

    // FromValue defines the inverse of ToValue.
    // Note that FromValue may be omitted if the [View.Set] method is never
    // called.
    FromValue func(ToV) V
}

// Bind returns a new View over a Go map collection.
func (m Map[K, V, ToK, ToV]) Bind(c map[K]V) View[ToK, ToV] {
    return Viewer[K, V, ToK, ToV]{
        Filterer: func(k K, v V) bool {
            if m.Filterer == nil { return true }
            return m.Filterer(k, v)
        },
        Getter: func(k K) (V, bool) {
            v, ok := c[k]
            return v, ok
        },
        Setter: func(k K, v V) {
            c[k] = v
        },
        Deleter: func(k K) {
            delete(c, k)
        },
        Iterer: func() iter.It[iter.Pair[K, V]] {
            return iter.FromMap(c)
        },
        ToKey:     m.ToKey,
        FromKey:   m.FromKey,
        ToValue:   m.ToValue,
        FromValue: m.FromValue,
    }
}

// Slice describes a way to create a View from a Go slice collection type that
// maps to a collection of mapped values.
//
// The type V defines element types in an original collection.
// The type ToV defines element types in a new collection formed by
// mapping v to ToV.
type Slice[V any, ToV any] struct {
    // Filterer defines a function that controls, given a value from the
    // underlying collection, whether that element is mapped to the new
    // collection in Get or Iter methods on a View. Values only appear if the
    // filter function returns true. If omitted, defaults to a function that
    // always returns true.
    Filterer func(V) bool

    // ToValue defines a function that maps V to ToV.
    ToValue func(V) ToV

    // FromValue defines the inverse of ToValue.
    // Note that FromValue may be omitted if the [View.Set] method is never
    // called.
    FromValue func(ToV) V
}

// Bind returns a new View over a Go slice collection.
func (s Slice[V, ToV]) Bind(c []V) View[int, ToV] {
    return Viewer[int, V, int, ToV]{
        Filterer: func(_ int, v V) bool {
            if s.Filterer == nil { return true }
            return  s.Filterer(v)
        },
        Getter: func(k int) (V, bool) {
            if (k < 0) || (k >= len(c)) {
                return operator.Zero[V](), false
            }
            return c[k], true
        },
        Setter: func(k int, v V) {
            if k == len(c) {
                c = append(c, v)
            } else {
                c[k] = v
            }
        },
        Deleter: func(k int) {
            if k == len(c) - 1 {
                c = c[0:k]
            } else {
                c = append(c[0:k], c[k+1:]...)
            }
        },
        Iterer: func() iter.It[iter.Pair[int, V]] {
            return iter.Enumerate(iter.FromSlice(c))
        },
        ToKey:     func (x int) int { return x },
        FromKey:   func (x int) int { return x },
        ToValue:   s.ToValue,
        FromValue: s.FromValue,
    }
}

func (m Viewer[K, V, ToK, ToV]) Get(key ToK) (ToV, bool) {
    k := m.FromKey(key)
    v, ok := m.Getter(k)
    if ok && m.Filterer(k, v) {
        return m.ToValue(v), true
    }
    return operator.Zero[ToV](), false
}

func (m Viewer[K, V, ToK, ToV]) Set(k ToK, v ToV) {
    m.Setter(m.FromKey(k), m.FromValue(v))
}

func (m Viewer[K, V, ToK, ToV]) Delete(k ToK) {
    m.Deleter(m.FromKey(k))
}

func (m Viewer[K, V, ToK, ToV]) pairFilterer(pair iter.Pair[K, V]) bool {
    return m.Filterer(pair.Key, pair.Value)
}

func (m Viewer[K, V, ToK, ToV]) Iter() iter.It[iter.Pair[ToK, ToV]] {
    mapper := func(i iter.Pair[K, V]) iter.Pair[ToK, ToV] {
        return iter.Pair[ToK, ToV]{m.ToKey(i.Key), m.ToValue(i.Value)}
    }

    return iter.Map[iter.Pair[K, V], iter.Pair[ToK, ToV]](mapper,
        iter.Filter[iter.Pair[K, V]](m.pairFilterer,
            m.Iterer(),
        ),
    )
}
