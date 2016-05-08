package cmv

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/BenLubar/df2014/cp437"
	"github.com/BenLubar/df2014/wtf23a"
)

func readStringList(r io.Reader, f func(int, byte) byte) ([]string, error) {
	var l []string

	var count uint32
	err := binary.Read(r, binary.LittleEndian, &count)
	if err != nil {
		return nil, err
	}

	for i := uint32(0); i < count; i++ {
		var l32 uint32
		err = binary.Read(r, binary.LittleEndian, &l32)
		if err != nil {
			return nil, err
		}

		var l16 uint16
		err = binary.Read(r, binary.LittleEndian, &l16)
		if err != nil {
			return nil, err
		}

		if l32 != uint32(l16) {
			return nil, fmt.Errorf("cmv: invalid string list (%d != %d)", l32, l16)
		}

		b := make([]byte, l16)
		_, err = io.ReadFull(r, b)
		if err != nil {
			return nil, err
		}

		if f != nil {
			for i := range b {
				b[i] = f(i, b[i])
			}
		}

		l = append(l, cp437.String(b))
	}

	return l, nil
}

// ReadStringList reads the format used by Dwarf Fortress's announcement,
// dipscript, and help files.
func ReadStringList(r io.Reader) ([]string, error) {
	return readStringList(NewCompression1Reader(r), nil)
}

// ReadStringListIndex reads the format used by Dwarf Fortress's index file.
func ReadStringListIndex(r io.Reader) ([]string, error) {
	return readStringList(NewCompression1Reader(r), func(i int, b byte) byte {
		return ^b - byte(i%5)
	})
}

// ReadStringListWTF23a reads the format used by Dwarf Fortress before 40d.
func ReadStringListWTF23a(r io.Reader) ([]string, error) {
	return readStringList(wtf23a.NewReader(r), nil)
}
