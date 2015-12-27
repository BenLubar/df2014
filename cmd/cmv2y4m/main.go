package main

import (
	"bufio"
	"errors"
	"flag"
	"github.com/BenLubar/df2014"
	"github.com/BenLubar/job"
	"image"
	_ "image/png"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	flagBuffer     = flag.Int("b", 0, "number of frames to go ahead")
	flagTileset    = flag.String("t", "", "path to a tileset")
	flagSkipHeader = flag.Bool("skip-header", false, "skip the output header and just encode the frames")
	flagSkipFrames = flag.Int("skip-frames", 0, "")
)

func main() {
	flag.Parse()

	moviech := make(chan *df2014.CMVStream)
	go func() {
		var (
			movie df2014.CMVStream
			err   error
		)
		if flag.NArg() == 0 {
			movie, err = df2014.StreamCMV(os.Stdin, *flagBuffer)
		} else {
			movie, err = CombineCMV(flag.Args()...)
		}
		if err != nil {
			log.Fatal(err)
		}

		log.Println("opened cmv")
		moviech <- &movie
	}()

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

	frameDuration := movie.Header.FrameTime
	delay := int(frameDuration / (10 * time.Millisecond))

	if delay <= 0 {
		delay, frameDuration = 2, 20*time.Millisecond
	}

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

func Multiplex(movie *df2014.CMVStream, completed chan<- *image.YCbCr, tileset *Tileset, frameDuration time.Duration) {
	input := make(chan interface{}, *flagBuffer)

	cols, rows := int(movie.Header.Columns), int(movie.Header.Rows)
	frameSize := image.Rect(0, 0, tileset.size.X*cols, tileset.size.Y*rows)
	tileSize := image.Rect(0, 0, tileset.size.X, tileset.size.Y)

	go func() {
		defer close(input)

		for i := 0; i < *flagSkipFrames; i++ {
			<-movie.Frames
			if i%10000 == 9999 {
				log.Println("skipped", i+1, "/", *flagSkipFrames, "frames")
			}
		}

		for frame := range movie.Frames {
			input <- frame
		}
	}()

	output := job.Work(input, &job.Options{
		Work: func(v interface{}) interface{} {
			in, out := v.(df2014.CMVFrame), image.NewYCbCr(frameSize, image.YCbCrSubsampleRatio444)

			for x, col := range in.Attributes {
				for y, attr := range col {
					rect := tileSize.Add(image.Pt(tileset.size.X*x, tileset.size.Y*y))
					tile := tileset.Tile(in.Characters[x][y], attr)
					fastDraw(out, rect, tile, tile.Bounds().Min)
				}
			}

			return out
		},
		NumWorkers: runtime.GOMAXPROCS(0),
		MaxWaiting: *flagBuffer,
	})

	defer close(completed)

	var lastLog time.Time
	var seq int
	for out := range output {
		seq++
		completed <- out.(*image.YCbCr)
		if time.Since(lastLog) >= time.Second {
			lastLog = time.Now()
			log.Println("encoding frames...", seq, "encoded,", time.Duration(seq)*frameDuration)
		}
	}
	log.Println("finished encoding", seq, "frames,", time.Duration(seq)*frameDuration)
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
