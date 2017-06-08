package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"io"
	"os"

	"github.com/BenLubar/df2014/cmv"
	"github.com/pkg/errors"
)

func main() {
	flagOutput := flag.String("o", "", "output filename (required)")
	flagSkipBlank := flag.Bool("skip-blank", false, "omit frames that are completely black, possibly with an FPS meter")
	flag.Parse()

	if *flagOutput == "" {
		flag.Usage()
		os.Exit(2)
		return
	}

	f, err := os.Create(*flagOutput)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var expected cmv.Header
	var writer *bufio.Writer
	check := func(h cmv.Header, s *cmv.Sounds) (io.Writer, error) {
		if writer == nil {
			expected = h
			if err := writeHeader(f, h, s); err != nil {
				return nil, err
			}
			writer = bufio.NewWriterSize(cmv.NewCompression1Writer(f), int(h.Width*h.Height)*2*200)
			return writer, nil
		}

		if h.Width != expected.Width || h.Height != expected.Height {
			return nil, errors.Errorf("frame size mismatch: %d×%d is not %d×%d", h.Width, h.Height, expected.Width, expected.Height)
		}

		return writer, nil
	}

	isFrameBlank := func(buf []byte) bool {
		half := len(buf) / 2
		seenFPS, inFPS := false, false
		for i := 0; i < half; i++ {
			if (buf[i] == ' ' || buf[half+i]&(1<<6|7) == 0) && (buf[half+i]>>3)&7 == 0 {
				inFPS = false
			} else if buf[i] == 'F' && buf[half+i] == (1<<6|3<<3|7) && !seenFPS {
				seenFPS = true
				inFPS = true
			} else if !inFPS || (buf[i] != ' ' && buf[half+i]&(1<<6|7) != (1<<6|7)) || (buf[half+i]>>3)&7 != 3 {
				return false
			}
		}
		return true
	}
	frame := func(buf []byte) bool {
		if *flagSkipBlank && isFrameBlank(buf) {
			return false
		}
		return true
	}

	for _, name := range flag.Args() {
		err = processFile(check, frame, name)
		if err != nil {
			panic(err)
		}
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}

func processFile(checkHeader func(cmv.Header, *cmv.Sounds) (io.Writer, error), checkFrame func([]byte) bool, name string) (err error) {
	f, err := os.Open(name)
	if err != nil {
		return errors.Wrap(err, "opening cmv")
	}
	defer func() {
		if e := f.Close(); err == nil {
			err = errors.Wrap(e, "closing cmv")
		}
	}()

	var h cmv.Header
	var s *cmv.Sounds
	err = cmv.ReadHeader(f, &h, &s)
	if err != nil {
		return errors.Wrap(err, "reading header")
	}

	w, err := checkHeader(h, s)
	if err != nil {
		return errors.Wrapf(err, "mismatch %q", name)
	}

	r := cmv.NewCompression1Reader(f)

	buf := make([]byte, int(h.Width*h.Height)*2)
	for {
		_, err := io.ReadFull(r, buf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "read frames from %q", name)
		}

		if checkFrame(buf) {
			n, err := w.Write(buf)
			if err == nil && n != len(buf) {
				err = io.ErrShortWrite
			}
			if err != nil {
				return errors.Wrapf(err, "write frames from %q", name)
			}
		}
	}
}

func writeHeader(w io.Writer, h cmv.Header, s *cmv.Sounds) error {
	if h.Version < 10000 || h.Version > 10001 {
		return errors.Errorf("cmv header: unhandled version %d", h.Version)
	}

	err := binary.Write(w, binary.LittleEndian, &h)
	if err != nil {
		return err
	}

	if h.Version >= 10001 {
		err = binary.Write(w, binary.LittleEndian, uint32(len(s.Files)))
		if err != nil {
			return errors.Wrap(err, "cmv sound count")
		}

		for i, f := range s.Files {
			var buf [50]byte
			copy(buf[:], f)
			err = binary.Write(w, binary.LittleEndian, &buf)
			if err != nil {
				return errors.Wrapf(err, "cmv sound name %d", i)
			}
		}

		err = binary.Write(w, binary.LittleEndian, &s.Timing)
		if err != nil {
			return errors.Wrap(err, "cmv sound timings")
		}
	}

	return nil
}
