package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/nattvara/dfb/internal/stats"

	"github.com/spf13/cobra"
)

var domainName string

var timeUnit string

var timeLength int

var aggregatorName string

var outputPath string

var shouldListMetrics bool

var shouldListTimeUnits bool

var shouldListAggregators bool

var cmd = &cobra.Command{
	Use:   "stats [group] [repo] [metric]",
	Short: "Make a chart for a backup metric",
	Long:  "The stats command allows a user to view metrics about the backed up data",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]
		repoName := args[1]
		metricName := args[2]

		var metric Metric
		var aggregator stats.Aggregator
		var err error

		if domainName == "" {
			domainName = stats.AllDomains
		}

		db := stats.NewDB()
		db.Load(groupName)

		if metric, err = stats.NewMetric(
			metricName,
			repoName,
			groupName,
			domainName,
			timeUnit,
			aggregatorName,
		); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		metric.FetchDataFromDB(db, timeUnit, timeLength)

		if aggregator, err = stats.NewAggregator(aggregatorName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		chart := stats.LineChart{
			Metric:     metric,
			Aggregator: aggregator,
		}

		err = chart.WriteToFile(outputPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func main() {
	cmd.Flags().StringVarP(&domainName, "domain", "d", "", "which domain to use for metric, not availiable for all metrics, optional/required for some metrics")
	cmd.Flags().StringVarP(&timeUnit, "time-unit", "u", stats.TimeUnitDays, "time unit to use for metric")
	cmd.Flags().IntVarP(&timeLength, "time-length", "l", 7, "how many time-units of history should be included")
	cmd.Flags().StringVarP(&aggregatorName, "aggregator", "a", "sum", "aggregation method to use for a metric")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "/tmp/dfb-metric.png", "output path for png image of metric")
	cmd.Flags().BoolVarP(&shouldListMetrics, "list-metrics", "", false, "list availiable metrics")
	cmd.Flags().BoolVarP(&shouldListTimeUnits, "list-time-units", "", false, "list availiable time units")
	cmd.Flags().BoolVarP(&shouldListAggregators, "list-aggregators", "", false, "list availiable aggregators")
	cmd.Execute()

	if shouldListMetrics {
		listMetrics()
		return
	}

	if shouldListTimeUnits {
		listTimeUnits()
		return
	}

	if shouldListAggregators {
		listAggregators()
		return
	}
}

func listMetrics() {
	var metrics []string

	for name := range stats.Metrics {
		metrics = append(metrics, name)
	}

	sort.Strings(metrics)
	fmt.Println("availible metrics are:")
	for _, name := range metrics {
		fmt.Printf("  %s\n", name)
	}
}

func listTimeUnits() {
	fmt.Println("availible time units are:")
	for _, unit := range stats.TimeUnits {
		fmt.Printf("  %s\n", unit)
	}
}

func listAggregators() {
	fmt.Println("availible aggregators are:")
	for name := range stats.Aggregators {
		fmt.Printf("  %s\n", name)
	}
	fmt.Println("\nnote: not all of these aggregators makes sense for all metrics")
}
