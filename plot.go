package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

// Converting plot to ebiten.Image
func PlotToImage(p *plot.Plot) *ebiten.Image {

	img := image.NewRGBA(image.Rect(0, 0, sW, sH)) //creating image.RGBA to store the plot

	c := vgimg.NewWith(vgimg.UseImage(img)) //creating plot drawer for the image

	p.Draw(draw.New(c)) //drawing plot on the image

	return ebiten.NewImageFromImage(c.Image()) //converting image.RGBA to ebiten.Image (doing in Draw)
	///Black screen issue: was giving "img" instead of "c.Image()" in the function.
}

// recreating plot with given data
func (a *App) updatePlot(k, b float64, px, py []float64) {

	p := plot.New() //initializing plot

	//##################################################

	//Line

	linePoints := plotter.XYs{
		{X: lineMin, Y: inference(lineMin, k, b)},
		{X: lineMax, Y: inference(lineMax, k, b)},
	}

	line, _ := plotter.NewLine(linePoints) //creating line

	p.Add(line) //adding line to the plot°

	//##################################################

	//Points

	var points plotter.XYs //initializing point plotter

	for i := 0; i < len(px); i++ {
		points = append(points, plotter.XY{X: px[i], Y: py[i]}) //Saving all points in plotter
	}

	scatter, _ := plotter.NewScatter(points) //creating new scatter from point data

	p.Add(scatter) //adding points to plot

	//##################################################

	//App

	a.plot = p //replacing old plot with new one

}
