package lexer

import (
	"testing"

	t "github.com/marcel/monkey/token"
	"github.com/stretchr/testify/suite"
)

type LexerTestSuite struct {
	suite.Suite
}

func TestLexerTestSuite(t *testing.T) {
	suite.Run(t, new(LexerTestSuite))
}

func (s *LexerTestSuite) TestNextToken() {
	input := `
		let five = 5;
		let ten = 10;

		let add = fn(x,y) {
			x + y
		};

		let result = add(five, ten);

		!-/*5;

		5 < 10 > 5;

		if (5 < 10) {
			return true;
		} else {
			return false;
		}

		10 == 10;

		10 != 9;
	`

	expectations := []struct {
		Type    t.Type
		Literal string
	}{
		// let five = 5;
		{t.LET, "let"}, {t.IDENT, "five"}, {t.ASSIGN, "="}, {t.INT, "5"}, {t.SEMICOLON, ";"},
		// let ten = 10;
		{t.LET, "let"}, {t.IDENT, "ten"}, {t.ASSIGN, "="}, {t.INT, "10"}, {t.SEMICOLON, ";"},
		// let add = fn(x,y) {
		{t.LET, "let"}, {t.IDENT, "add"}, {t.ASSIGN, "="}, {t.FUNCTION, "fn"}, {t.LPAREN, "("}, {t.IDENT, "x"}, {t.COMMA, ","}, {t.IDENT, "y"}, {t.RPAREN, ")"}, {t.LBRACE, "{"},
		//   x + y
		{t.IDENT, "x"}, {t.PLUS, "+"}, {t.IDENT, "y"},
		// };
		{t.RBRACE, "}"}, {t.SEMICOLON, ";"},
		// let result = add(five,ten);
		{t.LET, "let"}, {t.IDENT, "result"}, {t.ASSIGN, "="}, {t.IDENT, "add"}, {t.LPAREN, "("}, {t.IDENT, "five"}, {t.COMMA, ","}, {t.IDENT, "ten"}, {t.RPAREN, ")"}, {t.SEMICOLON, ";"},
		// !-/*5;
		{t.BANG, "!"}, {t.MINUS, "-"}, {t.SLASH, "/"}, {t.ASTERISK, "*"}, {t.INT, "5"}, {t.SEMICOLON, ";"},
		// 5 < 10 > 5;
		{t.INT, "5"}, {t.LT, "<"}, {t.INT, "10"}, {t.GT, ">"}, {t.INT, "5"}, {t.SEMICOLON, ";"},
		// if (5 < 10) {
		{t.IF, "if"}, {t.LPAREN, "("}, {t.INT, "5"}, {t.LT, "<"}, {t.INT, "10"}, {t.RPAREN, ")"}, {t.LBRACE, "{"},
		//   return true;
		{t.RETURN, "return"}, {t.TRUE, "true"}, {t.SEMICOLON, ";"},
		// } else {
		{t.RBRACE, "}"}, {t.ELSE, "else"}, {t.LBRACE, "{"},
		//  return false;
		{t.RETURN, "return"}, {t.FALSE, "false"}, {t.SEMICOLON, ";"},
		// }
		{t.RBRACE, "}"},
		// 10 == 10;
		{t.INT, "10"}, {t.EQ, "=="}, {t.INT, "10"}, {t.SEMICOLON, ";"},
		// 10 != 9;
		{t.INT, "10"}, {t.NOT_EQ, "!="}, {t.INT, "9"}, {t.SEMICOLON, ";"},
		// EOF
		{t.EOF, "\x00"},
	}

	l := New(input)

	for _, e := range expectations {
		tok := l.NextToken()

		s.Equal(e.Type, tok.Type)
		s.Equal(e.Literal, tok.Literal)
	}
}
