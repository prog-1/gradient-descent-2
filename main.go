package main

import (
	"encoding/csv"
	"fmt"
	"image"
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
	file, err := os.Open("C:/Programming/gradient-descent-2/data/house_prices.csv")
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
			}
		}
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Gradient descent")

	const (
		epochs                           = 2000
		printEveryNthEpochs              = 100
		learningRateW                    = 0.5e-4
		learningRateB                    = 0.7
		plotLoss                         = false
		startValueRange                  = 1
		inputPointsMinX, inputPointsMaxX = 0, 150
	)
	var xys plotter.XYs
	for i := 0; i < len(inputs); i++ {
		xys = append(xys, plotter.XY{X: inputs[i], Y: labels[i]})
	}
	inputsScatter, _ := plotter.NewScatter(xys)

	img := make(chan *image.RGBA, 1) // Have at most one image in the channel.
	render := func(x *image.RGBA) {
		select {
		case <-img: // Drain the channel.
			img <- x // Put the new image in.
		case img <- x: // Or just put the new image in.
		}
	}
	go func() {
		w := make([]float64, 6)
		for i := range w {
			w[i] = startValueRange - rand.Float64()*2*startValueRange
		}
		b := startValueRange - rand.Float64()*2*startValueRange
		var loss plotter.XYs
		for i := 0; i < epochs; i++ {
			y := inference(inputs, w, types, b)
			loss = append(loss, plotter.XY{
				X: float64(i),
				Y: msl(labels, y),
			})
			lossLines, _ := plotter.NewLine(loss)
			if plotLoss {
				render(Plot(lossLines))
			} else {
				const extra = (inputPointsMaxX - inputPointsMinX) / 10
				xs := []float64{inputPointsMinX - extra, inputPointsMaxX + extra}
				ys := inference(xs, w, types, b)
				resLine, _ := plotter.NewLine(plotter.XYs{{X: xs[0], Y: ys[0]}, {X: xs[1], Y: ys[1]}})
				render(Plot(inputsScatter, resLine))
			}
			dw, db := dmsl(inputs, labels, y)
			for i := range w {
				w[i] += dw * learningRateW
			}
			b += db * learningRateB
			if i%printEveryNthEpochs == 0 {
				fmt.Printf(`Epoch #%d
				loss: %.4f
				dw: %.4f, db: %.4f
				w : %.4f,  b: %.4f
				`, i, loss[len(loss)-1].Y, dw, db, w, b)
			}
		}
		fmt.Println(w, b)
	}()

	if err := ebiten.RunGame(&App{Img: img}); err != nil {
		log.Fatal(err)
	}
}

func inference(inputs, w []float64, t [][]float64, b float64) (res []float64) {
	for i, x := range inputs {
		res = append(res, w[0]*t[i][0]+w[1]*t[i][1]+w[2]*t[i][2]+w[3]*t[i][3]+w[4]*t[i][4]+w[5]*x+b)
	}
	return res
}

func msl(labels, y []float64) (loss float64) {
	for i := range labels {
		loss += (labels[i] - y[i]) * (labels[i] - y[i])
	}
	return loss / float64(len(labels))
}

func dmsl(inputs, labels, y []float64) (dw, db float64) {
	for i := range labels {
		diff := labels[i] - y[i]
		dw += inputs[i] * diff
		db += diff
	}
	return 2 * dw / float64(len(labels)), 2 * db / float64(len(labels))
}
