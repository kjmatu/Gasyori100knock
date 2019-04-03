package main

import (
	"image"
	"image/color"
	"image/jpeg"
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
	jpegImage, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	trackPoints := []image.Point{}
	for height := 0; height < jpegImage.Bounds().Size().Y; height++ {
		for width := 0; width < jpegImage.Bounds().Size().X; width++ {
			r32, g32, b32, _ := jpegImage.At(width, height).RGBA()
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

	maskImage := image.NewGray(jpegImage.Bounds())
	for _, trackPoint := range trackPoints {
		maskImage.Set(trackPoint.X, trackPoint.Y, color.Gray{uint8(1)})
	}

	for y := 0; y < maskImage.Bounds().Size().Y; y++ {
		for x := 0; x < maskImage.Bounds().Size().X; x++ {
			if maskImage.GrayAt(x, y).Y == 1 {
				maskImage.Set(x, y, color.Gray{uint8(0)})
			} else {
				maskImage.Set(x, y, color.Gray{uint8(1)})
			}
		}
	}

	maskedImage := image.NewNRGBA64(jpegImage.Bounds())
	for y := 0; y < maskedImage.Bounds().Size().Y; y++ {
		for x := 0; x < maskedImage.Bounds().Size().X; x++ {
			maskValue := maskImage.GrayAt(x, y).Y
			r, g, b, a := jpegImage.At(x, y).RGBA()
			r *= uint32(maskValue)
			g *= uint32(maskValue)
			b *= uint32(maskValue)
			maskColor := color.NRGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			maskedImage.SetNRGBA64(x, y, maskColor)
		}
	}

	maskedFile, err := os.Create("./answer_71.jpeg")
	defer maskedFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(maskedFile, maskedImage, &jpeg.Options{100})
}
