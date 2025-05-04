package application

import (
	"ImageClustering/internal/domain"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func CreatePalette(clusters []domain.Pixel, paletteHeight, cellWigth int, path string) (string, error) {
	rect := image.Rect(0, 0, cellWigth*len(clusters), paletteHeight)
	newImage := image.NewRGBA(rect)

	bounds := newImage.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			ind := x / cellWigth
			newCol := color.RGBA{clusters[ind].R, clusters[ind].G, clusters[ind].B, 1}
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
