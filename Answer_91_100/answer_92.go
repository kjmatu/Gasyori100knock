package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
	"reflect"
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

func assignIndex(jpegImage image.Image, clsColor [][]float64) ([][]int, []int) {
	H := jpegImage.Bounds().Size().Y
	W := jpegImage.Bounds().Size().X

	clsIndexCount := make([]int, len(clsColor))
	// clsIndexCount := [5]int{}
	clsIndexArray := make([][]int, H)
	for y := range clsIndexArray {
		clsIndexArray[y] = make([]int, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			distanceArray := [5]float64{}
			r, g, b, _ := jpegImage.At(x, y).RGBA()
			rf := float64(r) * 0xFF / 0xFFFF
			gf := float64(g) * 0xFF / 0xFFFF
			bf := float64(b) * 0xFF / 0xFFFF

			for i, color := range clsColor {
				rDiff := math.Abs(rf - color[0])
				gDiff := math.Abs(gf - color[1])
				bDiff := math.Abs(bf - color[2])
				distance := rDiff + gDiff + bDiff
				distanceArray[i] = distance
			}
			clsIndexArray[y][x] = searchMinIndex(distanceArray[:])
			clsIndexCount[clsIndexArray[y][x]]++
		}
	}
	return clsIndexArray, clsIndexCount
}

func calcAverageColorCls(jpegImage image.Image, clsIndexArray [][]int, clsIndexCount []int) [][]float64 {
	H := jpegImage.Bounds().Size().Y
	W := jpegImage.Bounds().Size().X

	clsColor := make([][]float64, len(clsIndexCount))
	for i := range clsColor {
		clsColor[i] = make([]float64, 3)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r, g, b, _ := jpegImage.At(x, y).RGBA()
			rf := float64(r) * 0xFF / 0xFFFF
			gf := float64(g) * 0xFF / 0xFFFF
			bf := float64(b) * 0xFF / 0xFFFF

			// fmt.Println("clsIndexArray[y][x]")
			clsColor[clsIndexArray[y][x]][0] += rf
			clsColor[clsIndexArray[y][x]][1] += gf
			clsColor[clsIndexArray[y][x]][2] += bf
		}
	}

	for i := range clsIndexCount {
		for j := range clsColor[i] {
			clsColor[i][j] /= float64(clsIndexCount[i])
		}
	}
	fmt.Println(clsColor)
	return clsColor
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
	randomSampleColorArray := [][]float64{
		{148, 121, 140},
		{122, 109, 135},
		{214, 189, 211},
		{84, 86, 135},
		{90, 102, 116},
	}

	H := jpegImage.Bounds().Size().Y
	W := jpegImage.Bounds().Size().X

	lastClsIndexArray := make([][]int, H)
	for i := range lastClsIndexArray {
		lastClsIndexArray[i] = make([]int, W)
	}

	tmpClsColor := make([][]float64, len(randomSampleColorArray))
	for i := range tmpClsColor {
		tmpClsColor[i] = make([]float64, len(randomSampleColorArray[0]))
	}

	clsColor := randomSampleColorArray
	// for _, elm := range clsColor {
	// 	fmt.Println(elm)
	// }

	for {
		clsIndexArray, clsIndexCount := assignIndex(jpegImage, randomSampleColorArray)

		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				r, g, b, _ := jpegImage.At(x, y).RGBA()
				rf := float64(r) * 0xFF / 0xFFFF
				gf := float64(g) * 0xFF / 0xFFFF
				bf := float64(b) * 0xFF / 0xFFFF

				// fmt.Println("clsIndexArray[y][x]")
				clsColor[clsIndexArray[y][x]][0] += rf
				clsColor[clsIndexArray[y][x]][1] += gf
				clsColor[clsIndexArray[y][x]][2] += bf
			}
		}

		for i := range clsIndexCount {
			for j := range clsColor[i] {
				clsColor[i][j] /= float64(clsIndexCount[i])
			}
		}

		if reflect.DeepEqual(tmpClsColor, clsColor) {
			break
		} else {
			for i := range clsColor {
				copy(tmpClsColor[i], clsColor[i])
			}
			for i := range clsIndexArray {
				copy(lastClsIndexArray[i], clsIndexArray[i])
			}
		}

	}

	for _, row := range clsColor {
		fmt.Println(row)
	}

	indexImage := image.NewRGBA(jpegImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			indexColor := clsColor[lastClsIndexArray[y][x]]
			// fmt.Println(indexColor)
			color := color.NRGBA{uint8(indexColor[0]), uint8(indexColor[1]), uint8(indexColor[2]), 255}
			indexImage.Set(x, y, color)
		}
	}
	indexFile, err := os.Create("./answer_92.jpg")
	defer indexFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(indexFile, indexImage, &jpeg.Options{100})

}
