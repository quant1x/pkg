package main

//go:generate go run main.go

import (
	"os"

	"gitee.com/quant1x/pkg/plot"
)

func main() {
	mainSeries := plot.ContinuousSeries{
		Name:    "A test series",
		XValues: plot.Seq{Sequence: plot.NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values(),
		YValues: plot.Seq{Sequence: plot.NewRandomSequence().WithLen(100).WithMin(50).WithMax(150)}.Values(),
	}

	minSeries := &plot.MinSeries{
		Style: plot.Style{
			StrokeColor:     plot.ColorAlternateGray,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}

	maxSeries := &plot.MaxSeries{
		Style: plot.Style{
			StrokeColor:     plot.ColorAlternateGray,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}

	graph := plot.Chart{
		Width:  1920,
		Height: 1080,
		YAxis: plot.YAxis{
			Name: "Random Values",
			Range: &plot.ContinuousRange{
				Min: 25,
				Max: 175,
			},
		},
		XAxis: plot.XAxis{
			Name: "Random Other Values",
		},
		Series: []plot.Series{
			mainSeries,
			minSeries,
			maxSeries,
			plot.LastValueAnnotationSeries(minSeries),
			plot.LastValueAnnotationSeries(maxSeries),
		},
	}

	graph.Elements = []plot.Renderable{plot.Legend(&graph)}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(plot.PNG, f)
}
