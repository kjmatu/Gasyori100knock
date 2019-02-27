package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func gaussian(x, y int, sigma float64) float64 {
	g := (math.Exp(-1 * float64(x*x+y*y) / float64(2.0*sigma*sigma))) / float64(2.0*math.Pi*sigma*sigma)
	return g
}

func main() {
	// 3x3ピクセル、標準偏差1.3の場合によく使われる近似行列
	gaussianMatrix := [3][3]float64{
		{1.0 / 16.0, 2.0 / 16.0, 1.0 / 16.0},
		{2.0 / 16.0, 4.0 / 16.0, 2.0 / 16.0},
		{1.0 / 16.0, 2.0 / 16.0, 1.0 / 16.0}}
	file, err := os.Open("./../Question_01_10/imori_noise.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// ガウシアンフィルタを定義式から計算
	sum := 0.0
	calcGaussian := [3][3]float64{}
	for width := -1; width < 2; width++ {
		for height := -1; height < 2; height++ {
			g := gaussian(width, height, 1.3)
			sum += g
			calcGaussian[width+1][height+1] = g
		}
	}

	for width := 0; width < 3; width++ {
		for height := 0; height < 3; height++ {
			calcGaussian[width][height] /= sum
		}
	}

	// ノイズ除去画像を作成
	denoiseImg := image.NewRGBA(jimg.Bounds())
	filterSize := 3
	for width := 0; width < jimg.Bounds().Size().X; width++ {
		for height := 0; height < jimg.Bounds().Size().Y; height++ {

			// ガウシアンフィルタの畳み込み処理
			var filR, filG, filB float64
			// 対象ピクセルを中心とした3x3ピクセルの画素値に対してガウシアンフィルタを適用する
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
						// 3x3ピクセル、標準偏差1.3の場合によく使われる近似行列を使った場合
						r8 *= gaussianMatrix[filterWidth][filterHeight]
						g8 *= gaussianMatrix[filterWidth][filterHeight]
						b8 *= gaussianMatrix[filterWidth][filterHeight]
						// 定義式から計算したガウシアンフィルタを使った場合
						// r8 *= calcGaussian[filterWidth][filterHeight]
						// g8 *= calcGaussian[filterWidth][filterHeight]
						// b8 *= calcGaussian[filterWidth][filterHeight]
					}
					filR += r8
					filG += g8
					filB += b8
				}
			}

			var denoiseColor color.NRGBA
			denoiseColor.R = uint8(filR)
			denoiseColor.G = uint8(filG)
			denoiseColor.B = uint8(filB)
			denoiseColor.A = uint8(255)
			denoiseImg.Set(width, height, denoiseColor)
		}
	}

	denoiseImgFile, err := os.Create("./answer_9.jpg")
	defer denoiseImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(denoiseImgFile, denoiseImg, &jpeg.Options{100})
}
