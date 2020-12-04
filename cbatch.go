// Package cbatch provides the Process function which allows easy processing
// of a set of data in a concurrent fashion.
package cbatch

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type runOptions struct {
	title        string
	headers      map[string]string
	concurrency  int
	report       bool
	reportWriter io.Writer
	progress     bool
}

// Option represents a given option for execution.
// Any given option can only be specified once unless
// otherwise stated in the documentation.
type Option func(o *runOptions)

// Title sets the title that prints at the top of the log
func Title(t string) func(o *runOptions) {
	return func(o *runOptions) { o.title = t }
}

// Header adds custom headers printed at the top of the log.
// This option can be used more than once.
func Header(k, v string) func(o *runOptions) {
	return func(o *runOptions) { o.headers[k] = v }
}

// Concurrency sets the maximum number of simultaneous processes
// being executed simultaneously.
func Concurrency(c int) func(o *runOptions) {
	return func(o *runOptions) { o.concurrency = c }
}

// Report prints simple ASCII report to the requested stream.
func Report(w io.Writer) func(o *runOptions) {
	return func(o *runOptions) {
		o.report = true
		o.reportWriter = w
	}
}

// Progress sends the progress of the process to `os.Stderr`
func Progress(o *runOptions) {
	o.progress = true
}

// Process takes a set of data and calls the exec function once for each entry
// in the parent array. Passing the child arrays as input. A few options can
// be provided for modifying concurrency, or outputting the results.
func Process(exec func(interface{}) error, data []interface{}, options ...Option) {
	start := time.Now()

	o := runOptions{
		title: "",
		headers: map[string]string{
			"Executed": start.String(),
		},
		concurrency: 1,
	}
	for _, option := range options {
		option(&o)
	}
	o.headers["Concurrency"] = strconv.Itoa(o.concurrency)
	printHeaders(o)
	print(o, "\n")

	failures := make([]error, 0)

	sem := make(chan bool, o.concurrency)
	success := make(chan bool)
	fail := make(chan error)
	quit := make(chan bool)

	recordCount := len(data)

	printLogHeader(o)

	go func() {
		done := 0

		var bar progressBar
		if o.progress {
			bar.New(0, int64(recordCount))
		}

		finished := false
		for !finished {
			select {
			case f := <-fail:
				failures = append(failures, f)
				done++
			case <-success:
				done++
			case <-quit:
				finished = true
			}

			if o.progress {
				bar.Play(int64(done))
				if finished {
					bar.Finish()
				}
			}
		}
	}()

	for _, record := range data {
		sem <- true
		go func(r interface{}) {
			defer func() { <-sem }()
			err := exec(r)
			if err != nil {
				fail <- err
				return
			}
			success <- true
		}(record)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	quit <- true

	print(o, "\n")

	if len(failures) > 0 {
		printFailures(o, failures)
		print(o, "\n")
	}

	printFooter(o, recordCount, failures, start)
}

func printHeaders(o runOptions) {
	maxLen := maxLabel(o.headers)

	if o.title != "" {
		bars := strings.Repeat("-", len(o.title)) + "\n"
		print(o, "  "+bars)
		print(o, "  %s\n", o.title)
		print(o, "  "+bars)
	}

	for k, v := range o.headers {
		print(o, " %*s: %s\n", maxLen, k, v)
	}
}

func maxLabel(labels map[string]string) int {
	maxLen := 0
	for k := range labels {
		if maxLen < len(k) {
			maxLen = len(k)
		}
	}
	return maxLen
}

func print(o runOptions, msg string, a ...interface{}) {
	if o.report {
		fmt.Fprintf(o.reportWriter, msg, a...)
	}
}

func printLogHeader(o runOptions) {
	print(o, "  --------\n")
	print(o, "    Logs\n")
	print(o, "  --------\n")
}

func printFailures(o runOptions, failures []error) {
	print(o, "  --------\n")
	print(o, "  Failures\n")
	print(o, "  --------\n")
	for _, failure := range failures {
		print(o, failure.Error()+"\n")
	}
}

func printFooter(o runOptions, recordCount int, failures []error, start time.Time) {
	print(o, "  --------\n")
	print(o, "  Results\n")
	print(o, "  --------\n")
	print(o, "  Finished: %v\n", time.Now())
	print(o, " Processed: %d\n", recordCount)
	print(o, "    Failed: %d\n", len(failures))
	print(o, "   Elapsed: %ds\n", time.Since(start)/time.Second)
}
