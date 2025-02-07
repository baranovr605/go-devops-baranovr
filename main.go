package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var endpointURL string = "http://srv.msk01.gigacorp.local/_stats"
var failedAttemps int = 0

// var secondSleep time.Duration = 1

func errorCheck() {
	failedAttemps += 1
	if failedAttemps == 3 {
		fmt.Println("Unable to fetch server statistic")
		failedAttemps = 0
	}
}

func dataCheck(url string) {
	for {
		res, err := http.Get(url)

		if err != nil {
			errorCheck()
		}

		if res.StatusCode != 200 {
			errorCheck()
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			errorCheck()
		}

		fullData := strings.Split(string(body), ",")

		dataSlice := make([]int, 7)

		for i, s := range fullData {
			num, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println("Error read data:", err)
				return
			}
			dataSlice[i] = num
		}

		if dataSlice[0] > 30 {
			fmt.Printf("Load Average is too high: %d\n", dataSlice[0])
		}

		percentMemoryUsage := dataSlice[2] * 100 / dataSlice[1]
		if percentMemoryUsage > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", percentMemoryUsage)
		}

		percentDiskUsage := (dataSlice[4] / dataSlice[3]) * 100
		if percentDiskUsage > 90 {
			mbFreeDiskSpace := (dataSlice[3] - dataSlice[4]) / (1024 * 1024)
			fmt.Printf("Free disk space is too low: %d Mb left\n", mbFreeDiskSpace)
		}

		percentBandwithNetwork := (dataSlice[6] / dataSlice[5]) * 100
		if percentBandwithNetwork > 90 {
			mbitfreeBandwidth := (dataSlice[5] - dataSlice[6]) / 1000000
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", mbitfreeBandwidth)
		}

		time.Sleep(time.Second)
	}
}

func main() {
	dataCheck(endpointURL)
}
