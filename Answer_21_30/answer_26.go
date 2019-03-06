package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func uint32Touint8(pixVal uint32) uint8 {
	return uint8(float64(pixVal*0xFF) / 0xFFFF)
}

func biLinearElm(pixVal00, pixVal01, pixVal10, pixVal11 uint32, dx, dy float64) uint8 {
	pixValUint8_00 := uint32Touint8(pixVal00)
	pixValUint8_01 := uint32Touint8(pixVal01)
	pixValUint8_10 := uint32Touint8(pixVal10)
	pixValUint8_11 := uint32Touint8(pixVal11)
	biLinearVal := (1.0-dx)*(1.0-dy)*float64(pixValUint8_00) + dx*(1-dy)*float64(pixValUint8_10) + (1.0-dx)*dy*float64(pixValUint8_01) + dx*dy*float64(pixValUint8_11)
	if biLinearVal > 255.0 {
		biLinearVal = 255.0
	}
	return uint8(biLinearVal)
}

func biLinear(I00, I01, I10, I11 color.Color, dx, dy float64) color.Color {
	var biLinearColor color.RGBA
	r00, g00, b00, _ := I00.RGBA()
	r01, g01, b01, _ := I01.RGBA()
	r10, g10, b10, _ := I10.RGBA()
	r11, g11, b11, _ := I11.RGBA()
	biLinearColor.R = biLinearElm(r00, r01, r10, r11, dx, dy)
	biLinearColor.G = biLinearElm(g00, g01, g10, g11, dx, dy)
	biLinearColor.B = biLinearElm(b00, b01, b10, b11, dx, dy)
	biLinearColor.A = 255
	return biLinearColor
}

func main() {
	file, err := os.Open("./../Question_21_30/imori.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	scale := 1.5

	// 拡大縮小画像を作成
	scaleBounds := jimg.Bounds()
	scaleBounds.Max.X = int(float64(scaleBounds.Max.X) * scale)
	scaleBounds.Max.Y = int(float64(scaleBounds.Max.Y) * scale)
	biLinearImg := image.NewRGBA(scaleBounds)

	for height := 0; height < biLinearImg.Bounds().Size().Y; height++ {
		for width := 0; width < biLinearImg.Bounds().Size().X; width++ {
			// 拡大画像のピクセル位置に対応する元画像位置を計算
			srcX := float64(width) / scale
			srcY := float64(height) / scale

			// 上記で計算した位置の周囲4点の画素値を取得
			I00 := jimg.At(int(srcX), int(srcY))
			I10 := jimg.At(int(srcX+1), int(srcY))
			I01 := jimg.At(int(srcX), int(srcY+1))
			I11 := jimg.At(int(srcX+1), int(srcY+1))

			// 周囲4点とピクセル対応点の距離を計算
			dx := srcX - math.Floor(srcX)
			dy := srcY - math.Floor(srcY)

			// Bi Linear補完
			biLinearColor := biLinear(I00, I01, I10, I11, dx, dy)
			biLinearImg.Set(width, height, biLinearColor)

			// 画像の境界を補完するときは最近傍補完を行う
			if width == biLinearImg.Bounds().Size().X-1 {
				biLinearImg.Set(width, height, jimg.At(int(srcX), int(srcY)))
			}
			if height == biLinearImg.Bounds().Size().Y-1 {
				biLinearImg.Set(width, height, jimg.At(int(srcX), int(srcY)))
			}
		}
	}

	biLinearImgFile, err := os.Create("./answer_26.jpg")
	defer biLinearImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(biLinearImgFile, biLinearImg, &jpeg.Options{100})

}
