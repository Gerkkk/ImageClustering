package domain

type Cluster struct {
	Centroid Pixel
	Members  *[]Pixel
}

func NewCluster(centroid Pixel) *Cluster {
	return &Cluster{Centroid: centroid, Members: &[]Pixel{}}
}
