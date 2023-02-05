package iter

import (
    "strings"
)

// Join uses a [Joiner] to build a result by walking over an iterator.
func Join[In any, Out any](j Joiner[In, Out], it It[In]) Out {
    WalkFinal(j.Join, it)
    return j.End()
}

// Joiner is a type for something that can build a result by walking over
// an iterator using [Join]. It must be possible for a Joiner to be used
// multiple times (although not concurrently) from multiple calls to [Join].
type Joiner[In any, Out any] interface {
    Join(x In, isFinal bool)
    End() Out
}

// StringJoiner returns a [Joiner] for concatenating strings with a (possibly
// empty) separator.
func StringJoiner(sep string) Joiner[string, string] {
    return &stringJoiner{sep: sep}
}

type stringJoiner struct {
    sb strings.Builder
    sep string
}

func (j *stringJoiner) Join(x string, isFinal bool) {
    j.sb.WriteString(x)
    if !isFinal { j.sb.WriteString(j.sep) }
}

func (j *stringJoiner) End() string {
    result := j.sb.String()
    j.sb.Reset()
    return result
}
