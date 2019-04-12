package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
)

// HLine draws a horizontal line
func hLine(img *image.RGBA, x1, y, x2 int, col color.RGBA) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// VLine draws a veritcal line
func vLine(img *image.RGBA, x, y1, y2 int, col color.RGBA) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func drawRectangle(img *image.RGBA, rect image.Rectangle, col color.RGBA) {
	x1 := rect.Min.X
	y1 := rect.Min.Y
	x2 := rect.Max.X
	y2 := rect.Max.Y
	hLine(img, x1, y1, x2, col)
	hLine(img, x1, y2, x2, col)
	vLine(img, x1, y1, y2, col)
	vLine(img, x2, y1, y2, col)
}

func calcIoU(rect1 image.Rectangle, rect2 image.Rectangle) float64 {
	area1 := rect1.Dx() * rect1.Dy()
	area2 := rect2.Dx() * rect2.Dy()

	overlapRect := rect1.Intersect(rect2)
	overlapArea := overlapRect.Dx() * overlapRect.Dy()

	iou := float64(overlapArea) / float64(area1+area2-overlapArea)

	return iou
}

func randomClip(srcImage *image.RGBA) (image.Image, image.Rectangle) {
	width, height := 60, 60
	H := srcImage.Bounds().Size().Y
	W := srcImage.Bounds().Size().X
	randomX := int(rand.Int63n(int64(W - 60)))
	randomY := int(rand.Int63n(int64(H - 60)))

	clipRect := image.Rect(randomX, randomY, randomX+width, randomY+height)
	clipImage := srcImage.SubImage(clipRect)
	return clipImage, clipRect
}

func main() {
	file, err := os.Open("./../Question_91_100/imori_1.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpegImage, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	imoriImage := image.NewRGBA(jpegImage.Bounds())
	draw.Draw(imoriImage, jpegImage.Bounds(), jpegImage, jpegImage.Bounds().Min, draw.Src)

	imoriFaceRect := image.Rect(47, 41, 129, 103)

	rand.Seed(0)
	for i := 0; i < 200; i++ {
		_, clipRect := randomClip(imoriImage)
		iou := calcIoU(imoriFaceRect, clipRect)
		if iou > 0.5 {
			drawRectangle(imoriImage, clipRect, color.RGBA{255, 0, 0, 255})
		} else {
			drawRectangle(imoriImage, clipRect, color.RGBA{0, 0, 255, 255})
		}
	}

	randomClipFIle, err := os.Create("./answer_94.jpg")
	defer randomClipFIle.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(randomClipFIle, imoriImage, &jpeg.Options{100})

}
