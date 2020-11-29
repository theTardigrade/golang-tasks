package tasks

import (
	"sync"
	"time"
)

type datum struct {
	task          Handler
	runInterval   time.Duration
	sleepInterval time.Duration
	lastRunTime   time.Time
	isNowRunning  bool
	isStopped     bool
	mutex         sync.Mutex
}

type datumCollection []*datum
