package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

// HLine draws a horizontal line
func hLine(img *image.NRGBA, x1, y, x2 int, col color.Color) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// VLine draws a veritcal line
func vLine(img *image.NRGBA, x1, y1, y2 int, col color.Color) {
	for ; y1 <= y2; y1++ {
		img.Set(x1, y1, col)
	}
}

func primaryFunc(x1, y1, x2, y2 int) {

}

func diagonalLine(img *image.Gray, x1, y1, x2, y2 int, col color.Gray) {
	a := float64(y2-y1) / float64(x2-x1)
	b := float64(y1)
	for x := x1; x < x2; x++ {
		y := a*float64(x-x1) + b
		img.SetGray(int(x), int(y), col)
	}
}

func main() {
	srcFile, err := os.Open("./../Question_61_70/imori.jpg")
	defer srcFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	srcImage, err := jpeg.Decode(srcFile)
	if err != nil {
		log.Fatal(err)
	}

	H := srcImage.Bounds().Size().Y
	W := srcImage.Bounds().Size().X

	// カラー画像をグレイスケール画像に変換
	grayArray := make([][]float64, H)
	for y := range grayArray {
		grayArray[y] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r, g, b, _ := srcImage.At(x, y).RGBA()
			gray := float64(r)*0.2126 + float64(g)*0.7152 + float64(b)*0.0722
			grayArray[y][x] = (gray * 0xFF) / 0xFFFF
		}
	}

	// x,y方向の輝度勾配を求める
	gx := make([][]float64, H)
	gy := make([][]float64, H)
	for y := range gx {
		gx[y] = make([]float64, W)
		gy[y] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			leftIndex := int(math.Max(float64(x-1), 0))
			rightIndex := int(math.Min(float64(x+1), float64(W-1)))
			gx[y][x] = grayArray[y][rightIndex] - grayArray[y][leftIndex]

			upIndex := int(math.Max(float64(y-1), 0))
			downIndex := int(math.Min(float64(y+1), float64(H-1)))
			gy[y][x] = grayArray[downIndex][x] - grayArray[upIndex][x]
		}
	}

	// x,y方向の輝度勾配から勾配強度と角度を求める
	mag := make([][]float64, H)
	angle := make([][]float64, H)
	for y := range mag {
		mag[y] = make([]float64, W)
		angle[y] = make([]float64, W)
	}

	magMax := -1.0
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			mag[y][x] = math.Hypot(gy[y][x], gx[y][x])
			if magMax < mag[y][x] {
				magMax = mag[y][x]
			}

			angle[y][x] = (180 * math.Atan(gy[y][x]/gx[y][x])) / math.Pi
			if angle[y][x] < 0 {
				angle[y][x] += 180
			} else if math.IsNaN(angle[y][x]) {
				angle[y][x] = 0
			}
		}
	}

	// 勾配角度を 0~180度で9分割した値に量子化する
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			ang := angle[y][x]
			if ang >= 0 && ang < 20 {
				angle[y][x] = 0
			} else if ang >= 20 && ang < 40 {
				angle[y][x] = 1
			} else if ang >= 40 && ang < 60 {
				angle[y][x] = 2
			} else if ang >= 60 && ang < 80 {
				angle[y][x] = 3
			} else if ang >= 80 && ang < 100 {
				angle[y][x] = 4
			} else if ang >= 100 && ang < 120 {
				angle[y][x] = 5
			} else if ang >= 120 && ang < 140 {
				angle[y][x] = 6
			} else if ang >= 140 && ang < 160 {
				angle[y][x] = 7
			} else if ang >= 160 && ang <= 180 {
				angle[y][x] = 8
			}
		}
	}

	colorList := []color.RGBA{color.RGBA{0, 0, 255, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 255, 255},
		color.RGBA{255, 0, 255, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{0, 127, 127, 255},
		color.RGBA{127, 0, 127, 255},
		color.RGBA{127, 127, 0, 255}}

	gradColorImage := image.NewRGBA(srcImage.Bounds())
	magnitudeImage := image.NewGray(srcImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			gradient := angle[y][x]
			magnitude := mag[y][x]
			gradColorImage.Set(x, y, colorList[int(gradient)])
			magnitudeImage.Set(x, y, color.Gray{uint8(magnitude * 255 / magMax)})
		}
	}

	// answer_67ここから
	cellSize := 8
	cellH := H / cellSize
	cellW := W / cellSize
	hist := make([][][]float64, cellH)
	for i := range hist {
		hist[i] = make([][]float64, cellW)
		for j := range hist[i] {
			hist[i][j] = make([]float64, 9)
		}
	}

	// cellX,Yは16x16の領域に分割されたCELLを指定するIndex
	for cellY := 0; cellY < cellH; cellY++ {
		for cellX := 0; cellX < cellW; cellX++ {
			// j,iはCELL内の8x8ピクセルを指定するIndex
			for j := 0; j < cellSize; j++ {
				for i := 0; i < cellSize; i++ {
					// CELL内のピクセルIndexが元画像のどのピクセルIndexに対応するか計算
					refX := cellX*cellSize + i
					refY := cellY*cellSize + j
					hist[cellY][cellX][int(angle[refY][refX])] += mag[refY][refX]
				}
			}
		}
	}

	// answer_68ここから
	// 3x3のセルを1ブロックとして扱う
	epsilon := 1.0

	count := 0
	for cellY := 0; cellY < cellH; cellY++ {
		for cellX := 0; cellX < cellW; cellX++ {
			sum := 0.0
			// 対象セルを中心とした3x3セルを1ブロックとして正規化を行う
			for i := -1; i < 2; i++ {
				for j := -1; j < 2; j++ {
					if cellY+i >= cellH {
						continue
					}

					if cellY+i < 0 {
						continue
					}

					if cellX+j >= cellW {
						continue
					}
					if cellX+j < 0 {
						continue
					}

					for _, elm := range hist[cellY+i][cellX+j] {
						sum += elm * elm
					}
				}
			}
			sum += epsilon
			sum = math.Sqrt(sum)
			// fmt.Println(count, sum)
			count++

			for i := range hist[cellY][cellX] {
				hist[cellY][cellX][i] /= sum
			}
			// fmt.Println(hist[cellY][cellX])
		}
	}

	// answer_69ここから
	arrowHist := make([][][]float64, cellH)
	for i := range arrowHist {
		arrowHist[i] = make([][]float64, cellW)
		for j := range arrowHist[i] {
			arrowHist[i][j] = make([]float64, 9)
		}
	}

	arrowImage := image.NewGray(srcImage.Bounds())
	for cellY := 0; cellY < cellH; cellY++ {
		for cellX := 0; cellX < cellW; cellX++ {
			// CELLの中心を通る線を引く
			cellCenterX := cellX*cellSize + cellSize/2
			cellCenterY := cellY*cellSize + cellSize/2

			histMax := math.Inf(-1)
			for _, elm := range hist[cellY][cellX] {
				if elm > histMax {
					histMax = elm
				}
			}
			for i := range hist[cellY][cellX] {
				hist[cellY][cellX][i] /= histMax
			}

			for i, elm := range hist[cellY][cellX] {
				angle := (20 * float64(i)) / 180 * math.Pi
				dx := math.Cos(angle) * float64(cellSize) / 2
				dy := math.Sin(angle) * float64(cellSize) / 2
				x1 := cellCenterX - int(dx)
				y1 := cellCenterY - int(dy)
				x2 := cellCenterX + int(dx)
				y2 := cellCenterY + int(dy)
				// fmt.Println(cellCenterX, cellCenterY)
				diagonalLine(arrowImage, x1, y1, x2, y2, color.Gray{uint8(elm * 255)})
			}
		}
	}

	arrowFile, err := os.Create("./answer_69.jpg")
	defer arrowFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(arrowFile, arrowImage, &jpeg.Options{100})

}
