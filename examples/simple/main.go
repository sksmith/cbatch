package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	cb "github.com/sksmith/cbatch"
)

func main() {
	// Get the data you would like to process
	data := []interface{}{
		"one", "two", "three", "four", "five", "six",
	}

	// Execute handler with the bare minimum
	cb.Process(handle, data)

	// Or speed things up with some concurrency
	results := cb.Process(handle, data, cb.Concurrency(2))
	if len(results.Errors) > 0 {
		// how would you like to handle the errors?
	}

	// Or if you would like some bells and whistles
	results = cb.Process(
		handle,                     // your handler function
		data,                       // the data you would like to process
		cb.Concurrency(2),          // how many records would you like to process simultaneously
		cb.Title("Some Batch Job"), // title of the output report
		cb.Progress)                // prints a progress bar to stderr

	results.Print(os.Stdout)
}

// Define how you would like each row handled
func handle(s interface{}) error {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	if s.(string) == "two" {
		return fmt.Errorf("failed to process %s", s)
	}
	return nil
}
