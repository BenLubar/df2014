package main

import (
	"bufio"
	"flag"
	"image"
	"image/color"
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
	flagSpeed      = flag.Float64("s", 1, "speed multiplier")
	flagTileset    = flag.String("t", "", "path to a tileset")
	flagInput      = flag.String("i", "input.cmv", "path to a cmv file")
	flagOutput     = flag.String("o", "output.y4m", "path to write the output")
	flagSkipHeader = flag.Bool("skip-header", false, "don't output Y4M header")
)

func main() {
	flag.Parse()

	switch flag.NArg() {
	case 0:
		// do nothing
	case 1:
		if *flagInput == "input.cmv" && *flagOutput == "output.y4m" && strings.HasSuffix(flag.Arg(0), ".cmv") {
			*flagInput = flag.Arg(0)
			*flagOutput = strings.TrimSuffix(flag.Arg(0), ".cmv") + ".y4m"
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

	palette := make([]color.YCbCr, len(tileset.Palette()))
	for i := range palette {
		palette[i] = color.YCbCrModel.Convert(tileset.Palette()[i]).(color.YCbCr)
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

	frames := make(chan *image.YCbCr)

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

			img := image.NewYCbCr(frameSize, image.YCbCrSubsampleRatio444)

			for x := 0; x < frame.Width(); x++ {
				for y := 0; y < frame.Height(); y++ {
					rect := tileSize.Add(image.Pt(size.X*x, size.Y*y))
					tile := tileset.Tile(frame.Byte(x, y), frame.Fg(x, y), frame.Bg(x, y))
					fastDraw(img, rect, tile, tile.Rect.Min, palette)
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

func fastDraw(dst *image.YCbCr, r image.Rectangle, src *image.Paletted, sp image.Point, palette []color.YCbCr) {
	pixY, strideY := dst.Y[dst.YOffset(r.Min.X, r.Min.Y):], dst.YStride
	pixCb, strideCb := dst.Cb[dst.COffset(r.Min.X, r.Min.Y):], dst.CStride
	pixCr, strideCr := dst.Cr[dst.COffset(r.Min.X, r.Min.Y):], dst.CStride
	pixSrc, strideSrc := src.Pix[src.PixOffset(sp.X, sp.Y):], src.Stride

	dx, dy := r.Dx(), r.Dy()
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			c := palette[pixSrc[strideSrc*y+x]]
			pixY[strideY*y+x] = c.Y
			pixCb[strideCb*y+x] = c.Cb
			pixCr[strideCr*y+x] = c.Cr
		}
	}
}
