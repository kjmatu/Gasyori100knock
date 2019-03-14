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

func siftDftValue(dftValue [][]complex128) [][]complex128 {
	W := len(dftValue)
	halfW := int(float64(W) / 2)
	H := len(dftValue[0][:])
	halfH := int(float64(H) / 2)

	siftDftResult := make([][]complex128, W)
	for index := 0; index < W; index++ {
		siftDftResult[index] = make([]complex128, H)
	}

	for y, rowArray := range dftValue {
		for x, dftValue := range rowArray {
			// x軸の左半分
			if x >= 0 && x < halfW {
				// y軸の上半分
				if y >= 0 && y < halfH {
					// 中心から左上の領域を右下へ移動
					siftDftResult[y+halfH][x+halfW] = dftValue
				} else {
					// y軸の下半分
					// 中心から左下の領域を右上へ移動
					siftDftResult[y-halfH][x+halfW] = dftValue
				}
			} else {
				// x軸の右半分
				if y >= 0 && y < halfH {
					// 中心から右上の領域を左下へ移動
					siftDftResult[y+halfH][x-halfW] = dftValue
				} else {
					// 中心から右下の領域を左上へ移動
					siftDftResult[y-halfH][x-halfW] = dftValue
				}
			}
		}
	}
	return siftDftResult
}

func highPassFileter(siftDftResult [][]complex128) [][]complex128 {
	H := len(siftDftResult)
	halfH := int(float64(H) / 2)
	W := len(siftDftResult[0][:])
	halfW := int(float64(W) / 2)

	filterdSiftDftResult := make([][]complex128, W)
	for index := 0; index < W; index++ {
		filterdSiftDftResult[index] = make([]complex128, H)
	}

	// ハイパスフィルタ処理
	for y, rowArray := range siftDftResult {
		for x, dftValue := range rowArray {
			dx := float64(x - halfW)
			dy := float64(y - halfH)
			distanceFromCenter := math.Hypot(dx, dy)
			// 中心部にある高周波成分をカットする
			if distanceFromCenter <= float64(halfW)*0.2 {
				filterdSiftDftResult[y][x] = 0
			} else {
				filterdSiftDftResult[y][x] = dftValue
			}
		}
	}
	return filterdSiftDftResult
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

	// 画像をフーリエ変換する
	dftResult := dft(grayImg)

	// 中心部が高周波領域の配列から中心部が低周波領域の配列にシフトする
	siftDftResult := siftDftValue(dftResult)

	// ハイパスフィルタをかける
	filteredSiftDftResult := highPassFileter(siftDftResult)

	// 中心部が低周波領域の配列から中心部が高周波領域の配列にシフトする
	filterdDftResult := siftDftValue(filteredSiftDftResult)

	// フーリエ逆変換をかけて画像に戻す
	filterdInvDftResult := invDft(filterdDftResult)

	highPassImg := image.NewGray(jimg.Bounds())
	for y, rowArray := range filterdInvDftResult {
		for x, invDftVal := range rowArray {
			var invDftGray color.Gray
			invDftGray.Y = uint8(invDftVal)
			highPassImg.Set(x, y, invDftGray)
		}
	}

	highPassFile, err := os.Create("./answer_34.jpg")
	defer highPassFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(highPassFile, highPassImg, &jpeg.Options{100})

}
