// Modified from go 1.3.1: src/pkg/image/gif/writer.go

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"compress/lzw"
	"errors"
	"image"
	"image/color"
	"io"
)

// Section indicators.
const (
	sExtension       = 0x21
	sImageDescriptor = 0x2C
	sTrailer         = 0x3B
)

// Graphic control extension fields.
const (
	gcLabel     = 0xF9
	gcBlockSize = 0x04
)

var log2Lookup = [8]int{2, 4, 8, 16, 32, 64, 128, 256}

func log2(x int) int {
	for i, v := range log2Lookup {
		if x <= v {
			return i
		}
	}
	return -1
}

// Little-endian.
func writeUint16(b []uint8, u uint16) {
	b[0] = uint8(u)
	b[1] = uint8(u >> 8)
}

// writer is a buffered writer.
type writer interface {
	Flush() error
	io.Writer
	io.ByteWriter
}

// encoder encodes an image to the GIF format.
type encoder struct {
	// w is the writer to write to. err is the first error encountered during
	// writing. All attempted writes after the first error become no-ops.
	w   writer
	err error
	// g is a reference to the data that is being encoded.
	//g *GIF
	// buf is a scratch buffer. It must be at least 768 so we can write the color map.
	buf [1024]byte
}

// blockWriter writes the block structure of GIF image data, which
// comprises (n, (n bytes)) blocks, with 1 <= n <= 255. It is the
// writer given to the LZW encoder, which is thus immune to the
// blocking.
type blockWriter struct {
	e *encoder
}

func (b blockWriter) Write(data []byte) (int, error) {
	if b.e.err != nil {
		return 0, b.e.err
	}
	if len(data) == 0 {
		return 0, nil
	}
	total := 0
	for total < len(data) {
		n := copy(b.e.buf[1:256], data[total:])
		total += n
		b.e.buf[0] = uint8(n)

		n, b.e.err = b.e.w.Write(b.e.buf[:n+1])
		if b.e.err != nil {
			return 0, b.e.err
		}
	}
	return total, b.e.err
}

func (e *encoder) flush() {
	if e.err != nil {
		return
	}
	e.err = e.w.Flush()
}

func (e *encoder) write(p []byte) {
	if e.err != nil {
		return
	}
	_, e.err = e.w.Write(p)
}

func (e *encoder) writeByte(b byte) {
	if e.err != nil {
		return
	}
	e.err = e.w.WriteByte(b)
}

func (e *encoder) writeHeader(pm *image.Paletted) {
	if e.err != nil {
		return
	}

	if len(pm.Palette) == 0 {
		e.err = errors.New("gif: cannot encode image with empty palette")
		return
	}

	_, e.err = io.WriteString(e.w, "GIF89a")
	if e.err != nil {
		return
	}

	//pm := e.g.Image[0]
	// Logical screen width and height.
	writeUint16(e.buf[0:2], uint16(pm.Bounds().Dx()))
	writeUint16(e.buf[2:4], uint16(pm.Bounds().Dy()))
	e.write(e.buf[:4])

	paddedSize := log2(len(pm.Palette)) // Size of Local Color Table: 2^(1+n).
	e.buf[0] = 0x80 | uint8(paddedSize)
	e.buf[1] = 0x00 // Background Color Index.
	e.buf[2] = 0x00 // Pixel Aspect Ratio.
	e.write(e.buf[:3])

	e.writeColorTable(pm.Palette, paddedSize)

	// Always add animation info.
	//if len(e.g.Image) > 1 {
	{
		e.buf[0] = 0x21 // Extension Introducer.
		e.buf[1] = 0xff // Application Label.
		e.buf[2] = 0x0b // Block Size.
		e.write(e.buf[:3])
		_, e.err = io.WriteString(e.w, "NETSCAPE2.0") // Application Identifier.
		if e.err != nil {
			return
		}
		e.buf[0] = 0x03            // Block Size.
		e.buf[1] = 0x01            // Sub-block Index.
		writeUint16(e.buf[2:4], 0) //writeUint16(e.buf[2:4], uint16(e.g.LoopCount))
		e.buf[4] = 0x00            // Block Terminator.
		e.write(e.buf[:5])
	}
}

func (e *encoder) writeColorTable(p color.Palette, size int) {
	if e.err != nil {
		return
	}

	for i := 0; i < log2Lookup[size]; i++ {
		if i < len(p) {
			r, g, b, _ := p[i].RGBA()
			e.buf[3*i+0] = uint8(r >> 8)
			e.buf[3*i+1] = uint8(g >> 8)
			e.buf[3*i+2] = uint8(b >> 8)
		} else {
			// Pad with black.
			e.buf[3*i+0] = 0x00
			e.buf[3*i+1] = 0x00
			e.buf[3*i+2] = 0x00
		}
	}
	e.write(e.buf[:3*log2Lookup[size]])
}

func (e *encoder) writeImageBlock(pm *image.Paletted, delay int) {
	if e.err != nil {
		return
	}

	b := pm.Bounds()
	if b.Dx() >= 1<<16 || b.Dy() >= 1<<16 || b.Min.X < 0 || b.Min.X >= 1<<16 || b.Min.Y < 0 || b.Min.Y >= 1<<16 {
		e.err = errors.New("gif: image block is too large to encode")
		return
	}

	e.buf[0] = sExtension                  // Extension Introducer.
	e.buf[1] = gcLabel                     // Graphic Control Label.
	e.buf[2] = gcBlockSize                 // Block Size.
	e.buf[3] = 0x00                        // No Transparency.
	writeUint16(e.buf[4:6], uint16(delay)) // Delay Time (1/100ths of a second)

	e.buf[6] = 0x00 // Transparent color index.
	e.buf[7] = 0x00 // Block Terminator.
	e.write(e.buf[:8])

	e.buf[0] = sImageDescriptor
	writeUint16(e.buf[1:3], uint16(b.Min.X))
	writeUint16(e.buf[3:5], uint16(b.Min.Y))
	writeUint16(e.buf[5:7], uint16(b.Dx()))
	writeUint16(e.buf[7:9], uint16(b.Dy()))
	e.write(e.buf[:9])

	paddedSize := log2(len(pm.Palette)) // Size of Local Color Table: 2^(1+n).
	// Interlacing is not supported.
	e.writeByte(0x00)

	litWidth := paddedSize + 1
	if litWidth < 2 {
		litWidth = 2
	}
	e.writeByte(uint8(litWidth)) // LZW Minimum Code Size.

	lzww := lzw.NewWriter(blockWriter{e: e}, lzw.LSB, litWidth)
	_, e.err = lzww.Write(pm.Pix)
	if e.err != nil {
		lzww.Close()
		return
	}
	lzww.Close()
	e.writeByte(0x00) // Block Terminator.
}

// EncodeAll writes the images in ch to w in GIF format with the
// given loop count and delay between frames.
func EncodeAll(w io.Writer, ch <-chan *image.Paletted, delay int) error {
	first, ok := <-ch
	if !ok {
		return errors.New("gif: must provide at least one image")
	}

	e := encoder{}
	if ww, ok := w.(writer); ok {
		e.w = ww
	} else {
		e.w = bufio.NewWriter(w)
	}

	e.writeHeader(first)
	e.writeImageBlock(first, delay)
	for pm := range ch {
		e.writeImageBlock(pm, delay)
	}
	e.writeByte(sTrailer)
	e.flush()
	return e.err
}
