package operator

// Code generated by (tawesoft.co.uk/go/operator) template-numbers.py: DO NOT EDIT.

// Some overflow checks with reference to stackoverflow.com/a/1514309/5654201

type int64Unary struct {
    Identity        func(int64) int64
    Abs             func(int64) int64
    Negation        func(int64) int64
    Zero            func(int64) bool
    NonZero         func(int64) bool
    Positive        func(int64) bool
    Negative        func(int64) bool
}

type int64UnaryChecked struct {
    Abs             func(int64) (int64, error)
    Negation        func(int64) (int64, error)
}

type int64Binary struct {
    Add             func(int64, int64) int64
    Sub             func(int64, int64) int64
    Mul             func(int64, int64) int64
    Div             func(int64, int64) int64
    Mod             func(int64, int64) int64
    
    Eq              func(int64, int64) bool
    Neq             func(int64, int64) bool
    Lt              func(int64, int64) bool
    Lte             func(int64, int64) bool
    Gt              func(int64, int64) bool
    Gte             func(int64, int64) bool
    
    And             func(int64, int64) int64
    Or              func(int64, int64) int64
    Xor             func(int64, int64) int64
    AndNot          func(int64, int64) int64
    
    Shl             func(int64, uint) int64
    Shr             func(int64, uint) int64
}

type int64BinaryChecked struct {
    Add             func(int64, int64) (int64, error)
    Sub             func(int64, int64) (int64, error)
    Mul             func(int64, int64) (int64, error)
    Div             func(int64, int64) (int64, error)
    
    Shl             func(int64, uint) (int64, error)
    Shr             func(int64, uint) (int64, error)
}

type int64Nary struct {
    Add             func(... int64) int64
    Sub             func(... int64) int64
    Mul             func(... int64) int64
}

type int64NaryChecked struct {
    Add             func(... int64) (int64, error)
    Sub             func(... int64) (int64, error)
    Mul             func(... int64) (int64, error)
}

// Int64 implements operations on one (unary), two (binary), or many (nary) arguments of type int64.
var Int64 = struct {
    Unary           int64Unary
    Binary          int64Binary
    Nary            int64Nary
    Reduce          func(operatorIdentity int64, operator func(int64, int64) int64, elements ... int64) int64
}{
    Unary:          int64Unary{
        Identity:   func(a int64) int64 { return a },
        Abs:        int64UnaryAbs,
        Negation:   func(a int64) int64 { return -a },
        Zero:       func(a int64) bool { return a == 0 },
        NonZero:    func(a int64) bool { return a != 0 },
        Positive:   int64UnaryPositive,
        Negative:   int64UnaryNegative,
    },
    
    Binary:          int64Binary{
        Add:        func(a int64, b int64) int64 { return a + b },
        Sub:        func(a int64, b int64) int64 { return a - b },
        Mul:        func(a int64, b int64) int64 { return a * b },
        Div:        func(a int64, b int64) int64 { return a / b },
        
        Eq:         func(a int64, b int64) bool { return a == b },
        Neq:        func(a int64, b int64) bool { return a != b },
        Lt:         func(a int64, b int64) bool { return a <  b },
        Lte:        func(a int64, b int64) bool { return a <= b },
        Gt:         func(a int64, b int64) bool { return a >  b },
        Gte:        func(a int64, b int64) bool { return a >= b },
        
        And:        func(a int64, b int64) int64 { return a & b },
        Or:         func(a int64, b int64) int64 { return a | b },
        Xor:        func(a int64, b int64) int64 { return a ^ b },
        AndNot:     func(a int64, b int64) int64 { return a &^ b },
        Mod:        func(a int64, b int64) int64 { return a % b },
        
        Shl:        func(a int64, b uint) int64 { return a << b },
        Shr:        func(a int64, b uint) int64 { return a >> b },
    },
    
    Nary:           int64Nary{
        Add:        int64NaryAdd,
        Mul:        int64NaryMul,
    },
    
    Reduce:         int64Reduce,
}

// Int64Checked implements operations on one (unary), two (binary), or many (nary) arguments of type int64, returning an
// error in cases such as overflow or an undefined operation.
var Int64Checked = struct {
    Unary           int64UnaryChecked
    Binary          int64BinaryChecked
    Nary            int64NaryChecked
    Reduce          func(operatorIdentity int64, operator func(int64, int64) (int64, error), elements ... int64) (int64, error)
}{
    Unary:          int64UnaryChecked{
        Abs:        int64UnaryCheckedAbs,
        Negation:   int64UnaryCheckedNegation,
    },
    
    Binary:         int64BinaryChecked{
        Add:        int64BinaryCheckedAdd,
        Sub:        int64BinaryCheckedSub,
        Mul:        int64BinaryCheckedMul,
        Div:        int64BinaryCheckedDiv,
        Shl:        int64BinaryCheckedShl,
    },
    
    Nary:           int64NaryChecked{
        Add:        int64NaryCheckedAdd,
        Mul:        int64NaryCheckedMul,
    },
    
    Reduce:         int64CheckedReduce,
}

func int64UnaryPositive(a int64) bool {
    return a > 0
}

func int64UnaryNegative(a int64) bool {
    return a < 0
}

func int64UnaryAbs(a int64) int64 {
    if a < 0 { return -a }
    return a
}

func int64UnaryCheckedAbs(a int64) (v int64, err error) {
    if a == minInt64 { return v, ErrorOverflow }
    if a < 0 { return -a, nil }
    return a, nil
}

func int64UnaryCheckedNegation(a int64) (v int64, err error) {
    if (a == minInt64) { return v, ErrorOverflow }
    return -a, nil
}

func int64BinaryCheckedAdd(a int64, b int64) (v int64, err error) {
    if (b > 0) && (a > (maxInt64 - b)) { return v, ErrorOverflow }
    if (b < 0) && (a < (minInt64 - b)) { return v, ErrorOverflow }
    return a + b, nil
}

func int64BinaryCheckedSub(a int64, b int64) (v int64, err error) {
    if (b < 0) && (a > (maxInt64 + b)) { return v, ErrorOverflow }
    if (b > 0) && (a < (minInt64 + b)) { return v, ErrorOverflow }
    return a - b, nil
}

func int64BinaryCheckedMul(a int64, b int64) (v int64, err error) {
    if (a == -1) && (b == minInt64) { return v, ErrorOverflow }
    if (b == -1) && (a == minInt64) { return v, ErrorOverflow }
    if (a > (maxInt64 / b)) { return v, ErrorOverflow }
    if (a < (minInt64 / b)) { return v, ErrorOverflow }
    
    return a * b, nil
}

func int64BinaryCheckedDiv(a int64, b int64) (v int64, err error) {
    if (b == -1) && (a == minInt64) { return v, ErrorOverflow }
    if (b == 0) { return v, ErrorUndefined }
    
    return a / b, nil
}

func int64BinaryCheckedShl(a int64, b uint) (v int64, err error) {
    if a < 0 { return v, ErrorUndefined }
    if b > uint(int64MostSignificantBit(maxInt64)) { return v, ErrorOverflow }
    return v, err
}

func int64MostSignificantBit(a int64) (result int) {
  for a > 0 {
      a >>= 1
      result++
  }
  return result;
}

func int64NaryAdd(xs ... int64) (result int64) {
    for i := 0; i < len(xs); i++ {
        result += xs[i]
    }
    return result
}

func int64NaryCheckedAdd(xs ... int64) (result int64, err error) {
    for i := 0; i < len(xs); i++ {
        result, err = int64BinaryCheckedAdd(result, xs[i])
        if err != nil { return result, err }
    }
    return result, nil
}

func int64NaryMul(xs ... int64) (result int64) {
    result = 1
    for i := 0; i < len(xs); i++ {
        result *= xs[i]
    }
    return result
}

func int64NaryCheckedMul(xs ... int64) (result int64, err error) {
    result = 1
    for i := 0; i < len(xs); i++ {
        result, err = int64BinaryCheckedMul(result, xs[i])
        if err != nil { return result, err }
    }
    return result, nil
}

func int64Reduce(operatorIdentity int64, operator func(int64, int64) int64, elements ... int64) (result int64) {
    result = operatorIdentity
    for i := 0; i < len(elements); i++ {
        result = operator(result, elements[i])
    }
    return result
}

func int64CheckedReduce(operatorIdentity int64, operator func(int64, int64) (int64, error), elements ... int64) (result int64, err error) {
    result = operatorIdentity
    for i := 0; i < len(elements); i++ {
        result, err = operator(result, elements[i])
        if err != nil { return result, err }
    }
    return result, err
}

