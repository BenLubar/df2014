package wtf23a

import (
	"encoding/binary"
	"errors"
	"io"
)

// Header is the header used to store the initial state of a wtf23a segment.
type Header struct {
	State [624]uint32
	Index uint32
}

// ErrInvalidIndex is returned if a header has an index greater than the size
// of the state array or is not an integer.
var ErrInvalidIndex = errors.New("wtf23a: invalid header index")

// ReadHeader reads a Header from an io.Reader.
func ReadHeader(r io.Reader) (Header, error) {
	var h Header

	if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
		return h, err
	}

	if h.Index%4 != 0 {
		return h, ErrInvalidIndex
	}

	h.Index /= 4

	if h.Index > uint32(len(h.State)) {
		return h, ErrInvalidIndex
	}

	return h, nil
}

// Next returns the next value. It modifies the State and Index of the Header.
func (h *Header) Next() uint32 {
	if h.Index == uint32(len(h.State)) {
		for i := range h.State {
			y := (h.State[i] & 0x80000000) | (h.State[(i+1)%624] & 0x7FFFFFFF)
			h.State[i] = h.State[(i+397)%624] ^ ((y >> 1) & 0x7FFFFFFF)
			if y&1 != 0 {
				h.State[i] ^= 0x9908b0df
			}
		}
		h.Index = 0
	}
	n := h.State[h.Index]
	h.Index++
	return n
}
