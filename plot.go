package main

import (
	"image"
	"image/color"

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

//###################################################################################

// Returns new plot with given data
func (a *App) updatePlot(k float64, b []float64, px, py [houseTypeCount][]float64) {

	//################# Initialization ##########################

	p := plot.New() //initializing plot

	//House type colors
	colors := []color.RGBA{
		0: {150, 0, 0, 255},
		1: {0, 150, 0, 255},
		2: {0, 0, 150, 255},
		3: {150, 150, 0, 255},
		4: {0, 150, 150, 255}}

	//##################### Line ##############################

	//Create line for every houseType
	for ht := 0; ht < houseTypeCount; ht++ { //for every houseType

		//Line points
		lp := plotter.XYs{
			{X: lineMin, Y: inference(lineMin, k, b[ht])},
			{X: lineMax, Y: inference(lineMax, k, b[ht])},
		}

		line, _ := plotter.NewLine(lp) //creating line
		line.Color = colors[ht]

		p.Add(line) //adding line to the plot
	}

	//#################### Points ##############################

	//Create line for every houseType
	for ht := 0; ht < houseTypeCount; ht++ { //for every houseType
		var points plotter.XYs //initializing point plotter

		for i := 0; i < len(px[ht]); i++ { //for every point in houseType
			points = append(points, plotter.XY{X: px[ht][i], Y: py[ht][i]}) //Saving all points in plotter
		}
		scatter, _ := plotter.NewScatter(points) //creating new scatter from point data°
		scatter.Color = colors[ht]

		p.Add(scatter) //adding points to plot
	}

	//##################### App #############################

	a.plot = p //replacing old plot with new one

}
