package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	e = 1e-5
)

func main() {
	l := NewLine(6)
	x, y := Get()
	go func() {
		// fmt.Println(x, y)
		err := l.Train(x, y, 0.5e-4, 0.7, 500000)
		fmt.Println(l.k, l.b)
		if err != nil {
			log.Fatal(err)
		}
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

	if err := ebiten.RunGame(&App{l, x[0], y}); err != nil {
		log.Fatal(err)
	}
}

func Abs(a float64) float64 {
	if a < 0 {
		return a * -1
	}
	return a
}
