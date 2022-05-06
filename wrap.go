package pht

import (
	"bytes"
	"io"
	"unicode/utf8"
)

// I tried to find a word wrapper that worked like this on godoc, but I
// couldn't find one.
//
// This tries to do unicode rune counting but

type WordWrapper struct {
	w io.Writer

	indent          int
	initialIndent   int
	maxLine         int
	currLine        int
	wordAccumulator []byte
	lineBuf         []byte
	inWhitespace    bool

	// actually, this is more like "have never printed", I think
	atBeginningOfLine bool
}

func NewWrapper(w io.Writer,
	initialIndent int,
	indent int,
	maxLine int) *WordWrapper {
	return &WordWrapper{
		w:             w,
		indent:        indent,
		initialIndent: initialIndent,
		maxLine:       maxLine,
		// This is just sort of a guess to try to prevent
		// thrashing, it doesn't really matter.
		wordAccumulator:   make([]byte, 0, (indent+maxLine)*2),
		lineBuf:           make([]byte, 0, (indent+maxLine)*2),
		atBeginningOfLine: true,
	}
}

func isWhitespace(b byte) bool {
	return b == ' ' ||
		b == '\t' ||
		b == '\n' ||
		b == '\r'
}

func WrapString(s string, initialIndent int, indent int, maxLine int) string {
	buf := &bytes.Buffer{}
	wrapper := NewWrapper(buf, initialIndent, indent, maxLine)
	wrapper.Write([]byte(s))
	wrapper.Close()
	return buf.String()
}

// Write accepts the incoming bytes as words to be wrapped.
func (ww *WordWrapper) Write(b []byte) (int, error) {
	// For every byte, do what we need to do.
	for _, b := range b {
		if isWhitespace(b) {
			if ww.inWhitespace {
				// Collapse all whitespace down
				continue
			}

			// Otherwise, this completes the word.
			ww.inWhitespace = true
			ww.writeWord()
		} else {
			ww.inWhitespace = false
			ww.wordAccumulator = append(ww.wordAccumulator, b)
		}
	}

	return len(b), nil
}

func (ww *WordWrapper) Close() error {
	ww.writeWord()
	return nil
}

func (ww *WordWrapper) writeWord() {
	if len(ww.wordAccumulator) > 0 {
		count := utf8.RuneCount(ww.wordAccumulator)
		if !ww.atBeginningOfLine && ww.currLine+1+count > ww.maxLine {
			ww.w.Write([]byte("\n"))
			writeRawIndent(ww.indent, ww.w)
			ww.w.Write(ww.wordAccumulator)
			ww.currLine = count
		} else {
			if ww.atBeginningOfLine {
				writeRawIndent(ww.initialIndent, ww.w)
				ww.w.Write(ww.wordAccumulator)
				ww.currLine += count
				ww.atBeginningOfLine = false
			} else {
				ww.w.Write([]byte(" "))
				ww.w.Write(ww.wordAccumulator)
				ww.currLine += 1 + count
			}
		}
		ww.wordAccumulator = ww.wordAccumulator[:0]
	}
}

func writeRawIndent(indent int, w io.Writer) {
	if indent <= 0 {
		return
	}
	b := make([]byte, indent)
	for i := 0; i < indent; i++ {
		b[i] = 32
	}
	w.Write(b)
}
