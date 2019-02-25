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
	poolingImg := image.NewRGBA(jimg.Bounds())
	grigImageSize := 8
	for height := 0; height < jimg.Bounds().Size().Y; height += grigImageSize {
		for width := 0; width < jimg.Bounds().Size().X; width += grigImageSize {
			// 8x8pixel画像の各RGBの平均値を計算する
			var aveR8, aveG8, aveB8, aveA8 float64
			gridImagePixSize := grigImageSize * grigImageSize
			for gridHeight := height; gridHeight < height+grigImageSize; gridHeight++ {
				for gridWidth := width; gridWidth < width+grigImageSize; gridWidth++ {
					gridImageColor := jimg.At(gridWidth, gridHeight)
					r32, g32, b32, a32 := gridImageColor.RGBA()
					// 32bit画像から8bit画像に変換
					aveR8 += (float64(r32) / 0xFFFF) * 0xFF
					aveG8 += (float64(g32) / 0xFFFF) * 0xFF
					aveB8 += (float64(b32) / 0xFFFF) * 0xFF
					aveA8 += (float64(a32) / 0xFFFF) * 0xFF
				}
			}
			aveR8 /= float64(gridImagePixSize)
			aveG8 /= float64(gridImagePixSize)
			aveB8 /= float64(gridImagePixSize)
			aveA8 /= float64(gridImagePixSize)

			// 計算したRGBの平均値をpoolingImgにセットする
			for poolingHeight := height; poolingHeight < height+grigImageSize; poolingHeight++ {
				for poolingWidth := width; poolingWidth < width+grigImageSize; poolingWidth++ {
					var poolingColor color.NRGBA
					poolingColor.R = uint8(aveR8)
					poolingColor.G = uint8(aveG8)
					poolingColor.B = uint8(aveB8)
					poolingColor.A = uint8(aveA8)
					poolingImg.Set(poolingWidth, poolingHeight, poolingColor)
				}
			}
		}
	}
	poolingColorFile, err := os.Create("./answer_7.jpg")
	defer poolingColorFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(poolingColorFile, poolingImg, &jpeg.Options{100})
}
