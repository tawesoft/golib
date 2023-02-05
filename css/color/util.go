package color

import (
    "math"

    "github.com/tawesoft/golib/v2/fun/maybe"
)

func modf(v float64) (int, float64) {
    d, f := math.Modf(v)
    inc := math.Copysign(0.5, f)
    i := int(d + inc) // int, truncated towards zero
    return i, f
}

func roughlyEqual(a float64, b float64) bool {
    // based on python "isClose"
    // https://peps.python.org/pep-0485/#proposed-implementation
    // TODO make this a "kitchen sink" function in the ks package?
    reltol := 1e-12
    abstol := 0.0
    return math.Abs(a-b) <= math.Max(reltol * math.Max(math.Abs(a), math.Abs(b)), abstol)
}

// componentPtrs make it easier to apply a function to each component in a
// Color.
type componentPtrs struct {
    start, end int
    components [4]*maybe.M[float64]
}

func (c *Color) loadPtrs(start int, end int) componentPtrs {
    var components [4]*maybe.M[float64]
    for i := start; i < end; i++ {
        if i < 3 {
            components[i] = &c.components[i]
        } else {
            components[i] = &c.alpha
        }
    }
    return componentPtrs{
        start: start,
        end: end,
        components: components,
    }
}

var (
    clamp_0_1 = clampFunc(0.0, 1.0)
)

func clampComponents(clampFunc clampFuncT, ptrs componentPtrs) {
    clampFuncM := maybe.Map(clampFunc)
    for _, x := range ptrs.components {
        *x = clampFuncM(*x)
    }
}

type clampFuncT func(x float64) float64
func clampFunc(
    min float64,
    max float64,
) clampFuncT {
    return func(x float64) float64 {
        return math.Min(math.Max(x, min), max)
    }
}
