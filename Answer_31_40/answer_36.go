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
	file, err := os.Open("./../assets/imori.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// グレイスケールに変換
	grayImg := image.NewGray(jimg.Bounds())
	for height := 0; height < jimg.Bounds().Size().Y; height++ {
		for width := 0; width < jimg.Bounds().Size().X; width++ {
			ycbcr := jimg.At(width, height)
			r, g, b, _ := ycbcr.RGBA()
			var graycolor color.Gray16
			graycolor.Y = uint16(0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b))
			grayImg.Set(width, height, graycolor)
		}
	}
	grayfile, err := os.Create("./imori_gray.jpg")
	defer grayfile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(grayfile, grayImg, &jpeg.Options{100})

	H := grayImg.Bounds().Size().Y
	W := grayImg.Bounds().Size().X

	dctResult := make([][]float64, H)
	for x := range dctResult {
		dctResult[x] = make([]float64, W)
	}

	// 離散コサイン変換
	blockSize := 8
	for vi := 0; vi < H; vi += blockSize {
		for ui := 0; ui < W; ui += blockSize {
			for v := 0; v < blockSize; v++ {
				for u := 0; u < blockSize; u++ {
					convVal := 0.0
					for y := 0; y < blockSize; y++ {
						for x := 0; x < blockSize; x++ {
							pixVal := grayImg.GrayAt(ui+x, vi+y).Y
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
					idctResult[yi+y][xi+x] = idctVal
				}
			}
		}
	}

	idctImg := image.NewGray(jimg.Bounds())
	for y, rowArray := range idctResult {
		for x, invDftVal := range rowArray {
			var invDftGray color.Gray
			invDftGray.Y = uint8(invDftVal)
			idctImg.Set(x, y, invDftGray)
		}
	}

	idctFile, err := os.Create("./answer_36.jpg")
	defer idctFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(idctFile, idctImg, &jpeg.Options{100})

}
