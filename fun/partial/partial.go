// Package partial provides helpers for partial function application.
//
// Each function, f, in this package has a suffix for arity that indicates the
// number of arguments to the provided input function, i. Each function f
// returns a function g with a single argument that can be used to "partially
// apply" function i with one argument already bound.
//
// Functions have the prefix "Left" if they bind the left-most argument first,
// or "Right" if they bind the right-most argument first.
//
// For example,
//
//     f(a, b, c) == Left3(f)(a)(b, c) == Right3(f)(c)(a, b)
//
package partial

// Single takes a function with a single argument and return value and
// constructs a function that takes a single argument and returns a function
// that takes no arguments and returns a single value.
//
// For example,
//
//    opener :=  partial.Single(result.WrapFunc(os.Open))
//    openFoo := opener("foo.txt")
//    f, err := openFoo().Unpack()
func Single[T any, Return any](
    f func(t T) Return,
) func (t T) func () Return {
    return func (t T) func () Return {
        return func() Return {
            return f(t)
        }
    }
}

func Left2[A any, B any, Return any](
    f func(a A, b B) Return,
) func (a A) func (b B) Return {
    return func (a A) func (b B) Return {
        return func(b B) Return {
            return f(a, b)
        }
    }
}

func Left3[A any, B any, C any, Return any](
    f func(a A, b B, c C) Return,
) func (a A) func (b B, c C) Return {
    return func (a A) func (b B, c C) Return {
        return func(b B, c C) Return {
            return f(a, b, c)
        }
    }
}

func Left4[A any, B any, C any, D any, Return any](
    f func(a A, b B, c C, d D) Return,
) func (a A) func (b B, c C, d D) Return {
    return func (a A) func (b B, c C, d D) Return {
        return func(b B, c C, d D) Return {
            return f(a, b, c, d)
        }
    }
}

func Right2[A any, B any, Return any](
    f func(a A, b B) Return,
) func (b B) func (a A) Return {
    return func (b B) func (a A) Return {
        return func(a A) Return {
            return f(a, b)
        }
    }
}

func Right3[A any, B any, C any, Return any](
    f func(a A, b B, c C) Return,
) func (c C) func (a A, b B) Return {
    return func (c C) func (a A, b B) Return {
        return func(a A, b B) Return {
            return f(a, b, c)
        }
    }
}

func Right4[A any, B any, C any, D any, Return any](
    f func(a A, b B, c C, d D) Return,
) func (d D) func (a A, b B, c C) Return {
    return func (d D) func (a A, b B, c C) Return {
        return func(a A, b B, c C) Return {
            return f(a, b, c, d)
        }
    }
}
