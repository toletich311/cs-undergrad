// Package future provides a Future type that can be used to
// represent a value that will be available in the future.
package future

import (
	"net/rpc"
	"strconv"
	"time"
)

// WeatherDataResult be used in GetWeatherData.
type WeatherDataResult struct {
	Value interface{}
	Err   error
}

// TemperatureRequest represents an RPC request with a station ID
type TemperatureRequest struct {
	StationID string
}

// TemperatureResponse represents an RPC response with the temperature value
type TemperatureResponse struct {
	Temperature float64
}

type Future struct {
	result chan interface{}
}

func NewFuture() *Future {
	return &Future{
		result: make(chan interface{}, 1),
	}
}

func (f *Future) CompleteFuture(res interface{}) {
	f.result <- res
	f.CloseFuture()
}

func (f *Future) GetResult() interface{} {
	return <-f.result
}

func (f *Future) CloseFuture() {
	close(f.result)
}

// Wait waits for the first n futures to return or for the timeout to expire,
// whichever happens first.
func Wait(futures []*Future, n int, timeout time.Duration, filter func(interface{}) bool) []interface{} {
	// STEP 1 - INTIALIZE
	// Initialize the results slice with a capacity of n
	// Create a buffered channel `data` to store results from all futures
	results := make([]interface{}, 0, n)
	data := make(chan interface{}, len(futures))

	//  then create timer and count to use in switch case to timeout
	timer := time.After(timeout)
	count := 0 //tracking responses, not to exceed n

	// STEP 2 - SEND GO ROUTINES
	//start for loop - use length of "futures"
	for _, future := range futures {
		//start goroutines for each future, retrieve result then sent to "data" channel
		go func(f *Future) {
			res := f.GetResult() //use "res", temp local var
			data <- res
		}(future)
	}

	// STEP 3 - COLLECT RESPONSES OR TIMEOUT
	//create for loop , only run n times
	for count < n {
		//start switch case
		select {
		case res := <-data:
			// case 1, recieve data before timeout

			//only count filtered results in case of HeatWaveWarning
			if filter == nil || filter(res) {
				results = append(results, res)
				count++ //increment count
			}

		case <-timer:
			// case 2, timeout and return whatever results recieved so far
			return results
		}
	}

	//case 3, didn't time out but reached cap "n" results
	return results
}

// User Defined Function Logic

// GetWeatherData implementation which immediately returns a Future.
func GetWeatherData(client *rpc.Client, id int) *Future {
	// STEP 1 - create future
	future := NewFuture()

	//STEP 2 - spawn goroutines
	go func() {
		//make instances of temp request and declare a response var
		request := TemperatureRequest{StationID: strconv.Itoa(id)}
		var response TemperatureResponse

		//make rpc call from client, check for errors
		err := client.Call("WeatherService.GetTemperature", request, &response)

		//error check in case of rpc failure
		if err != nil {
			//in case of error don't pass value
			future.CompleteFuture(WeatherDataResult{Value: nil, Err: err})
		} else {
			//else pass value
			future.CompleteFuture(WeatherDataResult{Value: response.Temperature, Err: nil})
		}
	}()
	//return response through future
	return future
}

// heatWaveWarning is the filter function for the received weatherData.
// Should be used to keep only temperatures > 35 degrees Celsius (return false otherwise).
func heatWaveWarning(res interface{}) bool {

	//get data from WeatherDataResult
	data, warning := res.(WeatherDataResult)

	//error check -- return false in case of rpc or other failure
	if !warning || data.Err != nil {
		return false
	}

	//create temp var for temperature
	temp, warning := data.Value.(float64)

	//return true if >=35, false otherwise
	if temp >= 35 && warning {
		return true //occurs when temp is <= 35 celsius
	} else {
		return false //occurs otherwise
	}
}
