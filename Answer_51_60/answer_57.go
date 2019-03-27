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

	// テンプレートマッチング ZNCC(Zero means Normalized Cross Correlation)
	// 零平均正規化相互相関
	meanIr := 0.0
	meanIg := 0.0
	meanIb := 0.0
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			srcR, srcG, srcB, _ := srcImage.At(x, y).RGBA()
			meanIr += float64(srcR)
			meanIg += float64(srcG)
			meanIb += float64(srcB)
		}
	}
	meanIr /= float64(H * W)
	meanIg /= float64(H * W)
	meanIb /= float64(H * W)

	meanTr := 0.0
	meanTg := 0.0
	meanTb := 0.0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			targetR, targetG, targetB, _ := targetImage.At(x, y).RGBA()
			meanTr += float64(targetR)
			meanTg += float64(targetG)
			meanTb += float64(targetB)
		}
	}
	meanTr /= float64(h * w)
	meanTg /= float64(h * w)
	meanTb /= float64(h * w)

	zncc := -math.MaxFloat64
	var detectPoint image.Point
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			var localZncc float64
			var localZnccR float64
			var localZnccG float64
			var localZnccB float64
			sumIr := 0.0
			sumTr := 0.0

			sumIg := 0.0
			sumTg := 0.0

			sumIb := 0.0
			sumTb := 0.0
			for v := 0; v < h; v++ {
				for u := 0; u < w; u++ {
					refX := x + u
					refY := y + v
					srcR, srcG, srcB, _ := srcImage.At(refX, refY).RGBA()
					targetR, targetG, targetB, _ := targetImage.At(u, v).RGBA()
					localZnccR += (float64(srcR) - meanIr) * (float64(targetR) - meanTr)
					sumIr += math.Pow((float64(srcR) - meanIr), 2)
					sumTr += math.Pow((float64(targetR) - meanTr), 2)

					localZnccG += (float64(srcG) - meanIg) * (float64(targetG) - meanTg)
					sumIg += math.Pow((float64(srcG) - meanIg), 2)
					sumTg += math.Pow((float64(targetG) - meanTg), 2)

					localZnccB += (float64(srcB) - meanIb) * (float64(targetB) - meanTb)
					sumIb += math.Pow((float64(srcB) - meanIb), 2)
					sumTb += math.Pow((float64(targetB) - meanTb), 2)
				}
			}
			localZnccR /= (math.Sqrt(sumIr) * math.Sqrt(sumTr))
			localZnccG /= (math.Sqrt(sumIg) * math.Sqrt(sumTg))
			localZnccB /= (math.Sqrt(sumIb) * math.Sqrt(sumTb))
			localZncc = (localZnccR + localZnccG + localZnccB) / 3
			if localZncc > zncc {
				zncc = localZncc
				detectPoint = image.Point{x, y}
			}
		}
	}

	fmt.Println("Point", detectPoint)
	tamplateMatchImage := image.NewNRGBA(srcImage.Bounds())
	draw.Draw(tamplateMatchImage, srcImage.Bounds(), srcImage, srcImage.Bounds().Min, draw.Src)
	rect(tamplateMatchImage, detectPoint.X, detectPoint.Y, detectPoint.X+w, detectPoint.Y+h, color.RGBA{255, 0, 0, 255})
	tamplateMatchFile, err := os.Create("./answer_57.jpg")
	defer tamplateMatchFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(tamplateMatchFile, tamplateMatchImage, &jpeg.Options{100})

}
