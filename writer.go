package qpencoding

import (
	"encoding/hex"
	"io"
	"strings"
)

type Writer struct {
	w    io.Writer
	line []byte
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	for n < len(p) {
		octet := p[n]
		switch {
		case octet == ' ', octet == '\t', 33 <= octet && octet != '=' && octet <= 126:
			w.line = append(w.line, octet)
		case octet == '\r', octet == '\n':
			if n, err = w.appendCRLF(p, n); err != nil {
				return
			}
		default:
			w.appendInHex(p[n : n+1])
		}

		if len(w.line) >= 75 {
			if n+1 < len(p) && (p[n+1] == '\n' || p[n+1] == '\r') {
				if n, err = w.appendCRLF(p, n+1); err != nil {
					return
				}
			} else {
				w.line = append(w.line, '=', '\r', '\n')
				if err = w.flush(); err != nil {
					return
				}
			}
		}

		n++
	}

	if len(w.line) > 0 {
		if w.line[len(w.line)-1] != '\n' {
			w.line = append(w.line, '=', '\r', '\n')
		}
		if err = w.flush(); err != nil {
			return
		}
	}

	return
}

func (w *Writer) appendCRLF(p []byte, n int) (int, error) {
	if w.endWithWhiteSpace() {
		sp := w.line[len(w.line)-1]
		w.line = w.line[:len(w.line)-1]
		if len(w.line) > 73 {
			if w.line[len(w.line)-1] != '\n' {
				w.line = append(w.line, '=')
			}
			w.line = append(w.line, '\r', '\n')
			if err := w.flush(); err != nil {
				return n, err
			}
		}
		w.appendInHex([]byte{sp})
	}

	if p[n] == '\n' {
		if n+1 < len(p) && p[n+1] == '\r' {
			n++
		}
	} else if p[n] == '\r' {
		if n+1 < len(p) && p[n+1] == '\n' {
			n++
		}
	}

	w.line = append(w.line, '\r', '\n')

	return n, w.flush()
}

func (w *Writer) endWithWhiteSpace() bool {
	return len(w.line) > 0 && (w.line[len(w.line)-1] == '\t' || w.line[len(w.line)-1] == ' ')
}

func (w *Writer) appendInHex(p []byte) {
	dump := hex.EncodeToString(p)
	dump = strings.ToUpper(dump)
	w.line = append(w.line, '=', dump[0], dump[1])
}

func (w *Writer) flush() (err error) {
	_, err = w.w.Write(w.line)
	if err != nil {
		return
	}

	w.line = []byte{}
	return
}
