package main

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"gitee.com/quant1x/pkg/plot"
	"gitee.com/quant1x/pkg/plot/drawing"
)

func drawChart(res http.ResponseWriter, req *http.Request) {

	viridisByY := func(xr, yr plot.Range, index int, x, y float64) drawing.Color {
		return plot.Viridis(y, yr.GetMin(), yr.GetMax())
	}

	graph := plot.Chart{
		Series: []plot.Series{
			plot.ContinuousSeries{
				Style: plot.Style{
					StrokeWidth:      plot.Disabled,
					DotWidth:         5,
					DotColorProvider: viridisByY,
				},
				XValues: plot.Seq{Sequence: plot.NewLinearSequence().WithStart(0).WithEnd(127)}.Values(),
				YValues: plot.Seq{Sequence: plot.NewRandomSequence().WithLen(128).WithMin(0).WithMax(1024)}.Values(),
			},
		},
	}

	res.Header().Set("Content-Type", plot.ContentTypePNG)
	err := graph.Render(plot.PNG, res)
	if err != nil {
		log.Println(err.Error())
	}
}

func unit(res http.ResponseWriter, req *http.Request) {
	graph := plot.Chart{
		Height: 50,
		Width:  50,
		Canvas: plot.Style{
			Padding: plot.BoxZero,
		},
		Background: plot.Style{
			Padding: plot.BoxZero,
		},
		Series: []plot.Series{
			plot.ContinuousSeries{
				XValues: plot.Seq{Sequence: plot.NewLinearSequence().WithStart(0).WithEnd(4)}.Values(),
				YValues: plot.Seq{Sequence: plot.NewLinearSequence().WithStart(0).WithEnd(4)}.Values(),
			},
		},
	}

	res.Header().Set("Content-Type", plot.ContentTypePNG)
	err := graph.Render(plot.PNG, res)
	if err != nil {
		log.Println(err.Error())
	}
}

func main() {
	http.HandleFunc("/", drawChart)
	http.HandleFunc("/unit", unit)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
