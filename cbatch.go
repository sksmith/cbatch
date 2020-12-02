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

// Progress sends the progress of the process to `os.Stderr`
func Progress(o *runOptions) {
	o.progress = true
}

type Results struct {
	Title       string
	StartedAt   time.Time
	FinishedAt  time.Time
	Concurrency int
	Headers     map[string]string
	Errors      []error
	RecordCount int
}

// Process takes a set of data and calls the exec function once for each entry
// in the parent array. Passing the child arrays as input. A few options can
// be provided for modifying concurrency, or outputting the results.
func Process(exec func(interface{}) error, data []interface{}, options ...Option) Results {
	results := Results{}
	results.StartedAt = time.Now()

	o := runOptions{
		title: "",
		headers: map[string]string{
			"Executed": results.StartedAt.String(),
		},
		concurrency: 1,
	}
	for _, option := range options {
		option(&o)
	}
	o.headers["Concurrency"] = strconv.Itoa(o.concurrency)

	results.Title = o.title
	results.Concurrency = o.concurrency
	results.Headers = o.headers
	results.Errors = make([]error, 0)

	sem := make(chan bool, o.concurrency)
	success := make(chan bool)
	fail := make(chan error)
	quit := make(chan bool)

	results.RecordCount = len(data)

	go func() {
		done := 0

		var bar progressBar
		if o.progress {
			bar.New(0, int64(results.RecordCount))
		}

		finished := false
		for !finished {
			select {
			case f := <-fail:
				results.Errors = append(results.Errors, f)
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

	results.FinishedAt = time.Now()

	return results
}

func (r *Results) Print(w io.Writer) {
	r.printHeaders(w)
	fmt.Fprintln(w)
	r.printErrors(w)
	fmt.Fprintln(w)
	r.printFooter(w)
}

func (r *Results) printHeaders(w io.Writer) {
	maxLen := maxLabel(r.Headers)

	if r.Title != "" {
		bars := strings.Repeat("-", len(r.Title)) + "\n"
		print(w, "  "+bars)
		print(w, "  %s\n", r.Title)
		print(w, "  "+bars)
	}

	for k, v := range r.Headers {
		print(w, " %*s: %s\n", maxLen, k, v)
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

func print(w io.Writer, msg string, a ...interface{}) {
	fmt.Fprintf(w, msg, a...)
}

func (r *Results) printErrors(w io.Writer) {
	print(w, "  --------\n")
	print(w, "   Errors\n")
	print(w, "  --------\n")
	for _, err := range r.Errors {
		print(w, err.Error()+"\n")
	}
}

func (r *Results) printFooter(w io.Writer) {
	elapsed := r.FinishedAt.Sub(r.StartedAt).Milliseconds()
	average := elapsed / int64(r.RecordCount)

	print(w, "  --------\n")
	print(w, "  Results\n")
	print(w, "  --------\n")
	print(w, "  Finished: %v\n", r.FinishedAt)
	print(w, " Processed: %d\n", r.RecordCount)
	print(w, "    Failed: %d\n", len(r.Errors))
	print(w, "   Elapsed: %dms\n", elapsed)
	print(w, "   Average: %dms\n", average)
}
