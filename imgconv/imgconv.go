// package imgconv contains functions for converting jpg files to png files
package imgconv

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

// ImageFile is a struct that contains the path and image.Image of a file.
type ImageFile struct {
	Path string
	Img  image.Image
}

// WalkJpg walks the directory tree rooted at root, converting all jpg files
func WalkJpg(root string) error {
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".jpg" {
				imageFile, err := NewImageFile(&ImageFile{Path: path})
				if err != nil {
					HandleError(path)
					return nil
				}
				err = ConvertToPng(imageFile)
				if err != nil {
					HandleError(path)
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

// NewImageFile returns a pointer to an ImageFile struct
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

// ConvertToPng converts the image.Image to a png file
func ConvertToPng(imageFile *ImageFile) error {
	pngFile, err := os.Create(GetFileNameWithoutExt(imageFile.Path) + ".png")
	if err != nil {
		return err
	}
	defer pngFile.Close()

	err = png.Encode(pngFile, imageFile.Img)
	if err != nil {
		return err
	}
	return nil
}

// GetFileNameWithoutExt returns the file name without the extension
func GetFileNameWithoutExt(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

// HandleError prints an error message to stderr
func HandleError(path string) {
	fmt.Fprintln(os.Stderr, "error:", path, "is not a valid file")
}
