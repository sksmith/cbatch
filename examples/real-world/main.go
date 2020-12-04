package main

import (
	"math/rand"
	"time"
)

func main() {
	// Get the data you would like to process
	orders := [][]interface{}{
		{"one"}, {"two"}, {"three"}, {"four"}, {"five"}, {"six"},
	}

	// Or uncomment this version if you would like some bells and whistles
	//
	// cb.Process(
	// 	handle,                     // your handler function
	// 	orders,                     // the data you would like to process
	// 	cb.Concurrency(2),          // how many records would you like to process simultaneously
	// 	cb.Title("Testing Report"), // title of the output report
	// 	cb.Report(os.Stdout),       // where would you like the report to write?
	// 	cb.Progress)                // prints progress to Stderr
}

// Define how you would like each row handled
func handle(r interface{}) error {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	return nil
}
