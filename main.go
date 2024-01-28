package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	lineMin, lineMax = 0, 120 //line lenght
)

func main() {

	//####################### Points #########################

	houses := getHouses()              //getting house data from CSV file
	px := make([]float64, len(houses)) //Point x coordinates
	py := make([]float64, len(houses)) //Point y coordinates

	//Saving point coordinates from house data
	for i, h := range houses {
		px[i] = h.square
		py[i] = h.price
	}

	//####################### Ebiten ####################################

	//Window
	ebiten.SetWindowSize(sW, sH)
	ebiten.SetWindowTitle("Linear Regression")

	//App instance
	a := NewApp(sW, sH)

	//Starting linear regression in another thread
	go func() {
		a.linearRegression(px, py)
	}()

	//Running game
	if err := ebiten.RunGame(a); err != nil {
		log.Fatal(err)
	}

}
