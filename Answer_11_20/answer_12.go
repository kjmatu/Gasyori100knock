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

	// モーションフィルタを作成
	motionFilter := [3][3]float64{
		{1.0 / 3.0, 0.0, 0.0},
		{0.0, 1.0 / 3.0, 0.0},
		{0.0, 0.0, 1.0 / 3.0},
	}

	// モーション画像を作成
	motionImg := image.NewRGBA(jimg.Bounds())
	filterSize := 3
	for width := 0; width < jimg.Bounds().Size().X; width++ {
		for height := 0; height < jimg.Bounds().Size().Y; height++ {

			// 対象ピクセルを中心とした3x3ピクセルの画素値とモーションフィルタを畳み込みする
			var motionR, motionG, motionB float64
			for filterWidth := 0; filterWidth < filterSize; filterWidth++ {
				for filterHeight := 0; filterHeight < filterSize; filterHeight++ {
					// 対象画像の位置を計算する
					srcPointX := filterWidth + width - 1
					srcPointY := filterHeight + height - 1
					var r8, g8, b8 float64
					if (srcPointX < 0) || (srcPointX >= jimg.Bounds().Size().X) ||
						(srcPointY < 0) || (srcPointY >= jimg.Bounds().Size().Y) {
						// 0パディング
						r8 = 0.0
						g8 = 0.0
						b8 = 0.0
					} else {
						r32, g32, b32, _ := jimg.At(srcPointX, srcPointY).RGBA()
						r8 = (float64(r32) * 0xFF) / 0xFFFF
						g8 = (float64(g32) * 0xFF) / 0xFFFF
						b8 = (float64(b32) * 0xFF) / 0xFFFF
					}
					// モーションフィルタを畳み込み
					motionR += r8 * motionFilter[filterWidth][filterHeight]
					motionG += g8 * motionFilter[filterWidth][filterHeight]
					motionB += b8 * motionFilter[filterWidth][filterHeight]
				}
			}

			var motionColor color.NRGBA
			motionColor.R = uint8(motionR)
			motionColor.G = uint8(motionG)
			motionColor.B = uint8(motionB)
			motionColor.A = uint8(255)
			motionImg.Set(width, height, motionColor)
		}
	}

	motionImgFile, err := os.Create("./answer_12.jpg")
	defer motionImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(motionImgFile, motionImg, &jpeg.Options{100})
}
