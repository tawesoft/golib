package legacy

import (
    "fmt"
    "io"
    "reflect"
    "unicode/utf8"

    "github.com/tawesoft/golib/v2/must"
)

func zero[T any]() T {
    var t T
    return t
}

// Cast converts an element of type X to an interface of type Y.
//
// For example, slices.Map(ks.Cast[*os.File, io.Reader], listOfFilePointers)
//
// This appears alone because reflect is a heavy import
func Cast[X any, Y any](x X) Y {
    ref := reflect.ValueOf(&x).Elem()
    return ref.Interface().(Y)
}


// Initsh casts an int to comparable. If the comparable is not an integer type,
// this will panic.
func intish[T any](i int) T {
    var t T
    ref := reflect.ValueOf(&t).Elem()
    ref.SetInt(int64(i))
    return t
}

// Range calls some function f(k, v) => err over any [Rangeable] of (K, V)s. If
// the return value of f is not nil, the iteration stops immediately, and
// returns (k, err) for the given k. Otherwise, returns (zero, nil).
//
// Caution: invalid key types will panic at runtime. The key type must be int
// for any type other than a map. See [Rangeable] for details. In a channel,
// the key is always zero.
func Range[K comparable, V any, R Rangeable[K, V]](
    f func(K, V) error,
    r R,
) (K, error) {
    switch ref := reflect.ValueOf(r); ref.Kind() {
        case reflect.Array:
            for i := 0; i < ref.Len(); i++ {
                k := intish[K](i)
                v := ref.Index(i).Interface().(V)
                if err := f(k, v); err != nil { return k, err }
            }
        case reflect.Chan:
            for {
                x, ok := ref.Recv()
                if !ok { break }
                k := intish[K](0)
                v := x.Interface().(V)
                if err := f(k, v); err != nil { return k, err }
            }
        case reflect.Map:
            iter := ref.MapRange()
            for iter.Next() {
                k, v := iter.Key().Interface().(K), iter.Value().Interface().(V)
                if err := f(k, v); err != nil { return k, err }
            }
        case reflect.Slice:
            for i := 0; i < ref.Len(); i++ {
                k := intish[K](i)
                v := ref.Index(i).Interface().(V)
                if err := f(k, v); err != nil { return k, err }
            }
        case reflect.String:
            runes := 0
            for i := 0; i < ref.Len(); i++ {
                k := intish[K](runes)
                runes++
                v := ref.Index(i).Interface().(byte)

                if utf8.FullRune([]byte{v}) {
                    if err := f(k, intish[V](int(v))); err != nil { return k, err }
                    continue
                }

                if utf8.RuneStart(v) {
                    buf := [4]byte{v, 0, 0, 0}
                    for j := 0; j < 3; j++ {
                        //if i + j + 1 >= ref.Len() { break }
                        v = ref.Index(i+j+1).Interface().(byte)
                        buf[1+j] = v
                        if utf8.FullRune(buf[0:2+j]) {
                            n, _ := utf8.DecodeRune(buf[0:2+j])
                            i += j + 1
                            if err := f(k, intish[V](int(n))); err != nil { return k, err }
                            break
                        }
                    }
                } else {
                    return zero[K](), fmt.Errorf("Unicode error at byte %d %x", i, v)
                }
            }
    }
    return zero[K](), nil
}

// Rangeable defines any type of value x where it is possible to range over
// using "for k, v := range x" or "v := range x" (in the case of a channel,
// only "v := range x" is permitted). For every Rangeable other than a map,
// K must always be int.
type Rangeable[K comparable, V any] interface {
    ~string | ~map[K]V | ~[]V | chan V
}

// WithCloser ...
//
// Deprecated: don't use
func WithCloser[T io.Closer](opener func() (T, error), do func(v T) error) error {
    var zero T

    f, err := opener()
    if err != nil { return fmt.Errorf("WithCloser[%T] open error: %w", zero, err) }

    doer := must.CatchFunc(func() error { return do(f) })
    err, panicErr := doer()
    if err != nil {
        err = fmt.Errorf("WithCloser[%T] error: %w", zero, err)
    } else if panicErr != nil {
        err = fmt.Errorf("WithCloser[%T] error: panic: %w", zero, panicErr)
    }

    errClose := f.Close()
    if errClose != nil {
        err = fmt.Errorf("WithCloser[%T] close error: %v; %w", zero, errClose, err)
    }

    return err
}
