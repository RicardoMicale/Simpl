package object

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJECT}
func (s *String) Inspect() string { return s.Value }
