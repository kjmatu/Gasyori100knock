package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func affineRotateImg(rotateAngle, tx, ty float64, x, y int) (rotateX, rotateY int) {
	rotateRadian := (rotateAngle * math.Pi) / 180.0
	affineMatrix := [3][3]float64{{math.Cos(rotateRadian), -math.Sin(rotateRadian), tx}, {math.Sin(rotateRadian), math.Cos(rotateRadian), ty}, {0.0, 0.0, 1.0}}
	moveArray := [3]float64{}
	for i, rowArray := range affineMatrix {
		moveArray[i] = float64(x)*rowArray[0] + float64(y)*rowArray[1] + 1.0*rowArray[2]
	}
	rotateX = int(moveArray[0])
	rotateY = int(moveArray[1])
	return rotateX, rotateY
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

	rotateImg := image.NewRGBA(jimg.Bounds())

	for height := 0; height < rotateImg.Bounds().Size().Y; height++ {
		for width := 0; width < rotateImg.Bounds().Size().X; width++ {
			x, y := affineRotateImg(30.0, 0.0, 0.0, width, height)
			if x < 0 || x >= jimg.Bounds().Size().X || y < 0 || y >= jimg.Bounds().Size().Y {
				black := color.RGBA{0, 0, 0, 255}
				rotateImg.Set(width, height, black)
			} else {
				rotateImg.Set(width, height, jimg.At(x, y))
			}
		}
	}

	rotateImgFile, err := os.Create("./answer_30_1.jpg")
	defer rotateImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(rotateImgFile, rotateImg, &jpeg.Options{100})

	rotateCenterImg := image.NewRGBA(jimg.Bounds())

	for height := 0; height < rotateCenterImg.Bounds().Size().Y; height++ {
		for width := 0; width < rotateCenterImg.Bounds().Size().X; width++ {
			x, y := affineRotateImg(0.0, -float64(rotateImg.Bounds().Size().X)/2.0, -float64(rotateImg.Bounds().Size().Y)/2.0, width, height)

			x, y = affineRotateImg(30.0, float64(rotateImg.Bounds().Size().X)/2.0, float64(rotateImg.Bounds().Size().Y)/2.0, x, y)

			if x < 0 || x >= jimg.Bounds().Size().X || y < 0 || y >= jimg.Bounds().Size().Y {
				black := color.RGBA{0, 0, 0, 255}
				rotateCenterImg.Set(width, height, black)
			} else {
				rotateCenterImg.Set(width, height, jimg.At(x, y))
			}
		}
	}

	rotateCenterImgFile, err := os.Create("./answer_30_2.jpg")
	defer rotateCenterImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(rotateCenterImgFile, rotateCenterImg, &jpeg.Options{100})

}
