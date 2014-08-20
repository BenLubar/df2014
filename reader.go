package df2014

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
)

type Reader struct {
	io.Reader
}

func (r *Reader) Int8() (n int8, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}

func (r *Reader) Uint8() (n uint8, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}

func (r *Reader) Int16() (n int16, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}

func (r *Reader) Uint16() (n uint16, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}

func (r *Reader) Int32() (n int32, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}

func (r *Reader) Uint32() (n uint32, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}

func (r *Reader) Int64() (n int64, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}

func (r *Reader) Uint64() (n uint64, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}

var cp437 = []rune("\x00☺☻♥♦♣♠•◘○◙♂♀♪♬☼►◄↕‼¶§▬↨↑↓→←∟↔▲▼ !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⌂ÇüéâäàåçêëèïîìÄÅÉæÆôöòûùÿÖÜ¢£¥₧ƒáíóúñÑªº¿⌐¬½¼¡«»░▒▓│┤╡╢╖╕╣║╗╝╜╛┐└┴┬├─┼╞╟╚╔╩╦╠═╬╧╨╤╥╙╘╒╓╫╪┘┌█▄▌▐▀αßΓπΣσµτΦΘΩδ∞φε∩≡±≥≤⌠⌡÷≈°∙·√ⁿ²■\xA0")

func (r *Reader) String() (string, error) {
	length, err := r.Uint16()
	if err != nil {
		return "", err
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}

	s := make([]rune, length)
	for i, b := range buf {
		s[i] = cp437[b]
	}

	return string(s), nil
}

func (r *Reader) Header() (version, compression uint32, err error) {
	version, err = r.Uint32()
	if err != nil {
		return 0, 0, err
	}

	compression, err = r.Uint32()
	if err != nil {
		return 0, 0, err
	}

	switch compression {
	case 0:
		// nothing to be done
	case 1:
		r.Reader = &compression1Reader{r: &Reader{r.Reader}}
	default:
		return 0, 0, fmt.Errorf("unhandled compression type %d", compression)
	}

	return
}

type Name struct {
	First    string
	Nick     string
	Index    [7]int32
	Form     [7]uint16
	Language uint32
	Unknown  int16
}

func (r *Reader) Name() (name Name, err error) {
	name.First, err = r.String()
	if err != nil {
		return
	}
	name.Nick, err = r.String()
	if err != nil {
		return
	}
	for i := range name.Index {
		name.Index[i], err = r.Int32()
		if err != nil {
			return
		}
	}
	for i := range name.Form {
		name.Form[i], err = r.Uint16()
		if err != nil {
			return
		}
	}
	name.Language, err = r.Uint32()
	if err != nil {
		return
	}
	name.Unknown, err = r.Int16()
	if err != nil {
		return
	}
	return
}

type compression1Reader struct {
	r   *Reader
	buf []byte
}

func (r *compression1Reader) Read(p []byte) (n int, err error) {
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

func (r *compression1Reader) fill() (err error) {
	length, err := r.r.Uint32()
	if err != nil {
		return
	}

	var buf bytes.Buffer
	z, err := zlib.NewReader(io.LimitReader(r.r, int64(length)))
	if err != nil {
		return
	}
	defer func() {
		e := z.Close()
		if err == nil {
			err = e
		}
	}()

	_, err = io.Copy(&buf, z)
	if err != nil {
		return
	}

	r.buf = buf.Bytes()
	return
}
