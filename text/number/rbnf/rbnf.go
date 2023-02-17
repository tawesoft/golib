// Package rbnf is a Go implementation of the Unicode Locale Data Markup
// Language (LDML) [Rule-Based Number Format (RBNF)].
//
// The RBNF can be used for complicated number formatting tasks, such as
// formatting a number of seconds as hours, minutes and seconds, or spelling
// out a number like 123 as "one hundred twenty-three", or adding an ordinal
// suffix to the end of a numeral like "123rd", or formatting numbers in a
// non-decimal number system such as Roman numerals or traditional Tamil
// numerals.
//
// This package does not implement any mapping from locale to specific rules.
// This must be handled at a higher layer.
//
// This package does not store any rules directly. You will have to obtain
// these from the Unicode Common Locale Data Repository (CLDR), or other
// sources, or define your own. Some rules CLDR rules for non-decimal number
// systems are implemented at [golib/v2/text/number/algorithmic].
//
// Rule-Based Number Format (RBNF): https://unicode.org/reports/tr35/tr35-numbers.html#6-rule-based-number-formatting
// [golib/v2/text/number/algorithmic]: https://github.com/tawesoft/golib/v2/text/number/algorithmic
//
// ## Security model
//
// It is assumed that the input rules come from a trusted author (e.g. the
// CLDR itself, or a trusted provider of localisation rules).
//
// ## Note of caution
//
// Quoting from the linked reference:
//
// "Where... CLDR plurals or ordinals can be used, their usage is recommended
// in preference to the RBNF data. First, the RBNF data is not completely
// fleshed out over all languages that otherwise have modern coverage.
// Secondly, the alternate forms are neither complete, nor useful without
// additional information. For example, for German there is
// spellout-cardinal-masculine, and spellout-cardinal-feminine. But a complete
// solution would have all genders (masculine/feminine/neuter), all cases
// (nominative, accusative, dative, genitive), plus context (with strong or
// weak determiner or none). Moreover, even for the alternate forms that do
// exist, CLDR does not supply any data for when to use one vs another (eg,
// when to use spellout-cardinal-masculine vs spellout-cardinal-feminine). So
// these data are inappropriate for general purpose software."
package rbnf

import (
    "errors"
    "fmt"
    "math"
    "strings"

    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/operator"
    "github.com/tawesoft/golib/v2/operator/checked/integer"
    "github.com/tawesoft/golib/v2/text/number/plurals"
    "github.com/tawesoft/golib/v2/text/number/rbnf/internal/body"
    "github.com/tawesoft/golib/v2/text/number/rbnf/internal/descriptor"
)

// Group defines a group of rule sets. Rule sets may refer to other rule sets
// in a Group by name, so think of a Group like a lexical scope in a
// programming language.
type Group struct {
    pluralRules plurals.Rules
    rulesets []ruleset
    rulesetNames map[string]int // index into rulesets
    stringData string
    descriptors []desc
    bodies []token
}

type ruleset struct {
    descriptorIdx int
    nRules int
    // isRegular bool
}

type desc struct {
    Base int64
    Divisor int64
    NumTokens uint8
    Type int8
    BodyIdx uint16
}

type token struct {
    Type uint8   // 0b00001111 Type + 0b0111000 SubstType
    Len uint8    // strlen or idx into rulesetNames
    Left uint16  // idx into stringData
}

// New returns a new rule-based number formatter formed from the group of
// rule sets described by the rules string.
//
// The plurals argument controls formatting of certain plural forms (cardinals
// and ordinals) used e.g. in spelling out "1st", "2nd", "3rd" or "1 cat",
// "2 cats", etc. If the ruleset does not contain any rules that use the
// cardinal syntax ("$(cardinal,plural syntax)$)") or ordinal syntax
// ("$(ordidinal,plural syntax)$)") then you may simply pass a nil Plural
// If specified, the methods implemented by the plural argument should
// usually match the same locale that the ruleset applies to.
//
// The rules string contains one or more rule sets in the format described by
// the International Components for Unicode (ICU) software implementations
// ([ICU4C RuleBasedNumberFormat]) and ([ICU4J RuleBasedNumberFormat]), e.g.:
// "%rulesetName: ruleName: ruleDescriptor; anotherRuleDescriptor: ruleBody;",
// with some differences:
//
// * In the ICU implementations, if a formatter only has one rule set, the name
//   may be omitted. In this implementation, the name is always required.
// * In the ICU implementations, a rule descriptor may be left out and have
//   an implicit meaning depending on the previous rule. In this implementation,
//   rule descriptors are always required (in any case, this doesn't appear
//   in the data files, regardless).
// * The ICU API documentation does not specify if a rule set name may appear
//   twice. In this implementation, this is treated as an error.
// * Only the following rule descriptors are supported (those not supported
//   do not seem to appear in the data files, regardless): "bv", "bv/rad",
//   "-x", "x.x", "0.x", "x.0", "Inf", "NaN".
// * For "x.x", "0.x", "x.0" rules, replacing the dot with a comma is not
//   supported (this does not seem to appear in the data files, regardless).
//   Note that this does not mean numbers cannot be *formatted* using commas,
//   only that they can not appear this way in a rule descriptor.
//
// Also note that a rule set is an ordered set.
//
// ICU4C RuleBasedNumberFormat: https://unicode-org.github.io/icu-docs/apidoc/released/icu4c/classicu_1_1RuleBasedNumberFormat.html
// ICU4J RuleBasedNumberFormat: https://unicode-org.github.io/icu-docs/apidoc/released/icu4j/com/ibm/icu/text/RuleBasedNumberFormat.html
func New(p plurals.Rules, rules string) (*Group, error) {
    g := &Group{
        pluralRules: p,
        rulesets: make([]ruleset, 0),
    }
    if err := g.parse(rules); err != nil {
        return nil, fmt.Errorf("error parsing rbnf ruleset group: %w", err)
    }
    return g, nil
}

func (g *Group) getRuleset(target string) (*ruleset, bool) {
    if n, ok := g.rulesetNames[target]; ok {
        return &(g.rulesets[n]), true
    } else {
        return nil, false
    }
}

func (g *Group) getRulesetIndex(target string) (int, bool) {
    if n, ok := g.rulesetNames[target]; ok {
        return n, true
    } else {
        return 0, false
    }
}

func (g *Group) addRuleset(name string) *ruleset {
    if _, exists := g.getRuleset(name); exists {
        panic(fmt.Errorf("duplicate ruleset name %q", name))
    }

    g.rulesetNames[name] = len(g.rulesets)
    g.rulesets = append(g.rulesets, ruleset{})
    return must.Ok(g.getRuleset(name))
}

func (g *Group) initRuleset(name string) *ruleset {
    rs := must.Ok(g.getRuleset(name))
    rs.descriptorIdx = len(g.descriptors)
    return rs
}

func (g *Group) addRuleDescriptor(rs *ruleset, str string) *desc {
    if len(g.bodies) > math.MaxUint16 {
        panic("too much rule body data in ruleset")
    }
    d := descriptor.Parse(str)
    rs.nRules++
    g.descriptors = append(g.descriptors, desc{
        Base:      d.Base,
        Divisor:   d.Divisor,
        Type:      int8(d.Type),
        BodyIdx:   uint16(len(g.bodies)),
    })
    return &(g.descriptors[len(g.descriptors) - 1])
}

func (g *Group) getString(tok token) string {
    if tok.Len == 0 { return "" }
    // TODO panic if wrong type
    return g.stringData[int(tok.Left) : int(tok.Left) + int(tok.Len)]
}

func (g *Group) getRule(rs *ruleset, condition func(d desc) bool) (desc, bool) {
    for i := 0; i < rs.nRules; i++ {
        d := g.descriptors[rs.descriptorIdx + i]
        if condition(d) {
            return g.descriptors[i], true
        }
    }
    return desc{}, false
}

// Errors returned by the Format methods.
var (
    ErrRange  = errors.New("value out of range")
    ErrNoRule = errors.New("no rule for this input")
    ErrNotImplemented = errors.New("rule logic not implemented for this input")
    ErrInvalidState = errors.New("invalid rule state")
)

func (g *Group) FormatInteger(rulesetName string, v int64) (string, error) {
    var sb strings.Builder
    rs := must.Ok(g.getRuleset(rulesetName))
    if err := g.formatInteger(&sb, rs, v, true); err == nil {
        return sb.String(), nil
    } else {
        return "", err
    }
}

func (g *Group) formatInteger(sb *strings.Builder, rs *ruleset, v int64, isRegular bool) error {
    if isRegular { // If the rule set is a regular rule set, do the following:

        // If the rule set includes a default rule (and the number was passed in as a
        // double), use the default rule. (If the number being formatted was passed in
        // as a long, the default rule is ignored.)
        // == ignored

        // If the number is negative, use the negative-number rule.
        if v < 0 {
            rule, ok := g.getRule(rs, func(d desc) bool {
                return descriptor.Type(d.Type) == descriptor.TypeNegativeNumber
            })
            if !ok { return ErrNoRule }
            return g.applyIntegerRule(sb, rs, rule, v, isRegular)
        }

        // If the number has a fractional part and is greater than 1, use the
        // improper fraction rule.
        // ...

        // If the number has a fractional part and is between 0 and
        // 1, use the proper fraction rule.
        // ...

        // Binary-search the rule list for the rule with the highest base value
        // less than or equal to the number.
        // TODO for now this is a linear search
        highestBaseValue := int64(-1)
        highestIdx := -1
        previousIdx := -1
        for i := 0; i < rs.nRules; i++ {
            d := g.descriptors[i]

            if !operator.In(descriptor.Type(d.Type),
                descriptor.TypeBaseValue, descriptor.TypeBaseValueAndRadix) { continue }

            if (d.Base <= v) && (d.Base > highestBaseValue) {
                highestBaseValue = d.Base
                highestIdx = i
                previousIdx = i-1
            }
        }
        if highestIdx < 0 { return ErrNoRule }
        rule := g.descriptors[highestIdx]

        // If that rule has two substitutions, its base value is not an even
        // multiple of its divisor, and the number is an even multiple of the
        // rule's divisor, use the rule that precedes it in the rule list.
        // Otherwise, use the rule itself.
        if (previousIdx >= 0) && (rule.Divisor > 0) {
            hasTwoSubs := func(r desc) bool {
                count := 0
                for i := 0; i < int(rule.NumTokens); i++ {
                    tok := g.bodies[int(rule.BodyIdx) + i]
                    _, st := decodeTokenType(tok.Type)
                    if body.SubstType(st) != body.SubstTypeNone {
                        count++
                    }
                }
                return count == 2
            }

            // base value is an even multiple?
            baseValueIsEvenMultiple := (rule.Base % rule.Divisor) == 0
            numberIsEvenMultuple := (v % rule.Divisor) == 0
            if hasTwoSubs(rule) && (!baseValueIsEvenMultiple) && numberIsEvenMultuple {
                rule = g.descriptors[previousIdx]
            }
        }
        return g.applyIntegerRule(sb, rs, rule, v, isRegular)
    }

    /*
    If the rule set is a fraction rule set, do the following:

    Ignore negative-number and fraction rules.

    For each rule in the list, multiply the number being formatted (which will
    always be between 0 and 1) by the rule's base value. Keep track of the distance
    between the result the nearest integer.

    Use the rule that produced the result closest to zero in the above calculation.
    In the event of a tie or a direct hit, use the first matching rule encountered.
    (The idea here is to try each rule's base value as a possible denominator of a
    fraction. Whichever denominator produces the fraction closest in value to the
    number being formatted wins.) If the rule following the matching rule has the
    same base value, use it if the numerator of the fraction is anything other than
    1; if the numerator is 1, use the original matching rule. (This is to allow
    singular and plural forms of the rule text without a lot of extra hassle.)
    */
    return nil
}

func divisor_int_log10(v int64) int64 {
    if v == 0 { return 1 }
    digits := int(math.Log10(float64(v))) // e.g. 900 => 2.95 => 2
    return int64(0.5 + math.Pow10(digits)) // e.g. 2 => 10^2 => 100
}

func isNormalRule(rule desc) bool {
    return operator.In(descriptor.Type(rule.Type),
        descriptor.TypeBaseValue, descriptor.TypeBaseValueAndRadix)
}

func (g *Group) applyIntegerRule(sb *strings.Builder, rs *ruleset, rule desc, v int64, isRegular bool) error {
    isOptional := false

    for i := 0; i < int(rule.NumTokens); i++ {
        tok := g.bodies[int(rule.BodyIdx) + i]

        ty, _ := decodeTokenType(tok.Type)
        if ty == body.TypeOptionalEnd {
            isOptional = false
            continue
        } else if isOptional {
            switch descriptor.Type(rule.Type) {
                case descriptor.TypeBaseValue:
                    // Omit the optional text if the number is
                    // an even multiple of the rule's divisor
                    d := divisor_int_log10(v)
                    if (v % d) == 0 { continue }
                case descriptor.TypeBaseValueAndRadix:
                    // Omit the optional text if the number is
                    // an even multiple of the rule's divisor
                    return ErrNotImplemented
                case descriptor.TypeNegativeNumber:
                    return ErrInvalidState
                case descriptor.TypeProperFraction:
                    return ErrInvalidState
                case descriptor.TypeDefault:
                    // Omit the optional text if the number is an integer (same
                    // as specifying both an x.x rule and an x.0 rule)
                    continue
                case descriptor.TypeImproperFraction:
                    // Omit the optional text if the number is between 0 and 1
                    // (same as specifying both an x.x rule and a 0.x rule)
                    if (v == 0) || (v == 1) { continue }
                default:
                    if !isRegular {
                        // Omit the optional text if multiplying the number by
                        // the rule's base value yields 1.
                        if rule.Base * v == 1 { continue }
                    }
            }
        }

        switch ty {
            case body.TypeOptionalStart:
                isOptional = true
            case body.TypeLiteral:
                sb.WriteString(g.getString(tok))

            case body.TypeSubstLeftArrow:
                if !isRegular {
                    // Multiply the number by the rule's base value and
                    // format the result.
                    n, ok := integer.Int64.Mul(v, rule.Base)
                    if !ok { return ErrRange }
                    err := g.formatInteger(sb, rs, n, isRegular)
                    if err != nil { return err }
                }
                switch descriptor.Type(rule.Type) {
                    case descriptor.TypeBaseValue:
                        // Divide the number by the rule's divisor and format
                        // the quotient
                        n := v / divisor_int_log10(v)
                        // todo subst type select ruleset
                        err := g.formatInteger(sb, rs, n, isRegular)
                        if err != nil { return err }

                    case descriptor.TypeBaseValueAndRadix:
                        // Divide the number by the rule's divisor and format the remainder
                        // The rule's divisor is the highest power of rad less than or equal to the base value.
                        return ErrNotImplemented
                    case descriptor.TypeNegativeNumber:
                        return ErrInvalidState
                    case descriptor.TypeDefault: fallthrough
                    case descriptor.TypeProperFraction: fallthrough
                    case descriptor.TypeImproperFraction:
                        // Isolate the number's integral part and format it.
                        // Here its already integer part.
                        return ErrInvalidState
                    default:
                        return ErrInvalidState
                }

            case body.TypeSubstRightArrow:
                switch descriptor.Type(rule.Type) {
                    case descriptor.TypeBaseValue:
                        // Divide the number by the rule's divisor and format the remainder
                        // The rule's divisor is the highest power of 10 less than or equal to the base value.
                        n := v % divisor_int_log10(v)
                        // todo subst type select ruleset
                        err := g.formatInteger(sb, rs, n, isRegular)
                        if err != nil { return err }

                    case descriptor.TypeBaseValueAndRadix:
                        // Divide the number by the rule's divisor and format the remainder
                        // The rule's divisor is the highest power of rad less than or equal to the base value.
                        return ErrNotImplemented

                    case descriptor.TypeNegativeNumber:
                        abs, ok := integer.Int64.Abs(v)
                        if !ok { return ErrRange }
                        // todo subst type select ruleset
                        err := g.formatInteger(sb, rs, abs, isRegular)
                        if err != nil { return err }

                    case descriptor.TypeDefault: fallthrough
                    case descriptor.TypeProperFraction:
                        // Isolate the number's fractional part and format it.
                        return ErrNotImplemented

                    // in rule in fraction rule set: Not allowed.

                    // other???
                    default: return ErrInvalidState
                }

            // other???
            default:
                return ErrNotImplemented
        }
    }

    return nil
}
