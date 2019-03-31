package main

import (
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
)

func ck(in int) int {
	if in == 1 {
		return 1
	}
	return 0
}

func asta(in int) int {
	return 1 - int(math.Abs(float64(in)))
}

func main() {
	thinningTargetFile, err := os.Open("./../Question_61_70/gazo.png")
	defer thinningTargetFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	thinningTargetImage, err := png.Decode(thinningTargetFile)
	if err != nil {
		log.Fatal(err)
	}

	H := thinningTargetImage.Bounds().Size().Y
	W := thinningTargetImage.Bounds().Size().X

	// カラー画像をグレイスケール画像に変換、画素値が0以上なら1をセットする
	graphicArray := make([][]int, H)
	thinningArray := make([][]int, H)
	for y := range graphicArray {
		graphicArray[y] = make([]int, W)
		thinningArray[y] = make([]int, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {

			gray := graphicArray[y][x]
			if gray > 0 {
				graphicArray[y][x] = 1
			} else {
				graphicArray[y][x] = 0
			}
		}
	}

	for {
		copy(thinningArray, graphicArray)

		counter := 0
		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				// 条件１:図形画素である
				centerVal := graphicArray[y][x]
				if centerVal == 0 {
					continue
				}

				leftIndex := int(math.Max(float64(x-1), 0))
				rightIndex := int(math.Min(float64(x+1), float64(W)-1))
				upIndex := int(math.Max(float64(y-1), 0))
				downIndex := int(math.Min(float64(y+1), float64(H)-1))

				// 条件２:境界点である
				leftPix := graphicArray[y][leftIndex]
				rightPix := graphicArray[y][rightIndex]
				upPix := graphicArray[upIndex][x]
				downPix := graphicArray[downIndex][x]
				boudary := asta(leftPix) + asta(rightPix) + asta(upPix) + asta(downPix)
				if boudary < 1 {
					continue
				}

				// 条件３:端点の保存
				x1 := thinningArray[y][rightIndex]
				x2 := thinningArray[upIndex][rightIndex]
				x3 := thinningArray[upIndex][x]
				x4 := thinningArray[upIndex][leftIndex]
				x5 := thinningArray[y][leftIndex]
				x6 := thinningArray[downIndex][leftIndex]
				x7 := thinningArray[downIndex][x]
				x8 := thinningArray[downIndex][rightIndex]
				endPoints := math.Abs(float64(x1)) + math.Abs(float64(x2)) + math.Abs(float64(x3)) + math.Abs(float64(x4)) + math.Abs(float64(x5)) + math.Abs(float64(x6)) + math.Abs(float64(x7)) + math.Abs(float64(x8))
				if endPoints < 2 {
					continue
				}

				// 条件４:孤立点の保存
				c1 := ck(x1)
				c2 := ck(x2)
				c3 := ck(x3)
				c4 := ck(x4)
				c5 := ck(x5)
				c6 := ck(x6)
				c7 := ck(x7)
				c8 := ck(x8)
				independentPoints := c1 + c2 + c3 + c4 + c5 + c6 + c7 + c8
				if independentPoints < 1 {
					continue
				}

				// 条件５:連結性の保存

			}
		}

		if counter == 0 {
			break
		} else {
			// 細線化した画像を図形画像とする
			draw.Draw(graphicImage, thinningImage.Bounds(), thinningImage, thinningImage.Bounds().Min, draw.Src)
		}
	}

	thinningFile, err := os.Create("./answer_63.png")
	defer thinningFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	for y := 0; y < thinningImage.Bounds().Size().Y; y++ {
		for x := 0; x < thinningImage.Bounds().Size().X; x++ {
			gray := color.Gray{thinningImage.GrayAt(x, y).Y * 255}
			thinningImage.Set(x, y, gray)
		}
	}
	png.Encode(thinningFile, thinningImage)

}
