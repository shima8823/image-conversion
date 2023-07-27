package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "error: invalid argument")
		os.Exit(1)
	}
	err := walkJpg(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func walkJpg(root string) error {
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".jpg" {
				jpgFile, err := os.Open(path)
				if handleError(path, err) != nil {
					return nil
				}
				defer jpgFile.Close()

				img, err := jpeg.Decode(jpgFile)
				if handleError(path, err) != nil {
					return nil
				}

				pngFile, err := os.Create(getFileNameWithoutExt(path) + ".png")
				if handleError(path, err) != nil {
					return nil
				}
				defer pngFile.Close()

				err = png.Encode(pngFile, img)
				if handleError(path, err) != nil {
					return nil
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

func getFileNameWithoutExt(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

func handleError(path string, err error) error {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", path, "is not a valid file")
		return err
	}
	return nil
}
