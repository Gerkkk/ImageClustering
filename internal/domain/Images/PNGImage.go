package Images

import (
	"ImageClustering/internal/domain"
	"image"
	"image/png"
	"io"
)

type PngImage struct {
	image  image.Image
	pixels []domain.Pixel
}

func NewPngImage(file io.Reader) (PngImage, error) {
	ret := new(PngImage)
	var err error
	ret.image, err = png.Decode(file)

	if err != nil {
		return PngImage{}, err
	}

	return *ret, nil
}

func (img *PngImage) GetPixels() ([]domain.Pixel, error) {
	img.pixels = domain.ImageToPixels(img.image)
	return img.pixels, nil
}

func (img *PngImage) CreateSimplified(clusters []domain.Pixel, path string) (string, error) {
	return "", nil
}
