package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"math"
	"os"
	"path/filepath"
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

func calcGravity(db [4][13]int) ([12]float64, [12]float64) {
	cls1Average := [12]float64{}
	cls1Num := 0
	cls2Average := [12]float64{}
	cls2Num := 0
	for _, dbRow := range db {
		for i := 0; i < 12; i++ {
			if dbRow[12] == 0 {
				cls1Average[i] += float64(dbRow[i])
			}
			if dbRow[12] == 1 {
				cls2Average[i] += float64(dbRow[i])
			}
		}

		if dbRow[12] == 0 {
			cls1Num++
		} else if dbRow[12] == 1 {
			cls2Num++
		}
	}

	for i := 0; i < 12; i++ {
		if cls1Num != 0 {
			cls1Average[i] /= float64(cls1Num)
		}

		if cls2Num != 0 {
			cls2Average[i] /= float64(cls2Num)
		}
	}
	return cls1Average, cls2Average
}

func main() {
	db := [4][13]int{}
	trainFiles, err := filepath.Glob("./../Question_81_90/dataset/test_*.jpg")
	if err != nil {
		panic(err)
	}

	// ランダムにクラスを割り当てる
	// しかし、解答との比較をするためにクラスを手動で設定する
	// rand.Seed(4)
	// th := 0.7
	for dbIndex, trainFileName := range trainFiles {
		histogram := createReduceColorHistogram(trainFileName)
		// if rand.Float64() > th {
		// 	histogram[12] = 1
		// } else {
		// 	histogram[12] = 0
		// }
		if dbIndex == 1 {
			histogram[12] = 1
		} else {
			histogram[12] = 0
		}
		db[dbIndex] = histogram
	}

	for {
		changeCount := 0
		// 重心を計算する
		cls1Average, cls2Average := calcGravity(db)
		// ラベリングを行う
		cls1Distance := 0.0
		cls2Distance := 0.0
		for fIndex, feature := range db {
			for i := 0; i < 12; i++ {
				cls1Distance += math.Pow(math.Abs(cls1Average[i]-float64(feature[i])), 2)
			}
			cls1Distance = math.Sqrt(cls1Distance)

			for i := 0; i < 12; i++ {
				cls2Distance += math.Pow(math.Abs(cls2Average[i]-float64(feature[i])), 2)
			}
			cls2Distance = math.Sqrt(cls2Distance)

			var pred int
			if cls1Distance < cls2Distance {
				pred = 0
			} else {
				pred = 1
			}

			if db[fIndex][12] != pred {
				db[fIndex][12] = pred
				changeCount++
			}

		}
		if changeCount == 0 {
			break
		}
	}

	for i, fileName := range trainFiles {
		fmt.Println(filepath.Base(fileName), "Pred:", db[i][12])
	}
}
