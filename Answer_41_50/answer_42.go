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

func color2Gray(jpegImage image.Image) ([][]float64, *image.Gray) {
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
	return grayArray, grayImage
}

func gaussianFilter(grayArray [][]float64, gaussianWidth, gaussianHeight int, sigma float64) ([][]float64, *image.Gray) {
	H := len(grayArray)
	W := len(grayArray[0][:])

	gaussianMatrix := make([][]float64, gaussianHeight)
	for x := range gaussianMatrix {
		gaussianMatrix[x] = make([]float64, gaussianHeight)
	}

	sum := 0.0
	for y, row := range gaussianMatrix {
		for x := range row {
			gy := y - gaussianHeight/2
			gx := x - gaussianHeight/2
			gaussianMatrix[y][x] = gaussian(gx, gy, sigma)
			sum += gaussianMatrix[y][x]
		}
	}

	for y, row := range gaussianMatrix {
		for x := range row {
			gaussianMatrix[y][x] /= sum
		}
	}

	grayGaussianImage := image.NewGray(image.Rect(0, 0, W, H))
	grayGaussianArray := make([][]float64, H)
	for x := range grayGaussianArray {
		grayGaussianArray[x] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < H; x++ {
			filterValue := 0.0
			for gy := 0; gy < gaussianHeight; gy++ {
				for gx := 0; gx < gaussianHeight; gx++ {
					refY := y + gy - gaussianHeight/2
					refX := x + gx - gaussianHeight/2
					gaussianElm := gaussianMatrix[gy][gx]
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
			grayGaussianArray[y][x] = filterValue
			grayGaussianImage.Set(x, y, color.Gray{uint8(filterValue)})
		}
	}
	return grayGaussianArray, grayGaussianImage
}

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

func calcGradientIntensityAndAngle(fxArray, fyArray [][]float64) (edgeArray, angleArray [][]float64, edgeImage, angleImage *image.Gray) {
	H := len(fxArray)
	W := len(fxArray[0][:])

	edgeImage = image.NewGray(image.Rect(0, 0, W, H))
	edgeArray = make([][]float64, H)

	angleImage = image.NewGray(image.Rect(0, 0, W, H))
	angleArray = make([][]float64, H)
	for x := range fxArray {
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
	return edgeArray, angleArray, edgeImage, angleImage
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

	// グレイスケース変換
	grayArray, _ := color2Gray(jpegImage)

	// ガウシアンフィルタを適用
	gH, gW := 5, 5
	sigma := 1.4
	grayGaussianArray, _ := gaussianFilter(grayArray, gH, gW, sigma)

	// 縦方向Sobelフィルタを作成
	sobelFilterV := [3][3]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1}}

	fyArray, _ := sobelFileter(grayGaussianArray, sobelFilterV)

	// 横方向Sobelフィルタを作成
	sobelFilterH := [3][3]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1}}

	fxArray, _ := sobelFileter(grayGaussianArray, sobelFilterH)

	edgeArray, angleArray, _, _ := calcGradientIntensityAndAngle(fxArray, fyArray)

	nonMaximumSuppressionEdgeArray := make([][]float64, len(edgeArray))
	for x := range nonMaximumSuppressionEdgeArray {
		nonMaximumSuppressionEdgeArray[x] = make([]float64, len(edgeArray[0][:]))
	}
	copy(nonMaximumSuppressionEdgeArray, edgeArray)
	nonMaximumSuppressionEdgeImage := image.NewGray(jpegImage.Bounds())
	W := nonMaximumSuppressionEdgeImage.Bounds().Size().X
	H := nonMaximumSuppressionEdgeImage.Bounds().Size().Y

	for y := 1; y < H-1; y++ {
		for x := 1; x < W-1; x++ {
			angle := angleArray[y][x]
			edge := edgeArray[y][x]
			nonMaximumSuppressionEdgeImage.Set(x, y, color.Gray{uint8(edge)})

			if angle == 0 {
				if edge < edgeArray[y][x-1] || edge < edgeArray[y][x+1] {
					nonMaximumSuppressionEdgeArray[y][x] = 0
					nonMaximumSuppressionEdgeImage.Set(x, y, color.Gray{0})
				}
			} else if angle == 45 {
				if edge < edgeArray[y+1][x-1] || edge < edgeArray[y-1][x+1] {
					nonMaximumSuppressionEdgeArray[y][x] = 0
					nonMaximumSuppressionEdgeImage.Set(x, y, color.Gray{0})
				}
			} else if angle == 90 {
				if edge < edgeArray[y-1][x] || edge < edgeArray[y+1][x] {
					nonMaximumSuppressionEdgeArray[y][x] = 0
					nonMaximumSuppressionEdgeImage.Set(x, y, color.Gray{0})
				}
			} else if angle == 135 {
				if edge < edgeArray[y-1][x-1] || edge < edgeArray[y+1][x+1] {
					nonMaximumSuppressionEdgeArray[y][x] = 0
					nonMaximumSuppressionEdgeImage.Set(x, y, color.Gray{0})
				}
			}
		}
	}

	edgeFile, err := os.Create("./answer_42.jpg")
	defer edgeFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(edgeFile, nonMaximumSuppressionEdgeImage, &jpeg.Options{100})

}
