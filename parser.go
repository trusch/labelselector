package labelselector

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Parse parses a label selector expression
func Parse(input io.Reader) (LabelSelector, error) {
	parser := NewParser(input)
	return parser.Parse()
}

// ParseString parses a label selector expression from a string
func ParseString(str string) (LabelSelector, error) {
	return Parse(strings.NewReader(str))
}

// NewParser creates a new parser instance
func NewParser(input io.Reader) *Parser {
	return &Parser{NewLexer(input)}
}

// Parser is capable of parsing a label selector expression
type Parser struct {
	lexer *Lexer
}

// scan gets the next token
func (p *Parser) scan() (Token, string) {
	tok, lit := p.lexer.Next()
	return tok, lit
}

// next returns the next non-whitespace token
func (p *Parser) next() (Token, string) {
	tok, lit := p.scan()
	for tok == WS {
		tok, lit = p.scan()
	}
	return tok, lit
}

// Parse actually parses the input and returns the resulting LabelSelector
func (p *Parser) Parse() (LabelSelector, error) {
	selector := LabelSelector{}

	// selectors are on the toplevel a comma separated list, so lets iterate over those items
	for {
		// discard the eventually leading comma
		tok, lit := p.next()
		if tok == COMMA {
			continue
		}
		// we are ready if we see EOF
		if tok == EOF {
			break
		}
		if tok == ILLEGAL {
			return selector, errors.New("illegal token")
		}
		// there are two cases now:
		// * we see a '!' -> this will be a not-exist requirement
		// * we see an identifier -> this will be a one of the other requirements
		switch tok {
		case EXCLAMATION_MARK:
			// its a not-exist requirement
			req, err := p.parseNotExistsRequirement()
			if err != nil {
				return selector, err
			}
			selector.Requirements = append(selector.Requirements, req)
		case IDENT:
			// we have a identifier so its one of
			// * equal requirement
			// * exists requirement
			// * in requirement
			// * not-in requirement
			var (
				key      = lit
				tok, lit = p.next()
				err      error
				req      Requirement
			)
			switch tok {
			case EQUAL:
				// its a equal requirement
				req, err = p.parseEqualRequirement(key)
			case NOT_EQUAL:
				// its a not-equal requirement
				req, err = p.parseNotEqualRequirement(key)
			case IN:
				// its a in requirement
				req, err = p.parseInRequirement(key)
			case LOWER_THAN:
				// its a lower than requirement
				req, err = p.parseLowerThanRequirement(key)
			case LOWER_THAN_EQUAL:
				// its a lower than equal requirement
				req, err = p.parseLowerThanEqualRequirement(key)
			case GREATER_THAN:
				// its a greater than requirement
				req, err = p.parseGreaterThanRequirement(key)
			case GREATER_THAN_EQUAL:
				// its a greater than equal requirement
				req, err = p.parseGreaterThanEqualRequirement(key)
			case NOT:
				// its a not-in requirement
				req, err = p.parseNotInRequirement(key)
			case COMMA, EOF:
				// its a exists requirement
				req = Requirement{
					Key:       key,
					Operation: OperationExists,
				}
			default:
				return selector, fmt.Errorf("unexpected token '%s'", lit)
			}
			if err != nil {
				return selector, err
			}
			selector.Requirements = append(selector.Requirements, req)
		}
	}
	return selector, nil
}

func (p *Parser) parseNotEqualRequirement(key string) (req Requirement, err error) {
	tok, lit := p.next()
	if tok != IDENT {
		return req, errors.New("expect identifier after not-equal operator")
	}
	req = Requirement{
		Key:       key,
		Operation: OperationNotEquals,
		Value:     lit,
	}
	return req, nil
}

func (p *Parser) parseEqualRequirement(key string) (req Requirement, err error) {
	tok, lit := p.next()
	if tok != IDENT {
		return req, errors.New("expect identifier after equal operator")
	}
	req = Requirement{
		Key:       key,
		Operation: OperationEquals,
		Value:     lit,
	}
	return req, nil
}

func (p *Parser) parseLowerThanRequirement(key string) (req Requirement, err error) {
	tok, lit := p.next()
	if tok != IDENT {
		return req, errors.New("expect identifier after < operator")
	}
	req = Requirement{
		Key:       key,
		Operation: OperationLowerThan,
		Value:     lit,
	}
	return req, nil
}

func (p *Parser) parseLowerThanEqualRequirement(key string) (req Requirement, err error) {
	tok, lit := p.next()
	if tok != IDENT {
		return req, errors.New("expect identifier after <= operator")
	}
	req = Requirement{
		Key:       key,
		Operation: OperationLowerThanEqual,
		Value:     lit,
	}
	return req, nil
}

func (p *Parser) parseGreaterThanRequirement(key string) (req Requirement, err error) {
	tok, lit := p.next()
	if tok != IDENT {
		return req, errors.New("expect identifier after > operator")
	}
	req = Requirement{
		Key:       key,
		Operation: OperationGreaterThan,
		Value:     lit,
	}
	return req, nil
}

func (p *Parser) parseGreaterThanEqualRequirement(key string) (req Requirement, err error) {
	tok, lit := p.next()
	if tok != IDENT {
		return req, errors.New("expect identifier after >= operator")
	}
	req = Requirement{
		Key:       key,
		Operation: OperationGreaterThanEqual,
		Value:     lit,
	}
	return req, nil
}

func (p *Parser) parseIdentList() (list []string, err error) {
	for {
		tok, lit := p.next()
		if tok == CLOSING_BRACKET {
			break
		} else if tok == COMMA {
			continue
		} else if tok == IDENT {
			list = append(list, lit)
		} else {
			return nil, fmt.Errorf("unexpected token in value list (%s)", lit)
		}
	}
	return list, nil
}

func (p *Parser) parseNotExistsRequirement() (req Requirement, err error) {
	tok, lit := p.next()
	if tok != IDENT {
		return req, errors.New("expect identifier after exclamation mark")
	}
	return Requirement{
		Key:       lit,
		Operation: OperationNotExists,
	}, nil
}

func (p *Parser) parseInRequirement(key string) (req Requirement, err error) {
	req = Requirement{
		Key:       key,
		Operation: OperationIn,
	}
	tok, _ := p.next()
	if tok != OPENING_BRACKET {
		return req, errors.New("expect opening bracket after in operator")
	}
	list, err := p.parseIdentList()
	if err != nil {
		return req, err
	}
	req.Values = list
	return
}

func (p *Parser) parseNotInRequirement(key string) (req Requirement, err error) {
	tok, lit := p.next()
	if tok != IN {
		return req, fmt.Errorf("require 'IN' after 'NOT' got '%s'", lit)
	}
	req, err = p.parseInRequirement(key)
	if err != nil {
		return req, err
	}
	req.Operation = OperationNotIn
	return req, nil
}
