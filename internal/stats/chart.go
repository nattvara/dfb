package stats

import (
	"bytes"
	"errors"
	"os"

	"github.com/nattvara/dfb/internal/fonts"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"github.com/wcharczuk/go-chart/util"
)

// LineChart is a line chart for a Metric
type LineChart struct {
	Metric     Metric
	Aggregator Aggregator
}

// WriteToFile writes LineChart c to file at given path
func (c *LineChart) WriteToFile(path string) error {
	graph := c.createGraph()

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return errors.New("failed to open file. " + err.Error())
	}

	buffer.WriteTo(file)
	return nil
}

// createGraph creates a graph for LineChart c
func (c *LineChart) createGraph() chart.Chart {
	latoRegular := fonts.GetFont(fonts.LatoRegular)
	latoBlack := fonts.GetFont(fonts.LatoBlack)

	return chart.Chart{
		Width:  2048,
		Height: 1024,
		Title:  c.Metric.GetTitle(),
		TitleStyle: chart.Style{
			Padding: chart.Box{
				Top: 50,
			},
			Show:      true,
			Font:      latoBlack,
			FontSize:  38,
			FontColor: chart.ColorWhite,
		},
		Background: chart.Style{
			FillColor: drawing.ColorFromHex("424242"),
			Padding: chart.Box{
				Top:    140,
				Left:   40,
				Right:  40,
				Bottom: 40,
			},
		},
		Canvas: chart.Style{
			FillColor: drawing.ColorFromHex("424242"),
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show:        true,
				Font:        latoRegular,
				FontColor:   drawing.ColorFromHex("fff"),
				FontSize:    18,
				StrokeWidth: 3,
			},
			TickStyle: chart.Style{
				Show:        true,
				StrokeColor: drawing.ColorFromHex("fff"),
				StrokeWidth: 3,
			},
			TickPosition: chart.TickPositionUnderTick,
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				typedDate := util.Time.FromFloat64(typed)
				return typedDate.Format(c.Metric.GetDateLayout())
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show:        true,
				Font:        latoRegular,
				FontColor:   drawing.ColorFromHex("fff"),
				FontSize:    18,
				StrokeWidth: 3,
			},
			TickStyle: chart.Style{
				Show:        true,
				StrokeColor: drawing.ColorFromHex("fff"),
				StrokeWidth: 3,
			},
			ValueFormatter: func(v interface{}) string {
				return c.Metric.GetFormatter().Format(v.(float64))
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: c.Metric.GetLabels(),
				YValues: c.Metric.GetValues(c.Aggregator),
				Style: chart.Style{
					Show:        true,
					StrokeColor: drawing.ColorFromHex("13c158"),
					FillColor:   drawing.ColorFromHex("13c158").WithAlpha(40),
					StrokeWidth: 4,
				},
			},
		},
	}
}
