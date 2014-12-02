package df2014

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
)

// Ported from http://dwarffortresswiki.org/index.php/User:Quietust/wtf23a.php

type wtf23a struct {
	State [624]uint32
	Index uint32
}

func (wtf *wtf23a) init(r io.Reader) error {
	err := binary.Read(r, binary.LittleEndian, wtf)
	if err != nil {
		return err
	}

	if wtf.Index%4 != 0 {
		return fmt.Errorf("df2014: cannot guess compression type")
	}

	wtf.Index /= 4

	if wtf.Index >= 624 {
		return fmt.Errorf("df2014: cannot guess compression type")
	}

	return nil
}

func (wtf *wtf23a) next(mod uint32) uint32 {
	if wtf.Index == 624 {
		for i := range wtf.State {
			y := (wtf.State[i] & 0x80000000) | (wtf.State[(i+1)%624] & 0x7FFFFFFF)
			wtf.State[i] = wtf.State[(i+397)%624] ^ ((y >> 1) & 0x7FFFFFFF)
			if y&1 != 0 {
				wtf.State[i] ^= 0x9908b0df
			}
		}
		wtf.Index = 0
	}
	n := wtf.State[wtf.Index]
	wtf.Index++
	return n % mod
}

type wtf23aReader struct {
	r   io.Reader
	buf []byte
}

func (r *wtf23aReader) Read(p []byte) (n int, err error) {
	if len(r.buf) == 0 {
		err = r.fill()
		if err != nil {
			return
		}
	}

	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	return
}

func (r *wtf23aReader) fill() (err error) {
	var initial wtf23a
	err = initial.init(r.r)
	if err != nil {
		return
	}

	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	var length int32
	err = binary.Read(r.r, binary.LittleEndian, &length)
	if err != nil {
		return
	}
	if length < 0 {
		return fmt.Errorf("df2014: negative length (%d)", length)
	}

	var skip [1]byte
	data := make([]byte, length)

	// read data, skipping dummy bytes
	rng := initial
	for i := length - 1; i >= 0; i-- {
		_, err = r.r.Read(data[i : i+1])
		if err != nil {
			return
		}

		if rng.next(2) == 1 {
			rng.next(256)
			_, err = r.r.Read(skip[:])
			if err != nil {
				return
			}
		}
	}

	// decode data
	rng = initial
	for i := length - 1; i >= 0; i-- {
		data[i] -= byte(rng.next(10))
	}

	z, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return
	}
	defer z.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, z)
	if err != nil {
		return
	}

	r.buf = buf.Bytes()

	return
}
