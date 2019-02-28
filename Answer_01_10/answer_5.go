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
	file, err := os.Open("./../assets/imori.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jImg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	invColorImg := image.NewRGBA(jImg.Bounds())
	for height := 0; height < jImg.Bounds().Size().Y; height++ {
		for width := 0; width < jImg.Bounds().Size().X; width++ {
			r32, g32, b32, _ := jImg.At(width, height).RGBA()
			// uint32bitからfloat64へ変換
			rfloat64 := float64(r32) / 0xFFFF
			gfloat64 := float64(g32) / 0xFFFF
			bfloat64 := float64(b32) / 0xFFFF

			// RGBからHSVへ変換
			array := []float64{rfloat64, gfloat64, bfloat64}
			min, max := getMinMaxFromArray(array)
			var h, s, v float64
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
			v = max
			s = max - min

			// 色相Hを反転
			h += 180.0

			// 色相Hが0.0~360.0の範囲になるように整形
			h = math.Mod(h, 360.0)

			// HSVからRGBへ変換
			c := s
			hh := h / 60.0
			x := c * (1.0 - math.Abs(math.Mod(hh, 2.0)-1.0))
			offset := v - c
			var r, g, b float64

			if (0.0 <= hh) && (hh < 1.0) {
				r = c
				g = x
				b = 0
			} else if (1.0 <= hh) && (hh < 2.0) {
				r = x
				g = c
				b = 0
			} else if (2.0 <= hh) && (hh < 3.0) {
				r = 0
				g = c
				b = x
			} else if (3.0 <= hh) && (hh < 4.0) {
				r = 0
				g = x
				b = c
			} else if (4.0 <= hh) && (hh < 5.0) {
				r = x
				g = 0
				b = c
			} else if (5.0 <= hh) && (hh < 6.0) {
				r = c
				g = 0
				b = x
			}
			r += offset
			g += offset
			b += offset

			// 色相を反転したRGB画像をセット
			var invColor color.RGBA
			invColor.R = uint8(r * 255)
			invColor.G = uint8(g * 255)
			invColor.B = uint8(b * 255)
			invColor.A = 255
			invColorImg.Set(width, height, invColor)
		}
	}
	invColorFile, err := os.Create("./answer_5.jpg")
	defer invColorFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(invColorFile, invColorImg, &jpeg.Options{100})

}
