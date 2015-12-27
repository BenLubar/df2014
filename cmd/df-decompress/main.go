package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BenLubar/df2014"
)

func main() {
	flag.Parse()

	for _, fn := range flag.Args() {
		ProcessFile(fn)
	}
}

func ProcessFile(fn string) {
	fi, err := os.Stat(fn)
	if err != nil {
		fmt.Println(fn, "error:", err)
		return
	}

	if fi.IsDir() {
		fmt.Println(fn, "is a directory")
		matches, err := filepath.Glob(filepath.Join(fn, "*.dat"))
		if err != nil {
			fmt.Println(fn, "error:", err)
			return
		}
		for _, match := range matches {
			ProcessFile(match)
		}
		return
	}

	fmt.Println(fn, "decompressing")
	err = DecompressFile(fn)
	if err != nil {
		fmt.Println(fn, "error:", err)
		_ = os.Remove(fn + ".decompressed")
		return
	}
	err = os.Rename(fn, fn+".compressed")
	if err != nil {
		fmt.Println(fn, "error:", err)
		_ = os.Remove(fn + ".decompressed")
		return
	}
	err = os.Rename(fn+".decompressed", fn)
	if err != nil {
		fmt.Println(fn, "error:", err)
		_ = os.Rename(fn+".compressed", fn)
		_ = os.Remove(fn + ".decompressed")
		return
	}
	err = os.Remove(fn + ".compressed")
	if err != nil {
		fmt.Println(fn, "error:", err)
		return
	}
	fmt.Println(fn, "done")
}

func DecompressFile(fn string) (err error) {
	in, err := os.Open(fn)
	if err != nil {
		return
	}
	defer func() {
		e := in.Close()
		if err == nil {
			err = e
		}
	}()

	r := &df2014.Reader{in}

	var h df2014.Header
	err = r.DecodeSimple(&h)
	if err != nil {
		return
	}

	if h.Compression != df2014.ZLib {
		err = fmt.Errorf("Unexpected compression type: %v", h.Compression)
		return
	}

	h.Compression = df2014.Uncompressed

	out, err := os.Create(fn + ".decompressed")
	if err != nil {
		return
	}
	defer func() {
		e := out.Close()
		if err == nil {
			err = e
		}
	}()

	err = binary.Write(out, binary.LittleEndian, &h)
	if err != nil {
		return
	}

	_, err = io.Copy(out, r)
	return
}
