package main

import (
	"fmt"
	"image/color"
	"image/jpeg"
	"log"
	"os"

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
	// trainFiles, err := filepath.Glob("./../Question_81_90/dataset/train_*.jpg")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(trainFiles)
	// for _, elm := range trainFiles {
	// 	fmt.Println(filepath.Base(elm))
	// }

	// for _, trainFileName := range trainFiles {
	// 	trainFile, err := os.Open(trainFileName)
	// 	defer trainFile.Close()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	trainImage, err := jpeg.Decode(trainFile)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	H := trainImage.Bounds().Size().Y
	// 	W := trainImage.Bounds().Size().X

	// 	reduceR := make([][]uint8, H)
	// 	reduceG := make([][]uint8, H)
	// 	reduceB := make([][]uint8, H)
	// 	for i := range reduceR {
	// 		reduceR[i] = make([]uint8, W)
	// 		reduceG[i] = make([]uint8, W)
	// 		reduceB[i] = make([]uint8, W)
	// 	}

	// 	for y := 0; y < H; y++ {
	// 		for x := 0; x < W; x++ {
	// 			r, g, b, _ := trainImage.At(x, y).RGBA()
	// 			reduceR[y][x] = reduceColors(uint8(float64(r*0xFF) / 0xFFFF))
	// 			reduceG[y][x] = reduceColors(uint8(float64(g*0xFF) / 0xFFFF))
	// 			reduceB[y][x] = reduceColors(uint8(float64(b*0xFF) / 0xFFFF))
	// 		}
	// 	}

	// 	hist1d := []float64{}
	// 	for _, row := range reduceB {
	// 		for _, elm := range row {
	// 			hist1d = append(hist1d, float64(elm))
	// 		}
	// 	}
	// 	for _, row := range reduceG {
	// 		for _, elm := range row {
	// 			hist1d = append(hist1d, float64(elm)+255)
	// 		}
	// 	}
	// 	for _, row := range reduceR {
	// 		for _, elm := range row {
	// 			hist1d = append(hist1d, float64(elm)+510)
	// 		}
	// 	}

	// 	// hist := [13]int{}
	// 	// for y := 0; y < H; y++ {
	// 	// 	for x := 0; x < W; x++ {
	// 	// 		if reduceB[y][x] == 32 {
	// 	// 			hist[1]++
	// 	// 		} else if reduceB[y][x] == 96 {
	// 	// 			hist[2]++
	// 	// 		} else if reduceB[y][x] == 160 {
	// 	// 			hist[3]++
	// 	// 		} else if reduceB[y][x] == 224 {
	// 	// 			hist[4]++
	// 	// 		}

	// 	// 		if reduceG[y][x] == 32 {
	// 	// 			hist[5]++
	// 	// 		} else if reduceG[y][x] == 96 {
	// 	// 			hist[6]++
	// 	// 		} else if reduceG[y][x] == 160 {
	// 	// 			hist[7]++
	// 	// 		} else if reduceG[y][x] == 224 {
	// 	// 			hist[8]++
	// 	// 		}

	// 	// 		if reduceR[y][x] == 32 {
	// 	// 			hist[9]++
	// 	// 		} else if reduceR[y][x] == 96 {
	// 	// 			hist[10]++
	// 	// 		} else if reduceR[y][x] == 160 {
	// 	// 			hist[11]++
	// 	// 		} else if reduceR[y][x] == 224 {
	// 	// 			hist[12]++
	// 	// 		}
	// 	// 	}
	// 	// }

	// 	// ここから下はヒストグラム表示
	// 	p, err := plot.New()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fileBaseName := filepath.Base(trainFileName)
	// 	p.Title.Text = fileBaseName

	// 	h, err := plotter.NewHist(plotter.Values(hist1d), 13)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	h.FillColor = color.Color(color.NRGBA{31, 119, 180, 255})
	// 	p.Add(h)

	// 	file := strings.Replace(fileBaseName, "jpg", "png", -1)
	// 	if err := p.Save(10*vg.Inch, 6*vg.Inch, file); err != nil {
	// 		panic(err)
	// 	}

	// }

	trainFile, err := os.Open("./../Question_81_90/dataset/train_akahara_1.jpg")
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
	for _, row := range reduceB {
		for _, elm := range row {
			hist1d = append(hist1d, float64(elm)+0)
		}
	}
	for _, row := range reduceG {
		for _, elm := range row {
			hist1d = append(hist1d, float64(elm)+255)
		}
	}
	for _, row := range reduceR {
		for _, elm := range row {
			hist1d = append(hist1d, float64(elm)+510)
		}
	}

	fmt.Println(hist1d)

	// hist := [13]int{}
	// for y := 0; y < H; y++ {
	// 	for x := 0; x < W; x++ {
	// 		if reduceB[y][x] == 32 {
	// 			hist[1]++
	// 		} else if reduceB[y][x] == 96 {
	// 			hist[2]++
	// 		} else if reduceB[y][x] == 160 {
	// 			hist[3]++
	// 		} else if reduceB[y][x] == 224 {
	// 			hist[4]++
	// 		}

	// 		if reduceG[y][x] == 32 {
	// 			hist[5]++
	// 		} else if reduceG[y][x] == 96 {
	// 			hist[6]++
	// 		} else if reduceG[y][x] == 160 {
	// 			hist[7]++
	// 		} else if reduceG[y][x] == 224 {
	// 			hist[8]++
	// 		}

	// 		if reduceR[y][x] == 32 {
	// 			hist[9]++
	// 		} else if reduceR[y][x] == 96 {
	// 			hist[10]++
	// 		} else if reduceR[y][x] == 160 {
	// 			hist[11]++
	// 		} else if reduceR[y][x] == 224 {
	// 			hist[12]++
	// 		}
	// 	}
	// }

	// ここから下はヒストグラム表示
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "train_akahara_1.jpg"

	h, err := plotter.NewHist(plotter.Values(hist1d), 13)
	if err != nil {
		panic(err)
	}
	h.FillColor = color.Color(color.NRGBA{31, 119, 180, 255})
	p.Add(h)

	file := "hist.png"
	if err := p.Save(10*vg.Inch, 6*vg.Inch, file); err != nil {
		panic(err)
	}

}
