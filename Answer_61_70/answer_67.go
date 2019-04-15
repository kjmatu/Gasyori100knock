package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

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
	// fmt.Println("angle")
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
			// fmt.Printf("%d ", int(angle[y][x]))
		}
		// fmt.Println()
	}

	for _, row := range angle {
		fmt.Println(row)
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
			// fmt.Printf("%f ", magnitude*255/magMax)
		}
		// fmt.Println()
	}

	// answer_67ここから
	N := 8
	HH := H / N
	HW := W / N
	hist := make([][][]float64, HH)
	for i := range hist {
		hist[i] = make([][]float64, HW)
		for j := range hist[i] {
			hist[i][j] = make([]float64, 9)
		}
	}
	fmt.Println(hist)
	fmt.Println(len(hist), len(hist[0]), len(hist[0][0]))

	for y := 0; y < HH; y++ {
		for x := 0; x < HW; x++ {
			for j := 0; j < N; j++ {
				for i := 0; i < N; i++ {
					hist[y][x][int(angle[y*4+j][x*4+i])] += mag[y*4+j][x*4+i]
				}
			}
		}
	}

	for _, row := range hist {
		for _, rowrow := range row {
			fmt.Println(rowrow)
		}
	}

	// 2D-Histgramの描画処理
	row, col := 3, 3
	plots := make([][]*plot.Plot, row)
	for j := range plots {
		plots[j] = make([]*plot.Plot, col)
	}

	i := 0
	for j := 0; j < row; j++ {
		for k := 0; k < col; k++ {
			histPlot := hplot.New()
			histPlot.Title.Text = fmt.Sprintf("HistIndex%d", i)
			// histPlot.X.Label.Text = "Cell IndexX"
			// histPlot.Y.Label.Text = "Cell IndexY"
			hist2D := hbook.NewH2D(HW, 0, float64(HW), HH, 0, float64(HH))
			for cellIndexY := 0; cellIndexY < HH; cellIndexY++ {
				for cellIndexX := 0; cellIndexX < HW; cellIndexX++ {
					histVal := hist[cellIndexY][cellIndexX][i]
					hist2D.Fill(float64(cellIndexX), float64(cellIndexY), histVal)
					histPlot.Add(hplot.NewH2D(hist2D, nil))
					histPlot.Add(plotter.NewGrid())
				}
			}
			plots[j][k] = histPlot.Plot
			i++
		}
	}

	img := vgimg.New(10*vg.Centimeter, 10*vg.Centimeter)
	dc := draw.New(img)
	t := draw.Tiles{
		Rows:      row,
		Cols:      col,
		PadX:      vg.Millimeter,
		PadY:      vg.Millimeter,
		PadTop:    vg.Points(2),
		PadBottom: vg.Points(2),
		PadLeft:   vg.Points(2),
		PadRight:  vg.Points(2),
	}
	canvases := plot.Align(plots, t, dc)
	for j := 0; j < row; j++ {
		for k := 0; k < col; k++ {
			plots[j][k].Draw(canvases[j][k])
		}
	}

	f, err := os.Create("answer_67.png")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(f); err != nil {
		panic(err)
	}

}
