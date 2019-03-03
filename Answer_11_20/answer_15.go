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

	// 縦方向Sobelフィルタを作成
	sobelFilterV := [3][3]int{
		{1, 0, -1},
		{2, 0, -2},
		{1, 0, -1}}

	// 横方向Sobelフィルタを作成
	sobelFilterH := [3][3]int{
		{1, 2, 1},
		{0, 0, 0},
		{-1, -2, -1}}

	// Sobelフィルタ適用済画像保存先を作成
	sobelImgV := image.NewGray(grayimg.Bounds())
	sobelImgH := image.NewGray(grayimg.Bounds())
	filterSize := 3

	for y := 0; y < grayimg.Bounds().Size().Y; y++ {
		for x := 0; x < grayimg.Bounds().Size().X; x++ {
			// 対象ピクセルを中心とした3x3ピクセルの画素値の最大値と最小値を取得する
			var sobelV, sobelH int
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
					sobelV += pixVal * sobelFilterV[filterY][filterX]
					sobelH += pixVal * sobelFilterH[filterY][filterX]
				}
			}

			if sobelV < 0 {
				sobelV = 0
			}
			if sobelV > 255 {
				sobelV = 255
			}

			if sobelH < 0 {
				sobelH = 0
			}
			if sobelH > 255 {
				sobelV = 255
			}

			var grayColorSobelV, grayColorSobelH color.Gray
			grayColorSobelV.Y = uint8(sobelV)
			sobelImgV.Set(x, y, grayColorSobelV)

			grayColorSobelH.Y = uint8(sobelH)
			sobelImgH.Set(x, y, grayColorSobelH)
		}
	}

	sobelFilterVImgFile, err := os.Create("./answer_15_v.jpg")
	defer sobelFilterVImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(sobelFilterVImgFile, sobelImgV, &jpeg.Options{100})

	sobelFilterHImgFile, err := os.Create("./answer_15_h.jpg")
	defer sobelFilterHImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(sobelFilterHImgFile, sobelImgH, &jpeg.Options{100})
}
