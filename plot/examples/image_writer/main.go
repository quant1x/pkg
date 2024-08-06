package main

import (
	"fmt"
	"log"

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
	collector := &plot.ImageWriter{}
	graph.Render(plot.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Final Image: %dx%d\n", image.Bounds().Size().X, image.Bounds().Size().Y)
}
