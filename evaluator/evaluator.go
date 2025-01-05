package evaluator

import (
	"fmt"
	"language/ast"
	"language/object"
)

var (
	NULL = &object.Null{}
	TRUE = &object.Boolean{ Value: true}
	FALSE = &object.Boolean{ Value: false }
)

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJECT || rt == object.ERROR_OBJECT {
				return result
			}
		}
	}

	return result
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJECT
	}

	return false
}

func evalStatements(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	//	evaluates the statements recursively
	for _, statement := range statements {
		result = Eval(statement, env)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func nativeBoolToBooleaObject(input bool) *object.Boolean {
	if input { return TRUE }
	return FALSE
}

func evalNotOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	//	check if the object passed is an integer
	if right.Type() != object.INTEGER_OBJECT {
		return newError("Unknown operator: -%s", right.Type())
	}
	//	retrieve the value passed on the right of the minus operand
	value := right.(*object.Integer).Value
	//	return a new object integer with the negative value
	return &object.Integer{ Value: -value }
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalNotOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("Unknown operator: %s%s", operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{ Value: leftValue + rightValue }
	case "-":
		return &object.Integer{ Value: leftValue - rightValue }
	case "*":
		return &object.Integer{ Value: leftValue * rightValue }
	case "/":
		return &object.Integer{ Value: leftValue / rightValue }
	case "==":
		return nativeBoolToBooleaObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleaObject(leftValue != rightValue)
	case ">":
		return nativeBoolToBooleaObject(leftValue > rightValue)
	case "<":
		return nativeBoolToBooleaObject(leftValue < rightValue)
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	//	gets both values
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case "==":
		return nativeBoolToBooleaObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleaObject(leftVal != rightVal)
	case "+":
		//	creates the object with the concatenated strings
		return &object.String{ Value: leftVal + rightVal }
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJECT && right.Type() == object.STRING_OBJECT:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleaObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleaObject(left != right)
	case left.Type() != right.Type():
		return newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	//	checks if the condition has an error
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{ Message: fmt.Sprintf(format, a...) }
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {

	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtIn, ok := builtins[node.Value]; ok {
		return builtIn
	}

	return newError("%s", "Identifier not found: " + node.Value)
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range expressions {
		evaluated := Eval(e, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch function := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(function, args)
		evaluated := Eval(function.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.BuiltIn:
		return function.Fn(args...)
	default:
		return newError("Not a function: %s", function.Type())
	}
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	//	Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		//	if an error is found on the val variable, return val, which contains an error
		if isError(val) {
			return val
		}
		return &object.ReturnValue{ Value: val }
	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return Eval(node.Name, env)
	case *ast.ConstStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return Eval(node.Value, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	//	Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{ Value: node.Value }
	case *ast.Boolean:
		return nativeBoolToBooleaObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.FunctionLiteral:
		parameters := node.Parameters
		body := node.Body
		return &object.Function{ Parameters: parameters, Body: body, Env: env }
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.StringLiteral:
		return &object.String{ Value: node.Value }
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{ Elements: elements }
	}

	return nil
}
