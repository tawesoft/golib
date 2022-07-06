package operator

// Code generated by (tawesoft.co.uk/go/operator) template-numbers.py: DO NOT EDIT.


// Some overflow checks with reference to stackoverflow.com/a/1514309/5654201

type uint8Unary struct {
    Identity        func(uint8) uint8
    Not             func(uint8) uint8
    Zero            func(uint8) bool
    NonZero         func(uint8) bool
    Positive        func(uint8) bool
    Negative        func(uint8) bool
}

type uint8UnaryChecked struct {
}

type uint8Binary struct {
    Add             func(uint8, uint8) uint8
    Sub             func(uint8, uint8) uint8
    Mul             func(uint8, uint8) uint8
    Div             func(uint8, uint8) uint8
    Mod             func(uint8, uint8) uint8
    
    Eq              func(uint8, uint8) bool
    Neq             func(uint8, uint8) bool
    Lt              func(uint8, uint8) bool
    Lte             func(uint8, uint8) bool
    Gt              func(uint8, uint8) bool
    Gte             func(uint8, uint8) bool
    
    And             func(uint8, uint8) uint8
    Or              func(uint8, uint8) uint8
    Xor             func(uint8, uint8) uint8
    AndNot          func(uint8, uint8) uint8
    
    Shl             func(uint8, uint) uint8
    Shr             func(uint8, uint) uint8
}

type uint8BinaryChecked struct {
    Add             func(uint8, uint8) (uint8, error)
    Sub             func(uint8, uint8) (uint8, error)
    Mul             func(uint8, uint8) (uint8, error)
    Div             func(uint8, uint8) (uint8, error)
    
    Shl             func(uint8, uint) (uint8, error)
    Shr             func(uint8, uint) (uint8, error)
}

type uint8Nary struct {
    Add             func(... uint8) uint8
    Sub             func(... uint8) uint8
    Mul             func(... uint8) uint8
}

type uint8NaryChecked struct {
    Add             func(... uint8) (uint8, error)
    Sub             func(... uint8) (uint8, error)
    Mul             func(... uint8) (uint8, error)
}

// Uint8 implements operations on one (unary), two (binary), or many (nary) arguments of type uint8.
var Uint8 = struct {
    Unary           uint8Unary
    Binary          uint8Binary
    Nary            uint8Nary
    Reduce          func(operatorIdentity uint8, operator func(uint8, uint8) uint8, elements ... uint8) uint8
}{
    Unary:          uint8Unary{
        Identity:   func(a uint8) uint8 { return a },
        Not:        func(a uint8) uint8 { return ^a },
        Zero:       func(a uint8) bool { return a == 0 },
        NonZero:    func(a uint8) bool { return a != 0 },
        Positive:   uint8UnaryPositive,
        Negative:   uint8UnaryNegative,
    },
    
    Binary:          uint8Binary{
        Add:        func(a uint8, b uint8) uint8 { return a + b },
        Sub:        func(a uint8, b uint8) uint8 { return a - b },
        Mul:        func(a uint8, b uint8) uint8 { return a * b },
        Div:        func(a uint8, b uint8) uint8 { return a / b },
        
        Eq:         func(a uint8, b uint8) bool { return a == b },
        Neq:        func(a uint8, b uint8) bool { return a != b },
        Lt:         func(a uint8, b uint8) bool { return a <  b },
        Lte:        func(a uint8, b uint8) bool { return a <= b },
        Gt:         func(a uint8, b uint8) bool { return a >  b },
        Gte:        func(a uint8, b uint8) bool { return a >= b },
        
        And:        func(a uint8, b uint8) uint8 { return a & b },
        Or:         func(a uint8, b uint8) uint8 { return a | b },
        Xor:        func(a uint8, b uint8) uint8 { return a ^ b },
        AndNot:     func(a uint8, b uint8) uint8 { return a &^ b },
        Mod:        func(a uint8, b uint8) uint8 { return a % b },
        
        Shl:        func(a uint8, b uint) uint8 { return a << b },
        Shr:        func(a uint8, b uint) uint8 { return a >> b },
    },
    
    Nary:           uint8Nary{
        Add:        uint8NaryAdd,
        Mul:        uint8NaryMul,
    },
    
    Reduce:         uint8Reduce,
}

// Uint8Checked implements operations on one (unary), two (binary), or many (nary) arguments of type uint8, returning an
// error in cases such as overflow or an undefined operation.
var Uint8Checked = struct {
    Unary           uint8UnaryChecked
    Binary          uint8BinaryChecked
    Nary            uint8NaryChecked
    Reduce          func(operatorIdentity uint8, operator func(uint8, uint8) (uint8, error), elements ... uint8) (uint8, error)
}{
    Unary:          uint8UnaryChecked{
    },
    
    Binary:         uint8BinaryChecked{
        Add:        uint8BinaryCheckedAdd,
        Sub:        uint8BinaryCheckedSub,
        Mul:        uint8BinaryCheckedMul,
        Div:        uint8BinaryCheckedDiv,
        Shl:        uint8BinaryCheckedShl,
    },
    
    Nary:           uint8NaryChecked{
        Add:        uint8NaryCheckedAdd,
        Mul:        uint8NaryCheckedMul,
    },
    
    Reduce:         uint8CheckedReduce,
}

func uint8UnaryPositive(a uint8) bool {
    return a > 0
}

func uint8UnaryNegative(a uint8) bool {
    return a < 0
}




func uint8BinaryCheckedAdd(a uint8, b uint8) (v uint8, err error) {
    if (b > 0) && (a > (maxUint8 - b)) { return v, ErrorOverflow }
    if (b < 0) && (a < (minUint8 - b)) { return v, ErrorOverflow }
    return a + b, nil
}

func uint8BinaryCheckedSub(a uint8, b uint8) (v uint8, err error) {
    if (b < 0) && (a > (maxUint8 + b)) { return v, ErrorOverflow }
    if (b > 0) && (a < (minUint8 + b)) { return v, ErrorOverflow }
    return a - b, nil
}

func uint8BinaryCheckedMul(a uint8, b uint8) (v uint8, err error) {
    if (a > (maxUint8 / b)) { return v, ErrorOverflow }
    if (a < (minUint8 / b)) { return v, ErrorOverflow }
    
    return a * b, nil
}

func uint8BinaryCheckedDiv(a uint8, b uint8) (v uint8, err error) {
    if (b == 0) { return v, ErrorUndefined }
    
    return a / b, nil
}

func uint8BinaryCheckedShl(a uint8, b uint) (v uint8, err error) {
    if b > uint(uint8MostSignificantBit(maxUint8)) { return v, ErrorOverflow }
    return v, err
}

func uint8MostSignificantBit(a uint8) (result int) {
  for a > 0 {
      a >>= 1
      result++
  }
  return result;
}

func uint8NaryAdd(xs ... uint8) (result uint8) {
    for i := 0; i < len(xs); i++ {
        result += xs[i]
    }
    return result
}

func uint8NaryCheckedAdd(xs ... uint8) (result uint8, err error) {
    for i := 0; i < len(xs); i++ {
        result, err = uint8BinaryCheckedAdd(result, xs[i])
        if err != nil { return result, err }
    }
    return result, nil
}

func uint8NaryMul(xs ... uint8) (result uint8) {
    result = 1
    for i := 0; i < len(xs); i++ {
        result *= xs[i]
    }
    return result
}

func uint8NaryCheckedMul(xs ... uint8) (result uint8, err error) {
    result = 1
    for i := 0; i < len(xs); i++ {
        result, err = uint8BinaryCheckedMul(result, xs[i])
        if err != nil { return result, err }
    }
    return result, nil
}

func uint8Reduce(operatorIdentity uint8, operator func(uint8, uint8) uint8, elements ... uint8) (result uint8) {
    result = operatorIdentity
    for i := 0; i < len(elements); i++ {
        result = operator(result, elements[i])
    }
    return result
}

func uint8CheckedReduce(operatorIdentity uint8, operator func(uint8, uint8) (uint8, error), elements ... uint8) (result uint8, err error) {
    result = operatorIdentity
    for i := 0; i < len(elements); i++ {
        result, err = operator(result, elements[i])
        if err != nil { return result, err }
    }
    return result, err
}
