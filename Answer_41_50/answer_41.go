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

	grayFile, err := os.Create("./answer_41_gray.jpg")
	defer grayFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(grayFile, grayImage, &jpeg.Options{100})

	// ガウシアンフィルタを定義式から計算
	gH, gW := 5, 5
	sigma := 1.4
	gaussianMatrix := make([][]float64, gH)
	for x := range gaussianMatrix {
		gaussianMatrix[x] = make([]float64, gW)
	}

	sum := 0.0
	for y, row := range gaussianMatrix {
		for x := range row {
			gy := y - gH/2
			gx := x - gW/2
			gaussianMatrix[y][x] = gaussian(gx, gy, sigma)
			sum += gaussianMatrix[y][x]
		}
	}

	for y, row := range gaussianMatrix {
		for x := range row {
			gaussianMatrix[y][x] /= sum
		}
	}

	grayGaussianImage := image.NewGray(jpegImage.Bounds())
	grayGaussianArray := make([][]float64, H)
	for x := range grayGaussianArray {
		grayGaussianArray[x] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < H; x++ {
			filterValue := 0.0
			for gy := 0; gy < gH; gy++ {
				for gx := 0; gx < gW; gx++ {
					refY := y + gy - gH/2
					refX := x + gx - gW/2
					// fmt.Printf("[%d][%d] ", refX, refY)
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
					// fmt.Println(gaussianElm, filterValue)
				}
				// fmt.Println()
			}
			// fmt.Println()
			grayGaussianArray[y][x] = filterValue
			grayGaussianImage.Set(x, y, color.Gray{uint8(filterValue)})
		}
	}

	// for _, row := range grayGaussianArray {
	// 	fmt.Println(row)
	// }

	grayGaussianFile, err := os.Create("./answer_41_gray_gau.jpg")
	defer grayGaussianFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(grayGaussianFile, grayGaussianImage, &jpeg.Options{100})

	// 縦方向Sobelフィルタを作成
	sobelFilterV := [3][3]float64{
		{1, 0, -1},
		{2, 0, -2},
		{1, 0, -1}}

	SH := len(sobelFilterV)
	SW := len(sobelFilterV[0][:])
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
					// fmt.Printf("[%d][%d] ", refX, refY)
					gaussianElm := sobelFilterV[sy][sx]
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
			if filterValue > 255.0 {
				filterValue = 255.0
			}
			if filterValue < 0.0 {
				filterValue = 0.0
			}
			fxArray[y][x] = filterValue
			fxImage.Set(x, y, color.Gray{uint8(filterValue)})
		}
	}
	grayGaussianFxFile, err := os.Create("./answer_41_gray_gau_fx.jpg")
	defer grayGaussianFxFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(grayGaussianFxFile, fxImage, &jpeg.Options{100})

	// 横方向Sobelフィルタを作成
	sobelFilterH := [3][3]float64{
		{1, 2, 1},
		{0, 0, 0},
		{-1, -2, -1}}

	SH = len(sobelFilterH)
	SW = len(sobelFilterH[0][:])
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
					// fmt.Printf("[%d][%d] ", refX, refY)
					gaussianElm := sobelFilterH[sy][sx]
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
			if filterValue > 255.0 {
				filterValue = 255.0
			}
			if filterValue < 0.0 {
				filterValue = 0.0
			}
			fxArray[y][x] = filterValue
			fyImage.Set(x, y, color.Gray{uint8(filterValue)})
		}
	}

	grayGaussianFyFile, err := os.Create("./answer_41_gray_gau_fy.jpg")
	defer grayGaussianFyFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(grayGaussianFyFile, fyImage, &jpeg.Options{100})

	edgeImage := image.NewGray(jpegImage.Bounds())
	angleImage := image.NewGray(jpegImage.Bounds())
	edgeArray := make([][]float64, H)
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

			angle := fxArray[y][x] / fyArray[y][x]
			if angle > -0.4142 && angle <= 0.4142 {
				angle = 0
			} else if angle > 0.4142 && angle < 2.4142 {
				angle = 45
			} else if angle <= -2.4142 || angle >= 2.4142 {
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
