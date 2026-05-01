package domain

import (
	"sync"
	"time"

	"github.com/afuradanime/backend/internal/core/domain/value"
)

const ONLINE_TIMEOUT = 5
const ACTIVITY_TIMEOUT = 10

type UserActivity struct {
	status value.ActivityStatus
	timer  time.Time
}

type ActivityTracker struct {
	mu      sync.RWMutex
	entries map[int]UserActivity
}

func NewActivityTracker() *ActivityTracker {
	tracker := ActivityTracker{
		entries: make(map[int]UserActivity),
	}

	tracker.StartCleanup(1 * time.Minute) // Start cleanup routine
	return &tracker
}

func (a *ActivityTracker) RecordActivity(userID int, status value.ActivityStatus) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.entries[userID] = UserActivity{
		status: status,
		timer:  time.Now(),
	}
}

func (a *ActivityTracker) IsActive(userID int) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	last, ok := a.entries[userID]
	if !ok {
		return 0
	}
	return int(last.status)
}

func (a *ActivityTracker) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			a.mu.Lock()
			for userID, last := range a.entries {

				if time.Since(last.timer) >= (ACTIVITY_TIMEOUT+ONLINE_TIMEOUT)*time.Minute {
					delete(a.entries, userID)
				} else if time.Since(last.timer) >= ONLINE_TIMEOUT*time.Minute && last.status == value.Online {
					
					a.entries[userID] = UserActivity{
						status: value.Idle,
						timer:  last.timer,
					}
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
		if time.Since(last.timer) < ACTIVITY_TIMEOUT*time.Minute {
			users = append(users, userID)
		}
	}
	return users
}
