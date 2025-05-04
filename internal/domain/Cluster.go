package domain

import (
	"sync"
)

type Cluster struct {
	Centroid Pixel
	Members  *[]Pixel
	Mutex    *sync.Mutex
}

func NewCluster(centroid Pixel) *Cluster {
	return &Cluster{Centroid: centroid, Members: &[]Pixel{}, Mutex: &sync.Mutex{}}
}

//func (c *Cluster) GeneratePalette(palHeight, ceilWidth int, filepath string) error {
//	image := image.NewRGBA(image.Rect(0, 0, ceilWidth * len(*c.Members),  palHeight))
//
//	for i := 0; i < image.Rect.Bounds().Dx(); i++ {
//		for j := 0; j < image.Rect.Bounds().Dy(); j++ {
//
//		}
//	}
//}
