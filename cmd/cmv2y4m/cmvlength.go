// +build ignore

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/BenLubar/df2014/cmv"
)

func main() {
	flag.Parse()

	var total time.Duration

	for _, name := range flag.Args() {
		fmt.Println(name)
		total += length(name)
	}

	fmt.Println("total:", total)
}

func length(name string) time.Duration {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r, err := cmv.NewReader(f)
	if err != nil {
		panic(err)
	}

	i := 0
	for {
		_, err := r.Frame()
		if err != nil {
			if err != io.EOF {
				fmt.Println("warning:", err)
			}
			break
		}
		i++
		fmt.Print(i, "\r")
	}

	t := time.Second / 50 * time.Duration(i)

	fmt.Println(i, "frames (", t, ")")

	return t
}
