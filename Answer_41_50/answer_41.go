package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func color2gray(colorImage image.Image) *image.Gray {
	grayImg := image.NewGray(colorImage.Bounds())
	for y := 0; y < colorImage.Bounds().Size().Y; y++ {
		for x := 0; x < colorImage.Bounds().Size().X; x++ {
			ycbcr := colorImage.At(x, y)
			r, g, b, _ := ycbcr.RGBA()
			var graycolor color.Gray16
			graycolor.Y = uint16(0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b))
			grayImg.Set(x, y, graycolor)
		}
	}
	return grayImg
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

func filterGaussian(grayImage *image.Gray, gaussianMatrix [][]float64) *image.Gray {
	gaussianImg := image.NewGray(grayImage.Bounds())

	W := gaussianImg.Bounds().Size().X
	H := gaussianImg.Bounds().Size().Y
	GW := len(gaussianMatrix)
	GH := len(gaussianMatrix[0][:])

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			filterledVal := 0.0
			for gy := 0; gy < GH; gy++ {
				for gx := 0; gx < GW; gx++ {
					grayX := x + (gx - GW/2)
					if grayX < 0 || grayX > W {
						continue
					}

					grayY := y + (gy - GH/2)
					if grayY < 0 || grayY > W {
						continue
					}

					pixVal := grayImage.GrayAt(grayX, grayY).Y
					filterledVal += float64(pixVal) * gaussianMatrix[gy][gx]
				}
			}
			if filterledVal > 255.0 {
				filterledVal = 255.0
			}

			grayColor := color.Gray{uint8(filterledVal)}
			gaussianImg.Set(x, y, grayColor)
		}
	}
	return gaussianImg
}

func fileterSobel(grayImage *image.Gray, sobelMatrix [3][3]int) *image.Gray {
	sobelImage := image.NewGray(grayImage.Bounds())

	H := grayImage.Bounds().Size().Y
	W := grayImage.Bounds().Size().X

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			// 対象ピクセルを中心とした3x3ピクセルの画素値に対してSobelフィルタを適用する
			sobelValue := 0
			for filterY := 0; filterY < len(sobelMatrix); filterY++ {
				for filterX := 0; filterX < len(sobelMatrix[0][:]); filterX++ {
					// 対象ピクセルの位置を計算する
					srcPointX := filterX + x - 1
					srcPointY := filterY + y - 1
					var pixVal int
					if (srcPointX < 0) || (srcPointX >= W) ||
						(srcPointY < 0) || (srcPointY >= H) {
						// 0パディング
						pixVal = 0
					} else {
						pixVal = int(grayImage.GrayAt(srcPointX, srcPointY).Y)
					}

					// Sobelフィルタ畳み込み
					sobelValue += pixVal * sobelMatrix[filterY][filterX]
				}
			}

			if sobelValue < 0 {
				sobelValue = 0
			}
			if sobelValue > 255 {
				sobelValue = 255
			}

			grayColorSobel := color.Gray{uint8(sobelValue)}
			sobelImage.Set(x, y, grayColorSobel)
		}
	}
	return sobelImage
}

func main() {
	file, err := os.Open("./../assets/imori.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	grayImage := color2gray(jimg)
	grayImageFile, err := os.Create("./answer_41_step1.jpg")
	defer grayImageFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(grayImageFile, grayImage, &jpeg.Options{100})

	gaussianMatrix := createGaussianFilter(5, 5, 1.4)
	sum := 0.0
	for _, row := range gaussianMatrix {
		for _, val := range row {
			sum += val
		}
	}

	W := grayImage.Bounds().Size().X
	H := grayImage.Bounds().Size().Y

	grayGaussianImg := filterGaussian(grayImage, gaussianMatrix)
	grayGaussianImageFile, err := os.Create("./answer_41_step2.jpg")
	defer grayGaussianImageFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(grayGaussianImageFile, grayGaussianImg, &jpeg.Options{100})

	// 縦方向Sobelフィルタを作成
	sobelFilterV := [3][3]int{
		{1, 0, -1},
		{2, 0, -2},
		{1, 0, -1}}
	// for _, row := range sobelFilterV {
	// fmt.Println(row)
	// }

	grayGaussianSobelvImg := fileterSobel(grayGaussianImg, sobelFilterV)
	grayGaussianSobelvImageFile, err := os.Create("./answer_41_step3v.jpg")
	defer grayGaussianSobelvImageFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(grayGaussianSobelvImageFile, grayGaussianSobelvImg, &jpeg.Options{100})

	// 横方向Sobelフィルタを作成
	sobelFilterH := [3][3]int{
		{1, 2, 1},
		{0, 0, 0},
		{-1, -2, -1}}

	grayGaussianSobelhImg := fileterSobel(grayGaussianImg, sobelFilterH)
	grayGaussianSobelhImageFile, err := os.Create("./answer_41_step3h.jpg")
	defer grayGaussianSobelhImageFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(grayGaussianSobelhImageFile, grayGaussianSobelhImg, &jpeg.Options{100})

	edgeImage := image.NewGray(grayGaussianImg.Bounds())
	tanImage := image.NewGray(grayGaussianImg.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			fx := grayGaussianSobelvImg.GrayAt(x, y).Y
			fy := grayGaussianSobelhImg.GrayAt(x, y).Y
			edgeValue := math.Hypot(float64(fx), float64(fy))
			if edgeValue > 255.0 {
				edgeValue = 255.0
			}
			// fmt.Printf("%f ", edgeValue)
			edgeGrayColor := color.Gray{uint8(edgeValue)}
			edgeImage.Set(x, y, edgeGrayColor)

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
			tanGrayColor := color.Gray{uint8(tanValue)}
			tanImage.Set(x, y, tanGrayColor)
		}
		// fmt.Printf("\n")
	}

	edgeImageFile, err := os.Create("./answer_41_1.jpg")
	defer edgeImageFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(edgeImageFile, edgeImage, &jpeg.Options{100})

	tanImageFile, err := os.Create("./answer_41_2.jpg")
	defer tanImageFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(tanImageFile, tanImage, &jpeg.Options{100})

}
