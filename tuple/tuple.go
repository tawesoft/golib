// Package tuple simplifies packing and unpacking function arguments and
// results into generic tuple types.
package tuple

type T2[A any, B any] struct{
    A A
    B B
}

func ToT2[A any, B any](a A, b B) T2[A, B] {
    return T2[A, B]{A: a, B: b}
}

func (t *T2[A, B]) Unpack() (A, B)  {
    return t.A, t.B
}

type T3[A any, B any, C any] struct{
    A A
    B B
    C C
}

func ToT3[A any, B any, C any](a A, b B, c C) T3[A, B, C] {
    return T3[A, B, C]{A: a, B: b, C: c}
}

func (t *T3[A, B, C]) Unpack() (A, B, C)  {
    return t.A, t.B, t.C
}

type T4[A any, B any, C any, D any] struct{
    A A
    B B
    C C
    D D
}

func ToT4[A any, B any, C any, D any](a A, b B, c C, d D) T4[A, B, C, D] {
    return T4[A, B, C, D]{A: a, B: b, C: c, D: d}
}

func (t *T4[A, B, C, D]) Unpack() (A, B, C, D)  {
    return t.A, t.B, t.C, t.D
}
