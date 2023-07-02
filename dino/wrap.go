package dino

import (
	"bytes"
	"unicode"
)

type Wrapper struct {
	MaxLineLength int
	buf           bytes.Buffer
	space         bytes.Buffer
	word          bytes.Buffer
	allLines      []string
	lineLen       int
}

func NewWrapper(limit int) *Wrapper {
	return &Wrapper{
		MaxLineLength: limit,
	}
}

func (w *Wrapper) splitString(s string) []string {
	for _, c := range s {
		if c == '\n' {
			// We've reached the end of line.
			// We don't have a word buffered.
			// Check if we can still add the content of the space buffer to the current line.
			if w.word.Len() == 0 {
				if w.lineLen+w.space.Len() > w.MaxLineLength {
					w.lineLen = 0
				} else {
					// Preserve existing whitespace.
					w.lineLen += w.space.Len()
					_, _ = w.buf.Write(w.space.Bytes())
				}
			} else {
				// Add the current word, the content of the space buffer, and a newline.
				w.lineLen += w.space.Len() + w.word.Len()
				_, _ = w.buf.Write(w.space.Bytes())
				_, _ = w.buf.Write(w.word.Bytes())
				w.word.Reset()
			}
			w.space.Reset()
			w.allLines = append(w.allLines, w.buf.String())
			w.buf.Reset()
			w.lineLen = 0
		} else if unicode.IsSpace(c) {
			// We've reached the end of current word.
			if w.space.Len() == 0 || w.word.Len() > 0 {
				w.lineLen += w.space.Len() + w.word.Len()
				_, _ = w.buf.Write(w.space.Bytes())
				w.space.Reset()
				_, _ = w.buf.Write(w.word.Bytes())
				w.word.Reset()
			}
			w.space.WriteRune(c)
		} else {
			// Any other character
			w.word.WriteRune(c)
			// If the current word would cause the current line to exceed the
			// maximum line length, add a line break.
			if w.lineLen+w.word.Len()+w.space.Len() > w.MaxLineLength && w.word.Len() < w.MaxLineLength {
				w.allLines = append(w.allLines, w.buf.String())
				w.buf.Reset()
				w.lineLen = 0
				w.space.Reset()
			}
		}
	}

	if w.word.Len() != 0 {
		// Add the current word, the content of the space buffer, and a newline.
		_, _ = w.buf.Write(w.space.Bytes())
		_, _ = w.buf.Write(w.word.Bytes())
	} else if w.lineLen+w.space.Len() <= w.MaxLineLength {
		// We don't have a word buffered.
		// Check if we can still add the content of the space buffer to the current line.
		_, _ = w.buf.Write(w.space.Bytes())
	}

	return append(w.allLines, w.buf.String())
}
