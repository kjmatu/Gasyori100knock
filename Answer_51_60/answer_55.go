package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
)

// HLine draws a horizontal line
func hLine(img *image.NRGBA, x1, y, x2 int, col color.Color) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// VLine draws a veritcal line
func vLine(img *image.NRGBA, x, y1, y2 int, col color.Color) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func rect(img *image.NRGBA, x1, y1, x2, y2 int, col color.Color) {
	hLine(img, x1, y1, x2, col)
	hLine(img, x1, y2, x2, col)
	vLine(img, x1, y1, y2, col)
	vLine(img, x2, y1, y2, col)
}

func main() {
	srcFile, err := os.Open("./../Question_51_60/imori.jpg")
	defer srcFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	srcImage, err := jpeg.Decode(srcFile)
	if err != nil {
		log.Fatal(err)
	}

	targetFile, err := os.Open("./../Question_51_60/imori_part.jpg")
	defer targetFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	targetImage, err := jpeg.Decode(targetFile)
	if err != nil {
		log.Fatal(err)
	}

	W := srcImage.Bounds().Size().X
	H := srcImage.Bounds().Size().Y

	w := targetImage.Bounds().Size().X
	h := targetImage.Bounds().Size().Y

	// テンプレートマッチング SAD(Sum of Absolute Difference)
	sad := math.MaxFloat64
	var detectPoint image.Point
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			var localSAD float64
			for v := 0; v < h; v++ {
				for u := 0; u < w; u++ {
					refX := x + u
					refY := y + v
					srcR, srcG, srcB, _ := srcImage.At(refX, refY).RGBA()
					targetR, targetG, targetB, _ := targetImage.At(u, v).RGBA()
					localSAD += math.Abs(float64(srcR) - float64(targetR))
					localSAD += math.Abs(float64(srcG) - float64(targetG))
					localSAD += math.Abs(float64(srcB) - float64(targetB))
				}
			}
			if localSAD < sad {
				sad = localSAD
				detectPoint = image.Point{x, y}
			}
		}
	}

	fmt.Println("Point", detectPoint)
	tamplateMatchImage := image.NewNRGBA(srcImage.Bounds())
	draw.Draw(tamplateMatchImage, srcImage.Bounds(), srcImage, srcImage.Bounds().Min, draw.Src)
	rect(tamplateMatchImage, detectPoint.X, detectPoint.Y, detectPoint.X+w, detectPoint.Y+h, color.RGBA{255, 0, 0, 255})
	tamplateMatchFile, err := os.Create("./answer_55.jpg")
	defer tamplateMatchFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(tamplateMatchFile, tamplateMatchImage, &jpeg.Options{100})

}
