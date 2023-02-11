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
    "fmt"

    "github.com/tawesoft/golib/v2/text/number/plurals"
)

// Group defines a group of rule sets. Rule sets may refer to other rule sets
// in a Group by name, so think of a Group like a lexical scope in a
// programming language.
type Group struct {
    pluralRules plurals.Rules
    rulesets map[string]ruleset
}

type ruleType int
const (
    ruleTypeDefault = ruleType(iota)
    ruleTypeBaseValue
    ruleTypeBaseValueAndRadix
    ruleTypeNegativeNumber
    ruleTypeProperFraction
    ruleTypeImproperFraction
    ruleTypeInfinity
    ruleTypeNaN
)

type rule struct {
    // rule descriptor:
    Type ruleType
    Base int64
    Divisor int64

    // rule body:
    Subs [3]string
    Optional string
    Literal string
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
        rulesets: make(map[string]ruleset),
    }
    if err := g.parseGroups(rules); err != nil {
        return nil, fmt.Errorf("error parsing rbnf ruleset: %w", err)
    }
    return g, nil
}

type ruleset []rule // ordered
