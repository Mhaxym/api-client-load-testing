package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var NET_CORE_URL string = "http://localhost:5000/Album/"
var GO_URL string = "http://localhost:8080/albums/"

var CONCURRENT_REQUESTS int = 10000
var MAX_REQUESTS int = CONCURRENT_REQUESTS * 100
var REQUESTS_PER_CONSUMER int = 1

var CURRENT_REQUESTS int = 0

var client *http.Client
var totalElapsedTime time.Duration = 0

func main() {

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		// Save logs into a file
		rand.Seed(time.Now().UnixNano())
		logFile, err := os.OpenFile(fmt.Sprintf("log%d.txt", rand.Intn(MAX_REQUESTS)), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer logFile.Close()
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)

		// Create a client with a custom transport
		tr := &http.Transport{
			MaxIdleConns:        CONCURRENT_REQUESTS,
			MaxIdleConnsPerHost: CONCURRENT_REQUESTS,
			IdleConnTimeout:     0,
		}
		client = &http.Client{Transport: tr}

		var url string
		switch argsWithoutProg[0] {
		case "GO":
			url = GO_URL
		case "NETCORE":
			url = NET_CORE_URL
		}

		log.Printf("Starting %s Load Test [URL %s | Concurrent Requests %d]", argsWithoutProg[0], url, CONCURRENT_REQUESTS)
		for CURRENT_REQUESTS < MAX_REQUESTS {
			LaunchTest(url)
			time.Sleep(1 * time.Second)
		}
		log.Printf("It tooks %s in total to Get %d", totalElapsedTime, CURRENT_REQUESTS)
	} else {
		log.Printf("Please specify which test you want to run. [GO, NETCORE]")
	}

}

func LaunchTest(URL string) {
	// Create a wait group to wait for all the consumers to finish
	var wg sync.WaitGroup
	start := time.Now()
	for consumerId := 0; consumerId < CONCURRENT_REQUESTS; consumerId++ {
		// Add a new goroutine to the wait group
		wg.Add(1)
		go func(wg *sync.WaitGroup, consumerId int) {
			// When the goroutine finishes, tell the wait group
			defer wg.Done()

			for i := 1; i <= REQUESTS_PER_CONSUMER; i++ {
				resp, err := client.Get(URL + strconv.Itoa(i) + "?consumerID=" + strconv.Itoa(consumerId) + "&totalRequests=" + strconv.Itoa(CURRENT_REQUESTS))

				if err != nil {
					log.Fatal(err)
				}

				defer resp.Body.Close()
				_, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
			}
		}(&wg, consumerId)
	}

	// Wait for all the consumers to finish
	wg.Wait()
	elapsed := time.Since(start)
	totalElapsedTime += elapsed
	CURRENT_REQUESTS += CONCURRENT_REQUESTS * REQUESTS_PER_CONSUMER
	log.Printf(
		"It tooks %s to get %d (Total Requests: %d/%d [%.f%%])",
		elapsed,
		CONCURRENT_REQUESTS*REQUESTS_PER_CONSUMER,
		CURRENT_REQUESTS,
		MAX_REQUESTS,
		(float64(CURRENT_REQUESTS)/float64(MAX_REQUESTS))*100,
	)
}
