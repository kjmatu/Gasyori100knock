package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
	"sort"
)

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

	H := jpegImage.Bounds().Size().Y
	W := jpegImage.Bounds().Size().X
	houghArray := make([][]int, H)
	nmsArray := make([][]int, H)
	for y := range houghArray {
		houghArray[y] = make([]int, W)
		nmsArray[y] = make([]int, W)
		for x := range houghArray[y] {
			r, _, _, _ := jpegImage.At(x, y).RGBA()
			houghArray[y][x] = int((r * 0xFF) / 0xFFFF)
		}
	}

	nmsFlattenArray := []int{}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			leftIndex := int(math.Max(float64(x-1), 0))
			rightIndex := int(math.Min(float64(x+1), float64(W-1)))
			upIndex := int(math.Max(float64(y-1), 0))
			downIndex := int(math.Min(float64(y+1), float64(H-1)))

			// 注目している得票数
			targetVotes := houghArray[y][x]
			// 8近傍の得票数
			v0 := houghArray[upIndex][leftIndex]
			v1 := houghArray[upIndex][x]
			v2 := houghArray[upIndex][rightIndex]
			v3 := houghArray[y][leftIndex]
			v4 := houghArray[y][rightIndex]
			v5 := houghArray[downIndex][leftIndex]
			v6 := houghArray[downIndex][x]
			v7 := houghArray[downIndex][rightIndex]
			neiborVotesArray := []int{v0, v1, v2, v3, v4, v5, v6, v7}

			zeroFlag := false
			for _, v := range neiborVotesArray {
				if targetVotes < v {
					zeroFlag = true
					break
				}
			}
			if zeroFlag {
				nmsArray[y][x] = 0
			} else {
				nmsArray[y][x] = houghArray[y][x]
			}

		}
		nmsFlattenArray = append(nmsFlattenArray, nmsArray[y][:]...)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(nmsFlattenArray)))
	top10Array := nmsFlattenArray[:10]

	maximumPointArray := []image.Point{}
	nmsImage := image.NewGray(jpegImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			houghVal := houghArray[y][x]
			top10Flag := false
			for _, top10 := range top10Array {
				if houghVal == top10 {
					top10Flag = true
					maximumPointArray = append(maximumPointArray, image.Point{x, y})
					break
				}
			}
			if top10Flag {
				nmsImage.SetGray(x, y, color.Gray{255})
			}
		}
	}

	thorinoFile, err := os.Open("./../Question_41_50/thorino.jpg")
	defer thorinoFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	thorinoImage, err := jpeg.Decode(thorinoFile)
	if err != nil {
		log.Fatal(err)
	}

	H = thorinoImage.Bounds().Size().Y
	W = thorinoImage.Bounds().Size().X

	houghInvImage := image.NewRGBA(thorinoImage.Bounds())
	draw.Draw(houghInvImage, thorinoImage.Bounds(), thorinoImage, thorinoImage.Bounds().Min, draw.Src)

	// Hought逆変換を行い、直線を描画する
	for _, point := range maximumPointArray {
		theta := math.Pi / 180 * float64(point.X)
		rho := float64(point.Y)

		for x := 0; x < W; x++ {
			if math.Sin(theta) != 0 {
				y := int(-(math.Cos(theta)/math.Sin(theta))*float64(x) + rho/math.Sin(theta))
				if y >= H || y < 0 {
					continue
				}
				houghInvImage.Set(x, y, color.RGBA{255, 0, 0, 255})

			}
		}

		for y := 0; y < H; y++ {
			if math.Cos(theta) != 0 {
				x := int(-(math.Sin(theta)/math.Cos(theta))*float64(y) + rho/math.Cos(theta))
				if x >= W || x <= 0 {
					continue
				}
				houghInvImage.Set(x, y, color.RGBA{255, 0, 0, 255})
			}

		}
	}
	houghInvFile, err := os.Create("./answer_46.jpg")
	defer houghInvFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(houghInvFile, houghInvImage, &jpeg.Options{100})

}
