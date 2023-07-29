package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

type ImageFile struct {
	Path string
	Img  image.Image
}

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
				imageFile, err := NewImageFile(&ImageFile{Path: path})
				if err != nil {
					handleError(path)
					return nil
				}
				err = ConvertToPng(imageFile)
				if err != nil {
					handleError(path)
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

func NewImageFile(imageFile *ImageFile) (*ImageFile, error) {
	jpgFile, err := os.Open(imageFile.Path)
	if err != nil {
		return nil, err
	}
	defer jpgFile.Close()

	img, err := jpeg.Decode(jpgFile)
	if err != nil {
		return nil, err
	}
	return &ImageFile{Path: imageFile.Path, Img: img}, nil
}

func ConvertToPng(imageFile *ImageFile) error {
	pngFile, err := os.Create(getFileNameWithoutExt(imageFile.Path) + ".png")
	if err != nil {
		return err
	}
	defer pngFile.Close()
	defer fmt.Println("converted:", pngFile.Name())

	err = png.Encode(pngFile, imageFile.Img)
	if err != nil {
		return err
	}
	return nil
}

func handleError(path string) {
	fmt.Fprintln(os.Stderr, "error:", path, "is not a valid file")
}
