package rulehandler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Neffats/wherecp/core"
)

const (
	
)

var (
	hostPattern = regexp.MustCompile("([1-2]?[0-9]?[0-9]\\.){3}([1-2]?[0-9]?[0-9])")
	networkPattern = regexp.MustCompile("([1-2]?[0-9]?[0-9]\\.){3}([1-2]?[0-9]?[0-9])\\/([1-3]?[0-9])")
	rangePattern = regexp.MustCompile("([1-2]?[0-9]?[0-9]\\.){3}([1-2]?[0-9]?[0-9])\\-([1-2]?[0-9]?[0-9]\\.){3}([1-2]?[0-9]?[0-9])")
	servicePattern = regexp.MustCompile("\\w*\\/\\d*")
)

type constructer interface {
	construct() filterFn
}

type Parser struct {
	s *Scanner
}

type boolOp struct {
	fn func(...filterFn) filterFn
	args []constructer
}

func (b *boolOp) construct() filterFn {
	constructedArgs := make([]filterFn, 0)
	for _, a := range b.args {
		constructedArgs = append(constructedArgs, a.construct())
	}
	return b.fn(constructedArgs...)
}

type hasOp struct {
	fn func(interface{}, func(*core.Rule) core.Haser) filterFn
	objArg interface{}
	compArg func() func(*core.Rule) core.Haser
}

func (h *hasOp) construct() filterFn {
	return h.fn(h.objArg, h.compArg())
}

func NewParser(s *Scanner) *Parser {
	return &Parser{s: s}
}

func Parse(input string) (filterFn, error) {
	s := NewScanner("Filter Scanner", input)
	p := NewParser(s)

	var filter filterFn
	var err error

	for tok := s.Next(); tok.Type != EOF; {
		switch tok.Type {
		case LeftParen:
			filter, err = p.parseKeyword()
			if err != nil {
				return nil, fmt.Errorf("failed to parse keyword: %v", err)
			}
		}
	}
	return filter, nil
}

type object struct {
	h *core.Host
	n *core.Network
	r *core.Range
	p *core.Port
	g *core.Group
}

func (p *Parser) parseKeyword() (filterFn, error) {
	var filter filterFn
	for tok := p.s.Next(); tok.Type != RightParen; {
		if tok.Type != Keyword {
			return nil, fmt.Errorf("expected a keyword but got: %s", tok.Value)
		}
		//keyword := tok.Value
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
	return filter, nil
}

func (p *Parser) parseHost(token string) (*core.Host, error) {
	return core.NewHost("filter host", token, "")
}

func (p *Parser) parseNetwork(token string) (*core.Network, error) {
	args := strings.Split(token, "/")
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid network string: %s", token)
	}
	mask, ok := masks[args[1]]
	if !ok {
		return nil, fmt.Errorf("invalid subnet mask: %s", token)
	}
	return core.NewNetwork("filter network", args[0], mask, "")
}

func (p *Parser) parseRange(token string) (*core.Range, error) {
	// Expected format: 192.168.1.1-192.168.1.5
	args := strings.Split(token, "-")
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid range string: %s", token)
	}
	return core.NewRange("filter range", args[0], args[1], "")
}

func (p *Parser) parseService(token string) (*core.Port, error) {
	// Expected format: tcp/443
	args := strings.Split(token, "/")
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid service string: %s", token)
	}
	portNo, err := strconv.ParseUint(args[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to convert port number to int: %v", err)
	}
	return core.NewPort("filter port", uint(portNo), args[0], "")
}

func (p *Parser) parseGroup(token string) (*core.Group, error) {
	return core.NewGroup(token, ""), nil
}
