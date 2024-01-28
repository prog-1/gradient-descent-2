package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/plot"
)

const (
	sW = 1250
	sH = 720
)

type App struct {
	width, height int        //screen width & height
	plot          *plot.Plot //global access to plot
}

func (a *App) Update() error {
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	if a.plot != nil { //to avoid crash at the start
		screen.DrawImage(PlotToImage(a.plot), &ebiten.DrawImageOptions{}) //drawing plot
	}
}

func (a *App) Layout(inWidth, inHeight int) (outWidth, outHeight int) {
	return a.width, a.height
}

func NewApp(width, height int) *App {
	return &App{width: width, height: height}
}
