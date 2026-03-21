package domain

import (
	"sync"
	"time"
)

const ACTIVITY_TIMEOUT = 5

type ActivityTracker struct {
	mu      sync.RWMutex
	entries map[int]time.Time
}

func NewActivityTracker() *ActivityTracker {
	tracker := ActivityTracker{
		entries: make(map[int]time.Time),
	}

	tracker.StartCleanup(ACTIVITY_TIMEOUT * time.Minute) // Start cleanup routine
	return &tracker
}

func (a *ActivityTracker) RecordActivity(userID int) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.entries[userID] = time.Now()
}

func (a *ActivityTracker) IsActive(userID int) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	last, ok := a.entries[userID]
	if !ok {
		return false
	}
	return time.Since(last) < ACTIVITY_TIMEOUT*time.Minute
}

func (a *ActivityTracker) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			a.mu.Lock()
			for userID, last := range a.entries {
				if time.Since(last) >= ACTIVITY_TIMEOUT*time.Minute {
					delete(a.entries, userID)
				}
			}
			a.mu.Unlock()
		}
	}()
}

func (a *ActivityTracker) GetActiveUsers() []int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	users := make([]int, 0, len(a.entries))
	for userID, last := range a.entries {
		if time.Since(last) < ACTIVITY_TIMEOUT*time.Minute {
			users = append(users, userID)
		}
	}
	return users
}
