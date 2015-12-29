package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"

	"github.com/BenLubar/df2014/cmv"
	"github.com/BenLubar/df2014/cp437"
	"github.com/nsf/termbox-go"
)

type Key struct {
	Index  int
	Width  int
	Height int
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

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var h cmv.Header

	err = binary.Read(f, binary.LittleEndian, &h)
	if err != nil {
		panic(err)
	}

	var sounds []byte

	if h.Version >= 10001 {
		var n uint32
		err = binary.Read(f, binary.LittleEndian, &n)
		if err != nil {
			panic(err)
		}

		sounds = make([]byte, 4+n*50+200*16*4)
		binary.LittleEndian.PutUint32(sounds[:4], n)

		_, err = io.ReadFull(f, sounds[4:])
		if err != nil {
			panic(err)
		}
	}

	rest, err := ioutil.ReadAll(cmv.NewCompression1Reader(f))
	if err != nil {
		panic(err)
	}

	size := &Key{
		Index:  0,
		Width:  int(h.Width),
		Height: int(h.Height),
	}
	sizes := []*Key{size}
	frame := 0
	index := 0

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		for y := 0; y < size.Height; y++ {
			for x := 0; x < size.Width; x++ {
				i := index + size.Height*x + y
				j := i + size.Width*size.Height

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
		termbox.Flush()

		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventError:
			panic(e.Err)

		case termbox.EventKey:
			if e.Ch == 0 {
				switch e.Key {
				case termbox.KeyEsc:
					f, err := os.Create("last_record.cmv")
					if err != nil {
						panic(err)
					}
					defer f.Close()
					save(f, h, sounds, rest, sizes)
					return

				case termbox.KeyHome:
					for i := 0; i < 10000; i++ {
						if frame == size.Index {
							break
						}
						frame--
						index -= size.Width * size.Height * 2
					}
				case termbox.KeyPgup:
					for i := 0; i < 100; i++ {
						if frame == size.Index {
							break
						}
						frame--
						index -= size.Width * size.Height * 2
					}
				case termbox.KeyInsert:
					if frame != size.Index {
						frame--
						index -= size.Width * size.Height * 2
					}
				case termbox.KeyDelete:
					frame++
					index += size.Width * size.Height * 2
				case termbox.KeyPgdn:
					frame += 100
					index += 100 * size.Width * size.Height * 2
				case termbox.KeyEnd:
					frame += 10000
					index += 10000 * size.Width * size.Height * 2

				case termbox.KeyArrowLeft:
					if size.Index != frame {
						size = &Key{
							Index:  frame,
							Width:  size.Width,
							Height: size.Height,
						}
						sizes = append(sizes, size)
					}
					if size.Width != 1 {
						size.Width--
					}
				case termbox.KeyArrowRight:
					if size.Index != frame {
						size = &Key{
							Index:  frame,
							Width:  size.Width,
							Height: size.Height,
						}
						sizes = append(sizes, size)
					}
					size.Width++
				case termbox.KeyArrowUp:
					if size.Index != frame {
						size = &Key{
							Index:  frame,
							Width:  size.Width,
							Height: size.Height,
						}
						sizes = append(sizes, size)
					}
					if size.Height != 1 {
						size.Height--
					}
				case termbox.KeyArrowDown:
					if size.Index != frame {
						size = &Key{
							Index:  frame,
							Width:  size.Width,
							Height: size.Height,
						}
						sizes = append(sizes, size)
					}
					size.Height++

				case termbox.KeyBackspace2:
					if len(sizes) != 1 {
						count := frame - size.Index
						index -= count * size.Width * size.Height
						sizes = sizes[:len(sizes)-1]
						size = sizes[len(sizes)-1]
						index += count * size.Width * size.Height
					}
				}
			}
		}
	}
	_ = sizes
}

func save(w io.Writer, h cmv.Header, sounds, rest []byte, sizes []*Key) error {
	for _, size := range sizes {
		if size.Width > int(h.Width) {
			h.Width = uint32(size.Width)
		}
		if size.Height > int(h.Height) {
			h.Height = uint32(size.Height)
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
	var key *Key
	for len(rest) != 0 {
		frame++
		if len(sizes) != 0 && sizes[0].Index == frame {
			key = sizes[0]
			sizes = sizes[1:]
		}

		var last bool

		for z := 0; z < 2; z++ {
			for x := 0; x < key.Width; x++ {
				for y := 0; y < key.Height; y++ {
					i := key.Height*x + y

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
				for y := key.Height; y < int(h.Height); y++ {
					err = bw.WriteByte(0)
					if err != nil {
						return err
					}
				}
			}
			for x := key.Width; x < int(h.Width); x++ {
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
				rest = rest[key.Width*key.Height:]
			}
		}
	}

	return bw.Flush()
}
