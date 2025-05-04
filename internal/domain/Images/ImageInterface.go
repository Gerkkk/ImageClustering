package Images

import (
	"ImageClustering/internal/domain"
	"errors"
	"strings"
)

type ImageI interface {
	GetPixels() ([]domain.Pixel, error)
	CreateSimplified(clusters []domain.Pixel, path string) (string, error)
	//createPalette(clusters []domain.Pixel, path string) string
}

func NewImage(input ImageConstructorData) (ImageI, error) {
	splittedName := strings.Split(input.FileName, ".")

	switch splittedName[len(splittedName)-1] {
	case "jpg", "jpeg":
		got, err := NewJpegImage(input.File)

		if err != nil {
			return nil, err
		}

		return &got, nil

	case "png":
		got, err := NewPngImage(input.File)
		if err != nil {
			return nil, err
		}
		return &got, nil
	case "gif":
		got, err := NewGifImage(input.File, input.FramesCount)
		if err != nil {
			return nil, err
		}
		return &got, nil
	}

	return nil, errors.New("invalid image format")
}
