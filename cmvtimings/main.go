package main

import (
	"flag"
	"fmt"
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
		fmt.Println(err)
		return
	}
	// f is closed by StreamCMV

	cmv, err := df2014.StreamCMV(f, 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	frameIndex := 0
	prevTitle, prevBody, prevBox := "", "", ""
	for frame := range cmv.Frames {
		frameIndex++

		if isTextViewer(frame) {
			title, body := transcribeTextViewer(frame)
			if title != prevTitle || body != prevBody {
				fmt.Println(time.Millisecond * 20 * time.Duration(frameIndex))
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
				fmt.Println(time.Millisecond * 20 * time.Duration(frameIndex))
				fmt.Println(box)
				fmt.Println()
				prevBox = box
			}
		} else {
			prevBox = ""
		}
	}

	fmt.Println(time.Millisecond * 20 * time.Duration(frameIndex))
}

func isTextViewer(frame df2014.CMVFrame) bool {
	// left side
	x := 0
	for y, c := range frame.Characters[x] {
		a := frame.Attributes[x][y]
		if c.Rune() != '█' || a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
			return false
		}
	}
	// right side
	x = len(frame.Characters) - 1
	for y, c := range frame.Characters[x][:len(frame.Characters[x])-1] {
		a := frame.Attributes[x][y]
		if c.Rune() != '█' || a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
			return false
		}
	}
	// bottom
	y := len(frame.Characters[0]) - 1
	rec := frame.Characters[len(frame.Characters)-3][y].Rune() == 'R'
	for x, col := range frame.Characters {
		c := col[y]
		a := frame.Attributes[x][y]
		if rec && x >= len(frame.Characters)-3 {
			if c.Rune() != []rune{'R', 'E', 'C'}[x-(len(frame.Characters)-3)] || a.Fg() != df2014.ColorLRed || a.Bg() != df2014.ColorRed {
				return false
			}
		} else if c.Rune() != '█' || a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
			return false
		}
	}
	// top
	y = 0
	startedTitle, endedTitle := false, false
	for x, col := range frame.Characters {
		c := col[y]
		a := frame.Attributes[x][y]
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
	for x, col := range frame.Attributes {
		if x == 0 || x == len(frame.Attributes)-1 {
			continue
		}
		for y, a := range col {
			if y == 0 || y == len(col)-1 {
				continue
			}

			if (a.Fg() != df2014.ColorLGray && a.Fg() != df2014.ColorBlack) || a.Bg() != df2014.ColorBlack {
				return false
			}
		}
	}
	return true
}

func transcribeTextViewer(frame df2014.CMVFrame) (string, string) {
	var title, body []rune

	for _, col := range frame.Characters {
		c := col[0]
		if c.Rune() != '█' && c.Rune() != 0 {
			title = append(title, c.Rune())
		}
	}

	for y := 1; y < len(frame.Characters[0])-1; y++ {
		var last int
		for x, col := range frame.Characters[1 : len(frame.Characters)-1] {
			c := col[y]

			if c.Rune() != 0 {
				if c.Rune() != ' ' || len(body) == 0 || body[len(body)-1] != ' ' {
					body = append(body, c.Rune())
				}
				last = x
			}
		}
		if last != len(frame.Characters)-4 {
			body = append(body, '\n')
		} else if len(body) != 0 && body[len(body)-1] != ' ' {
			body = append(body, ' ')
		}
	}

	return strings.TrimSpace(string(title)), strings.TrimSpace(string(body))
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

func isMegaBox(frame df2014.CMVFrame) (x1, y1, x2, y2 int, ok bool) {
	// ╔══════════════════════════════════╗ - border is fg:dgray bg:black
	// ║                                  ║ - inside is bg:black
	// ║                                  ║ - Press, to, close, window are
	// ║                                  ║   fg:white bg:black
	// ║                                  ║ - Enter is fg:lgreen bg:black
	// ║                                  ║ - spaces are bg:black
	// ╚═Press Enter to close window══════╝ - visible body is one fg color
	validate := func() bool {
		// don't check the left or top sides - they're already valid.

		// right side
		for y := y1 + 1; y < y2; y++ {
			c := frame.Characters[x2][y]
			a := frame.Attributes[x2][y]
			if c.Rune() != '║' || a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
				return false
			}
		}

		// bottom
		for x := x1 + 1; x < x2; x++ {
			c := frame.Characters[x][y2]
			a := frame.Attributes[x][y2]
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

	for x1 = range frame.Characters {
		for y1 = range frame.Characters[x1] {
			c := frame.Characters[x1][y1]
			a := frame.Attributes[x1][y1]
			if c.Rune() != '╔' || a.Fg() != df2014.ColorDGray || a.Bg() != df2014.ColorBlack {
				continue
			}
			// Found a possible top left corner.
			// Try to find the bottom left corner.
			for y2 = range frame.Characters[x1] {
				if y2 <= y1 {
					continue
				}
				c = frame.Characters[x1][y2]
				a = frame.Attributes[x1][y2]

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
				for x2 = range frame.Characters {
					if x2 <= x1 {
						continue
					}
					c = frame.Characters[x2][y1]
					a = frame.Attributes[x2][y1]

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
func transcribeMegaBox(frame df2014.CMVFrame, x1, y1, x2, y2 int) string {
	var box []rune
	for y := y1 + 1; y < y2; y++ {
		for _, col := range frame.Characters[x1+1 : x2] {
			c := col[y]

			if c.Rune() != 0 && (c.Rune() != ' ' || len(box) == 0 || box[len(box)-1] != ' ') {
				box = append(box, c.Rune())
			}
		}
	}
	return strings.TrimSpace(string(box))
}
