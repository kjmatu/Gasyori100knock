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

	// Embossフィルタを作成
	laplacianFilter := [3][3]int{
		{-2, -1, 0},
		{-1, 1, 1},
		{0, 1, 2}}

	// Embossフィルタ適用済画像保存先を作成
	embossImg := image.NewGray(grayimg.Bounds())
	filterSize := 3

	for y := 0; y < grayimg.Bounds().Size().Y; y++ {
		for x := 0; x < grayimg.Bounds().Size().X; x++ {
			// 対象ピクセルを中心とした3x3ピクセルの画素値の最大値と最小値を取得する
			var filteredValue int
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

					// Embossフィルタ畳み込み
					filteredValue += pixVal * laplacianFilter[filterY][filterX]
				}
			}

			if filteredValue < 0 {
				filteredValue = 0
			}

			if filteredValue > 255 {
				filteredValue = 255
			}

			var filteredGray color.Gray
			filteredGray.Y = uint8(filteredValue)
			embossImg.Set(x, y, filteredGray)
		}
	}

	embossFilterImgFile, err := os.Create("./answer_18.jpg")
	defer embossFilterImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(embossFilterImgFile, embossImg, &jpeg.Options{100})

}
