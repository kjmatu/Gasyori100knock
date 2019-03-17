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

	for y, rowArray := range gaussianMatrix {
		for x := range rowArray {
			gaussianMatrix[y][x] = gaussian(y-h/2, x-w/2, sigma)
		}
	}

	return gaussianMatrix
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
	// for _, row := range gaussianMatrix {
	// 	fmt.Println(row)
	// }

	W := grayImage.Bounds().Size().X
	H := grayImage.Bounds().Size().Y
	GW := len(gaussianMatrix)
	GH := len(gaussianMatrix[0][:])

	grayGaussianImg := image.NewGray(grayImage.Bounds())

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
			grayColor := color.Gray{uint8(filterledVal)}
			grayGaussianImg.Set(x, y, grayColor)
		}
	}

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

	grayGaussianSobelvImg := image.NewGray(grayGaussianImg.Bounds())

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			filterledVal := 0
			pixCount := 0
			for sobelYv, sobelRowV := range sobelFilterV {
				for sobelXv, sobelVal := range sobelRowV {
					grayX := x + (sobelXv - len(sobelRowV)/2)
					// fmt.Println("grayX", grayX)
					if grayX < 0 || grayX > W {
						continue
					}

					grayY := y + (sobelYv - len(sobelFilterV)/2)
					if grayY < 0 || grayY > W {
						continue
					}

					pixVal := grayGaussianImg.GrayAt(grayX, grayY).Y
					filterledVal += int(pixVal) * sobelVal
					pixCount++
				}
			}
			grayColor := color.Gray{uint8(float64(filterledVal) / float64(pixCount))}
			grayGaussianSobelvImg.Set(x, y, grayColor)
		}
	}
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

	grayGaussianSobelhImg := image.NewGray(grayGaussianImg.Bounds())

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			filterledVal := 0
			pixCount := 0
			for sobelYh, sobelRowH := range sobelFilterH {
				for sobelXh, sobelVal := range sobelRowH {
					grayX := x + (sobelXh - len(sobelRowH)/2)
					if grayX < 0 || grayX >= W {
						continue
					}

					grayY := y + (sobelYh - len(sobelFilterH)/2)
					if grayY < 0 || grayY >= W {
						continue
					}

					pixVal := grayGaussianImg.GrayAt(grayX, grayY).Y
					filterledVal += int(pixVal) * sobelVal
					pixCount++
				}
			}
			grayColor := color.Gray{uint8(float64(filterledVal) / float64(pixCount))}
			grayGaussianSobelhImg.Set(x, y, grayColor)
		}
	}
	grayGaussianSobelhImageFile, err := os.Create("./answer_41_step3h.jpg")
	defer grayGaussianSobelhImageFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(grayGaussianSobelhImageFile, grayGaussianSobelhImg, &jpeg.Options{100})

	edgeImage := image.NewGray(grayGaussianImg.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			fx := grayGaussianSobelvImg.GrayAt(x, y).Y
			fy := grayGaussianSobelhImg.GrayAt(x, y).Y
			edgeValue := math.Hypot(float64(fx), float64(fy))
			fmt.Printf("%f ", edgeValue)
			edgeGrayColor := color.Gray{uint8(edgeValue)}
			edgeImage.Set(x, y, edgeGrayColor)
		}
		fmt.Printf("\n")
	}
	edgeImageFile, err := os.Create("./answer_41_1.jpg")
	defer edgeImageFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(edgeImageFile, edgeImage, &jpeg.Options{100})

}
