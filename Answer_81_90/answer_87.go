package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func reduceColors(pixVal uint8) uint8 {
	if (0 <= pixVal) && (pixVal < 64) {
		return 32
	} else if (64 <= pixVal) && (pixVal < 128) {
		return 96
	} else if (128 <= pixVal) && (pixVal < 192) {
		return 160
	} else if (192 <= pixVal) && (pixVal <= 255) {
		return 224
	}
	return 0
}

func createReduceColorHistogram(fileName string) [13]int {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	loadImage, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	H := loadImage.Bounds().Size().Y
	W := loadImage.Bounds().Size().X

	reduceR := make([][]uint8, H)
	reduceG := make([][]uint8, H)
	reduceB := make([][]uint8, H)
	for i := range reduceR {
		reduceR[i] = make([]uint8, W)
		reduceG[i] = make([]uint8, W)
		reduceB[i] = make([]uint8, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r, g, b, _ := loadImage.At(x, y).RGBA()
			reduceR[y][x] = reduceColors(uint8(float64(r*0xFF) / 0xFFFF))
			reduceG[y][x] = reduceColors(uint8(float64(g*0xFF) / 0xFFFF))
			reduceB[y][x] = reduceColors(uint8(float64(b*0xFF) / 0xFFFF))
		}
	}

	hist := [13]int{}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			if reduceB[y][x] == 32 {
				hist[0]++
			} else if reduceB[y][x] == 96 {
				hist[1]++
			} else if reduceB[y][x] == 160 {
				hist[2]++
			} else if reduceB[y][x] == 224 {
				hist[3]++
			}

			if reduceG[y][x] == 32 {
				hist[4]++
			} else if reduceG[y][x] == 96 {
				hist[5]++
			} else if reduceG[y][x] == 160 {
				hist[6]++
			} else if reduceG[y][x] == 224 {
				hist[7]++
			}

			if reduceR[y][x] == 32 {
				hist[8]++
			} else if reduceR[y][x] == 96 {
				hist[9]++
			} else if reduceR[y][x] == 160 {
				hist[10]++
			} else if reduceR[y][x] == 224 {
				hist[11]++
			}

			if strings.Contains(fileName, "akahara") {
				hist[12] = 0
			} else if strings.Contains(fileName, "madara") {
				hist[12] = 1
			}
		}
	}
	return hist
}

func main() {
	db := [10][13]int{}
	trainFiles, err := filepath.Glob("./../Question_81_90/dataset/train_*.jpg")
	if err != nil {
		panic(err)
	}

	for dbIndex, trainFileName := range trainFiles {
		histogram := createReduceColorHistogram(trainFileName)
		db[dbIndex] = histogram
	}

	testFiles, err := filepath.Glob("./../Question_81_90/dataset/test_*.jpg")
	if err != nil {
		panic(err)
	}

	sucessCount := 0
	for _, testFileName := range testFiles {
		testFileHist := createReduceColorHistogram(testFileName)
		diff := [10]float64{}
		pred := [10]int{}
		for dbIndex, dbHist := range db {
			featureVal := 0.0
			for histIndex := 0; histIndex < len(dbHist)-1; histIndex++ {
				featureVal += math.Abs(float64(testFileHist[histIndex] - dbHist[histIndex]))
			}
			diff[dbIndex] = featureVal
			pred[dbIndex] = int(math.Abs(float64(testFileHist[len(dbHist)-1] - dbHist[len(dbHist)-1])))
		}

		diffMin := math.Inf(1)
		nearMap := map[float64]int{}
		diffMinIndex := 0
		for i, elm := range diff {
			nearMap[elm] = i
			if elm < diffMin {
				diffMin = elm
			}
			// go run で diffMinIndex declared and not usedが出るのを回避するための無駄な記述
			diffMinIndex = diffMinIndex
		}

		sort.Float64s(diff[:])
		similarFileName := ""
		akaharaCount := 0
		madaraCount := 0
		className := ""
		for _, diffElm := range diff[:3] {
			nearIndex := nearMap[diffElm]
			fileName := filepath.Base(trainFiles[nearIndex])
			similarFileName += fileName + ", "
			if strings.Contains(fileName, "akahara") {
				akaharaCount++
			}

			if strings.Contains(fileName, "madara") {
				madaraCount++
			}
		}

		if akaharaCount > madaraCount {
			className = "akahara"
		}

		if madaraCount > akaharaCount {
			className = "madara"
		}

		fmt.Println(filepath.Base(testFileName) + " is similar >> " + similarFileName)

		if strings.Contains(filepath.Base(testFileName), className) {
			sucessCount++
		}
	}

	fmt.Printf("Accuracy >> %f\n", float64(sucessCount)/float64(len(testFiles)))
}
