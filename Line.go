package main

import (
	"log"
)

type line struct {
	k float64
	b float64
}

func (l *line) y(x float64) float64 {
	return l.k*x + l.b
}

func NewLine() *line {
	return &line{0, 0}
}

func (l *line) Train(x, y []float64, lrk, lrb float64, epochs int) error {
	t := make([]float64, len(x))
	// fmt.Println(points)
	for i := 0; i < epochs; i++ {
		// time.Sleep(100 * time.Millisecond)
		for j, v := range x {
			t[j] = l.y(v)
		}
		// fmt.Println(y)

		var db, dk float64
		var loss float64
		for j := range x {
			dk += (y[j] - t[j]) * x[j]
			db += y[j] - t[j]
			loss += (y[j] - t[j]) * (y[j] - t[j])

		}

		l.k += lrk * (2 / float64(len(x))) * dk
		l.b += lrb * (2 / float64(len(x))) * db
		loss = loss / float64(len(x))
		if i%100 == 0 {
			log.Printf(`Epoch: %d/%d,
			Loss %f
			dk: %f, db: %f
			k: %f, b:%f\n`, i, epochs, loss, dk, db, l.k, l.b)
		}
		// fmt.Println(dk, db)
	}
	return nil
}
