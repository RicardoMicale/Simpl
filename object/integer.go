package object

import (
	"fmt"
)

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJECT }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (i *Integer) MapKey() MapKey {
	return MapKey{ Type: i.Type(), Value: uint64(i.Value) }
}
