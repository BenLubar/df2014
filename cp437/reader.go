package cp437

import (
	"bytes"
	"io"
)

type reader struct {
	r io.Reader
	b bytes.Buffer
}

// NewReader returns an io.Reader that reads from r and converts CP437-encoded
// data to UTF-8-encoded data.
func NewReader(r io.Reader) io.Reader {
	return &reader{r: r}
}

func (r *reader) Read(b []byte) (n int, err error) {
	n, err = r.b.Read(b)
	if n != 0 {
		return
	}
	n, err = r.r.Read(b)
	if n == 0 && err != nil {
		return
	}
	for _, c := range b[:n] {
		_, err = r.b.WriteRune(Rune(c))
		if err != nil {
			// bytes.Buffer should never return an error, but just
			// in case.
			panic(err)
		}
	}
	return r.b.Read(b)
}
