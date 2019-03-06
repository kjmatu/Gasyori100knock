package main

import (
	"image"
	"image/jpeg"
	"log"
	"os"
)

func affineMove(affineMatrix [3][3]int, x, y int) (moveX, moveY int) {
	moveArray := [3]int{}
	for i, rowArray := range affineMatrix {
		// fmt.Printf("[%d]%d\n", i, rowArray)
		moveArray[i] = rowArray[0]*x + rowArray[1]*y + rowArray[2]*1
	}
	// fmt.Println(moveArray)
	moveX = moveArray[0]
	moveY = moveArray[1]
	return moveX, moveY
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

	sx := 1.3
	sy := 0.8
	scaleRectangle := jimg.Bounds()
	scaleRectangle.Max.X = int(float64(scaleRectangle.Max.X) * sx)
	scaleRectangle.Max.Y = int(float64(scaleRectangle.Max.Y) * sy)

	affineScaleImg := image.NewRGBA(scaleRectangle)
	for height := 0; height < affineScaleImg.Bounds().Size().Y; height++ {
		for width := 0; width < affineScaleImg.Bounds().Size().X; width++ {
			// 拡大縮小画像に対応する元画像の位置を計算する
			srcX := int(float64(width) / sx)
			srcY := int(float64(height) / sy)
			affineScaleImg.Set(width, height, jimg.At(srcX, srcY))
		}
	}

	affineScaleImgFile, err := os.Create("./answer_29_1.jpg")
	defer affineScaleImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(affineScaleImgFile, affineScaleImg, &jpeg.Options{100})

	tx := 30
	ty := -30
	affineTranslationMatrix := [3][3]int{{1, 0, tx}, {0, 1, ty}, {0, 0, 1}}

	affineScaleMoveImg := image.NewRGBA(affineScaleImg.Bounds())
	for height := 0; height < affineScaleMoveImg.Bounds().Size().Y; height++ {
		for width := 0; width < affineScaleMoveImg.Bounds().Size().X; width++ {
			x, y := affineMove(affineTranslationMatrix, width, height)
			affineScaleMoveImg.Set(x, y, affineScaleImg.At(width, height))
		}
	}

	affineScaleMoveImgFile, err := os.Create("./answer_29_2.jpg")
	defer affineScaleMoveImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(affineScaleMoveImgFile, affineScaleMoveImg, &jpeg.Options{100})

}
