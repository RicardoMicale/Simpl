package object

type ObjectType string

type Object interface {
	Type() ObjectType
	Inpect() string
}

const (
	INTEGER_OBJECT = "INTEGER"
	BOOLEAN_OBJECT = "BOOLEAN"
	NULL_OBJECT = "NULL"
)

type Null struct {}

func (n *Null) Inspect() string { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJECT }
