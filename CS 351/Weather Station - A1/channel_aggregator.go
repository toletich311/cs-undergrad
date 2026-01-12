package main

import (
	"fmt"
	"net/rpc"
	"time"
)

// Channel-based aggregator that reports the global average temperature periodically
//
// Report the average temperature across all k weather stations every averagePeriod
// second by sending it to the out channel. The aggregator should terminate upon
// receiving a signal on the quit channel.
//
// Note! To receive credit, channelAggregator must not use mutexes.
func channelAggregator(
	k int,
	averagePeriod float64,
	out chan float64,
	quit chan struct{},
	client *rpc.Client,
) {
	//gettingWeatherData helper function to make rpc calls
	gettingWeatherData := func() chan float64 {

		//make new channel each function call!!!! -- creates "batch" functionality
		tempChan := make(chan float64, k)

		for j := 0; j < k; j++ {
			//call go routine in loop
			go func(stationID int) {
				//make rpc call to getweather data
				temp, err := GetWeatherData(client, stationID)
				if err != nil {
					fmt.Printf("Error fetching data for station %d: %v\n", stationID, err)
				}
				//send information through tempChan
				tempChan <- temp
			}(j)
		}
		return tempChan
	}

	//create ticker for average period
	ticker := time.NewTicker(time.Duration(averagePeriod * float64(time.Second)))
	tempChan := gettingWeatherData()

	//infinite loop
	for {
		//use select case
		select {
		//averagePeriod timout case
		case <-ticker.C:
			//create local tempSum and count vars
			var tempSum float64 = 0.0
			count := 0

			//loop through channel at this point
			for len(tempChan) > 0 {
				tempSum += <-tempChan
				count++
			}

			//send average through out channel
			out <- tempSum / float64(count)

			//make new batch of rpc calls
			//--resets the channel to ignore old batches because gettingWeatherData recreates tempChan
			tempChan = gettingWeatherData()

		case <-quit:
			return //recieve quit message through quit channel and return
		}
	}
}
