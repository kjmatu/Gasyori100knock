package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
)

func rgb2ycbcr(r, g, b float64) (y, cb, cr float64) {
	y = 0.299*r + 0.5870*g + 0.114*b
	cb = -0.1687*r - 0.3313*g + 0.5*b + 128
	cr = 0.5*r - 0.4187*g - 0.0813*b + 128
	return y, cb, cr
}

func ycbcr2rgb(y, cb, cr float64) (r, g, b uint8) {
	r = uint8(y + (cr-128)*1.402)
	g = uint8(y - (cb-128)*0.3441 - (cr-128)*0.7139)
	b = uint8(y + (cb-128)*1.7718)
	return r, g, b
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

	H := jimg.Bounds().Size().Y
	W := jimg.Bounds().Size().X

	darkImg := image.NewRGBA(jimg.Bounds())

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r32, g32, b32, _ := jimg.At(x, y).RGBA()
			r := float64(r32*0xFF) / 0xFFFF
			g := float64(g32*0xFF) / 0xFFFF
			b := float64(b32*0xFF) / 0xFFFF
			yc, cb, cr := rgb2ycbcr(r, g, b)
			yc *= 0.7
			ruint8, guint8, buint8 := ycbcr2rgb(yc, cb, cr)
			darkColor := color.RGBA{ruint8, guint8, buint8, 255}
			darkImg.Set(x, y, darkColor)
		}
	}

	darkFile, err := os.Create("./answer_39.jpg")
	defer darkFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(darkFile, darkImg, &jpeg.Options{100})

}
