package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func calcGaborFilter(x, y int, λ, θ, φ, σ, γ float64) float64 {
	xD := float64(x)*math.Cos(θ) + float64(y)*math.Sin(θ)
	yD := float64(x)*-math.Sin(θ) + float64(y)*math.Cos(θ)

	gabor := math.Exp(-(xD*xD + γ*γ*yD*yD) / (2 * σ * σ))
	gabor *= math.Cos(2*math.Pi*xD/λ + φ)
	return gabor
}

func createGaborFilter(width, height int, λ, θ, φ, σ, γ float64) [][]float64 {
	gaborFilter := make([][]float64, height)
	for i := range gaborFilter {
		gaborFilter[i] = make([]float64, width)
	}

	sum := 0.0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			refX := x - width/2
			refY := y - height/2
			gaborFilter[y][x] = calcGaborFilter(refX, refY, λ, θ, φ, σ, γ)
			sum += math.Abs(gaborFilter[y][x])
		}
	}

	// フィルタ値の絶対値の和が1になるように正規化
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gaborFilter[y][x] /= sum
		}
	}
	return gaborFilter
}

func main() {
	imoriFile, err := os.Open("./../Question_71_80/imori.jpg")
	defer imoriFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	imoriImage, err := jpeg.Decode(imoriFile)
	if err != nil {
		log.Fatal(err)
	}

	H := imoriImage.Bounds().Size().Y
	W := imoriImage.Bounds().Size().X

	// カラー画像をグレイスケール画像に変換
	grayImage := image.NewGray(imoriImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			c := color.GrayModel.Convert(imoriImage.At(x, y))
			gray, _ := c.(color.Gray)
			grayImage.SetGray(x, y, gray)
		}
	}

	K := 11
	λ := 3.0
	φ := 0.0
	σ := 1.5
	γ := 1.2

	featureArray := make([][]float64, H)
	for i := range featureArray {
		featureArray[i] = make([]float64, W)
	}

	for i := 0; i < 4; i++ {
		θ := 45 * i
		r := math.Pi * float64(θ) / 180
		gaborFilter := createGaborFilter(K, K, λ, r, φ, σ, γ)

		for y := 0; y < H; y++ {
			for x := 0; x < H; x++ {
				pixVal := 0.0
				for k := 0; k < K; k++ {
					for j := 0; j < K; j++ {
						gaborVal := gaborFilter[k][j]
						refX := x + j - K/2
						if refX < 0 {
							refX = 0
						} else if refX > W-1 {
							refX = W - 1
						}

						refY := y + k - K/2
						if refY < 0 {
							refY = 0
						} else if refY > H-1 {
							refY = H - 1
						}
						pixVal += gaborVal * float64(grayImage.GrayAt(refX, refY).Y)
					}
				}
				if pixVal > 255 {
					pixVal = 255
				} else if pixVal < 0 {
					pixVal = 0
				}
				featureArray[y][x] += pixVal
			}
		}
	}

	gaborFeatureImage := image.NewGray(imoriImage.Bounds())

	max := math.Inf(-1)
	for y := 0; y < H; y++ {
		for x := 0; x < H; x++ {
			if featureArray[y][x] > max {
				max = featureArray[y][x]
			}
		}
	}

	for y := 0; y < H; y++ {
		for x := 0; x < H; x++ {
			pixVal := float64(featureArray[y][x]) * 255 / max
			gaborFeatureImage.SetGray(x, y, color.Gray{uint8(pixVal)})
		}
	}

	gaborFeatureFile, err := os.Create("./answer_80.jpg")
	defer gaborFeatureFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(gaborFeatureFile, gaborFeatureImage, &jpeg.Options{100})

}
