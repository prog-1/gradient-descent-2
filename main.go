package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/plot/plotter"
)

const (
	screenWidth, screenHeight = 720, 480
	randMin, randMax          = 1500, 2000
	epochs, weightCount       = 1e6, 6
)

const lrb = 0.7e-3

var learningRates = [weightCount]float64{0.1e-3, lrb, lrb, lrb, lrb, lrb}

// Prediction(inference) for one argument
func p(x, w [weightCount]float64) float64 {
	return w[0]*x[0] + w[1]*x[1] + w[2]*x[2] + w[3]*x[3] + w[4]*x[4] + w[5]*x[5]
}

// Runs model on all the input data
func inference(xs [][weightCount]float64, w [weightCount]float64) (ys []float64) {
	for i := range xs {
		ys = append(ys, p(xs[i], w))
	}
	return
}

func loss(labels, y []float64) float64 {
	var errSum float64
	for i := range labels {
		errSum += (y[i] - labels[i]) * (y[i] - labels[i])
	}
	return errSum / float64(len(labels)) // For the sake of making numbers smaller -> better percievable
}

func gradient(labels, y, x []float64, inputs [][weightCount]float64) (dw [weightCount]float64) {
	// dw, db - Parial derivatives, w - weight, b - bias
	for i := 0; i < len(labels); i++ {
		dif := y[i] - labels[i]
		dw[0] += dif * x[0]
	}

	for t := 1; t < weightCount; t++ {
		for i := 0; i < len(labels); i++ {
			if inputs[i][t] == 1 {
				dif := y[i] - labels[i]
				dw[t] += dif * x[t]
			}
		}
	}

	n := float64(len(labels))
	for i := range dw {
		dw[i] *= 2 / n
	}

	return
}

type House struct {
	Square    float64
	Type      string
	Price     float64
	WallColor string
}

func readHousesFromCSV(path string) ([]House, error) {
	file, err1 := os.Open(path)
	if err1 != nil {
		log.Fatalf("Can't open file with path: %v", path)
	}
	defer file.Close()
	houses := []House{}

	reader := csv.NewReader(file)
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
			Type:      houseType,
			Price:     price,
			WallColor: wallColor,
		}

		houses = append(houses, house)
	}

	return houses, nil
}

func main() {
	houses, err2 := readHousesFromCSV("data/house_prices.csv")
	if err2 != nil {
		log.Fatalf("Can't read Houses from CSV: %v", err2)
	}

	var inputs [][weightCount]float64
	var squares []float64
	var labels []float64
	var points plotter.XYs
	for i, house := range houses {
		inputs = append(inputs, func() (res [weightCount]float64) {
			switch house.Type {
			case "Duplex":
				return [weightCount]float64{house.Square, 1, 0, 0, 0, 0}
			case "Detached":
				return [weightCount]float64{house.Square, 0, 1, 0, 0, 0}
			case "Townhouse":
				return [weightCount]float64{house.Square, 0, 0, 1, 0, 0}
			case "Semi-detached":
				return [weightCount]float64{house.Square, 0, 0, 0, 1, 0}
			case "Multi-family":
				return [weightCount]float64{house.Square, 0, 0, 0, 0, 1}
			default:
				log.Fatalf("Unknown house type: %v", house.Type)
			}
			return
		}())
		labels = append(labels, house.Price)
		squares = append(squares, house.Square)
		points = append(points, plotter.XY{X: inputs[i][0], Y: labels[i]})
	}

	img := make(chan *image.RGBA, 1)
	pointScatter, _ := plotter.NewScatter(points)

	go func() {
		var weights [weightCount]float64
		for i := range weights {
			weights[i] = randMin + rand.Float64()*(randMax-randMin)
		}
		var weightDerivatives [weightCount]float64 // Weight derivatives = Values of gradient projection onto the weight axis
		for epoch := 0; epoch < epochs; epoch++ {
			weightDerivatives = gradient(labels, inference(inputs, weights), squares, inputs)
			for j := 0; j < weightCount; j++ {
				weights[j] -= weightDerivatives[j] * learningRates[j]
			}
			if epoch%100 == 0 {
				fmt.Printf("Epoch: %v, loss gradient: {%v}\n", epoch, weightDerivatives)
				fmt.Printf("Weights: %v\n", weights)
				fmt.Println()
			}
			select {
			case img <- Plot(pointScatter,
				plotter.NewFunction(func(x float64) float64 { return weights[0]*x + weights[1] }),
				plotter.NewFunction(func(x float64) float64 { return weights[0]*x + weights[2] }),
				plotter.NewFunction(func(x float64) float64 { return weights[0]*x + weights[3] }),
				plotter.NewFunction(func(x float64) float64 { return weights[0]*x + weights[4] }),
				plotter.NewFunction(func(x float64) float64 { return weights[0]*x + weights[5] })):
			default:
			}
		}
	}()

	if err := ebiten.RunGame(&App{Img: img}); err != nil {
		log.Fatal(err)
	}
}
