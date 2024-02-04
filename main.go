package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	l := NewLine(10)
	x, y := Get()
	go func() {
		// fmt.Println(x, y)
		l.Train(x, y, 1e-4, 1e-4, 5500000)
		fmt.Println(l.w, l.b)

		var tmp float64
		var houseType string
		for {
			fmt.Print("Enter squares: ")
			fmt.Scan(&tmp)
			fmt.Println("Enter house type: ")
			fmt.Scan(&houseType)
			fmt.Println("You pobably can sell your house for:", l.y(append(enum(houseType), tmp)))
		}
	}()
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Gradient descent")
	p := make([]float64, len(x))
	for i := range p {
		p[i] = x[i][5]
	}

	if err := ebiten.RunGame(&App{l, p, y, 0}); err != nil {
		log.Fatal(err)
	}
}

func Abs(a float64) float64 {
	if a < 0 {
		return a * -1
	}
	return a
}
