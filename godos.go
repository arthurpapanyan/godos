package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

type Config struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    map[string]string `json:"body"`
}

var elapsedTimes []time.Duration
var statusList []int

func main() {
	concurrency := flag.Int("c", 1, "Concurrency level")
	numRequests := flag.Int("n", 1, "Number of requests")
	targetURL := flag.String("t", "", "Target URL")
	method := flag.String("m", "GET", "HTTP request method")
	requestBody := flag.String("d", "", "Request body")
	logFilePath := flag.String("logfile", "", "Log file path")
	configFile := flag.String("config", "", "JSON config file path")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	if *targetURL == "" || (*configFile == "" && *requestBody == "") {
		flag.Usage()
		os.Exit(1)
	}

	if *configFile != "" {
		config, err := parseConfig(*configFile)
		if err != nil {
			fmt.Printf("Error parsing config file: %v\n", err)
			return
		}
		*targetURL = config.URL
		*method = config.Method
		flag.Set("d", fmt.Sprintf("%s", config.Body))
	}

	if *targetURL == "" {
		fmt.Println("Error: Please provide a target URL.")
		return
	}

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < *numRequests; j++ {
				fmt.Printf(".")
				makeRequest(*targetURL, *method, *requestBody, *logFilePath)
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)
	retreiveCallStats()
	countStatusCodes(statusList)
	fmt.Printf("\nTotal time taken: %s\n", elapsed)

}

func mapToStringifiedJSON(data map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func makeRequest(url, method, body, logFilePath string) {
	client := &http.Client{}
	req_start := time.Now()

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))

	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	req_elapsed := time.Since(req_start)
	elapsedTimes = append(elapsedTimes, req_elapsed)

	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	if logFilePath != "" {
		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening log file: %v\n", err)
			return
		}
		defer logFile.Close()

		log.SetOutput(logFile)
		statusList = append(statusList, resp.StatusCode)
	}
}

func parseConfig(filePath string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func retreiveCallStats() {
	sort.Slice(elapsedTimes, func(i, j int) bool {
		return elapsedTimes[i] < elapsedTimes[j]
	})
	fmt.Printf("\n--------------|\n")
	fmt.Printf(" Total Requests sent: %d\n", len(statusList))
	fmt.Printf(" Avarage: %s\n", elapsedTimes[(len(elapsedTimes)/2)])
	fmt.Printf("\n--------------|\n")
	fmt.Println(" Fastest Calls:")
	fmt.Println("")
	for i := 0; i < 5 && i < len(elapsedTimes); i++ {
		fmt.Printf("%d: %v\n", i+1, elapsedTimes[i])
	}
	fmt.Printf("\n--------------|\n")
	fmt.Println("Slowest Calls:")
	fmt.Println("")

	for j := len(elapsedTimes); j > len(elapsedTimes)-5; j-- {
		fmt.Printf("%v\n", elapsedTimes[(len(elapsedTimes)-j)])
	}
	fmt.Printf("--------------\n")

}

func countStatusCodes(statusCodes []int) {
	counts := make(map[int]int)
	for _, statusCode := range statusCodes {
		counts[statusCode]++
	}
	for statusCode, count := range counts {
		fmt.Printf("Status Code: %d, Count: %d\n", statusCode, count)
	}
}
