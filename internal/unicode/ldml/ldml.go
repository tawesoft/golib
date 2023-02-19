// Package ldml parses LDML documents containing either language-dependent
// data or supplementary data.
//
// Note that this is not a full implementation.
//
// Current features:
// * Parsing ldml>identity (incomplete).
// * Parsing ldml>rbnf rules.
//
// TODO: ability to merge two documents with inheritance
//   see https://unicode.org/reports/tr35/#5-xml-format
package ldml

import (
    "encoding/xml"
    "fmt"
    "strings"
)

type DocType string
const (
    DocTypeLdml         = DocType("ldml")
    DocTypeSupplemental = DocType("supplementalData")
)

type String struct {
    Type string `xml:"type,attr"`
}

func (s String) String() string {
    return s.Type
}

type Document struct {
    Type                    DocType

    // One of the following, discriminated by Type...
    Ldml                    Ldml
    Supplemental            Supplemental
    LdmlBCP47               TODO // `xml:"ldmlBCP47"`
    Keyboard                TODO // `xml:"keyboard"`
    Platform                TODO // `xml:"platform"`
}

type TODO struct {}

type Supplemental struct {
    XMLName                 xml.Name `xml:"supplementalData"`
    Plurals                 Plurals `xml:"plurals"`
}

type Plurals struct {
    Type                    string `xml:"type,attr"`
    Rules                   []PluralRules `xml:"pluralRules"`
    Ranges                  []TODO `xml:"pluralRanges"`
}

type PluralRules struct {
    Locales                 string `xml:"locales,attr"`
    Rules                   []PluralRule `xml:"pluralRule"`
}

type PluralRule struct {
    Count                   string `xml:"count,attr"`
    Content                 string `xml:",chardata"`
}

type Ldml struct {
    XMLName                 xml.Name `xml:"ldml"`
    Language                String `xml:"identity>language"`
    Script                  String `xml:"identity>script"`
    Region                  String `xml:"identity>territory"`
    Variant                 String `xml:"identity>variant"`
    RbnfRulesetGroupings    []RbnfRulesetGrouping `xml:"rbnf>rulesetGrouping"`
    Numbers                 Numbers `xml:"numbers"`
}

type Numbers struct {
    Symbols                 []Symbols `xml:"symbols"`
    DecimalFormats          TODO
    ScientificFormats       TODO
    PercentFormats          TODO
    CurrencyFormats         TODO
    Currencies              TODO
}

type Alias struct {
    Source                  string `xml:"source,attr"`
    Path                    string `xml:"path,attr"`
}

type Symbols struct {
    Alias                   Alias  `xml:"alias"`
    NumberSystem            string `xml:"numberSystem,attr"`
    Decimal                 string `xml:"decimal"`
    Group                   string `xml:"group"`
    List                    string `xml:"list"`
    PercentSign             string `xml:"percentSign"`
    PlusSign                string `xml:"plusSign"`
    MinusSign               string `xml:"minusSign"`
    ApproximatelySign       string `xml:"approximatelySign"`
    Exponential             string `xml:"exponential"`
    SuperscriptingExponent  string `xml:"superscriptingExponent"`
    PerMille                string `xml:"perMille"`
    Infinity                string `xml:"infinity"`
    NaN                     string `xml:"nan"`
    CurrencyDecimal         string `xml:"currencyDecimal"`
    CurrencyGroup           string `xml:"currencyGroup"`
}

type RbnfRulesetGrouping struct {
    Type                    string `xml:"type,attr"`
    Rulesets                []RbnfRuleset `xml:"ruleset"`
}

type RbnfRuleset struct {
    Type                    string `xml:"type,attr"`
    Rules                   []RbnfRule `xml:"rbnfrule"`
    Access                  string `xml:"access,attr"`
}

func (rs RbnfRuleset) IsPrivate() bool {
    return strings.EqualFold(rs.Access, "private")
}

type RbnfRule struct {
    Value               string `xml:"value,attr"`
    Radix               string `xml:"radix,attr"`
    Content             string `xml:",chardata"`
}

// IcuStyle returns a rule formatted in the way described by the
// International Components for Unicode (ICU) software implementations
// ([ICU4C RuleBasedNumberFormat]) and ([ICU4J RuleBasedNumberFormat]), for
// example "name: body;".
//
// ICU4C RuleBasedNumberFormat: https://unicode-org.github.io/icu-docs/apidoc/released/icu4c/classicu_1_1RuleBasedNumberFormat.html
// ICU4J RuleBasedNumberFormat: https://unicode-org.github.io/icu-docs/apidoc/released/icu4j/com/ibm/icu/text/RuleBasedNumberFormat.html
func (r RbnfRule) IcuStyle() string {
    if r.Radix == "" {
        return fmt.Sprintf("%s: %s", r.Value, r.Content)
    } else {
        return fmt.Sprintf("%s/%s: %s", r.Value, r.Radix, r.Content)
    }
}

// Parse returns a Document parsed from an (XML-based) LDML or LDML
// Supplemental Data document.
func Parse(ldml []byte) (document Document, err error) {

    err = xml.Unmarshal(ldml, &document.Ldml)
    if err == nil { document.Type = DocTypeLdml; return }

    err = xml.Unmarshal(ldml, &document.Supplemental)
    if err == nil { document.Type = DocTypeSupplemental; return }

    err = fmt.Errorf("unrecognised document type")
    return
}
