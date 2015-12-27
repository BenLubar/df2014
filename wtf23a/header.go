package wtf23a

import (
	"encoding/binary"
	"errors"
	"io"
)

type Header struct {
	State [624]uint32
	Index uint32
}

var ErrInvalid = errors.New("wtf23a: invalid header")

func ReadHeader(r io.Reader) (Header, error) {
	var h Header

	if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
		return h, err
	}

	if h.Index%4 != 0 {
		return h, ErrInvalid
	}

	h.Index /= 4

	if h.Index > uint32(len(h.State)) {
		return h, ErrInvalid
	}

	return h, nil
}

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
