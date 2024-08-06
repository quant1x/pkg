package main

//go:generate go run main.go

import (
	"os"

	"gitee.com/quant1x/pkg/plot"
	"gitee.com/quant1x/pkg/plot/drawing"
)

func main() {
	f, _ := plot.GetDefaultFont()
	r, _ := plot.PNG(1024, 1024)

	plot.Draw.Text(r, "Test", 64, 64, plot.Style{
		FontColor: drawing.ColorBlack,
		FontSize:  18,
		Font:      f,
	})

	plot.Draw.Text(r, "Test", 64, 64, plot.Style{
		FontColor:           drawing.ColorBlack,
		FontSize:            18,
		Font:                f,
		TextRotationDegrees: 45.0,
	})

	tb := plot.Draw.MeasureText(r, "Test", plot.Style{
		FontColor: drawing.ColorBlack,
		FontSize:  18,
		Font:      f,
	}).Shift(64, 64)

	tbc := tb.Corners().Rotate(45)

	plot.Draw.BoxCorners(r, tbc, plot.Style{
		StrokeColor: drawing.ColorRed,
		StrokeWidth: 2,
	})

	tbcb := tbc.Box()
	plot.Draw.Box(r, tbcb, plot.Style{
		StrokeColor: drawing.ColorBlue,
		StrokeWidth: 2,
	})

	file, _ := os.Create("output.png")
	defer file.Close()
	r.Save(file)
}
