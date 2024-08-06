package main

//go:generate go run main.go

import (
	"fmt"
	"os"

	"gitee.com/quant1x/pkg/plot"
)

func main() {

	/*
	   In this example we add a second series, and assign it to the secondary y axis, giving that series it's own range.

	   We also enable all of the axes by setting the `Show` propery of their respective styles to `true`.
	*/

	graph := plot.Chart{
		XAxis: plot.XAxis{
			TickPosition: plot.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				typedDate := plot.TimeFromFloat64(typed)
				return fmt.Sprintf("%d-%d\n%d", typedDate.Month(), typedDate.Day(), typedDate.Year())
			},
		},
		Series: []plot.Series{
			plot.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
			plot.ContinuousSeries{
				YAxis:   plot.YAxisSecondary,
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{50.0, 40.0, 30.0, 20.0, 10.0},
			},
		},
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(plot.PNG, f)
}
