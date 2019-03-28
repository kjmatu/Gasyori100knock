package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
)

func main() {
	imageFile1, err := os.Open("./../Question_51_60/imori.jpg")
	defer imageFile1.Close()
	if err != nil {
		log.Fatal(err)
	}
	image1, err := jpeg.Decode(imageFile1)
	if err != nil {
		log.Fatal(err)
	}

	imageFile2, err := os.Open("./../Question_51_60/thorino.jpg")
	defer imageFile2.Close()
	if err != nil {
		log.Fatal(err)
	}
	image2, err := jpeg.Decode(imageFile2)
	if err != nil {
		log.Fatal(err)
	}

	blendImage := image.NewRGBA(image1.Bounds())
	alpha := 0.6
	for y := 0; y < blendImage.Bounds().Size().Y; y++ {
		for x := 0; x < blendImage.Bounds().Size().X; x++ {
			r1, g1, b1, _ := image1.At(x, y).RGBA()
			r2, g2, b2, _ := image2.At(x, y).RGBA()
			blendR := (float64(r1)*alpha + float64(r2)*(1-alpha)) * 0xFF / 0xFFFF
			if blendR > 255 {
				blendR = 255
			}
			blendG := (float64(g1)*alpha + float64(g2)*(1-alpha)) * 0xFF / 0xFFFF
			if blendG > 255 {
				blendG = 255
			}
			blendB := (float64(b1)*alpha + float64(b2)*(1-alpha)) * 0xFF / 0xFFFF
			if blendB > 255 {
				blendB = 255
			}
			blendColor := color.RGBA{uint8(blendR), uint8(blendG), uint8(blendB), 255}
			blendImage.Set(x, y, blendColor)
		}
	}
	blendFile, err := os.Create("./answer_60.jpg")
	defer blendFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(blendFile, blendImage, &jpeg.Options{100})

}
