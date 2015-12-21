package cmv

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/BenLubar/df2014/cp437"
)

func writeStringList(w io.Writer, l []string, f func(int, byte) byte) error {
	bw := bufio.NewWriter(NewCompression1Writer(w))

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

		err = binary.Write(bw, binary.LittleEndian, uint16(len(b)))
		if err != nil {
			return err
		}

		n, err := bw.Write(b)
		if err != nil {
			return err
		}
		if n != len(b) {
			return io.ErrShortWrite
		}
	}

	return bw.Flush()
}

// WriteStringList writes the format used by Dwarf Fortress's announcement,
// dipscript, and help files.
func WriteStringList(w io.Writer, l []string) error {
	return writeStringList(w, l, nil)
}

// WriteStringListIndex writes the format used by Dwarf Fortress's index file.
func WriteStringListIndex(w io.Writer, l []string) error {
	return writeStringList(w, l, func(i int, b byte) byte {
		return ^b - byte(i%5)
	})
}
