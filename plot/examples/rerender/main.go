package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"gitee.com/quant1x/pkg/plot"
)

var lock sync.Mutex
var graph *plot.Chart
var ts *plot.TimeSeries

func addData(t time.Time, e time.Duration) {
	lock.Lock()
	ts.XValues = append(ts.XValues, t)
	ts.YValues = append(ts.YValues, plot.TimeMillis(e))
	lock.Unlock()
}

func drawChart(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		addData(start, time.Since(start))
	}()
	if len(ts.XValues) == 0 {
		http.Error(res, "no data (yet)", http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "image/png")
	if err := graph.Render(plot.PNG, res); err != nil {
		log.Printf("%v", err)
	}
}

func main() {
	ts = &plot.TimeSeries{
		XValues: []time.Time{},
		YValues: []float64{},
	}
	graph = &plot.Chart{
		Series: []plot.Series{ts},
	}
	http.HandleFunc("/", drawChart)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
