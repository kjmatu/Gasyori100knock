package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
)

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

	scaleBounds := jimg.Bounds()
	scaleBounds.Max.X = int(float64(scaleBounds.Max.X) * scale)
	scaleBounds.Max.Y = int(float64(scaleBounds.Max.Y) * scale)
	fmt.Println(scaleBounds)
	nearestNeighborImg := image.NewRGBA(scaleBounds)

	for height := 0; height < nearestNeighborImg.Bounds().Size().Y; height++ {
		for width := 0; width < nearestNeighborImg.Bounds().Size().X; width++ {
			scaleX := int(float64(width) / scale)
			scaleY := int(float64(height) / scale)

			nearestNeighborImg.Set(width, height, jimg.At(scaleX, scaleY))
		}
	}

	nearestNeighborImgFile, err := os.Create("./answer_25.jpg")
	defer nearestNeighborImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(nearestNeighborImgFile, nearestNeighborImg, &jpeg.Options{100})

}
