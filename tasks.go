package tasks

import (
	"sync"
	"time"
)

var (
	data      datumCollection
	dataMutex sync.Mutex
)

const (
	iterateConcurrencyMax = 1 << 15
)

func iterate(lockDatumMutex bool, f func(*datum)) {
	var wg sync.WaitGroup
	var i int

	defer dataMutex.Unlock()
	dataMutex.Lock()

	for _, d := range data {
		wg.Add(1)

		go func(d *datum) {
			defer wg.Done()

			if lockDatumMutex {
				defer d.mutex.Unlock()
				d.mutex.Lock()
			}

			f(d)
		}(d)

		if i++; i%iterateConcurrencyMax == 0 {
			wg.Wait()
		}
	}

	if i%iterateConcurrencyMax != 0 {
		wg.Wait()
	}
}

func iterateWithConditionalRun() {
	iterate(false, func(d *datum) {
		var shouldRun bool
		var task Handler

		func() {
			defer d.mutex.Unlock()
			d.mutex.Lock()

			if !d.hasStatus(statusIsStopped) && !d.hasStatus(statusIsNowRunning) && time.Since(d.lastRunTime) >= d.runInterval {
				shouldRun = true
				task = d.task

				d.setStatus(statusIsNowRunning)
			}
		}()

		if shouldRun {
			go func() {
				lastRunTime := time.Now()

				defer func() {
					defer d.mutex.Unlock()
					d.mutex.Lock()

					d.unsetStatus(statusIsNowRunning)
					d.setStatus(statusHasRun)
					d.lastRunTime = lastRunTime
				}()

				task(d.identifier)
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

		iterate(true, func(d *datum) {
			if !d.hasStatus(statusIsStopped) {
				if i := d.sleepInterval; i < s {
					s = i
				}
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

	if runInitially {
		localDatum.lastRunTime = localDatum.setTime.Add(-interval)
	} else {
		localDatum.lastRunTime = localDatum.setTime
	}

	identifier := Identifier{&localDatum}
	localDatum.identifier = &identifier

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
