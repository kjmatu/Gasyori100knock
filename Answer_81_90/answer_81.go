package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func sobelFileter(grayArray [][]float64, sobelMatrix [3][3]float64) ([][]float64, *image.Gray) {
	SH := len(sobelMatrix)
	SW := len(sobelMatrix[0][:])

	H := len(grayArray)
	W := len(grayArray[0][:])

	sobelImage := image.NewGray(image.Rect(0, 0, W, H))
	sobelArray := make([][]float64, H)
	for x := range sobelArray {
		sobelArray[x] = make([]float64, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			filterValue := 0.0
			for sy := 0; sy < SW; sy++ {
				for sx := 0; sx < SH; sx++ {
					refY := y + sy - SH/2
					refX := x + sx - SW/2
					gaussianElm := sobelMatrix[sy][sx]
					if refX < 0 && refY < 0 || refX >= W && refY >= H || refX < 0 && refY >= H || refX >= W && refY < 0 {
						filterValue += gaussianElm * grayArray[y][x]
					} else if refY < 0 || refY >= H {
						filterValue += gaussianElm * grayArray[y][refX]
					} else if refX < 0 || refX >= W {
						filterValue += gaussianElm * grayArray[refY][x]
					} else {
						filterValue += gaussianElm * grayArray[refY][refX]
					}
				}
			}
			sobelArray[y][x] = filterValue / float64(SW*SH)

			// 画像保存用
			if filterValue > 255.0 {
				filterValue = 255.0
			}
			if filterValue < 0.0 {
				filterValue = 0.0
			}
			sobelImage.Set(x, y, color.Gray{uint8(filterValue / float64(SW*SH))})
		}
	}
	return sobelArray, sobelImage
}

func determinant(matrix [2][2]float64) float64 {
	// fmt.Println(matrix[0][0] * matrix[1][1])
	// fmt.Println(matrix[0][1] * matrix[1][0])
	return matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]
}

func calcGaussCurvature(hessianMatrix [2][2]float64, Ix, Iy float64) float64 {
	k := determinant(hessianMatrix)
	// fmt.Println("det(H)", k)
	k /= math.Pow((1 + Ix*Ix + Iy*Iy), 2)
	return k
}

func main() {
	file, err := os.Open("./../Question_81_90/thorino.jpg")
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

	// カラー画像をグレイスケール画像に変換
	grayArray := make([][]float64, H)
	// gaussCurvature := make([][]float64, H)
	dethArray := make([][]float64, H)
	for i := range grayArray {
		grayArray[i] = make([]float64, W)
		// gaussCurvature[i] = make([]float64, W)
		dethArray[i] = make([]float64, W)
	}

	grayImage := image.NewGray(jpegImage.Bounds())
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			c := color.GrayModel.Convert(jpegImage.At(x, y))
			gray, _ := c.(color.Gray)
			grayImage.SetGray(x, y, gray)
			grayArray[y][x] = float64(gray.Y)
		}
	}

	// 縦方向Sobelフィルタを作成
	sobelFilterV := [3][3]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1}}

	// Y軸方向の微分画像を作成
	Iy, _ := sobelFileter(grayArray, sobelFilterV) // 1次
	Iyy, _ := sobelFileter(Iy, sobelFilterV)       // 2次

	// 横方向Sobelフィルタを作成
	sobelFilterH := [3][3]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1}}

	// X軸方向の微分画像を作成
	Ix, _ := sobelFileter(grayArray, sobelFilterH) // 1次
	Ixx, _ := sobelFileter(Ix, sobelFilterH)       // 2次

	// Y軸X軸それぞれの微分画像を作成
	Iyx, _ := sobelFileter(Iy, sobelFilterH) //Y->X
	Ixy, _ := sobelFileter(Ix, sobelFilterV) // X->Y

	maxDetH := math.Inf(-1)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {

			// 2回微分を使ったHessian Matrix
			hessianMatrix := [2][2]float64{
				{Ixx[y][x], Ixy[y][x]},
				{Iyx[y][x], Iyy[y][x]},
			}

			// テイラー展開で近似したHessian Matrix
			hessianMatrix[0][0] = Ix[y][x] * Ix[y][x]
			hessianMatrix[1][1] = Iy[y][x] * Iy[y][x]
			hessianMatrix[0][1] = Ix[y][x] * Iy[y][x]
			hessianMatrix[1][0] = Ix[y][x] * Iy[y][x]

			// gaussCurvature[y][x] = calcGaussCurvature(hessianMatrix, Ix[y][x], Iy[y][x])

			dethArray[y][x] = determinant(hessianMatrix)
			if dethArray[y][x] > maxDetH {
				maxDetH = dethArray[y][x]
			}
		}
	}

	cornerPoint := []image.Point{}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			// 周囲8近傍のDeterminant値を取得する
			leftIndex := int(math.Max(float64(x-1), 0))
			rightIndex := int(math.Min(float64(x+1), float64(W)-1))
			upIndex := int(math.Max(float64(y-1), 0))
			downIndex := int(math.Min(float64(y+1), float64(H)-1))
			det0 := dethArray[y][x]
			det1 := dethArray[y][rightIndex]
			det2 := dethArray[upIndex][rightIndex]
			det3 := dethArray[upIndex][x]
			det4 := dethArray[upIndex][leftIndex]
			det5 := dethArray[y][leftIndex]
			det6 := dethArray[downIndex][leftIndex]
			det7 := dethArray[downIndex][x]
			det8 := dethArray[downIndex][rightIndex]

			// 周囲8近傍と中央のDeterminant値を比較して中央が最大だったら極大点とする
			kArray := []float64{det1, det2, det3, det4, det5, det6, det7, det8}
			maximumPointFlag := true
			for _, det := range kArray {
				if det > det0 {
					maximumPointFlag = false
				}
			}

			if maximumPointFlag {
				if det0 >= (maxDetH * 0.1) {
					cornerPoint = append(cornerPoint, image.Point{x, y})
				}
			}
		}
	}
	// fmt.Println(cornerPoint)

	cornerImage := image.NewNRGBA(grayImage.Bounds())
	draw.Draw(cornerImage, grayImage.Bounds(), grayImage, grayImage.Bounds().Min, draw.Src)
	for _, point := range cornerPoint {
		cornerImage.SetNRGBA(point.X, point.Y, color.NRGBA{0xFF, 0, 0, 0xFF})
	}
	cornerFile, err := os.Create("./answer_81.jpg")
	defer cornerFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(cornerFile, cornerImage, &jpeg.Options{100})

}
