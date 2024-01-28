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

	px, py := groupHouses() //get squares and prices grouped by houseType

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
