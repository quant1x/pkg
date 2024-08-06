package main

//go:generate go run main.go

import (
	"os"

	"gitee.com/quant1x/pkg/plot"
	"gitee.com/quant1x/pkg/plot/drawing"
)

func main() {
	/*
	   In this example we set some custom colors for the series and the plot background and canvas.
	*/
	graph := plot.Chart{
		Background: plot.Style{
			FillColor: drawing.ColorBlue,
		},
		Canvas: plot.Style{
			FillColor: drawing.ColorFromHex("efefef"),
		},
		Series: []plot.Series{
			plot.ContinuousSeries{
				Style: plot.Style{
					StrokeColor: drawing.ColorRed,               // will supercede defaults
					FillColor:   drawing.ColorRed.WithAlpha(64), // will supercede defaults
				},
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
		},
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(plot.PNG, f)
}
