package main

import (
	"flag"
	"fmt"
	"github.com/shima8823/image-conversion/imgconv"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "error: invalid argument")
		os.Exit(1)
	}

	converter := &imgconv.JpgToPngConverter{}
	err := imgconv.WalkJpg(args[0], converter)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
