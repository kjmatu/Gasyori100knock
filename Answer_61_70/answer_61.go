package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

func main() {
	connectivityNumberFile, err := os.Open("./../Question_61_70/renketsu.png")
	defer connectivityNumberFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	cnImage, err := png.Decode(connectivityNumberFile)
	if err != nil {
		log.Fatal(err)
	}

	// カラー画像をグレイスケール画像に変換し、画素値が0以上なら1をセットする
	cnGrayImage := image.NewGray(cnImage.Bounds())
	for y := 0; y < cnImage.Bounds().Size().Y; y++ {
		for x := 0; x < cnImage.Bounds().Size().X; x++ {
			c := color.GrayModel.Convert(cnImage.At(x, y))
			gray, _ := c.(color.Gray)
			if gray.Y > 0 {
				cnGrayImage.Set(x, y, color.Gray{1})
			}
		}
	}

	cnColorImage := image.NewNRGBA(cnImage.Bounds())

	for y := 0; y < cnGrayImage.Bounds().Size().Y; y++ {
		for x := 0; x < cnGrayImage.Bounds().Size().X; x++ {
			x0 := cnGrayImage.GrayAt(x, y).Y
			if x0 == 0 {
				cnColorImage.Set(x, y, color.RGBA{0, 0, 0, 255})
				continue
			}

			leftIndex := int(math.Max(float64(x-1), 0))
			rightIndex := int(math.Min(float64(x+1), float64(cnImage.Bounds().Size().X)-1))
			upIndex := int(math.Max(float64(y-1), 0))
			downIndex := int(math.Min(float64(y+1), float64(cnImage.Bounds().Size().Y)-1))
			x1 := int(cnGrayImage.GrayAt(rightIndex, y).Y)
			x2 := int(cnGrayImage.GrayAt(rightIndex, upIndex).Y)
			x3 := int(cnGrayImage.GrayAt(x, upIndex).Y)
			x4 := int(cnGrayImage.GrayAt(leftIndex, upIndex).Y)
			x5 := int(cnGrayImage.GrayAt(leftIndex, y).Y)
			x6 := int(cnGrayImage.GrayAt(leftIndex, downIndex).Y)
			x7 := int(cnGrayImage.GrayAt(x, downIndex).Y)
			x8 := int(cnGrayImage.GrayAt(rightIndex, downIndex).Y)

			s := (x1 - x1*x2*x3) + (x3 - x3*x4*x5) + (x5 - x5*x6*x7) + (x7 - x7*x8*x1)

			if s == 0 {
				cnColorImage.Set(x, y, color.RGBA{255, 0, 0, 255})
			} else if s == 1 {
				cnColorImage.Set(x, y, color.RGBA{0, 255, 0, 255})
			} else if s == 2 {
				cnColorImage.Set(x, y, color.RGBA{0, 0, 255, 255})
			} else if s == 3 {
				cnColorImage.Set(x, y, color.RGBA{0, 255, 255, 255})
			} else if s == 4 {
				cnColorImage.Set(x, y, color.RGBA{255, 0, 255, 255})
			}
		}
	}

	cnColorFile, err := os.Create("./answer_61.png")
	defer cnColorFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(cnColorFile, cnColorImage)

}
