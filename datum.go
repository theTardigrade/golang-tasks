package tasks

import (
	"sync"
	"time"
)

type datum struct {
	task          Handler
	identifier    *Identifier
	runInterval   time.Duration
	sleepInterval time.Duration
	setTime       time.Time
	lastRunTime   time.Time
	status        status
	mutex         sync.Mutex
}

type datumCollection []*datum
