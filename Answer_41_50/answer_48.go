package main

import (
	"image"
	"image/color"
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

	// モルフォロジー処理(縮小) 1回目
	erosionImage1 := image.NewGray(binImage.Bounds())
	filter := [3][3]int{
		{0, 1, 0},
		{1, 0, 1},
		{0, 1, 0}}
	for y := 1; y < binImage.Bounds().Size().Y-1; y++ {
		for x := 1; x < binImage.Bounds().Size().X-1; x++ {
			targetPix := binImage.GrayAt(x, y).Y
			// fmt.Println(targetPix)
			if targetPix == 255 {
				sum := 0
				for j, row := range filter {
					for k, value := range row {
						sum += int(binImage.GrayAt(x+k-1, y+j-1).Y) * value
					}
				}
				if sum < 255*4 {
					erosionImage1.Set(x, y, color.Gray{0})
				} else {
					erosionImage1.Set(x, y, binImage.GrayAt(x, y))
				}
			} else {
				erosionImage1.Set(x, y, binImage.GrayAt(x, y))
			}
		}
	}
	// erosionFile1, err := os.Create("./answer_48_1.jpg")
	// defer erosionFile1.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// jpeg.Encode(erosionFile1, erosionImage1, &jpeg.Options{100})

	// モルフォロジー処理(縮小) 2回目
	erosionImage2 := image.NewGray(binImage.Bounds())
	for y := 1; y < erosionImage1.Bounds().Size().Y-1; y++ {
		for x := 1; x < erosionImage1.Bounds().Size().X-1; x++ {
			targetPix := erosionImage1.GrayAt(x, y).Y
			// fmt.Println(targetPix)
			if targetPix == 255 {
				sum := 0
				for j, row := range filter {
					for k, value := range row {
						sum += int(erosionImage1.GrayAt(x+k-1, y+j-1).Y) * value
					}
				}
				if sum < 255*4 {
					erosionImage2.Set(x, y, color.Gray{0})
				} else {
					erosionImage2.Set(x, y, erosionImage1.GrayAt(x, y))
				}
			} else {
				erosionImage2.Set(x, y, erosionImage1.GrayAt(x, y))
			}
		}
	}

	erosionFile2, err := os.Create("./answer_48.jpg")
	defer erosionFile2.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(erosionFile2, erosionImage2, &jpeg.Options{100})

}
