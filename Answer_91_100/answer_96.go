package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"math/rand"
	"os"
)

func uint32ToUint8(pixVal uint32) uint8 {
	return uint8(float64(pixVal*0xFF) / 0xFFFF)
}

func biLinearElm(pixVal00, pixVal01, pixVal10, pixVal11 uint32, dx, dy float64) uint8 {
	pixValUint8_00 := uint32ToUint8(pixVal00)
	pixValUint8_01 := uint32ToUint8(pixVal01)
	pixValUint8_10 := uint32ToUint8(pixVal10)
	pixValUint8_11 := uint32ToUint8(pixVal11)
	biLinearVal := (1.0-dx)*(1.0-dy)*float64(pixValUint8_00) +
		dx*(1-dy)*float64(pixValUint8_10) +
		(1.0-dx)*dy*float64(pixValUint8_01) + dx*dy*float64(pixValUint8_11)
	if biLinearVal > 255.0 {
		biLinearVal = 255.0
	}
	if biLinearVal < 0.0 {
		biLinearVal = 0
	}

	return uint8(biLinearVal)
}

func biLinear(I00, I01, I10, I11 color.Color, dx, dy float64) color.Color {
	var biLinearColor color.RGBA
	r00, g00, b00, _ := I00.RGBA()
	r01, g01, b01, _ := I01.RGBA()
	r10, g10, b10, _ := I10.RGBA()
	r11, g11, b11, _ := I11.RGBA()
	biLinearColor.R = biLinearElm(r00, r01, r10, r11, dx, dy)
	biLinearColor.G = biLinearElm(g00, g01, g10, g11, dx, dy)
	biLinearColor.B = biLinearElm(b00, b01, b10, b11, dx, dy)
	biLinearColor.A = 255
	return biLinearColor
}

func bilinearScale(scaleImage *image.Gray, scale float64) *image.Gray {
	// 拡大縮小画像を作成
	scaleBounds := scaleImage.Bounds()
	if scaleBounds.Min.X != 0 || scaleBounds.Min.Y != 0 {
		scaleBounds.Max.X -= scaleBounds.Min.X
		scaleBounds.Max.Y -= scaleBounds.Min.Y
		scaleBounds.Min.X = 0
		scaleBounds.Min.Y = 0
	}
	scaleBounds.Max.X = int(float64(scaleBounds.Max.X) * scale)
	scaleBounds.Max.Y = int(float64(scaleBounds.Max.Y) * scale)
	biLinearImg := image.NewGray(scaleBounds)

	for height := 0; height < biLinearImg.Bounds().Size().Y; height++ {
		for width := 0; width < biLinearImg.Bounds().Size().X; width++ {
			// 拡大画像のピクセル位置に対応する元画像位置を計算
			srcX := float64(width) / scale
			srcY := float64(height) / scale

			// 上記で計算した位置の周囲4点の画素値を取得
			I00 := scaleImage.At(int(srcX), int(srcY))
			I10 := scaleImage.At(int(srcX+1), int(srcY))
			I01 := scaleImage.At(int(srcX), int(srcY+1))
			I11 := scaleImage.At(int(srcX+1), int(srcY+1))

			// 周囲4点とピクセル対応点の距離を計算
			dx := srcX - math.Floor(srcX)
			dy := srcY - math.Floor(srcY)

			// Bi Linear補完
			biLinearColor := biLinear(I00, I01, I10, I11, dx, dy)
			biLinearImg.Set(width, height, biLinearColor)

			// 画像の境界を補完するときは最近傍補完を行う
			if width == biLinearImg.Bounds().Size().X-1 {
				biLinearImg.Set(width, height, scaleImage.At(int(srcX), int(srcY)))
			}
			if height == biLinearImg.Bounds().Size().Y-1 {
				biLinearImg.Set(width, height, scaleImage.At(int(srcX), int(srcY)))
			}
		}
	}
	return biLinearImg
}

// HLine draws a horizontal line
func hLine(img *image.RGBA, x1, y, x2 int, col color.RGBA) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// VLine draws a veritcal line
func vLine(img *image.RGBA, x, y1, y2 int, col color.RGBA) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func drawRectangle(img *image.RGBA, rect image.Rectangle, col color.RGBA) {
	x1 := rect.Min.X
	y1 := rect.Min.Y
	x2 := rect.Max.X
	y2 := rect.Max.Y
	hLine(img, x1, y1, x2, col)
	hLine(img, x1, y2, x2, col)
	vLine(img, x1, y1, y2, col)
	vLine(img, x2, y1, y2, col)
}

func calcIoU(rect1 image.Rectangle, rect2 image.Rectangle) float64 {
	area1 := rect1.Dx() * rect1.Dy()
	area2 := rect2.Dx() * rect2.Dy()

	overlapRect := rect1.Intersect(rect2)
	overlapArea := overlapRect.Dx() * overlapRect.Dy()

	iou := float64(overlapArea) / float64(area1+area2-overlapArea)

	return iou
}

func randomClip(srcImage *image.RGBA) (image.Image, image.Rectangle) {
	width, height := 60, 60
	H := srcImage.Bounds().Size().Y
	W := srcImage.Bounds().Size().X
	randomX := int(rand.Int63n(int64(W - 60)))
	randomY := int(rand.Int63n(int64(H - 60)))

	clipRect := image.Rect(randomX, randomY, randomX+width, randomY+height)
	clipImage := srcImage.SubImage(clipRect)
	return clipImage, clipRect
}

func hog(srcImage *image.Gray) [][][]float64 {
	H := srcImage.Bounds().Size().Y
	W := srcImage.Bounds().Size().X

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
			gx[y][x] = float64(srcImage.GrayAt(rightIndex, y).Y - srcImage.GrayAt(leftIndex, y).Y)

			upIndex := int(math.Max(float64(y-1), 0))
			downIndex := int(math.Min(float64(y+1), float64(H-1)))
			gy[y][x] = float64(srcImage.GrayAt(x, downIndex).Y - srcImage.GrayAt(x, upIndex).Y)
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

	// 3x3のセルを1ブロックとして扱う
	epsilon := 1.0

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

			for i := range hist[cellY][cellX] {
				hist[cellY][cellX][i] /= sum
			}
		}
	}
	return hist
}

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

type NeuralNetwork struct {
	W1   [][]float64
	B1   []float64
	Wout [][]float64
	Bout []float64
	Lr   float64
	Z1   [][]float64
	Z2   [][]float64
	Out  [][]float64
}

func (nn *NeuralNetwork) dot(matrix1 [][]float64, matrix2 [][]float64) ([][]float64, error) {
	height1 := len(matrix1)
	width1 := len(matrix1[0])

	height2 := len(matrix2)
	width2 := len(matrix2[0])

	if width1 != height2 {
		return nil, errors.New("input matrix does not match size")
	}

	dotArray := make([][]float64, height1)
	for i := range dotArray {
		dotArray[i] = make([]float64, width2)
	}

	for y, row := range dotArray {
		for x := range row {
			dotValue := 0.0
			for i := 0; i < len(matrix1[y][:]); i++ {
				dotValue += matrix1[y][i] * matrix2[i][x]
			}
			dotArray[y][x] = dotValue
		}
	}
	return dotArray, nil
}

func (nn *NeuralNetwork) transpose(matrix [][]float64) [][]float64 {

	transposeMatrix := make([][]float64, len(matrix[0]))
	for i := range transposeMatrix {
		transposeMatrix[i] = make([]float64, len(matrix))
	}

	for i, row := range transposeMatrix {
		for j := range row {
			transposeMatrix[i][j] = matrix[j][i]
		}
	}
	return transposeMatrix
}

func (nn *NeuralNetwork) Forward(x [][]float64) [][]float64 {
	nn.Z1 = x

	// Z2
	dotValue, err := nn.dot(nn.Z1, nn.W1)
	if err != nil {
		fmt.Println("Dot error")
	}

	for i, row := range dotValue {
		for j := range row {
			dotValue[i][j] += nn.B1[j]
		}
	}

	nn.Z2 = make([][]float64, len(dotValue))
	for i := range nn.Z2 {
		nn.Z2[i] = make([]float64, len(dotValue[0]))
	}

	for i, row := range nn.Z2 {
		for j := range row {
			nn.Z2[i][j] = sigmoid(dotValue[i][j])
		}
	}

	// Out
	dotValue, err = nn.dot(nn.Z2, nn.Wout)
	if err != nil {
		fmt.Println("Dot error")
	}

	for i, row := range dotValue {
		for j := range row {
			dotValue[i][j] += nn.Bout[j]
		}
	}

	nn.Out = make([][]float64, len(dotValue))
	for i := range nn.Out {
		nn.Out[i] = make([]float64, len(dotValue[0]))
	}

	for i, row := range nn.Out {
		for j := range row {
			nn.Out[i][j] = sigmoid(dotValue[i][j])
		}
	}
	return nn.Out
}

func (nn *NeuralNetwork) Train(x, t [][]float64) {
	// backpropagation output layer
	En := make([][]float64, len(nn.Out))
	for i := range En {
		En[i] = make([]float64, len(nn.Out[0]))
	}

	for i, row := range nn.Out {
		for j := range row {
			En[i][j] = (nn.Out[i][j] - t[i][j]) * nn.Out[i][j] * (1 - nn.Out[i][j])
		}
	}

	// gradEn := append([][]float64{}, En...)

	gradWout, err := nn.dot(nn.transpose(nn.Z2), En)
	if err != nil {
		fmt.Println("Dot error")
	}

	ones := make([][]float64, len(En[0]))
	for i := range ones {
		ones[i] = make([]float64, len(En))
	}

	for i, row := range ones {
		for j := range row {
			ones[i][j] = 1.0
		}
	}

	gradBout, err := nn.dot(ones, En)
	if err != nil {
		fmt.Println("Dot error")
	}

	for i, row := range nn.Wout {
		for j := range row {
			nn.Wout[i][j] -= nn.Lr * gradWout[i][j]
		}
	}

	nn.Bout[0] -= nn.Lr * gradBout[0][0]

	// backpropagation inter layer
	// gradU1
	gradU1, err := nn.dot(En, nn.transpose(nn.Wout))
	for i, row := range gradU1 {
		for j := range row {
			gradU1[i][j] *= nn.Z2[i][j] * (1 - nn.Z2[i][j])
		}
	}

	// gradW1
	gradW1, err := nn.dot(nn.transpose(nn.Z1), gradU1)
	ones = make([][]float64, 1)
	for i := range ones {
		ones[i] = make([]float64, len(gradU1))
	}

	for i, row := range ones {
		for j := range row {
			ones[i][j] = 1.0
		}
	}

	// gradB1
	gradB1, err := nn.dot(ones, gradU1)

	for i, row := range nn.W1 {
		for j := range row {
			nn.W1[i][j] -= nn.Lr * gradW1[i][j]
		}
	}

	for i := range nn.B1 {
		nn.B1[i] -= nn.Lr * gradB1[0][i]
	}

}

func NewNeuralNetwork(index, weight, outd int, lr float64) *NeuralNetwork {
	rand.Seed(0)
	nn := new(NeuralNetwork)

	nn.W1 = make([][]float64, index)
	for i := range nn.W1 {
		nn.W1[i] = make([]float64, weight)
	}

	for i, row := range nn.W1 {
		for j := range row {
			nn.W1[i][j] = rand.NormFloat64()
		}
	}

	nn.B1 = make([]float64, weight)
	for i := range nn.B1 {
		nn.B1[i] = rand.NormFloat64()
	}

	nn.Wout = make([][]float64, weight)
	for i := range nn.Wout {
		nn.Wout[i] = make([]float64, outd)
	}

	for i, row := range nn.Wout {
		for j := range row {
			nn.Wout[i][j] = rand.NormFloat64()
		}
	}

	nn.Bout = make([]float64, outd)
	for i := range nn.Bout {
		nn.Bout[i] = rand.NormFloat64()
	}

	nn.Lr = lr

	return nn
}

func main() {
	file, err := os.Open("./../Question_91_100/imori_1.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpegImage, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	imoriImage := image.NewRGBA(jpegImage.Bounds())
	imoriFaceRect := image.Rect(47, 41, 129, 103)
	draw.Draw(imoriImage, jpegImage.Bounds(), jpegImage, jpegImage.Bounds().Min, draw.Src)

	rand.Seed(0)
	scaleSize := 32
	db := make([][]float64, 200)
	for i := range db {
		db[i] = make([]float64, (scaleSize/8)*(scaleSize/8)*9+1)
	}

	for i := 0; i < 200; i++ {
		var label float64
		clipImage, clipRect := randomClip(imoriImage)
		iou := calcIoU(imoriFaceRect, clipRect)
		if iou > 0.5 {
			label = 1.0
		} else {
			label = 0.0
		}
		// fmt.Println(label)

		clipGrayImage := image.NewGray(clipImage.Bounds())
		H := clipGrayImage.Bounds().Size().Y
		W := clipGrayImage.Bounds().Size().X
		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				c := color.GrayModel.Convert(clipImage.At(x, y))
				gray, _ := c.(color.Gray)
				clipGrayImage.SetGray(x, y, gray)
			}
		}

		resizeClipImage := bilinearScale(clipGrayImage, 32.0/60.0)

		hogData := hog(resizeClipImage)

		hogDataFlatten := []float64{}
		for _, rowY := range hogData {
			for _, hog := range rowY {
				hogDataFlatten = append(hogDataFlatten, hog...)
			}
		}
		db[i] = hogDataFlatten
		db[i] = append(db[i], label)
	}

	// train neural network
	nn := NewNeuralNetwork(4*4*9, 64, 1, 0.01)
	for i := 0; i < 100000; i++ {
		nn.Forward(db[:][:4*4*9])
		// nn.Train(db[:][:4*4*9], db[:][4*4*9+1])
	}
}
