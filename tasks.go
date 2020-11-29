package tasks

import (
	"sync"
	"time"
)

const (
	dataCapacityInitial = 1 << 12
)

var (
	data      = make(datumCollection, 0, dataCapacityInitial)
	dataMutex sync.Mutex
)

func iterate(f func(*datum)) {
	defer dataMutex.Unlock()
	dataMutex.Lock()

	for _, datum := range data {
		f(datum)
	}
}

func iterateWithConditionalRun() {
	iterate(func(d *datum) {
		var shouldRun bool
		var task Handler

		func() {
			defer d.mutex.Unlock()
			d.mutex.Lock()

			if !d.isNowRunning && !d.isStopped && (d.lastRunTime.IsZero() || time.Since(d.lastRunTime) >= d.runInterval) {
				shouldRun = true
				task = d.task

				d.isNowRunning = true
			}
		}()

		if shouldRun {
			go func() {
				defer func() {
					defer d.mutex.Unlock()
					d.mutex.Lock()

					d.lastRunTime = time.Now()
					d.isNowRunning = false
				}()

				task(&Identifier{d})
			}()
		}
	})
}

var (
	sleepDurationCached   *time.Duration
	sleepChanCancellation = make(chan struct{}, 1)
)

const (
	// SleepDurationMin is the minimum duration that the internal loop will pause before checking tasks.
	SleepDurationMin = time.Millisecond * 2
	// SleepDurationMax is the maximum duration that the internal loop will pause before checking tasks.
	SleepDurationMax = time.Hour * 24 * 100
)

func sleepDuration() (s time.Duration) {
	if sleepDurationCached != nil {
		s = *sleepDurationCached
	} else {
		s = SleepDurationMax

		iterate(func(d *datum) {
			defer d.mutex.Unlock()
			d.mutex.Lock()

			if i := d.sleepInterval; i < s {
				s = i
			}
		})

		if s < SleepDurationMin {
			s = SleepDurationMin
		}

		sleepDurationCached = &s
	}

	return
}

func sleep() {
	duration := sleepDuration()

	select {
	case <-time.After(duration):
	case <-sleepChanCancellation:
		sleepDurationCached = nil
	}
}

func sleepCancel() {
	select {
	case sleepChanCancellation <- struct{}{}:
	default:
		// sleepChanCancellation is already full
	}
}

const (
	sleepIntervalDivisor = 100
)

// Set ensures a callback function will run continuously at a given interval.
func Set(interval time.Duration, runInitially bool, task Handler) *Identifier {
	localDatum := datum{
		task:          task,
		runInterval:   interval,
		sleepInterval: interval / sleepIntervalDivisor,
		setTime:       time.Now(),
	}

	if !runInitially {
		localDatum.lastRunTime = localDatum.setTime
	}

	identifier := Identifier{&localDatum}

	defer dataMutex.Unlock()
	dataMutex.Lock()

	data = append(data, &localDatum)

	sleepCancel()

	return &identifier
}

func run() {
	for {
		iterateWithConditionalRun()
		sleep()
	}
}

func init() {
	go run()
}