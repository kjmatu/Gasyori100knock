package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"sort"
)

type any []interface{}

func flatten(in, result any) any {
	for _, x := range in {
		s, ok := x.(int)
		if ok {
			result = append(result, s)
		} else {
			result = flatten(x.(any), result)
		}
	}
	return result
}

func main() {
	file, err := os.Open("./../Question_41_50/answers/answer_44.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpegImage, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	houghArray := make([][]int, jpegImage.Bounds().Size().Y)
	nmsArray := make([][]int, jpegImage.Bounds().Size().Y)
	for y := range houghArray {
		houghArray[y] = make([]int, jpegImage.Bounds().Size().X)
		nmsArray[y] = make([]int, jpegImage.Bounds().Size().X)
		for x := range houghArray[y] {
			r, _, _, _ := jpegImage.At(x, y).RGBA()
			houghArray[y][x] = int((r * 0xFF) / 0xFFFF)
		}
	}

	// fmt.Println(len(houghArray)-1, len(houghArray[0][:])-1)
	for y := 1; y < len(houghArray)-1; y++ {
		for x := 1; x < len(houghArray[0][:])-1; x++ {
			center := houghArray[y][x]
			for neiborY := -1; neiborY <= 1; neiborY++ {
				for neiborX := -1; neiborX <= 1; neiborX++ {
					refY := y + neiborY
					refX := x + neiborX
					// fmt.Println(refX, refY)
					neibor := houghArray[refY][refX]
					if center > neibor {
						nmsArray[y][x] = center
					} else {
						nmsArray[y][x] = 0
					}
				}
			}
		}
	}

	nmsFlatten := make([]int, jpegImage.Bounds().Size().Y*jpegImage.Bounds().Size().X)
	for y, row := range nmsArray {
		for x, value := range row {
			nmsFlatten[y*len(nmsArray[0][:])+x] = value
		}
	}
	fmt.Println(nmsFlatten)
	sort.Sort(sort.Reverse(sort.Ints(nmsFlatten)))
	fmt.Println(nmsFlatten)

	nmsImage := image.NewGray(jpegImage.Bounds())
	for y := 0; y < nmsImage.Bounds().Size().Y; y++ {
		for x := 0; x < nmsImage.Bounds().Size().X; x++ {
			nmsColor := color.Gray{uint8(nmsArray[y][x])}
			nmsImage.Set(x, y, nmsColor)
		}
	}

	nmsFile, err := os.Create("./answer_45.jpg")
	defer nmsFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(nmsFile, nmsImage, &jpeg.Options{100})

}
