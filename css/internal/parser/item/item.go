// Package item defines CSS items in the tree produced by a parser.
package item

import (
    "fmt"

    "github.com/tawesoft/golib/v2/css/tokenizer/token"
    "github.com/tawesoft/golib/v2/fun/maybe"
    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/operator"
)

type Type string
const (
    TypeAtRule         = Type("at-rule")
    TypeQualifiedRule  = Type("qualified rule")
    TypeDeclaration    = Type("declaration")
    TypeComponentValue = Type("component value")
    TypePreservedToken = Type("preserved token")
    TypeFunction       = Type("function")
    TypeBlock          = Type("block")
)

// An Item is a value in a tree produced by the parser.
type Item struct {
    _type Type // discriminates
    _cvType Type // component value type

    // TODO WithPosition

    name string
    list []ComponentValue
    block Block
    flag bool
    token token.Token
}

func (i Item) String() string {
    switch i._type {
        case TypeAtRule: return i.ToAtRule().String()
        case TypeQualifiedRule: return "<QualifiedRule(TODO)>"
        case TypeDeclaration: return "<Declaration(TODO)>"
        case TypeComponentValue: return i.ToComponentValue().String()
        case TypePreservedToken: return i.token.String()
        case TypeFunction:
            f := i.ToFunction()
            return fmt.Sprintf("<Function{Name: %q, Value: %v}>", f.Name, f.Value)
        case TypeBlock:
            block := i.ToBlock()
            return fmt.Sprintf("<Block{Delim: %s, Value: %v}>", block.Delim, block.Value)
        default:
            return fmt.Sprintf("<%s>", i._type)
    }
}

// AtRule describes an Item with a name, a prelude consisting of a list of
// component values, and an optional block consisting of a simple {} block.
type AtRule struct {
    Name string
    Prelude []ComponentValue
    Block maybe.M[Block]
}

func (i Item) ToAtRule() AtRule {
    must.Equal(i._type, TypeAtRule)
    var block maybe.M[Block]
    if i.flag {
        block = maybe.Some[Block](i.block)
    } else {
        block = maybe.Nothing[Block]()
    }
    return AtRule{
        Name:    i.name,
        Prelude: i.list,
        Block:   block,
    }
}

func (i AtRule) ToItem() Item {
    return Item{
        _type: TypeAtRule,
        name:  i.Name,
        list:  i.Prelude,
        block: i.Block.Or(operator.Zero[Block]()),
        flag:  i.Block.Ok,
    }
}

func (i AtRule) String() string {
    block, hasBlock := i.Block.Unpack()
    if hasBlock {
        return fmt.Sprintf("<AtRule{Name: %q, Prelude: %v, Block: %s}>",
            i.Name, i.Prelude, block.ToItem().String())
    } else {
        return fmt.Sprintf("<AtRule{Name: %q, Prelude: %v}>",
            i.Name, i.Prelude)
    }
}

// QualifiedRule describes an Item with a prelude consisting of a list of
// component values, and a block consisting of a simple {} block.
//
// Most qualified rules will be style rules, where the prelude is a selector
// and the block a list of declarations.
type QualifiedRule struct {
    Prelude []ComponentValue
    Block Block
}

func (i Item) ToQualifiedRule() QualifiedRule {
    must.Equal(i._type, TypeQualifiedRule)
    return QualifiedRule{
        Prelude: i.list,
        Block:   i.block,
    }
}

func (i QualifiedRule) ToItem() Item {
    return Item{
        _type: TypeQualifiedRule,
        list:  i.Prelude,
        block: i.Block,
    }
}

// ComponentValue describes an Item that is either one of the preserved tokens,
// a function, or a simple block.
type ComponentValue struct {
    _type Type // discriminates
    token PreservedToken
    function Function
    block Block
}

func (i Item) ToComponentValue() ComponentValue {
    must.Equal(i._type, TypeComponentValue)
    return ComponentValue{
        _type:   i._cvType,
        token:   PreservedToken(i.token),
        function: Function{
            Name:  i.name,
            Value: i.list,
        },
        block: i.block,
    }
}

func (i ComponentValue) ToItem() Item {
    return Item{
        _type: TypeComponentValue,
        _cvType: i._type,
        token: token.Token(i.token),
        name: i.function.Name,
        list: i.function.Value,
        block: i.block,
    }
}

func (i ComponentValue) ToPreservedToken() PreservedToken {
    must.Equal(i._type, TypePreservedToken)
    return i.token
}

func (i ComponentValue) ToFunction() Function {
    must.Equal(i._type, TypeFunction)
    return i.function
}

func (i ComponentValue) ToBlock() Block {
    must.Equal(i._type, TypeBlock)
    return i.block
}

// IsPreservedToken returns true iff the component value is both of type
// PreservedToken and the preserved token is equal to the specified token.
func (i ComponentValue) IsPreservedToken(t token.Token) bool {
    return (i._type == TypePreservedToken) &&
        token.Equals(token.Token(i.token), t)
}

func (i ComponentValue) String() string {
    switch i._type {
        case TypePreservedToken:
            return fmt.Sprintf("<ComponentValue:%s>", token.Token(i.token))
        case TypeFunction:
            f := i.ToFunction()
            return fmt.Sprintf("<ComponentValue/Function{Name: %q, Value: %v}>", f.Name, f.Value)
        case TypeBlock:
            block := i.ToBlock()
            return fmt.Sprintf("<ComponentValue/Block{Delim: %s, Value: %v}>", block.Delim, block.Value)
    }
    return "<ComponentValue/Invalid>"
}

// Declaration describes an Item that associates a property or descriptor name
// with a value. It has a name, a value consisting of a list of component
// values, and an important flag which is initially unset.
//
// Declarations are further categorized as property declarations or descriptor
// declarations, with the former setting CSS properties and appearing most
// often in qualified rules and the latter setting CSS descriptors, which
// appear only in at-rules. However, this categorization does not occur at the
// Syntax level; instead, it is a product of where the declaration appears.
type Declaration struct {
    Name string
    Value []ComponentValue
    Important bool
}

func (i Item) ToDeclaration() Declaration {
    must.Equal(i._type, TypeDeclaration)
    return Declaration{
        Name:      i.name,
        Value:     i.list,
        Important: i.flag,
    }
}

func (i Declaration) ToItem() Item {
    return Item{
        _type: TypeDeclaration,
        name:  i.Name,
        list:  i.Value,
        flag:  i.Important,
    }
}

// PreservedToken describes an Item that is any token produced by the tokenizer
// except for <function-token>, <{-token>, <(-token>, or <[-token>.
//
// Note: The tokens <}-token>s, <)-token>s, <]-token>, <bad-string-token>, and
// <bad-url-token> are always parse errors, but they are preserved in the
// token stream by this specification to allow other specs, such as Media
// Queries, to define more fine-grained error-handling than just dropping an
// entire declaration or block.
type PreservedToken token.Token

func (i Item) ToPreservedToken() PreservedToken {
    must.Equal(i._type, TypePreservedToken)
    return PreservedToken(i.token)
}

func (i PreservedToken) ToItem() Item {
    return Item{
        _type: TypePreservedToken,
        token: token.Token(i),
    }
}

func (i PreservedToken) ToComponentValue() ComponentValue {
    return ComponentValue{
        _type:    TypePreservedToken,
        token:    i,
    }
}

// Function describes an Item with a name and a value consisting of a list of
// component values.
type Function struct {
    Name string
    Value []ComponentValue
}

func (i Item) ToFunction() Function {
    must.Equal(i._type, TypeFunction)
    return Function{
        Name:  i.name,
        Value: i.list,
    }
}

func (i Function) ToItem() Item {
    return Item{
        _type: TypeFunction,
        name:  i.Name,
        list:  i.Value,
    }
}

func (i Function) ToComponentValue() ComponentValue {
    return ComponentValue{
        _type:    TypeFunction,
        function: i,
    }
}

// Block describes an Item that is a {}, (), [] block. It has a token that
// is either a  <{-token>, <(-token>, or <[-token>, and a value consisting of
// a list of component values.
type Block struct {
    Delim token.Token
    Value []ComponentValue
}

func (i Item) ToBlock() Block {
    must.Equal(i._type, TypeBlock)
    return Block{
        Delim: i.token,
        Value: i.list,
    }
}

func (i Block) ToItem() Item {
    return Item{
        _type: TypeBlock,
        token: i.Delim,
        list:  i.Value,
    }
}

func (i Block) ToComponentValue() ComponentValue {
    return ComponentValue{
        _type: TypeBlock,
        block: i,
    }
}
