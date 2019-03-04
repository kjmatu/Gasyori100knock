package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

// Laplacian of Gaussian filter
func calcLoG(x, y int, s float64) float64 {
	fmt.Println("x", x)
	fmt.Println("y", y)
	fmt.Println("s", s)
	fmt.Println(-(float64(x*x+y*y) / (2 * s * s)))
	fmt.Println(math.Exp(-(float64(x*x+y*y) / (2 * s * s))))
	numerator := (float64(x*x) + float64(y*y) - (s * s)) * math.Exp(-(float64(x*x+y*y) / (2 * s * s)))
	fmt.Println("numerator", numerator)
	denominator := 2 * math.Pi * math.Pow(s, 6)
	fmt.Println("denominator", denominator)
	return numerator / denominator
}

func main() {
	file, err := os.Open("./../Question_11_20/imori_noise.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// グレースケール画像に変換
	grayimg := image.NewGray(jimg.Bounds())
	for height := 0; height < jimg.Bounds().Size().Y; height++ {
		for width := 0; width < jimg.Bounds().Size().X; width++ {
			ycbcr := jimg.At(width, height)
			r32, g32, b32, _ := ycbcr.RGBA()
			// 32bit画像から8bit画像に変換
			r8 := uint8((float64(r32) * 0xFF) / 0xFFFF)
			g8 := uint8((float64(g32) * 0xFF) / 0xFFFF)
			b8 := uint8((float64(b32) * 0xFF) / 0xFFFF)

			var graycolor color.Gray
			// カラーからグレースケールに変換
			grayfloat64 := 0.2126*float64(r8) + 0.7152*float64(g8) + 0.0722*float64(b8)
			graycolor.Y = uint8(grayfloat64)
			grayimg.Set(width, height, graycolor)
		}
	}

	filterSize := 3
	// LoGフィルタを作成 sigma=3
	logFilter := [3][3]float64{
		{2.0 / 16.0, 4.0 / 16.0, 2.0 / 16.0},
		{-6.0 / 16.0, -12.0 / 16.0, -6.0 / 16.0},
		{2.0 / 16.0, 4.0 / 16.0, 2.0 / 16.0}
	}
	// sigma := 3.0
	// sum := 0.0
	// for i := 0; i < filterSize; i++ {
	// 	for j := 0; j < filterSize; j++ {
	// 		logFilter[j][i] = calcLoG(i-1, j-1, sigma)
	// 		sum += logFilter[j][i]
	// 	}
	// }
	// fmt.Println("logFilter", logFilter)
	// fmt.Println("sum", sum)
	// for i := 0; i < filterSize; i++ {
	// 	for j := 0; j < filterSize; j++ {
	// 		logFilter[j][i] /= sum
	// 	}
	// }

	fmt.Println("logFilter", logFilter)

	// LoGフィルタ適用済画像保存先を作成
	logImg := image.NewGray(grayimg.Bounds())

	for y := 0; y < grayimg.Bounds().Size().Y; y++ {
		for x := 0; x < grayimg.Bounds().Size().X; x++ {
			// 対象ピクセルを中心とした3x3ピクセルの画素値の最大値と最小値を取得する
			var filteredValue float64
			for filterY := 0; filterY < filterSize; filterY++ {
				for filterX := 0; filterX < filterSize; filterX++ {
					// 対象ピクセルの位置を計算する
					srcPointX := filterX + x - 1
					srcPointY := filterY + y - 1
					var pixVal int
					if (srcPointX < 0) || (srcPointX >= grayimg.Bounds().Size().X) ||
						(srcPointY < 0) || (srcPointY >= grayimg.Bounds().Size().Y) {
						// 0パディング
						pixVal = 0
					} else {
						pixVal = int(grayimg.GrayAt(srcPointX, srcPointY).Y)
					}

					// LoGフィルタ畳み込み
					filteredValue += float64(pixVal) * logFilter[filterY][filterX]
				}
			}

			if filteredValue < 0 {
				filteredValue = 0
			}

			if filteredValue > 255 {
				filteredValue = 255
			}

			var filteredGray color.Gray
			filteredGray.Y = uint8(filteredValue)
			logImg.Set(x, y, filteredGray)
		}
	}

	logFilterImgFile, err := os.Create("./answer_19.jpg")
	defer logFilterImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(logFilterImgFile, logImg, &jpeg.Options{100})

}
