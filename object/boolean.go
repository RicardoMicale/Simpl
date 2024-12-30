package object

import (
	"fmt"
)

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJECT }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
