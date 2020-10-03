package compiler

import "strings"

type Reader struct {
	Sequence []string
	Index    int
}

func NewReader(sequence string) *Reader {
	return &Reader{
		Sequence: strings.Split(sequence, " "),
		Index:    0,
	}
}

func (r *Reader) Peek() string {
	return r.Sequence[r.Index]
}

func (r *Reader) HasNext() bool {
	return r.Index+1 < len(r.Sequence)
}

func (r *Reader) Next() string {
	r.Index++
	return r.Peek()
}

func (r *Reader) ValidKeyword() string {
	cur := r.Peek()
	for _, tokenType := range KindsList {
		re := tokenType.Regex()
		if token := re.FindString(cur); len(token) > 0 {
			return token
		}
	}
	return ""
}

func (r *Reader) CurrentTokenKind() TokenKind {
	token := r.ValidKeyword()
	if len(token) == 0 {
		return UNKNOWN
	}
	return ConvertTokenToKind(token)
}
