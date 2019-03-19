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

func colorImage2GrayArray(colorImage image.Image) [][]float64 {
	W := colorImage.Bounds().Size().X
	H := colorImage.Bounds().Size().Y

	grayArray := make([][]float64, H)
	for y := range grayArray {
		grayArray[y] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r32, g32, b32, _ := colorImage.At(x, y).RGBA()
			r8 := float64(r32)
			g8 := float64(g32)
			b8 := float64(b32)
			grayValue := 0.2126*r8 + 0.7152*g8 + 0.0722*b8
			grayValue = (grayValue * 0xFF) / 0xFFFF
			grayArray[y][x] = grayValue
		}
	}
	return grayArray
}

func gaussian(x, y int, sigma float64) float64 {
	g := (math.Exp(-1 * float64(x*x+y*y) / float64(2.0*sigma*sigma))) / float64(2.0*math.Pi*sigma*sigma)
	return g
}

func createGaussianFilter(w, h int, sigma float64) [][]float64 {
	gaussianMatrix := make([][]float64, h)
	for y := range gaussianMatrix {
		gaussianMatrix[y] = make([]float64, w)
	}

	sum := 0.0
	for y, rowArray := range gaussianMatrix {
		for x := range rowArray {
			gaussianMatrix[y][x] = gaussian(y-h/2, x-w/2, sigma)
			sum += gaussianMatrix[y][x]
		}
	}

	// ガウシアンフィルタを正規化する
	for y, rowArray := range gaussianMatrix {
		for x := range rowArray {
			gaussianMatrix[y][x] /= sum
		}
	}

	return gaussianMatrix
}

func filterGaussianArray(grayArray, gaussianMatrix [][]float64) [][]float64 {
	W := len(grayArray)
	H := len(grayArray[0][:])
	GW := len(gaussianMatrix)
	GH := len(gaussianMatrix[0][:])

	filteredArray := make([][]float64, H)
	for x := range filteredArray {
		filteredArray[x] = make([]float64, W)
	}
	copy(filteredArray, grayArray)

	for y := GH / 2; y < H-GH/2; y++ {
		for x := GW / 2; x < W-GW/2; x++ {
			// fmt.Println("x", x, "y", y)
			filterledVal := 0.0
			for gy := -GH / 2; gy <= GH/2; gy++ {
				for gx := -GW / 2; gx <= GW/2; gx++ {
					pixVal := grayArray[y+gy][x+gx]
					filterledVal += float64(pixVal) * gaussianMatrix[gy+GH/2][gx+GW/2]
				}
			}

			if filterledVal > 255.0 {
				filterledVal = 255.0
			}
			filteredArray[y][x] = filterledVal
		}
	}
	return filteredArray
}

func fileterSobelArray(grayArray [][]float64, sobelMatrix [3][3]int) [][]float64 {

	H := len(grayArray)
	W := len(grayArray[0])
	sobelArray := make([][]float64, H)
	for x := range sobelArray {
		sobelArray[x] = make([]float64, W)
	}

	for y := 1; y < 131; y++ {
		for x := 1; x < 131; x++ {
			sobelValue := 0.0
			for sy, row := range sobelMatrix {
				for sx, sobel := range row {
					sobelValue += float64(sobel) * grayArray[y+(sy-1)][x+(sx-1)]
				}
			}

			if sobelValue < 0.0 {
				sobelValue = 0.0
			}

			if sobelValue > 255.0 {
				sobelValue = 255.0
			}
			sobelArray[y][x] = sobelValue
		}
	}
	return sobelArray
}

func grayArray2grayImage(grayArray [][]float64) *image.Gray {
	H := len(grayArray)
	W := len(grayArray[0][:])
	grayImage := image.NewGray(image.Rect(0, 0, H, W))
	for y, row := range grayArray {
		for x, val := range row {
			color := color.Gray{uint8(val)}
			grayImage.Set(x, y, color)
		}
	}
	return grayImage
}

func printPixVal(grayImage *image.Gray) {
	W := grayImage.Bounds().Size().X
	H := grayImage.Bounds().Size().Y
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			fmt.Printf("%03d ", grayImage.GrayAt(x, y).Y)
		}
		fmt.Printf("\n")
	}
}

func main() {
	file, err := os.Open("./../Question_41_50/imori.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	grayArray := colorImage2GrayArray(jimg)

	// grayImage := grayArray2grayImage(grayArray)
	// grayImageFile, err := os.Create("./answer_41_step1.jpg")
	// defer grayImageFile.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// jpeg.Encode(grayImageFile, grayImage, &jpeg.Options{100})

	gaussianMatrix := createGaussianFilter(5, 5, 1.4)

	grayW := len(grayArray[0][:])
	grayH := len(grayArray)
	xPad := (len(gaussianMatrix[0][:]) / 2) * 2
	yPad := (len(gaussianMatrix) / 2) * 2
	paddingGrayArray := make([][]float64, grayH+yPad)
	for x := range paddingGrayArray {
		paddingGrayArray[x] = make([]float64, grayW+xPad)
	}

	for y, row := range paddingGrayArray {
		for x := range row {
			if y >= yPad/2 && y < grayH+yPad/2 && x >= xPad/2 && x < grayW+xPad/2 {
				paddingGrayArray[y][x] = grayArray[y-yPad/2][x-xPad/2]
			}
		}
	}

	for y, row := range paddingGrayArray {
		for x := range row {
			if x < xPad/2 {
				paddingGrayArray[y][x] = paddingGrayArray[y][xPad/2]
			}
			if x >= grayW+xPad/2 {
				paddingGrayArray[y][x] = paddingGrayArray[y][grayW+xPad/2-1]
			}
		}
	}

	for y, row := range paddingGrayArray {
		for x := range row {
			if y < yPad/2 {
				paddingGrayArray[y][x] = paddingGrayArray[yPad/2][x]
			}
			if y >= grayH+yPad/2 {
				paddingGrayArray[y][x] = paddingGrayArray[grayH+yPad/2-1][x]
			}
		}
	}

	// grayPaddingImageFile, err := os.Create("./answer_41_padding.jpg")
	// defer grayPaddingImageFile.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// jpeg.Encode(grayPaddingImageFile, grayArray2grayImage(paddingGrayArray), &jpeg.Options{100})

	grayGaussianArray := filterGaussianArray(paddingGrayArray, gaussianMatrix)

	// grayGaussianImg := grayArray2grayImage(grayGaussianArray)
	// grayGaussianImageFile, err := os.Create("./answer_41_step2.jpg")
	// defer grayGaussianImageFile.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// jpeg.Encode(grayGaussianImageFile, grayGaussianImg, &jpeg.Options{100})

	// 縦方向Sobelフィルタを作成
	sobelFilterV := [3][3]int{
		{1, 0, -1},
		{2, 0, -2},
		{1, 0, -1}}

	grayGaussianSobelvArray := fileterSobelArray(grayGaussianArray, sobelFilterV)

	// grayGaussianSobelvImg := grayArray2grayImage(grayGaussianSobelvArray)
	// grayGaussianSobelvImageFile, err := os.Create("./answer_41_step3v.jpg")
	// defer grayGaussianSobelvImageFile.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// jpeg.Encode(grayGaussianSobelvImageFile, grayGaussianSobelvImg, &jpeg.Options{100})

	// 横方向Sobelフィルタを作成
	sobelFilterH := [3][3]int{
		{1, 2, 1},
		{0, 0, 0},
		{-1, -2, -1}}

	grayGaussianSobelhArray := fileterSobelArray(grayGaussianArray, sobelFilterH)

	// grayGaussianSobelhImg := grayArray2grayImage(grayGaussianSobelhArray)
	// grayGaussianSobelhImageFile, err := os.Create("./answer_41_step3h.jpg")
	// defer grayGaussianSobelhImageFile.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// jpeg.Encode(grayGaussianSobelhImageFile, grayGaussianSobelhImg, &jpeg.Options{100})

	edgeArray := make([][]float64, len(grayGaussianSobelhArray))
	tanArray := make([][]float64, len(grayGaussianSobelhArray))
	for x := range edgeArray {
		edgeArray[x] = make([]float64, len(grayGaussianSobelhArray[0][:]))
		tanArray[x] = make([]float64, len(grayGaussianSobelhArray[0][:]))
	}

	for y, row := range edgeArray {
		for x := range row {
			fx := grayGaussianSobelvArray[y][x]
			fy := grayGaussianSobelhArray[y][x]
			edgeValue := math.Hypot(fx, fy)
			if edgeValue > 255.0 {
				edgeValue = 255.0
			}
			edgeArray[y][x] = edgeValue

			tanValue := math.Atan(float64(fy) / float64(fx))
			if tanValue > -0.4142 && tanValue <= 0.4142 {
				tanValue = 0
			} else if tanValue > 0.4142 && tanValue <= 2.4142 {
				tanValue = 45
			} else if tanValue <= -2.4142 || tanValue >= 2.4142 {
				tanValue = 90
			} else if tanValue > -2.4142 || tanValue <= -0.4142 {
				tanValue = 135
			}
			tanArray[y][x] = tanValue
		}
	}

	edgeImage := grayArray2grayImage(edgeArray)
	edgeImageFile, err := os.Create("./answer_41_1.jpg")
	defer edgeImageFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(edgeImageFile, edgeImage, &jpeg.Options{100})

	tanImage := grayArray2grayImage(tanArray)
	tanImageFile, err := os.Create("./answer_41_2.jpg")
	defer tanImageFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(tanImageFile, tanImage, &jpeg.Options{100})

}
