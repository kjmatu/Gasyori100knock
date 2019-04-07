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

func calcGaborFilter(x, y int, λ, θ, φ, σ, γ float64) float64 {
	xD := float64(x)*math.Cos(θ) + float64(y)*math.Sin(θ)
	yD := float64(x)*-math.Sin(θ) + float64(y)*math.Cos(θ)

	gabor := math.Exp(-(xD*xD + γ*γ*yD*yD) / (2 * σ * σ))
	gabor *= math.Cos(2*math.Pi*xD/λ + φ)
	return gabor
}

func createGaborFilter(width, height int, λ, θ, φ, σ, γ float64) *image.Gray {
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

	// 最小値を求める
	min := math.Inf(1)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if gaborFilter[y][x] < min {
				min = gaborFilter[y][x]
			}
		}
	}

	// 最小値を0にするように最小値を引く
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gaborFilter[y][x] -= min
		}
	}

	// 最大値を求める
	max := math.Inf(-1)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if gaborFilter[y][x] > max {
				max = gaborFilter[y][x]
			}
		}
	}

	// 0~255までに正規化した画像を保存
	gaborImage := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			grayVal := gaborFilter[y][x] * 255 / max
			gaborImage.SetGray(x, y, color.Gray{uint8(grayVal)})
		}
	}
	return gaborImage
}

func main() {

	K := 111
	λ := 10.0
	φ := 0.0
	σ := 10.0
	γ := 1.2

	for i := 0; i < 4; i++ {
		θ := 45 * i
		r := math.Pi * float64(θ) / 180
		gaborImage := createGaborFilter(K, K, λ, r, φ, σ, γ)
		gaborFile, err := os.Create(fmt.Sprintf("./answer_78_%d.jpg", θ))
		defer gaborFile.Close()
		if err != nil {
			log.Fatal(err)
		}
		jpeg.Encode(gaborFile, gaborImage, &jpeg.Options{100})
	}

}
