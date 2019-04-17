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
	file, err := os.Open("./../assets/imori.jpg")
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
	grayImage := image.NewGray(jpegImage.Bounds())
	grayArray := make([][]float64, H)
	for x := range grayArray {
		grayArray[x] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			// カラーからグレースケールに変換
			r32, g32, b32, _ := jpegImage.At(x, y).RGBA()
			grayfloat64 := 0.2126*float64(r32) + 0.7152*float64(g32) + 0.0722*float64(b32)
			grayfloat64 = (grayfloat64 * 0xFF) / 0xFFFF
			graycolor := color.Gray{uint8(grayfloat64)}
			grayImage.Set(x, y, graycolor)
			grayArray[y][x] = grayfloat64
		}
	}

	// grayFile, err := os.Create("./answer_41_gray.jpg")
	// defer grayFile.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// jpeg.Encode(grayFile, grayImage, &jpeg.Options{100})

	// ガウシアンフィルタを定義式から計算
	gH, gW := 5, 5
	sigma := 1.4
	gaussianMatrix := createGaussianMatrix(gH, gW, sigma)

	grayGaussianArray := applyGaussianFilter(grayArray, gaussianMatrix)

	// grayGaussianImage := image.NewGray(jpegImage.Bounds())
	// grayGaussianFile, err := os.Create("./answer_41_gray_gau.jpg")
	// defer grayGaussianFile.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// jpeg.Encode(grayGaussianFile, grayGaussianImage, &jpeg.Options{100})

	// 縦方向Sobelフィルタを作成
	sobelFilterV := [3][3]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1}}

	SH := len(sobelFilterV)
	SW := len(sobelFilterV[0][:])
	fyImage := image.NewGray(jpegImage.Bounds())
	fyArray := make([][]float64, H)
	for x := range fyArray {
		fyArray[x] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			filterValue := 0.0
			for sy := 0; sy < SW; sy++ {
				for sx := 0; sx < SH; sx++ {
					refY := y + sy - SH/2
					refX := x + sx - SW/2
					gaussianElm := sobelFilterV[sy][sx]
					if refX < 0 && refY < 0 || refX >= W && refY >= H || refX < 0 && refY >= H || refX >= W && refY < 0 {
						filterValue += gaussianElm * grayGaussianArray[y][x]
					} else if refY < 0 || refY >= H {
						filterValue += gaussianElm * grayGaussianArray[y][refX]
					} else if refX < 0 || refX >= W {
						filterValue += gaussianElm * grayGaussianArray[refY][x]
					} else {
						filterValue += gaussianElm * grayGaussianArray[refY][refX]
					}
				}
			}
			fyArray[y][x] = filterValue

			// 画像保存用
			if filterValue > 255.0 {
				filterValue = 255.0
			}
			if filterValue < 0.0 {
				filterValue = 0.0
			}
			fyImage.Set(x, y, color.Gray{uint8(filterValue)})
		}
	}

	// grayGaussianFyFile, err := os.Create("./answer_41_gray_gau_fy.jpg")
	// defer grayGaussianFyFile.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// jpeg.Encode(grayGaussianFyFile, fyImage, &jpeg.Options{100})

	// 横方向Sobelフィルタを作成
	sobelFilterH := [3][3]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1}}

	SH = len(sobelFilterH)
	SW = len(sobelFilterH[0][:])
	fxImage := image.NewGray(jpegImage.Bounds())
	fxArray := make([][]float64, H)
	for x := range fxArray {
		fxArray[x] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			filterValue := 0.0
			for sy := 0; sy < SW; sy++ {
				for sx := 0; sx < SH; sx++ {
					refY := y + sy - SH/2
					refX := x + sx - SW/2
					gaussianElm := sobelFilterH[sy][sx]
					if refX < 0 && refY < 0 || refX >= W && refY >= H || refX < 0 && refY >= H || refX >= W && refY < 0 {
						filterValue += gaussianElm * grayGaussianArray[y][x]
					} else if refY < 0 || refY >= H {
						filterValue += gaussianElm * grayGaussianArray[y][refX]
					} else if refX < 0 || refX >= W {
						filterValue += gaussianElm * grayGaussianArray[refY][x]
					} else {
						filterValue += gaussianElm * grayGaussianArray[refY][refX]
					}
				}
			}
			fxArray[y][x] = filterValue
			// 画像保存用
			if filterValue > 255.0 {
				filterValue = 255.0
			}
			if filterValue < 0.0 {
				filterValue = 0.0
			}
			fxImage.Set(x, y, color.Gray{uint8(filterValue)})
		}
	}

	// grayGaussianFxFile, err := os.Create("./answer_41_gray_gau_fx.jpg")
	// defer grayGaussianFxFile.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// jpeg.Encode(grayGaussianFxFile, fxImage, &jpeg.Options{100})

	edgeImage := image.NewGray(jpegImage.Bounds())
	edgeArray := make([][]float64, H)

	angleImage := image.NewGray(jpegImage.Bounds())
	angleArray := make([][]float64, H)
	for x := range grayGaussianArray {
		edgeArray[x] = make([]float64, W)
		angleArray[x] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			edgeValue := math.Hypot(fxArray[y][x], fyArray[y][x])
			if edgeValue > 255.0 {
				edgeValue = 255.0
			}
			edgeArray[y][x] = edgeValue
			edgeImage.Set(x, y, color.Gray{uint8(edgeValue)})

			fx := fxArray[y][x]
			fy := fyArray[y][x]
			if fx == 0.0 {
				fx = 1e-5
			}
			angle := math.Atan(fy / fx)
			if angle > -0.4142 && angle <= 0.4142 {
				angle = 0
			} else if angle > 0.4142 && angle < 2.4142 {
				angle = 45
			} else if math.Abs(angle) >= 2.4142 {
				angle = 90
			} else if angle > -2.4142 && angle <= -0.4142 {
				angle = 135
			}
			angleArray[y][x] = angle
			angleImage.Set(x, y, color.Gray{uint8(angle)})
		}
	}

	edgeFile, err := os.Create("./answer_41_1.jpg")
	defer edgeFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(edgeFile, edgeImage, &jpeg.Options{100})

	angleFile, err := os.Create("./answer_41_2.jpg")
	defer angleFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(angleFile, angleImage, &jpeg.Options{100})

}
