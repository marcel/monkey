package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/marcel/monkey/ast"
	"github.com/marcel/monkey/lexer"
	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
}

func TestParserTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}

func (s *ParserTestSuite) TestOperatorPrecedenceParsing() {
	expectations := []struct {
		Input    string
		Expected string
	}{
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, e := range expectations {
		p := New(lexer.New(e.Input))
		program := p.ParseProgram()
		s.NotNil(program)
		s.checkParserErrors(p)

		s.Equal(e.Expected, program.String())
	}
}

func (s *ParserTestSuite) TestIfExpression() {
	input := `if (x < y) { x }`

	p := New(lexer.New(input))
	program := p.ParseProgram()
	s.NotNil(program)
	s.checkParserErrors(p)

	s.Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.True(ok)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	s.True(ok)

	s.testInfixExpression(exp.Condition, "x", "<", "y")

	s.Len(exp.Consequence.Statements, 1)
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	s.True(ok)

	s.testIdentifier(consequence.Expression, "x")

	s.Nil(exp.Alternative)
}

func (s *ParserTestSuite) TestIfElseExpression() {
	input := `if (x < y) { x } else { y }`

	p := New(lexer.New(input))
	program := p.ParseProgram()
	s.NotNil(program)
	s.checkParserErrors(p)

	s.Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.True(ok)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	s.True(ok)

	s.testInfixExpression(exp.Condition, "x", "<", "y")

	s.Len(exp.Consequence.Statements, 1)
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	s.True(ok)

	s.testIdentifier(consequence.Expression, "x")

	s.NotNil(exp.Alternative)

	s.Len(exp.Alternative.Statements, 1)
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	s.True(ok)

	s.testIdentifier(alternative.Expression, "y")
}

func (s *ParserTestSuite) TestFunctionLiteralParsing() {
	input := `fn(x, y) { x + y; }`

	p := New(lexer.New(input))
	program := p.ParseProgram()
	s.NotNil(program)
	s.checkParserErrors(p)

	s.Len(program.Statements, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.True(ok)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	s.True(ok)

	s.Len(function.Parameters, 2)
	s.testLiteralExpression(function.Parameters[0], "x")
	s.testLiteralExpression(function.Parameters[1], "y")

	s.Len(function.Body.Statements, 1)
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	s.True(ok)
	s.testInfixExpression(bodyStmt.Expression, "x", "+", "y")
}

func (s *ParserTestSuite) TestFunctionParameterParsing() {
	expectations := []struct {
		Input  string
		Params []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, e := range expectations {
		p := New(lexer.New(e.Input))
		program := p.ParseProgram()
		s.NotNil(program)
		s.checkParserErrors(p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		s.True(ok)
		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		s.True(ok)

		s.Len(function.Parameters, len(e.Params))

		for i, ident := range e.Params {
			s.testLiteralExpression(function.Parameters[i], ident)
		}
	}
}

func (s *ParserTestSuite) TestCallExpressionParsing() {
	input := "add(1, 2 * 3, 4 + 5);"

	p := New(lexer.New(input))
	program := p.ParseProgram()
	s.NotNil(program)
	s.checkParserErrors(p)

	s.Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.True(ok)

	exp, ok := stmt.Expression.(*ast.CallExpression)
	s.True(ok)

	s.testIdentifier(exp.Function, "add")

	s.Len(exp.Arguments, 3)
	s.testLiteralExpression(exp.Arguments[0], 1)
	s.testInfixExpression(exp.Arguments[1], 2, "*", 3)
	s.testInfixExpression(exp.Arguments[2], 4, "+", 5)
}

func (s *ParserTestSuite) TestInfixExpressionsParsing() {
	expectations := []struct {
		Input    string
		LHS      interface{}
		Operator string
		RHS      interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, e := range expectations {
		p := New(lexer.New(e.Input))
		program := p.ParseProgram()
		s.NotNil(program)
		s.checkParserErrors(p)
		s.Len(program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		s.True(ok)

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		s.True(ok)

		s.testInfixExpression(exp, e.LHS, e.Operator, e.RHS)
	}
}

func (s *ParserTestSuite) TestParsingPrefixExpressions() {
	expectations := []struct {
		Input    string
		Operator string
		Value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, e := range expectations {
		l := lexer.New(e.Input)
		p := New(l)
		program := p.ParseProgram()
		s.checkParserErrors(p)
		s.NotNil(program)
		s.Len(program.Statements, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		s.True(ok)
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		s.True(ok)
		s.Equal(e.Operator, exp.Operator)
		s.testLiteralExpression(exp.Right, e.Value)
	}
}
func (s *ParserTestSuite) TestIntegerLiteralExpression() {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	s.checkParserErrors(p)
	s.NotNil(program)
	s.Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.True(ok)
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	s.True(ok)
	s.Equal(int64(5), literal.Value)
	s.Equal("5", literal.TokenLiteral())
}

func (s *ParserTestSuite) TestIdentifierExpression() {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	s.checkParserErrors(p)
	s.NotNil(program)
	s.Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.True(ok)
	ident, ok := stmt.Expression.(*ast.Identifier)
	s.True(ok)
	s.Equal("foobar", ident.Value)
	s.Equal("foobar", ident.TokenLiteral())
}

func (s *ParserTestSuite) TestBooleanExpression() {
	input := "true;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	s.checkParserErrors(p)
	s.NotNil(program)
	s.Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.True(ok)
	b, ok := stmt.Expression.(*ast.Boolean)
	s.True(ok)
	s.True(b.Value)
	s.Equal("true", b.TokenLiteral())
}

func (s *ParserTestSuite) TestReturnStatements() {
	expectations := []struct {
		Input       string
		ReturnValue int64
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return 9933322;", 9933322},
	}

	for _, e := range expectations {
		p := New(lexer.New(e.Input))
		program := p.ParseProgram()
		s.checkParserErrors(p)
		s.NotNil(program)
		s.Len(program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		s.True(ok)
		lit, ok := stmt.ReturnValue.(*ast.IntegerLiteral)
		s.True(ok)
		s.testIntegerLiteral(lit, e.ReturnValue)
	}
}

func (s *ParserTestSuite) TestLetStatements() {
	expectations := []struct {
		Input      string
		Identifier string
		Value      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, e := range expectations {
		p := New(lexer.New(e.Input))
		program := p.ParseProgram()
		s.NotNil(program)
		s.checkParserErrors(p)

		s.Len(program.Statements, 1)

		stmt := program.Statements[0]
		s.testLetStatement(stmt, e.Identifier)

		val, ok := stmt.(*ast.LetStatement)
		s.True(ok)
		s.testLiteralExpression(val.Value, e.Value)
	}
}

func (s *ParserTestSuite) testIntegerLiteral(il ast.Expression, value int64) {
	integ, ok := il.(*ast.IntegerLiteral)
	s.True(ok)

	s.Equal(value, integ.Value)

	s.Equal(strconv.Itoa(int(value)), integ.TokenLiteral())
}

func (s *ParserTestSuite) testLetStatement(stmt ast.Statement, name string) {
	s.Equal("let", stmt.TokenLiteral())

	letStmt, ok := stmt.(*ast.LetStatement)
	s.True(ok)

	s.Equal(name, letStmt.Name.Value)
	s.Equal(name, letStmt.Name.TokenLiteral())
}

func (s *ParserTestSuite) checkParserErrors(p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	messages := []interface{}{}
	for _, msg := range errors {
		messages = append(messages, fmt.Sprintf("parser error: %q", msg))
	}

	s.FailNow("parsing failed with errors", messages...)
}

func (s *ParserTestSuite) testIdentifier(exp ast.Expression, value string) {
	ident, ok := exp.(*ast.Identifier)
	s.True(ok)

	s.Equal(value, ident.Value)
	s.Equal(value, ident.TokenLiteral())
}

func (s *ParserTestSuite) testLiteralExpression(exp ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case int:
		s.testIntegerLiteral(exp, int64(v))
	case int64:
		s.testIntegerLiteral(exp, v)
	case string:
		s.testIdentifier(exp, v)
	case bool:
		s.testBooleanLiteral(exp, v)
	default:
		s.FailNow("can not test unhandled expression", "got=%T for type=%T", exp, expected)
	}
}

func (s *ParserTestSuite) testBooleanLiteral(exp ast.Expression, value bool) {
	bo, ok := exp.(*ast.Boolean)
	s.True(ok)

	s.Equal(value, bo.Value)
	s.Equal(fmt.Sprintf("%t", value), bo.TokenLiteral())
}

func (s *ParserTestSuite) testInfixExpression(
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) {
	opExp, ok := exp.(*ast.InfixExpression)
	s.True(ok)

	s.testLiteralExpression(opExp.Left, left)
	s.Equal(operator, opExp.Operator)
	s.testLiteralExpression(opExp.Right, right)
}
