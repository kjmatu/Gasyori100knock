package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func sobelFileter(grayArray [][]float64, sobelMatrix [3][3]float64) ([][]float64, *image.Gray) {
	SH := len(sobelMatrix)
	SW := len(sobelMatrix[0][:])

	H := len(grayArray)
	W := len(grayArray[0][:])

	sobelImage := image.NewGray(image.Rect(0, 0, W, H))
	sobelArray := make([][]float64, H)
	for x := range sobelArray {
		sobelArray[x] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			filterValue := 0.0
			for sy := 0; sy < SW; sy++ {
				for sx := 0; sx < SH; sx++ {
					refY := y + sy - SH/2
					refX := x + sx - SW/2
					gaussianElm := sobelMatrix[sy][sx]
					if refX < 0 && refY < 0 || refX >= W && refY >= H || refX < 0 && refY >= H || refX >= W && refY < 0 {
						filterValue += gaussianElm * grayArray[y][x]
					} else if refY < 0 || refY >= H {
						filterValue += gaussianElm * grayArray[y][refX]
					} else if refX < 0 || refX >= W {
						filterValue += gaussianElm * grayArray[refY][x]
					} else {
						filterValue += gaussianElm * grayArray[refY][refX]
					}
				}
			}
			sobelArray[y][x] = filterValue

			// 画像保存用
			if filterValue > 255.0 {
				filterValue = 255.0
			}
			if filterValue < 0.0 {
				filterValue = 0.0
			}
			sobelImage.Set(x, y, color.Gray{uint8(filterValue)})
		}
	}
	return sobelArray, sobelImage
}

func determinant(matrix [2][2]float64) float64 {
	return matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]
}

func calcGaussCurvature(hessianMatrix [2][2]float64, Ix, Iy float64) float64 {
	k := determinant(hessianMatrix)
	k /= math.Pow((1 + Ix*Ix + Iy*Iy), 2)
	return k
}

func gaussian(x, y int, sigma float64) float64 {
	g := (math.Exp(-1 * float64(x*x+y*y) / float64(2.0*sigma*sigma))) / float64(2.0*math.Pi*sigma*sigma)
	return g
}

func createGaussianMatrix(width, height int, σ float64) [][]float64 {
	// ガウシアンフィルタを定義式から計算
	sum := 0.0
	gaussianMatrix := make([][]float64, height)
	for i := range gaussianMatrix {
		gaussianMatrix[i] = make([]float64, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 行列の中央を原点とした座標系に変換する
			gx := x - width/2
			gy := y - height/2
			g := gaussian(gx, gy, σ)
			sum += g
			gaussianMatrix[y][x] = g
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gaussianMatrix[y][x] /= sum
		}
	}
	return gaussianMatrix
}

func applyGaussianFilter(grayArray, gaussianMatrix [][]float64) [][]float64 {
	height := len(grayArray)
	width := len(grayArray[0])
	filteredArray := make([][]float64, height)
	for i := range filteredArray {
		filteredArray[i] = make([]float64, width)
	}

	gHeight := len(gaussianMatrix)
	gWidth := len(gaussianMatrix[0])

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			filteredValue := 0.0
			for gY := 0; gY < gHeight; gY++ {
				for gX := 0; gX < gWidth; gX++ {
					refX := gX + x
					refY := gY + y
					if refX < 0 {
						refX = 0
					}
					if refY < 0 {
						refY = 0
					}
					if refX >= width {
						refX = width - 1
					}
					if refY >= height {
						refY = height - 1
					}
					filteredValue += gaussianMatrix[gY][gX] * grayArray[refY][refX]
				}
			}
			filteredArray[y][x] = filteredValue
		}
	}
	return filteredArray
}

func main() {
	file, err := os.Open("./../Question_81_90/thorino.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpegImage, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	H := jpegImage.Bounds().Size().Y
	W := jpegImage.Bounds().Size().X

	// カラー画像をグレイスケール画像に変換
	grayArray := make([][]float64, H)
	for i := range grayArray {
		grayArray[i] = make([]float64, W)
	}

	grayImage := image.NewGray(jpegImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			c := color.GrayModel.Convert(jpegImage.At(x, y))
			gray, _ := c.(color.Gray)
			grayImage.SetGray(x, y, gray)
			grayArray[y][x] = float64(gray.Y)
		}
	}

	Ix2 := make([][]float64, H)
	Iy2 := make([][]float64, H)
	Ixy := make([][]float64, H)
	for i := range grayArray {
		Ix2[i] = make([]float64, W)
		Iy2[i] = make([]float64, W)
		Ixy[i] = make([]float64, W)
	}

	// 縦方向Sobelフィルタを作成
	sobelFilterV := [3][3]float64{
		{1, 2, 1},
		{0, 0, 0},
		{-1, -2, -1}}

	// Y軸方向の微分画像を作成
	Iy, _ := sobelFileter(grayArray, sobelFilterV) // 1次

	// 横方向Sobelフィルタを作成
	sobelFilterH := [3][3]float64{
		{1, 0, -1},
		{2, 0, -2},
		{1, 0, -1}}

	// X軸方向の微分画像を作成
	Ix, _ := sobelFileter(grayArray, sobelFilterH) // 1次

	// 微分画像の2乗を計算
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			Ix2[y][x] = Ix[y][x] * Ix[y][x]
			Iy2[y][x] = Iy[y][x] * Iy[y][x]
			Ixy[y][x] = Ix[y][x] * Iy[y][x]
		}
	}

	// ガウシアンフィルタとの畳込み計算を行う
	gaussianFilter := createGaussianMatrix(3, 3, 3)

	Ix2Gaussian := applyGaussianFilter(Ix2, gaussianFilter)
	Iy2Gaussian := applyGaussianFilter(Iy2, gaussianFilter)
	IxyGaussian := applyGaussianFilter(Ixy, gaussianFilter)
	// 正規化するために最大値を取得しておく
	maxIx2 := math.Inf(-1)
	maxIy2 := math.Inf(-1)
	maxIxy := math.Inf(-1)
	for y, row := range Ix2Gaussian {
		for x := range row {
			valIx2 := Ix2Gaussian[y][x]
			if valIx2 > maxIx2 {
				maxIx2 = valIx2
			}

			valIy2 := Iy2Gaussian[y][x]
			if valIy2 > maxIy2 {
				maxIy2 = valIy2
			}

			valIxy := IxyGaussian[y][x]
			if valIxy > maxIxy {
				maxIxy = valIxy
			}
		}
	}

	Ix2GaussianImage := image.NewGray(grayImage.Bounds())
	Iy2GaussianImage := image.NewGray(grayImage.Bounds())
	IxyGaussianImage := image.NewGray(grayImage.Bounds())
	// ガウシアンフィルタを畳み込んだ配列を正規化して画像に保存する
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			Ix2GaussianImage.SetGray(x, y, color.Gray{uint8(Ix2Gaussian[y][x] * 0xFF / maxIx2)})
			Iy2GaussianImage.SetGray(x, y, color.Gray{uint8(Iy2Gaussian[y][x] * 0xFF / maxIy2)})
			IxyGaussianImage.SetGray(x, y, color.Gray{uint8(IxyGaussian[y][x] * 0xFF / maxIxy)})
		}
	}

	Ix2GaussianFile, err := os.Create("./answer_82_Ix2_gau.jpg")
	defer Ix2GaussianFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(Ix2GaussianFile, Ix2GaussianImage, &jpeg.Options{100})

	Iy2GaussianFile, err := os.Create("./answer_82_Iy2_gau.jpg")
	defer Iy2GaussianFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(Iy2GaussianFile, Iy2GaussianImage, &jpeg.Options{100})

	IxyGaussianFile, err := os.Create("./answer_82_Ixy_gau.jpg")
	defer IxyGaussianFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(IxyGaussianFile, IxyGaussianImage, &jpeg.Options{100})

}
