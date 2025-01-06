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
			case *object.Array:
				return &object.Integer{ Value: int64(len(arg.Elements)) }
			case *object.String:
				return &object.Integer{ Value: int64(len(arg.Value)) }
			default:
				return newError("Argument to `length` not supported, got %s", args[0].Type())
			}
		},
	},
	"firstElement": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Expected 1, got %d", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJECT {
				return newError("Argument to `firstElement` must be an Array, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"lastElement": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJECT {
				return newError("Argument to `lastElement` must be an Array, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[len(arr.Elements) - 1]
			}

			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Wrong number of arguments. Got %d, expected 2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJECT {
				return newError(
					"first Argument to `push` should be an Array type. Got %s",
					args[0].Type(),
				)
			}
			//	gets the array passe into the function
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			//	makes a new array with one extra slot for a  new element
			newElements := make([]object.Object, length + 1, length + 1)
			//	copies the contents of the elements of the original array into the newElements array
			copy(newElements, arr.Elements)
			//	adds the new element into the newElements array
			newElements[length] = args[1]
			//	reassigns the newElements array to the Elements attribute of the original array
			arr.Elements = newElements
			//	returns the same array
			return arr
		},
	},
	"removeLast": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJECT {
				return newError(
					"first Argument to `removeLast` should be an Array type. Got %s",
					args[0].Type(),
				)
			}

			//	gets the array passe into the function
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			//	makes a slice of the original array without the last element
			newElements := arr.Elements[:length - 1]
			//	reassigns the newElements array to the Elements attribute of the original array
			arr.Elements = newElements
			//	returns the same array
			return arr
		},
	},
	"removeAt": {},
	"copy": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJECT {
				return newError(
					"first Argument to `removeLast` should be an Array type. Got %s",
					args[0].Type(),
				)
			}
			//	gets the array from the arguments
			arr := args[0].(*object.Array)
			//	returns a new array with the same elements
			return &object.Array{ Elements: arr.Elements }
		},
	},
}
