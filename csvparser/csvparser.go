package csvparser

import (
	"errors"
	"io"
)

type CSVParser interface {
	ReadLine(r io.Reader) (*string, error)
	GetField(n int) (string, error) // 0-indexed
	GetNumberOfFields() int
}

var (
	ErrQuote      = errors.New("excess or missing \" in quoted-field")
	ErrFieldCount = errors.New("wrong number of fields")
)

type Parser struct {
	line        *string
	fields      *[]string
	numOfFields int
}

func NewParser() *Parser {
	return &Parser{numOfFields: -1}
}

func Join(s []string, sep string) string {
	if len(s) == 0 {
		return ""
	}
	res := s[0]
	for _, str := range s[1:] {
		res += sep + str
	}
	return res
}

func (pr *Parser) ReadLine(r io.Reader) (*string, error) {
	*pr.line = ""
	*pr.fields = (*pr.fields)[:0]
	var curField []string
	field := ""
	b := make([]byte, 1)
	inQuo := false
	quoField := false

	for {
		_, err := r.Read(b)
		if err == io.EOF {
			if inQuo {
				return nil, ErrQuote
			}
			if quoField && (*pr.line)[len(*pr.line)-1] != '"' {
			}
		} else if err != nil {
			return nil, err
		}
	}
}

func (pr *Parser) GetField(n int) (string, error) {
	if n < 0 || n >= len(*pr.fields) {
		return "", ErrFieldCount
	}
	return (*pr.fields)[n], nil
}

func (pr *Parser) GetNumberOfFields() int {
	return pr.numOfFields
}
