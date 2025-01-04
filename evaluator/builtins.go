package evaluator

import (
	"language/object"
)

var builtins = map[string]*object.BuiltIn{
	"length": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments, expected 1, got %d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{ Value: int64(len(arg.Value)) }
			default:
				return newError("Argument to `length` not supported, got %s", args[0].Type())
			}
		},
	},
}
