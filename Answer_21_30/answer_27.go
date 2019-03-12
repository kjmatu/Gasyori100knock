package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func calcBicubicWeight(d, a float64) float64 {
	if d < 0 {
		return math.NaN()
	}
	d = math.Abs(d)
	var weight float64
	if d >= 0 && d <= 1 {
		weight = (a+2.0)*d*d*d - (a+3.0)*d*d + 1.0
	} else if d > 1 && d <= 2 {
		weight = a*d*d*d - 5.0*a*d*d + 8.0*a*d - 4.0*a
	} else {
		weight = 0.0
	}
	return weight
}

func main() {
	file, err := os.Open("./../Question_21_30/imori.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	scale := 1.5

	// 拡大縮小画像を作成
	scaleBounds := jimg.Bounds()
	scaleBounds.Max.X = int(float64(scaleBounds.Max.X) * scale)
	scaleBounds.Max.Y = int(float64(scaleBounds.Max.Y) * scale)
	biCubicImg := image.NewRGBA(scaleBounds)

	// 拡大画像の各位置に対して、それに対応する元画像の周囲4ピクセルから画素値を計算する
	for height := 0; height < biCubicImg.Bounds().Size().Y; height++ {
		for width := 0; width < biCubicImg.Bounds().Size().X; width++ {
			weightColor := [3]float64{}
			weightSum := 0.0

			// 拡大画像のピクセル位置に対応する元画像位置を計算
			srcX := float64(width) / scale
			srcY := float64(height) / scale

			for y := -1; y < 3; y++ {
				for x := -1; x < 3; x++ {
					// 元画像の位置を中心とした4x4ピクセルの位置(xCurrent, yCurrent)を計算する
					xCurrent := int(srcX) + x
					yCurrent := int(srcY) + y
					if (xCurrent >= jimg.Bounds().Size().X) ||
						(yCurrent >= jimg.Bounds().Size().Y) {
						continue
					}

					// 距離を計算する
					dx := math.Abs(float64(xCurrent) - srcX)
					dy := math.Abs(float64(yCurrent) - srcY)

					// 重みを計算
					xweight := calcBicubicWeight(dx, -1.0)
					if xweight == 0.0 {
						continue
					}

					yweight := calcBicubicWeight(dy, -1.0)
					if yweight == 0.0 {
						continue
					}

					weight := xweight * yweight
					weightSum += weight

					// 重みと画素値から拡大した画像の画素値を計算する
					Ixy := jimg.At(xCurrent, yCurrent)
					r, g, b, _ := Ixy.RGBA()

					weightColor[0] += (weight * float64(r*0xFF)) / float64(0xFFFF)
					weightColor[1] += (weight * float64(g*0xFF)) / float64(0xFFFF)
					weightColor[2] += (weight * float64(b*0xFF)) / float64(0xFFFF)
				}
			}

			for k := 0; k < len(weightColor); k++ {
				weightColor[k] /= weightSum
				if weightColor[k] > 255.0 {
					weightColor[k] = 255.0
				}
			}

			biCUbicColor := color.RGBA{}
			biCUbicColor.R = uint8(weightColor[0])
			biCUbicColor.G = uint8(weightColor[1])
			biCUbicColor.B = uint8(weightColor[2])
			biCUbicColor.A = 255
			biCubicImg.Set(width, height, biCUbicColor)
		}
	}

	biCubicImgFile, err := os.Create("./answer_27.jpg")
	defer biCubicImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(biCubicImgFile, biCubicImg, &jpeg.Options{100})

}
