package main

import (
	"bufio"
	"fmt"
	"image"
)

func EncodeAll(w *bufio.Writer, frames <-chan *image.YCbCr, delay int) error {
	frame, ok := <-frames
	if !ok {
		return fmt.Errorf("No frames!")
	}

	if !*flagSkipHeader {
		_, err := fmt.Fprintf(w, "YUV4MPEG2 W%d H%d F50:1 Ip A1:1 C444\n", frame.Rect.Dx(), frame.Rect.Dy())
		if err != nil {
			return err
		}
	}

	d := delay
	for ok {
		_, err := w.WriteString("FRAME\n")
		if err != nil {
			return err
		}

		_, err = w.Write(frame.Y)
		if err != nil {
			return err
		}

		_, err = w.Write(frame.Cb)
		if err != nil {
			return err
		}

		_, err = w.Write(frame.Cr)
		if err != nil {
			return err
		}

		d -= 2
		for d < 2 {
			frame, ok = <-frames
			d += delay
		}
	}

	return nil
}
