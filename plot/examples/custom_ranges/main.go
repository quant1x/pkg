package main

//go:generate go run main.go

import (
	"os"

	"gitee.com/quant1x/pkg/plot"
)

func main() {
	/*
	   In this example we set a custom range for the y-axis, overriding the automatic range generation.
	   Note: the plot will still generate the ticks automatically based on the custom range, so the intervals may be a bit weird.
	*/

	graph := plot.Chart{
		YAxis: plot.YAxis{
			Range: &plot.ContinuousRange{
				Min: 0.0,
				Max: 10.0,
			},
		},
		Series: []plot.Series{
			plot.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
		},
	}
	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(plot.PNG, f)
}
