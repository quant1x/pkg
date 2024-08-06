package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gitee.com/quant1x/pkg/plot"
)

var (
	outputPath = flag.String("output", "", "The output file")

	inputFormat = flag.String("format", "csv", "The input format, either 'csv' or 'tsv' (defaults to 'csv')")
	inputPath   = flag.String("f", "", "The input file")
	reverse     = flag.Bool("reverse", false, "If we should reverse the inputs")

	hideLegend     = flag.Bool("hide-legend", false, "If we should omit the plot legend")
	hideSMA        = flag.Bool("hide-sma", false, "If we should omit simple moving average")
	hideLinreg     = flag.Bool("hide-linreg", false, "If we should omit linear regressions")
	hideLastValues = flag.Bool("hide-last-values", false, "If we should omit last values")
)

func main() {
	flag.Parse()
	log := plot.NewLogger()

	var rawData []byte
	var err error
	if *inputPath != "" {
		if *inputPath == "-" {
			rawData, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.FatalErr(err)
			}
		} else {
			rawData, err = ioutil.ReadFile(*inputPath)
			if err != nil {
				log.FatalErr(err)
			}
		}
	} else if len(flag.Args()) > 0 {
		rawData = []byte(flag.Args()[0])
	} else {
		flag.Usage()
		os.Exit(1)
	}

	var parts []string
	switch *inputFormat {
	case "csv":
		parts = plot.SplitCSV(string(rawData))
	case "tsv":
		parts = strings.Split(string(rawData), "\t")
	default:
		log.FatalErr(fmt.Errorf("invalid format; must be 'csv' or 'tsv'"))
	}

	yvalues, err := plot.ParseFloats(parts...)
	if err != nil {
		log.FatalErr(err)
	}

	if *reverse {
		yvalues = plot.ValueSequence(yvalues...).Reverse().Values()
	}

	var series []plot.Series
	mainSeries := plot.ContinuousSeries{
		Name:    "Values",
		XValues: plot.LinearRange(1, float64(len(yvalues))),
		YValues: yvalues,
	}
	series = append(series, mainSeries)

	smaSeries := &plot.SMASeries{
		Name: "SMA",
		Style: plot.Style{
			Hidden:          *hideSMA,
			StrokeColor:     plot.ColorRed,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}
	series = append(series, smaSeries)

	linRegSeries := &plot.LinearRegressionSeries{
		Name: "Values - Lin. Reg.",
		Style: plot.Style{
			Hidden: *hideLinreg,
		},
		InnerSeries: mainSeries,
	}
	series = append(series, linRegSeries)

	mainLastValue := plot.LastValueAnnotationSeries(mainSeries)
	mainLastValue.Style = plot.Style{
		Hidden: *hideLastValues,
	}
	series = append(series, mainLastValue)

	linregLastValue := plot.LastValueAnnotationSeries(linRegSeries)
	linregLastValue.Style = plot.Style{
		Hidden: (*hideLastValues || *hideLinreg),
	}
	series = append(series, linregLastValue)

	smaLastValue := plot.LastValueAnnotationSeries(smaSeries)
	smaLastValue.Style = plot.Style{
		Hidden: (*hideLastValues || *hideSMA),
	}
	series = append(series, smaLastValue)

	graph := plot.Chart{
		Background: plot.Style{
			Padding: plot.Box{
				Top: 50,
			},
		},
		Series: series,
	}

	if !*hideLegend {
		graph.Elements = []plot.Renderable{plot.LegendThin(&graph)}
	}

	var output *os.File
	if *outputPath != "" {
		output, err = os.Create(*outputPath)
		if err != nil {
			log.FatalErr(err)
		}
	} else {
		output, err = ioutil.TempFile("", "*.png")
		if err != nil {
			log.FatalErr(err)
		}
	}

	if err := graph.Render(plot.PNG, output); err != nil {
		log.FatalErr(err)
	}

	fmt.Fprintln(os.Stdout, output.Name())
	os.Exit(0)
}
