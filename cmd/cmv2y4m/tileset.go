package main

import (
	"bytes"
	"errors"
	"github.com/BenLubar/df2014"
	"image"
	"image/color"
	"image/draw"
	"io"
	"log"
	"os"
)

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
	set  [1 << 7]*image.YCbCr
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

	t.size = image.Pt(img.Bounds().Dx()/16, img.Bounds().Dy()/16)
	t.tile = image.Rectangle{image.ZP, t.size}.Add(img.Bounds().Min)

	log.Println("loading tileset... clearing alpha")

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

	log.Println("loading tileset... preparing tiles")

	base := image.NewRGBA(img.Bounds())
	draw.DrawMask(base, img.Bounds(), img, image.ZP, mask, image.ZP, draw.Src)

	for attr := range t.set {
		t.set[attr] = image.NewYCbCr(base.Bounds(), image.YCbCrSubsampleRatio444)
		tc := TileColor{
			Fg: colors[df2014.CMVAttribute(attr).Fg()],
			Bg: colors[df2014.CMVAttribute(attr).Bg()],
		}
		for x := base.Bounds().Min.X; x < base.Bounds().Max.X; x++ {
			for y := base.Bounds().Min.Y; y < base.Bounds().Max.Y; y++ {
				tc.Base = base.At(x, y)
				YCbCr{t.set[attr]}.Set(x, y, tc)
			}
		}
	}

	return &t, nil
}

func (t *Tileset) Tile(char df2014.CMVCharacter, attr df2014.CMVAttribute) *image.YCbCr {
	return t.set[attr].SubImage(t.tile.Add(image.Pt(t.size.X*int(char.Byte()&0xf), t.size.Y*int(char.Byte()>>4)))).(*image.YCbCr)
}
