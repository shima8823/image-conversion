package main

import (
	"flag"
	"fmt"
	"github.com/shima8823/image-conversion/imagefile"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "error: invalid argument")
		os.Exit(1)
	}
	err := imagefile.WalkJpg(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
