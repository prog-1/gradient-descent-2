package main

import (
	"image"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

//// Old one
// func Plot(xmin, xmax, xstep float64, f func(float64) float64) func(x float64) *image.RGBA {
// 	var pts plotter.XYs
// 	for x := xmin; x <= xmax; x += xstep {
// 		pts = append(pts, plotter.XY{X: x, Y: f(x)})
// 	}
// 	fn, err := plotter.NewLine(pts)
// 	if err != nil {
// 		log.Fatalf("Failed to NewLine: %v", err)
// 	}
// 	fn.Color = color.RGBA{B: 255, A: 255}
// 	return func(x float64) *image.RGBA {
// 		pts := plotter.XYs{{X: x, Y: f(x)}}
// 		xScatter, err := plotter.NewScatter(pts)
// 		if err != nil {
// 			log.Fatalf("Failed to NewScatter: %v", err)
// 		}
// 		xScatter.Color = color.RGBA{R: 255, A: 255}
//
// 		labels, err := plotter.NewLabels(plotter.XYLabels{
// 			XYs:    pts,
// 			Labels: []string{fmt.Sprintf("x = %.2f", x)},
// 		})
// 		labels.Offset = vg.Point{X: -10, Y: 15}
// 		if err != nil {
// 			log.Fatalf("Failed to NewLabels: %v", err)
// 		}
//
// 		p := plot.New()
// 		p.Add(
// 			plotter.NewGrid(),
// 			fn,
// 			xScatter,
// 			labels,
// 		)
// 		p.Legend.Add("f(x)", fn)
// 		p.Legend.Add("x", xScatter)
// 		p.X.Label.Text = "X"
// 		p.Y.Label.Text = "Y"
//
// 		img := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
// 		c := vgimg.NewWith(vgimg.UseImage(img))
// 		p.Draw(draw.New(c))
// 		return c.Image().(*image.RGBA)
// 	}
// }

func Plot(ps ...plot.Plotter) *image.RGBA {
	p := plot.New()
	p.Add(append([]plot.Plotter{
		plotter.NewGrid(),
	}, ps...)...)

	img := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	c := vgimg.NewWith(vgimg.UseImage(img))
	p.Draw(draw.New(c))
	return c.Image().(*image.RGBA)
}
