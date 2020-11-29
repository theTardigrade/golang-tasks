package tasks

import "time"

// Identifier is used in the Add function.
// It exposes public methods to control the task.
type Identifier struct {
	*datum
}

// Stop ensures that the identified task no longer runs.
func (i *Identifier) Stop() {
	d := i.datum

	defer d.mutex.Unlock()
	d.mutex.Lock()

	d.isStopped = true
}

// Unstop ensures that the identified task runs as normal.
func (i *Identifier) Unstop() {
	d := i.datum

	defer d.mutex.Unlock()
	d.mutex.Lock()

	d.isStopped = false
}

// DurationSinceSet returns the duration of time since
// the task was first set.
func (i *Identifier) DurationSinceSet() time.Duration {
	return time.Since(i.datum.setTime)
}
