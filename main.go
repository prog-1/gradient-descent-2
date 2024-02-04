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
	const (
		epochs                           = 10000
		learningRateW                    = 0.2e-3
		learningRateB                    = 0.5e-1
		plotLoss                         = false
		inputPointsMinX, inputPointsMaxX = 0, 150
		useWallColors                    = true
		startValueRange                  = 1
	)
	ebiten.SetWindowSize(640, 480)

	data, err := readDataFromCSV("data/house_prices.csv")
	if err != nil {
		log.Fatalf("%v", err)
	}

	var inputs, labels []float64
	types, colors := make([][]float64, len(data)), make([][]float64, len(data))

	typeInd := map[string]int{
		"Duplex":        0,
		"Detached":      1,
		"Semi-detached": 2,
		"Townhouse":     3,
		"Multi-family":  4,
	}
	colorInd := map[string]int{
		"brown":  0,
		"yellow": 1,
		"white":  2,
		"blue":   3,
		"green":  4,
	}

	for i, house := range data {
		labels = append(labels, house.Price)
		inputs = append(inputs, house.Square)

		typeIndex := typeInd[house.Type]
		types[i] = make([]float64, 5)
		types[i][typeIndex] = 1

		colorIndex := colorInd[house.WallColor]
		colors[i] = make([]float64, 5)
		colors[i][colorIndex] = 1
	}

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
	colorss := []color.RGBA{
		{0, 0, 0, 255},
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
		{255, 0, 255, 255},
		{0, 255, 255, 255},
	}

	for i := 0; i < 5; i++ {
		tmp, _ := plotter.NewScatter(xys[i])
		inputsScatter = append(inputsScatter, tmp)
		inputsScatter[i].Color = colorss[i]
	}

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
		wk, wb := w[6:], w[:6]

		var loss plotter.XYs
		for i := 0; i < epochs; i++ {
			y := inference(inputs, wk, wb, types)
			loss = append(loss, plotter.XY{
				X: float64(i),
				Y: msl(labels, y),
			})
			lossLines, _ := plotter.NewLine(loss)
			var dw, db []float64
			for i := 0; i < 5; i++ {
				tmpW, tmpB := dmslT(inputs, labels, y, types, i)
				dw, db = append(dw, tmpW), append(db, tmpB)
			}
			for i := 0; i < 5; i++ {
				wk[i] += dw[i] * learningRateW
				wk[5] += dw[i] * learningRateW
				wb[i] += db[i] * learningRateB
				wb[5] += db[i] * learningRateB
			}
			if i%100 == 0 {
				if plotLoss {
					render(Plot(lossLines))
				} else {
					var lines []*plotter.Line
					const extra = (inputPointsMaxX - inputPointsMinX) / 10
					xs := []float64{inputPointsMinX - extra, inputPointsMaxX + extra}
					for i := 0; i < 5; i++ {
						makeHouseTypes := func(ind int) [][]float64 {
							houseTypes := make([][]float64, 2)
							for j := range houseTypes {
								houseTypes[j] = make([]float64, 5)
								houseTypes[j][ind] = 1
							}
							return houseTypes
						}
						houseTypes := makeHouseTypes(i)
						ys := inference(xs, wk, wb, houseTypes)
						resLine, _ := plotter.NewLine(plotter.XYs{{X: xs[0], Y: ys[0]}, {X: xs[1], Y: ys[1]}})
						lines = append(lines, resLine)
						lines[i].LineStyle.Color = colorss[i]
					}
					render(Plot(inputsScatter[0], inputsScatter[1], inputsScatter[2], inputsScatter[3], inputsScatter[4], lines[0], lines[1], lines[2], lines[3], lines[4]))
				}
				fmt.Printf(`Epoch #%d
				loss: %.4f
				dw: %.4f, db: %.4f
				w : %.4f, b: %.4f
				`, i, loss[len(loss)-1].Y, dw, db, wk, wb)
			}
		}
		fmt.Println(w)
	}()
	if err := ebiten.RunGame(&App{Img: img}); err != nil {
		log.Fatal(err)
	}
}

func inference(inputs []float64, wk, wb []float64, t [][]float64) (res []float64) {
	var y float64
	for i, x := range inputs {
		y = wk[5]*x + wb[5]
		for j := range t[i] {
			y += wk[j]*t[i][j]*x + wb[j]*t[i][j]
		}
		res = append(res, y)
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

type House struct {
	Square    float64
	Type      string
	Price     float64
	WallColor string
}

func readDataFromCSV(house_prices string) ([]House, error) {
	file, err := os.Open(house_prices)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 4
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var houses []House
	for i, record := range data {
		if i == 0 {
			continue
		}
		square, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatalf("Invalid square: %v", record)
			continue
		}

		price, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatalf("Invalid price: %v", record)
			continue
		}

		houses = append(houses, House{
			Square:    square,
			Type:      record[1],
			Price:     price,
			WallColor: record[3],
		})
	}
	return houses, nil
}
