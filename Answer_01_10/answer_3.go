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
			grayuint8 := uint8(grayfloat64)
			// 2値化
			if grayuint8 < 128 {
				grayuint8 = 0
			} else {
				grayuint8 = 255
			}
			graycolor.Y = grayuint8
			grayimg.Set(width, height, graycolor)
		}
	}
	grayfile, err := os.Create("./answer_3.jpg")
	defer grayfile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(grayfile, grayimg, &jpeg.Options{100})

}
