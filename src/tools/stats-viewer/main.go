package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/nattvara/dfb/src/internal/stats"
)

func main() {
	if len(os.Args) != 8 {
		fmt.Println("usage: stats-viewer [metric] [time_unit] [length] [aggregator] [repo] [group] [domain]")
		os.Exit(1)
	}

	metric := os.Args[1]
	timeUnit := os.Args[2]
	timeLength, _ := strconv.Atoi(os.Args[3])
	aggregatorName := os.Args[4]
	repoName := os.Args[5]
	groupName := os.Args[6]
	domainName := os.Args[7]

	fmt.Println("metric: " + metric)
	fmt.Println("time_unit: " + timeUnit)
	fmt.Printf("length: %v\n", timeLength)
	fmt.Println("repo: " + repoName)
	fmt.Println("group: " + groupName)
	fmt.Println("domain: " + domainName)

	db := stats.NewDB()
	db.Load(groupName)

	m := stats.NewMetric(metric, repoName, groupName, domainName, timeUnit, aggregatorName)
	m.FetchDataFromDB(db, timeUnit, timeLength)

	chart := stats.LineChart{
		Metric:     m,
		Aggregator: stats.NewAggregator(aggregatorName),
	}
	chart.WriteToFile("foo.png")
}
