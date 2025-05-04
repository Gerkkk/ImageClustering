package main

import (
	"ImageClustering/internal/application"
	"ImageClustering/internal/domain/Images"
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"runtime"
)

func main() {
	var iterations int
	var maxprocs int
	var filepath string
	var clusters int
	var palette bool
	var simplified bool
	var frameNum = -1
	var batch int

	pflag.IntVarP(&iterations, "iterations", "i", 60, "Number of iterations")
	pflag.IntVarP(&maxprocs, "maxprocs", "m", 8, "Number of goroutines")
	pflag.IntVarP(&clusters, "clusters", "c", 20, "Number of clusters")
	pflag.IntVarP(&batch, "batch", "b", 1000, "Batch size")
	pflag.IntVarP(&frameNum, "frame", "", -1, "Number of frames that will be randomly chosen from gif file for clustering")
	pflag.StringVarP(&filepath, "file", "f", "", "File to read from")
	pflag.BoolVarP(&simplified, "simplified", "s", false, "Return simplified image or no")
	pflag.BoolVarP(&palette, "palette", "p", false, "Return palette or no")

	pflag.Parse()

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	fmt.Println(file.Name())
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}

	defer file.Close()

	runtime.GOMAXPROCS(maxprocs)

	input := Images.ImageConstructorData{File: file, FileName: file.Name(), FramesCount: frameNum}
	image, err := Images.NewImage(input)

	if err != nil {
		fmt.Printf("Error decoding file: %v\n", err)
		os.Exit(1)
	}

	pixels, err := image.GetPixels()

	if err != nil {
		fmt.Printf("Error getting pixels: %v\n", err)
		os.Exit(1)
	}

	clusteringObj := application.NewKMeansClustering(clusters, iterations, batch, &pixels)

	clusteringObj.InitClusters()
	clust := clusteringObj.DoClustering()

	if simplified == true {
		path, err := image.CreateSimplified(clust, "ret")

		if err != nil {
			fmt.Printf("Error creating simplified image: %v\n", err)
		} else {
			fmt.Printf("Created simplified image: %s\n", path)
		}
	}

	if palette == true {
		path, err := application.CreatePalette(clust, 200, 200, "ret-palette")

		if err != nil {
			fmt.Printf("Error creating palette of image: %v\n", err)
		} else {
			fmt.Printf("Created palette for image: %s\n", path)
		}
	}

}
