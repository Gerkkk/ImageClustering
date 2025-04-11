package Images

import (
	"ImageClustering/internal/domain"
	"errors"
	"io"
	"strings"
)

type ImageI interface {
	GetPixels() ([]domain.Pixel, error)
	CreateSimplified(clusters []domain.Pixel, path string) (string, error)
	//createPalette(clusters []domain.Pixel, path string) string
}

func NewImage(fileName string, file io.Reader) (ImageI, error) {
	splittedName := strings.Split(fileName, ".")

	switch splittedName[len(splittedName)-1] {
	case "jpg", "jpeg":
		got, err := NewJpegImage(file)

		if err != nil {
			return nil, err
		}

		return &got, nil

	case "png":
		got, err := NewPngImage(file)
		if err != nil {
			return nil, err
		}
		return &got, nil
	case "gif":
		got, err := NewGifImage(file)
		if err != nil {
			return nil, err
		}
		return &got, nil
	}

	return nil, errors.New("invalid image format")
}
