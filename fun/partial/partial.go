// Package partial provides helpers for partial function application.
//
// Each function in this package has a suffix for arity that indicates the
// number of arguments of the provided input function f.
//
// Each Left/Right function returns a function with a single argument that can
// be used to "partially apply" input function f with one argument already
// bound.
//
// Functions have the prefix "Left" if they bind the left-most argument first,
// or "Right" if they bind the right-most argument first.
//
// Each input function must return exactly one value. The "fun/result" and
// "fun/maybe" packages provide useful functions that can wrap functions that
// return (value, ok) or (value, error) result types into single-value return
// variants.
//
// For example,
//
// f(a, b, c) => x becomes g(b, c) => x, with a bound, by calling Left3(f)(a).
//
// f(a, b, c) => x becomes g(a, b) => x, with c bound, by calling Right3(f)(c).
//
// f(a, b, c) => x becomes g() => x, with a, b, and c bound, by calling
// All3(f)(a, b, c).
//
// f(a) => x becomes g() => x, with a bound, by calling Single(f)(a).
//
package partial

// Single takes a function with a single argument and return value and
// constructs a function that takes a single argument and returns a function
// that takes no arguments and returns a single value.
//
// For example,
//
//    opener :=  partial.Single(result.WrapFunc(os.Open))
//    fooOpener := opener("foo.txt")
//    f, err := fooOpener().Unpack()
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

func All2[A any, B any, Return any](
    f func(A, B) Return,
) func(A, B) func() Return {
    return func (a A, b B) func () Return {
        return func() Return {
            return f(a, b)
        }
    }
}

func All3[A any, B any, C any, Return any](
    f func(A, B, C) Return,
) func(A, B, C) func() Return {
    return func (a A, b B, c C) func () Return {
        return func() Return {
            return f(a, b, c)
        }
    }
}

func All4[A any, B any, C any, D any, Return any](
    f func(A, B, C, D) Return,
) func(A, B, C, D) func() Return {
    return func (a A, b B, c C, d D) func () Return {
        return func() Return {
            return f(a, b, c, d)
        }
    }
}
