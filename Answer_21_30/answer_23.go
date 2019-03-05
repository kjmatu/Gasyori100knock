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

func histFlatten(inputPixVal float64, imgSize int, pixValMax float64, histgram [256]int) uint8 {
	sumOfHistgram := 0
	pixValInt := int(inputPixVal)
	for index := 0; index < pixValInt; index++ {
		sumOfHistgram += histgram[index]
	}

	flattenPixVal := float64(pixValMax*float64(sumOfHistgram)) / float64(imgSize)
	return uint8(flattenPixVal)
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

	// RGB画像のヒストグラムと画素値の最大値を求める
	hist := [256]int{}
	i := 0
	max := 0.0
	for height := 0; height < jimg.Bounds().Size().Y; height++ {
		for width := 0; width < jimg.Bounds().Size().X; width++ {
			r, g, b, _ := jimg.At(width, height).RGBA()
			r8 := (float64(r*0xFF) / 0xFFFF)
			g8 := (float64(g*0xFF) / 0xFFFF)
			b8 := (float64(b*0xFF) / 0xFFFF)
			if max < r8 {
				max = r8
			}
			if max < g8 {
				max = g8
			}
			if max < b8 {
				max = b8
			}
			i += 3
			hist[int(r8)]++
			hist[int(g8)]++
			hist[int(b8)]++
		}
	}

	flattenImg := image.NewRGBA(jimg.Bounds())
	for height := 0; height < flattenImg.Bounds().Size().Y; height++ {
		for width := 0; width < flattenImg.Bounds().Size().X; width++ {
			r, g, b, _ := jimg.At(width, height).RGBA()
			r8 := float64(r*0xFF) / 0xFFFF
			g8 := float64(g*0xFF) / 0xFFFF
			b8 := float64(b*0xFF) / 0xFFFF

			var flattenColor color.RGBA
			flattenColor.R = histFlatten(r8, i, max, hist)
			flattenColor.G = histFlatten(g8, i, max, hist)
			flattenColor.B = histFlatten(b8, i, max, hist)
			flattenColor.A = 255
			flattenImg.Set(width, height, flattenColor)
		}
	}

	flattenImgFile, err := os.Create("./answer_23_1.jpg")
	defer flattenImgFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jpeg.Encode(flattenImgFile, flattenImg, &jpeg.Options{100})

	// plotを作成
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// 上限下限タイトルを設定
	p.Title.Text = "flatten imori.jpg histogram"
	p.X.Min = 0.0
	p.X.Max = 255.0
	p.Y.Min = 0.0
	p.Y.Max = 400.0

	// 平坦化した画像を1つの配列に格納する
	v := make(plotter.Values, jimg.Bounds().Size().Y*jimg.Bounds().Size().X*3)
	i = 0
	for height := 0; height < flattenImg.Bounds().Size().Y; height++ {
		for width := 0; width < flattenImg.Bounds().Size().X; width++ {
			r, g, b, _ := flattenImg.At(width, height).RGBA()
			r8 := (float64(r*0xFF) / 0xFFFF)
			g8 := (float64(g*0xFF) / 0xFFFF)
			b8 := (float64(b*0xFF) / 0xFFFF)
			if max < r8 {
				max = r8
			}
			if max < g8 {
				max = g8
			}
			if max < b8 {
				max = b8
			}
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
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "answer_23_2.png"); err != nil {
		panic(err)
	}

}
