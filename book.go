package df2014

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Book struct {
	Unk000 uint32
	Unk001 uint32
	Unk002 uint16
	Unk003 int32
	Unk004 uint16
	Unk005 uint32
	Unk006 uint32
	Unk007 *[74]byte
	Unk008 *[40]byte
	Unk009 string
	Unk010 *[6]byte
}

func (r *Reader) bookList() (l []Book, err error) {
	for {
		var signature uint64
		err = binary.Read(r, binary.LittleEndian, &signature)
		if err != nil {
			return
		}
		if signature != 0x8ad08ad08ad0 {
			// put the bytes we just read back on the buffer.
			buf := make([]byte, 8)
			binary.LittleEndian.PutUint64(buf, signature)
			r.Reader = io.MultiReader(bytes.NewReader(buf), r.Reader)
			return
		}

		var book Book
		err = binary.Read(r, binary.LittleEndian, &book.Unk000)
		if err != nil {
			return
		}

		{
			var n [3]uint16
			err = binary.Read(r, binary.LittleEndian, &n)
			if err != nil {
				return
			}
			if expected := [3]uint16{0, 0, 0}; n != expected {
				err = fmt.Errorf("df2014: book expectation failed: %v != %v", n, expected)
				return
			}
		}

		err = binary.Read(r, binary.LittleEndian, &book.Unk001)
		if err != nil {
			return
		}
		if len(l) > 0 && l[len(l)-1].Unk001 >= book.Unk001 {
			err = fmt.Errorf("df2014: book expectation failed: %d >= %d", l[len(l)-1].Unk001, book.Unk001)
			return
		}

		{
			var n [3]int32
			err = binary.Read(r, binary.LittleEndian, &n)
			if err != nil {
				return
			}
			if expected := [3]int32{0, 1, 1}; n != expected {
				err = fmt.Errorf("df2014: book expectation failed: %v != %v", n, expected)
				return
			}
		}

		{
			var n uint32
			err = binary.Read(r, binary.LittleEndian, &n)
			if err != nil {
				return
			}
			if expected := book.Unk001; n != expected {
				err = fmt.Errorf("df2014: book expectation failed: %v != %v", n, expected)
				return
			}
		}

		{
			var n [3]int32
			err = binary.Read(r, binary.LittleEndian, &n)
			if err != nil {
				return
			}
			if expected := [3]int32{-1, -1, 1}; n != expected {
				err = fmt.Errorf("df2014: book expectation failed: %v != %v", n, expected)
				return
			}
		}

		{
			var n [3]uint8
			err = binary.Read(r, binary.LittleEndian, &n)
			if err != nil {
				return
			}
			if expected := [3]uint8{0, 0, 0}; n != expected {
				err = fmt.Errorf("df2014: book expectation failed: %v != %v", n, expected)
				return
			}
		}

		{
			var n uint16
			err = binary.Read(r, binary.LittleEndian, &n)
			if err != nil {
				return
			}
			if expected := uint16(0x2742); n != expected {
				err = fmt.Errorf("df2014: book expectation failed: %v != %v", n, expected)
				return
			}
		}

		{
			var n [3]int32
			err = binary.Read(r, binary.LittleEndian, &n)
			if err != nil {
				return
			}
			if expected := [3]int32{0, 0, -1}; n != expected {
				err = fmt.Errorf("df2014: book expectation failed: %v != %v", n, expected)
				return
			}
		}

		err = binary.Read(r, binary.LittleEndian, &book.Unk002)
		if err != nil {
			return
		}

		err = binary.Read(r, binary.LittleEndian, &book.Unk003)
		if err != nil {
			return
		}

		{
			var n int16
			err = binary.Read(r, binary.LittleEndian, &n)
			if err != nil {
				return
			}
			if expected := int16(-1); n != expected {
				err = fmt.Errorf("df2014: book expectation failed: %v != %v", n, expected)
				return
			}
		}

		err = binary.Read(r, binary.LittleEndian, &book.Unk004)
		if err != nil {
			return
		}
		if a, b := book.Unk004 == 5, book.Unk000&0x800 == 0x800; a != b {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", a, b)
			return
		}

		{
			var n uint32
			err = binary.Read(r, binary.LittleEndian, &n)
			if err != nil {
				return
			}
			if expected := uint32(0); n != expected {
				err = fmt.Errorf("df2014: book expectation failed: %v != %v", n, expected)
				return
			}
		}

		err = binary.Read(r, binary.LittleEndian, &book.Unk005)
		if err != nil {
			return
		}

		{
			var n int32
			err = binary.Read(r, binary.LittleEndian, &n)
			if err != nil {
				return
			}
			if expected := int32(-1); n != expected {
				err = fmt.Errorf("df2014: book expectation failed: %v != %v", n, expected)
				return
			}
		}

		err = binary.Read(r, binary.LittleEndian, &book.Unk006)
		if err != nil {
			return
		}

		if book.Unk006&2 == 2 {
			book.Unk007 = new([74]byte)
			_, err = io.ReadFull(r, (*book.Unk007)[:])
			if err != nil {
				return
			}
		}
		if book.Unk006&1 == 1 {
			book.Unk008 = new([40]byte)
			_, err = io.ReadFull(r, (*book.Unk008)[:])
			if err != nil {
				return
			}
		}

		err = r.Decode(&book.Unk009)
		if err != nil {
			return
		}

		if book.Unk000&0x800 == 0x800 {
			book.Unk010 = new([6]byte)
			_, err = io.ReadFull(r, (*book.Unk010)[:])
			if err != nil {
				return
			}
		}

		l = append(l, book)
	}
}
