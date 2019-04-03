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

func ck(in int) int {
	if in == 1 {
		return 1
	}
	return 0
}

func asta(in int) int {
	ret := 1 - int(math.Abs(float64(in)))
	return ret
}

func calcHilditchConnectivity2(neiborPixels [8]int) int {
	neiborCkPixels := [8]int{}
	for i, neiborPix := range neiborPixels {
		neiborCkPixels[i] = ck(neiborPix)
	}

	neiborCkAstaPixels := [8]int{}
	for i, neiborCk := range neiborCkPixels {
		neiborCkAstaPixels[i] = asta(neiborCk)
	}

	nc8 := ((neiborCkAstaPixels[0] - neiborCkAstaPixels[0]*neiborCkAstaPixels[1]*neiborCkAstaPixels[2]) +
		(neiborCkAstaPixels[2] - neiborCkAstaPixels[2]*neiborCkAstaPixels[3]*neiborCkAstaPixels[4]) +
		(neiborCkAstaPixels[4] - neiborCkAstaPixels[4]*neiborCkAstaPixels[5]*neiborCkAstaPixels[6]) +
		(neiborCkAstaPixels[6] - neiborCkAstaPixels[6]*neiborCkAstaPixels[7]*neiborCkAstaPixels[0]))
	return nc8
}

func calcHilditchConnectivity(array [8]int) int {
	nc8 := ((asta(array[0]) - asta(array[0])*asta(array[1])*asta(array[2])) +
		(asta(array[2]) - asta(array[2])*asta(array[3])*asta(array[4])) +
		(asta(array[4]) - asta(array[4])*asta(array[5])*asta(array[6])) +
		(asta(array[6]) - asta(array[6])*asta(array[7])*asta(array[0])))
	return nc8
}

func main() {
	thinningTargetFile, err := os.Open("./../Question_61_70/gazo.png")
	defer thinningTargetFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	thinningTargetImage, err := png.Decode(thinningTargetFile)
	if err != nil {
		log.Fatal(err)
	}

	H := thinningTargetImage.Bounds().Size().Y
	W := thinningTargetImage.Bounds().Size().X

	// カラー画像をグレイスケール画像に変換、画素値が0以上なら1をセットする
	graphicArray := make([][]int, H)
	thinningArray := make([][]int, H)
	for y := range graphicArray {
		graphicArray[y] = make([]int, W)
		thinningArray[y] = make([]int, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			c := color.GrayModel.Convert(thinningTargetImage.At(x, y))
			gray, _ := c.(color.Gray)
			if gray.Y > 0 {
				graphicArray[y][x] = 1
			} else {
				graphicArray[y][x] = 0
			}
		}
	}

	for _, row := range graphicArray {
		fmt.Println(row)
	}

	copy(thinningArray, graphicArray)
	for {
		counter := 0
		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				// fmt.Println()
				// fmt.Println("x", x, "y", y)
				// 注目画素g0とその8近傍の画素値を取得する
				leftIndex := int(math.Max(float64(x-1), 0))
				rightIndex := int(math.Min(float64(x+1), float64(W)-1))
				upIndex := int(math.Max(float64(y-1), 0))
				downIndex := int(math.Min(float64(y+1), float64(H)-1))

				// 8近傍画素の位置
				// g4 g3 g2
				// g5 g0 g1
				// g6 g7 g8
				g4 := graphicArray[upIndex][leftIndex]
				g3 := graphicArray[upIndex][x]
				g2 := graphicArray[upIndex][rightIndex]
				g5 := graphicArray[y][leftIndex]
				g0 := graphicArray[y][x]
				g1 := graphicArray[y][rightIndex]
				g6 := graphicArray[downIndex][leftIndex]
				g7 := graphicArray[downIndex][x]
				g8 := graphicArray[downIndex][rightIndex]

				// 条件１:図形画素である x0(x,y)=1
				if g0 == 1 {
					// OK
				} else {
					continue
				}
				fmt.Println("g0", g0)

				// 条件２:境界点である  注目画素の4近傍の合計が1以上
				boudary := asta(g1) + asta(g3) + asta(g5) + asta(g7)
				fmt.Println("4-neibor")
				fmt.Printf("  %d  \n%d %d %d\n  %d  \n", g3, g5, g0, g1, g7)
				fmt.Println("boudary", boudary)
				if boudary >= 1 {
					// OK
				} else {
					fmt.Println("cond 2 false")
					continue
				}

				// 条件３:端点の保存 x1〜x8の絶対値の合計が2以上
				// neibor8Pix := [8]int{g1, g2, g3, g4, g5, g6, g7, g8}
				neibor8Pix := [8]int{g4, g3, g2, g5, g1, g6, g7, g8}
				endPoints := 0
				for _, g := range neibor8Pix {
					endPoints += int(math.Abs(float64(g)))
				}
				fmt.Println("8-neibor")
				fmt.Printf("%d %d %d\n%d %d %d\n%d %d %d\n", g4, g3, g2, g5, g0, g1, g6, g7, g8)
				if endPoints >= 2 {
					// OK
				} else {
					fmt.Println("cond 3 false")
					continue
				}

				// 条件４:孤立点の保存 x0の8近傍に1が1つ以上存在する
				c4 := ck(g4)
				c3 := ck(g3)
				c2 := ck(g2)
				c5 := ck(g5)
				c0 := ck(g0)
				c1 := ck(g1)
				c6 := ck(g6)
				c7 := ck(g7)
				c8 := ck(g8)
				independentPoints := c1 + c2 + c3 + c4 + c5 + c6 + c7 + c8
				fmt.Println("independentPoints", independentPoints)
				fmt.Printf("%d %d %d\n%d %d %d\n%d %d %d\n", c4, c3, c2, c5, c0, c1, c6, c7, c8)
				if independentPoints >= 1 {
					// OK
				} else {
					fmt.Println("cond 4 false")
					continue
				}

				// 条件５:連結性の保存
				connectArray := [8]int{c1, c2, c3, c4, c5, c6, c7, c8}

				// nc8 := (asta(c1) - asta(c1)*asta(c2)*asta(c3)) +
				// 	(asta(c3) - asta(c3)*asta(c4)*asta(c5)) +
				// 	(asta(c5) - asta(c5)*asta(c6)*asta(c7)) +
				// 	(asta(c7) - asta(c7)*asta(c8)*asta(c1))
				nc8 := calcHilditchConnectivity(connectArray)
				fmt.Println("nc8", nc8)
				fmt.Println("asta")
				fmt.Printf("%d %d %d\n%d %d %d\n%d %d %d\n", asta(c4), asta(c3),
					asta(c2), asta(c5), asta(c0), asta(c1), asta(c6), asta(c7), asta(c8))
				nc8_1 := calcHilditchConnectivity2(neibor8Pix)
				fmt.Println("nc8_1", nc8_1)

				if nc8 == 1 {
					// continue
				} else {
					fmt.Println("cond 5 false")
					continue
				}

				// 条件6 線幅２の線分の片側だけを削除する
				// i := 1~8の全てにおいて以下の6-1、6-2のどちらかの条件が成り立つ
				cond6 := 0
				for i, neiborPix := range neibor8Pix {
					// 条件6-1
					if neiborPix != -1 {
						cond6++
					} else {
						// 条件6-2
						neibor8PixCopy := neibor8Pix
						neibor8PixCopy[i] = -1
						nc8 = calcHilditchConnectivity2(neibor8PixCopy)
						if nc8 == 1 {
							// OK
							cond6++
						}
					}
				}

				if cond6 == 8 {
					thinningArray[y][x] = -1
				}

			}
		}

		// -1になった画素を0に変える
		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				if thinningArray[y][x] == -1 {
					thinningArray[y][x] = 0
					fmt.Println(x, y, "replace -1 to 0")
					counter++
				}
			}
		}

		thinningImageTrans := image.NewGray(thinningTargetImage.Bounds())
		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				gray := color.Gray{uint8(thinningArray[y][x]) * 255}
				thinningImageTrans.Set(x, y, gray)
			}
		}

		thinningTransFile, err := os.Create(fmt.Sprintf("./answer_64_%d.png", counter))
		defer thinningTransFile.Close()
		if err != nil {
			log.Fatal(err)
		}
		png.Encode(thinningTransFile, thinningImageTrans)

		fmt.Println("counter", counter)
		if counter == 0 {
			break
		} else {
			for _, row := range thinningArray {
				fmt.Println(row)
			}
			// 細線化した画像を図形画像とする
			copy(graphicArray, thinningArray)
		}
	}

	thinningImage := image.NewGray(thinningTargetImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			gray := color.Gray{uint8(thinningArray[y][x]) * 255}
			thinningImage.Set(x, y, gray)
		}
	}

	thinningFile, err := os.Create("./answer_64.png")
	defer thinningFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(thinningFile, thinningImage)

}
