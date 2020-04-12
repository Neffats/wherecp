package rulehandler

import (
	"fmt"
	"regexp"

	"github.com/Neffats/wherecp/core"
)

const (
	
)

var (
	hostPattern = regexp.MustCompile("([1-2]?[0-9]?[0-9]\.){3}([1-2]?[0-9]?[0-9])")
	networkPattern = regexp.MustCompile("([1-2]?[0-9]?[0-9]\.){3}([1-2]?[0-9]?[0-9])\/([1-3]?[0-9])")
	rangePattern = regexp.MustCompile("([1-2]?[0-9]?[0-9]\.){3}([1-2]?[0-9]?[0-9])\-([1-2]?[0-9]?[0-9]\.){3}([1-2]?[0-9]?[0-9])")
	servicePattern = regexp.MustCompile("\w*\/\d*")
)

type Parser struct {
	input string
	s Scanner
}

func Parse(input string) (filterFn, error) {
	s := NewScanner("Filter Scanner", input)
	p := NewParser(input, s)

	for tok := s.Next(); tok != EOF {
		switch tok.Type {
		case LeftParen:
			filter, err := p.parseKeyword()
			if err != nil {
				return nil, fmt.Errorf("failed to parse keyword: %v", err)
			}
		}
	}
}

type object struct {
	h *core.Host
	n *core.Network
	r *core.Range
	p *core.Port
	g *core.Group
}

func (p *Parser) parseKeyword() (filterFn, error) {
	for tok := p.s.Next(); tok.Type != RightParen {
		if tok.Type != Keyword {
			return nil, fmt.Errorf("expected a keyword but got: %s", tok.Value)
		}
		keyword := tok.Value
		tok = p.s.Next()

		obj := object{}
		var err error
		
		switch tok.Type {
		case Quote:
			tok = p.s.Next()
			if tok.Type != Parameter {
				return nil, fmt.Errorf("expected a parameter but got: %s", tok.Value)
			}
			switch {
			case hostPattern.MatchString(tok.Value):
				obj.h, err = p.parseHost(tok.Value)
				if err != nil {
					return nil, fmt.Errorf("failed to parse host: %v", err)
				}
			case networkPattern.MatchString(tok.Value):
				obj.n, err = p.parseNetwork(tok.Value)
				if err != nil {
					return nil, fmt.Errorf("failed to parse network: %v", err)
				}
			case rangePattern.MatchString(tok.Value):
				obj.r, err = p.parseRange(tok.Value)
				if err != nil {
					return nil, fmt.Errorf("failed to parse range: %v", err)
				}
			case servicePattern.MatchString(tok.Value):
				obj.p, err = p.parseService(tok.Value)
				if err != nil {
					return nil, fmt.Errorf("failed to parse service: %v", err)
				}
			default:
				obj.g, err = p.parseGroup(tok.Value)
				if err != nil {
					return nil, fmt.Errorf("failed to parse group: %v", err)
				}
			}
		}
	}
}

func (p *Parser) parseHost(token string) (*core.Host, error) {
	return core.NewHost("filter host", token, "")
}

func (p *Parser) parseNetwork(token string) (*core.Network, error) {
	args := strings.Split(token)
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid network string: %s", token)
	}
	mask, ok := masks[args[1]]
	if !ok {
		return nil, fmt.Errorf("invalid subnet mask: %s", token)
	}
	return core.NewNetwork("filter network", args[0], mask, "")
}
