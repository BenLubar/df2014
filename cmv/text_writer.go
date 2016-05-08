package cmv

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/BenLubar/df2014/cp437"
	"github.com/BenLubar/df2014/wtf23a"
)

func writeStringList(w io.Writer, l []string, f func(int, byte) byte, lengthTwice, nullTerminate bool) error {
	bw := bufio.NewWriter(w)

	err := binary.Write(bw, binary.LittleEndian, uint32(len(l)))
	if err != nil {
		return err
	}

	for _, s := range l {
		b := cp437.Bytes(s)

		if f != nil {
			for i := range b {
				b[i] = f(i, b[i])
			}
		}

		err = binary.Write(bw, binary.LittleEndian, uint32(len(b)))
		if err != nil {
			return err
		}

		if lengthTwice {
			err = binary.Write(bw, binary.LittleEndian, uint16(len(b)))
			if err != nil {
				return err
			}
		}

		n, err := bw.Write(b)
		if err != nil {
			return err
		}
		if n != len(b) {
			return io.ErrShortWrite
		}

		if nullTerminate {
			if f != nil {
				err = bw.WriteByte(f(len(b), 0))
			} else {
				err = bw.WriteByte(0)
			}
			if err != nil {
				return err
			}
		}
	}

	return bw.Flush()
}

// WriteStringList writes the format used by Dwarf Fortress's announcement,
// dipscript, and help files.
func WriteStringList(w io.Writer, l []string) error {
	return writeStringList(NewCompression1Writer(w), l, nil, true, false)
}

// WriteStringListIndex writes the format used by Dwarf Fortress's index file.
func WriteStringListIndex(w io.Writer, l []string) error {
	return writeStringList(NewCompression1Writer(w), l, func(i int, b byte) byte {
		return ^b - byte(i%5)
	}, true, false)
}

// WriteStringListWTF23a writes the format used by Dwarf Fortress 23a.
func WriteStringListWTF23a(w io.Writer, l []string, header func() wtf23a.Header) error {
	return writeStringList(wtf23a.NewWriter(w, header), l, nil, false, true)
}

// WriteStringList40d writes the format used by Dwarf Fortress 40d.
func WriteStringList40d(w io.Writer, l []string) error {
	return writeStringList(NewCompression1Writer(w), l, nil, false, true)
}

// WriteStringListIndex40d writes the format used by Dwarf Fortress 40d.
func WriteStringListIndex40d(w io.Writer, l []string) error {
	return writeStringList(NewCompression1Writer(w), l, func(i int, b byte) byte {
		return ^b - byte(i%5)
	}, false, true)
}
