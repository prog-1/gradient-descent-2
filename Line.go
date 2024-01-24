package main

import (
	"log"
)

type line struct {
	k []float64
	b float64
}

func (l *line) y(x []float64) float64 {
	var sum float64
	for i := range x {
		sum += l.k[i] * x[i]
	}
	return sum + l.b
}

func (l *line) abdabc(x float64, b int) float64 {
	a := make([]float64, 5)
	a[b] = 1
	return l.y(append(a, x))
}

func NewLine(num int) *line {
	return &line{make([]float64, num), 0}
}

func (l *line) Train(x [][]float64, y []float64, lrk, lrb float64, epochs int) error {
	t := make([]float64, len(x))
	// fmt.Println(points)
	for i := 0; i < epochs; i++ {
		// time.Sleep(500 * time.Millisecond)
		for j, v := range x {
			t[j] = l.y(v) //
		}
		// fmt.Println(y)

		var db float64
		var loss float64
		dk := make([]float64, len(x[0]))
		for j := range x {
			for k := range x[j] {
				if x[j][k] == 1 {
					dk[k] += (y[j] - t[j]) * x[j][k]
					dk[5] = (y[j] - t[j]) * x[j][5]
				}
			}
			db += y[j] - t[j]
			loss += (y[j] - t[j]) * (y[j] - t[j])

		}
		for j := range dk {
			l.k[j] += lrk * (2 / float64(len(x[j]))) * dk[j]
		}
		l.b += lrb * (2 / float64(len(x))) * db
		loss = loss / float64(len(x))
		if i%100 == 0 {
			log.Printf(`Epoch: %d/%d,
			Loss %f
			dk: %v, db: %f
			`, i, epochs, loss, dk, db)
		}
		// fmt.Println(dk, db)
	}
	return nil
}
