package evaluator

import (
	"language/ast"
	"language/object"
)

var (
	NULL = &object.Null{}
	TRUE = &object.Boolean{ Value: true}
	FALSE = &object.Boolean{ Value: false }
)

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object
	//	evaluates the statements recursively
	for _, statement := range statements {
		result = Eval(statement)
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

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalNotOperatorExpression(right)
	default:
		return NULL
	}
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	//	Statements
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	//	Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{ Value: node.Value }
	case *ast.Boolean:
		return nativeBoolToBooleaObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	}

	return nil
}
