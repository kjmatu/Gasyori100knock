package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
)

// モルフォロジー処理(膨張)
func dialation(binImage *image.Gray) *image.Gray {
	dialationImage := image.NewGray(binImage.Bounds())
	filter := [3][3]int{
		{0, 1, 0},
		{1, 0, 1},
		{0, 1, 0}}
	for y := 1; y < binImage.Bounds().Size().Y-1; y++ {
		for x := 1; x < binImage.Bounds().Size().X-1; x++ {
			targetPix := binImage.GrayAt(x, y).Y
			if targetPix == 0 {
				sum := 0
				for j, row := range filter {
					for k, value := range row {
						sum += int(binImage.GrayAt(x+k-1, y+j-1).Y) * value
					}
				}

				if sum >= 255 {
					dialationImage.Set(x, y, color.Gray{255})
				} else {
					dialationImage.Set(x, y, binImage.GrayAt(x, y))
				}
			} else {
				dialationImage.Set(x, y, binImage.GrayAt(x, y))
			}
			// fmt.Printf("[%d][%d] %d ", x, y, dialationImage.GrayAt(x, y).Y)
		}
	}
	return dialationImage
}

// モルフォロジー処理(縮小)
func erosion(binImage *image.Gray) *image.Gray {
	erosionImage := image.NewGray(binImage.Bounds())
	filter := [3][3]int{
		{0, 1, 0},
		{1, 0, 1},
		{0, 1, 0}}
	for y := 1; y < binImage.Bounds().Size().Y-1; y++ {
		for x := 1; x < binImage.Bounds().Size().X-1; x++ {
			targetPix := binImage.GrayAt(x, y).Y
			if targetPix == 255 {
				sum := 0
				for j, row := range filter {
					for k, value := range row {
						sum += int(binImage.GrayAt(x+k-1, y+j-1).Y) * value
					}
				}
				if sum < 255*4 {
					erosionImage.Set(x, y, color.Gray{0})
				} else {
					erosionImage.Set(x, y, binImage.GrayAt(x, y))
				}
			} else {
				erosionImage.Set(x, y, binImage.GrayAt(x, y))
			}
		}
	}
	return erosionImage
}

func main() {
	cannyEdgeFile, err := os.Open("./../Question_41_50/answers/answer_43.jpg")
	defer cannyEdgeFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	cannyEdgeImage, err := jpeg.Decode(cannyEdgeFile)
	if err != nil {
		log.Fatal(err)
	}

	cannyEdgeGrayImage := image.NewGray(cannyEdgeImage.Bounds())
	for y := 0; y < cannyEdgeGrayImage.Bounds().Size().Y; y++ {
		for x := 0; x < cannyEdgeGrayImage.Bounds().Size().X; x++ {
			r, _, _, _ := cannyEdgeImage.At(x, y).RGBA()
			cannyEdgeGrayImage.Set(x, y, color.Gray{uint8(r * 0xFF / 0xFFFF)})
		}
	}

	closingImage := dialation(cannyEdgeGrayImage)
	closingFile, err := os.Create("./answer_50_dia.jpg")
	defer closingFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(closingFile, closingImage, &jpeg.Options{100})

	closingImage = erosion(closingImage)

	closingFile, err = os.Create("./answer_50.jpg")
	defer closingFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(closingFile, closingImage, &jpeg.Options{100})

}
