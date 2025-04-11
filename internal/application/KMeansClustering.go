package application

import (
	"ImageClustering/internal/domain"
	"math"
	"math/rand"
	"sync"
)

type KMeansClustering struct {
	Data       *[]domain.Pixel
	K          int
	Iterations int
	Clusters   []domain.Cluster
	Batch      int
}

func NewKMeansClustering(K, Iterations, Batch int, Data *[]domain.Pixel) *KMeansClustering {
	return &KMeansClustering{K: K, Iterations: Iterations, Clusters: make([]domain.Cluster, 0),
		Data: Data, Batch: Batch}
}

func (km *KMeansClustering) AddClustering(lastPixel domain.Pixel) {
	newPixels := []domain.Pixel{}
	distPixels := []int{}

	for i := 0; i < km.Batch; i++ {
		nextI := rand.Intn(len(*km.Data))
		nextVal := (*km.Data)[nextI]
		newPixels = append(newPixels, nextVal)
		distPixels = append(distPixels, int(math.Pow(domain.EuclidianDistance(lastPixel, nextVal), 2)))
	}

	var total = 0
	for _, dist := range distPixels {
		total += dist
	}

	r := int(math.Ceil(rand.Float64() * float64(total)))
	for i, weight := range distPixels {
		r -= weight
		if r <= 0 {
			newClust := domain.NewCluster(newPixels[i])
			km.Clusters = append(km.Clusters, *newClust)
			break
		}
	}

	return
}

func (km *KMeansClustering) InitClusters() {
	first := domain.NewCluster((*km.Data)[rand.Intn(len(*km.Data))])
	km.Clusters = append(km.Clusters, *first)

	for i := 0; i < km.K-1; i++ {
		km.AddClustering(km.Clusters[len(km.Clusters)-1].Centroid)
	}
}

func (km *KMeansClustering) DoClustering() []domain.Pixel {

	for iteration := 0; iteration < km.Iterations; iteration++ {
		wg := &sync.WaitGroup{}
		wg.Add(len(*km.Data))

		for _, pixel := range *km.Data {

			go func(pixel domain.Pixel) {
				defer wg.Done()
				i := pixel.AssignCluster(&km.Clusters)
				*km.Clusters[i].Members = append(*km.Clusters[i].Members, pixel)
			}(pixel)
		}

		wg.Wait()

		for i, cluster := range km.Clusters {
			ii := i
			go func(i int) {
				var R, G, B = 0, 0, 0

				for _, pixel := range *km.Clusters[i].Members {
					R += int(pixel.R)
					G += int(pixel.G)
					B += int(pixel.B)
				}

				newR := uint8(math.Round(float64(R) / float64(len(*cluster.Members))))
				newG := uint8(math.Round(float64(G) / float64(len(*cluster.Members))))
				newB := uint8(math.Round(float64(B) / float64(len(*cluster.Members))))

				cluster.Centroid = domain.Pixel{R: newR, G: newG, B: newB}
				cluster.Members = &[]domain.Pixel{}
			}(ii)
		}
	}

	var ret []domain.Pixel
	for _, cluster := range km.Clusters {
		ret = append(ret, cluster.Centroid)
	}

	return ret
}
