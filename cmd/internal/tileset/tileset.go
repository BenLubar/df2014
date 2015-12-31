package tileset

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"io"
	"log"
	"os"

	"github.com/nsf/termbox-go"
)

var ErrInvalidDimensions = errors.New("tileset: image dimensions must be divisible by 16")

var colorToInt = map[termbox.Attribute]int{
	termbox.ColorBlack:                      0,
	termbox.ColorBlue:                       1,
	termbox.ColorGreen:                      2,
	termbox.ColorCyan:                       3,
	termbox.ColorRed:                        4,
	termbox.ColorMagenta:                    5,
	termbox.ColorYellow:                     6,
	termbox.ColorWhite:                      7,
	termbox.ColorBlack | termbox.AttrBold:   8,
	termbox.ColorBlue | termbox.AttrBold:    9,
	termbox.ColorGreen | termbox.AttrBold:   10,
	termbox.ColorCyan | termbox.AttrBold:    11,
	termbox.ColorRed | termbox.AttrBold:     12,
	termbox.ColorMagenta | termbox.AttrBold: 13,
	termbox.ColorYellow | termbox.AttrBold:  14,
	termbox.ColorWhite | termbox.AttrBold:   15,
}

var colors = [...]color.RGBA{
	{0, 0, 0, 255},
	{0, 0, 128, 255},
	{0, 128, 0, 255},
	{0, 128, 128, 255},
	{128, 0, 0, 255},
	{128, 0, 128, 255},
	{128, 128, 0, 255},
	{192, 192, 192, 255},
	{128, 128, 128, 255},
	{0, 0, 255, 255},
	{0, 255, 0, 255},
	{0, 255, 255, 255},
	{255, 0, 0, 255},
	{255, 0, 255, 255},
	{255, 255, 0, 255},
	{255, 255, 255, 255},
}

type Tileset struct {
	size    image.Point
	tile    image.Rectangle
	set     [len(colors)][len(colors) / 2]*image.Paletted
	palette color.Palette
}

func NewTilesetFromFile(filename string) (*Tileset, error) {
	if filename == "" {
		return NewTileset(bytes.NewReader(curses800x600Png))
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewTileset(f)
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

	base := image.NewRGBA(img.Bounds())
	draw.DrawMask(base, img.Bounds(), img, image.ZP, mask, image.ZP, draw.Src)

	var colorized [len(t.set)][len(t.set[0])]*image.RGBA

	palette := make(map[color.RGBA]bool, 256)

	log.Println("loading tileset... computing palette")

	for fg := range colorized {
		for bg := range colorized[fg] {
			colorized[fg][bg] = image.NewRGBA(base.Bounds())
			for x := base.Bounds().Min.X; x < base.Bounds().Max.X; x++ {
				for y := base.Bounds().Min.Y; y < base.Bounds().Max.Y; y++ {
					colorized[fg][bg].SetRGBA(x, y, tileColor(base.RGBAAt(x, y), colors[fg], colors[bg]))
					if c := colorized[fg][bg].RGBAAt(x, y); !palette[c] {
						t.palette = append(t.palette, c)
						palette[c] = true
					}
				}
			}
		}
	}

	for fg := range colorized {
		for bg, c := range colorized[fg] {
			converted := image.NewPaletted(c.Rect, t.palette)
			draw.Draw(converted, converted.Rect, c, c.Rect.Min, draw.Src)
			t.set[fg][bg] = converted
		}
	}

	return &t, nil
}

func (t *Tileset) Tile(ch byte, fg, bg termbox.Attribute) *image.Paletted {
	b := int(ch)
	lower := b & 0xf
	upper := b >> 4
	pt := image.Pt(t.size.X*lower, t.size.Y*upper)
	sheet := t.set[colorToInt[fg]][colorToInt[bg]]
	return sheet.SubImage(t.tile.Add(pt)).(*image.Paletted)
}

func (t *Tileset) Size() image.Point { return t.size }

func (t *Tileset) Palette() color.Palette { return t.palette }

func tileColor(base, fg, bg color.RGBA) color.RGBA {
	a := 0xff - uint16(base.A)*uint16(fg.A)/0xff

	return color.RGBA{
		R: uint8(uint16(base.R)*uint16(fg.R)/0xff + uint16(bg.R)*a/0xff),
		G: uint8(uint16(base.G)*uint16(fg.G)/0xff + uint16(bg.G)*a/0xff),
		B: uint8(uint16(base.B)*uint16(fg.B)/0xff + uint16(bg.B)*a/0xff),
		A: uint8(uint16(base.A)*uint16(fg.A)/0xff + uint16(bg.A)*a/0xff),
	}
}
