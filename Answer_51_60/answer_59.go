package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

func main() {
	seqFile, err := os.Open("./../Question_51_60/seg.png")
	defer seqFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	seqImage, err := png.Decode(seqFile)
	if err != nil {
		log.Fatal(err)
	}

	H := seqImage.Bounds().Size().Y
	W := seqImage.Bounds().Size().X

	lookUpTable := [20]int{}
	label := 0
	labelImage := image.NewGray(seqImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r, _, _, _ := seqImage.At(x, y).RGBA()
			if r == 0 {
				continue
			}
			upYindex := int(math.Max(float64(y-1), 0))
			leftXindex := int(math.Max(float64(x-1), 0))

			// 左、上、左上、右上のラベルを取得する
			upLabel := labelImage.GrayAt(x, upYindex).Y
			leftLabel := labelImage.GrayAt(leftXindex, y).Y
			upperLeftLabel := labelImage.GrayAt(leftXindex, upYindex).Y
			upperRightLabel := labelImage.GrayAt(x+1, upYindex).Y

			if upLabel == 0 && leftLabel == 0 && upperLeftLabel == 0 && upperRightLabel == 0 {
				label++
				labelImage.Set(x, y, color.Gray{uint8(label)})
				lookUpTable[label] = label
			} else {
				minLabel1 := math.Min(float64(upLabel), float64(leftLabel))
				minLabel2 := math.Min(float64(upperLeftLabel), float64(upperRightLabel))
				minLabel := math.Min(minLabel1, minLabel2)

				maxLabel1 := math.Max(float64(upLabel), float64(leftLabel))
				maxLabel2 := math.Max(float64(upperLeftLabel), float64(upperRightLabel))
				maxLabel := math.Max(maxLabel1, maxLabel2)

				if minLabel == 0 {
					labelImage.Set(x, y, color.Gray{uint8(maxLabel)})
				} else if int(minLabel) == int(maxLabel) {
					labelImage.Set(x, y, color.Gray{uint8(minLabel)})
				} else {
					lookUpTable[int(maxLabel)] = int(minLabel)
					labelImage.Set(x, y, color.Gray{uint8(minLabel)})
				}
			}
		}
	}

	changeLookUpTable := [20]int{}

	// LookUpTableで変更されるラベルが収束するまで変更する
	for {
		for i := 0; i < len(lookUpTable); i++ {
			// dst labelが指すsrc labelをchangeLookUpTableに格納する
			changeLookUpTable[i] = lookUpTable[lookUpTable[i]]
		}

		if lookUpTable != changeLookUpTable {
			lookUpTable = changeLookUpTable
			continue
		} else {
			break
		}
	}

	// labelを詰め直す
	count := 0
	for i := 0; i < len(changeLookUpTable)-1; i++ {
		lookUpTable[i] = count
		if changeLookUpTable[i] != changeLookUpTable[i+1] {
			count++
		}
		if changeLookUpTable[i+1] == 0 {
			break
		}
	}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			labelValue := labelImage.GrayAt(x, y).Y
			labelImage.Set(x, y, color.Gray{uint8(lookUpTable[int(labelValue)])})
		}
	}
	color := []color.NRGBA{color.NRGBA{0, 0, 0, 255},
		color.NRGBA{255, 0, 0, 255}, color.NRGBA{0, 255, 0, 255},
		color.NRGBA{0, 0, 255, 255}, color.NRGBA{0, 255, 255, 255},
		color.NRGBA{255, 0, 255, 255}, color.NRGBA{0, 255, 255, 255},
		color.NRGBA{0, 255, 255, 255}}

	labeledImage := image.NewNRGBA(seqImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			colorIndex := labelImage.GrayAt(x, y).Y
			labeledImage.Set(x, y, color[colorIndex])
		}
	}
	coloredFile, err := os.Create("./answer_59.png")
	defer coloredFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(coloredFile, labeledImage)
}
