package main

//go:generate go run main.go

import (
	"os"

	"gitee.com/quant1x/pkg/plot"
)

func main() {
	graph := plot.Chart{
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
