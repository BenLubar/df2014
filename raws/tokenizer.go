package raws

import (
	"bufio"
	"io"
	"strings"
)

type Tokenizer struct {
	r *bufio.Reader
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{bufio.NewReader(r)}
}

func (t *Tokenizer) Next() ([]string, error) {
	_, err := t.r.ReadString('[')
	if err != nil {
		return nil, err
	}
	s, err := t.r.ReadString(']')
	if err != nil {
		return nil, err
	}

	s = s[:len(s)-1]

	return strings.Split(s, ":"), nil
}
