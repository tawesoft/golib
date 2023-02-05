// Package color implements parsing, formatting, and manipulating CSS colours
// based on the [CSS Color Module Level 4] (W3C Candidate Recommendation
// Draft), 1 November 2022.
//
// NOTE - INCOMPLETE! DO NOT USE YET.
//
// Note that named colors and system colors are not implemented.
//
// [CSS Color Module Level 4]: https://www.w3.org/TR/css-color-4/
//
// Disclaimer: although this software runs against a thorough and diverse set
// of test cases, no claims are made of this software's performance or
// conformance against the W3C Specification itself because it has not yet
// been tested against the relevant W3C test suite.
//
// This software includes material derived from CSS Color Module Level 4,
// W3C Candidate Recommendation Draft, 1 November 2022. Copyright © 2021 W3C®
// (MIT, ERCIM, Keio, Beihang). See LICENSE-PARTS.txt and TRADEMARKS.md.
//
// ## "Specified", "Computed", "Actual" and "Used" values in CSS
//
// CSS has the concept of "specified", "computed", "actual" and "used"
// values. For example, "rgb(512, 128, 0)" may be the specified value. CSS
// says it is not invalid, but is clamped to "rgb(255, 128, 0)" at compute
// time. Similarly, color(display-p3 1.0844 0.43 0.1) is "valid", remains so
// at compute time, but is out of gamut and so the "actual" value is the color
// given by gamut mapping the color for the display. This distinction is up to
// the caller to manage, but function documentation will give notes on how to
// achieve this.
//
// ##  “Missing” Color Components and the none Keyword
//
// Sometimes, a color component can be "missing". In CSS this can be manually
// specified by the keyword "none". In this Go implementation, the distinction
// between missing and present values is implemented with the [maybe.M] type.
//
// Color components can also be "powerless". For example, in hsl(), the hue
// component is powerless when the saturation component is 0% - this is a
// grayscale color, and hue has no effect on it, no matter its value. A
// powerless component automatically produced by color space conversion will
// be set to missing.
//
// See section 4.4 of the specification for more details.
//
// TODO only actually need float16 precision...
package color

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/tawesoft/golib/v2/fun/maybe"
)

const (
    typeHex = "#hexidecimal"
    typeRGB = "rgb()"
    typeHSL = "hsl()"
    typeHWB = "hwb()"
    typeLab = "lab()"
    typeLch = "lch()"
    typeColor = "color()"
)

type Color struct {
    _type string
    space Space
    components [3]maybe.M[float64]
    alpha maybe.M[float64]
}

// String returns a serialised representation of the color.
//
// Note that serialised colours are normally derived from the computed, not
// specified, colour. Therefore, this function is only compliant with the spec
// when serialising a computed colour. The [Color.Norm] method can help with
// calculating a computed colour.
//
// Even though the spec does not define behaviour for serialising non-computed
// values, this function still returns a useful serialisation for such values
// which is generally the same as the serialisation of a computed value without
// clamping or rgb conversion applied.
//
// Note that serialisation generally uses a fallback legacy format as far as
// possible. For example, the color returned by parsing "rgb(128 64 32 / 50%)"
// will always be serialised into the legacy format "rgba(128, 64, 32, 0.5)".
func (c Color) String() string {
    const decimalSigFigs = 4
    var sb strings.Builder
    var f string // e.g. "rgb" or "rgba"
    var cC [4]float64
    var cPc [4]bool // '%' percent

    alpha := c.alpha.Or(1.0)
    omitAlpha := roughlyEqual(alpha, 1.0)
    cC[3] = alpha

    switch c._type {
        case typeHex:
            r := int(0.5 + (c.components[0].Or(0.0) * 255.0))
            g := int(0.5 + (c.components[1].Or(0.0) * 255.0))
            b := int(0.5 + (c.components[2].Or(0.0) * 255.0))
            a := int(0.5 + (alpha * 255.0))
            if omitAlpha {
                return fmt.Sprintf("#%02x%02x%02x", r, g, b)
            } else {
                return fmt.Sprintf("#%02x%02x%02x%02x", r, g, b, a)
            }
        case typeRGB:
            cC[0] = c.components[0].Or(0.0) * 255.0
            cC[1] = c.components[1].Or(0.0) * 255.0
            cC[2] = c.components[2].Or(0.0) * 255.0
            if omitAlpha {
                f = "rgb"
            } else {
                f = "rgba"
                cC[3] = alpha
            }
        case typeHSL:
            cC[0] = c.components[0].Or(0.0) * 360.0
            cC[1] = c.components[1].Or(0.0) * 100.0; cPc[1] = true
            cC[2] = c.components[2].Or(0.0) * 100.0; cPc[2] = true
            if omitAlpha {
                f = "hsl"
            } else {
                f = "hsla"
            }
        default:
            return "color() /* error */"
    }

    sb.WriteString(f)
    sb.WriteByte('(')

    n := 4
    if omitAlpha { n = 3 }
    for i := 0; i < n; i++ {
        d, r := modf(cC[i])
        s := strconv.FormatFloat(r, 'g', decimalSigFigs, 64)
        if len(s) < 3 { // just '0'
            sb.WriteString(fmt.Sprintf("%d", d))
        } else {
            sb.WriteString(fmt.Sprintf("%d.%s", d, s[2:]))
        }
        if cPc[i] { sb.WriteByte('%') }

        if (i + 1) < n {
            sb.WriteByte(',')
            sb.WriteByte(' ')
        }
    }

    sb.WriteByte(')')
    return sb.String()
}

// Hexadecimal returns a color as if specified in a CSS RGB hexadecimal
// notation (e.g. #FFCC77).
func Hexadecimal(red, green, blue, alpha uint8) Color {
    f := maybe.Lift(func(x uint8) float64 {
        return float64(x) / 255.0
    })
    return Color{
        _type: typeHex,
        space: SpaceSRGB,
        components: [3]maybe.M[float64]{
            f(red),
            f(green),
            f(blue),
        },
        alpha: f(alpha),
    }
}

// RGB returns a color as if specified by the CSS rgb() function. However,
// each argument here is specified in the normalized range [0,1] (but may lie
// outside this range until computed). The computed value needs to clamp the
// input to the allowed range with the [Color.Norm] method.
func RGB(
    red maybe.M[float64],
    green maybe.M[float64],
    blue maybe.M[float64],
    alpha maybe.M[float64],
) Color {
    return Color{
        _type: typeRGB,
        space: SpaceSRGB,
        components: [3]maybe.M[float64]{red, green, blue},
        alpha: alpha,
    }
}

// HSL returns a color as if specified by the CSS hsl() function. However,
// each argument here is specified in the normalized range [0,1] (but may lie
// outside this range until computed). The computed value needs to clamp the
// input to the allowed range with the [Color.Norm] method, which also converts
// a HSL color to RGB.
func HSL(
    hue maybe.M[float64],
    saturation maybe.M[float64],
    lightness maybe.M[float64],
    alpha maybe.M[float64],
) Color {
    return Color{
        _type: typeRGB,
        space: SpaceSRGB,
        components: [3]maybe.M[float64]{hue, saturation, lightness},
        alpha: alpha,
    }
}

// Norm performs some steps to "normalise" a Color as part of turning a
// "specified" value into a "computed" value.
//
// For colours specified as hexadecimal, rgb(), rgba(), hsl(), hsla(), hwb(),
// or named colours, this involves clamping to the normalised range for each
// colour component, and converting to the rgb representation.
//
// For colours specified as lab() or lch(),  oklab(), or oklch(), this involves
// clamping to the normalised range for each colour component where a bound
// exists.
//
// For colours specified by the color() function, values are not clamped (this
// means that, although valid, they may still be out of gamut).
//
// For colours specified by the color() function using the "xyz" color space
// (which is an alias of the xyz-d65 color space), the computed and used value
// is in the xyz-d65 color space.
//
// All alpha values are clamped to the range [0, 1].
func (c Color) Norm() Color {

    // clamp all except color() function defined
    switch c._type {
        case typeHex: fallthrough
        case typeRGB: fallthrough
        case typeHSL:
            clampComponents(clamp_0_1, c.loadPtrs(0, 4))
    }

    // convert hexadecimal, hsl, hsla, hwb, named colors to rgb
    switch c._type {
        case typeHex: fallthrough
        case typeRGB:
            c._type = typeRGB
    }

    return c
}

// Map performs [gamut mapping] of a colour using the CSS gamut mapping
// algorithm. It returns a gamut mapped colour still in its original colour
// space, but representable in the destination color space.
//
// This is necessary when a colour (in any color space) is to be represented in
// a destination colour space where it would be could not be physically
// produced on a display device (e.g. a screen).
//
// The destination colour space may be the same colour space as the input
// colour's space. For example, color(srgb 200% 0% 0%) is a valid colour, and
// could be used as a gradient stop, but is out of gamut for the sRGB color
// space and needs gamut mapping.
//
// Some destination color spaces are not intended for display, and therefore
// have no gamut limits, and therefore perform no gamut mapping. These are the
// XYZ, Lab, LCH, Oklab, Oklch spaces.
//
// Note that for other purposes, such as print, different gamut mapping
// functions (not specified by CSS) may be more appropriate.
//
// [gamut mapping]: https://www.w3.org/TR/css-color-4/#gamut-mapping
func Map(dest Space, c Color) Color {
    // return gamutMap(x, dest)
    return Color{} // TODO
}

// Equal returns true if colour a is the same as colour b, irrespective of the
// type of each color or what color space they are defined in. That is, the two
// colours are perceptually equal to the standard observer (assuming correctly
// calibrated displays or printers!)
func Equal(a Color, b Color) bool {
    // TODO special case when a.Color == b.Color?
    // TODO convert both to XYZ and compare
    return false // TODO
}
