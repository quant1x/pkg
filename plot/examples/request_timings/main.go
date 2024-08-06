package main

//go:generate go run main.go

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"gitee.com/quant1x/pkg/plot"
)

func main() {
	log := plot.NewLogger()
	drawChart(log)
}

func parseInt(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}

func parseFloat64(str string) float64 {
	v, _ := strconv.ParseFloat(str, 64)
	return v
}

func readData() ([]time.Time, []float64) {
	var xvalues []time.Time
	var yvalues []float64
	err := plot.ReadLines("requests.csv", func(line string) error {
		parts := plot.SplitCSV(line)
		year := parseInt(parts[0])
		month := parseInt(parts[1])
		day := parseInt(parts[2])
		hour := parseInt(parts[3])
		elapsedMillis := parseFloat64(parts[4])
		xvalues = append(xvalues, time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC))
		yvalues = append(yvalues, elapsedMillis)
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	return xvalues, yvalues
}

func releases() []plot.GridLine {
	return []plot.GridLine{
		{Value: plot.TimeToFloat64(time.Date(2016, 8, 1, 9, 30, 0, 0, time.UTC))},
		{Value: plot.TimeToFloat64(time.Date(2016, 8, 2, 9, 30, 0, 0, time.UTC))},
		{Value: plot.TimeToFloat64(time.Date(2016, 8, 2, 15, 30, 0, 0, time.UTC))},
		{Value: plot.TimeToFloat64(time.Date(2016, 8, 4, 9, 30, 0, 0, time.UTC))},
		{Value: plot.TimeToFloat64(time.Date(2016, 8, 5, 9, 30, 0, 0, time.UTC))},
		{Value: plot.TimeToFloat64(time.Date(2016, 8, 6, 9, 30, 0, 0, time.UTC))},
	}
}

func drawChart(log plot.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		xvalues, yvalues := readData()
		mainSeries := plot.TimeSeries{
			Name: "Prod Request Timings",
			Style: plot.Style{
				StrokeColor: plot.ColorBlue,
				FillColor:   plot.ColorBlue.WithAlpha(100),
			},
			XValues: xvalues,
			YValues: yvalues,
		}

		linreg := &plot.LinearRegressionSeries{
			Name: "Linear Regression",
			Style: plot.Style{
				StrokeColor:     plot.ColorAlternateBlue,
				StrokeDashArray: []float64{5.0, 5.0},
			},
			InnerSeries: mainSeries,
		}

		sma := &plot.SMASeries{
			Name: "SMA",
			Style: plot.Style{
				StrokeColor:     plot.ColorRed,
				StrokeDashArray: []float64{5.0, 5.0},
			},
			InnerSeries: mainSeries,
		}

		graph := plot.Chart{
			Log:    log,
			Width:  1280,
			Height: 720,
			Background: plot.Style{
				Padding: plot.Box{
					Top: 50,
				},
			},
			YAxis: plot.YAxis{
				Name: "Elapsed Millis",
				TickStyle: plot.Style{
					TextRotationDegrees: 45.0,
				},
				ValueFormatter: func(v interface{}) string {
					return fmt.Sprintf("%d ms", int(v.(float64)))
				},
			},
			XAxis: plot.XAxis{
				ValueFormatter: plot.TimeHourValueFormatter,
				GridMajorStyle: plot.Style{
					StrokeColor: plot.ColorAlternateGray,
					StrokeWidth: 1.0,
				},
				GridLines: releases(),
			},
			Series: []plot.Series{
				mainSeries,
				linreg,
				plot.LastValueAnnotationSeries(linreg),
				sma,
				plot.LastValueAnnotationSeries(sma),
			},
		}

		graph.Elements = []plot.Renderable{plot.LegendThin(&graph)}

		f, _ := os.Create("output.png")
		defer f.Close()
		graph.Render(plot.PNG, f)
	}
}
