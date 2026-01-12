package main

import (
	"math"
	"net/rpc"
	"sync"
	"time"
)

// Mutex-based aggregator that reports the global average temperature periodically
//
// Report the average temperature across all k weather stations every averagePeriod
// second by sending it to the out channel. The aggregator should terminate upon
// receiving a signal on the quit channel.
//
// Note! To receive credit, mutexAggregator must implement a mutex based solution.
func mutexAggregator(
	k int,
	averagePeriod float64,
	out chan float64,
	quit chan struct{},
	client *rpc.Client,
) {
	//declare variables
	var tempSum float64
	var count int
	var batch int = 0    //counting overall batch
	var currentBatch = 0 //current batch of each rpc

	//create mutex's
	var mu sync.Mutex
	var batchMu sync.Mutex

	//create ticker
	ticker := time.NewTicker(time.Duration(averagePeriod * float64(time.Second)))
	defer ticker.Stop()

	//for loop - select case
	for {

		//rpc calls
		for j := 0; j < k; j++ {
			go func(stationID int, batch int) {

				//set current batch of this rpc -- use lock to avoid datarace
				batchMu.Lock()
				currentBatch = batch
				batchMu.Unlock()

				//make rpc call
				temp, err := GetWeatherData(client, stationID)
				if err != nil {
					return
				}

				//update tempSum and count (w/ locking to avoid data race)
				mu.Lock()
				if !math.IsNaN(temp) && batch == currentBatch {
					tempSum += temp
					count++
				}
				mu.Unlock()

			}(j, batch)
		}

		select {
		//averagePeriod timeout protocol
		case <-ticker.C:
			mu.Lock() //aquire lock to avoid data race

			//incremement batch
			batch++
			//send current average through out channel
			out <- tempSum / float64(count)

			//reset tempSum and count
			tempSum, count = 0, 0

			mu.Unlock() //now we can unlock

		case <-quit:
			return // recieve quite message and return

		}
	}

}
