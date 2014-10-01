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

func fastDraw(dst *image.YCbCr, r image.Rectangle, src *image.YCbCr, sp image.Point) {
	yoff0, ystride0 := dst.YOffset(r.Min.X, r.Min.Y), dst.YStride
	coff0, cstride0 := dst.COffset(r.Min.X, r.Min.Y), dst.CStride
	y0, cb0, cr0 := dst.Y[yoff0:], dst.Cb[coff0:], dst.Cr[coff0:]

	yoff1, ystride1 := src.YOffset(sp.X, sp.Y), src.YStride
	coff1, cstride1 := src.COffset(sp.X, sp.Y), src.CStride
	y1, cb1, cr1 := src.Y[yoff1:], src.Cb[coff1:], src.Cr[coff1:]

	dx, dy := r.Dx(), r.Dy()
	for y := 0; y < dy; y++ {
		copy(y0[ystride0*y:ystride0*y+dx], y1[ystride1*y:ystride1*y+dx])
		copy(cb0[cstride0*y:cstride0*y+dx], cb1[cstride1*y:cstride1*y+dx])
		copy(cr0[cstride0*y:cstride0*y+dx], cr1[cstride1*y:cstride1*y+dx])
	}
}
