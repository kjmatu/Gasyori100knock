package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func searchMinIndex(array []float64) int {
	min := math.Inf(1)
	minIndex := 0
	for i, elm := range array {
		if elm < min {
			min = elm
			minIndex = i
		}
	}
	return minIndex
}

func main() {
	file, err := os.Open("./../Question_91_100/imori.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpegImage, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// K := 5
	randomSampleColorArray := [5][3]int{
		{148, 121, 140},
		{122, 109, 135},
		{214, 189, 211},
		{84, 86, 135},
		{90, 102, 116},
	}

	H := jpegImage.Bounds().Size().Y
	W := jpegImage.Bounds().Size().X

	clsIndexArray := make([][]int, H)
	for i := range clsIndexArray {
		clsIndexArray[i] = make([]int, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			distanceArray := [5]float64{}
			r, g, b, _ := jpegImage.At(x, y).RGBA()
			rf := float64(r) * 0xFF / 0xFFFF
			gf := float64(g) * 0xFF / 0xFFFF
			bf := float64(b) * 0xFF / 0xFFFF

			for i, color := range randomSampleColorArray {
				rDiff := math.Abs(rf - float64(color[0]))
				gDiff := math.Abs(gf - float64(color[1]))
				bDiff := math.Abs(bf - float64(color[2]))
				distance := rDiff + gDiff + bDiff
				distanceArray[i] = distance
			}
			clsIndexArray[y][x] = searchMinIndex(distanceArray[:])
		}
	}

	indexImage := image.NewGray(jpegImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			indexImage.SetGray(x, y, color.Gray{uint8(clsIndexArray[y][x] * 50)})
		}
	}
	indexFile, err := os.Create("./answer_91.jpg")
	defer indexFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(indexFile, indexImage, &jpeg.Options{100})

}
