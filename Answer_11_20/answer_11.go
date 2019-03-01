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

	// 平滑化画像を作成
	averageImg := image.NewRGBA(jimg.Bounds())
	filterSize := 3
	for width := 0; width < jimg.Bounds().Size().X; width++ {
		for height := 0; height < jimg.Bounds().Size().Y; height++ {

			averageR := 0.0
			averageG := 0.0
			averageB := 0.0
			// 対象ピクセルを中心とした3x3ピクセルの画素値を配列に格納する
			for filterWidth := 0; filterWidth < filterSize; filterWidth++ {
				for filterHeight := 0; filterHeight < filterSize; filterHeight++ {
					// 対象画像の位置を計算する
					srcPointX := filterWidth + width - 1
					srcPointY := filterHeight + height - 1
					var r8, g8, b8 float64
					if (srcPointX < 0) || (srcPointX >= jimg.Bounds().Size().X) ||
						(srcPointY < 0) || (srcPointY >= jimg.Bounds().Size().Y) {
						r8 = 0.0
						g8 = 0.0
						b8 = 0.0
					} else {
						r32, g32, b32, _ := jimg.At(srcPointX, srcPointY).RGBA()
						r8 = (float64(r32) * 0xFF) / 0xFFFF
						g8 = (float64(g32) * 0xFF) / 0xFFFF
						b8 = (float64(b32) * 0xFF) / 0xFFFF
					}
					averageR += r8
					averageG += g8
					averageB += b8
				}
			}

			var averageColor color.NRGBA
			averageColor.R = uint8((averageR / float64(filterSize*filterSize)))
			averageColor.G = uint8((averageG / float64(filterSize*filterSize)))
			averageColor.B = uint8((averageB / float64(filterSize*filterSize)))
			averageColor.A = uint8(255)
			averageImg.Set(width, height, averageColor)
		}
	}

	averageImgFile, err := os.Create("./answer_11.jpg")
	defer averageImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(averageImgFile, averageImg, &jpeg.Options{100})
}
