package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

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

	graphicArray := make([][]int, H)
	thinningArray := make([][]int, H)
	for y := range graphicArray {
		graphicArray[y] = make([]int, W)
		thinningArray[y] = make([]int, W)
	}

	// カラー画像をグレイスケール画像に変換
	// 画素値が0以上なら0それ以外は1
	// つまり、0が線、1が背景
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			c := color.GrayModel.Convert(thinningTargetImage.At(x, y))
			gray, _ := c.(color.Gray)
			if gray.Y > 0 {
				graphicArray[y][x] = 0
			} else {
				graphicArray[y][x] = 1
			}
		}
	}

	for {
		// Step1
		changePixelPointsStep1 := []image.Point{}
		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				// 注目画素g0とその8近傍の画素値を取得する
				// 8近傍画素の位置
				// g9 g2 g3
				// g8 g1 g4
				// g7 g6 g5

				// 条件1:黒画素である
				g1 := graphicArray[y][x]
				if g1 != 0 {
					continue
				}

				leftIndex := int(math.Max(float64(x-1), 0))
				rightIndex := int(math.Min(float64(x+1), float64(W)-1))
				upIndex := int(math.Max(float64(y-1), 0))
				downIndex := int(math.Min(float64(y+1), float64(H)-1))
				g9 := graphicArray[upIndex][leftIndex]
				g2 := graphicArray[upIndex][x]
				g3 := graphicArray[upIndex][rightIndex]
				g8 := graphicArray[y][leftIndex]
				g4 := graphicArray[y][rightIndex]
				g7 := graphicArray[downIndex][leftIndex]
				g6 := graphicArray[downIndex][x]
				g5 := graphicArray[downIndex][rightIndex]

				// 条件2:g2, g3, ..., g9, g2と時計まわりに見て、0から1に変わる回数がちょうど1回
				neiborPixels := []int{g2, g3, g4, g5, g6, g7, g8, g9}
				invertCount := 0
				for i, pix := range neiborPixels {
					nextPixIndex := i + 1
					if nextPixIndex >= len(neiborPixels) {
						nextPixIndex = 0
					}
					nextPix := neiborPixels[nextPixIndex]
					if pix == 0 && nextPix == 1 {
						invertCount++
					}
				}
				if invertCount == 1 {
					// OK
				} else {
					continue
				}

				// 条件3:g2, g3, ..., g9の中で1の個数が2以上6以下
				oneCount := 0
				for _, pix := range neiborPixels {
					if pix == 1 {
						oneCount++
					}
				}
				if oneCount >= 2 && oneCount <= 6 {
					// OK
				} else {
					continue
				}

				// 条件4:g2, g4, g6のどれかが1
				if g2 == 1 || g4 == 1 || g6 == 1 {
					// OK
				} else {
					continue
				}

				// 条件5:g4, g6, g8のどれかが1
				if g4 == 1 || g6 == 1 || g8 == 1 {
					// OK
				} else {
					continue
				}

				changePixelPointsStep1 = append(changePixelPointsStep1, image.Point{x, y})
			}
		}

		for _, changePix := range changePixelPointsStep1 {
			// fmt.Println(i, changePix)
			graphicArray[changePix.Y][changePix.X] = 1
		}

		// Step2
		changePixelPointsStep2 := []image.Point{}
		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				// 注目画素g0とその8近傍の画素値を取得する
				// 8近傍画素の位置
				// g9 g2 g3
				// g8 g1 g4
				// g7 g6 g5

				// 条件1:黒画素である
				g1 := graphicArray[y][x]
				if g1 != 0 {
					continue
				}

				leftIndex := int(math.Max(float64(x-1), 0))
				rightIndex := int(math.Min(float64(x+1), float64(W)-1))
				upIndex := int(math.Max(float64(y-1), 0))
				downIndex := int(math.Min(float64(y+1), float64(H)-1))
				g9 := graphicArray[upIndex][leftIndex]
				g2 := graphicArray[upIndex][x]
				g3 := graphicArray[upIndex][rightIndex]
				g8 := graphicArray[y][leftIndex]
				g4 := graphicArray[y][rightIndex]
				g7 := graphicArray[downIndex][leftIndex]
				g6 := graphicArray[downIndex][x]
				g5 := graphicArray[downIndex][rightIndex]

				// 条件2:g2, g3, ..., g9, g2と時計まわりに見て、0から1に変わる回数がちょうど1回
				neiborPixels := []int{g2, g3, g4, g5, g6, g7, g8, g9}
				invertCount := 0
				for i, pix := range neiborPixels {
					nextPixIndex := i + 1
					if nextPixIndex >= len(neiborPixels) {
						nextPixIndex = 0
					}
					nextPix := neiborPixels[nextPixIndex]
					if pix == 0 && nextPix == 1 {
						invertCount++
					}
				}
				if invertCount == 1 {
					// OK
				} else {
					continue
				}

				// 条件3:g2, g3, ..., g9の中で1の個数が2以上6以下
				oneCount := 0
				for _, pix := range neiborPixels {
					if pix == 1 {
						oneCount++
					}
				}
				if oneCount >= 2 && oneCount <= 6 {
					// OK
				} else {
					continue
				}

				// 条件4:g2, g4, g8のどれかが1
				if g2 == 1 || g4 == 1 || g8 == 1 {
					// OK
				} else {
					continue
				}

				// 条件5:g2, g6, g8のどれかが1
				if g2 == 1 || g6 == 1 || g8 == 1 {
					// OK
				} else {
					continue
				}

				changePixelPointsStep2 = append(changePixelPointsStep2, image.Point{x, y})
			}
		}

		for _, changePix := range changePixelPointsStep2 {
			// fmt.Println(i, changePix)
			graphicArray[changePix.Y][changePix.X] = 1
		}

		if len(changePixelPointsStep1) == 0 && len(changePixelPointsStep1) == 0 {
			break
		}
	}

	thinningImage := image.NewRGBA(thinningTargetImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			val := graphicArray[y][x]
			if val == 0 {
				thinningImage.Set(x, y, color.Gray{255})
			} else {
				thinningImage.Set(x, y, color.Gray{0})
			}
		}
	}

	thinningFile, err := os.Create("./answer_65.png")
	defer thinningFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(thinningFile, thinningImage)

}
