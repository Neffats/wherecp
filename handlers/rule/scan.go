package rulehandler

import (
	"unicode"
	"unicode/utf8"
)

const (
	EOF = iota
	Error
	LeftParen
	RightParen
	Keyword
	Parameter
	Quote
	Newline
)

type Token struct {
	Type TokenType
	Value string
}

type TokenType int

type stateFn func(*Scanner) stateFn

type Scanner struct {
	name string
	tokens <-chan Token
	input string
	state stateFn
	
	// Used as a buffer for the input.
	start int
	pos int
	width int
}

func NewScanner(name, input string) *Scanner {
	return Scanner{
		name: name,
		tokens: make(chan Token, 2),
		input: input,
		state: lexAny,
		start: 0,
		pos: 0,
		width: 0,
	}
}

func (s *Scanner) Next() Token {
	for s.state != nil {
		select {
		case tok := <-s.tokens:
			return tok
		default:
			s.state = s.state(s)
		}
	}
	if s.tokens != nil {
		close(s.tokens)
		s.tokens = nil
	}

	return Token{EOF, "EOF"}
}

// Emit takes a TokenType and creates a new Token of that type and the
// contents of the current buffer, the new token is then put on the
// Tokens channel ready to be consumed. The buffer is then cleared and
// set to the current position in the input.
func (s *Scanner) emit(t TokenType) {
	s.Tokens <- Token{
		Type: t,
		Value: s.Input[s.Start: this.Pos],
	}
	// Reset the buffer.
	s.Start = s.Pos
}

func (s *Scanner) next() rune {
	if s.Pos >= len(s.Input) {
		s.width = 0
		return eof
	}
	rune, s.width = utf8.DecodeInString(l.input[l.pos:])
	s.pos += s.width
	return rune
}

func (s *Scanner) ignore() {
	s.start = s.pos
}

func (s *Scanner) backup() {
	s.pos -= s.width
}

func (s *Scanner) peek() int {
	rune := s.next()
	s.backup()
	return rune
}

func (s *Scanner) accept(valid string) bool {
	if strings.IndexRune(valid, s.next()) >= 0 {
		return true
	}
	s.backup()
	return false
}

func (s *Scanner) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	s.backup()
}

func lexAny(s *Scanner) stateFn {
	switch r := s.next(); {
	case r == eof:
		return nil
	case r == '\n':
		s.emit(Newline)
		return lexAny
	case r == "(":
		s.emit(LeftParen)
		return lexKeyword
	case r == ")":
		s.emit(RightParen)
		return lexAny
	case isSpace(r):
		s.skipWhitespace()
		return lexAny
	default:
		return lexParam
	}
}

func lexKeyword(s *Scanner) stateFn {
	for !isSpace(s.peek()) {
		s.next()
	}
	s.emit(Keyword)
	return lexParam
}

func lexParam(s *Scanner) stateFn {
	s.skipWhitespace()
	switch r := s.next() {
	case r == "(":
		s.emit(LeftParen)
		return lexParam
	case r == '"':
		s.emit(Quote)
		return lexInsideQuote
	case isAlphanumeric(r):
		for !isSpace(s.peek()) {
			s.next()
		}
		s.emit(Parameter)
		return lexAny
	default:
		s.errorf("")
	}
}

func lexInsideQuote(s *Scanner) stateFn {
	for s.peek() != '"' {
		s.next()
	}
	s.emit(Parameter)
	return lexClosingQuote
}

func lexClosingQuote(s *Scanner) stateFn {
	r := s.next()
	if r != '"' {
		return s.errorf("expected closing quote but got: %s", r)
	}
	s.emit(Quote)
	return lexAny
}
	
func (s *Scanner) skipWhitespace() {
	for isSpace(s.peek()) {
		s.next()
	}
	s.ignore()
}

func (s *Scanner) errorf(format string, args ...interface{}) stateFn {
	s.tokens <- Token{Error, fmt.Sprintf(format, args...)}
	return lexAny
}

func isAlphaNumeric(r rune) bool {
	return r == "_" || unicode.IsLetter(r) || unicode.IsDigit(r)
}
	
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}
