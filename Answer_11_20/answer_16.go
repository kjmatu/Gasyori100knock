package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
)

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
			r8 := uint8((float64(r32) * 0xFF) / 0xFFFF)
			g8 := uint8((float64(g32) * 0xFF) / 0xFFFF)
			b8 := uint8((float64(b32) * 0xFF) / 0xFFFF)

			var graycolor color.Gray
			// カラーからグレースケールに変換
			grayfloat64 := 0.2126*float64(r8) + 0.7152*float64(g8) + 0.0722*float64(b8)
			graycolor.Y = uint8(grayfloat64)
			grayimg.Set(width, height, graycolor)
		}
	}

	for y := 0; y < grayimg.Bounds().Size().Y; y++ {
		for x := 0; x < grayimg.Bounds().Size().X; x++ {
		}
	}

	// 縦方向Prewittフィルタを作成
	prewittFilterV := [3][3]int{
		{-1, -1, -1},
		{0, 0, 0},
		{1, 1, 1}}

	// 横方向Prewittフィルタを作成
	prewittFilterH := [3][3]int{
		{-1, 0, 1},
		{-1, 0, 1},
		{-1, 0, 1}}

	// Prewittフィルタ適用済画像保存先を作成
	prewittImgV := image.NewGray(grayimg.Bounds())
	prewittImgH := image.NewGray(grayimg.Bounds())
	filterSize := 3

	for y := 0; y < grayimg.Bounds().Size().Y; y++ {
		for x := 0; x < grayimg.Bounds().Size().X; x++ {
			// 対象ピクセルを中心とした3x3ピクセルの画素値の最大値と最小値を取得する
			var prewittV, prewittH int
			for filterY := 0; filterY < filterSize; filterY++ {
				for filterX := 0; filterX < filterSize; filterX++ {
					// 対象ピクセルの位置を計算する
					srcPointX := filterX + x - 1
					srcPointY := filterY + y - 1
					var pixVal int
					if (srcPointX < 0) || (srcPointX >= grayimg.Bounds().Size().X) ||
						(srcPointY < 0) || (srcPointY >= grayimg.Bounds().Size().Y) {
						// 0パディング
						pixVal = 0
					} else {
						pixVal = int(grayimg.GrayAt(srcPointX, srcPointY).Y)
					}

					// Sobelフィルタ畳み込み
					prewittV += pixVal * prewittFilterV[filterY][filterX]
					prewittH += pixVal * prewittFilterH[filterY][filterX]
				}
			}

			// fmt.Printf("%04d ", prewittV)
			if prewittV < 0 {
				prewittV = 0
			}

			if prewittV > 255 {
				prewittV = 255
			}

			if prewittH < 0 {
				prewittH = 0
			}

			if prewittH > 255 {
				prewittH = 255
			}

			var grayColorPrewittV, grayColorPrewittH color.Gray
			grayColorPrewittV.Y = uint8(prewittV)
			prewittImgV.Set(x, y, grayColorPrewittV)

			grayColorPrewittH.Y = uint8(prewittH)
			prewittImgH.Set(x, y, grayColorPrewittH)
		}

	}

	prewittFilterVImgFile, err := os.Create("./answer_16_v.jpg")
	defer prewittFilterVImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(prewittFilterVImgFile, prewittImgV, &jpeg.Options{100})

	prewittFilterHImgFile, err := os.Create("./answer_16_h.jpg")
	defer prewittFilterHImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(prewittFilterHImgFile, prewittImgH, &jpeg.Options{100})
}
