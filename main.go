package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func Plot(ps ...plot.Plotter) *image.RGBA {
	p := plot.New()
	p.Add(append([]plot.Plotter{
		plotter.NewGrid(),
	}, ps...)...)
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	c := vgimg.NewWith(vgimg.UseImage(img))
	p.Draw(draw.New(c))
	return c.Image().(*image.RGBA)
}

func main() {
	file, err := os.Open("data/house_prices.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	var inputs, labels []float64
	types := make([][]float64, len(data)-1)
	colors := make([][]float64, len(data)-1)
	for i, row := range data {
		if i == 0 {
			continue
		}
		for j, col := range row {
			if j == 0 {
				v, err := strconv.ParseFloat(col, 64)
				if err != nil {
					log.Fatal(err)
				}
				inputs = append(inputs, v)
			} else if j == 1 {
				var v int
				switch col {
				case "Duplex":
					v = 0
				case "Detached":
					v = 1
				case "Semi-detached":
					v = 2
				case "Townhouse":
					v = 3
				case "Multi-family":
					v = 4
				}
				types[i-1] = make([]float64, 5)
				types[i-1][v] = 1
			} else if j == 2 {
				v, err := strconv.ParseFloat(col, 64)
				if err != nil {
					log.Fatal(err)
				}
				labels = append(labels, v)
			} else if j == 3 {
				var v int
				switch col {
				case "brown":
					v = 0
				case "yellow":
					v = 1
				case "white":
					v = 2
				case "blue":
					v = 3
				case "green":
					v = 4
				}
				colors[i-1] = make([]float64, 5)
				colors[i-1][v] = 1
			}
		}
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Gradient descent")

	const (
		epochs                           = 3000
		printEveryNthEpochs              = 100
		learningRateW                    = 0.2e-3
		learningRateB                    = 0.5e-1
		plotLoss                         = false
		startValueRange                  = 1
		inputPointsMinX, inputPointsMaxX = 0, 150
		useWallColors                    = true
	)
	if useWallColors {
		types = colors
	}
	xys := make([]plotter.XYs, 5)
	for i := 0; i < len(inputs); i++ {
		for j := 0; j < 5; j++ {
			if types[i][j] == 1 {
				xys[j] = append(xys[j], plotter.XY{X: inputs[i], Y: labels[i]})
			}
		}
	}
	var inputsScatter []*plotter.Scatter
	for i := 0; i < 5; i++ {
		tmp, _ := plotter.NewScatter(xys[i])
		inputsScatter = append(inputsScatter, tmp)
	}
	inputsScatter[0].Color = color.RGBA{255, 0, 0, 255}
	inputsScatter[1].Color = color.RGBA{0, 255, 0, 255}
	inputsScatter[2].Color = color.RGBA{0, 0, 255, 255}
	inputsScatter[3].Color = color.RGBA{255, 255, 0, 255}
	inputsScatter[4].Color = color.RGBA{255, 0, 255, 255}

	img := make(chan *image.RGBA, 1) // Have at most one image in the channel.
	render := func(x *image.RGBA) {
		select {
		case <-img: // Drain the channel.
			img <- x // Put the new image in.
		case img <- x: // Or just put the new image in.
		}
	}
	go func() {
		w := make([]float64, 12)
		for i := range w {
			w[i] = startValueRange - rand.Float64()*2*startValueRange
		}
		var loss plotter.XYs
		for i := 0; i < epochs; i++ {
			y := inference(inputs, w, types)
			loss = append(loss, plotter.XY{
				X: float64(i),
				Y: msl(labels, y),
			})
			lossLines, _ := plotter.NewLine(loss)
			if plotLoss {
				render(Plot(lossLines))
			} else {
				var resLines []*plotter.Line
				const extra = (inputPointsMaxX - inputPointsMinX) / 10
				xs := []float64{inputPointsMinX - extra, inputPointsMaxX + extra}
				for i := 0; i < 5; i++ {
					houseTypes := make([][]float64, 2)
					houseTypes[0], houseTypes[1] = make([]float64, 5), make([]float64, 5)
					houseTypes[0][i], houseTypes[1][i] = 1, 1
					ys := inference(xs, w, houseTypes)
					resLine, _ := plotter.NewLine(plotter.XYs{{X: xs[0], Y: ys[0]}, {X: xs[1], Y: ys[1]}})
					resLines = append(resLines, resLine)
				}
				resLines[0].LineStyle.Color = color.RGBA{255, 0, 0, 255}
				resLines[1].LineStyle.Color = color.RGBA{0, 255, 0, 255}
				resLines[2].LineStyle.Color = color.RGBA{0, 0, 255, 255}
				resLines[3].LineStyle.Color = color.RGBA{255, 255, 0, 255}
				resLines[4].LineStyle.Color = color.RGBA{255, 0, 255, 255}
				render(Plot(inputsScatter[0], inputsScatter[1], inputsScatter[2], inputsScatter[3], inputsScatter[4], resLines[0], resLines[1], resLines[2], resLines[3], resLines[4]))
			}
			var dw, db []float64
			for i := 0; i < 5; i++ {
				tmpW, tmpB := dmslT(inputs, labels, y, types, i)
				dw = append(dw, tmpW)
				db = append(db, tmpB)
			}
			for i := 0; i < 5; i++ {
				w[i] += db[i] * learningRateB
				w[5] += db[i] * learningRateB
				w[i+6] += dw[i] * learningRateW
				w[11] += dw[i] * learningRateW
			}
			if i%printEveryNthEpochs == 0 {
				fmt.Printf(`Epoch #%d
				loss: %.4f
				dw: %.4f, db: %.4f
				w : %.4f, b: %.4f
				`, i, loss[len(loss)-1].Y, dw, db, w[:6], w[6:])
			}
		}
		fmt.Println(w)
	}()

	if err := ebiten.RunGame(&App{Img: img}); err != nil {
		log.Fatal(err)
	}
}

func inference(inputs, w []float64, t [][]float64) (res []float64) {
	for i, x := range inputs {
		res = append(res, (w[0]*t[i][0]+w[1]*t[i][1]+w[2]*t[i][2]+w[3]*t[i][3]+w[4]*t[i][4]+w[5])+(w[6]*t[i][0]+w[7]*t[i][1]+w[8]*t[i][2]+w[9]*t[i][3]+w[10]*t[i][4]+w[11])*x)
	}
	return res
}

func msl(labels, y []float64) (loss float64) {
	for i := range labels {
		loss += (labels[i] - y[i]) * (labels[i] - y[i])
	}
	return loss / float64(len(labels))
}

func dmslT(inputs, labels, y []float64, types [][]float64, t int) (dw, db float64) {
	for i := range labels {
		if types[i][t] == 1 {
			diff := labels[i] - y[i]
			dw += inputs[i] * diff
			db += diff
		}
	}
	return 2 * dw / float64(len(labels)), 2 * db / float64(len(labels))
}
