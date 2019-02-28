package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"sort"
)

func getMedianPixVal(pixArray []float64) float64 {
	sort.Slice(pixArray, func(i, j int) bool {
		return pixArray[i] < pixArray[j]
	})
	medianValue := 0.0
	if len(pixArray)%2 == 0 {
		leftValue := pixArray[len(pixArray)/2]
		rightValue := pixArray[len(pixArray)/2+1]
		medianValue = (leftValue + rightValue) / 2
	} else {
		medianValue = pixArray[len(pixArray)/2]
	}
	return medianValue
}

func main() {
	file, err := os.Open("./../Question_01_10/imori_noise.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// 平滑化画像を作成
	medianImg := image.NewRGBA(jimg.Bounds())
	filterSize := 3
	for width := 0; width < jimg.Bounds().Size().X; width++ {
		for height := 0; height < jimg.Bounds().Size().Y; height++ {

			// 中央値を求めるための配列
			arrayR := make([]float64, filterSize*filterSize)
			arrayG := make([]float64, filterSize*filterSize)
			arrayB := make([]float64, filterSize*filterSize)
			index := 0
			// 対象ピクセルを中心とした3x3ピクセルの画素値を配列に格納する
			for filterWidth := 0; filterWidth < filterSize; filterWidth++ {
				for filterHeight := 0; filterHeight < filterSize; filterHeight++ {
					// 対象画像の位置を計算する
					srcPointX := filterWidth + width - 1
					srcPointY := filterHeight + height - 1
					var r8, g8, b8 float64
					if (srcPointX < 0) || (srcPointX >= jimg.Bounds().Size().X) ||
						(srcPointY < 0) || (srcPointY >= jimg.Bounds().Size().Y) {
						r8 = 0.0
						g8 = 0.0
						b8 = 0.0
					} else {
						r32, g32, b32, _ := jimg.At(srcPointX, srcPointY).RGBA()
						r8 = (float64(r32) * 0xFF) / 0xFFFF
						g8 = (float64(g32) * 0xFF) / 0xFFFF
						b8 = (float64(b32) * 0xFF) / 0xFFFF
					}
					arrayR[index] = r8
					arrayG[index] = g8
					arrayB[index] = b8
					index++
				}
			}

			medianR := uint8(getMedianPixVal(arrayR))
			medianG := uint8(getMedianPixVal(arrayG))
			medianB := uint8(getMedianPixVal(arrayB))
			var medianColor color.NRGBA
			medianColor.R = medianR
			medianColor.G = medianG
			medianColor.B = medianB
			medianColor.A = uint8(255)
			medianImg.Set(width, height, medianColor)
		}
	}

	medianImgFile, err := os.Create("./answer_10.jpg")
	defer medianImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(medianImgFile, medianImg, &jpeg.Options{100})
}
