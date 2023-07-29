package imagefile

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

func getFileNameWithoutExt(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

func handleError(path string) {
	fmt.Fprintln(os.Stderr, "error:", path, "is not a valid file")
}
