package rulehandler

import (
	"reflect"
	"testing"
)

func TestScanner(t *testing.T) {
	input := "(and (has \"192.168.1.1\" (in src)) (contains \"8.8.8.8\" (in dst)))"
	s := NewScanner("Test Scanner", input)
	var expected = []Token{
		{LeftParen, "("},
		{Keyword, "and"},
		{LeftParen, "("},
		{Keyword, "has"},
		{Quote, "\""},
		{Parameter, "192.168.1.1"},
		{Quote, "\""},
		{LeftParen, "("},
		{Keyword, "in"},
		{Parameter, "src"},
		{RightParen, ")"},
		{RightParen, ")"},
		{LeftParen, "("},
		{Keyword, "contains"},
		{Quote, "\""},
		{Parameter, "8.8.8.8"},
		{Quote, "\""},
		{LeftParen, "("},
		{Keyword, "in"},
		{Parameter, "dst"},
		{RightParen, ")"},
		{RightParen, ")"},
		{RightParen, ")"},
	}
	pos := 0
	for {
		tok := s.Next()
		if tok.Type == EOF {
			if pos != len(expected)-1 {
				t.Fatalf("got EOF before end of test")
			}
			return
		}
		if !reflect.DeepEqual(tok, expected[pos]) {
			t.Fatalf("\nexpected: %v\ngot: %v\npos: %d", expected[pos], tok, pos)
		}
		pos++
	}
}

func TestScannerOther(t *testing.T) {
	input := "(and (has \"192.168.1.1\" in src) (contains \"8.8.8.8\" in dst)))"
	s := NewScanner("Test Scanner", input)
	var expected = []Token{
		{LeftParen, "("},
		{Keyword, "and"},
		{LeftParen, "("},
		{Keyword, "has"},
		{Quote, "\""},
		{Parameter, "192.168.1.1"},
		{Quote, "\""},
		{Parameter, "in"},
		{Parameter, "src"},
		{RightParen, ")"},
		{LeftParen, "("},
		{Keyword, "contains"},
		{Quote, "\""},
		{Parameter, "8.8.8.8"},
		{Quote, "\""},
		{Parameter, "in"},
		{Parameter, "dst"},
		{RightParen, ")"},
		{RightParen, ")"},
		{RightParen, ")"},
	}
	pos := 0
	for {
		tok := s.Next()
		if tok.Type == EOF {
			if pos != len(expected)-1 {
				t.Fatalf("got EOF before end of test")
			}
			return
		}
		if !reflect.DeepEqual(tok, expected[pos]) {
			t.Fatalf("\nexpected: %v\ngot: %v\npos: %d", expected[pos], tok, pos)
		}
		pos++
	}
}
