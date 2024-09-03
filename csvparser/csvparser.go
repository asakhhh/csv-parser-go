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
	endedOnCR   bool
}

func NewParser() *Parser {
	return &Parser{numOfFields: -1, endedOnCR: false}
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

func (pr *Parser) readQuote(r io.Reader) (*string, error) {
	b := make([]byte, 1)
	line := ""
	for {
		_, err := r.Read(b)
		if err == io.EOF {
			return nil, ErrQuote
		} else if err != nil {
			return nil, err
		}
		if b[0] == '"' {
			return &line, nil
		}
		line += string(b)
	}
}

func last(i any) any {
	switch v := i.(type) {
	case string:
		return v[len(v)-1]
	case []int:
		return v[len(v)-1]
	case []string:
		return v[len(v)-1]
	case []byte:
		return v[len(v)-1]
	case []rune:
		return v[len(v)-1]
	case []bool:
		return v[len(v)-1]
	case *string:
		return (*v)[len(*v)-1]
	case *[]int:
		return (*v)[len(*v)-1]
	case *[]string:
		return (*v)[len(*v)-1]
	case *[]byte:
		return (*v)[len(*v)-1]
	case *[]rune:
		return (*v)[len(*v)-1]
	case *[]bool:
		return (*v)[len(*v)-1]
	default:
		return nil
	}
}

func (pr *Parser) ReadLine(r io.Reader) (*string, error) {
	b := make([]byte, 1)
	_, err := r.Read(b)
	if err != nil {
		return nil, err
	}

	if b[0] == '\n' && pr.endedOnCR {
		pr.endedOnCR = false
		return pr.ReadLine(r)
	}
	if b[0] == '\n' || b[0] == '\r' {
		*pr.line = ""
		*pr.fields = []string{}
		pr.endedOnCR = (b[0] == '\r')

		if pr.numOfFields == -1 {
			pr.numOfFields = 0
			return pr.line, nil
		}
		if pr.numOfFields != 0 {
			return nil, ErrFieldCount
		}
		return pr.line, nil
	}

	*pr.line = ""
	*pr.fields = []string{""}
	for {
		if b[0] == '"' {
			if len(*pr.line) == 0 || last(pr.line) == "" {
				line, err := pr.readQuote(r)
				if err != nil {
					return nil, err
				}
				*pr.fields = append(*pr.fields, *line)
				*pr.line += *line
			} else if last(pr.line) != '"' {
				return nil, ErrQuote
			} else {
				line, err := pr.readQuote(r)
				if err != nil {
					return nil, err
				}
				(*pr.fields)[len(*pr.fields)-1] += "\"" + *line
				*pr.line += *line
			}
		} else if b[0] == ',' {
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
