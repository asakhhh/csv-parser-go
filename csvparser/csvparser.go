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
	curField := []string{""}
	b := make([]byte, 1)
	inQuo := false
	quoField := false
	started := false

	for {
		_, err := r.Read(b)
		if err == io.EOF {
			if inQuo {
				return nil, ErrQuote
			}
			// todo
		} else if err != nil {
			return nil, err
		}

		if b[0] == '"' {
			if !started {
				quoField = true
				started = true
			} else if !quoField {
				return nil, ErrQuote
			}
			inQuo = !inQuo
			curField[len(curField)-1] += string(b)
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
