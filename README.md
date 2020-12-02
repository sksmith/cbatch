# Concurrent Batch (CBatch)

This utility package is useful for creating quick batch jobs that execute concurrently.  
Call `cbatch.Process` in your program and supplying the handler function and required data.

## Example

Execute the below example like so...

```shell
cd example/simple
go run .
```

If you uncomment the version with the bells and whistles, you'll want to send your output
to a log file...

```shell
cd example/simple
go run . > results.log
```

```golang
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
```

```golang
// Or speed things up with some concurrency
r = strings.NewReader(data)
results := cb.Process(handle, r, cb.Concurrency(2))
if len(results.Errors) > 0 {
  // how would you like to handle the errors?
}
```

```golang
// Or if you would like some bells and whistles
r = strings.NewReader(data)
results = cb.Process(handle, r,
  cb.Concurrency(2),
  cb.Title("Some Batch Job"), // title of the output report
  cb.Progress)       // prints a progress bar to stderr
```