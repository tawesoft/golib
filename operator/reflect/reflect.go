package reflect

import (
    "reflect"
)

// Cast converts between values and interfaces.
func Cast[X any, Y any](x X) Y {
    ref := reflect.ValueOf(&x).Elem()
    return ref.Interface().(Y)
}
