package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func gammaCorrection(pixVal, c, g float64) uint8 {
	pixVal /= 255.0
	correctValue := math.Pow((pixVal / c), (1.0 / g))
	return uint8(correctValue * 255.0)
}

func main() {
	file, err := os.Open("./../Question_21_30/imori_gamma.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	gammaCorrectionImg := image.NewRGBA(jimg.Bounds())
	c := 1.0
	gamma := 2.2
	for height := 0; height < gammaCorrectionImg.Bounds().Size().Y; height++ {
		for width := 0; width < gammaCorrectionImg.Bounds().Size().X; width++ {
			r, g, b, _ := jimg.At(width, height).RGBA()
			r8 := float64(r*0xFF) / 0xFFFF
			g8 := float64(g*0xFF) / 0xFFFF
			b8 := float64(b*0xFF) / 0xFFFF

			var gammaCorrectionColor color.RGBA
			gammaCorrectionColor.R = gammaCorrection(r8, c, gamma)
			gammaCorrectionColor.G = gammaCorrection(g8, c, gamma)
			gammaCorrectionColor.B = gammaCorrection(b8, c, gamma)
			gammaCorrectionColor.A = 255
			gammaCorrectionImg.Set(width, height, gammaCorrectionColor)
		}
	}

	gammaCorrectionImgFile, err := os.Create("./answer_24.jpg")
	defer gammaCorrectionImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(gammaCorrectionImgFile, gammaCorrectionImg, &jpeg.Options{100})

}
