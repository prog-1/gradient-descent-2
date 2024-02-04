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
	useWallColours, wallColourCount = true, typeCount
)

// var typeColors [typeCount]color.RGBA = [typeCount]color.RGBA{{0, 0, 255, 255}, /*Blue*/
// 	{139, 69, 19, 255} /*Brown*/, {100, 100, 100, 255} /*White*/, {0, 255, 0, 255} /*Green*/, {255, 255, 0, 255} /*Yellow*/}

var red = color.RGBA{255, 0, 0, 255}
var typeColours = [typeCount]color.RGBA{red, red, red, red, red}

type weightRelated struct {
	w0, b   float64
	ws, bws [typeCount]float64
}

func randWeights() *weightRelated {
	randFloat := func() float64 {
		return randMin + rand.Float64()*(randMax-randMin)
	}
	randTypeCountArr := func() (arr [typeCount]float64) {
		for i := 0; i < typeCount; i++ {
			arr[i] = randFloat()
		}
		return arr
	}
	return &weightRelated{w0: randFloat(), b: randFloat(), ws: randTypeCountArr(), bws: randTypeCountArr()}
}

func (weights *weightRelated) adjustWeights(derivatives, learningRates *weightRelated) {
	weights.w0 -= derivatives.w0 * learningRates.w0
	weights.b -= derivatives.b * learningRates.b
	for i := 0; i < len(weights.ws); i++ {
		weights.ws[i] -= derivatives.ws[i] * learningRates.ws[i]
		weights.bws[i] -= derivatives.bws[i] * learningRates.bws[i]
	}
}

var learningRates weightRelated = weightRelated{
	w0: 1e-4, b: 8e-7,
	ws:  [typeCount]float64{1e-5, 1e-5, 1e-4, 1e-5, 1e-5},
	bws: [typeCount]float64{0.9, 0.5, 0.5, 0.5, 0.5},
}

func prediction(x float64, t [typeCount]float64, weights *weightRelated) float64 {
	y := weights.w0*x + weights.b
	for i := range t {
		y += weights.ws[i]*t[i]*x + weights.bws[i]*t[i]
	}

	return y
}

func inference(xs []float64, types [][typeCount]float64, weights *weightRelated) (ys []float64) {
	for i := range xs {
		ys = append(ys, prediction(xs[i], types[i], weights))
	}
	return
}

func gradient(labels, y, x []float64, types [][typeCount]float64) *weightRelated {
	var ds weightRelated // weight partial DerivativeS
	n := float64(len(labels))
	for i := 0; i < len(labels); i++ {
		dif := y[i] - labels[i]
		ds.w0 += 2 / n * dif * x[i]
		ds.b += 2 / n * dif
		for t := 0; t < typeCount; t++ {
			if types[i][t] == 1 {
				ds.ws[t] += 2 / n * dif * x[i]
				ds.bws[t] += 2 / n * dif
			}
		}
	}

	return &ds
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
				// Arrays are pointless here since slices do not allocate more memory than is needed
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
		pointScatter[i].Color = typeColours[i]
	}

	FunctionWithColor := func(f func(x float64) float64, color color.RGBA) *plotter.Function {
		res := plotter.NewFunction(f)
		res.Color = color
		return res
	}
	go func() {
		weights := randWeights()
		var derivatives *weightRelated
		for epoch := 0; epoch < epochs; epoch++ {
			derivatives = gradient(labels, inference(squares, types, weights), squares, types)
			weights.adjustWeights(derivatives, &learningRates)
			if epoch%1000 == 0 {
				//time.Sleep(50 * time.Millisecond)
				select {
				case img <- Plot(
					pointScatter[0], pointScatter[1], pointScatter[2], pointScatter[3], pointScatter[4],
					FunctionWithColor(func(x float64) float64 { return (weights.w0+weights.ws[0])*x + weights.b + weights.bws[0] }, typeColours[0]),
					FunctionWithColor(func(x float64) float64 { return (weights.w0+weights.ws[1])*x + weights.b + weights.bws[1] }, typeColours[1]),
					FunctionWithColor(func(x float64) float64 { return (weights.w0+weights.ws[2])*x + weights.b + weights.bws[2] }, typeColours[2]),
					FunctionWithColor(func(x float64) float64 { return (weights.w0+weights.ws[3])*x + weights.b + weights.bws[3] }, typeColours[3]),
					FunctionWithColor(func(x float64) float64 { return (weights.w0+weights.ws[4])*x + weights.b + weights.bws[4] }, typeColours[4])):
				default:
				}
				fmt.Printf("Epoch: %v\n\n", epoch)
				fmt.Printf("Weights:\nw0 = %v, ws = %v\nb = %v, bws = %v\n\n", weights.w0, weights.ws, weights.b, weights.bws)
				fmt.Printf("Derivatives:\nw0 = %v, ws = %v\nb = %v, bws = %v\n\n", derivatives.w0, derivatives.ws, derivatives.b, derivatives.bws)
				fmt.Println()
			}
		}
	}()

	if err := ebiten.RunGame(&App{Img: img}); err != nil {
		log.Fatal(err)
	}
}
