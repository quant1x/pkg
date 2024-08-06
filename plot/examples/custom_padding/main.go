package main

//go:generate go run main.go

import (
	"os"

	"gitee.com/quant1x/pkg/plot"
	"gitee.com/quant1x/pkg/plot/drawing"
)

func main() {
	graph := plot.Chart{
		Background: plot.Style{
			Padding: plot.Box{
				Top:    50,
				Left:   25,
				Right:  25,
				Bottom: 10,
			},
			FillColor: drawing.ColorFromHex("efefef"),
		},
		Series: []plot.Series{
			plot.ContinuousSeries{
				XValues: plot.Seq{Sequence: plot.NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values(),
				YValues: plot.Seq{Sequence: plot.NewRandomSequence().WithLen(100).WithMin(100).WithMax(512)}.Values(),
			},
		},
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(plot.PNG, f)
}
