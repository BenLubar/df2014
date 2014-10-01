package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"errors"
	"flag"
	"github.com/BenLubar/df2014"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	flagBuffer     = flag.Int("b", 0, "number of frames to go ahead")
	flagTileset    = flag.String("t", "", "path to a tileset")
	flagSkipHeader = flag.Bool("skip-header", false, "skip the output header and just encode the frames")
)

func main() {
	flag.Parse()

	moviech := make(chan *df2014.CMVStream)
	switch flag.NArg() {
	case 0:
		go func() {
			movie, err := df2014.StreamCMV(os.Stdin, *flagBuffer)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("opened cmv")
			moviech <- &movie
		}()
	default:
		go func() {
			movie, err := CombineCMV(flag.Args()...)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("opened cmv")
			moviech <- &movie
		}()
	}

	tilesetch := make(chan *Tileset)
	go func() {
		tileset, err := NewTilesetFromFile(*flagTileset)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("loaded tileset", *flagTileset)
		tilesetch <- tileset
	}()

	tileset, movie := <-tilesetch, <-moviech

	delay := int(movie.Header.FrameTime / (10 * time.Millisecond))
	if delay <= 0 {
		delay = 2
	}

	frameDuration := time.Duration(delay) * 10 * time.Millisecond

	log.Println("time per frame:", frameDuration)

	frames := make(chan *image.YCbCr, *flagBuffer)

	go Multiplex(movie, frames, tileset, frameDuration)

	w := bufio.NewWriter(os.Stdout)

	err := EncodeAll(w, frames, delay)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Flush()
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

func Multiplex(movie *df2014.CMVStream, completed chan<- *image.YCbCr, tileset *Tileset, frameDuration time.Duration) {
	input, output := make(chan *job, *flagBuffer), make(chan *job, *flagBuffer)

	procs := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup

	wg.Add(procs)
	go func() {
		wg.Wait()
		close(output)
	}()

	jobs := new(jobHeap)
	heap.Init(jobs)

	limiter := make(chan struct{}, procs+*flagBuffer)
	go JobCreator(movie.Frames, input, limiter)

	cols, rows := int(movie.Header.Columns), int(movie.Header.Rows)
	frameSize := image.Rect(0, 0, tileset.size.X*cols, tileset.size.Y*rows)
	tileSize := image.Rect(0, 0, tileset.size.X, tileset.size.Y)

	for i := 0; i < procs; i++ {
		go Worker(tileset, input, output, &wg, frameSize, tileSize)
	}

	var lastLog time.Time
	seq, done := 0, false
	var img *image.YCbCr
	for !done || jobs.Len() > 0 {
		var (
			ready <-chan *job
			send  chan<- *image.YCbCr
		)

		if !done {
			ready = output
		}
		if img != nil {
			send = completed
		}

		select {
		case j, ok := <-ready:
			if ok {
				heap.Push(jobs, j)
			} else {
				done = true
			}

		case send <- img:
			<-limiter
			seq++
			img = nil

			if time.Since(lastLog) >= time.Second {
				lastLog = time.Now()
				log.Println("encoding frames...", seq, "encoded,", time.Duration(seq)*frameDuration)
			}
		}

		if img == nil && jobs.Len() > 0 && (*jobs)[0].seq == seq {
			img = heap.Pop(jobs).(*job).out
		}
	}
	log.Println("finished encoding", seq, "frames,", time.Duration(seq)*frameDuration)
	close(completed)
}

func JobCreator(frames <-chan df2014.CMVFrame, input chan<- *job, limiter chan struct{}) {
	i := 0

	for frame := range frames {
		limiter <- struct{}{}
		input <- &job{
			seq: i,
			in:  frame,
		}

		i++
	}
	close(input)
}

func Worker(tileset *Tileset, input <-chan *job, output chan<- *job, wg *sync.WaitGroup, frameSize, tileSize image.Rectangle) {
	for j := range input {
		j.out = image.NewYCbCr(frameSize, image.YCbCrSubsampleRatio444)

		for x, col := range j.in.Attributes {
			for y, attr := range col {
				rect := tileSize.Add(image.Pt(tileset.size.X*x, tileset.size.Y*y))
				tile := tileset.Tile(j.in.Characters[x][y], attr)
				fastDraw(j.out, rect, tile, tile.Bounds().Min)
			}
		}

		output <- j
	}

	wg.Done()
}

type job struct {
	seq int
	in  df2014.CMVFrame
	out *image.YCbCr
}

type jobHeap []*job

func (h *jobHeap) Len() int {
	return len(*h)
}

func (h *jobHeap) Less(i, j int) bool {
	return (*h)[i].seq < (*h)[j].seq
}

func (h *jobHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *jobHeap) Push(x interface{}) {
	*h = append(*h, x.(*job))
}

func (h *jobHeap) Pop() interface{} {
	l := len(*h)
	j := (*h)[l-1]
	*h = (*h)[:l-1]
	return j
}

var ErrHeaderMismatch = errors.New("CMV file headers differ in dimensions, timing, or version")

func CombineCMV(names ...string) (cmv df2014.CMVStream, err error) {
	movies := make([]df2014.CMVStream, len(names))
	for i, fn := range names {
		var f *os.File
		f, err = os.Open(fn)
		if err != nil {
			return
		}

		// f is closed by StreamCMV
		movies[i], err = df2014.StreamCMV(f, *flagBuffer)
		if err != nil {
			return
		}
	}

	cmv.Header = movies[0].Header

	for _, m := range movies[1:] {
		if m.Header != cmv.Header {
			err = ErrHeaderMismatch
			return
		}
	}

	frames := make(chan df2014.CMVFrame, *flagBuffer)

	cmv.Frames = frames

	go func() {
		defer close(frames)
		for _, m := range movies {
			for f := range m.Frames {
				frames <- f
			}
		}
	}()

	return
}
