package color

// WhitePoint defines the colour white under a certain viewing condition, such
// as daylight, as chromaticity coordinates (x,y) of a perfectly reflecting
// diffuser.
//
// For the purposes of CSS, these are only the coordinates for a 2 degree
// standard observer.
//
// Wikipedia says "if the color of an object is recorded under one illuminant,
// then it is possible to estimate the color of that object under another
// illuminant, given only the white points of the two illuminants."
type WhitePoint struct {
    X, Y float64
}

var D50 = WhitePoint{0.345700, 0.358500}
var D65 = WhitePoint{0.312700, 0.329000}

// Space is a color space. It describes "an organization of colours with respect
// to an underlying colorimetric model, such that there is a clear,
// objectively-measurable meaning for any colour in that colour space. This
// also means that the same color can be expressed in multiple color spaces,
// or transformed from one color space to another, while still looking the
// same."
type Space string

func (s Space) Name() string {
    return string(s)
}

const (
    // SpaceSRGB is the colour space used by RGB, HSL, HWB.
    SpaceSRGB = Space("sRGB")

    // SpaceSRGBLinear predefined colour space is "the same as sRGB except that
    // the transfer function is linear-light (there is no gamma-encoding)."
    SpaceSRGBLinear = Space("sRGB-linear")

    // SpaceDisplayP3 has the same transfer curve as sRGB but a wider gamut.
    // "Modern displays, TVs, laptop screens and phone screens are able to
    // display all, or nearly all, of the display-p3 gamut."
    SpaceDisplayP3 = Space("display-p3")

    // SpaceA98RGB is inaccurate, but with a wider gamut than sRGB. It was
    // developed by Adobe and often used in Photoshop and with some
    // professional displays.
    SpaceA98RGB = Space("a98-RGB")

    // SpaceProPhotoRGB (also known as ROMM RGB) was developed by Kodak and
    // is often used in digital photography due to its very wide gamut.
    SpaceProPhotoRGB = Space("prophoto-RGB")

    // SpaceRec2020 (the ITU-R BT.2020-2 colour space) has a very wide
    // gamut and is used in ultra-high-definition television.
    SpaceRec2020 = Space("rec2020")

    // SpaceXYZ is a CIE XYZ colour space that covers the full gamut of
    // human vision. The Y component represents luminosity and the XZ plane
    // contains all possible chromaticities at that luminance.
    SpaceXYZ = SpaceXYZD65

    // SpaceXYZD50 is the CIE XYZ color space with a D50 white point.
    SpaceXYZD50 = Space("xyz-d50")

    // SpaceXYZD65 is the CIE XYZ color space with a D65 white point.
    SpaceXYZD65 = Space("xyz-d65")
)
