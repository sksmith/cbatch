# Concurrent Batch (CBatch)

This utility package is useful for creating quick batch jobs that execute concurrently.  
Call `cbatch.Process` in your program and supplying the handler function and required data.

## Example

Execute the below example like so...

```shell
cd example
go run .
```

If you uncomment the version with the bells and whistles, you'll want to send your output
to a log file...

```shell
cd example
go run . > results.log
```

```golang
package main

import (
  "math/rand"
  "time"

  cb "github.com/sksmith/cbatch"
)

func main() {
  // Get the data you would like to process
  orders := [][]interface{}{
    {"one"}, {"two"}, {"three"}, {"four"}, {"five"}, {"six"},
  }

  // Execute handler with the bare minimum
  cb.Process(handle, orders)

  // Speed things up with some concurrency
  cb.Process(handle, orders, cb.Concurrency(2))

  // Or uncomment this version if you would like some bells and whistles
  //
  // cb.Process(
  //  handle,                     // your handler function
  //  orders,                     // the data you would like to process
  //  cb.Concurrency(2),          // how many records would you like to process simultaneously
  //  cb.Title("Testing Report"), // title of the output report
  //  cb.Report(os.Stdout),       // where would you like the report to write?
  //  cb.Progress)                // prints progress to Stderr
}

// Define how you would like each row handled
func handle(s interface{}) error {
  time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
  return nil
}
```

### A More "Real World" Example

