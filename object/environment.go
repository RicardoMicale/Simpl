package object

import "fmt"

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{ store: s, outer: nil }
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, value Object) Object {
	prevValue, ok := e.store[name]

	if ok {
		if prevValue.Type() != value.Type() {
			return &Error{ Message: fmt.Sprintf("Cannot reassign different types. Passed %s type to %s type variable", value.Type(), prevValue.Type()) }
		}
	}

	e.store[name] = value
	return value
}
