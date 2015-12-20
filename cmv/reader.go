package cmv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/BenLubar/df2014/cp437"
	"github.com/nsf/termbox-go"
)

type Reader struct {
	Header
	*Sounds
	n int
	r io.Reader
}

func NewReader(r io.Reader) (*Reader, error) {
	var cmv Reader

	err := binary.Read(r, binary.LittleEndian, &cmv.Header)
	if err != nil {
		return nil, err
	}

	if cmv.Version < 10000 || cmv.Version > 10001 {
		return nil, fmt.Errorf("cmv: unhandled version %d", cmv.Version)
	}

	if cmv.Version >= 10001 {
		cmv.Sounds = new(Sounds)

		var n uint32
		err = binary.Read(r, binary.LittleEndian, &n)
		if err != nil {
			return nil, err
		}

		cmv.Sounds.Files = make([]string, n)

		var buf [50]byte
		for i := range cmv.Sounds.Files {
			_, err = io.ReadFull(r, buf[:])
			if err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				return nil, err
			}
			cmv.Sounds.Files[i] = string(buf[:bytes.IndexByte(buf[:], 0)])
		}

		err = binary.Read(r, binary.LittleEndian, &cmv.Sounds.Timing)
		if err != nil {
			return nil, err
		}
	}

	cmv.r = NewCompression1Reader(r)

	return &cmv, nil
}

type Header struct {
	Version    uint32
	Width      uint32
	Height     uint32
	FrameTicks uint32
}

func (h *Header) FrameTime() time.Duration {
	raw := h.FrameTicks
	if raw == 0 {
		raw = 2
	}
	return time.Duration(raw) * time.Second / 100
}

type Sounds struct {
	Files  []string
	Timing [200][16]uint32
}

func (r *Reader) Frame() (*Frame, error) {
	buf := make([]byte, r.Width*r.Height*2)
	_, err := io.ReadFull(r.r, buf)
	if err != nil {
		return nil, err
	}
	f := &Frame{
		r: r,
		n: r.n,
		b: buf,
	}
	r.n++
	return f, nil
}

type Frame struct {
	r *Reader
	n int
	b []byte
}

func (f *Frame) Width() int  { return int(f.r.Width) }
func (f *Frame) Height() int { return int(f.r.Height) }

func (f *Frame) Rune(x, y int) rune { return cp437.Rune(f.b[x*f.Height()+y]) }

func (f *Frame) Equal(o *Frame) bool { return f.r == o.r && bytes.Equal(f.b, o.b) }

var termboxColors = [...]termbox.Attribute{
	termbox.ColorBlack,
	termbox.ColorBlue,
	termbox.ColorGreen,
	termbox.ColorCyan,
	termbox.ColorRed,
	termbox.ColorMagenta,
	termbox.ColorYellow,
	termbox.ColorWhite,
}

func (f *Frame) attr(x, y int) uint8 {
	return f.b[f.Width()*f.Height()+x*f.Height()+y]
}

func (f *Frame) Fg(x, y int) termbox.Attribute {
	a := f.attr(x, y)
	if a&(1<<6) != 0 {
		return termboxColors[a&7] | termbox.AttrBold
	}
	return termboxColors[a&7]
}
func (f *Frame) Bg(x, y int) termbox.Attribute {
	return termboxColors[(f.attr(x, y)>>3)&7]
}
