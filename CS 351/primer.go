// Package primer provides practice exercises to get you started with Go programming.
package primer

import (
	"bufio"
	"os"
	"strconv"
	"sync"
	"time"
)

// SumFromFile should read a file with the given filename, and
// return the sum of all line-separated integers in the file.
func SumFromFile(filename string) (int, error) {
	// TODO: Your code here

	//open the file and throw exception
	file, ferr := os.Open(filename)

	if ferr != nil {
		panic(ferr)
	}

	//create scanner and sum variable
	scanner := bufio.NewScanner(file)
	var sum int = 0

	//iterate through lines in file
	for scanner.Scan() {
		line := scanner.Text()
		//convert each line to an int, throwing exceptions when needed
		num, nerr := strconv.Atoi(line)
		if nerr != nil {
			panic(nerr)
		}
		//fmt.Println(sum) debugging
		sum += num //add to sum count
		//fmt.Println(num) debugging

	}

	//return sum count
	return sum, ferr
}

// RunConcurrently should run the given function n times concurrently.
func RunConcurrently(n int, f func()) {
	// TODO: Your code here

	//call goroutine n times with for loop
	//commented out all testing lines as well

	//testing
	//startNow := time.Now()

	for i := 0; i < n; i++ {
		go f()
		//testing
		//fmt.Printf(" call %d ", i)
	}

	//testing
	//fmt.Println("This took: ", time.Since(startNow))
}

// SafeIncrement should increment the given counter safely using a mutex.
func SafeIncrement(counter *int, mu *sync.Mutex) {
	// TODO: Your code here
	//lock mutex, increment counter, then unlock mutex
	mu.Lock()
	*counter++
	mu.Unlock()

}

// Ticker returns a channel that sends a signal every `d` duration.
func Ticker(d time.Duration) chan struct{} {
	// TODO: Your code here

	//make ticker chan
	ticker := make(chan struct{})

	//goroutine
	go func() {
		for {
			//wait for duration then call
			time.Sleep(d)
			ticker <- struct{}{}
		}
	}()

	//return ticker chan
	return ticker
}
