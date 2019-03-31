package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func main() {
	srcFile, err := os.Open("./../Question_61_70/imori.jpg")
	defer srcFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	srcImage, err := jpeg.Decode(srcFile)
	if err != nil {
		log.Fatal(err)
	}

	H := srcImage.Bounds().Size().Y
	W := srcImage.Bounds().Size().X

	// カラー画像をグレイスケール画像に変換
	grayArray := make([][]float64, H)
	for y := range grayArray {
		grayArray[y] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r, g, b, _ := srcImage.At(x, y).RGBA()
			gray := float64(r)*0.2126 + float64(g)*0.7152 + float64(b)*0.0722
			grayArray[y][x] = (gray * 0xFF) / 0xFFFF
		}
	}

	// x,y方向の輝度勾配を求める
	gx := make([][]float64, H)
	gy := make([][]float64, H)
	for y := range gx {
		gx[y] = make([]float64, W)
		gy[y] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			leftIndex := int(math.Max(float64(x-1), 0))
			rightIndex := int(math.Min(float64(x+1), float64(W-1)))
			gx[y][x] = grayArray[y][rightIndex] - grayArray[y][leftIndex]

			upIndex := int(math.Max(float64(y-1), 0))
			downIndex := int(math.Min(float64(y+1), float64(H-1)))
			gy[y][x] = grayArray[downIndex][x] - grayArray[upIndex][x]
		}
	}

	// x,y方向の輝度勾配から勾配強度と角度を求める
	mag := make([][]float64, H)
	angle := make([][]float64, H)
	for y := range mag {
		mag[y] = make([]float64, W)
		angle[y] = make([]float64, W)
	}

	magMax := -1.0
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			mag[y][x] = math.Hypot(gy[y][x], gx[y][x])
			if magMax < mag[y][x] {
				magMax = mag[y][x]
			}

			angle[y][x] = (180 * math.Atan(gy[y][x]/gx[y][x])) / math.Pi
			if angle[y][x] < 0 {
				angle[y][x] += 180
			} else if math.IsNaN(angle[y][x]) {
				angle[y][x] = 0
			}

		}
	}

	// 勾配角度を 0~180度で9分割した値に量子化する
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			ang := angle[y][x]
			if ang >= 0 && ang < 20 {
				angle[y][x] = 0
			} else if ang >= 20 && ang < 40 {
				angle[y][x] = 1
			} else if ang >= 40 && ang < 60 {
				angle[y][x] = 2
			} else if ang >= 60 && ang < 80 {
				angle[y][x] = 3
			} else if ang >= 80 && ang < 100 {
				angle[y][x] = 4
			} else if ang >= 100 && ang < 120 {
				angle[y][x] = 5
			} else if ang >= 120 && ang < 140 {
				angle[y][x] = 6
			} else if ang >= 140 && ang < 160 {
				angle[y][x] = 7
			} else if ang >= 160 && ang <= 180 {
				angle[y][x] = 8
			}
		}
	}

	colorList := []color.RGBA{color.RGBA{0, 0, 255, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 255, 255},
		color.RGBA{255, 0, 255, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{0, 127, 127, 255},
		color.RGBA{127, 0, 127, 255},
		color.RGBA{127, 127, 0, 255}}

	gradColorImage := image.NewRGBA(srcImage.Bounds())
	magnitudeImage := image.NewGray(srcImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			gradient := angle[y][x]
			magnitude := mag[y][x]
			gradColorImage.Set(x, y, colorList[int(gradient)])
			magnitudeImage.Set(x, y, color.Gray{uint8(magnitude * 255 / magMax)})
		}
	}

	magnitudeFile, err := os.Create("./answer_66_mag.jpg")
	defer magnitudeFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(magnitudeFile, magnitudeImage, &jpeg.Options{100})

	gradientFile, err := os.Create("./answer_66_gra.jpg")
	defer gradientFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(gradientFile, gradColorImage, &jpeg.Options{100})

}
