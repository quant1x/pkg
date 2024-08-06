package main

import (
	"os"

	"gitee.com/quant1x/pkg/plot"
)

func main() {
	sbc := plot.StackedBarChart{
		Title: "Test Stacked Bar Chart",
		Background: plot.Style{
			Padding: plot.Box{
				Top: 40,
			},
		},
		Height: 512,
		Bars: []plot.StackedBar{
			{
				Name: "This is a very long string to test word break wrapping.",
				Values: []plot.Value{
					{Value: 5, Label: "Blue"},
					{Value: 5, Label: "Green"},
					{Value: 4, Label: "Gray"},
					{Value: 3, Label: "Orange"},
					{Value: 3, Label: "Test"},
					{Value: 2, Label: "??"},
					{Value: 1, Label: "!!"},
				},
			},
			{
				Name: "Test",
				Values: []plot.Value{
					{Value: 10, Label: "Blue"},
					{Value: 5, Label: "Green"},
					{Value: 1, Label: "Gray"},
				},
			},
			{
				Name: "Test 2",
				Values: []plot.Value{
					{Value: 10, Label: "Blue"},
					{Value: 6, Label: "Green"},
					{Value: 4, Label: "Gray"},
				},
			},
		},
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	sbc.Render(plot.PNG, f)
}
