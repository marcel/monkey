package ast

import (
	"testing"

	"github.com/marcel/monkey/token"
	"github.com/stretchr/testify/suite"
)

type ASTTestSuite struct {
	suite.Suite
}

func TestASTTestSuite(t *testing.T) {
	suite.Run(t, new(ASTTestSuite))
}

func (s *ASTTestSuite) TestString() {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.LET.Token("let"),
				Name: &Identifier{
					Token: token.IDENT.Token("myVar"),
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.IDENT.Token("anotherVar"),
					Value: "anotherVar",
				},
			},
			&ReturnStatement{
				Token: token.RETURN.Token("return"),
				ReturnValue: &Identifier{
					Token: token.IDENT.Token("myVar"),
					Value: "myVar",
				},
			},
		},
	}

	s.Equal("let myVar = anotherVar;return myVar;", program.String())
}
