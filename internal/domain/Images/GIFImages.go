package Images

import (
	"ImageClustering/internal/domain"
	"image/gif"
	"io"
)

type GifImage struct {
	gif    *gif.GIF
	pixels []domain.Pixel
}

func NewGifImage(file io.Reader) (GifImage, error) {
	ret := new(GifImage)
	var err error
	ret.gif, err = gif.DecodeAll(file)

	if err != nil {
		return GifImage{}, err
	}

	return *ret, nil
}

func (img *GifImage) GetPixels() ([]domain.Pixel, error) {
	img.pixels = make([]domain.Pixel, 0)

	for _, image := range img.gif.Image {
		img.pixels = append(img.pixels, domain.ImageToPixels(image)...)
	}

	return img.pixels, nil
}

func (img *GifImage) CreateSimplified(clusters []domain.Pixel, path string) (string, error) {
	return "", nil
}
