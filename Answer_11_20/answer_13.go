package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func getArrayMaxMin(array []float64) (float64, float64) {
	min := math.Inf(0)
	max := math.Inf(-1)
	for index := 0; index < len(array); index++ {
		if min > array[index] {
			min = array[index]
		}

		if max < array[index] {
			max = array[index]
		}
	}
	return max, min
}

func main() {
	file, err := os.Open("./../assets/imori.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// グレースケール画像に変換
	grayimg := image.NewGray(jimg.Bounds())
	for height := 0; height < jimg.Bounds().Size().Y; height++ {
		for width := 0; width < jimg.Bounds().Size().X; width++ {
			ycbcr := jimg.At(width, height)
			r32, g32, b32, _ := ycbcr.RGBA()
			// 32bit画像から8bit画像に変換
			r8 := uint8((float64(r32) / 0xFFFF) * 0xFF)
			g8 := uint8((float64(g32) / 0xFFFF) * 0xFF)
			b8 := uint8((float64(b32) / 0xFFFF) * 0xFF)

			var graycolor color.Gray
			// カラーからグレースケールに変換
			grayfloat64 := 0.2126*float64(r8) + 0.7152*float64(g8) + 0.0722*float64(b8)
			graycolor.Y = uint8(grayfloat64)
			grayimg.Set(width, height, graycolor)
		}
	}

	// MAX-MINフィルタ画像を作成
	maxMinImg := image.NewGray(grayimg.Bounds())
	filterSize := 3
	for width := 0; width < grayimg.Bounds().Size().X; width++ {
		for height := 0; height < grayimg.Bounds().Size().Y; height++ {

			// 対象ピクセルを中心とした3x3ピクセルの画素値の最大値と最小値を取得する
			grayArray := make([]float64, filterSize*filterSize)
			index := 0
			for filterWidth := 0; filterWidth < filterSize; filterWidth++ {
				for filterHeight := 0; filterHeight < filterSize; filterHeight++ {
					// 対象ピクセルの位置を計算する
					srcPointX := filterWidth + width - 1
					srcPointY := filterHeight + height - 1
					var gray float64
					if (srcPointX < 0) || (srcPointX >= grayimg.Bounds().Size().X) ||
						(srcPointY < 0) || (srcPointY >= grayimg.Bounds().Size().Y) {
						// 0パディング
						gray = 0.0
					} else {
						gray = float64(grayimg.GrayAt(srcPointX, srcPointY).Y)
					}
					grayArray[index] = gray
					index++
				}
			}

			var maxMinColor color.Gray
			max, min := getArrayMaxMin(grayArray)
			maxMinColor.Y = uint8(max - min)
			maxMinImg.Set(width, height, maxMinColor)
		}
	}

	maxMinImgFile, err := os.Create("./answer_13.jpg")
	defer maxMinImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(maxMinImgFile, maxMinImg, &jpeg.Options{100})
}
