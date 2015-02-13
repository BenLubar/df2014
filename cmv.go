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

func RawCMV(r io.Reader) (header CMVHeader, sounds *CMVSounds, frames io.Reader, err error) {
	in := &Reader{bufio.NewReader(r)}

	err = binary.Read(in, binary.LittleEndian, &header.Version)
	if err != nil {
		return
	}

	if header.Version < 10000 || header.Version > 10001 {
		err = fmt.Errorf("df2014: unhandled version %d", header.Version)
		return
	}

	err = binary.Read(in, binary.LittleEndian, &header.Columns)
	if err != nil {
		return
	}

	err = binary.Read(in, binary.LittleEndian, &header.Rows)
	if err != nil {
		return
	}

	var frameTimeRaw uint32
	err = binary.Read(in, binary.LittleEndian, &frameTimeRaw)
	if err != nil {
		return
	}
	if frameTimeRaw == 0 {
		frameTimeRaw = 2
	}
	header.FrameTime = time.Duration(frameTimeRaw) * time.Second / 100

	if header.Version >= 10001 {
		sounds = new(CMVSounds)

		var soundsRaw [][50]byte
		err = in.DecodeSimple(&soundsRaw)
		if err != nil {
			return
		}

		for _, sound := range soundsRaw {
			sounds.Files = append(sounds.Files, string(sound[:bytes.IndexByte(sound[:], 0)]))
		}

		err = in.DecodeSimple(&sounds.Timing)
		if err != nil {
			return
		}
	}

	in.Reader = &compression1Reader{r: in.Reader}
	frames = in
	return
}

func StreamCMV(in io.ReadCloser, buffer int) (cmv CMVStream, err error) {
	header, sounds, r, err := RawCMV(in)
	if err != nil {
		in.Close()
		return
	}

	frames := make(chan CMVFrame, buffer)
	cmv.Header, cmv.Sounds, cmv.Frames = header, sounds, frames

	size := cmv.Header.Columns * cmv.Header.Rows

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
