package filehelper

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/shoobyban/mxj"
	"github.com/shoobyban/slog"
)

// ParserFunc is to parse a []byte into an interface{}
type ParserFunc func([]byte) (interface{}, error)

// Parser is the main type
type Parser struct {
	parsers map[string]ParserFunc
}

// NewParser defines a new parser
func NewParser() *Parser {
	return &Parser{
		parsers: map[string]ParserFunc{
			"xml": func(content []byte) (interface{}, error) { return mxj.NewMapXml(content) },
			"json": func(content []byte) (interface{}, error) {
				var out interface{}
				err := json.Unmarshal(content, &out)
				return out, err
			},
			"csv": func(content []byte) (interface{}, error) {
				r := csv.NewReader(bytes.NewBuffer(content))
				return r.ReadAll()
			},
		},
	}
}

// RegisterParser registers or overrides a format parser func. Indices are lower case.
func (l *Parser) RegisterParser(format string, parser ParserFunc) {
	l.parsers[format] = parser
}

// ReadStruct reads from given file, parsing into structure
func (l *Parser) ReadStruct(filename, format string) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		slog.Infof("Can't open file %s", filename)
		return nil, err
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)
	return l.ParseStruct(byteValue, format)
}

// ParseStruct parses byte slice into map or slice
func (l *Parser) ParseStruct(content []byte, format string) (interface{}, error) {
	var out interface{}
	var err error
	if parser, ok := l.parsers[format]; ok {
		out, err = parser(content)
	} else {
		return nil, errors.New("Unknown file")
	}
	if err != nil {
		return nil, fmt.Errorf("Can't parse %s: %v", format, err)
	}
	return out, nil
}
