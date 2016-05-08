package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/BenLubar/df2014/cmv"
)

func main() {
	index := flag.Bool("i", false, "index mode")
	encode := flag.Bool("e", false, "encode from JSON to DF format")
	decode := flag.Bool("d", false, "decode from DF format to JSON")
	lines := flag.Bool("l", false, "use lines instead of JSON")

	flag.Parse()

	if *encode == *decode {
		flag.Usage()
		os.Exit(1)
	}

	var s []string
	var err error

	if *decode {
		if *index {
			s, err = cmv.ReadStringListIndex(os.Stdin)
		} else {
			s, err = cmv.ReadStringList(os.Stdin)
		}
	} else {
		if *lines {
			in := bufio.NewScanner(os.Stdin)
			for in.Scan() {
				s = append(s, in.Text())
			}
			err = in.Err()
		} else {
			err = json.NewDecoder(os.Stdin).Decode(&s)
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if *encode {
		if *index {
			err = cmv.WriteStringListIndex(os.Stdout, s)
		} else {
			err = cmv.WriteStringList(os.Stdout, s)
		}
	} else {
		if *lines {
			for _, l := range s {
				_, err = fmt.Fprintln(os.Stdout, l)
				if err != nil {
					break
				}
			}
		} else {
			err = json.NewEncoder(os.Stdout).Encode(&s)
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}
