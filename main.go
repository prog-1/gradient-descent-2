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
	epochs, typeCount         = 1e6, 5
	lrw, lrb                  = 0.1e-3, 0.7e-3
)

var learningRates = [typeCount * 2]float64{lrw, lrw, lrw, lrw, lrw, lrb, lrb, lrb, lrb, lrb}

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
				ds[t] += dif * x[t]
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
	houses, err := readHousesFromCSV("data/house_prices.csv")
	if err != nil {
		log.Fatalf("Can't read Houses from CSV: %v", err)
	}

	var types [][typeCount]float64
	var squares []float64
	var labels []float64
	var points plotter.XYs
	for i, house := range houses {
		labels = append(labels, house.Price)
		squares = append(squares, house.Square)
		points = append(points, plotter.XY{X: squares[i], Y: labels[i]})
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
	}

	img := make(chan *image.RGBA, 1)
	pointScatter, _ := plotter.NewScatter(points)

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
