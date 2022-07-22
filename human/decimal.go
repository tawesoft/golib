package human

import (
    "fmt"
    "math"
    "strconv"

    "github.com/tawesoft/golib/v2/numbers"
    "golang.org/x/text/language"
    "golang.org/x/text/message"
    "golang.org/x/text/number"
)

// DecimalParser is a type of [NumberParser] that parses strings encoding a
// number in a positional base-10 number system (i.e. decimal).
//
// For example, a million is represented in British English as "1,000,000", but
// in many European countries as "1.000.000". Similarly, the constant pi is
// represented in British English as "3.141..." but in many European countries
// as "3,141". Numbers aren't always represented by the Western-Arabic
// numerals ("0123456789"), either: The number "1,234.56" (in British English)
// is written as "١٬٢٣٤٫٥٦" in Eastern-Arabic numerals.
type DecimalParser struct {
    // GroupSeparator is a set of all runes that may encode (in the locale) a
    // digit separator. For example, commas are used in British English to
    // separate groups of three digits. If blank, ignored. A single whitespace
    // digit is always treated as a group separator.
    GroupSeparators []rune

    // DecimalSeparator is a set of all runes that may encode (in the locale)
    // a separator between the integer and fractional part of a decimal number.
    // For example, in British English, this is a dot '.', but in many European
    // countries, it is a comma ','. Required.
    DecimalSeparators []rune

    // Digits are an ascending list of decimal digit runes (in the locale). For
    // example, in British English, "0123456789". Each entry is a set of all
    // runes that may encode that digit. For example, some Persian numerals
    // look identical to some Eastern Arabic numerals, but still have different
    // code points. A robust parser might accept either.
    DecimalDigits [10][]rune
}

// NewDecimalParser automatically configures and returns a [DecimalParser] with
// defaults for a given locale.
//
// Caution: this behaviour is undefined for locales that have non-base-10
// number systems. Unfortunately, this information, while it exists, isn't
// exposed from go's /x/text/internal.
func NewDecimalParser(tag language.Tag) NumberParser {
    var dp DecimalParser
    mp := message.NewPrinter(tag)

    if gs, ok := dp.guessGroupSeparator(mp); ok { // optional
        dp.GroupSeparators = []rune{gs}
    }

    if ds, ok := dp.guessDecimalSeparator(mp); ok {
        dp.DecimalSeparators = []rune{ds}
    } else {
        err := fmt.Errorf("unknown decimal separator")
        return invalidNumberParser("decimal", tag, err)
    }

    if ok := dp.guessDecimalDigits(mp, &dp.DecimalDigits); !ok {
        err := fmt.Errorf("unknown digits")
        return invalidNumberParser("decimal", tag, err)
    }

    return &dp
}

func (p *DecimalParser) AcceptInt(s string) (int64, int, error) {
    return acceptNumberPart[int64](
        s,
        p.GroupSeparators,
        p.guessRuneValue,
        numbers.Int64.CheckedMul,
        numbers.Int64.CheckedAdd,
    )
}

func (p *DecimalParser) AcceptSignedInt(s string) (int64, int, error) {
    if (len(s) > 0) && (s[0] == '-') {
        value, length, err := p.AcceptInt(s[1:])
        return -value, length + 1, err
    } else {
        return p.AcceptInt(s)
    }
}

func (p *DecimalParser) acceptFloatPart(s string) (float64, int, error) {
    return acceptNumberPart[float64](
        s,
        p.GroupSeparators,
        p.guessRuneValue,
        numbers.Float64.CheckedMul,
        numbers.Float64.CheckedAdd,
    )
}

func (p *DecimalParser) AcceptFloat(s string) (float64, int, error) {
    var (
        offset int
        leftValue float64
        rightValue float64
        ok bool
    )

    // leading decimal point e.g. ".0123"?
    if leadingPoint := acceptRuneFromSet(p.DecimalSeparators, s); leadingPoint != 0 {
        // yes, ignore for now
    } else {
        // no, so parse the left-hand side e.g. "123.456".
        value, length, err := p.acceptFloatPart(s)
        if err != nil { return 0, 0, err }
        leftValue = value
        s = s[length:]
        offset += length
    }

    // leading decimal point e.g. ".0123" after left (if any)?
    if leadingPoint := acceptRuneFromSet(p.DecimalSeparators, s); leadingPoint != 0 {
        s = s[leadingPoint:]

        offset, ok = numbers.Int.CheckedAdd(offset, leadingPoint)
        if !ok { return 0, 0, strconv.ErrRange }
    } else {
        // no, so we're done
        return float64(leftValue), 0, nil
    }

    // leading zeroes e.g. "0123" on right hand side?
    leadingZeros, leadingZerosBytes := acceptLeading(p.DecimalDigits[0], s)
    s = s [leadingZerosBytes:]

    offset, ok = numbers.Int.CheckedAdd(offset, leadingZerosBytes)
    if !ok { return 0, 0, strconv.ErrRange }

    // right-hand-side
    value, length, err := p.acceptFloatPart(s)
    if err != nil { return 0, 0, err }

    offset, ok = numbers.Int.CheckedAdd(offset, length)
    if !ok { return 0, 0, strconv.ErrRange }

    rightValue = value

    if rightValue > 0.0 {
        places := 1.0 + math.Floor(math.Log10(rightValue)) + float64(leadingZeros)
        rightValue *= math.Pow(0.1, places)
    }

    result, ok := numbers.Float64.CheckedAdd(leftValue, rightValue)
    if !ok { return 0, 0, strconv.ErrRange }

    return result, offset, nil
}

func (p *DecimalParser) AcceptSignedFloat(s string) (float64, int, error) {
    if (len(s) > 0) && (s[0] == '-') {
        value, length, err := p.AcceptFloat(s[1:])
        return -value, length + 1, err
    } else {
        return p.AcceptFloat(s)
    }
}

func (_ *DecimalParser) guessGroupSeparator(mp *message.Printer) (rune, bool) {
    // heuristic: any rune that appears at least twice is probably a comma
    s := mp.Sprint(number.Decimal(1234567890))
    return findRepeatingRune(s)
}

func (_ *DecimalParser) guessDecimalSeparator(p *message.Printer) (rune, bool) {
    // heuristic: any rune that is common to both these strings is probably a
    // decimal point. Concat the strings and find any repeated rune.
    s1 := p.Sprint(number.Decimal(1.23))
    s2 := p.Sprint(number.Decimal(4.56))
    s := s1 + s2
    return findRepeatingRune(s)
}

func (_ *DecimalParser) guessDecimalDigits(p *message.Printer, out *[10][]rune) bool {
    for i := 0; i < 10; i++ {
        s := []rune(p.Sprint(number.Decimal(i)))
        if len(s) != 1 { return false }
        out[i] = []rune{s[0]}
    }
    return true
}

// guessRuneValue guesses, for a rune representing a digit in a given locale,
// its value. For example, the rune '1' has integer value 1.
func (p *DecimalParser) guessRuneValue(d rune) (int, bool) {
    for i := 0; i < 10; i++ {
        if runeInSet(p.DecimalDigits[i], d) { return i, true }
    }
    return 0, false
}
