package main

//go:generate go run main.go

import (
	"os"

	"gitee.com/quant1x/pkg/plot"
)

func main() {

	/*
	   In this example we set the primary YAxis to have logarithmic range.
	*/

	graph := plot.Chart{
		Background: plot.Style{
			Padding: plot.Box{
				Top:  20,
				Left: 20,
			},
		},
		Series: []plot.Series{
			plot.ContinuousSeries{
				Name:    "A test series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1, 10, 100, 1000, 10000},
			},
		},
		YAxis: plot.YAxis{
			Style:     plot.Shown(),
			NameStyle: plot.Shown(),
			Range:     &plot.LogarithmicRange{},
		},
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(plot.PNG, f)
}
