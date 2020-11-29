package tasks

// Identifier is used in the Add function; it allows tasks to be stopped and unstopped.
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
