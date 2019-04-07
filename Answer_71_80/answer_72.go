package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
)

func getMinMaxFromArray(array []float64) (float64, float64) {
	min := math.Inf(1)
	max := math.Inf(-1)
	for index := 0; index < len(array); index++ {
		if array[index] < min {
			min = array[index]
		}

		if array[index] > max {
			max = array[index]
		}
	}
	return min, max
}

// モルフォロジー処理(膨張)
func dialation(binImage *image.Gray) *image.Gray {
	dialationImage := image.NewGray(binImage.Bounds())
	draw.Draw(dialationImage, binImage.Bounds(), binImage, binImage.Bounds().Min, draw.Src)

	filter := [3][3]int{
		{0, 1, 0},
		{1, 0, 1},
		{0, 1, 0}}
	for y := 0; y < binImage.Bounds().Size().Y; y++ {
		for x := 0; x < binImage.Bounds().Size().X; x++ {
			if binImage.GrayAt(x, y).Y == 255 {
				continue
			}

			sum := 0
			for j, row := range filter {
				for k, filterValue := range row {
					refX := x + k - 1
					refY := y + j - 1
					if refX < 0 {
						refX = 0
					}

					if refX >= binImage.Bounds().Size().X {
						refX = binImage.Bounds().Size().X - 1
					}

					if refY < 0 {
						refY = 0
					}

					if refY >= binImage.Bounds().Size().Y {
						refY = binImage.Bounds().Size().Y - 1
					}
					pixVal := binImage.GrayAt(refX, refY).Y
					sum += int(pixVal) * filterValue
				}
			}
			if sum >= 255 {
				dialationImage.Set(x, y, color.Gray{255})
			}
		}
	}
	return dialationImage
}

// モルフォロジー処理(縮小)
func erosion(binImage *image.Gray) *image.Gray {
	erosionImage := image.NewGray(binImage.Bounds())
	draw.Draw(erosionImage, binImage.Bounds(), binImage, binImage.Bounds().Min, draw.Src)

	filter := [3][3]int{
		{0, 1, 0},
		{1, 0, 1},
		{0, 1, 0}}
	for y := 0; y < binImage.Bounds().Size().Y; y++ {
		for x := 0; x < binImage.Bounds().Size().X; x++ {
			if binImage.GrayAt(x, y).Y == 0 {
				continue
			}

			sum := 0
			for j, row := range filter {
				for k, value := range row {
					refX := x + k - 1
					refY := y + j - 1
					if refX < 0 {
						refX = 0
					}

					if refX >= binImage.Bounds().Size().X {
						refX = binImage.Bounds().Size().X - 1
					}

					if refY < 0 {
						refY = 0
					}

					if refY >= binImage.Bounds().Size().Y {
						refY = binImage.Bounds().Size().Y - 1
					}
					sum += int(binImage.GrayAt(refX, refY).Y) * value
				}
			}
			if sum < 255*4 {
				erosionImage.Set(x, y, color.Gray{0})
			}
		}
	}
	return erosionImage
}

func opening(binImage *image.Gray, iteration int) *image.Gray {
	openingImage := image.NewGray(binImage.Bounds())
	draw.Draw(openingImage, binImage.Bounds(), binImage, binImage.Bounds().Min, draw.Src)

	for iter := 0; iter < iteration; iter++ {
		openingImage = erosion(openingImage)
	}

	for iter := 0; iter < iteration; iter++ {
		openingImage = dialation(openingImage)
	}

	return openingImage
}

func closing(binImage *image.Gray, iteration int) *image.Gray {

	closingImage := image.NewGray(binImage.Bounds())
	draw.Draw(closingImage, binImage.Bounds(), binImage, binImage.Bounds().Min, draw.Src)

	for iter := 0; iter < iteration; iter++ {
		closingImage = dialation(closingImage)
	}

	for iter := 0; iter < iteration; iter++ {
		closingImage = erosion(closingImage)
	}

	return closingImage
}

func main() {
	file, err := os.Open("./../Question_61_70/imori.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpegImage, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// マスキング画像作成
	trackPoints := []image.Point{}
	for height := 0; height < jpegImage.Bounds().Size().Y; height++ {
		for width := 0; width < jpegImage.Bounds().Size().X; width++ {
			r32, g32, b32, _ := jpegImage.At(width, height).RGBA()
			// uint32bitからfloat64へ変換
			rfloat64 := float64(r32) / 0xFFFF
			gfloat64 := float64(g32) / 0xFFFF
			bfloat64 := float64(b32) / 0xFFFF

			// RGBからHSVへ変換
			array := []float64{rfloat64, gfloat64, bfloat64}
			min, max := getMinMaxFromArray(array)
			var h float64
			if min == max {
				h = 0.0
			} else if min == bfloat64 {
				h = 60.0*(gfloat64-rfloat64)/(max-min) + 60.0
			} else if min == rfloat64 {
				h = 60.0*(bfloat64-gfloat64)/(max-min) + 180.0
			} else if min == gfloat64 {
				h = 60.0*(rfloat64-bfloat64)/(max-min) + 300.0
			} else {
				h = math.NaN()
			}

			// 色相Hが0.0~360.0の範囲になるように整形
			h = math.Mod(h, 360.0)

			// 青色の範囲を抽出する
			if h >= 180 && h <= 260 {
				trackPoints = append(trackPoints, image.Point{width, height})
			}
		}
	}

	// 青色の位置が1になるマスク画像を作成
	maskImage := image.NewGray(jpegImage.Bounds())
	for _, trackPoint := range trackPoints {
		maskImage.Set(trackPoint.X, trackPoint.Y, color.Gray{uint8(255)})
	}

	// クロージング処理 膨張 -> 収縮
	maskImage = closing(maskImage, 5)

	// オープニング処理 収縮 -> 膨張
	maskImage = opening(maskImage, 5)

	// マスキング処理
	maskedImage := image.NewNRGBA64(jpegImage.Bounds())
	for y := 0; y < maskedImage.Bounds().Size().Y; y++ {
		for x := 0; x < maskedImage.Bounds().Size().X; x++ {
			maskValue := maskImage.GrayAt(x, y).Y
			if maskValue == 255 {
				maskValue = 0
			} else {
				maskValue = 1
			}
			r, g, b, a := jpegImage.At(x, y).RGBA()
			r *= uint32(maskValue)
			g *= uint32(maskValue)
			b *= uint32(maskValue)
			maskColor := color.NRGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			maskedImage.SetNRGBA64(x, y, maskColor)
		}
	}

	maskedFile, err := os.Create("./answer_72.jpeg")
	defer maskedFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(maskedFile, maskedImage, &jpeg.Options{100})

	maskFile, err := os.Create("./answer_72_mask.png")
	defer maskFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(maskFile, maskImage)

}
