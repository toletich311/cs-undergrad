package main

import (
	"log"
	"math"
	"net/rpc"
	"strconv"
)

// TemperatureRequest represents an RPC request with a station ID
type TemperatureRequest struct {
	StationID string
}

// TemperatureResponse represents an RPC response with the temperature value
type TemperatureResponse struct {
	Temperature float64
}

// GetWeatherData fetches the temperature reading for a given weather station ID over RPC
func GetWeatherData(client *rpc.Client, id int) (float64, error) {

	// make instances of temp request and response
	//convert station id from int to string
	req := TemperatureRequest{StationID: strconv.Itoa(id)}
	var res TemperatureResponse

	//make rpc call from client, check for errors
	err := client.Call("WeatherService.GetTemperature", req, &res)
	if err != nil {
		return math.NaN(), nil
	}

	//return Temperature field of res instance of TemperatureResponse object
	return res.Temperature, nil
}

// Test GetWeatherData implementation
func main() {
	// Connect to the RPC server
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error connecting to RPC server:", err)
	}
	defer client.Close()

	// Fetch weather data for station ID 1
	temp, err := GetWeatherData(client, 1)
	if err != nil {
		log.Fatal("Error fetching weather data:", err)
	}

	log.Println("Temperature:", temp)
}
