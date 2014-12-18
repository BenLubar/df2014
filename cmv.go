package df2014

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type CMVStream struct {
	Header CMVHeader
	Sounds *CMVSounds
	Frames <-chan CMVFrame
}

type CMVHeader struct {
	Version uint32

	Columns   uint32
	Rows      uint32
	FrameTime time.Duration
}

type CMVSounds struct {
	Files  []string
	Timing [200][16]uint32
}

type CMVFrame struct {
	Characters [][]CMVCharacter
	Attributes [][]CMVAttribute
}

type CMVCharacter uint8

func (c CMVCharacter) Byte() byte {
	return byte(c)
}

func (c CMVCharacter) Rune() rune {
	return cp437[c]
}

func (c CMVCharacter) String() string {
	return string(c.Rune())
}

type CMVAttribute uint8

func (c CMVAttribute) Fg() CMVColor {
	return CMVColor((c<<1)&0xE | (c >> 6))
}

func (c CMVAttribute) Bg() CMVColor {
	return CMVColor((c >> 2) & 0xE)
}

type CMVColor uint8

const (
	ColorBlack CMVColor = iota
	ColorDGray
	ColorBlue
	ColorLBlue
	ColorGreen
	ColorLGreen
	ColorCyan
	ColorLCyan
	ColorRed
	ColorLRed
	ColorMagenta
	ColorLMagenta
	ColorBrown
	ColorYellow
	ColorLGray
	ColorWhite
)

func StreamCMV(in io.ReadCloser, buffer int) (cmv CMVStream, err error) {
	r := &Reader{bufio.NewReader(in)}

	err = binary.Read(r, binary.LittleEndian, &cmv.Header.Version)
	if err != nil {
		in.Close()
		return
	}

	if cmv.Header.Version < 10000 || cmv.Header.Version > 10001 {
		err = fmt.Errorf("df2014: unhandled version %d", cmv.Header.Version)
		in.Close()
		return
	}

	err = binary.Read(r, binary.LittleEndian, &cmv.Header.Columns)
	if err != nil {
		in.Close()
		return
	}

	err = binary.Read(r, binary.LittleEndian, &cmv.Header.Rows)
	if err != nil {
		in.Close()
		return
	}

	size := cmv.Header.Columns * cmv.Header.Rows

	var frameTimeRaw uint32
	err = binary.Read(r, binary.LittleEndian, &frameTimeRaw)
	if err != nil {
		in.Close()
		return
	}
	cmv.Header.FrameTime = time.Duration(frameTimeRaw) * time.Second / 100

	if cmv.Header.Version >= 10001 {
		cmv.Sounds = new(CMVSounds)

		var soundsRaw [][50]byte
		err = r.DecodeSimple(&soundsRaw)
		if err != nil {
			in.Close()
			return
		}

		for _, sound := range soundsRaw {
			cmv.Sounds.Files = append(cmv.Sounds.Files, string(sound[:bytes.IndexByte(sound[:], 0)]))
		}

		err = r.DecodeSimple(&cmv.Sounds.Timing)
		if err != nil {
			in.Close()
			return
		}
	}

	r.Reader = &compression1Reader{r: r.Reader}

	frames := make(chan CMVFrame, buffer)
	cmv.Frames = frames

	go func() {
		defer close(frames)
		defer in.Close()

		for {
			buf := make([]byte, size)
			_, err := io.ReadFull(r, buf)
			if err != nil {
				if err == io.EOF {
					return
				}
				// TODO: some sort of error handling
				return
			}

			var frame CMVFrame

			characters := make([]CMVCharacter, size)
			frame.Characters = make([][]CMVCharacter, cmv.Header.Columns)
			for i := range frame.Characters {
				frame.Characters[i] = characters[i*int(cmv.Header.Rows) : (i+1)*int(cmv.Header.Rows)]
			}

			for i, b := range buf {
				characters[i] = CMVCharacter(b)
			}

			_, err = io.ReadFull(r, buf)
			if err != nil {
				// TODO: some sort of error handling
				return
			}

			attributes := make([]CMVAttribute, size)
			frame.Attributes = make([][]CMVAttribute, cmv.Header.Columns)
			for i := range frame.Attributes {
				frame.Attributes[i] = attributes[i*int(cmv.Header.Rows) : (i+1)*int(cmv.Header.Rows)]
			}

			for i, b := range buf {
				attributes[i] = CMVAttribute(b)
			}

			frames <- frame
		}
	}()
	return
}
