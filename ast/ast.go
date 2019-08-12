package ast

import (
	"bytes"
	"strings"

	"github.com/marcel/monkey/token"
)

type (
	Node interface {
		TokenLiteral() string
		String() string
	}

	Statement interface {
		Node
		statementNode()
	}

	Expression interface {
		Node
		expressionNode()
	}

	Program struct {
		Statements []Statement
	}

	Identifier struct {
		token.Token
		Value string
	}

	IntegerLiteral struct {
		token.Token
		Value int64
	}

	FunctionLiteral struct {
		token.Token
		Parameters []*Identifier
		Body       *BlockStatement
	}

	Boolean struct {
		token.Token
		Value bool
	}

	ExpressionStatement struct {
		token.Token
		Expression Expression
	}

	LetStatement struct {
		token.Token
		Name  *Identifier
		Value Expression
	}

	BlockStatement struct {
		token.Token
		Statements []Statement
	}

	PrefixExpression struct {
		token.Token
		Operator string
		Right    Expression
	}

	InfixExpression struct {
		token.Token
		Left     Expression
		Operator string
		Right    Expression
	}

	IfExpression struct {
		token.Token
		Condition   Expression
		Consequence *BlockStatement
		Alternative *BlockStatement
	}

	CallExpression struct {
		token.Token
		Function  Expression
		Arguments []Expression
	}

	ReturnStatement struct {
		token.Token
		ReturnValue Expression
	}
)

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (*Identifier) expressionNode()  {}
func (i *Identifier) String() string { return i.Value }

func (*IntegerLiteral) expressionNode()   {}
func (il *IntegerLiteral) String() string { return il.Literal }

func (*FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

func (*Boolean) expressionNode()  {}
func (b *Boolean) String() string { return b.Literal }

func (*ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

func (*BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (*LetStatement) statementNode() {}
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (*PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

func (*InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

func (*IfExpression) expressionNode() {}
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

func (*CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (*ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
