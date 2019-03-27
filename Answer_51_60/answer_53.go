package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func colorImage2GrayImage(colorImage image.Image) *image.Gray {
	grayImage := image.NewGray(colorImage.Bounds())

	W := grayImage.Bounds().Size().X
	H := grayImage.Bounds().Size().Y

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r32, g32, b32, _ := colorImage.At(x, y).RGBA()
			r8 := float64(r32)
			g8 := float64(g32)
			b8 := float64(b32)
			grayValue := 0.2126*r8 + 0.7152*g8 + 0.0722*b8
			grayValue = (grayValue * 0xFF) / 0xFFFF
			grayColor := color.Gray{uint8(grayValue)}
			grayImage.Set(x, y, grayColor)
		}
	}
	return grayImage
}

func discriminantAnalysisMethod(grayImage *image.Gray) *image.Gray {
	// 大津の2値化
	// ヒストグラムを計算
	var hist [255]uint16
	pixVal := 0
	W := grayImage.Bounds().Size().X
	H := grayImage.Bounds().Size().Y

	for height := 0; height < H; height++ {
		for width := 0; width < W; width++ {
			grayIndex := grayImage.GrayAt(width, height).Y
			hist[grayIndex]++
			pixVal += int(grayIndex)
		}
	}

	// ヒストグラムから最適なしきい値ootu_threshを計算
	var thresh int
	pAll := W * H
	sb2Max := 0.0
	ootuThresh := 0
	for thresh = 0; thresh < 255; thresh++ {
		p0 := 0
		m0Sum := 0
		for threshIndex := 0; threshIndex < thresh; threshIndex++ {
			p0 += int(hist[threshIndex])
			m0Sum += threshIndex * int(hist[threshIndex])
		}
		m0 := float64(m0Sum) / float64(p0)
		r0 := float64(p0) / float64(pAll)

		p1 := 0
		m1Sum := 0
		for threshIndex := thresh; threshIndex < 255; threshIndex++ {
			p1 += int(hist[threshIndex])
			m1Sum += threshIndex * int(hist[threshIndex])
		}
		m1 := float64(m1Sum) / float64(p1)
		r1 := float64(p1) / float64(pAll)

		sb2 := r0 * r1 * math.Pow(m0-m1, 2)
		if sb2 > sb2Max {
			sb2Max = sb2
			ootuThresh = thresh
		}
	}

	binImage := image.NewGray(grayImage.Bounds())
	for height := 0; height < H; height++ {
		for width := 0; width < W; width++ {
			gray := grayImage.GrayAt(width, height)
			if gray.Y > uint8(ootuThresh) {
				gray.Y = 255
			} else {
				gray.Y = 0
			}
			binImage.Set(width, height, gray)
		}
	}
	return binImage
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
			targetPix := binImage.GrayAt(x, y).Y
			if targetPix == 0 {
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
			targetPix := binImage.GrayAt(x, y).Y
			if targetPix == 255 {
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
	}
	return erosionImage
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

	// グレイスケール化
	grayImage := colorImage2GrayImage(jimg)

	// 大津の2値化
	binImage := discriminantAnalysisMethod(grayImage)

	// クロージング処理
	closingImage := image.NewGray(binImage.Bounds())
	// copy binImage -> closingImage
	draw.Draw(closingImage, binImage.Bounds(), binImage, binImage.Bounds().Min, draw.Src)
	for i := 0; i < 3; i++ {
		// モルフォロジー処理(膨張)
		closingImage = dialation(closingImage)
	}

	for i := 0; i < 3; i++ {
		// モルフォロジー処理(縮小)
		closingImage = erosion(closingImage)
	}
	// closingFile, err := os.Create("./answer_53_closing.jpg")
	// defer closingFile.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// jpeg.Encode(closingFile, closingImage, &jpeg.Options{100})

	// ブラックハット変換
	// 差分画像
	diffImage := image.NewGray(grayImage.Bounds())

	for y := 0; y < binImage.Bounds().Size().Y; y++ {
		for x := 0; x < binImage.Bounds().Size().X; x++ {
			binVal := binImage.GrayAt(x, y).Y
			closeVal := closingImage.GrayAt(x, y).Y
			diff := int(closeVal) - int(binVal)
			if diff < 0 {
				diff = 0
			}
			// fmt.Println("binVal", binVal, "closeVal", closeVal, "diff", diff)
			diffImage.Set(x, y, color.Gray{uint8(diff)})
		}
	}

	diffFile, err := os.Create("./answer_53.jpg")
	defer diffFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(diffFile, diffImage, &jpeg.Options{100})

}
