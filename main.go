package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/plot/plotter"
)

const (
	screenWidth, screenHeight       = 720, 480
	randMin, randMax                = 1500, 2000
	epochs, typeCount               = 1e6, 5
	lrw                             = 0.1e-2
	useWallColours, wallColourCount = true, typeCount
)

var typeColors [typeCount]color.RGBA = [typeCount]color.RGBA{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 255, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{100, 200, 0, 255}}

var learningRates = [typeCount * 2]float64{lrw, 0.1e-9, lrw, lrw, lrw, 0.8e-1, 0.9e-3, 0.5e-1, 0.9, 0.8e-1}

// Prediction(inference) for one argument
func p(x float64, t [typeCount]float64, w [typeCount * 2]float64) (y float64) {
	for i := range t {
		y += w[i]*t[i]*x + w[typeCount+i]*t[i]
	}

	return
}

func inference(xs []float64, types [][typeCount]float64, weights [typeCount * 2]float64) (ys []float64) {
	for i := range xs {
		ys = append(ys, p(xs[i], types[i], weights))
	}
	return
}

func gradient(labels, y, x []float64, types [][typeCount]float64) (ds [typeCount * 2]float64) {
	// ds - weight partial DerivativeS
	for i := 0; i < len(labels); i++ {
		for t := 0; t < typeCount; t++ {
			if types[i][t] == 1 {
				dif := y[i] - labels[i]
				ds[t] += dif * x[i]
				ds[t+typeCount] += dif
			}
		}
	}

	n := float64(len(labels))
	for i := range ds {
		ds[i] *= 2 / n
	}

	return
}

func main() {
	houses, err := readHousesFromCSV("data/house_prices.csv")
	if err != nil {
		log.Fatalf("Can't read Houses from CSV: %v", err)
	}

	var types [][typeCount]float64
	var squares []float64
	var labels []float64
	var wallColours [][wallColourCount]float64
	var points [typeCount]plotter.XYs
	for i, house := range houses {
		labels = append(labels, house.Price)
		squares = append(squares, house.Square)
		types = append(types, func() (res [typeCount]float64) {
			switch house.Type {
			case "Duplex":
				return [typeCount]float64{1, 0, 0, 0, 0}
			case "Detached":
				return [typeCount]float64{0, 1, 0, 0, 0}
			case "Townhouse":
				return [typeCount]float64{0, 0, 1, 0, 0}
			case "Semi-detached":
				return [typeCount]float64{0, 0, 0, 1, 0}
			case "Multi-family":
				return [typeCount]float64{0, 0, 0, 0, 1}
			default:
				log.Fatalf("Unknown house type: %v", house.Type)
			}
			return
		}())
		wallColours = append(wallColours, func() (selectors [wallColourCount]float64) {
			switch house.WallColor {
			case "blue":
				return [wallColourCount]float64{1, 0, 0, 0, 0}
			case "brown":
				return [wallColourCount]float64{0, 1, 0, 0, 0}
			case "white":
				return [wallColourCount]float64{0, 0, 1, 0, 0}
			case "green":
				return [wallColourCount]float64{0, 0, 0, 1, 0}
			case "yellow":
				return [wallColourCount]float64{0, 0, 0, 0, 1}
			}
			return

		}())
		for t := 0; t < typeCount; t++ {
			if types[i][t] == 1 {
				points[t] = append(points[t], plotter.XY{X: squares[i], Y: labels[i]})
			}
		}
	}

	if useWallColours {
		types = wallColours
	}

	img := make(chan *image.RGBA, 1)
	var pointScatter [typeCount]*plotter.Scatter
	for i := range pointScatter {
		pointScatter[i], err = plotter.NewScatter(points[i])
		if err != nil {
			log.Fatalf("Failed to create scatter: %v", err)
		}
		pointScatter[i].Color = typeColors[i]
	}

	go func() {
		var weights [typeCount * 2]float64
		for i := range weights {
			weights[i] = randMin + rand.Float64()*(randMax-randMin)
		}
		var weightDerivatives [typeCount * 2]float64 // Weight derivatives = Values of gradient projection onto the weight axis
		for epoch := 0; epoch < epochs; epoch++ {
			weightDerivatives = gradient(labels, inference(squares, types, weights), squares, types)
			for j := 0; j < typeCount*2; j++ {
				weights[j] -= weightDerivatives[j] * learningRates[j]
			}
			if epoch%100 == 0 {
				fmt.Printf("Epoch: %v, loss gradient: {%v}\n", epoch, weightDerivatives)
				fmt.Printf("Weights: %v\n", weights)
				fmt.Println()
			}

			FunctionWithColor := func(f func(x float64) float64, color color.RGBA) *plotter.Function {
				res := plotter.NewFunction(f)
				res.Color = color
				return res
			}

			select {
			case img <- Plot(
				pointScatter[0], pointScatter[1], pointScatter[2], pointScatter[3], pointScatter[4],
				FunctionWithColor(func(x float64) float64 { return weights[0]*x + weights[5] }, typeColors[0]),
				FunctionWithColor(func(x float64) float64 { return weights[1]*x + weights[6] }, typeColors[1]),
				FunctionWithColor(func(x float64) float64 { return weights[2]*x + weights[7] }, typeColors[2]),
				FunctionWithColor(func(x float64) float64 { return weights[3]*x + weights[8] }, typeColors[3]),
				FunctionWithColor(func(x float64) float64 { return weights[4]*x + weights[9] }, typeColors[4])):
			default:
			}
		}
	}()

	if err := ebiten.RunGame(&App{Img: img}); err != nil {
		log.Fatal(err)
	}
}
