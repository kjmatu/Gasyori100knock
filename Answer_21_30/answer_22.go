package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func histOperation(pixVal, m, s, m0, s0 float64) uint8 {
	operation := s0*(pixVal-m)/s + m0
	return uint8(operation)
}

func main() {
	file, err := os.Open("./../Question_21_30/imori.jpg")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	jimg, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// RGB画像の平均画素値と標準偏差を求める
	m := 0.0
	for height := 0; height < jimg.Bounds().Size().Y; height++ {
		for width := 0; width < jimg.Bounds().Size().X; width++ {
			r, g, b, _ := jimg.At(width, height).RGBA()
			r8 := (float64(r*0xFF) / 0xFFFF)
			g8 := (float64(g*0xFF) / 0xFFFF)
			b8 := (float64(b*0xFF) / 0xFFFF)

			m += r8 + g8 + b8
		}
	}
	// fmt.Println("m sum", m)
	totalPixNum := jimg.Bounds().Size().Y * jimg.Bounds().Size().X * 3
	m /= float64(totalPixNum)

	s := 0.0
	for height := 0; height < jimg.Bounds().Size().Y; height++ {
		for width := 0; width < jimg.Bounds().Size().X; width++ {
			r, g, b, _ := jimg.At(width, height).RGBA()
			r8 := (float64(r*0xFF) / 0xFFFF)
			g8 := (float64(g*0xFF) / 0xFFFF)
			b8 := (float64(b*0xFF) / 0xFFFF)
			s += math.Pow((r8-m), 2) + math.Pow((g8-m), 2) + math.Pow((b8-m), 2)
		}
	}
	s /= float64(totalPixNum)
	s = math.Sqrt(s)

	histOperationImg := image.NewRGBA(jimg.Bounds())
	m0 := 128.0
	s0 := 52.0
	for height := 0; height < histOperationImg.Bounds().Size().Y; height++ {
		for width := 0; width < histOperationImg.Bounds().Size().X; width++ {
			r, g, b, _ := jimg.At(width, height).RGBA()
			r8 := float64(r*0xFF) / 0xFFFF
			g8 := float64(g*0xFF) / 0xFFFF
			b8 := float64(b*0xFF) / 0xFFFF

			var histOperationColor color.RGBA
			histOperationColor.R = histOperation(r8, m, s, m0, s0)
			histOperationColor.G = histOperation(g8, m, s, m0, s0)
			histOperationColor.B = histOperation(b8, m, s, m0, s0)
			histOperationColor.A = 255
			histOperationImg.Set(width, height, histOperationColor)
		}
	}

	histOperationImgFile, err := os.Create("./answer_22_1.jpg")
	defer histOperationImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(histOperationImgFile, histOperationImg, &jpeg.Options{100})

	// plotを作成
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// 上限下限タイトルを設定
	p.Title.Text = "operation imori.jpg histogram"
	p.X.Min = 0.0
	p.X.Max = 255.0
	p.Y.Min = 0.0
	p.Y.Max = 400.0

	// ヒストグラム操作した画像を1つの配列に格納する
	v := make(plotter.Values, histOperationImg.Bounds().Size().Y*histOperationImg.Bounds().Size().X*3)
	i := 0
	for height := 0; height < histOperationImg.Bounds().Size().Y; height++ {
		for width := 0; width < histOperationImg.Bounds().Size().X; width++ {
			r, g, b, _ := histOperationImg.At(width, height).RGBA()
			r8 := (float64(r*0xFF) / 0xFFFF)
			g8 := (float64(g*0xFF) / 0xFFFF)
			b8 := (float64(b*0xFF) / 0xFFFF)
			v[i] = r8
			v[i+1] = g8
			v[i+2] = b8
			i += 3
		}
	}

	h, err := plotter.NewHist(v, 255)
	if err != nil {
		panic(err)
	}
	h.Color = plotutil.Color(2)
	p.Add(h)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "answer_22_2.png"); err != nil {
		panic(err)
	}

}
