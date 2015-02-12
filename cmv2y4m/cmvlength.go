// +build ignore

package main

import (
	"fmt"
	"github.com/BenLubar/df2014"
	"os"
	"time"
)

func main() {
	cmv, err := df2014.StreamCMV(os.Stdin, 0)
	if err != nil {
		panic(err)
	}

	i := 0
	for range cmv.Frames {
		i++
		fmt.Print(i, "\r")
	}

	fmt.Println(i, "frames (", time.Second/50*time.Duration(i), ")")
}
