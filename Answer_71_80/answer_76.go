package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func uint32Touint8(pixVal uint32) uint8 {
	return uint8(float64(pixVal*0xFF) / 0xFFFF)
}

func biLinearElm(pixVal00, pixVal01, pixVal10, pixVal11 uint32, dx, dy float64) uint8 {
	pixValUint8_00 := uint32Touint8(pixVal00)
	pixValUint8_01 := uint32Touint8(pixVal01)
	pixValUint8_10 := uint32Touint8(pixVal10)
	pixValUint8_11 := uint32Touint8(pixVal11)
	biLinearVal := (1.0-dx)*(1.0-dy)*float64(pixValUint8_00) +
		dx*(1-dy)*float64(pixValUint8_10) +
		(1.0-dx)*dy*float64(pixValUint8_01) + dx*dy*float64(pixValUint8_11)
	if biLinearVal > 255.0 {
		biLinearVal = 255.0
	}
	if biLinearVal < 0.0 {
		biLinearVal = 0
	}

	return uint8(biLinearVal)
}

func biLinear(I00, I01, I10, I11 color.Color, dx, dy float64) color.Color {
	var biLinearColor color.RGBA
	r00, g00, b00, _ := I00.RGBA()
	r01, g01, b01, _ := I01.RGBA()
	r10, g10, b10, _ := I10.RGBA()
	r11, g11, b11, _ := I11.RGBA()
	biLinearColor.R = biLinearElm(r00, r01, r10, r11, dx, dy)
	biLinearColor.G = biLinearElm(g00, g01, g10, g11, dx, dy)
	biLinearColor.B = biLinearElm(b00, b01, b10, b11, dx, dy)
	biLinearColor.A = 255
	return biLinearColor
}

func bilinearScale(scaleImage *image.Gray, scale float64) *image.Gray {
	// 拡大縮小画像を作成
	scaleBounds := scaleImage.Bounds()
	scaleBounds.Max.X = int(float64(scaleBounds.Max.X) * scale)
	scaleBounds.Max.Y = int(float64(scaleBounds.Max.Y) * scale)
	biLinearImg := image.NewGray(scaleBounds)

	H := biLinearImg.Bounds().Size().Y
	W := biLinearImg.Bounds().Size().X
	for height := 0; height < H; height++ {
		for width := 0; width < W; width++ {
			// 拡大画像のピクセル位置に対応する元画像位置を計算
			srcX := float64(width) / scale
			srcY := float64(height) / scale

			// 上記で計算した位置の周囲4点の画素値を取得
			I00 := scaleImage.At(int(srcX), int(srcY))
			rightIndex := int(srcX) + 1
			if rightIndex > W {
				rightIndex = W
			}
			downIndex := int(srcY) + 1
			if downIndex > H {
				downIndex = H
			}
			I10 := scaleImage.At(rightIndex, int(srcY))
			I01 := scaleImage.At(int(srcX), downIndex)
			I11 := scaleImage.At(rightIndex, downIndex)

			// 周囲4点とピクセル対応点の距離を計算
			dx := srcX - math.Floor(srcX)
			dy := srcY - math.Floor(srcY)

			// Bi Linear補完
			biLinearColor := biLinear(I00, I01, I10, I11, dx, dy)
			biLinearImg.Set(width, height, biLinearColor)
		}
	}
	return biLinearImg
}

func imageDiff(image1, image2 *image.Gray) *image.Gray {
	diffImage := image.NewGray(image1.Bounds())
	H := image1.Bounds().Size().Y
	W := image1.Bounds().Size().X
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			diff := int(image1.GrayAt(x, y).Y) - int(image2.GrayAt(x, y).Y)
			diffImage.SetGray(x, y, color.Gray{uint8(math.Abs(float64(diff)))})
		}
	}
	return diffImage
}

func main() {
	file, err := os.Open("./../Question_71_80/imori.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpegImage, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// カラー画像をグレイスケール画像に変換
	H := jpegImage.Bounds().Size().Y
	W := jpegImage.Bounds().Size().X
	grayImage := image.NewGray(jpegImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			c := color.GrayModel.Convert(jpegImage.At(x, y))
			gray, _ := c.(color.Gray)
			grayImage.Set(x, y, gray)
		}
	}

	grayImageArray := [6]*image.Gray{}
	grayResizeImageArray := [6]*image.Gray{}

	grayImageArray[0] = grayImage

	for i := 1; i <= 5; i++ {
		scaleDenomi := math.Pow(2, float64(i))
		grayScaleImage := bilinearScale(grayImage, 1/scaleDenomi)

		grayImageArray[i] = grayScaleImage
	}

	for i, gray := range grayImageArray {
		scale := math.Pow(2, float64(i))
		grayResizeImage := bilinearScale(gray, scale)
		// grayResizeFile, err := os.Create(fmt.Sprintf("./answer_76_resize_%d.jpg", int(scale)))
		// defer grayResizeFile.Close()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// jpeg.Encode(grayResizeFile, grayResizeImage, &jpeg.Options{100})

		grayResizeImageArray[i] = grayResizeImage
	}

	diffImage0_1 := imageDiff(grayResizeImageArray[0], grayResizeImageArray[1])
	diffImage0_3 := imageDiff(grayResizeImageArray[0], grayResizeImageArray[3])
	diffImage0_5 := imageDiff(grayResizeImageArray[0], grayResizeImageArray[5])
	diffImage1_4 := imageDiff(grayResizeImageArray[1], grayResizeImageArray[4])
	diffImage2_3 := imageDiff(grayResizeImageArray[2], grayResizeImageArray[3])
	diffImage3_5 := imageDiff(grayResizeImageArray[3], grayResizeImageArray[5])

	diffSumArray := make([][]int, H)
	for i := range diffSumArray {
		diffSumArray[i] = make([]int, W)
	}

	diffSumMax := -0xFFFF
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			diffSum := int(diffImage0_1.GrayAt(x, y).Y) +
				int(diffImage0_3.GrayAt(x, y).Y) +
				int(diffImage0_5.GrayAt(x, y).Y) +
				int(diffImage1_4.GrayAt(x, y).Y) +
				int(diffImage2_3.GrayAt(x, y).Y) +
				int(diffImage3_5.GrayAt(x, y).Y)
			if diffSum > diffSumMax {
				diffSumMax = diffSum
			}
			diffSumArray[y][x] = diffSum
		}
	}

	saliencyImage := image.NewGray(jpegImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			gray := float64(diffSumArray[y][x]) * 0xFF / float64(diffSumMax)
			// fmt.Printf("[%d][%d] %f", x, y, gray)
			saliencyImage.Set(x, y, color.Gray{uint8(gray)})
		}
	}

	saliencyFile, err := os.Create("./answer_76.jpg")
	defer saliencyFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(saliencyFile, saliencyImage, &jpeg.Options{100})
}
