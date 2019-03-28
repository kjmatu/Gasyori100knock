package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

func main() {
	// seqFile, err := os.Open("./../sample/label.png")
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
	// lookUpTable := make([]int, 20)
	lookUpTable := map[int]int{}
	for i := 0; i < len(lookUpTable); i++ {
		lookUpTable[i] = i
	}
	label := 0
	labelImage := image.NewGray(seqImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r, _, _, _ := seqImage.At(x, y).RGBA()
			if r == 0 {
				continue
			}

			upLabel := labelImage.GrayAt(x, y-1).Y
			leftLabel := labelImage.GrayAt(x-1, y).Y
			if x == 55 && y == 100 {
				fmt.Println("label", labelImage.GrayAt(x, y).Y)
				fmt.Println("upLabel", upLabel)
				fmt.Println("leftLabel", leftLabel)
			}

			if upLabel == 0 && leftLabel == 0 {
				label++
				labelImage.Set(x, y, color.Gray{uint8(label)})
				lookUpTable[label] = label
			} else {
				minLabel := math.Min(float64(upLabel), float64(leftLabel))
				maxLabel := math.Max(float64(upLabel), float64(leftLabel))
				if minLabel == 0 {
					labelImage.Set(x, y, color.Gray{uint8(maxLabel)})
				} else {
					if int(maxLabel) != int(minLabel) {
					}
					_, ok := lookUpTable[int(maxLabel)]
					if ok {
						lookUpTable[int(maxLabel)] = int(minLabel)
					}
					labelImage.Set(x, y, color.Gray{uint8(minLabel)})
					// if int(maxLabel) == 3 && minLabel == 3 {
					// 	fmt.Println("label", labelImage.GrayAt(x, y).Y)
					// 	fmt.Println("upLabel", upLabel)
					// 	fmt.Println("leftLabel", leftLabel)
					// 	os.Exit(0)
					// }
					// labelImage.Set(x, y, color.Gray{uint8(minLabel)})
					// lookUpTable[int(maxLabel)] = int(minLabel)
				}
			}
			// if x == 95 && y == 76 {
			// 	fmt.Println("LUT", lookUpTable)
			// }
		}
	}

	fmt.Println(lookUpTable)

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			labelValue := labelImage.GrayAt(x, y).Y
			if labelValue == 2 || labelValue == 3 {
				// fmt.Println(lookUpTable)
				// fmt.Printf("%d -> %d \n", labelValue, lookUpTable[labelValue])
			}
			labelImage.Set(x, y, color.Gray{uint8(lookUpTable[int(labelValue)])})
		}
	}
	color := []color.NRGBA{color.NRGBA{0, 0, 0, 255},
		color.NRGBA{255, 0, 0, 255}, color.NRGBA{0, 255, 0, 255},
		color.NRGBA{0, 0, 255, 255}, color.NRGBA{255, 255, 0, 255},
		color.NRGBA{255, 0, 255, 255}, color.NRGBA{0, 255, 255, 255},
		color.NRGBA{0, 255, 255, 255}}

	labeledImage := image.NewNRGBA(seqImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			colorIndex := labelImage.GrayAt(x, y).Y
			// fmt.Println("colorIndex", colorIndex)
			labeledImage.Set(x, y, color[colorIndex])
		}
	}
	coloredFile, err := os.Create("./answer_58.png")
	defer coloredFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(coloredFile, labeledImage)
}
