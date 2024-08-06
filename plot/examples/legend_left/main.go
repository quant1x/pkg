package main

//go:generate go run main.go

import (
	"os"

	"gitee.com/quant1x/pkg/plot"
)

func main() {

	/*
	   In this example we add a `Renderable` or a custom component to the `Elements` array.
	   In this specific case it is a pre-built renderable (`CreateLegend`) that draws a legend for the plot's series.
	   If you like, you can use `CreateLegend` as a template for writing your own renderable, or even your own legend.
	*/

	graph := plot.Chart{
		Background: plot.Style{
			Padding: plot.Box{
				Top:  20,
				Left: 260,
			},
		},
		Series: []plot.Series{
			plot.ContinuousSeries{
				Name:    "A test series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			plot.ContinuousSeries{
				Name:    "Another test series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			plot.ContinuousSeries{
				Name:    "Yet Another test series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			plot.ContinuousSeries{
				Name:    "Even More series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			plot.ContinuousSeries{
				Name:    "Foo Bar",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			plot.ContinuousSeries{
				Name:    "Bar Baz",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			plot.ContinuousSeries{
				Name:    "Moo Bar",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			plot.ContinuousSeries{
				Name:    "Zoo Bar Baz",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			plot.ContinuousSeries{
				Name:    "Fast and the Furious",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{5.0, 4.0, 3.0, 2.0, 1.0},
			},

			plot.ContinuousSeries{
				Name:    "2 Fast 2 Furious",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{5.0, 4.0, 3.0, 2.0, 1.0},
			},

			plot.ContinuousSeries{
				Name:    "They only get more fast and more furious",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{5.0, 4.0, 3.0, 2.0, 1.0},
			},
		},
	}

	//note we have to do this as a separate step because we need a reference to graph
	graph.Elements = []plot.Renderable{
		plot.LegendLeft(&graph),
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(plot.PNG, f)
}
