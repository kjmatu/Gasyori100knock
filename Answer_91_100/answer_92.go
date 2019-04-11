package main

import (
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

func kMeansReduceColor(randomSampleColorClass [][]float64, reduceImage image.Image) ([][]float64, [][]int) {
	H := reduceImage.Bounds().Size().Y
	W := reduceImage.Bounds().Size().X

	clsIndexArray := make([][]int, H)
	for i := range clsIndexArray {
		clsIndexArray[i] = make([]int, W)
	}

	clsColor := randomSampleColorClass
	backUpClsColor := [][]float64{}
	for {
		clsIndexCount := [5]int{}
		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				distanceArray := [5]float64{}
				r, g, b, _ := reduceImage.At(x, y).RGBA()
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

		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				r, g, b, _ := reduceImage.At(x, y).RGBA()
				rf := float64(r) * 0xFF / 0xFFFF
				gf := float64(g) * 0xFF / 0xFFFF
				bf := float64(b) * 0xFF / 0xFFFF

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
		// fmt.Println(clsColor)
		if reflect.DeepEqual(backUpClsColor, clsColor) {
			// if backUpClsColor == clsColor {
			break
		} else {
			backUpClsColor = clsColor
		}
	}
	return clsColor, clsIndexArray
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
	randomSampleColorArray := [5][3]float64{
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

	backUpClsColor := [5][3]float64{}
	clsColor := randomSampleColorArray
	for {
		clsIndexCount := [5]int{}
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

		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				r, g, b, _ := jpegImage.At(x, y).RGBA()
				rf := float64(r) * 0xFF / 0xFFFF
				gf := float64(g) * 0xFF / 0xFFFF
				bf := float64(b) * 0xFF / 0xFFFF

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
		// fmt.Println(clsColor)
		if backUpClsColor == clsColor {
			break
		} else {
			backUpClsColor = clsColor
		}

	}

	// randomSampleColorSlice := make([][]float64, len(randomSampleColorArray))
	// for i, row := range randomSampleColorArray {
	// 	fmt.Println(i)
	// 	randomSampleColorSlice[i] = row[:]
	// }

	// clsColor, clsIndexArray := kMeansReduceColor(randomSampleColorSlice, jpegImage)

	// for _, row := range clsColor {
	// 	fmt.Println(row)
	// }

	indexImage := image.NewRGBA(jpegImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			indexColor := clsColor[clsIndexArray[y][x]]
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
