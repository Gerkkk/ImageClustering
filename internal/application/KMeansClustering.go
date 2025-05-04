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
	Cd         *domain.ColorDistance
}

func NewKMeansClustering(K, Iterations, Batch int, Data *[]domain.Pixel) *KMeansClustering {
	return &KMeansClustering{K: K, Iterations: Iterations, Clusters: make([]domain.Cluster, 0),
		Data: Data, Batch: Batch, Cd: &domain.ColorDistance{}}
}

func (km *KMeansClustering) AddClustering(lastPixel domain.Pixel) {
	newPixels := []domain.Pixel{}
	distPixels := []int{}

	for i := 0; i < km.Batch; i++ {
		nextI := rand.Intn(len(*km.Data))
		nextVal := (*km.Data)[nextI]
		newPixels = append(newPixels, nextVal)
		distPixels = append(distPixels, int(math.Pow(domain.EuclidianDistance(lastPixel, nextVal, km.Cd), 2)))
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

//func (km *KMeansClustering) InitClusters() {
//	// Шаг 1: случайный первый центр
//	first := domain.NewCluster((*km.Data)[rand.Intn(len(*km.Data))])
//	km.Clusters = append(km.Clusters, *first)
//
//	for i := 1; i < km.K; i++ {
//		distances := make([]float64, len(*km.Data))
//
//		for idx, p := range *km.Data {
//			minDist := math.MaxFloat64
//			for _, c := range km.Clusters {
//				dist := domain.EuclidianDistance(p, c.Centroid, km.Cd)
//				if dist < minDist {
//					minDist = dist
//				}
//			}
//			distances[idx] = minDist * minDist // k-means++ использует квадрат расстояния
//		}
//
//		// Выбор нового центра по вероятности
//		sum := 0.0
//		for _, d := range distances {
//			sum += d
//		}
//
//		r := rand.Float64() * sum
//		cumulative := 0.0
//		for idx, d := range distances {
//			cumulative += d
//			if cumulative >= r {
//				newClust := domain.NewCluster((*km.Data)[idx])
//				km.Clusters = append(km.Clusters, *newClust)
//				break
//			}
//		}
//	}
//}

func (km *KMeansClustering) InitClusters() {
	first := domain.NewCluster((*km.Data)[rand.Intn(len(*km.Data))])
	km.Clusters = append(km.Clusters, *first)

	for i := 0; i < km.K-1; i++ {
		km.AddClustering(km.Clusters[len(km.Clusters)-1].Centroid)
	}
}

func (km *KMeansClustering) DoClustering() []domain.Pixel {
	for iteration := 0; iteration < km.Iterations; iteration++ {
		//fmt.Println(iteration)
		wg := &sync.WaitGroup{}
		wg.Add(len(*km.Data))

		for _, pixel := range *km.Data {
			go func(pixel domain.Pixel) {
				defer wg.Done()
				i := pixel.AssignCluster(&km.Clusters)

				km.Clusters[i].Mutex.Lock()
				*km.Clusters[i].Members = append(*km.Clusters[i].Members, pixel)
				km.Clusters[i].Mutex.Unlock()
			}(pixel)
		}

		wg.Wait()

		wg = &sync.WaitGroup{}
		wg.Add(len(km.Clusters))
		for i, _ := range km.Clusters {
			ii := i
			go func(i int) {
				defer wg.Done()

				km.Clusters[i].Mutex.Lock()
				defer km.Clusters[i].Mutex.Unlock()

				var R, G, B = 0, 0, 0

				for _, pixel := range *km.Clusters[i].Members {
					R += int(pixel.R)
					G += int(pixel.G)
					B += int(pixel.B)
				}

				var newR, newG, newB uint8
				if len(*km.Clusters[i].Members) != 0 {
					newR = uint8(math.Round(float64(R) / float64(len(*km.Clusters[i].Members))))
					newG = uint8(math.Round(float64(G) / float64(len(*km.Clusters[i].Members))))
					newB = uint8(math.Round(float64(B) / float64(len(*km.Clusters[i].Members))))
				} else {
					newR = 0
					newG = 0
					newB = 0
				}

				km.Clusters[i].Centroid = domain.Pixel{R: newR, G: newG, B: newB}
				km.Clusters[i].Members = &[]domain.Pixel{}
			}(ii)
		}

		wg.Wait()
	}

	var ret []domain.Pixel
	for _, cluster := range km.Clusters {
		ret = append(ret, cluster.Centroid)
	}

	return ret
}
