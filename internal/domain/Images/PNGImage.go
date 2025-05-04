package Images

import (
	"ImageClustering/internal/domain"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
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
	newImage := image.NewRGBA(img.image.Bounds())

	bounds := newImage.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.image.At(x, y).RGBA()
			newPic := domain.Pixel{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
			newPic = newPic.Predict(&clusters)
			newCol := color.RGBA{newPic.R, newPic.G, newPic.B, uint8(a >> 8)}
			newImage.Set(x, y, newCol)
		}
	}

	outFile, err := os.Create(path + ".png")
	if err != nil {
		return "", err
	}

	defer outFile.Close()

	err = png.Encode(outFile, newImage)
	if err != nil {
		return "", err
	}

	return path, nil
}
