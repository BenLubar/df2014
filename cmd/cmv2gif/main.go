package main

import (
	"bufio"
	"flag"
	"image"
	_ "image/png"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BenLubar/df2014/cmd/internal/tileset"
	"github.com/BenLubar/df2014/cmv"
)

var (
	flagSpeed   = flag.Float64("s", 1, "speed multiplier")
	flagTileset = flag.String("t", "", "path to a tileset")
	flagInput   = flag.String("i", "input.cmv", "path to a cmv file")
	flagOutput  = flag.String("o", "output.gif", "path to write the output")
)

func main() {
	flag.Parse()

	switch flag.NArg() {
	case 0:
		// do nothing
	case 1:
		if *flagInput == "input.cmv" && *flagOutput == "output.gif" && strings.HasSuffix(flag.Arg(0), ".cmv") {
			*flagInput = flag.Arg(0)
			*flagOutput = strings.TrimSuffix(flag.Arg(0), ".cmv") + ".gif"
			break
		}
		fallthrough
	default:
		flag.Usage()
		os.Exit(1)
	}

	tileset, err := tileset.NewTilesetFromFile(*flagTileset)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("loaded tileset", *flagTileset)

	in, err := os.Open(*flagInput)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	movie, err := cmv.NewReader(in)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("opened cmv", *flagInput)

	delay := int(movie.Header.FrameTime().Seconds() * 100 / *flagSpeed)
	if delay <= 0 {
		delay = 1
	}

	frameDuration := time.Duration(delay) * 10 * time.Millisecond

	log.Println("time per frame:", frameDuration)

	cols, rows := int(movie.Width), int(movie.Height)
	size := tileset.Size()
	frameSize := image.Rect(0, 0, size.X*cols, size.Y*rows)
	tileSize := image.Rect(0, 0, size.X, size.Y)

	frames := make(chan *image.Paletted)

	go func() {
		var lastLog time.Time

		i := 0

		defer func() {
			log.Println("finished encoding", i, "frames,", time.Duration(i)*frameDuration)
			close(frames)
		}()

		for {
			frame, err := movie.Frame()
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				return
			}
			i++

			img := image.NewPaletted(frameSize, tileset.Palette())

			for x := 0; x < frame.Width(); x++ {
				for y := 0; y < frame.Height(); y++ {
					rect := tileSize.Add(image.Pt(size.X*x, size.Y*y))
					tile := tileset.Tile(frame.Byte(x, y), frame.Fg(x, y), frame.Bg(x, y))
					fastDraw(img, rect, tile, tile.Rect.Min)
				}
			}
			frames <- img

			if time.Since(lastLog) >= time.Second {
				lastLog = time.Now()
				log.Println("encoding frames...", i, "encoded,", time.Duration(i)*frameDuration)
			}
		}
	}()

	f, err := os.Create(*flagOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	err = EncodeAll(w, frames, delay)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

// Assumptions:
// dst and src do not overlap
// dst and src have the same palette
// the locations given are valid for both images
//
func fastDraw(dst *image.Paletted, r image.Rectangle, src *image.Paletted, sp image.Point) {
	pix0, stride0 := dst.Pix[dst.PixOffset(r.Min.X, r.Min.Y):], dst.Stride
	pix1, stride1 := src.Pix[src.PixOffset(sp.X, sp.Y):], src.Stride

	dx, dy := r.Dx(), r.Dy()
	for y := 0; y < dy; y++ {
		copy(pix0[stride0*y:stride0*y+dx], pix1[stride1*y:stride1*y+dx])
	}
}
