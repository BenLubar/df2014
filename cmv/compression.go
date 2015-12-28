package cmv

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
)

type compression1Reader struct {
	r io.Reader
	z io.ReadCloser // also implements zlib.Resetter
}

// NewCompression1Reader wraps an io.Reader to decode Dwarf Fortress's chunked
// zlib format. The format consists of a stream of a uint32 length followed by
// that many bytes of compressed data.
func NewCompression1Reader(r io.Reader) io.Reader {
	return &compression1Reader{r: r}
}

func (r *compression1Reader) Read(b []byte) (n int, err error) {
	if r.z == nil {
		if err = r.fill(); err != nil {
			return
		}
	}
	n, err = r.z.Read(b)

	if err == io.EOF {
		if n != 0 {
			err = nil
			return
		}

		if err = r.fill(); err != nil {
			return
		}

		return r.z.Read(b)
	}

	return
}

func (r *compression1Reader) fill() error {
	var length uint32
	err := binary.Read(r.r, binary.LittleEndian, &length)
	if err != nil {
		return err
	}

	section := io.LimitReader(r.r, int64(length))
	if r.z == nil {
		r.z, err = zlib.NewReader(section)
	} else {
		err = r.z.(zlib.Resetter).Reset(section, nil)
	}

	return err
}

type compression1Writer struct {
	w io.Writer
	b bytes.Buffer
	z *zlib.Writer
}

// NewCompression1Writer wraps an io.Writer to output Dwarf Fortress's chunked
// zlib format. Each call to Write outputs its own frame, so the returned
// io.Writer should be wrapped in bufio.NewWriter if there will be many small
// writes.
func NewCompression1Writer(w io.Writer) io.Writer {
	return &compression1Writer{w: w}
}

func (w *compression1Writer) Write(b []byte) (int, error) {
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
