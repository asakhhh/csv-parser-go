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
	return &Parser{line: new(string), fields: new([]string), numOfFields: -1, endedOnCR: false}
}

func (pr *Parser) scanLine() error {
	if len(*pr.line) == 0 {
		if pr.numOfFields == -1 {
			pr.numOfFields = 0
		} else if pr.numOfFields != 0 {
			return ErrFieldCount
		}
		return nil
	}

	line := *pr.line
	inQuote := false
	*pr.fields = append(*pr.fields, "")
	for _, c := range line {
		if c == ',' && !inQuote {
			*pr.fields = append(*pr.fields, "")
		} else {
			(*pr.fields)[len(*pr.fields)-1] += string(c)
			if c == '"' {
				inQuote = !inQuote
			}
		}
	}
	if inQuote {
		return ErrQuote
	}
	if pr.numOfFields == -1 {
		pr.numOfFields = len(*pr.fields)
	} else if pr.numOfFields != len(*pr.fields) {
		return ErrFieldCount
	}

	var isQ []bool
	for i := range *pr.fields {
		if len((*pr.fields)[i]) == 0 {
			isQ = append(isQ, false)
			continue
		}
		if ((*pr.fields)[i][0] == '"') != ((*pr.fields)[i][len((*pr.fields)[i])-1] == '"') {
			return ErrQuote
		}
		isQ = append(isQ, (*pr.fields)[i][0] == '"')
		if !isQ[i] && contains((*pr.fields)[i], "\"") {
			return ErrQuote
		}
	}

	for i := range *pr.fields {
		if !isQ[i] {
			continue
		}
		(*pr.fields)[i] = (*pr.fields)[i][1 : len((*pr.fields)[i])-1]
		var sl []string
		t := ""
		for j := 0; j < len((*pr.fields)[i]); {
			if (*pr.fields)[i][j] == '"' {
				sl = append(sl, t)
				t = ""
				for j < len((*pr.fields)[i]) && (*pr.fields)[i][j] == '"' {
					t += "\""
					j++
				}
				if len(t)%2 == 1 {
					return ErrQuote
				}
				sl = append(sl, t)
				t = ""
			} else {
				t += string((*pr.fields)[i][j])
				j++
			}
		}
		sl = append(sl, t)
		(*pr.fields)[i] = ""
		for j := range sl {
			(*pr.fields)[i] += sl[j][:len(sl[j])/(j%2+1)]
		}
	}

	return nil
}

func (pr *Parser) ReadLine(r io.Reader) (*string, error) {
	b := make([]byte, 1)
	*pr.line = ""
	*pr.fields = (*pr.fields)[:0]
	_, err := r.Read(b)
	if err != nil {
		return nil, err
	}
	if b[0] == '\n' && pr.endedOnCR {
		pr.endedOnCR = false
		return pr.ReadLine(r)
	}

	inQuote := false
	for (b[0] != '\n' && b[0] != '\r') || inQuote {
		*pr.line += string(b)
		if b[0] == '"' {
			inQuote = !inQuote
		}
		_, err = r.Read(b)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	if err != io.EOF {
		pr.endedOnCR = (b[0] == '\r')
	}
	if inQuote {
		*pr.line = ""
		return nil, err
	}
	err = pr.scanLine()
	if err != nil {
		*pr.line = ""
		*pr.fields = (*pr.fields)[:0]
	}
	return pr.line, err
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
