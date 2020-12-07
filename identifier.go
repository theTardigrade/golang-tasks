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

	d.setStatus(statusIsStopped)

	sleepCancel()
}

// Unstop ensures that the identified task runs as normal.
func (i *Identifier) Unstop() {
	d := i.datum

	defer d.mutex.Unlock()
	d.mutex.Lock()

	d.unsetStatus(statusIsStopped)

	sleepCancel()
}

// ChangeInterval allows the interval between runs of the
// task to be modified dynamically.
func (i *Identifier) ChangeInterval(interval time.Duration) {
	d := i.datum

	defer d.mutex.Unlock()
	d.mutex.Lock()

	d.runInterval = interval
	d.sleepInterval = interval / sleepIntervalDivisor

	sleepCancel()
}

// DurationSinceSet returns the duration of time since
// the task was first set.
func (i *Identifier) DurationSinceSet() time.Duration {
	d := i.datum

	t := func() time.Time {
		defer d.mutex.Unlock()
		d.mutex.Lock()

		return d.setTime
	}()

	return time.Since(t)
}

// DurationSinceLastRun returns a bool value indicating if
// the task has ever run and, if so, the duration of time since
// it last did so.
func (i *Identifier) DurationSinceLastRun() (hasRun bool, duration time.Duration) {
	d := i.datum

	var t time.Time
	func() {
		defer d.mutex.Unlock()
		d.mutex.Lock()

		if d.hasStatus(statusHasRun) {
			hasRun = true
			t = d.lastRunTime
		}
	}()

	if hasRun {
		duration = time.Since(t)
	}

	return
}
