package labelselector

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// Token represents a lexical token.
type Token int

const (
	ILLEGAL Token = iota
	EOF
	WS
	IDENT            // identifier
	COMMA            // ,
	EXCLAMATION_MARK // !
	IN               // in
	NOT              // not
	EQUAL            // = or ==
	NOT_EQUAL        // !=
	OPENING_BRACKET  // (
	CLOSING_BRACKET  // )
)

var eof = rune(0)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isValidIdentRune(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		(ch == '_') ||
		(ch == '-') ||
		(ch == '.') ||
		(ch == '/')
}

type Lexer struct {
	r *bufio.Reader
}

// NewLexer returns a new instance of Lexer.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Lexer) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Lexer) unread() { _ = s.r.UnreadRune() }

// Next returns the next token and literal value.
func (s *Lexer) Next() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()
	switch {
	case isWhitespace(ch):
		s.unread()
		return s.scanWhitespace()
	case isValidIdentRune(ch):
		s.unread()
		return s.scanIdent()
	case ch == '"':
		return s.scanQuotedIdent()
	case ch == eof:
		return EOF, ""
	case ch == '!':
		ch := s.read()
		if ch == '=' {
			return NOT_EQUAL, "!="
		}
		s.unread()
		return EXCLAMATION_MARK, string(ch)
	case ch == '=':
		ch := s.read()
		if ch == '=' {
			return EQUAL, "=="
		}
		s.unread()
		return EQUAL, string(ch)
	case ch == ',':
		return COMMA, string(ch)
	case ch == '(':
		return OPENING_BRACKET, string(ch)
	case ch == ')':
		return CLOSING_BRACKET, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Lexer) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Lexer) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isValidIdentRune(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.ToUpper(buf.String()) {
	case "NOT":
		return NOT, buf.String()
	case "IN":
		return IN, buf.String()
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}

func (s *Lexer) scanQuotedIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var (
		buf  bytes.Buffer
		last rune
	)

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof || (ch == '"' && last != '\\') {
			break
		} else {
			_, _ = buf.WriteRune(ch)
			last = ch
		}
	}
	return IDENT, buf.String()
}
