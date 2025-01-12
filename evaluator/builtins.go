package evaluator

import (
	"fmt"
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
	"removeAt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Wrong number of arguments. Got %d, expected 2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJECT {
				return newError(
					"First argument to `removeAt should be an array type. Got %s instead.",
					args[0].Type(),
				)
			}

			if args[1].Type() != object.INTEGER_OBJECT {
				return newError(
					"Second argument to `removeAt` should be an Integer/ Got %s instead",
					args[1].Type(),
				)
			}
			//	gets both the aray and the index
			arrayToModify := args[0].(*object.Array)
			indexToRemove := args[1].(*object.Integer).Value
			//	creates a slice copy of the elements of the array to modify
			updatedArray := arrayToModify.Elements

			//	checks if the index passed is bigger than the amount of elements of the array
			if int(indexToRemove) >= len(updatedArray) || int(indexToRemove) < 0 {
				return newError(
					"Index out of range. Received %d. should be between 0 and %d",
					indexToRemove,
					len(updatedArray) - 1,
				)
			}

			//	concatenates the array to modify, skipping the position to remove
			updatedArray = append(
				updatedArray[:indexToRemove], updatedArray[indexToRemove + 1:]...,
			)
			//	updates the elements on the original array to be modified
			arrayToModify.Elements = updatedArray
			//	returns the updated array
			return arrayToModify
		},
	},
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
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
	"range": {
		Fn: func(args ...object.Object) object.Object {
			//	check that the correct number of arguments is received
			if len(args) > 2 || len(args) == 0 {
				return newError("Wrong number of arguments. Expected 1 or 2, got %d", len(args))
			}
			//	type check the first argument to be an Integer
			_, ok := args[0].(*object.Integer)

			if !ok {
				return newError(
					"Expected the first argument to be an Integer value type. GOt %T instead",
					args[0],
				)
			}
			//	if there is more than one argument check the second element of the args array
			if len(args) == 2 {
				//	type check the second argument to be an Integer
				_, ok := args[1].(*object.Integer)

				if !ok {
					return newError(
						"Expected the second argument to be an Integer value type. Got %T instead",
						args[1],
					)
				}
			}
			//	declare the variables to navigate the range
			//	if the start is not provided, default it at zero
			var start int64 = 0
			var end int64
			/**
			*
			*	If there are 2 arguments
			*	the first one is the start of the range
			*	and the second one is the end of the range
			* if there is only one argument, define it as the end of the range
			*
			*/
			if len(args) == 2 {
				//	redefines start to the passed argument
				start = args[0].(*object.Integer).Value
				end = args[1].(*object.Integer).Value
			} else {
				end = args[0].(*object.Integer).Value
			}
			//	makes an array of size end - start + 1
			rangeArray := make([]object.Object, end - start + 1)
			//	loops through that array and assigns each position the start + the index they are at
			for i := range rangeArray {
				rangeArray[i] = &object.Integer{ Value: start + int64(i) }
			}
			//	returns the array object
			return &object.Array{ Elements: rangeArray }
		},
	},
}
