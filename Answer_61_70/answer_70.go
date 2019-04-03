package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
)

func getMinMaxFromArray(array []float64) (float64, float64) {
	min := math.Inf(1)
	max := math.Inf(-1)
	for index := 0; index < len(array); index++ {
		if array[index] < min {
			min = array[index]
		}

		if array[index] > max {
			max = array[index]
		}
	}
	return min, max
}

func main() {
	file, err := os.Open("./../Question_61_70/imori.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpenImage, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	trackPoints := []image.Point{}
	for height := 0; height < jpenImage.Bounds().Size().Y; height++ {
		for width := 0; width < jpenImage.Bounds().Size().X; width++ {
			r32, g32, b32, _ := jpenImage.At(width, height).RGBA()
			// uint32bitからfloat64へ変換
			rfloat64 := float64(r32) / 0xFFFF
			gfloat64 := float64(g32) / 0xFFFF
			bfloat64 := float64(b32) / 0xFFFF

			// RGBからHSVへ変換
			array := []float64{rfloat64, gfloat64, bfloat64}
			min, max := getMinMaxFromArray(array)
			var h float64
			// var h, s, v float64
			if min == max {
				h = 0.0
			} else if min == bfloat64 {
				h = 60.0*(gfloat64-rfloat64)/(max-min) + 60.0
			} else if min == rfloat64 {
				h = 60.0*(bfloat64-gfloat64)/(max-min) + 180.0
			} else if min == gfloat64 {
				h = 60.0*(rfloat64-bfloat64)/(max-min) + 300.0
			} else {
				h = math.NaN()
			}
			// v = max
			// s = max - min

			// 色相Hが0.0~360.0の範囲になるように整形
			h = math.Mod(h, 360.0)
			if h >= 180 && h <= 260 {
				trackPoints = append(trackPoints, image.Point{width, height})
			}
		}
	}

	trackImage := image.NewGray(jpenImage.Bounds())
	for _, trackPoint := range trackPoints {
		trackImage.Set(trackPoint.X, trackPoint.Y, color.Gray{uint8(255)})
	}

	trackFile, err := os.Create("./answer_70.png")
	defer trackFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	png.Encode(trackFile, trackImageg)
}
