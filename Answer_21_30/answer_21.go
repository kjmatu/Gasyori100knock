package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func normalizeValue(pixVal, min, max, normMin, normMax float64) float64 {
	normValue := 0.0
	if pixVal <= min {
		normValue = normMin
	} else if pixVal >= max {
		normValue = normMax
	} else {
		normValue = ((normMax-normMin)*(pixVal-min))/(max-min) + normMin
	}
	return normValue
}

func main() {
	file, err := os.Open("./../Question_21_30/imori_dark.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	v := make(plotter.Values, jimg.Bounds().Size().Y*jimg.Bounds().Size().X*3)

	// RGB画像を配列に格納
	i := 0
	for height := 0; height < jimg.Bounds().Size().Y; height++ {
		for width := 0; width < jimg.Bounds().Size().X; width++ {
			r, g, b, _ := jimg.At(width, height).RGBA()
			v[i] = float64((r * 0xFF / 0xFFFF))
			v[i+1] = float64((g * 0xFF / 0xFFFF))
			v[i+2] = float64((b * 0xFF / 0xFFFF))
			i += 3
		}
	}

	// RGB画像の最大最小値を求める
	min := 255.0
	max := 0.0
	for _, value := range v {
		if min > value {
			min = value
		}

		if max < value {
			max = value
		}
	}

	// ヒストグラム正規化画像を作成
	normImg := image.NewRGBA(jimg.Bounds())
	for height := 0; height < normImg.Bounds().Size().Y; height++ {
		for width := 0; width < normImg.Bounds().Size().X; width++ {
			r, g, b, _ := jimg.At(width, height).RGBA()
			rfloat64 := float64((r * 0xFF / 0xFFFF))
			gfloat64 := float64((g * 0xFF / 0xFFFF))
			bfloat64 := float64((b * 0xFF / 0xFFFF))

			var normColor color.RGBA
			normColor.R = uint8(normalizeValue(rfloat64, min, max, 0, 255))
			normColor.G = uint8(normalizeValue(gfloat64, min, max, 0, 255))
			normColor.B = uint8(normalizeValue(bfloat64, min, max, 0, 255))
			normColor.A = uint8(255)
			normImg.Set(width, height, normColor)
		}
	}

	normImgFile, err := os.Create("./answer_21_1.jpg")
	defer normImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(normImgFile, normImg, &jpeg.Options{100})

	// plotを作成
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	// 上限下限タイトルを設定
	p.Title.Text = "imori_dark.jpg histogram normalization"
	p.X.Min = 0.0
	p.X.Max = 255.0
	p.Y.Min = 0.0
	p.Y.Max = 1400.0

	// 正規化されたヒストグラムを格納する配列を宣言
	vNorm := make(plotter.Values, normImg.Bounds().Size().Y*normImg.Bounds().Size().X*3)

	// ヒストグラム正規化画像を配列に格納
	i = 0
	for height := 0; height < normImg.Bounds().Size().Y; height++ {
		for width := 0; width < normImg.Bounds().Size().X; width++ {
			r, g, b, _ := normImg.At(width, height).RGBA()
			vNorm[i] = float64((r * 0xFF / 0xFFFF))
			vNorm[i+1] = float64((g * 0xFF / 0xFFFF))
			vNorm[i+2] = float64((b * 0xFF / 0xFFFF))
			i += 3
		}
	}

	// ヒストグラムを作成
	h, err := plotter.NewHist(vNorm, 255)
	if err != nil {
		panic(err)
	}
	h.Color = plotutil.Color(2)
	p.Add(h)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "answer_21_2.png"); err != nil {
		panic(err)
	}
}
