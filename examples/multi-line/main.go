package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	cb "github.com/sksmith/cbatch"
)

// Sometimes you'll want to read more than a single line at a time per thread.
// cbatch offers an implementation of bufio.SplitFunc that allows consuming of
// multiple lines at a time.
func main() {
	results := cb.Process(
		handle,
		os.Stdin,
		cb.Concurrency(3),
		cb.Title("Multiline Execution"),
		cb.Progress,
		cb.Split(cb.ScanMultiLines(8)))

	results.Print(os.Stdout)
}

// Define how you would like each row handled. Bare in mind that
// multiple handlers will be running simultaneously if you're using
// concurrency. So beware of sharing data between goroutines!
func handle(r []byte) error {
	// let's sleep a bit while processing records to simulate latency
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	s := string(r)
	v := strings.Split(s, "\n")
	p := ""
	for _, c := range v {
		p += c + "-"
	}
	fmt.Printf("Processed: %s\n", p)
	return nil
}
