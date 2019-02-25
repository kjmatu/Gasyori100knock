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
	maxPoolingImg := image.NewRGBA(jimg.Bounds())
	gridSize := 8
	for height := 0; height < jimg.Bounds().Size().Y; height += gridSize {
		for width := 0; width < jimg.Bounds().Size().X; width += gridSize {

			var maxR8, maxG8, maxB8, maxA8 float64
			maxR8 = 0.0
			maxG8 = 0.0
			maxB8 = 0.0
			maxA8 = 0.0
			// 8x8pixel画像の各RGBの最大値を計算する
			for grigHeight := height; grigHeight < height+gridSize; grigHeight++ {
				for gridWidth := width; gridWidth < width+gridSize; gridWidth++ {
					gridImageColor := jimg.At(gridWidth, grigHeight)
					r32, g32, b32, a32 := gridImageColor.RGBA()
					// 32bit画像から8bit画像に変換
					r8 := (float64(r32) / 0xFFFF) * 0xFF
					if maxR8 < r8 {
						maxR8 = r8
					}
					g8 := (float64(g32) / 0xFFFF) * 0xFF
					if maxG8 < g8 {
						maxG8 = g8
					}
					b8 := (float64(b32) / 0xFFFF) * 0xFF
					if maxB8 < b8 {
						maxB8 = b8
					}
					a8 := (float64(a32) / 0xFFFF) * 0xFF
					if maxA8 < a8 {
						maxA8 = a8
					}
				}
			}

			// maxP算したRGBの最大値ををpoolingImgにセットする
			for poolingHeight := height; poolingHeight < height+gridSize; poolingHeight++ {
				for poolingWidth := width; poolingWidth < width+gridSize; poolingWidth++ {
					var poolingColor color.NRGBA
					poolingColor.R = uint8(maxR8)
					poolingColor.G = uint8(maxG8)
					poolingColor.B = uint8(maxB8)
					poolingColor.A = uint8(maxA8)

					maxPoolingImg.Set(poolingWidth, poolingHeight, poolingColor)
				}
			}
		}
	}
	maxPoolingColorFile, err := os.Create("./answer_8.jpg")
	defer maxPoolingColorFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(maxPoolingColorFile, maxPoolingImg, &jpeg.Options{100})

}
