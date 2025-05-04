package domain

import (
	"image"
)

type Pixel struct {
	R, G, B uint8
}

func NewPixel(r, g, b uint8) *Pixel {
	return &Pixel{R: r, G: g, B: b}
}

func EuclidianDistance(pixel1, pixel2 Pixel, cd *ColorDistance) float64 {
	return cd.CIE76(pixel1, pixel2)
	//return math.Sqrt(math.Pow(float64(pixel1.R-pixel2.R), 2) + math.Pow(float64(pixel1.G-pixel2.G), 2) + math.Pow(float64(pixel1.B-pixel2.B), 2))
}

func (p Pixel) AssignCluster(clusters *[]Cluster) int {
	cd := &ColorDistance{}
	var mi = -1.0
	var mii = -1

	for i, cluster := range *clusters {
		cur := EuclidianDistance(cluster.Centroid, p, cd)
		if mi == -1.0 || cur < mi {
			mi = cur
			mii = i
		}
	}

	return mii
}

func (p *Pixel) Predict(clusters *[]Pixel) Pixel {
	var mi = -1.0
	var mii = -1
	cd := &ColorDistance{}

	for i, cluster := range *clusters {
		cur := EuclidianDistance(cluster, *p, cd)
		if mi == -1.0 || cur < mi {
			mi = cur
			mii = i
		}
	}

	return (*clusters)[mii]
}

// TODO: this also may be paralleled
func ImageToPixels(image image.Image) []Pixel {
	pixels := make([]Pixel, 0)
	bounds := image.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := image.At(x, y).RGBA()
			pixels = append(pixels, Pixel{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
			})
		}
	}

	return pixels
}
