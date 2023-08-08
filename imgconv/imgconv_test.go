package imgconv_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/shima8823/image-conversion/imgconv"
)

type mockConverterSuccess struct{}

func (c *mockConverterSuccess) Convert(path string) error {
	return nil
}

type mockConverterFail struct{}

func (c *mockConverterFail) Convert(path string) error {
	return errors.New("mock converter error")
}

func TestWalkJpg(t *testing.T) {
	t.Cleanup(func() {
		cleanupTestData(t)
	})
	t.Parallel()
	tests := []struct {
		name      string
		converter imgconv.ImageConverter
		root      string
		wantErr   bool
	}{
		{
			name:      "successful conversion",
			converter: &mockConverterSuccess{},
			root:      "../testdata/sub_dir",
			wantErr:   false,
		},
		{
			name:      "failed conversion",
			converter: &mockConverterFail{},
			root:      "./test_directory",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := imgconv.WalkJpg(tt.root, tt.converter)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalkJpg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewImageFile(t *testing.T) {
	t.Parallel()
	t.Cleanup(func() {
		cleanupTestData(t)
	})
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid jpg file",
			path:    "../testdata/rainbow.jpg",
			wantErr: false,
		},
		{
			name:    "invalid jpg file",
			path:    "../testdata/error2.jpg",
			wantErr: true,
		},
		{
			name:    "invalid path",
			path:    "../testdata/invalid.jpg",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := imgconv.NewImageFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewImageFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertToPng(t *testing.T) {
	t.Parallel()
	t.Cleanup(func() {
		cleanupTestData(t)
	})
	tests := []struct {
		name    string
		imgFile *imgconv.ImageFile
		wantErr bool
	}{
		{
			name: "valid image file",
			imgFile: &imgconv.ImageFile{
				Path: "../testdata/hitode.jpg",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt.imgFile, _ = imgconv.NewImageFile(tt.imgFile.Path)
		t.Run(tt.name, func(t *testing.T) {
			err := imgconv.ConvertToPng(tt.imgFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToPng() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFileNameWithoutExt(t *testing.T) {
	t.Parallel()
	t.Cleanup(func() {
		cleanupTestData(t)
	})
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "JPEG file", path: "/path/to/file.jpg", want: "/path/to/file"},
		{name: "PNG file", path: "/path/to/file.png", want: "/path/to/file"},
		{name: "No extension", path: "/path/to/file", want: "/path/to/file"},
		{name: "Empty string", path: "", want: ""},
		{name: "Multiple extensions", path: "/path/to/file.tar.gz", want: "/path/to/file.tar"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := imgconv.GetFileNameWithoutExt(tt.path); got != tt.want {
				t.Errorf("GetFileNameWithoutExt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func cleanupTestData(t *testing.T) {
	t.Helper()
	filepath.Walk("../testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ディレクトリは無視
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".png" {
			return os.Remove(path)
		}

		return nil
	})
}
