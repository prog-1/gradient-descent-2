package main

import (
	"fmt"
	"time"
)

const (
	lrB    = 0.1
	lrK    = 0.0002
	epochs = 1000
)

//#######################################################################

func loss(k, b float64, px, py []float64) float64 {
	totalE := 0.0 // error
	for i := range px {
		x := px[i]
		y := py[i]
		totalE += (y - (k*x + b)) * 2
	}
	totalE /= float64(len(px))

	return totalE
}

func inference(x, k, b float64) float64 {
	return k*x + b
}

func gradientDescent(k, b float64, px, py []float64, epoch int) (float64, float64) {
	dk, db := 0.0, 0.0 // gradients for coefficients
	n := float64(len(px))
	for i := range px {
		x := px[i]
		y := py[i]
		dk -= (2 / n) * (y - (k*x + b)) * x
		db -= (2 / n) * (y - (k*x + b))
	}
	k -= dk * lrK
	b -= db * lrB
	if epoch%100 == 0 {
		fmt.Println("dk:", dk, "db:", db, "\n")
	}
	return k, b
}

func (a *App) linearRegression(px, py []float64) (k, b float64) {
	for epoch := 1; epoch <= epochs; epoch++ {
		if epoch%100 == 0 {
			fmt.Println("Epoch:", epoch, "Loss:", loss(k, b, px, py))
		}
		k, b = gradientDescent(k, b, px, py, epoch)
		a.updatePlot(k, b, px, py)
		time.Sleep(time.Millisecond)
	}

	return k, b
}

//#######################################################################
