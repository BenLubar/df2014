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
	_ "image/png"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	flagBuffer  = flag.Int("b", 0, "number of frames to go ahead")
	flagSpeed   = flag.Float64("s", 1, "speed multiplier")
	flagTileset = flag.String("t", "", "path to a tileset")
	flagInput   = flag.String("i", "input.cmv", "path to a cmv file")
	flagOutput  = flag.String("o", "output.gif", "path to write the output")
	Palette     = palette.WebSafe
)

func main() {
	flag.Parse()

	switch len(flag.Args()) {
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
		flag.PrintDefaults()
		os.Exit(1)
	}

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

	delay := int(float64(movie.Header.FrameTime/(time.Second/10)) / *flagSpeed)

	cols, rows := int(movie.Header.Columns), int(movie.Header.Rows)
	frameSize := image.Rect(0, 0, tileset.size.X*cols, tileset.size.Y*rows)
	tileSize := image.Rect(0, 0, tileset.size.X, tileset.size.Y)

	frames := make(chan *image.Paletted, *flagBuffer)

	go func() {
		for _, frame := range movie.Frames {
			img := image.NewPaletted(frameSize, Palette)

			for x, col := range frame.Attributes {
				for y, attr := range col {
					rect := tileSize.Add(image.Point{tileset.size.X * x, tileset.size.Y * y})
					tile := tileset.Tile(frame.Characters[x][y], attr)
					fastDraw(img, rect, tile, tile.Bounds().Min)
				}
			}
			frames <- img
		}

		close(frames)
	}()

	f, err := os.Create(*flagOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = EncodeAll(f, frames, delay)
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
	set  [1 << 7]*image.Paletted
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

	for attr := range t.set {
		colorized := image.NewPaletted(base.Bounds(), Palette)
		for x := base.Bounds().Min.X; x < base.Bounds().Max.X; x++ {
			for y := base.Bounds().Min.Y; y < base.Bounds().Max.Y; y++ {
				colorized.Set(x, y, TileColor{base.At(x, y), colors[df2014.CMVAttribute(attr).Fg()], colors[df2014.CMVAttribute(attr).Bg()]})
			}
		}

		t.set[attr] = colorized
	}

	return &t, nil
}

func (t *Tileset) Tile(char df2014.CMVCharacter, attr df2014.CMVAttribute) *image.Paletted {
	return t.set[attr].SubImage(t.tile.Add(image.Point{t.size.X * int(char.Byte()&0xf), t.size.Y * int(char.Byte()>>4)})).(*image.Paletted)
}

type TileColor struct {
	Base, Fg, Bg color.Color
}

func (c TileColor) RGBA() (r, g, b, a uint32) {
	r0, g0, b0, a0 := c.Base.RGBA()
	r1, g1, b1, a1 := c.Fg.RGBA()
	r2, g2, b2, a2 := c.Bg.RGBA()

	a3 := 0xffff - a0*a1/0xffff

	r = (r0*r1/0xffff + r2*a3/0xffff)
	g = (g0*g1/0xffff + g2*a3/0xffff)
	b = (b0*b1/0xffff + b2*a3/0xffff)
	a = (a0*a1/0xffff + a2*a3/0xffff)

	return
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
