package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/BenLubar/df2014"
)

func main() {
	flag.Parse()

	for _, fn := range flag.Args() {
		if flag.NArg() != 1 {
			fmt.Println("==>", fn, "<==")
		}
		process(fn)
	}
}

func process(fn string) {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer f.Close()

	header, _, r, err := df2014.RawCMV(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	frame := &Frame{Header: header}
	prevTitle, prevBody, prevBox := "", "", ""
	for {
		if _, err := frame.ReadFrom(r); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "%s (near frame %d)\n", err, frame.Index)
			return
		}

		if frame.IsTextViewer() {
			title, body := frame.TranscribeTextViewer()
			if title != prevTitle || body != prevBody {
				fmt.Println(frame.Timestamp())
				fmt.Println(title)
				fmt.Println(body)
				fmt.Println()
				prevTitle, prevBody = title, body
			}
		} else {
			prevTitle, prevBody = "", ""
		}

		if x1, y1, x2, y2, ok := frame.IsMegaBox(); ok {
			box := frame.TranscribeMegaBox(x1, y1, x2, y2)
			if box != prevBox {
				fmt.Println(frame.Timestamp())
				fmt.Println(box)
				fmt.Println()
				prevBox = box
			}
		} else {
			prevBox = ""
		}
	}

	fmt.Println(frame.Timestamp())
}

type Frame struct {
	Index  int
	Header df2014.CMVHeader
	Data   []byte
}

func (f *Frame) Timestamp() time.Duration {
	return f.Header.FrameTime * time.Duration(f.Index)
}

func (f *Frame) ReadFrom(r io.Reader) (n int64, err error) {
	if f.Data == nil {
		f.Data = make([]byte, f.Header.Columns*f.Header.Rows*2)
	}
	n_, err := io.ReadFull(r, f.Data)
	n = int64(n_)
	if err == nil {
		f.Index++
	}
	return
}

func (f *Frame) Character(x, y int) df2014.CMVCharacter {
	return df2014.CMVCharacter(f.Data[x*int(f.Header.Rows)+y])
}

func (f *Frame) Attribute(x, y int) df2014.CMVAttribute {
	return df2014.CMVAttribute(f.Data[(int(f.Header.Columns)+x)*int(f.Header.Rows)+y])
}

func (f *Frame) IsTextViewer() bool {
	cols, rows := int(f.Header.Columns)-1, int(f.Header.Rows)-1
	for y := 0; y <= rows-1; y++ {
		// left side
		c := f.Character(0, y)
		a := f.Attribute(0, y)
		if c.Rune() != '█' || a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
			return false
		}

		// right side
		c = f.Character(cols, y)
		a = f.Attribute(cols, y)
		if c.Rune() != '█' || a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
			return false
		}
	}

	rec := f.Character(cols-2, rows).Rune() == 'R'
	startedTitle, endedTitle := false, false
	for x := 0; x <= cols; x++ {
		// bottom
		c := f.Character(x, rows)
		a := f.Attribute(x, rows)
		if rec && x > cols-3 {
			if c.Rune() != []rune{'R', 'E', 'C'}[x-cols+2] || a.Fg() != df2014.ColorLRed || a.Bg() != df2014.ColorRed {
				return false
			}
		} else if c.Rune() != '█' || a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
			return false
		}

		// top
		c = f.Character(x, 0)
		a = f.Attribute(x, 0)
		if c.Rune() != '█' || a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
			if endedTitle {
				return false
			}
			startedTitle = true
			if c.Rune() == 0 || a.Fg() != df2014.ColorBlack || a.Bg() != df2014.ColorLGray {
				return false
			}
		} else if startedTitle {
			endedTitle = true
		}
	}
	if !startedTitle || !endedTitle {
		return false
	}

	// body
	for x := 1; x <= cols-1; x++ {
		for y := 1; y <= rows-1; y++ {
			a := f.Attribute(x, y)
			if (a.Fg() != df2014.ColorLGray && a.Fg() != df2014.ColorBlack) || a.Bg() != df2014.ColorBlack {
				return false
			}
		}
	}

	return true
}

func (f *Frame) TranscribeTextViewer() (title string, body string) {
	cols, rows := int(f.Header.Columns)-1, int(f.Header.Rows)-1

	var buf []rune

	buf = buf[:0]
	for x := 0; x <= cols; x++ {
		c := f.Character(x, 0)
		if c.Rune() != '█' && c.Rune() != 0 {
			buf = append(buf, c.Rune())
		}
	}
	title = strings.TrimSpace(string(buf))

	buf = buf[:0]
	for y := 1; y <= rows-1; y++ {
		var last int
		for x := 1; x <= cols-1; x++ {
			c := f.Character(x, y)

			if c.Rune() != 0 {
				if c.Rune() != ' ' || len(buf) == 0 || buf[len(buf)-1] != ' ' {
					buf = append(buf, c.Rune())
				}
				last = x
			}
		}
		if last != cols-2 {
			buf = append(buf, '\n')
		} else if len(buf) != 0 && buf[len(buf)-1] != ' ' {
			buf = append(buf, ' ')
		}
	}
	body = strings.TrimSpace(string(buf))

	return
}

var megaBoxBottom = []struct {
	Rune   rune
	Fg, Bg df2014.CMVColor
}{
	{'P', df2014.ColorWhite, df2014.ColorBlack},
	{'r', df2014.ColorWhite, df2014.ColorBlack},
	{'e', df2014.ColorWhite, df2014.ColorBlack},
	{'s', df2014.ColorWhite, df2014.ColorBlack},
	{'s', df2014.ColorWhite, df2014.ColorBlack},
	{' ', df2014.ColorWhite, df2014.ColorBlack},
	{'E', df2014.ColorLGreen, df2014.ColorBlack},
	{'n', df2014.ColorLGreen, df2014.ColorBlack},
	{'t', df2014.ColorLGreen, df2014.ColorBlack},
	{'e', df2014.ColorLGreen, df2014.ColorBlack},
	{'r', df2014.ColorLGreen, df2014.ColorBlack},
	{' ', df2014.ColorWhite, df2014.ColorBlack},
	{'t', df2014.ColorWhite, df2014.ColorBlack},
	{'o', df2014.ColorWhite, df2014.ColorBlack},
	{' ', df2014.ColorWhite, df2014.ColorBlack},
	{'c', df2014.ColorWhite, df2014.ColorBlack},
	{'l', df2014.ColorWhite, df2014.ColorBlack},
	{'o', df2014.ColorWhite, df2014.ColorBlack},
	{'s', df2014.ColorWhite, df2014.ColorBlack},
	{'e', df2014.ColorWhite, df2014.ColorBlack},
	{' ', df2014.ColorWhite, df2014.ColorBlack},
	{'w', df2014.ColorWhite, df2014.ColorBlack},
	{'i', df2014.ColorWhite, df2014.ColorBlack},
	{'n', df2014.ColorWhite, df2014.ColorBlack},
	{'d', df2014.ColorWhite, df2014.ColorBlack},
	{'o', df2014.ColorWhite, df2014.ColorBlack},
	{'w', df2014.ColorWhite, df2014.ColorBlack},
}

func (f *Frame) IsMegaBox() (x1, y1, x2, y2 int, ok bool) {
	// ╔══════════════════════════════════╗ - border is fg:dgray bg:black
	// ║                                  ║ - inside is bg:black
	// ║                                  ║ - Press, to, close, window are
	// ║                                  ║   fg:white bg:black
	// ║                                  ║ - Enter is fg:lgreen bg:black
	// ║                                  ║ - spaces are bg:black
	// ╚═Press Enter to close window══════╝ - visible body is one fg color
	validate := func() bool {
		// don't check the left or top sides - they're already valid.

		if x2-x1 < len(megaBoxBottom) {
			return false
		}

		// right side
		for y := y1 + 1; y < y2; y++ {
			c := f.Character(x2, y)
			a := f.Attribute(x2, y)
			if c.Rune() != '║' || a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
				return false
			}
		}

		// bottom
		for x := x1 + 1; x < x2; x++ {
			c := f.Character(x, y2)
			a := f.Attribute(x, y2)
			if i := x - x1 - 2; i >= 0 && i < len(megaBoxBottom) {
				if c.Rune() != megaBoxBottom[i].Rune || a.Fg() != megaBoxBottom[i].Fg || a.Bg() != megaBoxBottom[i].Bg {
					return false
				}
			} else {
				if c.Rune() != '═' || a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
					return false
				}
			}
		}

		return true
	}

	cols, rows := int(f.Header.Columns)-1, int(f.Header.Rows)-1

	for x1 = 0; x1 <= cols; x1++ {
		for y1 = 0; y1 <= rows; y1++ {
			c := f.Character(x1, y1)
			a := f.Attribute(x1, y1)
			if c.Rune() != '╔' || a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
				continue
			}

			// Found a possible top left corner.
			// Try to find the bottom left corner.
			for y2 = y1 + 1; y2 <= rows; y2++ {
				c = f.Character(x1, y2)
				a = f.Attribute(x1, y2)

				if a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
					break
				}
				if c.Rune() == '║' {
					continue
				}
				if c.Rune() != '╚' {
					break
				}

				// Found a possible bottom left corner.
				// Try to find the top right corner.
				for x2 = x1 + 1; x2 <= cols; x2++ {
					c = f.Character(x2, y1)
					a = f.Attribute(x2, y1)

					if a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
						break
					}
					if c.Rune() == '═' {
						continue
					}
					if c.Rune() != '╗' {
						break
					}

					if validate() {
						ok = true
						return
					}
				}
			}
		}
	}
	return
}

func (f *Frame) TranscribeMegaBox(x1, y1, x2, y2 int) string {
	var box []rune
	for y := y1 + 1; y < y2; y++ {
		for x := x1 + 1; x < x2; x++ {
			c := f.Character(x, y)

			if c.Rune() != 0 && (c.Rune() != ' ' || len(box) == 0 || box[len(box)-1] != ' ') {
				box = append(box, c.Rune())
			}
		}
	}
	return strings.TrimSpace(string(box))
}
