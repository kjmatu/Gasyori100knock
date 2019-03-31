package main

import (
	"image"
	"image/color"
	"image/draw"
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

	// カラー画像をグレイスケール画像に変換、画素値が0以上なら1をセットする
	graphicImage := image.NewGray(thinningTargetImage.Bounds())
	for y := 0; y < thinningTargetImage.Bounds().Size().Y; y++ {
		for x := 0; x < thinningTargetImage.Bounds().Size().X; x++ {
			c := color.GrayModel.Convert(thinningTargetImage.At(x, y))
			gray, _ := c.(color.Gray)
			if gray.Y > 0 {
				graphicImage.Set(x, y, color.Gray{1})
			} else {
				graphicImage.Set(x, y, color.Gray{0})
			}
		}
	}

	thinningImage := image.NewGray(graphicImage.Bounds())
	for {
		draw.Draw(thinningImage, graphicImage.Bounds(), graphicImage, graphicImage.Bounds().Min, draw.Src)

		counter := 0
		for y := 0; y < graphicImage.Bounds().Size().Y; y++ {
			for x := 0; x < graphicImage.Bounds().Size().X; x++ {
				centerVal := graphicImage.GrayAt(x, y).Y
				if centerVal == 0 {
					continue
				}

				leftIndex := int(math.Max(float64(x-1), 0))
				rightIndex := int(math.Min(float64(x+1), float64(graphicImage.Bounds().Size().X)-1))
				upIndex := int(math.Max(float64(y-1), 0))
				downIndex := int(math.Min(float64(y+1), float64(graphicImage.Bounds().Size().Y)-1))

				// 境界条件 注目画素の4近傍に0が１つ以上存在
				// 境界条件は処理前の配列で判定
				boudaryFlag := false
				leftPix := graphicImage.GrayAt(leftIndex, y).Y
				rightPix := graphicImage.GrayAt(rightIndex, y).Y
				upPix := graphicImage.GrayAt(x, upIndex).Y
				downPix := graphicImage.GrayAt(x, downIndex).Y
				if leftPix*rightPix*upPix*downPix == 0 {
					boudaryFlag = true
				}

				// 連結性条件及び非端点条件は処理後の配列を用いて判定
				// 連結性条件 注目画素の4連結数が1
				connectivityFlag := false
				x1 := int(thinningImage.GrayAt(rightIndex, y).Y)
				x2 := int(thinningImage.GrayAt(rightIndex, upIndex).Y)
				x3 := int(thinningImage.GrayAt(x, upIndex).Y)
				x4 := int(thinningImage.GrayAt(leftIndex, upIndex).Y)
				x5 := int(thinningImage.GrayAt(leftIndex, y).Y)
				x6 := int(thinningImage.GrayAt(leftIndex, downIndex).Y)
				x7 := int(thinningImage.GrayAt(x, downIndex).Y)
				x8 := int(thinningImage.GrayAt(rightIndex, downIndex).Y)
				s := (x1 - x1*x2*x3) + (x3 - x3*x4*x5) + (x5 - x5*x6*x7) + (x7 - x7*x8*x1)
				if s == 1 {
					connectivityFlag = true
				}

				// 非端点条件 注目画素の8近傍に1が３つ以上存在
				nonEndPointsFlag := false
				if (x1 + x2 + x3 + x4 + x5 + x6 + x7 + x8) >= 3 {
					nonEndPointsFlag = true
				}

				if boudaryFlag && connectivityFlag && nonEndPointsFlag {
					thinningImage.Set(x, y, color.Gray{0})
					counter++
				}
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
