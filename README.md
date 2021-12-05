# golang-tasks

This package allows tasks to be set, whereby handler functions will run at given intervals.

## Example

```golang
package main

import (
	"fmt"
	"time"

	tasks "github.com/theTardigrade/golang-tasks"
)

func main() {
	// set up a handler function to run once every second;
	// do not call the function on initialization
	tasks.Set(time.Second, false, func(id *tasks.Identifier) {
		fmt.Println("ONE SECOND HAS PASSED")

		// stop the task after five seconds
		if id.DurationSinceSet() >= time.Second*5 {
			id.Stop()
		}
	})

	// keep the main function running
	select {}
}
```