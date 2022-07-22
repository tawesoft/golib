// Package human is an elegant, general-purpose, extensible, modular,
// locale-aware way to format and parse numbers and quantities - like distances,
// bytes, and time - in a human-readable way ideal for config files, user-input,
// website scraping, and as a building-block for fully translated ergonomic
// user interfaces.
//
// Package human integrates very nicely with the localisation features at
// golang.org/x/text.
//
// For example:
//
//    dutch := human.New(language.Dutch)
//    v, err := dutch.ParseDecimalInt("1.234.567")
//    fmt.Println(v) // prints 1234567
//
//    eng := human.New(language.English)
//    fmt.Println(h.FormatDecimal("%d", 123456789)) // prints 123,456,789
//    fmt.Println(h.FormatDistance(1500)) // prints 1.5km
//
package human

import (
    "strconv"
    "strings"

    "golang.org/x/text/language"
    "golang.org/x/text/message"
    "golang.org/x/text/number"
)

// NumberParser specifies the interface for parsing numbers in a specific
// format (e.g. decimal) and a specific locale (e.g. British English).
//
// The AcceptInt and AcceptFloat functions parse as many bytes as possible from
// the start of the input string. Their first returned value is the parsed
// number. The second return value is the number of bytes (not runes)
// successfully parsed. The third return value is non-nil only if the parser
// cannot parse using the configured language, or strconv.ErrRange. These
// methods do not parse negative numbers, "NaN" or (signed or unsigned)
// infinity.
//
// Only the AcceptSignedInt and AcceptSignedFloat function variants parse
// negative numbers.
type NumberParser interface {
    AcceptInt(string) (int64, int, error)
    AcceptFloat(string) (float64, int, error)

    AcceptSignedInt(string) (int64, int, error)
    AcceptSignedFloat(string) (float64, int, error)
}

func ParseInt(p NumberParser, s string) (int64, error) {
    s = strings.TrimSpace(s)
    value, length, err := p.AcceptSignedInt(s)
    if (err == nil) && (length != len(s)) {
        return 0, strconv.ErrSyntax
    }
    return value, err
}

func ParseFloat(p NumberParser, s string) (float64, error) {
    s = strings.TrimSpace(s)
    value, length, err := p.AcceptSignedFloat(s)
    if (err == nil) && (length != len(s)) {
        return 0, strconv.ErrSyntax
    }
    return value, err
}

// Humanizer is a parser and formatter for numbers and quantities in a
// specific language.
//
// Initialise with [Humanizer.New], or for greater control initialise the
// struct yourself.
type Humanizer struct {
    // Locale specifies language or locale e.g. British English, Welsh, etc.
    Locale language.Tag

    // Printer is used to format numbers. See [message.NewPrinter].
    Printer *message.Printer

    // DecimalNumberFormatter is used to format decimal numbers. May be
    // customised with [number.Option] arguments to [number.NewFormat].
    DecimalNumberFormatter number.FormatFunc

    // DecimalNumberParser is used to parse decimal numbers.
    DecimalNumberParser NumberParser
}

// New initialises a humanizer for a given language.
//
// For the tag argument, try a named builtin, like language.English, load
// one with language.Parse("cy"), or see [language.Matcher].
func New(tag language.Tag) *Humanizer {

    // TODO add Options for DecimalParser

    mp := message.NewPrinter(tag)

    return &Humanizer{
        Locale: tag,
        Printer: mp,
        DecimalNumberFormatter: number.NewFormat(number.Decimal),
        DecimalNumberParser: NewDecimalParser(tag),
    }
}

// AcceptDecimalInt parses part of a string representing a decimal number.
// See [NumberParser]. For example, "1,000,000" => 1000000.
func (h Humanizer) AcceptDecimalInt(x string) (int64, int, error) {
    return h.DecimalNumberParser.AcceptInt(x)
}

// AcceptDecimalFloat parses part of a string representing a decimal number.
// See [NumberParser]. For example, "1,299.99" => 1299.99.
func (h Humanizer) AcceptDecimalFloat(x string) (float64, int, error) {
    return h.DecimalNumberParser.AcceptFloat(x)
}

// FormatDecimal formats a number as a string representing a decimal number in
// the Humanizer's locale. The format string is a printf-style format string (like
// "%d", "%.2f", etc). For example, 1000000 => "1,000,000".
func (h Humanizer) FormatDecimal(format string, n any, opts ... number.Option) string {
    return h.Printer.Sprintf(format, h.DecimalNumberFormatter(n, opts...))
}

// ParseDecimalInt parses a string representing a decimal number in the
// Humanizer's locale. For example, "1,000,000" => 1000000.
func (h Humanizer) ParseDecimalInt(x string) (int64, error) {
    return ParseInt(h.DecimalNumberParser, x)
}

// ParseDecimalFloat parses a string representing a decimal number in the
// Humanizer's locale. For example, "1,299.99" => 1299.99.
func (h Humanizer) ParseDecimalFloat(x string) (float64, error) {
    return ParseFloat(h.DecimalNumberParser, x)
}


/*
func (h Humanizer) Accept(str string, unit Unit, factors Factors) {
}

    // Format is a general purpose locale-aware way to format any quantity
    // with a defined set of factors. The unit argument is the base unit
    // e.g. s for seconds, m for meters, B for bytes.
    Format(value float64, unit Unit, factors Factors) String

    FormatNumber(number float64) String           // e.g. 12 k
    FormatDistance(meters float64) String         // e.g. 10 µm, 10 km
    FormatDuration(duration time.Duration) string // e.g. 1 h 50 min
    FormatSeconds(seconds float64) string         // e.g. 1 h 50 min
    FormatBytesJEDEC(bytes int64) string          // e.g. 12 KB, 5 MB
    FormatBytesIEC(bytes int64) string            // e.g. 12 kB, 5 MB
    FormatBytesSI(bytes int64) string             // e.g. 12 KiB, 5 MiB

    // Accept is a general purpose locale-aware way to parse any quantity
    // with a defined set of factors from the start of the string str. The
    // provided unit is optional and is accepted if it appears in str.
    //
    // Accept returns the value, the number of bytes successfully parsed (which
    // may be zero), or an error.
    Accept(str string, unit Unit, factors Factors) (float64, int, error)

    // Parse is a general purpose locale-aware way to parse any quantity
    // with a defined set of factors.  The provided unit is optional and is
    // accepted if it appears in str.
    Parse(str string, unit Unit, factors Factors) (float64, error)
 */
