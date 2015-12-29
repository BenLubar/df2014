package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/BenLubar/df2014/cmv"
	"github.com/nsf/termbox-go"
)

func main() {
	flag.Parse()

	for _, fn := range flag.Args() {
		if flag.NArg() != 1 {
			fmt.Println("==>", fn, "<==")
		}
		if err := process(fn); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func process(fn string) (err error) {
	f, err := os.Open(fn)
	if err != nil {
		return
	}
	defer func() {
		if e := f.Close(); err == nil {
			err = e
		}
	}()

	r, err := cmv.NewReader(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	index := 0
	prevTitle, prevBody, prevBox := "", "", ""
	for {
		var frame *cmv.Frame
		frame, err = r.Frame()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return fmt.Errorf("%s (near frame %d)\n", err, index)
		}
		index++

		if isTextViewer(frame) {
			title, body := transcribeTextViewer(frame)
			if title != prevTitle || body != prevBody {
				fmt.Println(time.Duration(index) * r.FrameTime())
				fmt.Println(title)
				fmt.Println(body)
				fmt.Println()
				prevTitle, prevBody = title, body
			}
		} else {
			prevTitle, prevBody = "", ""
		}

		if x1, y1, x2, y2, ok := isMegaBox(frame); ok {
			box := transcribeMegaBox(frame, x1, y1, x2, y2)
			if box != prevBox {
				fmt.Println(time.Duration(index) * r.FrameTime())
				fmt.Println(box)
				fmt.Println()
				prevBox = box
			}
		} else {
			prevBox = ""
		}
	}

	fmt.Println(time.Duration(index) * r.FrameTime())
	return
}

func isTextViewer(f *cmv.Frame) bool {
	cols, rows := f.Width()-1, f.Height()-1
	for y := 0; y <= rows-1; y++ {
		// left side
		if f.Rune(0, y) != '█' ||
			f.Fg(0, y) != termbox.ColorBlack|termbox.AttrBold ||
			f.Bg(0, y) != termbox.ColorBlack {
			return false
		}

		// right side
		if f.Rune(cols, y) != '█' ||
			f.Fg(cols, y) != termbox.ColorBlack|termbox.AttrBold ||
			f.Bg(cols, y) != termbox.ColorBlack {
			return false
		}
	}

	rec := f.Rune(cols-2, rows) == 'R'
	startedTitle, endedTitle := false, false
	for x := 0; x <= cols; x++ {
		// bottom
		if rec && x > cols-3 {
			if f.Rune(x, rows) != []rune{'R', 'E', 'C'}[x-cols+2] ||
				f.Fg(x, rows) != termbox.ColorRed|termbox.AttrBold ||
				f.Bg(x, rows) != termbox.ColorRed {
				return false
			}
		} else if f.Rune(x, rows) != '█' ||
			f.Fg(x, rows) != termbox.ColorBlack|termbox.AttrBold ||
			f.Bg(x, rows) != termbox.ColorBlack {
			return false
		}

		// top
		if f.Rune(x, 0) != '█' ||
			f.Fg(x, 0) != termbox.ColorBlack|termbox.AttrBold ||
			f.Bg(x, 0) != termbox.ColorBlack {
			if endedTitle {
				return false
			}
			startedTitle = true
			if f.Rune(x, 0) == 0 ||
				f.Fg(x, 0) != termbox.ColorBlack ||
				f.Bg(x, 0) != termbox.ColorWhite {
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
			if (f.Fg(x, y) != termbox.ColorWhite &&
				f.Fg(x, y) != termbox.ColorBlack) ||
				f.Bg(x, y) != termbox.ColorBlack {
				return false
			}
		}
	}

	return true
}

func transcribeTextViewer(f *cmv.Frame) (title, body string) {
	cols, rows := f.Width()-1, f.Height()-1

	var buf []rune

	for x := 0; x <= cols; x++ {
		if ch := f.Rune(x, 0); ch != '█' && ch != 0 {
			buf = append(buf, ch)
		}
	}
	title = strings.TrimSpace(string(buf))

	buf = buf[:0]
	for y := 1; y <= rows-1; y++ {
		var last int
		for x := 1; x <= cols-1; x++ {
			if ch := f.Rune(x, y); ch != 0 {
				if ch != ' ' || len(buf) == 0 || buf[len(buf)-1] != ' ' {
					buf = append(buf, ch)
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
	Fg, Bg termbox.Attribute
}{
	{'P', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'r', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'e', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'s', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'s', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{' ', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'E', termbox.ColorGreen | termbox.AttrBold, termbox.ColorBlack},
	{'n', termbox.ColorGreen | termbox.AttrBold, termbox.ColorBlack},
	{'t', termbox.ColorGreen | termbox.AttrBold, termbox.ColorBlack},
	{'e', termbox.ColorGreen | termbox.AttrBold, termbox.ColorBlack},
	{'r', termbox.ColorGreen | termbox.AttrBold, termbox.ColorBlack},
	{' ', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'t', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'o', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{' ', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'c', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'l', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'o', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'s', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'e', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{' ', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'w', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'i', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'n', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'d', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'o', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
	{'w', termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack},
}

func isMegaBox(f *cmv.Frame) (x1, y1, x2, y2 int, ok bool) {
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
			if f.Rune(x2, y) != '║' ||
				f.Fg(x2, y) != termbox.ColorBlack|termbox.AttrBold ||
				f.Bg(x2, y) != termbox.ColorBlack {
				return false
			}
		}

		// bottom
		for x := x1 + 1; x < x2; x++ {
			if i := x - x1 - 2; i >= 0 && i < len(megaBoxBottom) {
				if f.Rune(x, y2) != megaBoxBottom[i].Rune ||
					f.Fg(x, y2) != megaBoxBottom[i].Fg ||
					f.Bg(x, y2) != megaBoxBottom[i].Bg {
					return false
				}
			} else {
				if f.Rune(x, y2) != '═' ||
					f.Fg(x, y2) != termbox.ColorBlack|termbox.AttrBold ||
					f.Bg(x, y2) != termbox.ColorBlack {
					return false
				}
			}
		}

		return true
	}

	cols, rows := f.Width()-1, f.Height()-1

	for x1 = 0; x1 <= cols; x1++ {
		for y1 = 0; y1 <= rows; y1++ {
			if f.Rune(x1, y1) != '╔' ||
				f.Fg(x1, y1) != termbox.ColorBlack|termbox.AttrBold ||
				f.Bg(x1, y1) != termbox.ColorBlack {
				continue
			}

			// Found a possible top left corner.
			// Try to find the bottom left corner.
			for y2 = y1 + 1; y2 <= rows; y2++ {
				if f.Fg(x1, y2) != termbox.ColorBlack|termbox.AttrBold ||
					f.Bg(x1, y2) != termbox.ColorBlack {
					break
				}
				if f.Rune(x1, y2) == '║' {
					continue
				}
				if f.Rune(x1, y2) != '╚' {
					break
				}

				// Found a possible bottom left corner.
				// Try to find the top right corner.
				for x2 = x1 + 1; x2 <= cols; x2++ {
					if f.Fg(x2, y1) != termbox.ColorBlack|termbox.AttrBold ||
						f.Bg(x2, y1) != termbox.ColorBlack {
						break
					}
					if f.Rune(x2, y1) == '═' {
						continue
					}
					if f.Rune(x2, y1) != '╗' {
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

func transcribeMegaBox(f *cmv.Frame, x1, y1, x2, y2 int) string {
	var box []rune
	for y := y1 + 1; y < y2; y++ {
		for x := x1 + 1; x < x2; x++ {
			if ch := f.Rune(x, y); ch != 0 && (ch != ' ' || len(box) == 0 || box[len(box)-1] != ' ') {
				box = append(box, ch)
			}
		}
	}
	return strings.TrimSpace(string(box))
}
