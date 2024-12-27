package ast

import (
	"language/token"
	"testing"
)

func TestString(t * testing.T) {
	program := &Program{
		Statements: []Statement{
			&ConstStatement{
				Token: token.Token{ Type: token.CONST, Literal: "const" },
				Name: &Identifier{
					Token: token.Token{ Type: token.IDENTIFIER, Literal: "myConst" },
					Value: "myConst",
				},
				Value: &Identifier{
					Token: token.Token{ Type: token.IDENTIFIER, Literal: "anotherConst" },
					Value: "anotherConst",
				},
			},
		},
	}

	if program.String()!= "const myConst = anotherConst;" {
		t.Errorf("program.String() got wrong. Got=%q", program.String())
	}
}
