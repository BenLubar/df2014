package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BenLubar/commander"
	"github.com/BenLubar/df2014/cmv"
	"github.com/nsf/termbox-go"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: ./cmv2asciicast name.cmv")
		flag.PrintDefaults()
	}

	commander.RegisterFlags(flag.CommandLine)

	flag.Parse()

	if err := commander.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer commander.Close()

	if flag.NArg() == 0 {
		flag.Usage()
	}

	for _, name := range flag.Args() {
		err := convert(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		}
	}
}

func convert(name string) (err error) {
	fin, err := os.Open(name)
	if err != nil {
		return
	}
	defer fin.Close()

	r, err := cmv.NewReader(fin)
	if err != nil {
		return
	}

	w, err := os.Create(strings.TrimSuffix(name, ".cmv") + ".json")
	if err != nil {
		return
	}
	defer w.Close()

	type whoops struct{ error }
	defer func() {
		if r := recover(); r != nil {
			err = r.(whoops).error
		}
	}()
	printf := func(format string, args ...interface{}) {
		_, err := fmt.Fprintf(w, format, args...)
		if err != nil {
			panic(whoops{err})
		}
	}

	jsonName, err := json.Marshal(&name)
	if err != nil {
		panic(err)
	}

	printf(`{"version":1,"width":%d,"height":%d,"command":"Dwarf_Fortress","title":%s,"env":{"TERM":"xterm"},"stdout":[`, r.Width, r.Height, jsonName)
	var prev *cmv.Frame
	frames, skip := 0, 0
	for {
		frame, err := r.Frame()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		frames++
		if prev != nil && prev.Equal(frame) {
			skip++
			continue
		}
		prev = frame
		if frames == 1 {
			printf(`[0,"\u001b[?25l`)
		} else {
			printf(`,[%v,"`, r.FrameTime().Seconds()*float64(skip+1))
		}
		skip = 0
		prevBright, prevFg, prevBg := -1, -1, -1
		for y := 0; y < int(r.Height); y++ {
			printf(`\u001b[%dH`, y+1)
			for x := 0; x < int(r.Width); x++ {
				bright := 22
				fg := int(frame.Fg(x, y))
				bg := int(frame.Bg(x, y))
				if fg&int(termbox.AttrBold) != 0 {
					fg &^= int(termbox.AttrBold)
					bright = 1
				}
				fg--
				bg--

				first := true
				if prevBright != bright {
					if first {
						first = false
						printf(`\u001b[`)
					} else {
						printf(`;`)
					}
					printf(`%d`, bright)
				}
				if prevFg != fg {
					if first {
						first = false
						printf(`\u001b[`)
					} else {
						printf(`;`)
					}
					printf(`3%d`, fg)
				}
				if prevBg != bg {
					if first {
						first = false
						printf(`\u001b[`)
					} else {
						printf(`;`)
					}
					printf(`4%d`, bg)
				}
				if !first {
					printf(`m`)
				}
				prevBright, prevFg, prevBg = bright, fg, bg

				ch := frame.Rune(x, y)
				if ch < ' ' {
					ch = ' '
				}
				if ch == '\\' || ch == '"' {
					printf(`\`)
				}
				printf(`%c`, ch)
			}
		}
		printf("\"]\n")
	}
	if frames != 0 {
		printf(`,[0,"\u001b[?25h\u001b[22;39;49m"]`)
	}
	printf(`],"duration":%v}`, r.FrameTime().Seconds()*float64(frames))
	return nil
}
