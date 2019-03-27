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

func ncc() {

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

	// テンプレートマッチング 正規化相互相関 NCC(Normalized Cross Correlation)
	ncc := -math.MaxFloat64
	var detectPoint image.Point
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			var localNcc float64
			var localNccR float64
			var localNccG float64
			var localNccB float64
			sumIrR := 0.0
			sumTrR := 0.0
			sumIrG := 0.0
			sumTrG := 0.0
			sumIrB := 0.0
			sumTrB := 0.0
			for v := 0; v < h; v++ {
				for u := 0; u < w; u++ {
					refX := x + u
					refY := y + v
					srcR, srcG, srcB, _ := srcImage.At(refX, refY).RGBA()
					targetR, targetG, targetB, _ := targetImage.At(u, v).RGBA()
					localNccR += float64(srcR * targetR)
					sumIrR += float64(srcR * srcR)
					sumTrR += float64(targetR * targetR)
					localNccG += float64(srcG * targetG)
					sumIrG += float64(srcG * srcG)
					sumTrG += float64(targetG * targetG)
					localNccB += float64(srcB * targetB)
					sumIrB += float64(srcB * srcB)
					sumTrB += float64(targetB * targetB)
				}
			}
			localNccR /= (math.Sqrt(sumIrR) * math.Sqrt(sumTrR))
			localNccG /= (math.Sqrt(sumIrG) * math.Sqrt(sumTrG))
			localNccB /= (math.Sqrt(sumIrB) * math.Sqrt(sumTrB))
			localNcc = (localNccR + localNccG + localNccB) / 3
			if localNcc > ncc {
				ncc = localNcc
				detectPoint = image.Point{x, y}
			}
		}
	}

	fmt.Println("Point", detectPoint)
	tamplateMatchImage := image.NewNRGBA(srcImage.Bounds())
	draw.Draw(tamplateMatchImage, srcImage.Bounds(), srcImage, srcImage.Bounds().Min, draw.Src)
	rect(tamplateMatchImage, detectPoint.X, detectPoint.Y, detectPoint.X+w, detectPoint.Y+h, color.RGBA{255, 0, 0, 255})
	tamplateMatchFile, err := os.Create("./answer_56.jpg")
	defer tamplateMatchFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(tamplateMatchFile, tamplateMatchImage, &jpeg.Options{100})

}
