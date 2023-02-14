# golang-tasks

This package allows tasks to be set, whereby handler functions will run at given intervals.

[![Go Reference](https://pkg.go.dev/badge/github.com/theTardigrade/golang-tasks.svg)](https://pkg.go.dev/github.com/theTardigrade/golang-tasks) [![Go Report Card](https://goreportcard.com/badge/github.com/theTardigrade/golang-tasks)](https://goreportcard.com/report/github.com/theTardigrade/golang-tasks)

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

## Support

If you use this package, or find any value in it, please consider donating:

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/S6S2EIRL0)