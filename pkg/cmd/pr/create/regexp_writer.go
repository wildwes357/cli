package create

import (
	"bytes"
	"io"
	"regexp"
)

func NewRegexpWriter(out io.Writer, re *regexp.Regexp, repl string) *RegexpWriter {
	return &RegexpWriter{out: out, re: *re, repl: repl, buf: []byte{}}
}

type RegexpWriter struct {
	out  io.Writer
	re   regexp.Regexp
	repl string
	buf  []byte
}

func (s *RegexpWriter) Write(data []byte) (int, error) {
	dataLength := len(data)
	if dataLength == 0 {
		return 0, nil
	}

	mod := data
	filtered := []byte{}
	repl := []byte(s.repl)

	for {
		if i := bytes.IndexByte(mod, '\n'); i >= 0 {
			b := append(s.buf, mod[0:i]...)
			f := s.re.ReplaceAll(b, repl)
			if len(f) > 0 {
				filtered = append(filtered, f...)
				filtered = append(filtered, '\n')
			}
			mod = mod[i+1:]
			s.buf = []byte{}
		} else {
			s.buf = append(s.buf, mod...)
			break
		}
	}

	if len(filtered) != 0 {
		_, err := s.out.Write(filtered)
		if err != nil {
			return 0, err
		}
	}

	return dataLength, nil
}

func (s *RegexpWriter) Flush() (int, error) {
	if len(s.buf) > 0 {
		repl := []byte(s.repl)
		filtered := s.re.ReplaceAll(s.buf, repl)
		if len(filtered) > 0 {
			return s.out.Write(filtered)
		}
	}

	return 0, nil
}
