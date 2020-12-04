package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	cb "github.com/sksmith/cbatch"
)

func main() {
	results := cb.Process(
		handle,                     // your handler function
		os.Stdin,                   // the data you would like to process
		cb.Concurrency(20),         // how many records would you like to process simultaneously
		cb.Title("Testing Report"), // title of the output report
		cb.Progress)                // print progress to Stderr

	results.Print(os.Stdout)
}

// Define how you would like each row handled. Bare in mind that
// multiple handlers will be running simultaneously if you're using
// concurrency. So beware of sharing data between goroutines!
func handle(r []byte) error {
	// let's sleep a bit while processing records to simulate latency
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	// not all of the records in the "numbers" file are numbers
	_, err := strconv.Atoi(string(r))
	if err != nil {
		fmt.Printf("something interesting while processing value=[%v]\n", r)
		// errors end up in the results at the end of the process
		return err
	}
	return nil
}
