package main

import (
	"bufio"
	"github.com/labstack/gommon/color"
	"log"
	"os"
	"sync"
)

var checkPre = color.Yellow("[") + color.Green("✓") + color.Yellow("]")
var tildPre = color.Yellow("[") + color.Green("~") + color.Yellow("]")
var crossPre = color.Yellow("[") + color.Red("✗") + color.Yellow("]")

func main() {
	var worker sync.WaitGroup
	var count, index int

	// Parse arguments
	parseArgs(os.Args)

	file, err := os.Open(arguments.LinksFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		worker.Add(1)
		count++
		go download(&worker, scanner.Text())
		index++
		if count == arguments.Concurrency {
			worker.Wait()
			count = 0
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
