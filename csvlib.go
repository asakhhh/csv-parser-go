package main

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

// func (pr *Parser) ReadLine(r io.Reader) (*string, error) {
// }

func (pr *Parser) GetField(n int) (string, error) {
	if n < 0 || n >= len(*pr.fields) {
		return "", ErrFieldCount
	}
	return (*pr.fields)[n], nil
}

func (pr *Parser) GetNumberOfFields() int {
	return pr.numOfFields
}
