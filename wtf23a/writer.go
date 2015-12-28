package wtf23a

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"math/rand"
)

type writer struct {
	w io.Writer
	h func() Header
}

// NewWriter wraps an io.Writer to output Dwarf Fortress 23a's obfuscated file
// encoding. ZeroHeader and RandomHeader(r) are predefined header generators.
func NewWriter(w io.Writer, generateHeader func() Header) io.Writer {
	return &writer{w: w, h: generateHeader}
}

func (w *writer) Write(b []byte) (int, error) {
	n := len(b)
	if n == 0 {
		return 0, nil
	}

	initial := w.h()
	if initial.Index > uint32(len(initial.State)) {
		return 0, ErrInvalidIndex
	}

	var buf bytes.Buffer
	z := zlib.NewWriter(&buf)
	if _, err := z.Write(b); err != nil {
		return 0, err
	}
	if err := z.Close(); err != nil {
		return 0, err
	}
	data := buf.Bytes()

	rng := initial
	for i := n - 1; i >= 0; i-- {
		data[i] += byte(rng.Next() % 10)
	}

	filled := make([]byte, 0, n)

	rng = initial
	for i := n - 1; i >= 0; i-- {
		filled = append(filled, data[i])

		if rng.Next()%2 == 1 {
			filled = append(filled, byte(rng.Next()))
		}
	}

	initial.Index *= 4

	if err := binary.Write(w.w, binary.LittleEndian, &initial); err != nil {
		return 0, err
	}

	if err := binary.Write(w.w, binary.LittleEndian, uint32(n)); err != nil {
		return 0, err
	}

	if _, err := w.w.Write(filled); err != nil {
		return 0, err
	}

	return n, nil
}

// ZeroHeader is a function that can be used with NewWriter to generate headers
// that do not obfuscate the data.
func ZeroHeader() Header {
	return Header{}
}

// RandomHeader returns a function that can be used with NewWriter to generate
// random valid header values.
func RandomHeader(r *rand.Rand) func() Header {
	return func() Header {
		var h Header
		h.Index = uint32(r.Intn(len(h.State) + 1))
		for i := range h.State {
			h.State[i] = r.Uint32()
		}
		return h
	}
}
