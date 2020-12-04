package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	cb "github.com/sksmith/cbatch"
)

func main() {
	data := `one
		two
		three
		four
		five`

	// Execute handler with the bare minimum
	r := strings.NewReader(data)
	cb.Process(handle, r)

	// Or speed things up with some concurrency
	r = strings.NewReader(data)
	results := cb.Process(handle, r, cb.Concurrency(2))
	if len(results.Errors) > 0 {
		// how would you like to handle the errors?
	}

	// Or if you would like some bells and whistles
	r = strings.NewReader(data)
	results = cb.Process(
		handle,                     // your handler function
		r,                          // the data you would like to process
		cb.Concurrency(2),          // how many records would you like to process simultaneously
		cb.Title("Some Batch Job"), // title of the output report
		cb.Progress)                // prints a progress bar to stderr

	results.Print(os.Stdout)
}

// Define how you would like each row handled
func handle(v []byte) error {
	s := strings.TrimSpace(string(v))
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	if s == "two" {
		return fmt.Errorf("failed to process %s", s)
	}
	return nil
}
