package object

import (
	"bytes"
	"fmt"
	"strings"
)

type MapPair struct {
	Key Object
	Value Object
}

type Map struct {
	Pairs map[MapKey]MapPair
}

func (m *Map) Type() ObjectType { return MAP_OBJECT }
func (m *Map) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
