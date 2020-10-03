package compiler

import "strings"

type Wordifier struct {
	Reader   *Reader
	Keywords []string
}

func NewWordifier(inputSequence string) *Wordifier {
	return &Wordifier{
		Reader:   NewReader(inputSequence),
		Keywords: make([]string, 0),
	}
}

func (w *Wordifier) IsolateKeywords() []string {
	cur := w.Reader.Peek()
	for w.Reader.HasNext() {
		cur = strings.TrimSpace(cur)

		w.Keywords = append(w.Keywords, cur)

		cur = w.Reader.Next()
	}
	return w.Keywords
}
