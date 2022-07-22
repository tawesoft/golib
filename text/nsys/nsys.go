// Package nsys provides information about Unicode numbering systems
package nsys

type Type int

type System struct {

}


func (s System) IsDecimal() {}

// system can be a named numbering system e.g. "latn", or a type such as
// "default", "native", "traditiono", "finance".
func New(base String, region String, system string) System {
}

func NewFromTag(tag language.Tag) System {
}


