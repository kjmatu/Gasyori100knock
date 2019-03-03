package main

import (
	"image/jpeg"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	file, err := os.Open("./../Question_11_20/imori_dark.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	v := make(plotter.Values, jimg.Bounds().Size().Y*jimg.Bounds().Size().X*3)

	// RGB画像を配列に格納する
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

	// plotを作成
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	// 上限下限タイトルを設定
	p.Title.Text = "imori_dark.jpg histogram"
	p.X.Min = 0.0
	p.X.Max = 255.0
	p.Y.Min = 0.0
	p.Y.Max = 1400.0

	// ヒストグラムを作成
	h, err := plotter.NewHist(v, 255)
	if err != nil {
		panic(err)
	}
	h.Color = plotutil.Color(2)
	p.Add(h)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "answer_20.png"); err != nil {
		panic(err)
	}

}
