package parser

import (
    "fmt"
    "io"

    "github.com/tawesoft/golib/v2/css/internal/parser/item"
    "github.com/tawesoft/golib/v2/css/tokenizer"
    "github.com/tawesoft/golib/v2/css/tokenizer/token"
)

var (
    ErrUnexpectedEOF = fmt.Errorf("unexpected end of file")
)

type Parser struct {
    t *tokenizer.Tokenizer
    errors []error
}

func (p *Parser) Tokenizer() *tokenizer.Tokenizer {
    return p.t
}

func New(r io.Reader) *Parser {
    return &Parser{
        t: tokenizer.New(r),
    }
}

func (p *Parser) Errors() []error {
    return p.errors
}

func (p *Parser) ConsumeComponentValue(k *tokenizer.Tokenizer) item.ComponentValue {
    // Consume the next input token.
    t := k.Next()

    // If the current input token is a <{-token>, <[-token>, or <(-token>,
    // consume a simple block and return it.
    if TokenIsBlockStart(t) {
        return p.ConsumeBlock(k, t).ToComponentValue()
    }

    // Otherwise, if the current input token is a <function-token>, consume a function and return it.
    if t.Is(token.TypeFunction) {
        return p.ConsumeFunction(k, t).ToComponentValue()
    }

    // Otherwise, return the current input token.
    return item.PreservedToken(t).ToComponentValue()
}

// ConsumeBlock consumes a simple block. Note that this algorithm assumes that
// the current input token has already been checked to be an <{-token>,
// <[-token>, or <(-token>.
func (p *Parser) ConsumeBlock(k *tokenizer.Tokenizer, start token.Token) item.Block {
    // The ending token is the mirror variant of the current input token.
    // (E.g. if it was called with <[-token>, the ending token is <]-token>.)
    end := mirror(start)

    // Create a simple block with its associated token set to the current
    // input token and with its value initially set to an empty list.
    block := item.Block{
        Delim: start,
    }

    // Repeatedly consume the next input token and process it as follows:
    for {
        t := k.Next()
        if token.Equals(t, end) {
            return block
        } else if t.Is(token.TypeEOF) {
            // This is a parse error. Return the block.
            return block
        } else {
            // Reconsume the current input token.
            k.Push(t)
            // Consume a component value and append it to the value of the block.
            block.Value = append(block.Value, p.ConsumeComponentValue(k))
        }
    }
}

// ConsumeFunction consumes a function. Note: This algorithm assumes that the
// current input token has already been checked to be a <function-token>.
func (p *Parser) ConsumeFunction(k *tokenizer.Tokenizer, name token.Token) item.Function {
    // Create a function with its name equal to the value of the current input
    // token and with its value initially set to an empty list.
    f := item.Function{
        Name: name.StringValue(),
    }

    // Repeatedly consume the next input token and process it as follows:
    for {
        t := k.Next()
        if t.Is(token.TypeRightParen) {
            return f
        } else if t.Is(token.TypeEOF) {
            // This is a parse error. Return the function.
            return f
        } else {
            // Reconsume the current input token.
            k.Push(t)
            // Consume a component value and append the returned value to the
            // function’s value.
            f.Value = append(f.Value, p.ConsumeComponentValue(k))
        }
    }
}
