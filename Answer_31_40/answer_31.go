package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
)

func calcXsharingSourcePos(a, tx, ty float64, x, y int) (srcX, srcY int) {
	invAffineXsharingMatrix := [3][3]float64{
		{1.0, -a, a*ty - tx},
		{0.0, 1.0, -ty},
		{0.0, 0.0, 1.0}}
	moveArray := [3]float64{}
	for i, rowArray := range invAffineXsharingMatrix {
		moveArray[i] = float64(x)*rowArray[0] + float64(y)*rowArray[1] + 1.0*rowArray[2]
	}
	srcX = int(moveArray[0])
	srcY = int(moveArray[1])
	return srcX, srcY
}

func calcYsharingSourcePos(a, tx, ty float64, x, y int) (srcX, srcY int) {
	invAffineYsharingMatrix := [3][3]float64{
		{1.0, 0.0, -tx},
		{-a, 1.0, a*ty - tx},
		{0.0, 0.0, 1.0}}
	moveArray := [3]float64{}
	for i, rowArray := range invAffineYsharingMatrix {
		moveArray[i] = float64(x)*rowArray[0] + float64(y)*rowArray[1] + 1.0*rowArray[2]
	}
	srcX = int(moveArray[0])
	srcY = int(moveArray[1])
	return srcX, srcY
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

	dx := 30
	xsharingBound := jimg.Bounds()
	xsharingBound.Max.X += dx
	xsharingImg := image.NewRGBA(xsharingBound)
	a := float64(dx) / float64(xsharingImg.Bounds().Size().Y)

	for height := 0; height < xsharingImg.Bounds().Size().Y; height++ {
		for width := 0; width < xsharingImg.Bounds().Size().X; width++ {
			srcX, srcY := calcXsharingSourcePos(a, 0.0, 0.0, width, height)
			if srcX < 0 || srcY < 0 || srcX >= jimg.Bounds().Size().X || srcY >= jimg.Bounds().Size().Y {
				xsharingImg.Set(width, height, color.RGBA{0, 0, 0, 255})
			} else {
				xsharingImg.Set(width, height, jimg.At(srcX, srcY))

			}
		}
	}

	xsharingImgFile, err := os.Create("./answer_31_1.jpg")
	defer xsharingImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(xsharingImgFile, xsharingImg, &jpeg.Options{100})

	dy := 30
	ysharingBound := jimg.Bounds()
	ysharingBound.Max.Y += dy
	ysharingImg := image.NewRGBA(ysharingBound)
	a = float64(dy) / float64(ysharingImg.Bounds().Size().X)

	for height := 0; height < ysharingImg.Bounds().Size().Y; height++ {
		for width := 0; width < ysharingImg.Bounds().Size().X; width++ {
			srcX, srcY := calcYsharingSourcePos(a, 0.0, 0.0, width, height)
			if srcX < 0 || srcY < 0 || srcX >= jimg.Bounds().Size().X || srcY >= jimg.Bounds().Size().Y {
				ysharingImg.Set(width, height, color.RGBA{0, 0, 0, 255})
			} else {
				ysharingImg.Set(width, height, jimg.At(srcX, srcY))

			}
		}
	}

	ysharingImgFile, err := os.Create("./answer_31_2.jpg")
	defer xsharingImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(ysharingImgFile, ysharingImg, &jpeg.Options{100})

	xysharingBound := jimg.Bounds()
	xysharingBound.Max.Y += dy
	xysharingBound.Max.X += dx
	xysharingImg := image.NewRGBA(xysharingBound)
	a = float64(dy) / float64(xysharingImg.Bounds().Size().X)
	b := float64(dx) / float64(xysharingImg.Bounds().Size().Y)

	for height := 0; height < xysharingImg.Bounds().Size().Y; height++ {
		for width := 0; width < xysharingImg.Bounds().Size().X; width++ {
			srcX, srcY := calcYsharingSourcePos(a, 0.0, 0.0, width, height)
			srcX, srcY = calcXsharingSourcePos(b, 0.0, 0.0, srcX, srcY)
			if srcX < 0 || srcY < 0 || srcX >= jimg.Bounds().Size().X || srcY >= jimg.Bounds().Size().Y {
				xysharingImg.Set(width, height, color.RGBA{0, 0, 0, 255})
			} else {
				xysharingImg.Set(width, height, jimg.At(srcX, srcY))

			}
		}
	}

	xysharingImgFile, err := os.Create("./answer_31_3.jpg")
	defer xsharingImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(xysharingImgFile, xysharingImg, &jpeg.Options{100})

}
