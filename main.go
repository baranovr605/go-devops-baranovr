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

var secondSleep time.Duration = 5

func errorCheck() {
	failedAttemps += 1
	if failedAttemps == 3 {
		fmt.Println("Unable to fetch server statistic")
		failedAttemps = 0
	}
}

func dataCheck(url string, sleepTime time.Duration) {
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

		percentMemoryUsage := float64(dataSlice[2]) * 100 / float64(dataSlice[1])
		if percentMemoryUsage > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", int(percentMemoryUsage))
		}

		percentDiskUsage := float64(dataSlice[4]) * 100 / float64(dataSlice[3])
		if percentDiskUsage > 90 {
			mbFreeDiskSpace := (float64(dataSlice[3]) - float64(dataSlice[4])) / (1024 * 1024)
			fmt.Printf("Free disk space is too low: %d Mb left\n", int(mbFreeDiskSpace))
		}

		percentBandwithNetwork := float64(dataSlice[6]) * 100 / float64(dataSlice[5])
		if percentBandwithNetwork > 90 {
			mbitfreeBandwidth := (float64(dataSlice[5]) - float64(dataSlice[6])) / 1000000
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", int(mbitfreeBandwidth))
		}

		time.Sleep(sleepTime * time.Second)
	}
}

func main() {
	dataCheck(endpointURL, secondSleep)
}
