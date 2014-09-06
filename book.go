package df2014

import (
	"encoding/binary"
	"fmt"
)

type Book struct {
	Unk000 uint32
	ID     uint32
	Unk002 uint16
	Unk003 int32
	Unk004 uint16
	Unk005 uint32
	Unk006 uint32
	Unk007 *BookUnk007
	Unk008 *BookUnk008
	Unk009 string
	Unk010 *BookUnk010
}

type BookUnk007 struct {
	Unk000 uint32
	Unk001 uint32
	Unk002 uint16
	Unk003 int32
	Unk004 int32
	Unk005 uint32
	Unk006 uint32
	Unk007 uint16
	Unk008 uint32
	Unk009 uint32
	Unk010 uint16
	Unk011 int32
	Unk012 int32
	Unk013 uint32
	Unk014 uint32
	Unk015 uint16
	Unk016 int16
	Unk017 int32
	Unk018 int32
	Unk019 uint32
	Unk020 uint32
}

type BookUnk008 struct {
	Unk000 uint32
	Unk001 uint16
	Unk002 uint32
	Unk003 uint32
	Unk004 int32
	Unk005 [5]uint16
	Unk006 uint32
	Unk007 uint32
	Unk008 uint32
}

type BookUnk010 struct {
	Unk000 uint32
	Unk001 uint16
}

func (r *Reader) book() (book Book, err error) {
	var signature uint64
	err = binary.Read(r, binary.LittleEndian, &signature)
	if err != nil {
		return
	}
	if signature != 0x8ad08ad08ad0 {
		err = fmt.Errorf("df2014: invalid book signature: %016x", signature)
		return
	}

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

	err = binary.Read(r, binary.LittleEndian, &book.ID)
	if err != nil {
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
		if expected := book.ID; n != expected {
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
		book.Unk007 = new(BookUnk007)
		err = binary.Read(r, binary.LittleEndian, book.Unk007)
		if err != nil {
			return
		}
		if actual, expected := book.Unk007.Unk000, uint32(7); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk002, uint16(0); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk003, int32(-1); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk004, int32(-1); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk005, uint32(0); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk006, uint32(0); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk007, uint16(0); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk008, uint32(6); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk009, book.Unk007.Unk001; actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk010, uint16(0); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk011, int32(-1); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk012, int32(-1); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk013, uint32(0); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk014, uint32(0); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk015, uint16(0); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk016, int16(-1); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk017, int32(-1); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk018, int32(-1); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk019, uint32(0); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk007.Unk020, uint32(0); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
	}
	if book.Unk006&1 == 1 {
		book.Unk008 = new(BookUnk008)
		err = binary.Read(r, binary.LittleEndian, book.Unk008)
		if err != nil {
			return
		}
		if actual, expected := book.Unk008.Unk000, uint32(9); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, min, max := book.Unk008.Unk001, uint16(36), uint16(38); actual < min || actual > max {
			err = fmt.Errorf("df2014: book expectation failed: %v < %v || %v > %v", actual, min, actual, max)
			return
		}
		if actual, expected := book.Unk008.Unk004, int32(-1); actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
		if actual, expected := book.Unk008.Unk005, [5]uint16{0, 0, 0, 0, 0}; actual != expected {
			err = fmt.Errorf("df2014: book expectation failed: %v != %v", actual, expected)
			return
		}
	}

	err = r.Decode(&book.Unk009)
	if err != nil {
		return
	}

	if book.Unk000&0x800 == 0x800 {
		book.Unk010 = new(BookUnk010)
		err = binary.Read(r, binary.LittleEndian, book.Unk010)
		if err != nil {
			return
		}
	}

	return
}
