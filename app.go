package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type App struct {
	l *line
	x []float64
	y []float64
}

func Plot(ps ...plot.Plotter) *image.RGBA {
	p := plot.New()
	p.X.Min = 0
	p.X.Max = 300

	p.Add(append([]plot.Plotter{
		plotter.NewGrid(),
	}, ps...)...)

	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	c := vgimg.NewWith(vgimg.UseImage(img))
	p.Draw(draw.New(c))
	return c.Image().(*image.RGBA)
}
func (app *App) Update() error { return nil }

func (app *App) Draw(screen *ebiten.Image) {
	var points plotter.XYs
	for i := 0; i < len(app.x); i++ {
		points = append(points, plotter.XY{X: app.x[i], Y: app.y[i]})
	}
	pointsScatter, _ := plotter.NewScatter(points)
	// fp := plotter.NewFunction(app.l.y())
	*screen = *ebiten.NewImageFromImage(Plot(pointsScatter))
}

func (app *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
