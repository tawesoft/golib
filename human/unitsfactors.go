package human

import (
    "math"
)

// String holds an Utf8-encoded and an Ascii-compatible encoding of a string.
type String struct {
    // Utf8 is the native Utf8-encoded Unicode representation
    Utf8 string

    // Ascii is an alternative version accepted for non-Unicode inputs (such
    // as when a user does not know how to enter µ on their keyboard) or for
    // non-Unicode output (such as legacy display systems).
    Ascii  string
}

// Unit describes some quantity. For example, "m" is a Unit for length in
// metres, "bps" is a Unit for download speed, and "k" is a Unit for the
// SI unit prefix "kilo-".
type Unit String

// Cat concatenates two units (u + v) and returns the result. For example,
// concatenating the Unit "k" for the SI unit prefix "kilo-" with the Unit "m"
// for length in metres gives Unit "km" for kilometres.
func (u Unit) Cat(v Unit) Unit {
    return Unit{
        u.Utf8 + v.Utf8,
        u.Ascii + v.Ascii,
    }
}

// Part describes some component of a formatting result e.g. the time
// representing one hour and twenty minutes is made up of two parts: (1, hour)
// and (20, minute).
type Part struct {
    Magnitude float64
    Unit      Unit
}

func partEqual(a Part, b Part, epsilon float64) bool {
    if a.Unit != b.Unit { return false }
    return math.Abs(a.Magnitude - b.Magnitude) < epsilon
}

func partsEqual(a []Part, b []Part, epsilon float64) bool {
    if len(a) != len(b) { return false }

    for i := 0; i < len(a); i++ {
        if !partEqual(a[i], b[i], epsilon) { return false }
    }

    return true
}

// Factors is a list of rules that specifies a way to format a quantity with
// units. For example, the value 195 with [Unit] seconds can be factored into
// the parts 3 (with Unit minutes) and 15 (with Unit seconds).
//
// See [CommonFactors] for ready-built implementations.
type Factors struct{
    // Factors is a list of Factor entries in ascending order of magnitude.
    Factors []Factor

    // Components controls how the formatting is broken up. If set to zero or
    // one, formatting has a single component e.g. "1.5 M". If set to two or
    // more, formatting is broken up into previous factors e.g. "1 h 50 min"
    // (with 2 components) or "1 h 50 min 25 s" (with 3 components).
    Components int
}

// Factor defines one entry in an ordered list of [Factors].
type Factor struct{

    // Magnitude defines the absolute size of the factor e.g. 1000000 for the
    // SI unit prefix "M", 0.001 for the SI unit prefix "m".
    Magnitude float64

    // Label describes the magnitude, usually as a unit prefix (like SI "M")
    // or as a replacement (like "min"), controlled by Mode.
    Unit Unit

    // Mode controls how this factor rule is applied
    Mode FactorMode
}

// FactorMode controls how a factor rule is applied.
type FactorMode int

const (
    // FactorModeIdentity indicates that the given factor label represents the
    // unit with no changes. For example, if you're talking about distance in
    // metres, FactorModeIdentity indicates that the current factor is measured
    // in metres (the matching factor Magnitude must be equal 1).
    FactorModeIdentity   = FactorMode(0)

    // FactorModeUnitPrefix indicates the given factor label is a unit prefix
    // e.g. "Ki" is a prefix that can be applied to Bytes giving a new "KiB"
    // (Kibibytes) unit.
    FactorModeUnitPrefix = FactorMode(1)

    // FactorModeReplace indicates the given factor label replaces the current
    // unit completely. For example, the duration of time given as 90 seconds
    // becomes 1.30 min. This is in contrast to FactorModeUnitPrefix, because
    // the minute unit replaces the second unit instead of adding a prefix:
    // there is no unit like "ks" ("Kiloseconds")!
    FactorModeReplace    = FactorMode(2)

    // FactorModeInputCompat indicates that the given factor label should only
    // be considered on input.  For example, when measuring file size,
    // you might accept based on context that "K" probably refers to the
    // "kilo"-prefix, not a Kelvin-byte!
    //
    // This mode may be combined with any other FactorMode by a bitwise OR
    // operation.
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

// CommonUnits are defined [Unit] values for commonly-used measurements.
var CommonUnits = struct{
    None           Unit
    Second         Unit
    Meter          Unit
    Byte           Unit
    Bit            Unit
    BitsPerSecond  Unit
}{
    None:          Unit{"", ""},
    Second:        Unit{"s", "s"},
    Meter:         Unit{"m", "m"},
    Byte:          Unit{"B", "B"},
    Bit:           Unit{"b", "b"},
    BitsPerSecond: Unit{"bps", "bps"},
}

// CommonFactors are defined [Factors] for commonly-used measurements.
var CommonFactors = struct{
    // Time is time units in seconds, minutes, hours, days and years as min, h,
    // d, and y. These are non-SI units but generally accepted in context.
    // For times smaller than a second (e.g. nanoseconds), use SI instead.
    // The expected unit is a second (Unit{"s", "s"} or CommonUnits.Second)
    Time Factors

    // Distance are SI units that stop at kilo (because nobody uses
    // megametres or gigametres!) but includes centi. The expected unit is the
    // SI unit for distance, the metre (Unit{"m", "m"} or CommonUnits.Meter)
    Distance Factors

    // IEC are the "ibi" unit prefixes for bytes e.g. Ki, Mi, Gi with a
    // factor of 1024.
    IEC Factors

    // JEDEC are the old unit prefixes for bytes: K, M, G (only) with a factor
    // of 1024.
    JEDEC Factors

    // SIBytes are the SI unit prefixes for bytes e.g. k, M, G with a
    // factor of 1000. Unlike the normal SI Factors, it is assumed based on
    // context that when a "K" is input this is intended to mean the "k" SI
    // unit prefix instead of Kelvin - I've never heard of a Kelvin-Byte!
    SIBytes Factors

    // SIUncommon are the SI unit prefixes including deci, deca, and hecto
    SIUncommon Factors

    // SI are the SI unit prefixes except centi, deci, deca, and hecto
    SI Factors
}{
    Time:       timeFactors,
    Distance:   distanceFactors,
    IEC:        iecFactors,
    JEDEC:      jdecFactors,
    SIBytes:    siByteFactors,
    SIUncommon: siUncommonFactors,
    SI:         siFactors,
}

var timeFactors = Factors{
    Factors: []Factor{
        {1,                                 Unit{"s", "s"},     FactorModeReplace},
        {60,                                Unit{"min", "min"}, FactorModeReplace},
        {60 * 60,                           Unit{"h", "h"},     FactorModeReplace},
        {24 * 60 * 60,                      Unit{"d", "d"},     FactorModeReplace},
        {365.2422 * 24 * 60 * 60,           Unit{"y", "y"},     FactorModeReplace},
    },
    Components: 2,
}

var distanceFactors = Factors{
    Factors: []Factor{
        {1E-9,                              Unit{"n", "n"},     FactorModeUnitPrefix}, // nano
        {1E-6,                              Unit{"μ", "u"},     FactorModeUnitPrefix}, // micro
        {1E-3,                              Unit{"m", "m"},     FactorModeUnitPrefix}, // milli
        {1E-2,                              Unit{"c", "c"},     FactorModeUnitPrefix}, // centi
        {1,                                 Unit{ "",  ""},     FactorModeIdentity},
        {1000,                              Unit{"k", "k"},     FactorModeUnitPrefix}, // kilo
    },
}

var iecFactors = Factors{
    Factors: []Factor{
        {1,                                 Unit{ "",  ""},     FactorModeUnitPrefix},
        {1024,                              Unit{"Ki", "Ki"},   FactorModeUnitPrefix},
        {1024 * 1024,                       Unit{"Mi", "Mi"},   FactorModeUnitPrefix},
        {1024 * 1024 * 1024,                Unit{"Gi", "Gi"},   FactorModeUnitPrefix},
        {1024 * 1024 * 1024 * 1024,         Unit{"Ti", "Ti"},   FactorModeUnitPrefix},
    },
}

var jdecFactors = Factors{
    Factors: []Factor{
        {1,                                 Unit{ "",  ""},     FactorModeIdentity},
        {1024,                              Unit{"K", "K"},     FactorModeUnitPrefix},
        {1024 * 1024,                       Unit{"M", "M"},     FactorModeUnitPrefix},
        {1024 * 1024 * 1024,                Unit{"G", "G"},     FactorModeUnitPrefix},
    },
}

var siByteFactors = Factors{
    Factors: []Factor{
        {1,                                 Unit{ "",  ""},     FactorModeIdentity},
        { 1E3,                              Unit{"k", "k"},     FactorModeUnitPrefix},
        { 1E3,                              Unit{"K", "K"},     FactorModeUnitPrefix | FactorModeInputCompat}, // Kelvin-Bytes(!)
        { 1E6,                              Unit{"M", "M"},     FactorModeUnitPrefix},
        { 1E9,                              Unit{"G", "G"},     FactorModeUnitPrefix},
        {1E12,                              Unit{"T", "T"},     FactorModeUnitPrefix},
    },
}

var siUncommonFactors = Factors{
    Factors: []Factor{
        {1E-9,                              Unit{"n", "n"},     FactorModeUnitPrefix}, // nano
        {1E-6,                              Unit{"μ", "u"},     FactorModeUnitPrefix}, // micro
        {1E-3,                              Unit{"m", "m"},     FactorModeUnitPrefix}, // milli
        {1E-2,                              Unit{"c", "c"},     FactorModeUnitPrefix}, // centi
        {1E-1,                              Unit{"d", "d"},     FactorModeUnitPrefix}, // deci
        {1,                                 Unit{ "",  ""},     FactorModeIdentity},
        { 1E1,                              Unit{"da", "da"},   FactorModeUnitPrefix}, // deca
        { 1E2,                              Unit{"h", "h"},     FactorModeUnitPrefix}, // hecto
        { 1E3,                              Unit{"k", "k"},     FactorModeUnitPrefix}, // kilo
        { 1E6,                              Unit{"M", "M"},     FactorModeUnitPrefix},
        { 1E9,                              Unit{"G", "G"},     FactorModeUnitPrefix},
        {1E12,                              Unit{"T", "T"},     FactorModeUnitPrefix},
    },
}

var siFactors = Factors{
    Factors: []Factor{
        {1E-9,                              Unit{"n", "n"},     FactorModeUnitPrefix}, // nano
        {1E-6,                              Unit{"μ", "u"},     FactorModeUnitPrefix}, // micro
        {1E-3,                              Unit{"m", "m"},     FactorModeUnitPrefix}, // milli
        {1,                                 Unit{ "",  ""},     FactorModeIdentity},
        { 1E3,                              Unit{"k", "k"},     FactorModeUnitPrefix}, // kilo
        { 1E6,                              Unit{"M", "M"},     FactorModeUnitPrefix},
        { 1E9,                              Unit{"G", "G"},     FactorModeUnitPrefix},
        {1E12,                              Unit{"T", "T"},     FactorModeUnitPrefix},
    },
}
