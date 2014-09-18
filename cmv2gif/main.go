package main

import (
	"bytes"
	"errors"
	"flag"
	"github.com/BenLubar/df2014"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "image/png"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var (
	flagTileset = flag.String("t", "", "path to a tileset")
	flagInput   = flag.String("i", "input.cmv", "path to a cmv file")
	flagOutput  = flag.String("o", "output.gif", "path to write the output")
)

func main() {
	flag.Parse()

	tileset, err := NewTilesetFromFile(*flagTileset)
	if err != nil {
		log.Fatal(err)
	}

	var movie df2014.CMV
	{
		f, err := os.Open(*flagInput)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		err = (&df2014.Reader{f}).Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
	}

	output := &gif.GIF{LoopCount: -1}
	output.Delay = make([]int, len(movie.Frames))
	for i := range output.Delay {
		output.Delay[i] = int(movie.Header.FrameTime / (time.Second / 10))
	}

	cols, rows := int(movie.Header.Columns), int(movie.Header.Rows)
	frameSize := image.Rect(0, 0, tileset.size.X*cols, tileset.size.Y*rows)
	tileSize := image.Rect(0, 0, tileset.size.X, tileset.size.Y)

	var wg sync.WaitGroup
	wg.Add(len(movie.Frames))

	render := func(i int, frame df2014.CMVFrame, out **image.Paletted) {
		defer wg.Done()

		img := image.NewPaletted(frameSize, palette.WebSafe)

		for x, col := range frame.Attributes {
			for y, attr := range col {
				tile := tileSize.Add(image.Point{tileset.size.X * x, tileset.size.Y * y})
				draw.Draw(img, tile, tileset.Bg(attr), image.ZP, draw.Src)
				fg := tileset.Fg(frame.Characters[x][y], attr)
				draw.Draw(img, tile, fg, fg.Bounds().Min, draw.Over)
			}
		}

		*out = img
	}

	output.Image = make([]*image.Paletted, len(movie.Frames))
	for i, frame := range movie.Frames {
		go render(i, frame, &output.Image[i])
	}
	wg.Wait()

	f, err := os.Create(*flagOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = gif.EncodeAll(f, output)
	if err != nil {
		log.Fatal(err)
	}
}

var ErrInvalidDimensions = errors.New("image dimensions are invalid for a tileset")

var colors = map[df2014.CMVColor]color.RGBA{
	df2014.ColorBlack:    color.RGBA{0, 0, 0, 255},
	df2014.ColorDGray:    color.RGBA{128, 128, 128, 255},
	df2014.ColorBlue:     color.RGBA{0, 0, 128, 255},
	df2014.ColorLBlue:    color.RGBA{0, 0, 255, 255},
	df2014.ColorGreen:    color.RGBA{0, 128, 0, 255},
	df2014.ColorLGreen:   color.RGBA{0, 255, 0, 255},
	df2014.ColorCyan:     color.RGBA{0, 128, 128, 255},
	df2014.ColorLCyan:    color.RGBA{0, 255, 255, 255},
	df2014.ColorRed:      color.RGBA{128, 0, 0, 255},
	df2014.ColorLRed:     color.RGBA{255, 0, 0, 255},
	df2014.ColorMagenta:  color.RGBA{128, 0, 128, 255},
	df2014.ColorLMagenta: color.RGBA{255, 0, 255, 255},
	df2014.ColorBrown:    color.RGBA{128, 128, 0, 255},
	df2014.ColorYellow:   color.RGBA{255, 255, 0, 255},
	df2014.ColorLGray:    color.RGBA{192, 192, 192, 255},
	df2014.ColorWhite:    color.RGBA{255, 255, 255, 255},
}

type Tileset struct {
	size image.Point
	tile image.Rectangle
	set  map[df2014.CMVColor]*image.RGBA
}

func NewTilesetFromFile(filename string) (*Tileset, error) {
	var r io.Reader

	if filename == "" {
		r = bytes.NewReader(Curses800x600Png)
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}

	return NewTileset(r)
}

func NewTileset(r io.Reader) (*Tileset, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	if img.Bounds().Empty() || img.Bounds().Dx()%16 != 0 || img.Bounds().Dy()%16 != 0 {
		return nil, ErrInvalidDimensions
	}

	var t Tileset

	t.size = image.Point{img.Bounds().Dx() / 16, img.Bounds().Dy() / 16}
	t.tile = image.Rectangle{image.ZP, t.size}.Add(img.Bounds().Min)

	mask := image.NewAlpha16(img.Bounds())
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			if r, g, b, _ := img.At(x, y).RGBA(); r == 0xffff && g == 0 && b == 0xffff {
				mask.SetAlpha16(x, y, color.Transparent)
			} else {
				mask.SetAlpha16(x, y, color.Opaque)
			}
		}
	}

	base := image.NewRGBA(img.Bounds())
	draw.DrawMask(base, img.Bounds(), img, image.ZP, mask, image.ZP, draw.Src)

	t.set = make(map[df2014.CMVColor]*image.RGBA, len(colors))
	for k, v := range colors {

		colorized := image.NewRGBA(base.Bounds())
		for x := base.Bounds().Min.X; x < base.Bounds().Max.X; x++ {
			for y := base.Bounds().Min.Y; y < base.Bounds().Max.Y; y++ {
				colorized.Set(x, y, MultipliedColor{v, base.At(x, y)})
			}
		}

		t.set[k] = colorized
	}

	return &t, nil
}

func (t *Tileset) Fg(char df2014.CMVCharacter, attr df2014.CMVAttribute) image.Image {
	return t.set[attr.Fg()].SubImage(t.tile.Add(image.Point{t.size.X * int(char.Byte()&0xf), t.size.Y * int(char.Byte()>>4)}))
}

func (t *Tileset) Bg(attr df2014.CMVAttribute) image.Image {
	return image.NewUniform(colors[attr.Bg()])
}

type MultipliedColor struct {
	A, B color.Color
}

func (c MultipliedColor) RGBA() (r, g, b, a uint32) {
	r0, g0, b0, a0 := c.A.RGBA()
	r1, g1, b1, a1 := c.B.RGBA()

	return r0 * r1 / 0xffff, g0 * g1 / 0xffff, b0 * b1 / 0xffff, a0 * a1 / 0xffff
}
