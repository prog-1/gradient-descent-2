package main

import "math/rand"

type point struct {
	x, y float64
}

func NewRandomPoints(num int) (x, y []float64) {
	x = make([]float64, num)
	y = make([]float64, num)
	ff := func(x float64) float64 {
		return 10*x - 5
	}
	for i := 0; i < num; i++ {
		x[i] = rand.Float64() * 300
		y[i] = ff(x[i])
	}
	return
}
