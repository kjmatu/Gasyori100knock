package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"math/cmplx"
	"os"
)

func dft(grayImage *image.Gray) [][]complex128 {
	W := grayImage.Bounds().Size().X
	H := grayImage.Bounds().Size().Y
	dftResult := make([][]complex128, W)
	for index := 0; index < W; index++ {
		dftResult[index] = make([]complex128, H)
	}

	N := math.Sqrt(float64(W * H))
	for v, dftRow := range dftResult {
		for u := range dftRow {
			var dftVal complex128
			for y := 0; y < H; y++ {
				for x := 0; x < W; x++ {
					pixVal := grayImage.GrayAt(x, y).Y
					imag := -2 * math.Pi * (float64(u*x)/float64(W) + float64(v*y)/float64(H))
					exponentVal := complex(0, imag)
					eulerVal := cmplx.Exp(exponentVal)
					dftVal += complex(float64(pixVal), 0) * eulerVal
				}
			}
			dftResult[v][u] = dftVal / complex(N, 0)
		}
	}
	return dftResult
}

func invDft(imageDft [][]complex128) [][]float64 {
	W := len(imageDft[0][:])
	H := len(imageDft)
	N := math.Sqrt(float64(W * H))

	invDftResult := make([][]float64, W)
	for index := 0; index < W; index++ {
		invDftResult[index] = make([]float64, H)
	}

	for y := 0; y < W; y++ {
		for x := 0; x < H; x++ {
			var invDftVal complex128
			for v, dftRow := range imageDft {
				for u, dftVal := range dftRow {
					imag := 2 * math.Pi * (float64(u*x)/float64(W) + float64(v*y)/float64(H))
					exponentVal := complex(0, imag)
					eulerVal := cmplx.Exp(exponentVal)
					invDftVal += dftVal * eulerVal
				}
			}
			invDftAbs := cmplx.Abs(invDftVal) / N
			if invDftAbs > 255.0 {
				invDftAbs = 255.0
			}
			invDftResult[y][x] = invDftAbs
		}
	}
	return invDftResult
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

	dftResult := dft(grayImg)

	ampMax := math.Inf(-1)
	for _, dftRow := range dftResult {
		for _, dftVal := range dftRow {
			dftAmp := cmplx.Abs(dftVal)
			if dftAmp > ampMax {
				ampMax = dftAmp
			}
		}
	}

	ampImg := image.NewGray(jimg.Bounds())

	for height := 0; height < ampImg.Bounds().Size().Y; height++ {
		for width := 0; width < ampImg.Bounds().Size().X; width++ {
			var graycolor color.Gray
			graycolor.Y = uint8(cmplx.Abs(dftResult[height][width]) * 255 / ampMax)
			ampImg.Set(width, height, graycolor)
		}
	}

	ampFile, err := os.Create("./answer_32_ps.jpg")
	defer ampFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(ampFile, ampImg, &jpeg.Options{100})

	invDftResult := invDft(dftResult)

	invDftImg := image.NewGray(jimg.Bounds())

	for y, rowArray := range invDftResult {
		for x, invDftVal := range rowArray {
			var invDftGray color.Gray
			invDftGray.Y = uint8(invDftVal)
			invDftImg.Set(x, y, invDftGray)
		}
	}

	invDftFile, err := os.Create("./answer_32.jpg")
	defer invDftFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(invDftFile, invDftImg, &jpeg.Options{100})

}
