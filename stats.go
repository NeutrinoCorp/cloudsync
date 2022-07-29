package cloudsync

import "sync/atomic"

// Stats contains counters used by internal processes to keep track of its operations.
//
// This struct is goroutine-safe as it relies on atomic operations.
type Stats struct {
	currentUploadJobs uint64
	totalUploadJobs   uint64
	totalFailedJobs   uint64
}

var DefaultStats = &Stats{}

func (s *Stats) increaseUploadJobs() {
	atomic.AddUint64(&s.currentUploadJobs, 1)
	atomic.AddUint64(&s.totalUploadJobs, 1)
}

func (s *Stats) decreaseUploadJobs() {
	prev := atomic.LoadUint64(&s.currentUploadJobs)
	atomic.StoreUint64(&s.currentUploadJobs, prev-1)
}

func (s Stats) GetTotalUploadJobs() uint64 {
	return atomic.LoadUint64(&s.totalUploadJobs)
}

func (s Stats) GetCurrentUploadJobs() uint64 {
	return atomic.LoadUint64(&s.currentUploadJobs)
}

func (s *Stats) increaseFailedJobs() {
	atomic.AddUint64(&s.totalFailedJobs, 1)
}

func (s Stats) GetTotalFailedJobs() uint64 {
	return atomic.LoadUint64(&s.totalFailedJobs)
}
