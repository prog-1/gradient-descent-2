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
	l := NewLine()
	x, y := Get()
	go func() {
		// fmt.Println(x, y)
		err := l.Train(x, y, 1e-4, 0.0003, 500000)
		fmt.Println(l.k, l.b)
		if err != nil {
			log.Fatal(err)
		}
		var tmp float64
		for {
			fmt.Print("Enter squares: ")
			fmt.Scan(&tmp)
			fmt.Println("You pobably can sell your hiuse for:", l.y(tmp))
		}
	}()
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Gradient descent")

	if err := ebiten.RunGame(&App{l, x, y}); err != nil {
		log.Fatal(err)
	}
}

func Abs(a float64) float64 {
	if a < 0 {
		return a * -1
	}
	return a
}
