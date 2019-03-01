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

	// 縦方向微分フィルタを作成
	diffFilterV := [3][3]int{
		{0, -1, 0},
		{0, 1, 0},
		{0, 0, 0}}

	// 横方向微分フィルタを作成
	diffFilterH := [3][3]int{
		{0, 0, 0},
		{-1, 1, 0},
		{0, 0, 0}}

	// 微分フィルタ適用済画像保存先を作成
	diffImgV := image.NewGray(grayimg.Bounds())
	diffImgH := image.NewGray(grayimg.Bounds())
	filterSize := 3

	for y := 0; y < grayimg.Bounds().Size().Y; y++ {
		for x := 0; x < grayimg.Bounds().Size().X; x++ {
			// 対象ピクセルを中心とした3x3ピクセルの画素値の最大値と最小値を取得する
			var diffV, diffH int
			for filterX := 0; filterX < filterSize; filterX++ {
				for filterY := 0; filterY < filterSize; filterY++ {
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

					// 微分フィルタ畳み込み
					diffV += pixVal * diffFilterV[filterY][filterX]
					diffH += pixVal * diffFilterH[filterY][filterX]
				}
			}
			if diffV < 0 {
				diffV = 0
			}
			if diffH < 0 {
				diffH = 0
			}

			var grayColorDiffV, grayColorDiffH color.Gray
			grayColorDiffV.Y = uint8(diffV)
			diffImgV.Set(x, y, grayColorDiffV)

			grayColorDiffH.Y = uint8(diffH)
			diffImgH.Set(x, y, grayColorDiffH)
		}
	}

	diffFilterVImgFile, err := os.Create("./answer_14_v.jpg")
	defer diffFilterVImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(diffFilterVImgFile, diffImgV, &jpeg.Options{100})

	diffFilterHImgFile, err := os.Create("./answer_14_h.jpg")
	defer diffFilterHImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(diffFilterHImgFile, diffImgH, &jpeg.Options{100})
}
