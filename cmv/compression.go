package cmv

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
)

type Compression1Reader struct {
	r io.Reader
	z io.ReadCloser // also implements zlib.Resetter
}

func NewCompression1Reader(r io.Reader) io.Reader {
	return &Compression1Reader{r: r}
}

func (r *Compression1Reader) Read(b []byte) (int, error) {
	if r.z == nil {
		l, err := r.readLength()
		if err != nil {
			return 0, err
		}

		r.z, err = zlib.NewReader(io.LimitReader(r.r, l))
		if err != nil {
			return 0, err
		}
	}
	n, err := r.z.Read(b)
	if err == io.EOF {
		if n != 0 {
			return n, nil
		}

		l, err := r.readLength()
		if err != nil {
			return 0, err
		}

		err = r.z.(zlib.Resetter).Reset(io.LimitReader(r.r, l), nil)
		if err != nil {
			return 0, err
		}

		return r.z.Read(b)
	}
	return n, err
}

func (r *Compression1Reader) readLength() (int64, error) {
	var n uint32
	err := binary.Read(r.r, binary.LittleEndian, &n)
	return int64(n), err
}

type Compression1Writer struct {
	w io.Writer
	b bytes.Buffer
	z *zlib.Writer
}

func NewCompression1Writer(w io.Writer) io.Writer {
	return &Compression1Writer{w: w}
}

func (w *Compression1Writer) Write(b []byte) (int, error) {
	w.b.Reset()

	if w.z == nil {
		w.z = zlib.NewWriter(&w.b)
	} else {
		w.z.Reset(&w.b)
	}

	n, err := w.z.Write(b)
	if err == nil && n != len(b) {
		err = io.ErrShortWrite
	}
	if err != nil {
		return 0, err
	}

	err = w.z.Close()
	if err != nil {
		return 0, err
	}

	err = binary.Write(w.w, binary.LittleEndian, uint32(w.b.Len()))
	if err != nil {
		return 0, err
	}

	n, err = w.w.Write(w.b.Bytes())
	if err == nil && n != w.b.Len() {
		err = io.ErrShortWrite
	}
	if err != nil {
		return 0, err
	}

	return len(b), nil
}
