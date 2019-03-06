package main

import (
	"image"
	"image/jpeg"
	"log"
	"os"
)

func affineScaleMoveImg(affineMatrix [3][3]int, x, y int) (moveX, moveY int) {
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

	tx := 30
	ty := -30
	affineTranslationMatrix := [3][3]int{{1, 0, tx}, {0, 1, ty}, {0, 0, 1}}

	translationImg := image.NewRGBA(jimg.Bounds())

	for height := 0; height < translationImg.Bounds().Size().Y; height++ {
		for width := 0; width < translationImg.Bounds().Size().X; width++ {
			x, y := affineScaleMoveImg(affineTranslationMatrix, width, height)
			translationImg.Set(x, y, jimg.At(width, height))
		}
	}

	translationImgFile, err := os.Create("./answer_28.jpg")
	defer translationImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(translationImgFile, translationImg, &jpeg.Options{100})

}
