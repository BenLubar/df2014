package main

import (
	"image"
	"image/color"
)

type YCbCr struct {
	*image.YCbCr
}

func (p YCbCr) Set(x, y int, c color.Color) {
	p.SetYCbCr(x, y, p.ColorModel().Convert(c).(color.YCbCr))
}

func (p YCbCr) SetYCbCr(x, y int, c color.YCbCr) {
	if !image.Pt(x, y).In(p.Rect) {
		return
	}

	yi := p.YOffset(x, y)
	ci := p.COffset(x, y)

	p.Y[yi] = c.Y
	p.Cb[ci] = c.Cb
	p.Cr[ci] = c.Cr
}
