package main

import (
	"log"
)

type model struct {
	w [6]float64
	b [6]float64
}

func (m *model) infrence(x [][]float64) []float64 {
	res := make([]float64, len(x))
	for i := range res {
		res[i] = m.y(x[i])
	}
	return res
}

// y = (x0 * w0 + x1 * w1 + x2 * w2 + x3 *w3 + x4 * w4 + w5)x5 + x1 * w6 + x2 * w7 + x3 * 8 + x4 * w9 + w5
func (l *model) y(x []float64) float64 {
	var sum float64
	for i := 0; i < 5; i++ {
		sum += l.w[i] * x[i] * x[5]
	}
	for i := 0; i < 5; i++ {
		sum += l.b[i] * x[i]
	}
	return sum
}

func (l *model) abdabc(x float64, b int) float64 {
	a := make([]float64, 5)
	a[b] = 1
	return l.y(append(a, x))
}

func NewLine(num int) *model {
	return &model{[6]float64{}, [6]float64{}}
}

func (m *model) Train(x [][]float64, y []float64, lrk, lrb float64, epochs int) {
	for i := 0; i < epochs; i++ {
		t := m.infrence(x)
		dk, db, loss := m.dmsl(x, y, t)
		for j := range dk {
			m.w[j] += lrk * (2 / float64(len(x[j]))) * dk[j]
		}
		for j := range db {
			m.b[j] += lrb * (2 / float64(len(x))) * db[j]
		}
		log.Printf(`Epoch: %d/%d,
			Loss %f 
			dk: %v, db: %v
			`, i, epochs, loss, dk, db)
	}
}

func (m *model) dmsl(x [][]float64, y []float64, lables []float64) (dk, db []float64, loss float64) {
	dk = make([]float64, 6)
	db = make([]float64, 6)
	for j := range x {
		for k := range x[j] {
			if x[j][k] == 1 {
				dk[k] += (y[j] - lables[j])
				dk[5] += (y[j] - lables[j]) * x[j][5]
			}
		}
		for k := range x[j] {
			if x[j][k] == 1 {
				db[k] += y[j] - lables[j]
				db[5] += y[j] - lables[j]
			}
		}
		// db += y[j] - lables[j]
		loss += (y[j] - lables[j]) * (y[j] - lables[j])

	}
	return dk, db, loss

}

// func (l *model) Train(x [][]float64, y []float64, lrk, lrb float64, epochs int) {
// 	t := make([]float64, len(x))
// 	// fmt.Println(points)
// 	for i := 0; i < epochs; i++ {
// 		// time.Sleep(500 * time.Millisecond)
// 		for j, v := range x {
// 			t[j] = l.y(v) //
// 		}
// 		// fmt.Println(y)

// 		var db float64
// 		var loss float64
// 		dk := make([]float64, 10)
// 		for j := range x {
// 			for k := 0; k < 5; k++ {
// 				if x[j][k] == 1 {
// 					dk[k] += (y[j] - t[j])
// 					dk[k+5] = (y[j] - t[j]) * x[j][5]
// 				}
// 			}
// 			db += y[j] - t[j]
// 			loss += (y[j] - t[j]) * (y[j] - t[j])

// 		}
// 		for j := range dk {
// 			l.k[j] += lrk * (2 / float64(len(x[j]))) * dk[j]
// 		}
// 		l.b += lrb * (2 / float64(len(x))) * db
// 		loss = loss / float64(len(x))
// 		if i%100 == 0 {
// 			log.Printf(`Epoch: %d/%d,
// 			Loss %f
// 			dk: %v, db: %v
// 			`, i, epochs, loss, dk, db)
// 		}
// 		// fmt.Println(dk, db)
// 	}
// }
