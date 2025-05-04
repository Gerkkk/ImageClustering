package Images

import (
	"ImageClustering/internal/domain"
	"image"
	"image/color"
	"image/gif"
	"io"
	"math/rand"
	"os"
	"sync"
)

type GifImage struct {
	gif       *gif.GIF
	pixels    []domain.Pixel
	framesNum int
}

type frameResult struct {
	index int
	frame *image.Paletted
	delay int
}

func NewGifImage(file io.Reader, framesNum int) (GifImage, error) {
	ret := new(GifImage)
	var err error
	ret.gif, err = gif.DecodeAll(file)
	ret.framesNum = framesNum

	if err != nil {
		return GifImage{}, err
	}

	return *ret, nil
}

func (img *GifImage) GetPixels() ([]domain.Pixel, error) {
	img.pixels = make([]domain.Pixel, 0)

	indeces := []int{}
	for i := 0; i < len(img.gif.Image); i++ {
		indeces = append(indeces, i)
	}

	rand.Shuffle(len(indeces), func(i, j int) {
		indeces[i], indeces[j] = indeces[j], indeces[i]
	})

	for u, i := range indeces {
		if u < img.framesNum {
			img.pixels = append(img.pixels, domain.ImageToPixels(img.gif.Image[i])...)
		} else {
			break
		}
	}

	return img.pixels, nil
}

func (img *GifImage) CreateSimplified(clusters []domain.Pixel, path string) (string, error) {
	newImage := gif.GIF{
		LoopCount: img.gif.LoopCount,
	}

	newFrames := make([]*image.Paletted, len(img.gif.Image))
	delays := make([]int, len(img.gif.Image))

	results := make(chan frameResult, len(img.gif.Image))

	wg := new(sync.WaitGroup)
	for u, frame := range img.gif.Image {
		wg.Add(1)

		go func(u int, frame *image.Paletted) {
			defer wg.Done()

			bounds := frame.Bounds()
			newFrame := image.NewPaletted(bounds, frame.Palette)

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, _ := frame.At(x, y).RGBA()
					pixel := domain.Pixel{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8)}
					pixel = pixel.Predict(&clusters)
					newFrame.Set(x, y, color.RGBA{pixel.R, pixel.G, pixel.B, 255})
				}
			}

			results <- frameResult{index: u, frame: newFrame, delay: img.gif.Delay[u]}
		}(u, frame)
	}

	wg.Wait()
	close(results)

	for res := range results {
		newFrames[res.index] = res.frame
		delays[res.index] = res.delay
	}

	newImage.Image = newFrames
	newImage.Delay = delays

	outFile, err := os.Create(path + ".gif")
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	err = gif.EncodeAll(outFile, &newImage)
	if err != nil {
		return "", err
	}

	return path, nil
}
