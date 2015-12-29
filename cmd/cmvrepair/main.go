package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/BenLubar/df2014/cmv"
	"github.com/BenLubar/df2014/cp437"
	"github.com/nsf/termbox-go"
)

type keyframe struct {
	frame  int
	width  int
	height int
}

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

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "filename.cmv")
	fmt.Fprintln(os.Stderr, "use arrow keys to move the bottom right corner of the current frame.")
	fmt.Fprintln(os.Stderr, "press esc to save to last_record.cmv and exit.")
	fmt.Fprintln(os.Stderr, "delete/page down/end go forward by 1/100/10000 frames.")
	fmt.Fprintln(os.Stderr, "insert/page up/home go back by 1/100/10000 frames.")
	fmt.Fprintln(os.Stderr, "backspace deletes the last keyframe, allowing you to move back if you made a mistake.")
	os.Exit(2)
}

type fatal struct {
	text string
	err  error
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			f := r.(fatal)
			fmt.Fprintln(os.Stderr, f.text, f.err)
			os.Exit(1)
		}
	}()

	if len(os.Args) != 2 {
		usage()
	}

	h, sounds, rest := load(os.Args[1])

	size := &keyframe{
		frame:  0,
		width:  int(h.Width),
		height: int(h.Height),
	}
	sizes := []*keyframe{size}
	frame := 0
	index := 0

	if err := termbox.Init(); err != nil {
		panic(fatal{"error initializing UI:", err})
	}
	defer termbox.Close()

	step := func(diff int) {
		if frame+diff < size.frame {
			diff = size.frame - frame
		}
		frame += diff
		index += diff * size.width * size.height * 2
	}
	move := func(dx, dy int) {
		if size.width+dx < 1 || size.height+dy < 1 {
			return
		}
		if size.frame != frame {
			size = &keyframe{
				frame:  frame,
				width:  size.width,
				height: size.height,
			}
			sizes = append(sizes, size)
		}
		size.width += dx
		size.height += dy
	}

	for {
		render(size, index, rest)

		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventError:
			panic(fatal{"termbox error:", e.Err})

		case termbox.EventKey:
			if e.Ch == 0 {
				switch e.Key {
				case termbox.KeyEsc:
					f, err := os.Create("last_record.cmv")
					if err == nil {
						err = save(f, h, sounds, rest, sizes)
					}
					if err == nil {
						err = f.Close()
					}
					if err != nil {
						panic(fatal{"could not write output:", err})
					}
					return

				case termbox.KeyHome:
					step(-10000)
				case termbox.KeyPgup:
					step(-100)
				case termbox.KeyInsert:
					step(-1)
				case termbox.KeyDelete:
					step(1)
				case termbox.KeyPgdn:
					step(100)
				case termbox.KeyEnd:
					step(10000)

				case termbox.KeyArrowLeft:
					move(-1, 0)
				case termbox.KeyArrowRight:
					move(1, 0)
				case termbox.KeyArrowUp:
					move(0, -1)
				case termbox.KeyArrowDown:
					move(0, 1)

				case termbox.KeyBackspace2:
					if len(sizes) != 1 {
						count := frame - size.frame
						index -= count * size.width * size.height
						sizes = sizes[:len(sizes)-1]
						size = sizes[len(sizes)-1]
						index += count * size.width * size.height
					}
				}
			}
		}
	}
}

func load(name string) (h cmv.Header, sounds, rest []byte) {
	f, err := os.Open(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		usage()
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(fatal{"error closing CMV:", err})
		}
	}()

	err = binary.Read(f, binary.LittleEndian, &h)
	if err != nil {
		panic(fatal{"CMV format error in header:", err})
	}

	if h.Version >= 10001 {
		var n uint32
		err = binary.Read(f, binary.LittleEndian, &n)
		if err != nil {
			panic(fatal{"CMV format error in sound list:", err})
		}

		sounds = make([]byte, 4+n*50+200*16*4)
		binary.LittleEndian.PutUint32(sounds[:4], n)

		_, err = io.ReadFull(f, sounds[4:])
		if err != nil {
			panic(fatal{"CMV format error in sound list:", err})
		}
	}

	rest, err = ioutil.ReadAll(cmv.NewCompression1Reader(f))
	if err != nil {
		panic(fatal{"CMV decompression error:", err})
	}

	return
}

func save(w io.Writer, h cmv.Header, sounds, rest []byte, sizes []*keyframe) error {
	for _, size := range sizes {
		if size.width > int(h.Width) {
			h.Width = uint32(size.width)
		}
		if size.height > int(h.Height) {
			h.Height = uint32(size.height)
		}
	}

	err := binary.Write(w, binary.LittleEndian, &h)
	if err != nil {
		return err
	}
	_, err = w.Write(sounds)
	if err != nil {
		return err
	}

	bw := bufio.NewWriterSize(cmv.NewCompression1Writer(w), int(h.Width*h.Height)*2*200)

	frame := -1
	var key *keyframe
	for len(rest) != 0 {
		frame++
		if len(sizes) != 0 && sizes[0].frame == frame {
			key = sizes[0]
			sizes = sizes[1:]
		}
		if frame%100 == 0 {
			render(key, 0, rest)
		}

		var last bool

		for z := 0; z < 2; z++ {
			for x := 0; x < key.width; x++ {
				for y := 0; y < key.height; y++ {
					i := key.height*x + y

					if i < len(rest) {
						err = bw.WriteByte(rest[i])
					} else {
						err = bw.WriteByte(0)
						last = true
					}
					if err != nil {
						return err
					}
				}
				for y := key.height; y < int(h.Height); y++ {
					err = bw.WriteByte(0)
					if err != nil {
						return err
					}
				}
			}
			for x := key.width; x < int(h.Width); x++ {
				for y := 0; y < int(h.Height); y++ {
					err = bw.WriteByte(0)
					if err != nil {
						return err
					}
				}
			}
			if last {
				rest = nil
			} else {
				rest = rest[key.width*key.height:]
			}
		}
	}

	return bw.Flush()
}

func render(size *keyframe, index int, rest []byte) {
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(fatal{"graphics error:", err})
	}
	for y := 0; y < size.height; y++ {
		for x := 0; x < size.width; x++ {
			i := index + size.height*x + y
			j := i + size.width*size.height

			ch, fg, bg := rune(0), termbox.ColorDefault, termbox.ColorDefault
			if i < len(rest) {
				ch = cp437.Rune(rest[i])
			}
			if j < len(rest) {
				fg = termboxColors[rest[j]&7]
				if rest[j]&(1<<6) != 0 {
					fg |= termbox.AttrBold
				}
				bg = termboxColors[(rest[j]>>3)&7]
			}

			termbox.SetCell(x, y, ch, fg, bg)
		}
	}
	if err := termbox.Flush(); err != nil {
		panic(fatal{"graphics error:", err})
	}
}
