package main

import (
	"fmt"
	"image"
	"math"
)

func calcIoU(rect1 image.Rectangle, rect2 image.Rectangle) float64 {
	area1 := rect1.Dx() * rect1.Dy()
	area2 := rect2.Dx() * rect2.Dy()

	overlapX1 := int(math.Max(float64(rect1.Min.X), float64(rect2.Min.X)))
	overlapY1 := int(math.Max(float64(rect1.Min.Y), float64(rect2.Min.Y)))

	overlapX2 := int(math.Min(float64(rect1.Max.X), float64(rect2.Max.X)))
	overlapY2 := int(math.Min(float64(rect1.Max.Y), float64(rect2.Max.Y)))

	overlapArea := (overlapX2 - overlapX1) * (overlapY2 - overlapY1)

	iou := float64(overlapArea) / float64(area1+area2-overlapArea)

	return iou
}

func main() {
	rect1 := image.Rect(50, 50, 150, 150)
	rect2 := image.Rect(60, 60, 170, 160)
	fmt.Println(calcIoU(rect1, rect2))
}
