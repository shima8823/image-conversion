package imgconv

import (
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

type ImageConverter interface {
	Convert(path string) error
}

type JpgToPngConverter struct{}

// WalkJpg walks the file tree rooted at root, calling converter.Convert for each jpg file in the tree.
func WalkJpg(root string, converter ImageConverter) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".jpg" {
			return nil
		}
		return converter.Convert(path)
	})
}

// Convert converts a jpg file to a png file.
func (c *JpgToPngConverter) Convert(path string) error {
	imageFile, err := NewImageFile(path)
	if err != nil {
		HandleError(path)
		return nil
	}
	return ConvertToPng(imageFile)
}

// NewImageFile returns a new ImageFile.
func NewImageFile(path string) (*ImageFile, error) {
	jpgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jpgFile.Close()

	img, err := jpeg.Decode(jpgFile)
	if err != nil {
		return nil, err
	}
	return &ImageFile{Path: path, Img: img}, nil
}

// ConvertToPng converts a jpg file to a png file.
func ConvertToPng(imageFile *ImageFile) error {
	pngFile, err := os.Create(GetFileNameWithoutExt(imageFile.Path) + ".png")
	if err != nil {
		return err
	}
	defer pngFile.Close()

	return png.Encode(pngFile, imageFile.Img)
}

// GetFileNameWithoutExt returns the file name without the extension
func GetFileNameWithoutExt(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

// HandleError prints an error message to stderr
func HandleError(path string) {
	fmt.Fprintln(os.Stderr, "error:", path, "is not a valid file")
}
