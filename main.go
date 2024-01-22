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
	epochs, lrw, lrb          = 1e6, 0.5e-4, 0.7
)

// Function points are spawed along
func f(x float64) float64 {
	// return 0.5*x + 2
	return 10*x - 5
}

// Inference for 1 argument(x)
func i(x, w, b float64) float64 { return w*x + b }

// Runs model on all the input data
func inference(x []float64, w, b float64) (out []float64) {
	for _, v := range x {
		out = append(out, i(v, w, b))
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

func gradient(labels, y, x []float64) (dw, db float64) {
	// dw, db - Parial derivatives, w - weight, b - bias
	for i := 0; i < len(labels); i++ {
		dif := y[i] - labels[i]
		dw += dif * x[i]
		db += dif
	}
	n := float64(len(labels))
	dw *= 2 / n
	db *= 2 / n

	return
}

func train(epochs int, inputs, labels []float64) (w, b float64) {
	randFloat64 := func() float64 {
		return randMin + rand.Float64()*(randMax-randMin)
	}
	w, b = randFloat64(), randFloat64()
	// w, b = 1, 0
	var dw, db float64
	for i := 0; i < epochs; i++ {
		dw, db = gradient(labels, inference(inputs, w, b), inputs)
		w -= dw * lrw
		b -= db * lrb
	}
	return
}

type House struct {
	Square    float64
	HouseType string
	Price     float64
	WallColor string
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

func main() {
	path := "data/house_prices.csv"
	file, err1 := os.Open(path)
	if err1 != nil {
		log.Fatalf("Can't open file with path: %v", path)
	}
	defer file.Close()
	houses, err2 := readHousesFromCSV(file)
	if err2 != nil {
		log.Fatal("Can't read Houses from CSV: %v", err2)
	}

	var inputs, labels []float64
	var points plotter.XYs
	for i, house := range houses {
		inputs = append(inputs, house.Square)
		labels = append(labels, house.Price)
		points = append(points, plotter.XY{X: inputs[i], Y: labels[i]})
	}

	img := make(chan *image.RGBA, 1)
	pointsScatter, _ := plotter.NewScatter(points)
	fp := plotter.NewFunction(f) // f plot

	go func() {
		randFloat64 := func() float64 {
			return randMin + rand.Float64()*(randMax-randMin)
		}
		w, b := randFloat64(), randFloat64()
		var dw, db float64
		for i := 0; i < epochs; i++ {
			// time.Sleep(1 * time.Millisecond)
			dw, db = gradient(labels, inference(inputs, w, b), inputs)
			w -= dw * lrw
			b -= db * lrb
			if i%100 == 0 {
				fmt.Printf("Epoch: %v, loss gradient: {%v,%v}\n", i, dw, db)
			}
			ap := plotter.NewFunction(func(x float64) float64 { return w*x + b }) // approximating function plot
			// Channels in go are blocking, i.e. until file are not read, new ones cannot be pasted
			// If there is something in channel we are not rendering anything
			select {
			case img <- Plot(pointsScatter, fp, ap): // Executes on successful pasting into channel
			default:
			} // In case of just adding writing to channel, we'll wait until we can write. Here we ignore it.
		}
	}()

	if err := ebiten.RunGame(&App{Img: img}); err != nil {
		log.Fatal(err)
	}
}
