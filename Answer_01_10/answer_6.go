package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
)

func reduceColors(pixVal uint8) uint8 {
	if (0 <= pixVal) && (pixVal < 63) {
		return 32
	} else if (63 <= pixVal) && (pixVal < 127) {
		return 96
	} else if (127 <= pixVal) && (pixVal < 191) {
		return 160
	} else if (191 <= pixVal) && (pixVal <= 255) {
		return 224
	} else {
		return 255
	}
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
	subtractiveColorImg := image.NewRGBA(jimg.Bounds())
	for height := 0; height < jimg.Bounds().Size().Y; height++ {
		for width := 0; width < jimg.Bounds().Size().X; width++ {
			ycbcr := jimg.At(width, height)
			r32, g32, b32, a32 := ycbcr.RGBA()
			// 32bit画像から8bit画像に変換
			r8 := uint8((float64(r32) / 0xFFFF) * 0xFF)
			g8 := uint8((float64(g32) / 0xFFFF) * 0xFF)
			b8 := uint8((float64(b32) / 0xFFFF) * 0xFF)
			a8 := uint8((float64(a32) / 0xFFFF) * 0xFF)

			var subtractiveColor color.NRGBA
			r8 = reduceColors(r8)
			g8 = reduceColors(g8)
			b8 = reduceColors(b8)
			subtractiveColor.R = r8
			subtractiveColor.G = g8
			subtractiveColor.B = b8
			subtractiveColor.A = a8

			subtractiveColorImg.Set(width, height, subtractiveColor)
		}
	}

	subtractiveColorFile, err := os.Create("./answer_6.jpg")
	defer subtractiveColorFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(subtractiveColorFile, subtractiveColorImg, &jpeg.Options{100})

}
