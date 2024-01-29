package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type House struct {
	Square    float64
	HouseType string
	Price     float64
	WallColor string
}

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
	csvFile, err := os.Open("data/house_prices.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer csvFile.Close()

	houses, err := readHousesFromCSV(csvFile)
	if err != nil {
		fmt.Println("Error reading houses from CSV:", err)
		return
	}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Gradient descent")

	const (
		epochs              = 10000
		printEveryNthEpochs = 500
		learningRateW       = 0.1e-4
		learningRateB       = 0.7

		plotLoss = false // Loss curve: true, Resulting line: false.

		inputPointsMinX, inputPointsMaxX = 0, 125
		startValueRange                  = 1 // Start values for weights are in range [-startValueRange, startValueRange].
	)

	var (
		inputs, labels []float64
	)
	types := make([][]float64, len(houses))

	xys := make([]plotter.XYs, 5)
	for i := 0; i < len(inputs); i++ {
		for j := 0; j < 5; j++ {
			if types[i][j] == 1 {
				xys[j] = append(xys[j], plotter.XY{X: inputs[i], Y: labels[i]})
			}
		}
	}
	for i, house := range houses {
		inputs = append(inputs, house.Square)
		labels = append(labels, house.Price)
		switch house.HouseType {
		case "Duplex":
			types[i] = append(types[i], []float64{1, 0, 0, 0, 0}...)
		case "Detached":
			types[i] = append(types[i], []float64{0, 1, 0, 0, 0}...)
		case "Semi-detached":
			types[i] = append(types[i], []float64{0, 0, 1, 0, 0}...)
		case "Townhouse":
			types[i] = append(types[i], []float64{0, 0, 0, 1, 0}...)
		case "Multi-family":
			types[i] = append(types[i], []float64{0, 0, 0, 0, 1}...)
		}
		for j := 0; j < 5; j++ {
			if types[i][j] == 1 {
				xys[j] = append(xys[j], plotter.XY{X: inputs[i], Y: labels[i]})
			}
		}
	}

	img := make(chan *image.RGBA, 1) // Have at most one image in the channel.
	render := func(x *image.RGBA) {
		select {
		case <-img: // Drain the channel.
			img <- x // Put the new image in.
		case img <- x: // Or just put the new image in.
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

	go func() {
		var loss plotter.XYs
		w := make([]float64, 10)
		for i := 0; i < 10; i++ {
			w[i] = startValueRange - rand.Float64()*2*startValueRange
		}

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

			for i := 0; i < 5; i++ {
				dw, db := dmslTypes(inputs, labels, y, types, i)
				w[i] += dw * learningRateW
				w[i+5] += db * learningRateB
			}

			if i%printEveryNthEpochs == 0 {
				fmt.Printf(`Epoch #%d
	loss: %.4f
	w : %.4f, 
`, i, loss[len(loss)-1].Y, w)
			}
		}
		fmt.Println(w)
	}()

	if err := ebiten.RunGame(&App{Img: img}); err != nil {
		log.Fatal(err)
	}
}

func inference(inputs, w []float64, types [][]float64) (res []float64) {
	for i, x := range inputs {
		res = append(res, w[0]*x*types[i][0]+w[1]*x*types[i][1]+w[2]*x*types[i][2]+w[3]*x*types[i][3]+w[4]*x*types[i][4]+w[5]*types[i][0]+w[6]*types[i][1]+w[7]*types[i][2]+w[8]*types[i][3]+w[9]*types[i][4])
	}
	return res
}

func msl(labels, y []float64) (loss float64) {
	for i := range labels {
		loss += (labels[i] - y[i]) * (labels[i] - y[i])
	}
	return loss / float64(len(labels))
}

func dmslTypes(inputs, labels, y []float64, types [][]float64, t int) (dw, db float64) {
	for i := range labels {
		if types[i][t] == 1 {
			diff := labels[i] - y[i]
			dw += inputs[i] * diff
			db += diff
		}
	}
	return 2 * dw / float64(len(labels)), 2 * db / float64(len(labels))
}

func readHousesFromCSV(csvFile io.Reader) ([]House, error) {
	houses := []House{}

	reader := csv.NewReader(csvFile)
	reader.Comma = ','

	for i := 0; ; {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if i == 0 {
			i++
			continue
		}

		square, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return nil, err
		}

		houseType := record[1]

		price, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}

		wallColor := record[3]

		house := House{
			Square:    square,
			HouseType: houseType,
			Price:     price,
			WallColor: wallColor,
		}

		houses = append(houses, house)
	}

	return houses, nil
}
