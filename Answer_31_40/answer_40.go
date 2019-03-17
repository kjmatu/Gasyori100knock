package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func dct(grayImg [][]float64, blockSize int) [][]float64 {
	H := len(grayImg)
	W := len(grayImg[0][:])

	dctResult := make([][]float64, H)
	for x := range dctResult {
		dctResult[x] = make([]float64, W)
	}

	// 離散コサイン変換
	for vi := 0; vi < H; vi += blockSize {
		for ui := 0; ui < W; ui += blockSize {
			for v := 0; v < blockSize; v++ {
				for u := 0; u < blockSize; u++ {
					convVal := 0.0
					for y := 0; y < blockSize; y++ {
						for x := 0; x < blockSize; x++ {
							pixVal := grayImg[vi+y][ui+x]
							cosVal1 := math.Cos((float64((2*x+1)*u) * math.Pi) / 16)
							cosVal2 := math.Cos((float64((2*y+1)*v) * math.Pi) / 16)

							cu := 1.0
							if u == 0 {
								cu = 1 / math.Sqrt2
							}

							cv := 1.0
							if v == 0 {
								cv = 1 / math.Sqrt2
							}

							convVal += cu * cv * float64(pixVal) * cosVal1 * cosVal2
						}
					}

					dctVal := convVal / 4
					dctResult[vi+v][ui+u] = dctVal
				}
			}
		}
	}
	return dctResult
}

func idct(array2D [][]float64, blockSize int) [][]float64 {
	H := len(array2D)
	W := len(array2D[0][:])

	idctResult := make([][]float64, H)
	for x := range idctResult {
		idctResult[x] = make([]float64, W)
	}

	for yi := 0; yi < H; yi += blockSize {
		for xi := 0; xi < H; xi += blockSize {
			for y := 0; y < blockSize; y++ {
				for x := 0; x < blockSize; x++ {
					iconvVal := 0.0
					for v := 0; v < blockSize; v++ {
						for u := 0; u < blockSize; u++ {
							pixVal := array2D[yi+v][xi+u]
							cosVal1 := math.Cos((float64((2*x+1)*u) * math.Pi) / 16)
							cosVal2 := math.Cos((float64((2*y+1)*v) * math.Pi) / 16)
							cu := 1.0
							if u == 0 {
								cu = 1 / math.Sqrt2
							}

							cv := 1.0
							if v == 0 {
								cv = 1 / math.Sqrt2
							}
							iconvVal += cu * cv * float64(pixVal) * cosVal1 * cosVal2
						}
					}
					idctVal := (iconvVal) / 4
					if idctVal > 255.0 {
						idctVal = 255.0
					}
					if idctVal < 0.0 {
						idctVal = 0.0
					}
					idctResult[yi+y][xi+x] = idctVal
				}
			}
		}
	}
	return idctResult
}

func rgb2ycbcr(r, g, b float64) (y, cb, cr float64) {
	y = 0.299*r + 0.5870*g + 0.114*b
	cb = -0.1687*r - 0.3313*g + 0.5*b + 128
	cr = 0.5*r - 0.4187*g - 0.0813*b + 128
	return y, cb, cr
}

func ycbcr2rgb(y, cb, cr float64) (r, g, b uint8) {
	rfloat64 := y + (cr-128)*1.402
	if rfloat64 > 255.0 {
		rfloat64 = 255.0
	}
	if rfloat64 < 0.0 {
		rfloat64 = 0.0
	}
	r = uint8(rfloat64)

	gfloat64 := y - (cb-128)*0.3441 - (cr-128)*0.7139
	if gfloat64 > 255.0 {
		gfloat64 = 255.0
	}
	if gfloat64 < 0.0 {
		gfloat64 = 0.0
	}
	g = uint8(gfloat64)

	bfloat64 := y + (cb-128)*1.7718
	if bfloat64 > 255.0 {
		bfloat64 = 255.0
	}
	if bfloat64 < 0.0 {
		bfloat64 = 0.0
	}
	b = uint8(bfloat64)
	return r, g, b
}

func quantization(inputArray, quantizationMatrix [][]float64, blockSize int) [][]float64 {
	H := len(inputArray)
	W := len(inputArray[0][:])

	quantizedArray := make([][]float64, H)
	for x := range quantizedArray {
		quantizedArray[x] = make([]float64, W)
	}

	for vi := 0; vi < H; vi += blockSize {
		for ui := 0; ui < W; ui += blockSize {
			for v := 0; v < blockSize; v++ {
				for u := 0; u < blockSize; u++ {
					quantizedArray[vi+v][ui+u] = math.Round(inputArray[vi+v][ui+u]/quantizationMatrix[v][u]) * quantizationMatrix[v][u]
				}
			}
		}
	}
	return quantizedArray
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

	H := jimg.Bounds().Size().Y
	W := jimg.Bounds().Size().X

	yValueArray := make([][]float64, H)
	cbValueArray := make([][]float64, H)
	crValueArray := make([][]float64, H)
	for x := range yValueArray {
		yValueArray[x] = make([]float64, W)
		cbValueArray[x] = make([]float64, W)
		crValueArray[x] = make([]float64, W)
	}

	// RGBをYCbCrに変換
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r32, g32, b32, _ := jimg.At(x, y).RGBA()
			r := float64(r32*0xFF) / 0xFFFF
			g := float64(g32*0xFF) / 0xFFFF
			b := float64(b32*0xFF) / 0xFFFF
			yc, cb, cr := rgb2ycbcr(r, g, b)
			yValueArray[y][x] = yc
			cbValueArray[y][x] = cb
			crValueArray[y][x] = cr
		}
	}

	quantizationMatrix1 := [][]float64{
		{16, 11, 10, 16, 24, 40, 51, 61},
		{12, 12, 14, 19, 26, 58, 60, 55},
		{14, 13, 16, 24, 40, 57, 69, 56},
		{14, 17, 22, 29, 51, 87, 80, 62},
		{18, 22, 37, 56, 68, 109, 103, 77},
		{24, 35, 55, 64, 81, 104, 113, 92},
		{49, 64, 78, 87, 103, 121, 120, 101},
		{72, 92, 95, 98, 112, 100, 103, 99},
	}

	quantizationMatrix2 := [][]float64{
		{17, 18, 24, 47, 99, 99, 99, 99},
		{18, 21, 26, 66, 99, 99, 99, 99},
		{24, 26, 56, 99, 99, 99, 99, 99},
		{47, 66, 99, 99, 99, 99, 99, 99},
		{99, 99, 99, 99, 99, 99, 99, 99},
		{99, 99, 99, 99, 99, 99, 99, 99},
		{99, 99, 99, 99, 99, 99, 99, 99},
		{99, 99, 99, 99, 99, 99, 99, 99},
	}

	blockSize := 8
	// 離散コサイン変換
	yDctResult := dct(yValueArray, blockSize)
	cbDctResult := dct(cbValueArray, blockSize)
	crDctResult := dct(crValueArray, blockSize)

	yDctResultQuan := quantization(yDctResult, quantizationMatrix1, blockSize)
	cbDctResultQuan := quantization(cbDctResult, quantizationMatrix2, blockSize)
	crDctResultQuan := quantization(crDctResult, quantizationMatrix2, blockSize)

	// 離散コサイン逆変換
	yIdctResultQuan := idct(yDctResultQuan, blockSize)
	cbIdctResultQuan := idct(cbDctResultQuan, blockSize)
	crIdctResultQuan := idct(crDctResultQuan, blockSize)

	quantinizedImg := image.NewRGBA(jimg.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r, g, b := ycbcr2rgb(yIdctResultQuan[y][x], cbIdctResultQuan[y][x], crIdctResultQuan[y][x])
			// r, g, b := ycbcr2rgb(yIdctResult[y][x], cbIdctResult[y][x], crIdctResult[y][x])
			quanColor := color.RGBA{r, g, b, 255}
			quantinizedImg.Set(x, y, quanColor)
		}
	}

	quanFile, err := os.Create("./answer_40.jpg")
	defer quanFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(quanFile, quantinizedImg, &jpeg.Options{100})

}
