package tasks

type status uint8

const (
	statusHasRun = 1 << iota
	statusIsNowRunning
	statusIsStopped
)

func (d *datum) setStatus(s status) {
	d.status |= s
}

func (d *datum) unsetStatus(s status) {
	d.status &= ^s
}

func (d *datum) hasStatus(s status) bool {
	return (d.status & s) != 0
}
