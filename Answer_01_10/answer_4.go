package main

import (
	"fmt"
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
	grayImg := image.NewGray(jimg.Bounds())
	for height := 0; height < jimg.Bounds().Size().Y; height++ {
		for width := 0; width < jimg.Bounds().Size().X; width++ {
			ycbcr := jimg.At(width, height)
			r32, g32, b32, _ := ycbcr.RGBA()
			// 32bit画像から8bit画像に変換
			r8 := uint8((float64(r32) / 0xFFFF) * 0xFF)
			g8 := uint8((float64(g32) / 0xFFFF) * 0xFF)
			b8 := uint8((float64(b32) / 0xFFFF) * 0xFF)

			var graycolor color.Gray
			// カラーからグレースケールに変換
			grayfloat64 := 0.2126*float64(r8) + 0.7152*float64(g8) + 0.0722*float64(b8)
			grayuint8 := uint8(grayfloat64)
			graycolor.Y = grayuint8
			grayImg.Set(width, height, graycolor)
		}
	}

	// 大津の2値化
	// ヒストグラムを計算
	var hist [255]uint16
	pixVal := 0
	for height := 0; height < grayImg.Bounds().Size().Y; height++ {
		for width := 0; width < grayImg.Bounds().Size().X; width++ {
			grayIndex := grayImg.GrayAt(width, height).Y
			hist[grayIndex]++
			pixVal += int(grayIndex)
		}
	}

	// ヒストグラムから最適なしきい値ootu_threshを計算
	var thresh int
	pAll := grayImg.Bounds().Size().X * grayImg.Bounds().Size().Y
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
		// fmt.Println("pixVal", pixVal)
		// fmt.Println("pixValAve", float64(pixVal)/float64(pAll))
		// fmt.Println("p0", p0, "m0", m0, "r0", r0)
		// fmt.Println("p1", p1, "m1", m1, "r1", r1)
		// fmt.Println("r0 + r1", r0+r1)
		// fmt.Println("thresh", thresh, "sb2", sb2)
		if sb2 > sb2Max {
			sb2Max = sb2
			ootuThresh = thresh
		}
	}
	fmt.Println("ootu thresh", ootuThresh)
	for height := 0; height < grayImg.Bounds().Size().Y; height++ {
		for width := 0; width < grayImg.Bounds().Size().X; width++ {
			gray := grayImg.GrayAt(width, height)
			if gray.Y > uint8(ootuThresh) {
				gray.Y = 255
			} else {
				gray.Y = 0
			}
			grayImg.Set(width, height, gray)
		}
	}

	grayfile, err := os.Create("./answer_4.jpg")
	defer grayfile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(grayfile, grayImg, &jpeg.Options{100})

}
