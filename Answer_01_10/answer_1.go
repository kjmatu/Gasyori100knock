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
	bgrimg := image.NewRGBA(jimg.Bounds())
	for height := 0; height < jimg.Bounds().Size().Y; height++ {
		for width := 0; width < jimg.Bounds().Size().X; width++ {
			ycbcr := jimg.At(width, height)
			r, g, b, a := ycbcr.RGBA()
			var rotatecolor color.RGBA64
			rotatecolor.R = uint16(b)
			rotatecolor.G = uint16(g)
			rotatecolor.B = uint16(r)
			rotatecolor.A = uint16(a)
			bgrimg.Set(width, height, rotatecolor)
		}
	}
	bgrfile, err := os.Create("./answer_1.jpg")
	defer bgrfile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(bgrfile, bgrimg, &jpeg.Options{100})
}
