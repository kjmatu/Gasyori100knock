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

func main() {

	K := 111
	gaborFilter := make([][]float64, K)
	for i := range gaborFilter {
		gaborFilter[i] = make([]float64, K)
	}

	λ := 10.0
	θ := 0.0
	φ := 0.0
	σ := 10.0
	γ := 1.2

	sum := 0.0
	for y := 0; y < K; y++ {
		for x := 0; x < K; x++ {
			refX := x - K/2
			refY := y - K/2
			gaborFilter[y][x] = calcGaborFilter(refX, refY, λ, θ, φ, σ, γ)
			sum += math.Abs(gaborFilter[y][x])
		}
	}

	// フィルタ値の絶対値の和が1になるように正規化
	for y := 0; y < K; y++ {
		for x := 0; x < K; x++ {
			gaborFilter[y][x] /= sum
		}
	}

	// 最小値を求める
	min := math.Inf(1)
	for y := 0; y < K; y++ {
		for x := 0; x < K; x++ {
			if gaborFilter[y][x] < min {
				min = gaborFilter[y][x]
			}
		}
	}

	// 最小値を0にするように最小値を引く
	for y := 0; y < K; y++ {
		for x := 0; x < K; x++ {
			gaborFilter[y][x] -= min
		}
	}

	// 最大値を求める
	max := math.Inf(-1)
	for y := 0; y < K; y++ {
		for x := 0; x < K; x++ {
			if gaborFilter[y][x] > max {
				max = gaborFilter[y][x]
			}
		}
	}

	// 0~255までに正規化した画像を保存
	gaborImage := image.NewGray(image.Rect(0, 0, K, K))
	for y := 0; y < K; y++ {
		for x := 0; x < K; x++ {
			grayVal := gaborFilter[y][x] * 255 / max
			gaborImage.SetGray(x, y, color.Gray{uint8(grayVal)})
		}
	}
	gaborFile, err := os.Create("./answer_77.jpg")
	defer gaborFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(gaborFile, gaborImage, &jpeg.Options{100})

}
