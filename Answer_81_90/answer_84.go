package main

import (
	"fmt"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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

func main() {
	db := [10][13]int{}
	trainFiles, err := filepath.Glob("./../Question_81_90/dataset/train_*.jpg")
	if err != nil {
		panic(err)
	}

	for dbIndex, trainFileName := range trainFiles {
		trainFile, err := os.Open(trainFileName)
		defer trainFile.Close()
		if err != nil {
			log.Fatal(err)
		}
		trainImage, err := jpeg.Decode(trainFile)
		if err != nil {
			log.Fatal(err)
		}

		H := trainImage.Bounds().Size().Y
		W := trainImage.Bounds().Size().X

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
				r, g, b, _ := trainImage.At(x, y).RGBA()
				reduceR[y][x] = reduceColors(uint8(float64(r*0xFF) / 0xFFFF))
				reduceG[y][x] = reduceColors(uint8(float64(g*0xFF) / 0xFFFF))
				reduceB[y][x] = reduceColors(uint8(float64(b*0xFF) / 0xFFFF))
			}
		}

		hist1d := []float64{}
		blueMap := map[uint8]int{32: 0, 96: 1, 160: 2, 224: 3}
		greenMap := map[uint8]int{32: 4, 96: 5, 160: 6, 224: 7}
		redMap := map[uint8]int{32: 8, 96: 9, 160: 10, 224: 11}

		for _, row := range reduceB {
			for _, elm := range row {
				hist1d = append(hist1d, float64(blueMap[elm]))
			}
		}
		for _, row := range reduceG {
			for _, elm := range row {
				hist1d = append(hist1d, float64(greenMap[elm]))
			}
		}
		for _, row := range reduceR {
			for _, elm := range row {
				hist1d = append(hist1d, float64(redMap[elm]))
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
			}
		}
		db[dbIndex] = hist

		// ここから下はヒストグラム表示
		p, err := plot.New()
		if err != nil {
			panic(err)
		}

		fileBaseName := filepath.Base(trainFileName)
		p.Title.Text = fileBaseName

		h, err := plotter.NewHist(plotter.Values(hist1d), 12)
		if err != nil {
			panic(err)
		}
		h.FillColor = color.Color(color.NRGBA{31, 119, 180, 255})
		p.Add(h)

		file := strings.Replace("answer_84_"+fileBaseName, "jpg", "png", -1)
		if err := p.Save(10*vg.Inch, 6*vg.Inch, file); err != nil {
			panic(err)
		}

	}

	for _, elm := range db {
		fmt.Println(elm)
	}

}
