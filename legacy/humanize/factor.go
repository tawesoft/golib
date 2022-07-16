package humanize

// Factors describes a way to format a quantity with units.
type Factors struct{
    // Factors is a list of Factor entries in ascending order of magnitude.
    Factors []Factor

    // Components controls how the formatting is broken up. If set to zero or
    // one, formatting has a single component e.g. "1.5 M". If set to two or
    // more, formatting is broken up into previous factors e.g. "1 h 50 min"
    // (with 2 components) or "1 h 50 min 25 s" (with 3 components).
    Components int
}

// Factor defines one entry in an ordered list of Factors.
type Factor struct{

    // Magnitude defines the absolute size of the factor e.g. 1000000 for the
    // SI unit prefix "M".
    Magnitude float64

    // Label describes the magnitude, usually as a unit prefix (like SI "M")
    // or as a replacement (like "min"), controlled by Mode.
    Unit Unit

    // Mode controls the formatting of this factor
    Mode FactorMode
}

// FactorMode controls the formatting of a factor.
type FactorMode int

const (
    // FactorModeIdentity indicates that the given factor label represents the
    // unit with no changes. For example, if you're talking about distance in
    // metres, FactorModeIdentity means that the current factor is measured in
    // metres and not millimetres, or kilometres.
    FactorModeIdentity   = FactorMode(0)

    // FactorModeUnitPrefix indicates the given factor label is a unit prefix
    // e.g. "Ki" is a byte prefix giving "KiB".
    FactorModeUnitPrefix = FactorMode(1)

    // FactorModeReplace indicates the given factor label replaces the current
    // unit e.g. the duration of time 100 s becomes 1 min 40 s, not 1 hs
    // (which would read as a "hectosecond"!).
    FactorModeReplace    = FactorMode(2)

    // FactorModeInputCompat indicates that the given factor label should only
    // be considered on input. This may be combined with any other FactorMode
    // by a bitwise OR operation. For example, when measuring distance,
    // you might accept based on context that "K" probably refers to the
    // "kilo"-prefix, not a Kelvin-metre!
    FactorModeInputCompat = FactorMode(4)
)

// Min returns the index of the first Factor greater or equal to n. If n is
// smaller than the first Factor, returns the first Factor instead. Ignores all
// factors that have mode FactorModeInputCompat.
func (f Factors) Min(n float64) int {
    if len(f.Factors) == 0 { panic("empty list of factors") }

    if n < f.Factors[0].Magnitude { return 0 }

    for i, factor := range f.Factors[1:] {
        if factor.Mode & FactorModeInputCompat == FactorModeInputCompat {
            continue // skip
        }

        if n < factor.Magnitude {
            return i
        }
    }

    return len(f.Factors) - 1
}
