// Package wtf23a implements the encoding used by Dwarf Fortress for save files
// created before version 40d.
//
// Ported from http://dwarffortresswiki.org/index.php/User:Quietust/wtf23a.php
package wtf23a

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"io/ioutil"
)

type reader struct {
	r   io.Reader
	buf []byte
}

func NewReader(r io.Reader) io.Reader {
	return &reader{r: r}
}

func (r *reader) Read(p []byte) (n int, err error) {
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

func (r *reader) fill() (err error) {
	initial, err := ReadHeader(r.r)
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

	var skip [1]byte
	data := make([]byte, length)

	// read data, skipping dummy bytes
	rng := initial
	for i := length - 1; i >= 0; i-- {
		_, err = r.r.Read(data[i : i+1])
		if err != nil {
			return
		}

		if rng.Next()%2 == 1 {
			rng.Next()
			_, err = r.r.Read(skip[:])
			if err != nil {
				return
			}
		}
	}

	// decode data
	rng = initial
	for i := length - 1; i >= 0; i-- {
		data[i] -= byte(rng.Next() % 10)
	}

	z, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return
	}
	defer z.Close()

	r.buf, err = ioutil.ReadAll(z)

	return
}

// The Aristocrats
