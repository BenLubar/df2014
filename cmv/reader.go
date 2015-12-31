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

// Reader reads CMV files.
type Reader struct {
	Header
	*Sounds
	n int
	r io.Reader
}

// NewReader returns a Reader that reads CMV files. If there is an error
// parsing the header, the error will be returned instead.
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

// Header is the CMV header.
type Header struct {
	// Version is either 10000 (0x2710) or 10001 (0x2701). The latter
	// includes Sounds.
	Version uint32
	// Width is the number of columns in each frame of the CMV.
	Width uint32
	// Height is the number of rows in each frame of the CMV.
	Height uint32
	// FrameTicks is the frame rate in hundredths of a second per frame.
	// It can be 0, in which case the default frame rate is used.
	FrameTicks uint32
}

// FrameTime converts FrameTicks to a duration. It uses a default of 50 fps if
// FrameTicks is zero.
func (h *Header) FrameTime() time.Duration {
	raw := h.FrameTicks
	if raw == 0 {
		raw = 2
	}
	return time.Duration(raw) * time.Second / 100
}

// Sounds is the CMV 10001 sounds header, used by the intro videos for Dwarf
// Fortress.
type Sounds struct {
	Files  []string
	Timing [200][16]uint32
}

// Frame reads and returns the next frame. Either the Frame or the error will
// be non-nil. io.EOF signals the end of the CMV.
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

// Frame represents a single frame of a CMV.
type Frame struct {
	r *Reader
	n int
	b []byte
}

// Width returns the Width field of the Reader that created this Frame.
func (f *Frame) Width() int { return int(f.r.Width) }

// Height returns the Height field of the Reader that created this Frame.
func (f *Frame) Height() int { return int(f.r.Height) }

// Byte returns the CP437-encoded character on tile (x, y), where 0≤x<Width and
// 0≤y<Height.
func (f *Frame) Byte(x, y int) byte { return f.b[x*f.Height()+y] }

// Rune returns the character on tile (x, y), where 0≤x<Width and 0≤y<Height.
func (f *Frame) Rune(x, y int) rune { return cp437.Rune(f.Byte(x, y)) }

// Equal returns true if f and o are from the same reader and have the same
// characters and colors for each tile.
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

// Fg returns the foreground color on tile (x, y), where 0≤x<Width and
// 0≤y<Height. It is guaranteed to be one of the termbox.Color* constants or
// one of the termbox.Color* constants ORed with termbox.AttrBold.
func (f *Frame) Fg(x, y int) termbox.Attribute {
	a := f.attr(x, y)
	if a&(1<<6) != 0 {
		return termboxColors[a&7] | termbox.AttrBold
	}
	return termboxColors[a&7]
}

// Bg returns the background color on tile (x, y), where 0≤x<Width and
// 0≤y<Height. It is guaranteed to be one of the termbox.Color* constants.
func (f *Frame) Bg(x, y int) termbox.Attribute {
	return termboxColors[(f.attr(x, y)>>3)&7]
}
