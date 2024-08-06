package main

//go:generate go run main.go

import (
	"os"

	"gitee.com/quant1x/pkg/plot"
)

func main() {

	/*
	   The below will draw the same plot as the `basic` example, except with both the x and y axes turned on.
	   In this case, both the x and y axis ticks are generated automatically, the x and y ranges are established automatically,
	   the canvas "box" is adjusted to fit the space the axes occupy so as not to clip.
	   Additionally, it shows how you can use the "Descending" property of continuous ranges to change the ordering of
	   how values (including ticks) are drawn.
	*/

	graph := plot.Chart{
		Height: 500,
		Width:  500,
		XAxis:  plot.XAxis{
			/*Range: &plot.ContinuousRange{
				Descending: true,
			},*/
		},
		YAxis: plot.YAxis{
			Range: &plot.ContinuousRange{
				Descending: true,
			},
		},
		Series: []plot.Series{
			plot.ContinuousSeries{
				Style: plot.Style{
					StrokeColor: plot.GetDefaultColor(0).WithAlpha(64),
					FillColor:   plot.GetDefaultColor(0).WithAlpha(64),
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
