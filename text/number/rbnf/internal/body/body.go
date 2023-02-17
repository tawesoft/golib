// Package body implements parsing of a rbnf rule body.
//
// See: https://unicode-org.github.io/icu-docs/apidoc/released/icu4c/classicu_1_1RuleBasedNumberFormat.html#details
//
// Token format:
//     substitution = bracket-start substitution-descriptor bracket-end
//     part = substitution or literal
//     item = part OR ('[' part* ']') *
//     substitution = a rule set name starting with %
//                    OR a DecimalFormat pattern starting with '0' or '#'.
//                    OR empty
//
// bracket-start, bracket-end pairs are:
//     →, →
//     ←, ←
//     =, =
//     $(, )$
package body

import (
    "fmt"
    "strings"
    "unicode/utf8"

    "github.com/tawesoft/golib/v2/must"
)

type Type uint8
const (
    _typeNone = Type(iota)
    TypeLiteral             // literal text data e.g. " million"
    TypeSubstRightArrow     // →...→, simple
    TypeSubstLeftArrow      // ←...←, simple
    TypeSubstEqualsSign     // =...=, simple
    TypeOptionalStart       // "[" in [...], may contain substitutions and literals
    TypeOptionalEnd         // "]" in [...], may contain substitutions and literals
    TypeSubstPluralCardinal // $(cardinal,plural-syntax)$
    TypeSubstPluralOrdinal  // $(ordinal,plural-syntax)$
    TypeEOF
)

type SubstType uint8
const (
    SubstTypeNone = SubstType(iota)
    SubstTypeEmpty
    SubstTypeRulesetName
    SubstTypeDecimalFormat
    SubstTypeInvalid
)

type Slice [2]int // byte offsets into the input string
func (s Slice) Len() int {
    return s[1] - s[0]
}
func (s Slice) Of(str string) string {
    return str[s[0]:s[1]]
}

// Token represents one broken-down component of a rule body. Where it includes
// string data (such as a substitution descriptor or literal text), this is
// represented by a Slice containing indexes into the original string initially
// passed to the tokenizer.
//
// Tokens produced by a stream either contain an empty slice, or the slice
// indexes are non-overlapping and the start index is greater than or equal to
// any previous token's slice end index.
type Token struct {
    Type Type
    Content Slice
}

// SimpleSubstType returns the type of substitution descriptor appearing
// between delimiters e.g. "→→" (empty), "=%RulesetName=" "=#DecimalFormat=".
// The string argument is the rule body string that the token content slice
// applies to.
func (t Token) SimpleSubstType(s string) SubstType {
    for {
        switch t.Type {
            case TypeSubstRightArrow: fallthrough
            case TypeSubstLeftArrow:  fallthrough
            case TypeSubstEqualsSign: {
                if t.Content.Len() == 0 {
                    if (t.Type == TypeSubstEqualsSign) {
                        return SubstTypeInvalid // not allowed here
                    }
                    return SubstTypeEmpty
                }
                c := s[t.Content[0]] // only need first byte
                switch c {
                    case '%': return SubstTypeRulesetName
                    case '0': fallthrough
                    case '#': return SubstTypeDecimalFormat
                    default: return SubstTypeInvalid
                }
            }
            default: return SubstTypeNone
        }
    }
}

func decode(s string) (rune, int) {
    r, z := utf8.DecodeRuneInString(s)
    if (r == utf8.RuneError) && (z == 1) {
        panic(fmt.Errorf("unexpected end"))
    }
    return r, z
}

type tokenizer struct {
    str string
    start int

    // "A rule's body consists of... zero, one, or two substitution tokens,
    // and a range of text in brackets. The brackets denote optional text
    // (and may also include one or both substitutions). The rest of the
    // text... is literal text."
    seenBrackets int // 0: none, 1: seen opening, 2: seen close
    seenSubs     int // max 2
}

func NewTokenizer(s string) (next func() (Token, error)) {
    t := tokenizer{str: s}
    return t.Next
}

// Next returns the next token, also validating that the token is syntactically
// valid at this point in the stream. Note that a token may still be erroneous
// depending on the rule descriptor but that requires more context.
func (t *tokenizer) Next() (rettok Token, reterr error) {
    defer func() {
        if r := recover(); r != nil {
            rettok = Token{}
            if err, ok := r.(error); ok {
                reterr = err
            } else {
                reterr = fmt.Errorf("unknown error: %v", r)
            }
        }
    }()

    tok := t.next()
    switch tok.Type {
        case TypeEOF:
            if t.seenBrackets == 1 {
                panic(fmt.Errorf("expected ']'"))
            }

        case TypeOptionalStart:
            if t.seenBrackets != 0 {
                panic(fmt.Errorf("unexpected '['"))
            }
            t.seenBrackets = 1

        case TypeOptionalEnd:
            if t.seenBrackets != 1 {
                panic(fmt.Errorf("unexpected ']'"))
            }
            t.seenBrackets = 2

        case TypeLiteral:
            break

        case TypeSubstEqualsSign:
            fallthrough
        case TypeSubstLeftArrow:
            fallthrough
        case TypeSubstRightArrow:
            if t.seenSubs > 2 {
                panic(fmt.Errorf("too many substitutions"))
            }
            if tok.SimpleSubstType(t.str) == SubstTypeInvalid {
                panic(fmt.Errorf("invalid substitution descriptor %q", tok.Content.Of(t.str)))
            }
            t.seenSubs++

        case TypeSubstPluralCardinal: fallthrough
        case TypeSubstPluralOrdinal:
            // TODO validate, parse
            if t.seenSubs > 2 {
                panic(fmt.Errorf("too many substitutions"))
            }
            t.seenSubs++
        default:
            must.Never()
    }
    return tok, nil
}

// next returns the next token without additional validation.
func (t *tokenizer) next() Token {
    tok, end := next(t.str, t.start)
    t.start = end
    return tok
}

func next(s string, start int) (Token, int) {
    current, z := decode(s[start:])
    if z == 0 { return Token{Type: TypeEOF}, 0 }

    for {
        switch current {
            case '→': fallthrough
            case '←': fallthrough
            case '=': return consumeSimpleSubstitution(s, start + z, current)
            case '[': return Token{Type: TypeOptionalStart}, start + z
            case ']': return Token{Type: TypeOptionalEnd},   start + z
            default: {
                // $(cardinal|ordinal,plural-syntax)$
                if current == '$' {
                    peek, _ := decode(s[start + z:])
                    if peek == '(' {
                        return consumePluralSubstitution(s, start + z + 1)
                    }
                }
                // literal
                return consumeLiteral(s, start)
            }
        }
    }
}

// consumes e.g. " million", stopping at any special sequence
func consumeLiteral(s string, start int) (Token, int) {
    end := start
    for {
        next, z := decode(s[end:])
        if z == 0 { break }

        if strings.ContainsRune("→←=[", next) { break }

        peek, _ := decode(s[end+z:])
        if (next == '$') && (peek == '(') { break }

        end += z
    }
    return Token{Type: TypeLiteral, Content: Slice{start, end}}, end
}

// consumes e.g. →...→
func consumeSimpleSubstitution(s string, start int, terminator rune) (Token, int) {
    subst := func(terminator rune) Type {
        switch terminator {
            case '=': return TypeSubstEqualsSign
            case '→': return TypeSubstRightArrow
            case '←': return TypeSubstLeftArrow
            default: must.Never(); return 0
        }
    }(terminator)

    end := start
    for {
        next, z := decode(s[end:])
        if z == 0 { panic(fmt.Errorf("expected closing %c", terminator)) }
        if next == terminator { break }
        end += z
    }
    return Token{Type: subst, Content: Slice{start, end}}, end + utf8.RuneLen(terminator)
}

// consumes e.g. $(cardinal|ordinal,plural-syntax)$
func consumePluralSubstitution(s string, start int) (Token, int) {
    end := start
    for {
        next, z := decode(s[end:])
        if z == 0 { panic(fmt.Errorf("expected closing )$")) }

        peek, _ := decode(s[end+z:])
        if (next == ')') && (peek == '$') { break }
        end += z
    }

    var ptype Type
    str := s[start:end]
    idx := strings.IndexRune(str, ',')
    if idx > 0 {
        left := str[0:idx]
        if strings.EqualFold(left, "cardinal") {
            ptype = TypeSubstPluralCardinal
        } else if strings.EqualFold(left, "ordinal") {
            ptype = TypeSubstPluralOrdinal
        }
        start += idx + 1
    }

    if ptype == Type(0) {
        panic(fmt.Errorf("invalid plural rule type %q", str))
    }

    return Token{Type: ptype, Content: Slice{start, end}}, end + len(")$")
}
