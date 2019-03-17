package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func rgb2ycbcr(r, g, b float64) (y, cb, cr float64) {
	y = 0.299*r + 0.5870*g + 0.114*b
	cb = -0.1687*r - 0.3313*g + 0.5*b + 128
	cr = 0.5*r - 0.4187*g - 0.0813*b + 128
	return y, cb, cr
}

func ycbcr2rgb(y, cb, cr float64) (r, g, b uint8) {
	r = uint8(y + (cr-128)*1.402)
	g = uint8(y - (cb-128)*0.3441 - (cr-128)*0.7139)
	b = uint8(y + (cb-128)*1.7718)
	return r, g, b
}

func dct(array2D [][]float64, blockSize int) [][]float64 {
	H := len(array2D)
	W := len(array2D[0][:])

	// H := array2D.Bounds().Size().Y
	// W := array2D.Bounds().Size().X

	// 離散コサイン変換
	dctResult := make([][]float64, H)
	for x := range dctResult {
		dctResult[x] = make([]float64, W)
	}

	// blockSize := 8
	for vi := 0; vi < H; vi += blockSize {
		for ui := 0; ui < W; ui += blockSize {
			for v := 0; v < blockSize; v++ {
				for u := 0; u < blockSize; u++ {
					convVal := 0.0
					for y := 0; y < blockSize; y++ {
						for x := 0; x < blockSize; x++ {
							// pixVal := array2D.GrayAt(ui+x, vi+y).Y
							pixVal := array2D[vi+y][ui+x]
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

func idct(dctResult [][]float64, blockSize int) [][]float64 {
	H := len(dctResult)
	W := len(dctResult[0][:])

	idctResult := make([][]float64, H)
	for x := range idctResult {
		idctResult[x] = make([]float64, W)
	}

	// 離散コサイン逆変換
	for yi := 0; yi < H; yi += blockSize {
		for xi := 0; xi < H; xi += blockSize {
			for y := 0; y < blockSize; y++ {
				for x := 0; x < blockSize; x++ {
					iconvVal := 0.0
					for v := 0; v < blockSize; v++ {
						for u := 0; u < blockSize; u++ {
							pixVal := dctResult[yi+v][xi+u]
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
					idctResult[yi+y][xi+x] = idctVal
				}
			}
		}
	}
	return idctResult
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

	// グレイスケールに変換
	grayArray := make([][]float64, H)
	for x := range grayArray {
		grayArray[x] = make([]float64, W)
	}

	grayImg := image.NewGray(jimg.Bounds())
	for height := 0; height < jimg.Bounds().Size().Y; height++ {
		for width := 0; width < jimg.Bounds().Size().X; width++ {
			ycbcr := jimg.At(width, height)
			r, g, b, _ := ycbcr.RGBA()
			var graycolor color.Gray16
			graycolor.Y = uint16(0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b))
			grayImg.Set(width, height, graycolor)
			grayArray[height][width] = 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
		}
	}

	dctResult := dct(grayArray, 8)
	idctResult := idct(dctResult, 8)

	idctImg := image.NewGray(jimg.Bounds())
	for y, rowArray := range idctResult {
		for x, invDftVal := range rowArray {
			var invDftGray color.Gray
			invDftGray.Y = uint8(invDftVal)
			idctImg.Set(x, y, invDftGray)
		}
	}

	degreIdctFile, err := os.Create("./answer_40.jpg")
	defer degreIdctFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(degreIdctFile, idctImg, &jpeg.Options{100})

}
