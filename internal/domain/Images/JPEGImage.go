package Images

import (
	"ImageClustering/internal/domain"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"
)

type JpegImage struct {
	image  image.Image
	pixels []domain.Pixel
}

func NewJpegImage(file io.Reader) (JpegImage, error) {
	ret := new(JpegImage)
	var err error
	ret.image, err = jpeg.Decode(file)

	if err != nil {
		return JpegImage{}, err
	}

	return *ret, nil
}

func (img *JpegImage) GetPixels() ([]domain.Pixel, error) {
	img.pixels = domain.ImageToPixels(img.image)
	return img.pixels, nil
}

func (img *JpegImage) CreateSimplified(clusters []domain.Pixel, path string) (string, error) {
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

	outFile, err := os.Create(path + ".jpeg")
	if err != nil {
		return "", err
	}

	defer outFile.Close()

	err = jpeg.Encode(outFile, newImage, nil)
	if err != nil {
		return "", err
	}

	return path, nil
}
