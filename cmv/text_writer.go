package cmv

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/BenLubar/df2014/cp437"
	"github.com/BenLubar/df2014/wtf23a"
)

func writeStringList(w io.Writer, l []string, f func(int, byte) byte) error {
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
	return writeStringList(NewCompression1Writer(w), l, nil)
}

// WriteStringListIndex writes the format used by Dwarf Fortress's index file.
func WriteStringListIndex(w io.Writer, l []string) error {
	return writeStringList(NewCompression1Writer(w), l, func(i int, b byte) byte {
		return ^b - byte(i%5)
	})
}

// WriteStringListWTF23a writes the format used by Dwarf Fortress before 40d.
func WriteStringListWTF23a(w io.Writer, l []string, header func() wtf23a.Header) error {
	return writeStringList(wtf23a.NewWriter(w, header), l, nil)
}
