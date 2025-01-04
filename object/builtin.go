package object

type BuiltInFunction func(args ...Object) Object

type BuiltIn struct {
	Fn BuiltInFunction
}

func (bi *BuiltIn) Type() ObjectType { return BUILTIN_OBJECT }
func (bi *BuiltIn) Inspect() string { return "builtin function" }
