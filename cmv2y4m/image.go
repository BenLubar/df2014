package main

import (
	"image"
	"image/color"
)

type YCbCr struct {
	*image.YCbCr
}

func (p YCbCr) Set(x, y int, c color.Color) {
	if !image.Pt(x, y).In(p.Rect) {
		return
	}

	cc := p.ColorModel().Convert(c).(color.YCbCr)
	yi := p.YOffset(x, y)
	ci := p.COffset(x, y)

	p.Y[yi] = cc.Y
	p.Cb[ci] = cc.Cb
	p.Cr[ci] = cc.Cr
}
